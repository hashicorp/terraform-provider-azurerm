// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

// This file begins to define the abstraction layer for allowing the provider to be migrated from PluginSDK to the
// Plugin Framework.

type ResourceData interface {
	// Get returns a value from either the config/state depending on where this is called
	// in Create and Update functions this will return from the config
	// in Read, Exists and Import functions this will return from the state
	// NOTE: this should not be called from Delete functions.
	Get(key string) interface{}

	GetOk(key string) (interface{}, bool)

	GetOkExists(key string) (interface{}, bool)

	// GetFromConfig(key string) interface{} // TODO - Should we support this?
	//
	// GetFromState(key string) interface{} // TODO - Should we support this?

	GetChange(key string) (interface{}, interface{})

	HasChange(key string) bool

	HasChanges(keys ...string) bool

	HasChangesExcept(keys ...string) bool

	Id() string

	Set(key string, value interface{}) error

	SetConnInfo(v map[string]string)

	SetId(id string)

	// TODO: add Get/Set helpers for each type?
	// e.g. GetString(key string) *string
}
