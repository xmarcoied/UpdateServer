package model

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // for database
)

type Rule struct {
	gorm.Model
	Name      string
	ReleaseID int
	TimeRule  TimeRule
}

type TimeRule struct {
	RuleID    int
	StartTime time.Time
	EndTime   time.Time
}
