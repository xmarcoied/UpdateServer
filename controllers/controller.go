package controllers

import "code.videolan.org/GSoC2017/Marco/UpdateServer/model"

var db model.Impl

func SetDB(DB model.Impl) {
	db = DB
}

type (
	ReleaseController struct{}
	RequestController struct{}
	ChannelController struct{}
	RulesController   struct{}
)
