package helper

import (
	"ai-dekadns/request"
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

func InfoToAuditLog(userId, serviceName, serviceType, methodName, action, projectId, organizationId, ipAddress string) {
	uri := os.Getenv("GATEWAY_API_URL") + "/user/activity/create"
	request := request.RequestAuditlog{
		UserID:         userId,
		ServiceName:    serviceName,
		ServiceType:    serviceType,
		MethodName:     methodName,
		Action:         action,
		Level:          "INFO",
		ProjectID:      projectId,
		OrganizationID: organizationId,
		IPAddress:      ipAddress,
	}
	payload, err := json.Marshal(request)
	if err != nil {
		log.Errorf("Error create InfoAuditlog error: %s", err.Error())
		return
	}
	_, err = PostWithHeader(uri, payload, map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	})
	if err != nil {
		log.Errorf("Error create InfoAuditlog error: %s", err.Error())
		return
	}
	log.Info("Successfully created InfoAuditLog")
}

func ErrorToAuditLog(userId, serviceName, serviceType, methodName, action, projectId, organizationId, ipAddress string) {
	uri := os.Getenv("GATEWAY_API_URL") + "/user/activity/create"
	request := request.RequestAuditlog{
		UserID:         userId,
		ServiceName:    serviceName,
		ServiceType:    serviceType,
		MethodName:     methodName,
		Action:         action,
		Level:          "ERROR",
		ProjectID:      projectId,
		OrganizationID: organizationId,
		IPAddress:      ipAddress,
	}
	payload, err := json.Marshal(request)
	if err != nil {
		log.Errorf("Error create ErrorAuditlog error: %s", err.Error())
		return
	}
	_, err = PostWithHeader(uri, payload, map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	})
	if err != nil {
		log.Errorf("Error create ErrorAuditlog error: %s", err.Error())
		return
	}
	log.Info("Successfully created ErrorAuditLog")
}
