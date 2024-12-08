package parsing

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

// Unit test suite for all Time parsing functions
type CimClassKeyValUnitTestSuite struct {
	suite.Suite
	// Fixtures
	testJson      string
	testExpected  TestCimClass
	testJson2     string
	testExpected2 TestCimClass
}

// Fixture objects
type TestCimClass struct {
	CimClass TestCimClassValues `json:"CimClass"`
}
type TestCimClassValues struct {
	CimSuperClassName  string         `json:"CimSuperClassName"`
	CimClassQualifiers CimClassKeyVal `json:"CimClassQualifiers"`
}

// Setup test suite for CimClassKeyVal
func (suite *CimClassKeyValUnitTestSuite) SetupSuite() {
	// Fixtures
	suite.testJson = `{
        "CimClass":  {
            "CimSuperClassName":  "DnsDomain",
		    "CimSuperClass":  "ROOT/Microsoft/Windows/DNS:DnsDomain",
		    "CimClassProperties":  "DistinguishedName HostName RecordClass RecordData RecordType Timestamp TimeToLive Type",
		    "CimClassQualifiers":  "dynamic = True provider = \"DnsServerPSProvider\" ClassVersion = \"1.0.0\" locale = 1033",
		    "CimClassMethods":  "",
		    "CimSystemProperties":  "Microsoft.Management.Infrastructure.CimSystemProperties"
        }
    }`
	suite.testExpected = TestCimClass{
		TestCimClassValues{
			CimSuperClassName: "DnsDomain",
			CimClassQualifiers: CimClassKeyVal{
				"dynamic":      "True",
				"provider":     "DnsServerPSProvider",
				"ClassVersion": "1.0.0",
				"locale":       "1033",
			},
		},
	}
	suite.testJson2 = `{
        "CimClass":  {
            "CimSuperClassName":  "DnsDomain",
		    "CimSuperClass":  "ROOT/Microsoft/Windows/DNS:DnsDomain",
		    "CimClassProperties":  "DistinguishedName HostName RecordClass RecordData RecordType Timestamp TimeToLive Type",
		    "CimClassQualifiers":  ["HostNameAlias = \"test.local.\""],
		    "CimClassMethods":  "",
		    "CimSystemProperties":  "Microsoft.Management.Infrastructure.CimSystemProperties"
        }
    }`
	suite.testExpected2 = TestCimClass{
		TestCimClassValues{
			CimSuperClassName: "DnsDomain",
			CimClassQualifiers: CimClassKeyVal{
				"HostNameAlias": "test.local.",
			},
		},
	}
}

func TestCimClassKeyValUnitTestSuite(t *testing.T) {
	suite.Run(t, &CimClassKeyValUnitTestSuite{})
}

func (suite *CimClassKeyValUnitTestSuite) TestUnmarshalJSON() {
	suite.T().Parallel()

	suite.Run("should unmarshal the whole CimClass correctly", func() {
		testCimClass := TestCimClass{}
		err := json.Unmarshal([]byte(suite.testJson), &testCimClass)
		suite.NoError(err)
		suite.Equal(suite.testExpected, testCimClass)
	})

	suite.Run("should unmarshal the whole CimClass correctly with json array", func() {
		testCimClass := TestCimClass{}
		err := json.Unmarshal([]byte(suite.testJson2), &testCimClass)
		suite.NoError(err)
		suite.Equal(suite.testExpected2, testCimClass)
	})

	suite.Run("should unmarshal the CimClassQualifiers a key-value map", func() {
		cimClassKeyVal := CimClassKeyVal{}
		err := cimClassKeyVal.UnmarshalJSON([]byte(`"dynamic = True provider = \"DnsServerPSProvider\" ClassVersion = \"1.0.0\" locale = 1033"`))
		suite.NoError(err)
		suite.Equal(CimClassKeyVal{
			"dynamic":      "True",
			"provider":     "DnsServerPSProvider",
			"ClassVersion": "1.0.0",
			"locale":       "1033",
		}, cimClassKeyVal)
	})

	suite.Run("should unmarshal to an empty key-value map", func() {
		cimClassKeyVal := CimClassKeyVal{}
		err := cimClassKeyVal.UnmarshalJSON([]byte(`""`))
		suite.NoError(err)
		suite.Equal(CimClassKeyVal{}, cimClassKeyVal)
	})

	suite.Run("should unmarshal single field unquoted correctly", func() {
		cimClassKeyVal := CimClassKeyVal{}
		err := cimClassKeyVal.UnmarshalJSON([]byte(`"dynamic = True"`))
		suite.NoError(err)
		suite.Equal(CimClassKeyVal{"dynamic": "True"}, cimClassKeyVal)
	})

	suite.Run("should unmarshal single field quoted correctly", func() {
		cimClassKeyVal := CimClassKeyVal{}
		err := cimClassKeyVal.UnmarshalJSON([]byte(`"dynamic = \"True\""`))
		suite.NoError(err)
		suite.Equal(CimClassKeyVal{"dynamic": "True"}, cimClassKeyVal)
	})

	suite.Run("should unmarshal the json array to key-value map", func() {
		cimClassKeyVal := CimClassKeyVal{}
		err := cimClassKeyVal.UnmarshalJSON([]byte(`["HostNameAlias = \"test.local.\""],`))
		suite.NoError(err)
		suite.Equal(CimClassKeyVal{
			"HostNameAlias": "test.local.",
		}, cimClassKeyVal)
	})
}
