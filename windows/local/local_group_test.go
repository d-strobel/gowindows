package local

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupRead(t *testing.T) {
	client := &Client{}

	// Test with empty parameters should fail
	t.Run("EmptyParams", func(t *testing.T) {
		params := GroupParams{}
		group, err := client.GroupRead(context.Background(), params)

		assert.Error(t, err, "Error should not be nil")
		assert.Nil(t, group, "Group should be nil")

		expectedError := "Name or SID must be set"
		assert.EqualError(t, err, expectedError, "Error message should match")
	})
}

func TestGroupCreate(t *testing.T) {
	client := &Client{}

	// Test with empty parameters should fail
	t.Run("EmptyParams", func(t *testing.T) {
		params := GroupParams{}
		group, err := client.GroupCreate(context.Background(), params)

		assert.Error(t, err, "Error should not be nil")
		assert.Nil(t, group, "Group should be nil")

		expectedError := "Name must be set"
		assert.EqualError(t, err, expectedError, "Error message should match")
	})
}

func TestGroupUpdate(t *testing.T) {
	client := &Client{}

	// Test with empty parameters should fail
	t.Run("EmptyParams", func(t *testing.T) {
		params := GroupParams{}
		group, err := client.GroupUpdate(context.Background(), params)

		assert.Error(t, err, "Error should not be nil")
		assert.Nil(t, group, "Group should be nil")

		expectedError := "Name or SID must be set to change a group"
		assert.EqualError(t, err, expectedError, "Error message should match")
	})

	// Test without description parameters should fail
	t.Run("WithoutDescription", func(t *testing.T) {
		params := GroupParams{
			Name: "User",
		}
		group, err := client.GroupUpdate(context.Background(), params)

		assert.Error(t, err, "Error should not be nil")
		assert.Nil(t, group, "Group should be nil")

		expectedError := "Description must be set"
		assert.EqualError(t, err, expectedError, "Error message should match")
	})
}

func TestGroupDelete(t *testing.T) {
	client := &Client{}

	// Test with empty parameters should fail
	t.Run("EmptyParams", func(t *testing.T) {
		params := GroupParams{}
		err := client.GroupDelete(context.Background(), params)

		assert.Error(t, err, "Error should not be nil")

		expectedError := "Name or SID must be set to delete a group"
		assert.EqualError(t, err, expectedError, "Error message should match")
	})
}
