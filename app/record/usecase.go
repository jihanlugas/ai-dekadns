package record

import (
	"ai-dekadns/app/organization"
	"ai-dekadns/app/project"
	"ai-dekadns/app/superadminrole"
	"ai-dekadns/app/user"
	"ai-dekadns/app/zone"
	"ai-dekadns/constant"
	"ai-dekadns/helper"
	"ai-dekadns/request"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joeig/go-powerdns/v3"
)

type Usecase interface {
	Create(c *gin.Context, req request.CreateRecord) (err error)
	Update(c *gin.Context, req request.UpdateRecord) (err error)
	Delete(c *gin.Context, req request.DeleteRecord) (err error)
}

type usecase struct {
	organizationRepo   organization.Repository
	projectRepo        project.Repository
	userRepo           user.Repository
	superadminroleRepo superadminrole.Repository
	zoneRepo           zone.Repository
}

func (u usecase) Create(c *gin.Context, req request.CreateRecord) (err error) {
	var userId = helper.GetUserID(c)
	var role = helper.GetRole(c)

	tZone, err := u.zoneRepo.GetById(req.ZoneID)
	if err != nil {
		return err
	}

	project, err := u.projectRepo.GetById(tZone.ProjectId)
	if err != nil {
		return err
	}

	user, err := u.userRepo.GetById(userId, "Organization", "OrganizationRole")
	if err != nil {
		return err
	}

	if role != constant.RoleSuperadmin {
		// check IDOR
		if user.OrganizationID != project.OrganizationID {
			return errors.New("you are not allowed to perform this operation")
		}

		// check user role permission
		// When organization role is not set, the user have no privileges.
		if user.OrganizationRole == nil {
			return errors.New("you are not allowed to perform this operation")
		}

		ok := user.OrganizationRole.CheckForPrivilege(constant.PrivilegeDns, true)
		if !ok {
			return errors.New("you are not allowed to perform this operation")
		}
	} else {
		rolePrivileges, err := u.superadminroleRepo.GetBySuperadminId(user.ID)
		if err != nil {
			return err
		}

		isAllow := false
		for _, role := range rolePrivileges {
			if role.CheckForPrivilege(constant.SuperadminPrivilegeDns, true) {
				isAllow = true
				break
			}
		}
		if !isAllow {
			return errors.New("you are not allowed to perform this operation")
		}
	}

	pdns := powerdns.New(os.Getenv("DNS_HOST"), "localhost", powerdns.WithHeaders(map[string]string{
		"X-API-Key": os.Getenv("DNS_API_KEY_VALUE"),
	}))

	pdnsZone, err := pdns.Zones.Get(c, tZone.Name)
	if err != nil {
		go helper.ErrorToAuditLog(userId, fmt.Sprintf("Create Record %s", tZone.Name), "DNS", "Create Record", err.Error(), project.ID, project.OrganizationID, helper.GetIP(c))
		return err
	}

	var existingContents []string
	for _, rrset := range pdnsZone.RRsets {
		reqType := powerdns.RRType(req.Type)
		reqName := fmt.Sprintf("%s.", req.Name)
		if *rrset.Name == reqName && *rrset.Type == reqType {
			for _, rec := range rrset.Records {
				existingContents = append(existingContents, *rec.Content)
			}
			break
		}
	}

	existingContents = append(existingContents, helper.ToPowerDNS(powerdns.RRType(req.Type), req.Content))

	err = pdns.Records.Add(c, tZone.Name, req.Name, powerdns.RRType(req.Type), req.TTL, existingContents)
	if err != nil {
		go helper.ErrorToAuditLog(userId, fmt.Sprintf("Create Record %s", tZone.Name), "DNS", "Create Record", err.Error(), project.ID, project.OrganizationID, helper.GetIP(c))
		return fmt.Errorf("failed to add record: %w", err)
	}

	go helper.InfoToAuditLog(userId, fmt.Sprintf("Create Record %s", tZone.Name), "DNS", "Create Record", "Success Create Record", project.ID, project.OrganizationID, helper.GetIP(c))

	return err
}

