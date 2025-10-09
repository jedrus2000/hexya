// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package web

import (
	"github.com/jedrus2000/hexya/hexya/src/models/security"
	"github.com/jedrus2000/hexya/pool/h"
)

func init() {
	h.Filter().Methods().AllowAllToGroup(security.GroupEveryone)
}
