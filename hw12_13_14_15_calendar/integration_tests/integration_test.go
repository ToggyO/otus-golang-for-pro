//go:build integration
// +build integration

package integration_tests

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/integration_tests/suites"
)

type IntegrationTest struct {
	suites.ApplicationActionsSuite
}

func TestIntegrationTest(t *testing.T) {
	suite.Run(t, new(IntegrationTest))
}
