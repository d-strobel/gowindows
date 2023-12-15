package parser

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// Unit test suite for all CLIXML parsing functions
type CLIXMLUnitTestSuite struct {
	suite.Suite
	// Fixtures
	cliXMLError               string
	expectedString            string
	expectedUnmarshaledCLIXML *clixml
	expectedStringSlice       []string
}

func (suite *CLIXMLUnitTestSuite) SetupTest() {
	// Fixtures
	suite.cliXMLError = `#< CLIXML
	<Objs Version="1.1.0.1" xmlns="http://schemas.microsoft.com/powershell/2004/04"><Obj S="progress" RefId="0">
    <TN RefId="0"><T>System.Management.Automation.PSCustomObject</T><T>System.Object</T></TN><MS>
    <I64 N="SourceId">1</I64><PR N="Record"><AV>Loading Active Directory module for Windows PowerShell with default drive 'AD:'</AV>
	<AI>0</AI><Nil /><PI>-1</PI><PC>0</PC><T>Processing</T><SR>-1</SR><SD> </SD></PR></MS></Obj><Obj S="progress" RefId="1">
	<TNRef RefId="0" /><MS><I64 N="SourceId">1</I64><PR N="Record"><AV>Loading Active Directory module for Windows PowerShell with default drive 'AD:'</AV>
	<AI>0</AI><Nil /><PI>-1</PI><PC>25</PC><T>Processing</T><SR>-1</SR><SD> </SD></PR></MS></Obj><Obj S="progress" RefId="2">
	<TNRef RefId="0" /><MS><I64 N="SourceId">1</I64><PR N="Record"><AV>Loading Active Directory module for Windows PowerShell with default drive 'AD:'</AV>
	<AI>0</AI><Nil /><PI>-1</PI><PC>50</PC><T>Processing</T><SR>-1</SR><SD> </SD></PR></MS></Obj><Obj S="progress" RefId="3">
	<TNRef RefId="0" /><MS><I64 N="SourceId">1</I64><PR N="Record"><AV>Loading Active Directory module for Windows PowerShell with default drive 'AD:'</AV>
	<AI>0</AI><Nil /><PI>-1</PI><PC>75</PC><T>Processing</T><SR>-1</SR><SD> </SD></PR></MS></Obj><Obj S="progress" RefId="4">
	<TNRef RefId="0" /><MS><I64 N="SourceId">1</I64><PR N="Record"><AV>Loading Active Directory module for Windows PowerShell with default drive 'AD:'</AV>
	<AI>0</AI><Nil /><PI>-1</PI><PC>100</PC><T>Processing</T><SR>-1</SR><SD> </SD></PR></MS></Obj><Obj S="progress" RefId="5">
	<TNRef RefId="0" /><MS><I64 N="SourceId">1</I64><PR N="Record"><AV>Loading Active Directory module for Windows PowerShell with default drive 'AD:'</AV>
	<AI>0</AI><Nil /><PI>-1</PI><PC>100</PC><T>Completed</T><SR>-1</SR><SD> </SD></PR></MS></Obj>
	<S S="Error">Set-ADOrganizationalUnit : A parameter cannot be found that matches parameter _x000D__x000A_</S>
	<S S="Error">name 'Path'._x000D__x000A_</S><S S="Error">At line:1 char:101_x000D__x000A_</S>
	<S S="Error">+ ... e description" -Path "DC=yourdomain,DC=com" _x000D__x000A_</S>
	<S S="Error">-ProtectedFromAccidentalDeletion $tr ..._x000D__x000A_</S>
	<S S="Error">+                    ~~~~~_x000D__x000A_</S>
	<S S="Error">    + CategoryInfo          : InvalidArgument: (:) [Set-ADOrganizationalUnit], _x000D__x000A_</S
	><S S="Error">    ParameterBindingException_x000D__x000A_</S>
	<S S="Error">    + FullyQualifiedErrorId : NamedParameterNotFound,Microsoft.ActiveDirectory _x000D__x000A_</S>
	<S S="Error">   .Management.Commands.SetADOrganizationalUnit_x000D__x000A_</S><S S="Error"> _x000D__x000A_</S>
	</Objs>`

	suite.expectedString = `Set-ADOrganizationalUnit : A parameter cannot be found that matches parameter name 'Path'.At line:1 char:101
... e description" -Path "DC=yourdomain,DC=com" -ProtectedFromAccidentalDeletion $tr ...
                   ~~~~~
CategoryInfo          : InvalidArgument: (:) [Set-ADOrganizationalUnit], ParameterBindingException
FullyQualifiedErrorId : NamedParameterNotFound,Microsoft.ActiveDirectory .Management.Commands.SetADOrganizationalUnit`

	suite.expectedUnmarshaledCLIXML = &clixml{
		XML: []string{
			"Set-ADOrganizationalUnit : A parameter cannot be found that matches parameter _x000D__x000A_",
			"name 'Path'._x000D__x000A_",
			"At line:1 char:101_x000D__x000A_",
			"+ ... e description\" -Path \"DC=yourdomain,DC=com\" _x000D__x000A_",
			"-ProtectedFromAccidentalDeletion $tr ..._x000D__x000A_",
			"+                    ~~~~~_x000D__x000A_",
			"    + CategoryInfo          : InvalidArgument: (:) [Set-ADOrganizationalUnit], _x000D__x000A_",
			"    ParameterBindingException_x000D__x000A_",
			"    + FullyQualifiedErrorId : NamedParameterNotFound,Microsoft.ActiveDirectory _x000D__x000A_",
			"   .Management.Commands.SetADOrganizationalUnit_x000D__x000A_",
			" _x000D__x000A_",
		},
	}

	suite.expectedStringSlice = []string{
		"Set-ADOrganizationalUnit : A parameter cannot be found that matches parameter ",
		"name 'Path'.",
		"At line:1 char:101",
		"\n... e description\" -Path \"DC=yourdomain,DC=com\" ",
		"-ProtectedFromAccidentalDeletion $tr ...",
		"\n                   ~~~~~",
		"\nCategoryInfo          : InvalidArgument: (:) [Set-ADOrganizationalUnit], ",
		"ParameterBindingException",
		"\nFullyQualifiedErrorId : NamedParameterNotFound,Microsoft.ActiveDirectory ",
		".Management.Commands.SetADOrganizationalUnit",
		"",
	}
}

