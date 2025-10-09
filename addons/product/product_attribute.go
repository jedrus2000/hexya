// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"fmt"
	"log"
	"strings"

	"github.com/jedrus2000/hexya/addons/decimalPrecision"
	"github.com/jedrus2000/hexya/hexya/src/models"
	"github.com/jedrus2000/hexya/hexya/src/models/operator"
	"github.com/jedrus2000/hexya/pool/h"
	"github.com/jedrus2000/hexya/pool/m"
	"github.com/jedrus2000/hexya/pool/q"
)

func init() {

	h.ProductAttribute().DeclareModel()
	h.ProductAttribute().SetDefaultOrder("Sequence", "Name")

	h.ProductAttribute().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{Required: true, Translate: true},
		"Values": models.One2ManyField{RelationModel: h.ProductAttributeValue(), ReverseFK: "Attribute",
			JSON: "value_ids", Copy: true},
		"Sequence": models.IntegerField{Help: "Determine the display order"},
		"AttributeLines": models.One2ManyField{String: "Lines", RelationModel: h.ProductAttributeLine(),
			ReverseFK: "Attribute", JSON: "attribute_line_ids"},
		"CreateVariant": models.BooleanField{Default: models.DefaultValue(true),
			Help: "Check this if you want to create multiple variants for this attribute."},
	})

	h.ProductAttributeValue().DeclareModel()
	h.ProductAttributeValue().SetDefaultOrder("Sequence")

	h.ProductAttributeValue().AddFields(map[string]models.FieldDefinition{
		"Name":     models.CharField{String: "Value", Required: true, Translate: true},
		"Sequence": models.IntegerField{Help: "Determine the display order"},
		"Attribute": models.Many2OneField{RelationModel: h.ProductAttribute(), OnDelete: models.Cascade,
			Required: true},
		"Products": models.Many2ManyField{String: "Variants", RelationModel: h.ProductProduct(),
			JSON: "product_ids"},
		"PriceExtra": models.FloatField{String: "Attribute Price Extra",
			Compute: h.ProductAttributeValue().Methods().ComputePriceExtra(),
			Inverse: h.ProductAttributeValue().Methods().InversePriceExtra(),
			Default: models.DefaultValue(0), Digits: decimalPrecision.GetPrecision("Product Price"),
			Help: "Price Extra: Extra price for the variant with this attribute value on sale price. eg. 200 price extra, 1000 + 200 = 1200."},
		"Prices": models.One2ManyField{String: "Attribute Prices", RelationModel: h.ProductAttributePrice(),
			ReverseFK: "Value", JSON: "price_ids", ReadOnly: true},
	})

	// TODO Convert to constrains method
	//h.ProductAttributeValue().AddSQLConstraint("ValueCompanyUniq", "unique (name,attribute_id)", "This attribute value already exists !")

	h.ProductAttributeValue().Methods().ComputePriceExtra().DeclareMethod(
		`ComputePriceExtra returns the price extra for this attribute for the product
		template passed as 'active_id' in the context. Returns 0 if there is not 'active_id'.`,
		func(rs m.ProductAttributeValueSet) m.ProductAttributeValueData {
			var priceExtra float64
			if rs.Env().Context().HasKey("active_id") {
				productTmpl := h.ProductTemplate().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("active_id")})
				price := rs.Prices().Search(q.ProductAttributePrice().ProductTmpl().Equals(productTmpl))
				priceExtra = price.PriceExtra()
			}
			return h.ProductAttributeValue().NewData().SetPriceExtra(priceExtra)
		})

	h.ProductAttributeValue().Methods().InversePriceExtra().DeclareMethod(
		`InversePriceExtra sets the price extra based on the product
		template passed as 'active_id'. Does nothing if there is not 'active_id'.`,
		func(rs m.ProductAttributeValueSet, value float64) {
			if !rs.Env().Context().HasKey("active_id") {
				return
			}
			productTmpl := h.ProductTemplate().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("active_id")})
			prices := h.ProductAttributePrice().Search(rs.Env(),
				q.ProductAttributePrice().Value().In(rs).And().ProductTmpl().Equals(productTmpl))
			if !prices.IsEmpty() {
				prices.SetPriceExtra(value)
				return
			}
			updated := h.ProductAttributeValue().NewSet(rs.Env())
			for _, price := range prices.Records() {
				updated = updated.Union(price.Value())
			}
			for _, val := range rs.Subtract(updated).Records() {
				h.ProductAttributePrice().Create(rs.Env(), h.ProductAttributePrice().NewData().
					SetProductTmpl(productTmpl).
					SetValue(val).
					SetPriceExtra(value))
			}
		})

	h.ProductAttributeValue().Methods().NameGet().Extend("",
		func(rs m.ProductAttributeValueSet) string {
			if rs.Env().Context().HasKey("show_attribute") && !rs.Env().Context().GetBool("show_attribute") {
				return rs.Super().NameGet()
			}
			return fmt.Sprintf("%s: %s", rs.Attribute().Name(), rs.Name())
		})

	h.ProductAttributeValue().Methods().Unlink().Extend("",
		func(rs m.ProductAttributeValueSet) int64 {
			linkedProducts := h.ProductProduct().NewSet(rs.Env()).WithContext("active_test", false).Search(
				q.ProductProduct().AttributeValues().In(rs))
			if !linkedProducts.IsEmpty() {
				log.Panic(rs.T(`The operation cannot be completed:
You are trying to delete an attribute value with a reference on a product variant.`))
			}
			return rs.Super().Unlink()
		})

	h.ProductAttributeValue().Methods().VariantName().DeclareMethod(
		`VariantName returns a comma separated list of this product's
		attributes values of the given variable attributes'`,
		func(rs m.ProductAttributeValueSet, variableAttribute m.ProductAttributeSet) string {
			var names []string
			rSet := rs.Sorted(func(rs1, rs2 m.ProductAttributeValueSet) bool {
				return rs1.Attribute().Name() < rs2.Attribute().Name()
			})
			for _, attrValue := range rSet.Records() {
				if attrValue.Attribute().Intersect(variableAttribute).IsEmpty() {
					continue
				}
				names = append(names, attrValue.Name())
			}
			return strings.Join(names, ", ")
		})

	h.ProductAttributePrice().DeclareModel()

	h.ProductAttributePrice().AddFields(map[string]models.FieldDefinition{
		"ProductTmpl": models.Many2OneField{String: "Product Template", RelationModel: h.ProductTemplate(),
			OnDelete: models.Cascade, Required: true},
		"Value": models.Many2OneField{String: "Product Attribute Value", RelationModel: h.ProductAttributeValue(),
			OnDelete: models.Cascade, Required: true},
		"PriceExtra": models.FloatField{String: "Price Extra", Digits: decimalPrecision.GetPrecision("Product Price")},
	})

	h.ProductAttributeLine().DeclareModel()

	h.ProductAttributeLine().AddFields(map[string]models.FieldDefinition{
		"ProductTmpl": models.Many2OneField{String: "Product Template", RelationModel: h.ProductTemplate(),
			OnDelete: models.Cascade, Required: true},
		"Attribute": models.Many2OneField{RelationModel: h.ProductAttribute(),
			OnDelete: models.Restrict, Required: true,
			Constraint: h.ProductAttributeLine().Methods().CheckValidAttribute()},
		"Values": models.Many2ManyField{String: "Attribute Values", RelationModel: h.ProductAttributeValue(),
			JSON: "value_ids", Constraint: h.ProductAttributeLine().Methods().CheckValidAttribute()},
		"Name": models.CharField{Compute: h.ProductAttributeLine().Methods().ComputeName(), Stored: true,
			Depends: []string{"Attribute", "Attribute.Name", "Values", "Values.Name"}},
	})

	h.ProductAttributeLine().Methods().ComputeName().DeclareMethod(
		`Name returns a standard name with the attribute name and the values for searching`,
		func(rs m.ProductAttributeLineSet) m.ProductAttributeLineData {
			var values []string
			for _, value := range rs.Values().Records() {
				values = append(values, value.Name())
			}
			return h.ProductAttributeLine().NewData().
				SetName(rs.Attribute().Name() + ": " + strings.Join(values, ", "))
		})

	h.ProductAttributeLine().Methods().CheckValidAttribute().DeclareMethod(
		`CheckValidAttribute check that attributes values are valid for the given attributes.`,
		func(rs m.ProductAttributeLineSet) {
			for _, line := range rs.Records() {
				if !line.Values().Subtract(line.Attribute().Values()).IsEmpty() {
					log.Panic(rs.T("Error ! You cannot use this attribute with the following value."))
				}
			}
		})

	h.ProductAttributeLine().Methods().NameGet().Extend("",
		func(rs m.ProductAttributeLineSet) string {
			return rs.Attribute().NameGet()
		})

	h.ProductAttributeLine().Methods().SearchByName().Extend("",
		func(rs m.ProductAttributeLineSet, name string, op operator.Operator, additionalCond q.ProductAttributeLineCondition, limit int) m.ProductAttributeLineSet {
			// TDE FIXME: currently overriding the domain; however as it includes a
			// search on a m2o and one on a m2m, probably this will quickly become
			// difficult to compute - check if performance optimization is required
			if name != "" && op.IsPositive() {
				additionalCond = q.ProductAttributeLine().Attribute().AddOperator(op, name).
					Or().Values().AddOperator(op, name)
			}
			return rs.Super().SearchByName(name, op, additionalCond, limit)
		})

}
