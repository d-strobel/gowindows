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
		assert.IsType(t, Group{}, group)
		assert.ErrorContains(t, err, "GroupRead: group parameter 'Name' or 'SID' must be set")
	})
}

func TestGroupCreate(t *testing.T) {
	client := &Client{}

	// Test with empty parameters should fail
	t.Run("EmptyParams", func(t *testing.T) {
		params := GroupParams{}
		group, err := client.GroupCreate(context.Background(), params)

		assert.Error(t, err, "Error should not be nil")
		assert.IsType(t, Group{}, group)
		assert.ErrorContains(t, err, "GroupCreate: group parameter 'Name' must be set")
	})
}

func TestGroupUpdate(t *testing.T) {
	client := &Client{}

	// Test with empty parameters should fail
	t.Run("EmptyParams", func(t *testing.T) {
		params := GroupParams{}
		group, err := client.GroupUpdate(context.Background(), params)

		assert.Error(t, err, "Error should not be nil")
		assert.IsType(t, Group{}, group)
		assert.ErrorContains(t, err, "GroupUpdate: group parameter 'Name' or 'SID' must be set")
	})

	// Test without description parameters should fail
	t.Run("WithoutDescription", func(t *testing.T) {
		params := GroupParams{
			Name: "User",
		}
		group, err := client.GroupUpdate(context.Background(), params)

		assert.Error(t, err, "Error should not be nil")
		assert.IsType(t, Group{}, group)
		assert.ErrorContains(t, err, "GroupUpdate: group parameter 'Description' must be set")
	})
}

func TestGroupDelete(t *testing.T) {
	client := &Client{}

	// Test with empty parameters should fail
	t.Run("EmptyParams", func(t *testing.T) {
		params := GroupParams{}
		err := client.GroupDelete(context.Background(), params)

		assert.Error(t, err, "Error should not be nil")
		assert.ErrorContains(t, err, "GroupDelete: group parameter 'Name' or 'SID' must be set")
	})
}
