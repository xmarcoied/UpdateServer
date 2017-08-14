package core

import "code.videolan.org/GSoC2017/Marco/UpdateServer/database"

var db database.Impl

func SetDB(DB database.Impl) {
	db = DB
}