func (u usecase) Update(c *gin.Context, req request.UpdateRecord) error {
	var userId = helper.GetUserID(c)
	var role = helper.GetRole(c)

	// --- Authorization & permission checks (tetap seperti punyamu) ---
	tZone, err := u.zoneRepo.GetById(req.ZoneID)
	if err != nil {
		return err
	}

	project, err := u.projectRepo.GetById(tZone.ProjectId)
	if err != nil {
		return err
	}

	user, err := u.userRepo.GetById(userId, "Organization", "OrganizationRole")
	if err != nil {
		return err
	}

	if role != constant.RoleSuperadmin {
		if user.OrganizationID != project.OrganizationID {
			return errors.New("you are not allowed to perform this operation")
		}
		if user.OrganizationRole == nil {
			return errors.New("you are not allowed to perform this operation")
		}
		ok := user.OrganizationRole.CheckForPrivilege(constant.PrivilegeDns, true)
		if !ok {
			return errors.New("you are not allowed to perform this operation")
		}
	} else {
		rolePrivileges, err := u.superadminroleRepo.GetBySuperadminId(user.ID)
		if err != nil {
			return err
		}
		isAllow := false
		for _, role := range rolePrivileges {
			if role.CheckForPrivilege(constant.SuperadminPrivilegeDns, true) {
				isAllow = true
				break
			}
		}
		if !isAllow {
			return errors.New("you are not allowed to perform this operation")
		}
	}

	// --- PDNS client v3 ---
	pdns := powerdns.New(os.Getenv("DNS_HOST"), "localhost", powerdns.WithAPIKey(os.Getenv("DNS_API_KEY_VALUE")))

	// Pastikan zone reachable
	pdnsZone, err := pdns.Zones.Get(c, tZone.Name)
	if err != nil {
		go helper.ErrorToAuditLog(userId, fmt.Sprintf("Update Record %s", req.OldName), "DNS", "Update Record", err.Error(), project.ID, project.OrganizationID, helper.GetIP(c))
		return fmt.Errorf("zone get failed: %w", err)
	}

	oldName := helper.EnsureDot(req.OldName)
	newName := helper.EnsureDot(req.Name)
	oldType := powerdns.RRType(strings.ToUpper(req.OldType))
	newType := powerdns.RRType(strings.ToUpper(req.Type))

	oldContent := helper.ToPowerDNS(powerdns.RRType(req.OldType), req.OldContent)
	newContent := helper.ToPowerDNS(powerdns.RRType(req.Type), req.Content)

	if oldType == newType && oldName == newName {
		contents := []string{}
		for _, rrset := range pdnsZone.RRsets {
			if *rrset.Name == oldName && *rrset.Type == oldType {
				for _, record := range rrset.Records {
					if *record.Content != oldContent {
						contents = append(contents, *record.Content)
					}
				}
			}
		}
		contents = append(contents, newContent)
		err = pdns.Records.Change(c, tZone.Name, newName, newType, req.TTL, contents)
		if err != nil {
			go helper.ErrorToAuditLog(userId, fmt.Sprintf("Update Record %s", req.OldName), "DNS", "Update Record", err.Error(), project.ID, project.OrganizationID, helper.GetIP(c))
			return fmt.Errorf("failed to change records: %w", err)
		}
	} else {
		// cari data yang lama lalu update
		contents := []string{}
		for _, rrset := range pdnsZone.RRsets {
			if *rrset.Name == oldName && *rrset.Type == oldType {
				for _, record := range rrset.Records {
					if *record.Content != oldContent {
						contents = append(contents, *record.Content)
					}
				}
			}
		}
		err = pdns.Records.Change(c, tZone.Name, oldName, oldType, req.OldTTL, contents)
		if err != nil {
			go helper.ErrorToAuditLog(userId, fmt.Sprintf("Update Record %s", req.OldName), "DNS", "Update Record", err.Error(), project.ID, project.OrganizationID, helper.GetIP(c))
			return fmt.Errorf("failed to change records: %w", err)
		}

		// add data yang baru
		contents = []string{}
		for _, rrset := range pdnsZone.RRsets {
			if *rrset.Name == newName && *rrset.Type == newType {
				for _, record := range rrset.Records {
					contents = append(contents, *record.Content)
				}
			}
		}
		contents = append(contents, newContent)
		err = pdns.Records.Change(c, tZone.Name, newName, newType, req.TTL, contents)
		if err != nil {
			go helper.ErrorToAuditLog(userId, fmt.Sprintf("Update Record %s", req.OldName), "DNS", "Update Record", err.Error(), project.ID, project.OrganizationID, helper.GetIP(c))
			return fmt.Errorf("failed to change records: %w", err)
		}
	}

	go helper.InfoToAuditLog(userId, fmt.Sprintf("Update Record %s", req.OldName), "DNS", "Update Record", "Success Update Record", project.ID, project.OrganizationID, helper.GetIP(c))

	return nil
}

