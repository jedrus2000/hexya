// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package saleTeams

import (
	"github.com/jedrus2000/hexya/hexya/src/models"
	"github.com/jedrus2000/hexya/pool/h"
)

func init() {

	h.User().AddFields(map[string]models.FieldDefinition{
		"SaleTeam": models.Many2OneField{String: "Sales Team", RelationModel: h.CRMTeam(),
			Help: `Sales Team the user is member of.
Used to compute the members of a sales team through the inverse one2many`},
	})

}
