// Core package provides the main functionalities for the update server
package core

import "code.videolan.org/GSoC2017/Marco/UpdateServer/database"

var db database.Impl

// SetDB set db for dialect
func SetDB(DB database.Impl) {
	db = DB
}
