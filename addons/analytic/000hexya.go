// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package analytic

import (
	// Import dependencies
	_ "github.com/jedrus2000/hexya/addons/product"
	"github.com/jedrus2000/hexya/hexya/src/models/security"
	"github.com/jedrus2000/hexya/hexya/src/server"
)

const MODULE_NAME string = "analytic"

// GroupAnalyticAccounting is the group of the people allowed manage analytic accounting
var GroupAnalyticAccounting *security.Group

func init() {
	server.RegisterModule(&server.Module{
		Name:     MODULE_NAME,
		PostInit: func() {},
	})

	GroupAnalyticAccounting = security.Registry.NewGroup("analytic_group_analytic_accounting", "Analytic Accounting")
}
