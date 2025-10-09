// Copyright 2018 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package base

import (
	"testing"

	"github.com/jedrus2000/hexya/hexya/src/models"
	"github.com/jedrus2000/hexya/hexya/src/models/security"
	"github.com/jedrus2000/hexya/pool/h"
	"github.com/jedrus2000/hexya/pool/q"
	. "github.com/smartystreets/goconvey/convey"
)

func TestWithEnvironment(t *testing.T) {
	Convey("Test cases with SimulateNewEnvironment", t, func() {
		So(models.SimulateInNewEnvironment(security.SuperUserID, func(env models.Environment) {
			Convey("Create a partner.", func() {
				h.Partner().Create(env, h.Partner().NewData().SetName("test_per_class_teardown_partner"))
				partners := h.Partner().Search(env, q.Partner().Name().Equals("test_per_class_teardown_partner"))
				So(partners.Len(), ShouldEqual, 1)
			})
			Convey("Find the created partner.", func() {
				partners := h.Partner().Search(env, q.Partner().Name().Equals("test_per_class_teardown_partner"))
				So(partners.Len(), ShouldEqual, 0)
			})
		}), ShouldBeNil)
	})
	Convey("Test cases with ExecuteInNewEnvironment", t, func() {
		So(models.ExecuteInNewEnvironment(security.SuperUserID, func(env models.Environment) {
			Convey("Create a partner.", func() {
				h.Partner().Create(env, h.Partner().NewData().SetName("test_per_class_teardown_partner"))
				partners := h.Partner().Search(env, q.Partner().Name().Equals("test_per_class_teardown_partner"))
				So(partners.Len(), ShouldEqual, 1)
			})
			Convey("Find the created partner.", func() {
				partners := h.Partner().Search(env, q.Partner().Name().Equals("test_per_class_teardown_partner"))
				So(partners.Len(), ShouldEqual, 1)
			})
			Convey("Deleted the created partner.", func() {
				partners := h.Partner().Search(env, q.Partner().Name().Equals("test_per_class_teardown_partner"))
				So(partners.Unlink(), ShouldEqual, 1)
			})
		}), ShouldBeNil)
	})
}
