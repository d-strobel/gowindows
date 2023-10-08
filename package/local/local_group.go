package local

import (
	"context"
	"encoding/json"
	"fmt"
)

type SID struct {
	Value string `json:"Value"`
}

type Group struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Sid         SID    `json:"SID"`
}

var g Group

// GetLocalGroupByName will return a Group by a given name.
// If the group is not present, it will return an error.
func (c *Client) GetGroupByName(ctx context.Context, name string) (*Group, error) {

	cmd := fmt.Sprintf("Get-LocalGroup -Name \"%s\" | ConvertTo-Json", name)

	result, err := c.Connection.Run(ctx, cmd)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(result), &g)
	if err != nil {
		return nil, err
	}

	return &g, nil
}

// GetLocalGroupBySID will return a Group by a given SID.
// If the group is not present, it will return an error.
func (c *Client) GetGroupBySID(ctx context.Context, sid string) (*Group, error) {

	cmd := fmt.Sprintf("Get-LocalGroup -SID \"%s\" | ConvertTo-Json", sid)

	result, err := c.Connection.Run(ctx, cmd)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(result), &g)
	if err != nil {
		return nil, err
	}

	return &g, nil
}
