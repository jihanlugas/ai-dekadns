package dns

import (
	"ai-dekadns/app/organization"
	"ai-dekadns/app/project"
	"ai-dekadns/app/record"
	"ai-dekadns/app/superadminrole"
	"ai-dekadns/app/user"
	"ai-dekadns/app/zone"
	"ai-dekadns/constant"
	"ai-dekadns/helper"
	"ai-dekadns/model"
	"ai-dekadns/request"
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joeig/go-powerdns/v3"
	"github.com/lib/pq"
)

type Usecase interface {
	Create(c *gin.Context, req request.CreateDns) (err error)
}

type usecase struct {
	dnsRepo            Repository
	organizationRepo   organization.Repository
	projectRepo        project.Repository
	userRepo           user.Repository
	superadminroleRepo superadminrole.Repository
	zoneRepo           zone.Repository
	recordRepo         record.Repository
}

func (u usecase) Create(c *gin.Context, req request.CreateDns) (err error) {
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

		ok := user.OrganizationRole.CheckForPrivilege(constant.PrivilegeSsl, true)
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
			if role.CheckForPrivilege(constant.SuperadminPrivilegeSsl, true) {
				isAllow = true
				break
			}
		}
		if !isAllow {
			return errors.New("you are not allowed to perform this operation")
		}
	}

	// todo check dns is duplicate

	pdns := powerdns.New(os.Getenv("DNS_HOST"), "localhost", powerdns.WithHeaders(map[string]string{
		"X-API-Key": os.Getenv("DNS_API_KEY_VALUE"),
	}))

	_, err = pdns.Zones.AddMaster(c, req.Name+".", false, "d", false, "d", "d", false, []string{"ns3.cloudeka.id.", "ns4.cloudeka.id."})
	if err != nil {
		return err
	}

	err = pdns.Records.Add(c, req.Name, req.Name, powerdns.RRTypeNS, 30, []string{"ns3.cloudeka.id.", "ns4.cloudeka.id."})
	if err != nil {
		return err
	}

	// todo insert table zone
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

	// todo get status and IsCustomNs

	err = u.zoneRepo.Create(tZone)
	if err != nil {
		return err
	}

	// todo insert table record
	tRecord1 := model.Record{
		ID:            helper.GetUniqueID(),
		ZoneId:        tZone.ID,
		TypeId:        "NS",
		Status:        false,
		Cache:         false,
		Name:          tZone.Name,
		Content:       pq.StringArray{"ns3.cloudeka.id."},
		TTL:           30,
		IsHidden:      false,
		PricePerMonth: 0,
		PricePerHour:  0,
		CreatedBy:     userId,
		UpdatedBy:     userId,
	}
	err = u.recordRepo.Create(tRecord1)
	if err != nil {
		return err
	}

	tRecord2 := model.Record{
		ID:            helper.GetUniqueID(),
		ZoneId:        tZone.ID,
		TypeId:        "NS",
		Status:        false,
		Cache:         false,
		Name:          tZone.Name,
		Content:       pq.StringArray{"ns4.cloudeka.id."},
		TTL:           30,
		IsHidden:      false,
		PricePerMonth: 0,
		PricePerHour:  0,
		CreatedBy:     userId,
		UpdatedBy:     userId,
	}

	err = u.recordRepo.Create(tRecord2)
	if err != nil {
		return err
	}

	go helper.InfoToAuditLog(userId, "Create DNS"+req.Name, "DNS", "Create DNS", "Success Create DNS", projectId, orgId, helper.GetIP(c))
	return err
}

func NewUsecase(dnsRepo Repository, organizationRepo organization.Repository, projectRepo project.Repository, userRepo user.Repository, superadminroleRepo superadminrole.Repository, zoneRepo zone.Repository, recordRepo record.Repository) Usecase {
	return &usecase{
		dnsRepo:            dnsRepo,
		organizationRepo:   organizationRepo,
		projectRepo:        projectRepo,
		userRepo:           userRepo,
		superadminroleRepo: superadminroleRepo,
		zoneRepo:           zoneRepo,
		recordRepo:         recordRepo,
	}
}
