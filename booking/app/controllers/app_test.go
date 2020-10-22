package controllers_test

import (
	"testing"

	"github.com/revel/examples/booking/app/tmp/run"
	"github.com/revel/modules/server-engine/gohttptest/testsuite"
)

//  go test -coverprofile=coverage.out github.com/revel/examples/booking/app/controllers/  -args -revel.importPath=github.com/revel/examples/booking
func TestMain(m *testing.M) {
	testsuite.RevelTestHelper(m, "dev", run.Run)
}

func TestIndex(t *testing.T) {
	tester := testsuite.NewTestSuite(t)
	tester.Get("/").AssertOk()
}
