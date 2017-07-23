package model

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // for database
)

type Rule struct {
	gorm.Model
	ReleaseID   int
	TimeRule    TimeRule
	OsRule      OsRule
	VersionRule VersionRule
	IPRule      IPRule
}

type TimeRule struct {
	RuleID    int
	StartTime time.Time
	EndTime   time.Time
}

type OsRule struct {
	RuleID    int
	OsVersion string
}

type VersionRule struct {
	RuleID         int
	ProductVersion string
}

type IPRule struct {
	RuleID int
	IP     string
}
