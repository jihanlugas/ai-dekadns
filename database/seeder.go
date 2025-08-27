package database

import (
	"ai-dekadns/model"
	"log"
)

func SeedTypes() {
	var err error
	db := GetCorePostsqlConn()

	// defind ID biar gak ada yang duplicate

	types := []model.Type{
		{
			ID:   "a73b8707-4f97-445f-a71f-82fcfb5ae6b4",
			Name: "A",
		},
		{
			ID:   "0ebbdec1-1995-4a80-b4c0-e22d3fff730e",
			Name: "AAAA",
		},
		{
			ID:   "442906bc-bb25-4cf1-a751-400d67cfe138",
			Name: "MX",
		},
		{
			ID:   "bf9cea16-71f9-4fe4-b30d-e8d11b6c7b32",
			Name: "TXT",
		},
		{
			ID:   "d56d3d2f-71b2-4f49-a695-b92d1eb740d3",
			Name: "SRV",
		},
		{
			ID:   "52cb6eaf-983c-4b97-92f8-cc58ccbe3586",
			Name: "CNAME",
		},
		{
			ID:   "00273ca3-0c09-478b-90da-e7c37187036e",
			Name: "SPF",
		},
		{
			ID:   "d627fc18-fc1d-4083-89fc-446c2067a778",
			Name: "DKIM",
		},
		{
			ID:   "fe96162c-3682-43f2-9e0f-1a5cb90b5b80",
			Name: "PTR",
		},
		{
			ID:   "c4b109f2-5299-4d9c-98e0-04a665d9fdb0",
			Name: "SOA",
		},
		{
			ID:   "72248214-1bb1-4b61-b164-3e73fb6047c7",
			Name: "CAA",
		},
		{
			ID:   "52df8f59-ae3c-4360-84bf-e6682ad1c057",
			Name: "NS",
		},
	}
	err = db.Create(&types).Error
	if err != nil {
		log.Fatal(err)
	}
}
