// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package procurement

import (
	"github.com/jedrus2000/hexya/addons/base"
	"github.com/jedrus2000/hexya/pool/h"
)

func init() {

	h.ProcurementOrder().Methods().AllowAllToGroup(base.GroupUser)
	h.ProcurementGroup().Methods().AllowAllToGroup(base.GroupUser)
	h.ProcurementRule().Methods().AllowAllToGroup(base.GroupUser)

}
