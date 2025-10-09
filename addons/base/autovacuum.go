// Copyright 2020 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package base

import (
	"github.com/jedrus2000/hexya/hexya/src/models"
	"github.com/jedrus2000/hexya/hexya/src/models/security"
	"github.com/jedrus2000/hexya/pool/h"
	"github.com/jedrus2000/hexya/pool/m"
)

// GCUserLogs garbage collects old user log instances
func autoVacuum_GCUserLogs(rs m.AutoVacuumSet) {
	res := rs.Env().Cr().Execute(`
DELETE FROM user_log log1 WHERE EXISTS (
	SELECT 1 FROM user_log log2
	WHERE log1.create_uid = log2.create_uid
	AND log1.create_date < log2.create_date
)`)
	n, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	log.Info("GC'd %d user log entries", n)
}

// PowerOn executes a vacuum of internal resources.
// Override this method to add your own garbage collections.
func autoVacuum_PowerOn(rs m.AutoVacuumSet) {
	if rs.Env().Uid() != security.SuperUserID {
		panic("Access Denied")
	}
	h.Attachment().NewSet(rs.Env()).FileGC()
	rs.GCUserLogs()
}

func init() {
	models.NewManualModel("AutoVacuum")
	h.AutoVacuum().NewMethod("GCUserLogs", autoVacuum_GCUserLogs)
	h.AutoVacuum().NewMethod("PowerOn", autoVacuum_PowerOn)
}
