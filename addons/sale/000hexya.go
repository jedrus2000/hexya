// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	_ "github.com/jedrus2000/hexya/addons/account"
	_ "github.com/jedrus2000/hexya/addons/procurement"
	_ "github.com/jedrus2000/hexya/addons/saleTeams"
	"github.com/jedrus2000/hexya/addons/web/controllers"
	"github.com/jedrus2000/hexya/hexya/src/server"
)

const MODULE_NAME = "sale"

func init() {
	server.RegisterModule(&server.Module{
		Name:     MODULE_NAME,
		PostInit: func() {},
	})

	controllers.BackendCSS = append(controllers.BackendCSS, "/static/sale/src/css/sale.css")
	controllers.BackendJS = append(controllers.BackendJS, "/static/sale/src/js/sale.js")
}
