package suites

type ApplicationActionsSuite struct {
	ServiceFixtureSuite
}

func (ap *ApplicationActionsSuite) SetupSuite() {
	ap.Init()
	ap.InitEventsTests()
}
