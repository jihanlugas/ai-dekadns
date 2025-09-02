package zone

import (
	"ai-dekadns/app/organization"
	"ai-dekadns/app/project"
	"ai-dekadns/app/superadminrole"
	"ai-dekadns/app/user"
	"ai-dekadns/constant"
	"ai-dekadns/helper"
	"ai-dekadns/model"
	"ai-dekadns/request"
	"ai-dekadns/response"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joeig/go-powerdns/v3"
)

type Usecase interface {
	Page(c *gin.Context, req request.PageZone) (pagination *model.Pagination, err error)
	GetById(c *gin.Context, id string) (res response.Zone, err error)
	Create(c *gin.Context, req request.CreateZone) (err error)
	Delete(c *gin.Context, id string) (err error)
}

type usecase struct {
	organizationRepo   organization.Repository
	projectRepo        project.Repository
	userRepo           user.Repository
	superadminroleRepo superadminrole.Repository
	zoneRepo           Repository
}

func (u usecase) Page(c *gin.Context, req request.PageZone) (pagination *model.Pagination, err error) {
	var userId = helper.GetUserID(c)
	var role = helper.GetRole(c)
	var projectId = req.ProjectId

	project, err := u.projectRepo.GetById(projectId)
	if err != nil {
		return pagination, err
	}

	user, err := u.userRepo.GetById(userId, "Organization", "OrganizationRole")
	if err != nil {
		return pagination, err
	}

	//orgId := project.OrganizationID

	if role != constant.RoleSuperadmin {
		// check IDOR
		if user.OrganizationID != project.OrganizationID {
			return pagination, errors.New("you are not allowed to perform this operation")
		}

		// check user role permission
		// When organization role is not set, the user have no privileges.
		if user.OrganizationRole == nil {
			return pagination, errors.New("you are not allowed to perform this operation")
		}

		ok := user.OrganizationRole.CheckForPrivilege(constant.PrivilegeDns, false)
		if !ok {
			return pagination, errors.New("you are not allowed to perform this operation")
		}
	} else {
		rolePrivileges, err := u.superadminroleRepo.GetBySuperadminId(user.ID)
		if err != nil {
			return pagination, err
		}

		isAllow := false
		for _, role := range rolePrivileges {
			if role.CheckForPrivilege(constant.SuperadminPrivilegeDns, false) {
				isAllow = true
				break
			}
		}
		if !isAllow {
			return pagination, errors.New("you are not allowed to perform this operation")
		}
	}

	page, _ := strconv.Atoi(req.Page)
	limit, _ := strconv.Atoi(req.Limit)
	pageReq := model.Pagination{
		Page:  page,
		Limit: limit,
	}

	pagination, err = u.zoneRepo.Page(req, pageReq)
	if err != nil {
		return pagination, err
	}

	return pagination, nil
}

func (u usecase) GetById(c *gin.Context, id string) (res response.Zone, err error) {
	var userId = helper.GetUserID(c)
	var role = helper.GetRole(c)
	var tZone model.Zone

	tZone, err = u.zoneRepo.GetById(id)
	if err != nil {
		return res, err
	}

	project, err := u.projectRepo.GetById(tZone.ProjectId)
	if err != nil {
		return res, err
	}

	user, err := u.userRepo.GetById(userId, "Organization", "OrganizationRole")
	if err != nil {
		return res, err
	}

	//orgId := project.OrganizationID

	if role != constant.RoleSuperadmin {
		// check IDOR
		if user.OrganizationID != project.OrganizationID {
			return res, errors.New("you are not allowed to perform this operation")
		}

		// check user role permission
		// When organization role is not set, the user have no privileges.
		if user.OrganizationRole == nil {
			return res, errors.New("you are not allowed to perform this operation")
		}

		ok := user.OrganizationRole.CheckForPrivilege(constant.PrivilegeDns, false)
		if !ok {
			return res, errors.New("you are not allowed to perform this operation")
		}
	} else {
		rolePrivileges, err := u.superadminroleRepo.GetBySuperadminId(user.ID)
		if err != nil {
			return res, err
		}

		isAllow := false
		for _, role := range rolePrivileges {
			if role.CheckForPrivilege(constant.SuperadminPrivilegeDns, false) {
				isAllow = true
				break
			}
		}
		if !isAllow {
			return res, errors.New("you are not allowed to perform this operation")
		}
	}

	pdns := powerdns.New(os.Getenv("DNS_HOST"), "localhost", powerdns.WithHeaders(map[string]string{
		"X-API-Key": os.Getenv("DNS_API_KEY_VALUE"),
	}))

	resRecord := []response.Record{}
	pdnszone, err := pdns.Zones.Get(c, tZone.Name)
	if err != nil {
		return res, err
	}

	for _, rrset := range pdnszone.RRsets {
		for _, record := range rrset.Records {
			content := helper.FromPowerDNS(*rrset.Type, *record.Content)
			name := strings.TrimSuffix(*rrset.Name, ".")
			newRecord := response.Record{
				Name:    &name,
				Type:    rrset.Type,
				TTL:     rrset.TTL,
				Content: &content,
			}
			resRecord = append(resRecord, newRecord)
		}
	}

	res = response.Zone{
		Zone:    tZone,
		Records: resRecord,
	}

	return res, err
}

