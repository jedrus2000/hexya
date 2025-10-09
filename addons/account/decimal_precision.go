package account

import (
	"github.com/jedrus2000/hexya/addons/decimalPrecision"
	"github.com/jedrus2000/hexya/hexya/src/tools/nbutils"
)

func init() {
	decimalPrecision.Precisions["Payment Terms"] = nbutils.Digits{Precision: 6, Scale: 2}
}
