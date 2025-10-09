// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package saleTeams

import (
	"github.com/jedrus2000/hexya/addons/base"
	"github.com/jedrus2000/hexya/addons/web/controllers"
	//
	// github.com/jedrus2000/hexya/addons/webKanban v0.0.30
	// This project is archived. All features of this module are now in the web module
	//
	"github.com/jedrus2000/hexya/hexya/src/models/security"
	"github.com/jedrus2000/hexya/hexya/src/server"
)

const MODULE_NAME = "saleTeams"

func init() {
	server.RegisterModule(&server.Module{
		Name:     MODULE_NAME,
		PostInit: func() {},
	})

	GroupSaleSalesman = security.Registry.NewGroup("sale_teams_group_sale_salesman", "User: Own Documents Only", base.GroupUser)
	GroupSaleSalesmanAllLeads = security.Registry.NewGroup("sale_teams_group_sale_salesman_all_leads", "User: All Documents", GroupSaleSalesman)
	GroupSaleManager = security.Registry.NewGroup("sale_teams_group_sale_manager", "Manager", GroupSaleSalesmanAllLeads)

	controllers.BackendLess = append(controllers.BackendLess, "/static/saleTeams/src/less/sales_team_dashboard.less")
	controllers.BackendCSS = append(controllers.BackendCSS, "/static/saleTeams/src/css/sales_team.css")
	controllers.BackendJS = append(controllers.BackendJS,
		"/static/saleTeams/src/js/sales_team.js",
		"/static/saleTeams/src/js/sales_team_dashboard.js",
	)
}
