package core

import (
	"time"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
)

// GetRequests return all requests recorded at the database associated with a channel and product
func GetRequests(query string) []database.UpdateRequest {
	var requests []database.UpdateRequest
	db.DB.Where(query).Order("created_at desc").Find(&requests)
	ProcessCreatedSince(&requests)
	return requests
}

func NewRequest(request database.UpdateRequest) {
	db.DB.Create(&request)
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

// GetSignature function return the signature of a given release.
func GetSignature(release_id int) string {
	var release database.Release
	db.DB.First(&release, "id = ?", release_id)
	return release.Signature
}

// ReleaseMap function map the incoming update_request with the most suitable update release
// Ordered
func ReleaseMap(request database.UpdateRequest) (database.Release, bool) {
	var emptyrelease database.Release
	// First , find and count the available releases match the request specs
	var releasescount int
	var releases []database.Release

	db.DB.Where("product = ? AND channel = ? AND os = ? AND os_arch = ? AND os_ver >= ? AND active = true",
		request.Product, request.Channel, request.OS, request.OsArch, request.OsVer).
		Order("product_version desc").Find(&releases).Count(&releasescount)

	if releasescount == 0 {
		return emptyrelease, false
	} else {
		for _, release := range releases {
			rules := GetRules(release)
			if len(rules) == 0 {
				return releases[0], true
			}
			if CheckTimeRule(release) == false {
				return emptyrelease, false
			}
			if CheckOsRule(release, request) == false {
				return emptyrelease, false
			}
			if CheckVersionRule(release, request) == false {
				return emptyrelease, false
			}
			if CheckRollRule(release) == false {
				return emptyrelease, false
			}
			found, check := CheckIPRule(release, request)

			if found == true && check == true {
				return release, true
			} else {
				return emptyrelease, false
			}

			return releases[0], true
		}

		return emptyrelease, false
	}
}
