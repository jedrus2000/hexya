// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/jedrus2000/hexya/hexya/src/models"
	"github.com/jedrus2000/hexya/hexya/src/models/types"
	"github.com/jedrus2000/hexya/pool/h"
	"github.com/jedrus2000/hexya/pool/m"
)

func init() {

	h.AccountCommonAccountReport().DeclareMixinModel()
	h.AccountCommonAccountReport().InheritModel(h.AccountCommonReport())

	h.AccountCommonAccountReport().AddFields(map[string]models.FieldDefinition{
		"DisplayAccount": models.SelectionField{
			String: "Display Accounts",
			Selection: types.Selection{
				"all":      "All",
				"movement": "With movements",
				"not_zero": "With balance is not equal to 0"},
			Required: true,
			Default:  models.DefaultValue("movement")},
	})
	h.AccountCommonAccountReport().Methods().PrePrintReport().DeclareMethod(
		`PrePrintReport`,
		func(rs m.AccountCommonAccountReportSet, data map[string]interface{}) map[string]interface{} {
			data["form"] = rs.Read([]string{"display_account"})[0]
			return data
		})

}
