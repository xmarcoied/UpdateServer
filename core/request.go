package core

import (
	"time"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
)

// GetRequests return all requests recorded at the database associated with a channel and product
func GetRequests(channel, product string) []database.UpdateRequest {
	var requests []database.UpdateRequest
	db.DB.Where("product = ? AND channel = ?", product, channel).Order("created_at desc").Find(&requests)
	ProcessCreatedSince(&requests)
	return requests
}

// ProcessCreatedSince initiate the CreateSince Section at update_requests
func ProcessCreatedSince(requests *[]database.UpdateRequest) {
	TimeNow := time.Now().UTC()
	for i := 0; i < len(*requests); i++ {
		Duration := TimeNow.Sub((*requests)[i].CreatedAt.UTC())
		(*requests)[i].CreatedSince.Month = int(Duration.Hours() / (24 * 30))
		(*requests)[i].CreatedSince.Day = TimeNow.Day() - (*requests)[i].CreatedAt.UTC().Day()
		(*requests)[i].CreatedSince.Hour = TimeNow.Hour() - (*requests)[i].CreatedAt.UTC().Hour()
		(*requests)[i].CreatedSince.Minute = TimeNow.Minute() - (*requests)[i].CreatedAt.UTC().Minute()
		(*requests)[i].CreatedSince.Second = TimeNow.Second() - (*requests)[i].CreatedAt.UTC().Second()
	}
}
