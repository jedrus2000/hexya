// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package saleTeams

import (
	"github.com/jedrus2000/hexya/hexya/src/models"
	"github.com/jedrus2000/hexya/pool/h"
)

func init() {

	h.Partner().AddFields(map[string]models.FieldDefinition{
		"Team": models.Many2OneField{String: "Sales Team", RelationModel: h.CRMTeam(),
			Help: "If set, sale team used notably for sales and assignations related to this partner"},
	})

}
