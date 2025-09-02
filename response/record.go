package response

import "github.com/joeig/go-powerdns/v3"

type Record struct {
	Name    *string          `json:"name"`
	Type    *powerdns.RRType `json:"type"`
	TTL     *uint32          `json:"ttl"`
	Content *string          `json:"content"`
}