func (u usecase) Delete(c *gin.Context, req request.DeleteRecord) (err error) {
	var userId = helper.GetUserID(c)
	var role = helper.GetRole(c)

	tZone, err := u.zoneRepo.GetById(req.ZoneID)
	if err != nil {
		return err
	}

	project, err := u.projectRepo.GetById(tZone.ProjectId)
	if err != nil {
		return err
	}

	user, err := u.userRepo.GetById(userId, "Organization", "OrganizationRole")
	if err != nil {
		return err
	}

	if role != constant.RoleSuperadmin {
		// check IDOR
		if user.OrganizationID != project.OrganizationID {
			return errors.New("you are not allowed to perform this operation")
		}

		// check user role permission
		// When organization role is not set, the user have no privileges.
		if user.OrganizationRole == nil {
			return errors.New("you are not allowed to perform this operation")
		}

		ok := user.OrganizationRole.CheckForPrivilege(constant.PrivilegeDns, true)
		if !ok {
			return errors.New("you are not allowed to perform this operation")
		}
	} else {
		rolePrivileges, err := u.superadminroleRepo.GetBySuperadminId(user.ID)
		if err != nil {
			return err
		}

		isAllow := false
		for _, role := range rolePrivileges {
			if role.CheckForPrivilege(constant.SuperadminPrivilegeDns, true) {
				isAllow = true
				break
			}
		}
		if !isAllow {
			return errors.New("you are not allowed to perform this operation")
		}
	}

	pdns := powerdns.New(os.Getenv("DNS_HOST"), "localhost", powerdns.WithHeaders(map[string]string{
		"X-API-Key": os.Getenv("DNS_API_KEY_VALUE"),
	}))

	pdnsZone, err := pdns.Zones.Get(c, tZone.Name)
	if err != nil {
		go helper.ErrorToAuditLog(userId, fmt.Sprintf("Delete Record %s", req.Name), "DNS", "Delete Record", err.Error(), project.ID, project.OrganizationID, helper.GetIP(c))
		return err
	}

	var newContents []string
	ttl := &req.TTL
	found := false

	reqName := helper.EnsureDot(req.Name)
	reqType := powerdns.RRType(req.Type)
	reqContent := helper.ToPowerDNS(reqType, req.Content)

	// Cari RRSet sesuai name + type
	for _, rrset := range pdnsZone.RRsets {

		if *rrset.Name == reqName && *rrset.Type == reqType {
			found = true
			ttl = rrset.TTL

			for _, rec := range rrset.Records {
				if *rec.Content != reqContent {
					newContents = append(newContents, *rec.Content)
				}
			}
			break
		}
	}

	if !found {
		err = fmt.Errorf("record not found: %s %s", req.Name, req.Type)
		go helper.ErrorToAuditLog(userId, fmt.Sprintf("Delete Record %s", req.Name), "DNS", "Delete Record", err.Error(), project.ID, project.OrganizationID, helper.GetIP(c))
		return err
	}

	// Kalau masih ada content lain → update RRSet
	if len(newContents) > 0 {
		err = pdns.Records.Add(c, tZone.Name, req.Name, powerdns.RRType(req.Type), *ttl, newContents)
		if err != nil {
			go helper.ErrorToAuditLog(userId, fmt.Sprintf("Delete Record %s", req.Name), "DNS", "Delete Record", err.Error(), project.ID, project.OrganizationID, helper.GetIP(c))
			return fmt.Errorf("failed to update record: %w", err)
		}
	} else {
		// Kalau sudah kosong → hapus seluruh RRSet
		err = pdns.Records.Delete(c, tZone.Name, req.Name, powerdns.RRType(req.Type))
		if err != nil {
			go helper.ErrorToAuditLog(userId, fmt.Sprintf("Delete Record %s", req.Name), "DNS", "Delete Record", err.Error(), project.ID, project.OrganizationID, helper.GetIP(c))
			return fmt.Errorf("failed to delete record: %w", err)
		}
	}

	go helper.InfoToAuditLog(userId, fmt.Sprintf("Delete Record %s", req.Name), "DNS", "Delete Record", "Success Delete Record", project.ID, project.OrganizationID, helper.GetIP(c))
	return err
}

func NewUsecase(organizationRepo organization.Repository, projectRepo project.Repository, userRepo user.Repository, superadminroleRepo superadminrole.Repository, zoneRepo zone.Repository) Usecase {
	return &usecase{
		organizationRepo:   organizationRepo,
		projectRepo:        projectRepo,
		userRepo:           userRepo,
		superadminroleRepo: superadminroleRepo,
		zoneRepo:           zoneRepo,
	}
}
