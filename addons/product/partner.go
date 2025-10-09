// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"github.com/jedrus2000/hexya/addons/base"
	"github.com/jedrus2000/hexya/hexya/src/models"
	"github.com/jedrus2000/hexya/pool/h"
	"github.com/jedrus2000/hexya/pool/m"
	"github.com/jedrus2000/hexya/pool/q"
)

func init() {

	h.Partner().AddFields(map[string]models.FieldDefinition{
		"PropertyProductPricelist": models.Many2OneField{String: "Sale Pricelist", RelationModel: h.ProductPricelist(),
			Compute: h.Partner().Methods().ComputeProductPricelist(),
			Depends: []string{"Country"},
			Inverse: h.Partner().Methods().InverseProductPricelist(),
			Help:    "This pricelist will be used instead of the default one for sales to the current partner"},
		"ProductPricelist": models.Many2OneField{String: "Stored Pricelist", RelationModel: h.ProductPricelist(),
			Contexts: base.CompanyDependent},
	})

	h.Partner().Methods().ComputeProductPricelist().DeclareMethod(
		`ComputeProductPricelist returns the price list applicable for this partner`,
		func(rs m.PartnerSet) m.PartnerData {
			if rs.ID() == 0 {
				// We are processing an Onchange
				return h.Partner().NewData()
			}
			company := h.User().NewSet(rs.Env()).CurrentUser().Company()
			return h.Partner().NewData().SetPropertyProductPricelist(
				h.ProductPricelist().NewSet(rs.Env()).GetPartnerPricelist(rs, company))
		})

	h.Partner().Methods().InverseProductPricelist().DeclareMethod(
		`InverseProductPricelist sets the price list for this partner to the given list`,
		func(rs m.PartnerSet, priceList m.ProductPricelistSet) {
			var defaultForCountry m.ProductPricelistSet
			if !rs.Country().IsEmpty() {
				defaultForCountry = h.ProductPricelist().Search(rs.Env(),
					q.ProductPricelist().CountryGroupsFilteredOn(
						q.CountryGroup().CountriesFilteredOn(
							q.Country().Code().Equals(rs.Country().Code())))).Limit(1)
			} else {
				defaultForCountry = h.ProductPricelist().Search(rs.Env(),
					q.ProductPricelist().CountryGroups().IsNull()).Limit(1)
			}
			actual := rs.PropertyProductPricelist()
			if !priceList.IsEmpty() || (!actual.IsEmpty() && !defaultForCountry.Equals(actual)) {
				if priceList.IsEmpty() {
					rs.SetProductPricelist(defaultForCountry)
					return
				}
				rs.SetProductPricelist(priceList)
			}
		})

	h.Partner().Methods().CommercialFields().Extend(
		`CommercialFields`,
		func(rs m.PartnerSet) []models.FieldNamer {
			return append(rs.Super().CommercialFields(), q.Partner().PropertyProductPricelist())
		})
}
