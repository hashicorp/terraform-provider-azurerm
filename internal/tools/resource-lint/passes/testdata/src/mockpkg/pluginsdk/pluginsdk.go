// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package pluginsdk

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// This file is to mock pluginsdk in azurerm

const (
	TypeString = schema.TypeString
	TypeBool   = schema.TypeBool
	TypeInt    = schema.TypeInt
	TypeMap    = schema.TypeMap
	TypeList   = schema.TypeList
)

type (
	Resource     = schema.Resource
	Schema       = schema.Schema
	ResourceData = schema.ResourceData
)

// StateChangeConf mocks the pluginsdk.StateChangeConf type
type StateChangeConf struct {
	Pending      []string
	Target       []string
	Refresh      StateRefreshFunc
	Timeout      time.Duration
	PollInterval time.Duration
}

// StateRefreshFunc is a function that refreshes the state
type StateRefreshFunc func() (interface{}, string, error)

// WaitForStateContext waits for the state to reach a target
func (conf *StateChangeConf) WaitForStateContext(ctx context.Context) (interface{}, error) {
	return nil, nil
}

// GetWriteOnly retrieves a write-only attribute value from the ResourceData.
func GetWriteOnly(d *ResourceData, key string, valType interface{}) (interface{}, error) {
	return d.Get(key), nil
}
