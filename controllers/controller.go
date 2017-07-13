package controllers

import "github.com/xmarcoied/go-updater/model"

var db model.Impl

func SetDB(DB model.Impl) {
	db = DB
}
