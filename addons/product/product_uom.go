// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"log"

	"github.com/jedrus2000/hexya/hexya/src/models"
	"github.com/jedrus2000/hexya/hexya/src/models/types"
	"github.com/jedrus2000/hexya/hexya/src/tools/nbutils"
	"github.com/jedrus2000/hexya/pool/h"
	"github.com/jedrus2000/hexya/pool/m"
)

func init() {

	h.ProductUomCategory().DeclareModel()

	h.ProductUomCategory().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Name", Required: true, Translate: true},
	})

	h.ProductUom().DeclareModel()
	h.ProductUom().SetDefaultOrder("Name")

	h.ProductUom().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Unit of Measure", Required: true, Translate: true},
		"Category": models.Many2OneField{RelationModel: h.ProductUomCategory(), Required: true, OnDelete: models.Cascade,
			Help: `Conversion between Units of Measure can only occur if they belong to the same category.
The conversion will be made based on the ratios.`},
		"Factor": models.FloatField{String: "Ratio", Default: models.DefaultValue(1.0), Required: true,
			Help: `How much bigger or smaller this unit is compared to the reference Unit of Measure for this category:
1 * (reference unit) = ratio * (this unit)`},
		"FactorInv": models.FloatField{String: "Bigger Ratio", Compute: h.ProductUom().Methods().ComputeFactorInv(),
			Required: true,
			Help: `How many times this Unit of Measure is bigger than the reference Unit of Measure in this category:
1 * (this unit) = ratio * (reference unit)`,
			Depends: []string{"Factor"}},
		"Rounding": models.FloatField{String: "Rounding Precision", Default: models.DefaultValue(0.01),
			Required: true, Help: `The computed quantity will be a multiple of this value.
Use 1.0 for a Unit of Measure that cannot be further split, such as a piece.`},
		"Active": models.BooleanField{Default: models.DefaultValue(true), Required: true,
			Help: "Uncheck the active field to disable a unit of measure without deleting it."},
		"UomType": models.SelectionField{String: "Type", Selection: types.Selection{
			"bigger":    "Bigger than the reference Unit of Measure",
			"reference": "Reference Unit of Measure for this category",
			"smaller":   "Smaller than the reference Unit of Measure",
		}, Default: models.DefaultValue("reference"), Required: true,
			OnChange: h.ProductUom().Methods().OnchangeUomType()},
	})

	h.ProductUom().AddSQLConstraint("FactorGtZero", "CHECK (factor!=0)", "The conversion ratio for a unit of measure cannot be 0!")
	h.ProductUom().AddSQLConstraint("RoundingGtZero", "CHECK (rounding>0)", "The rounding precision must be greater than 0!")

	h.ProductUom().Methods().ComputeFactorInv().DeclareMethod(
		`ComputeFactorInv computes the inverse factor`,
		func(rs m.ProductUomSet) m.ProductUomData {
			var factorInv float64
			if rs.Factor() != 0 {
				factorInv = 1 / rs.Factor()
			}
			return h.ProductUom().NewData().SetFactorInv(factorInv)
		})

	h.ProductUom().Methods().OnchangeUomType().DeclareMethod(
		`OnchangeUomType updates factor when the UoM type is changed`,
		func(rs m.ProductUomSet) m.ProductUomData {
			res := h.ProductUom().NewData()
			if rs.UomType() == "reference" {
				res.SetFactor(1)
			}
			return res

		})

	h.ProductUom().Methods().Create().Extend("",
		func(rs m.ProductUomSet, data m.ProductUomData) m.ProductUomSet {
			if data.FactorInv() != 0 {
				data.SetFactor(1 / data.FactorInv())
				data.SetFactorInv(0)
			}
			return rs.Super().Create(data)
		})

	h.ProductUom().Methods().Write().Extend("",
		func(rs m.ProductUomSet, vals m.ProductUomData) bool {
			if vals.HasFactorInv() {
				var factor float64
				if vals.FactorInv() != 0 {
					factor = 1 / vals.FactorInv()
				}
				vals.SetFactor(factor)
				vals.SetFactorInv(0)
			}
			return rs.Super().Write(vals)
		})

	h.ProductUom().Methods().ComputeQuantity().DeclareMethod(
		`ComputeQuantity converts the given qty from this UoM to toUnit UoM. If round is true,
		the result will be rounded to toUnit rounding.

		It panics if both units are not from the same category`,
		func(rs m.ProductUomSet, qty float64, toUnit m.ProductUomSet, round bool) float64 {
			if rs.IsEmpty() {
				return qty
			}
			rs.EnsureOne()
			if !rs.Category().Equals(toUnit.Category()) {
				log.Panic(rs.T("Conversion from Product UoM %s to Default UoM %s is not possible as they both belong to different Category!.", rs.Name(), toUnit.Name()))
			}
			amount := qty / rs.Factor()
			if toUnit.IsEmpty() {
				return amount
			}
			amount = amount * toUnit.Factor()
			if round {
				amount = nbutils.Round(amount, toUnit.Rounding())
			}
			return amount
		})

	h.ProductUom().Methods().ComputePrice().DeclareMethod(
		`ComputePrice computes the price per 'toUnit' from the given price per this unit`,
		func(rs m.ProductUomSet, price float64, toUnit m.ProductUomSet) float64 {
			rs.EnsureOne()
			if price == 0 || toUnit.IsEmpty() || rs.Equals(toUnit) {
				return price
			}
			if !rs.Category().Equals(toUnit.Category()) {
				return price
			}
			amount := price * rs.Factor()
			return amount / toUnit.Factor()
		})

}
