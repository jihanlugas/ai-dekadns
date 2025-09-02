package response

import "ai-dekadns/model"

type Zone struct {
	model.Zone
	Records []Record `json:"records"`
}
