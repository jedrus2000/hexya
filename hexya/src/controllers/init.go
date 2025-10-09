// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package controllers

import (
	"github.com/jedrus2000/hexya/hexya/src/server"
	"github.com/jedrus2000/hexya/hexya/src/tools/logging"
)

var log logging.Logger

// BootStrap creates the actual controllers from the controllers registry.
// This function must be called before starting the http server.
func BootStrap() {
	Registry.createRoutes(server.GetServer().Group("/"))
}

func init() {
	log = logging.GetLogger("controllers")
	Registry = newGroup("/")
}