func TestCLIXMLUnitTestSuite(t *testing.T) {
	suite.Run(t, &CLIXMLUnitTestSuite{})
}

func (suite *CLIXMLUnitTestSuite) TestUnmarshal() {
	suite.T().Parallel()

	suite.Run("should unmarshal correctly", func() {
		actualResult := &clixml{}
		err := actualResult.unmarshal(suite.cliXMLError)
		suite.Require().NoError(err)
		suite.Equal(suite.expectedUnmarshaledCLIXML, actualResult)
	})
	suite.Run("should return error when empty string", func() {
		actualResult := &clixml{}
		err := actualResult.unmarshal("")
		suite.Error(err)
	})
}

func (suite *CLIXMLUnitTestSuite) TestStringSlice() {
	suite.T().Parallel()

	suite.Run("should return the expected string slice", func() {
		actualResult := suite.expectedUnmarshaledCLIXML.stringSlice()
		suite.Equal(suite.expectedStringSlice, actualResult)
	})
	suite.Run("should not panic with empty slice", func() {
		clixml := &clixml{
			XML: []string{},
		}
		actualResult := clixml.stringSlice()
		suite.Require().NotPanics(func() { clixml.stringSlice() })
		suite.Equal([]string{}, actualResult)
	})
	suite.Run("should not panic with empty string inside slice", func() {
		clixml := &clixml{
			XML: []string{"", ""},
		}
		actualResult := clixml.stringSlice()
		suite.Require().NotPanics(func() { clixml.stringSlice() })
		suite.Equal([]string{"", ""}, actualResult)
	})
}

func (suite *CLIXMLUnitTestSuite) TestDecodeCLIXML() {
	suite.T().Parallel()

	suite.Run("should return expected result", func() {
		actualResult, err := DecodeCLIXML(suite.cliXMLError)
		suite.Require().NoError(err)
		suite.Equal(suite.expectedString, actualResult)
	})
	suite.Run("should return error if not a clixml string", func() {
		actualResult, err := DecodeCLIXML("")
		suite.Error(err)
		suite.Equal("", actualResult)
	})
}
