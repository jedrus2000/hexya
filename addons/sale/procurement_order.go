// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"github.com/jedrus2000/hexya/hexya/src/models"
	"github.com/jedrus2000/hexya/pool/h"
)

func init() {

	h.ProcurementOrder().AddFields(map[string]models.FieldDefinition{
		"SaleLine": models.Many2OneField{String: "Sale Order Line", RelationModel: h.SaleOrderLine()},
	})

}
