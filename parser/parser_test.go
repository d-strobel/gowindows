package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewParser(t *testing.T) {
	expectedResult := &Parser{}
	actualResult := NewParser()

	assert.Equal(t, expectedResult, actualResult)
}