func (u usecase) Create(c *gin.Context, req request.CreateZone) (err error) {
	var userId = helper.GetUserID(c)
	var role = helper.GetRole(c)
	var projectId = req.ProjectId

	project, err := u.projectRepo.GetById(projectId)
	if err != nil {
		return err
	}

	user, err := u.userRepo.GetById(userId, "Organization", "OrganizationRole")
	if err != nil {
		return err
	}

	orgId := project.OrganizationID

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

	_, err = pdns.Zones.AddMaster(c, req.Name+".", false, "d", false, "d", "d", false, []string{"ns3.cloudeka.id.", "ns4.cloudeka.id."})
	if err != nil {
		go helper.ErrorToAuditLog(userId, fmt.Sprintf("Create Zone %s", req.Name), "DNS", "Create Zone", err.Error(), projectId, orgId, helper.GetIP(c))
		return err
	}

	err = pdns.Records.Add(c, req.Name, req.Name, powerdns.RRTypeNS, 30, []string{"ns3.cloudeka.id.", "ns4.cloudeka.id."})
	if err != nil {
		go helper.ErrorToAuditLog(userId, fmt.Sprintf("Create Zone %s", req.Name), "DNS", "Create Zone", err.Error(), projectId, orgId, helper.GetIP(c))
		return err
	}

	tZone := model.Zone{
		ID:             helper.GetUniqueID(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		Name:           req.Name,
		Status:         "",
		IsCustomNs:     false,
		IsDnssec:       "",
		CreatedBy:      userId,
		UpdatedBy:      userId,
	}

	err = u.zoneRepo.Create(tZone)
	if err != nil {
		go helper.ErrorToAuditLog(userId, fmt.Sprintf("Create Zone %s", req.Name), "DNS", "Create Zone", err.Error(), projectId, orgId, helper.GetIP(c))
		return err
	}

	go helper.InfoToAuditLog(userId, fmt.Sprintf("Create Zone %s", req.Name), "DNS", "Create Zone", "Success Create Zone", projectId, orgId, helper.GetIP(c))
	return err
}

func (u usecase) Delete(c *gin.Context, id string) (err error) {
	var userId = helper.GetUserID(c)
	var role = helper.GetRole(c)
	var tZone model.Zone

	tZone, err = u.zoneRepo.GetById(id)
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

	//orgId := project.OrganizationID

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

		ok := user.OrganizationRole.CheckForPrivilege(constant.PrivilegeDns, false)
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
			if role.CheckForPrivilege(constant.SuperadminPrivilegeDns, false) {
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

	err = pdns.Zones.Delete(c, tZone.Name)
	if err != nil {
		go helper.ErrorToAuditLog(userId, fmt.Sprintf("Delete Zone %s", tZone.Name), "DNS", "Delete Zone", err.Error(), project.ID, project.OrganizationID, helper.GetIP(c))
		return err
	}

	go helper.InfoToAuditLog(userId, fmt.Sprintf("Delete Zone %s", tZone.Name), "DNS", "Delete Zone", "Success Delete Zone", project.ID, project.OrganizationID, helper.GetIP(c))

	return nil
}

func NewUsecase(organizationRepo organization.Repository, projectRepo project.Repository, userRepo user.Repository, superadminroleRepo superadminrole.Repository, zoneRepo Repository) Usecase {
	return &usecase{
		organizationRepo:   organizationRepo,
		projectRepo:        projectRepo,
		userRepo:           userRepo,
		superadminroleRepo: superadminroleRepo,
		zoneRepo:           zoneRepo,
	}
}
