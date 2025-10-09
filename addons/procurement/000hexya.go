// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package procurement

import (
	"github.com/jedrus2000/hexya/hexya/src/server"
	// Procurement dependencies
	_ "github.com/jedrus2000/hexya/addons/product"
)

const MODULE_NAME string = "procurement"

func init() {
	server.RegisterModule(&server.Module{
		Name:     MODULE_NAME,
		PostInit: func() {},
	})
}
