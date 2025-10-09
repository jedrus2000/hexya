// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package decimalPrecision

import (
	// decimalPrecision depends on base module
	_ "github.com/jedrus2000/hexya/addons/base"
	"github.com/jedrus2000/hexya/hexya/src/server"
	"github.com/jedrus2000/hexya/hexya/src/tools/nbutils"
)

const MODULE_NAME string = "decimalPrecision"

// Precisions is Hexya's global map for storing
// precision data for each application.
var Precisions = map[string]nbutils.Digits{}

func init() {
	server.RegisterModule(&server.Module{
		Name:     MODULE_NAME,
		PostInit: func() {},
	})
}

// GetPrecision returns the precision for the given application
func GetPrecision(app string) nbutils.Digits {
	res, exists := Precisions[app]
	if !exists {
		return nbutils.Digits{Precision: 16, Scale: 2}
	}
	return res
}
