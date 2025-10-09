// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/jedrus2000/hexya/hexya/src/actions"
	"github.com/jedrus2000/hexya/pool/h"
	"github.com/jedrus2000/hexya/pool/m"
)

func init() {

	h.AccountUnreconcile().DeclareTransientModel()
	h.AccountUnreconcile().Methods().TransUnrec().DeclareMethod(
		`TransUnrec`,
		func(rs m.AccountUnreconcileSet) *actions.Action {
			if ids := rs.Env().Context().GetIntegerSlice("active_ids"); len(ids) > 0 {
				h.AccountMoveLine().Browse(rs.Env(), ids).RemoveMoveReconcile()
			}
			return &actions.Action{
				Type: actions.ActionCloseWindow,
			}
		})

}
