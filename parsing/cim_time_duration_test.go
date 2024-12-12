package parsing

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// Unit test suite for all CimTimeDuration parsing functions
type CimTimeDurationUnitTestSuite struct {
	suite.Suite
	// Fixtures
	testJson     string
	testExpected testCimTime
}

// Fixture objects
type testCimTime struct {
	Test            string          `json:"Test"`
	CimTimeDuration CimTimeDuration `json:"LeaseDuration"`
}

func (suite *CimTimeDurationUnitTestSuite) SetupSuite() {
	// Fixtures
	suite.testJson = `{
	    "Test": "Test String",
		"LeaseDuration":  {
            "Ticks":  6912000000000,
            "Days":  8,
            "Hours":  0,
            "Milliseconds":  0,
            "Minutes":  0,
            "Seconds":  0,
            "TotalDays":  8,
            "TotalHours":  192,
            "TotalMilliseconds":  691200000,
            "TotalMinutes":  11520,
            "TotalSeconds":  691200
        }
    }`
	suite.testExpected = testCimTime{
		Test: "Test String",
		CimTimeDuration: CimTimeDuration{
			Duration: 8 * 24 * time.Hour,
		},
	}
}

func TestCimTimeDurationUnitTestSuite(t *testing.T) {
	suite.Run(t, &CimTimeDurationUnitTestSuite{})
}

func (suite *CimTimeDurationUnitTestSuite) TestUnmarshalJSON() {
	suite.T().Parallel()

	suite.Run("should unmarshal the CimInstance duration json to CimTimeDuration", func() {
		cimTime := CimTimeDuration{}
		expectedTimeDuration, err := time.ParseDuration("1h30m")
		suite.Require().NoError(err)
		expectedCimTimeDuration := CimTimeDuration{Duration: expectedTimeDuration}

		err = cimTime.UnmarshalJSON([]byte(`{"Days":0,"Hours":1,"Minutes":30,"Seconds":0,"Milliseconds":0}`))
		suite.NoError(err)
		suite.Equal(expectedCimTimeDuration, cimTime)
	})

	suite.Run("should unmarshal the CimInstance duration json to CimTimeDuration with all possible fields", func() {
		cimTime := CimTimeDuration{}
		expectedTimeDuration, err := time.ParseDuration("98h30m5s10ms")
		suite.Require().NoError(err)
		expectedCimTimeDuration := CimTimeDuration{Duration: expectedTimeDuration}

		err = cimTime.UnmarshalJSON([]byte(`{"Days":4,"Hours":2,"Minutes":30,"Seconds":5,"Milliseconds":10}`))
		suite.NoError(err)
		suite.Equal(expectedCimTimeDuration, cimTime)
	})

	suite.Run("should unmarshal the whole CimTimeDuration correctly", func() {
		testCimTime := testCimTime{}
		err := json.Unmarshal([]byte(suite.testJson), &testCimTime)
		suite.NoError(err)
		suite.Equal(suite.testExpected, testCimTime)
	})
}
