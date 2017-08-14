package database

import (
	"time"

	_ "github.com/lib/pq" // for database
)

type Rule struct {
	ID          int
	ReleaseID   int
	TimeRule    TimeRule
	OsRule      OsRule
	VersionRule VersionRule
	IPRule      IPRule
	RollRule    RollRule
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

type RollRule struct {
	RuleID            int
	RollingPercentage int
}

type IPRule struct {
	RuleID int
	IP     string
}
