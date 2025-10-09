// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/jedrus2000/hexya/hexya/src/actions"
	"github.com/jedrus2000/hexya/hexya/src/tools/strutils"
	"github.com/jedrus2000/hexya/pool/h"
	"github.com/jedrus2000/hexya/pool/m"
)

func init() {

	h.AccountInvoiceConfirm().DeclareTransientModel()
	h.AccountInvoiceConfirm().Methods().InvoiceConfirm().DeclareMethod(
		`InvoiceConfirm`,
		func(rs m.AccountInvoiceConfirmSet) *actions.Action {
			for _, rec := range h.AccountInvoice().Browse(rs.Env(), rs.Env().Context().GetIntegerSlice("active_ids")).Records() {
				if !strutils.IsIn(rec.State(), "draft", "proforma", "proforma2") {
					panic(rs.T(`Selected invoice(s) cannot be confirmed as they are not in 'Draft' or 'Pro-Forma' state.`))
				}
				rec.ActionInvoiceOpen()
			}
			return &actions.Action{
				Type: actions.ActionCloseWindow,
			}
		})

	h.AccountInvoiceCancel().DeclareTransientModel()
	h.AccountInvoiceCancel().Methods().InvoiceCancel().DeclareMethod(
		`InvoiceCancel`,
		func(rs m.AccountInvoiceCancelSet) *actions.Action {
			for _, rec := range h.AccountInvoice().Browse(rs.Env(), rs.Env().Context().GetIntegerSlice("active_ids")).Records() {
				if strutils.IsIn(rec.State(), "cancel", "paid") {
					panic(rs.T(`Selected invoice(s) cannot be cancelled as they are already in 'Cancelled' or 'Done' state.`))
				}
				rec.ActionInvoiceCancel()
			}
			return &actions.Action{
				Type: actions.ActionCloseWindow,
			}
		})

}
