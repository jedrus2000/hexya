// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/jedrus2000/hexya/hexya/src/models"
	"github.com/jedrus2000/hexya/pool/h"
	"github.com/jedrus2000/hexya/pool/m"
)

func init() {

	h.AccountCommonJournalReport().DeclareMixinModel()
	h.AccountCommonJournalReport().InheritModel(h.AccountCommonReport())

	h.AccountCommonJournalReport().AddFields(map[string]models.FieldDefinition{
		"AmountCurrency": models.BooleanField{
			String: "With Currency",
			Help:   "Print Report with the currency column if the currency differs from the company currency."},
	})
	h.AccountCommonJournalReport().Methods().PrePrintReport().DeclareMethod(
		`PrePrintReport`,
		func(rs m.AccountCommonJournalReportSet, args struct {
			Data interface{}
		}) {
			//@api.multi
			/*def pre_print_report(self, data):
			  data['form'].update({'amount_currency': self.amount_currency})
			  return data
			*/
		})

}
