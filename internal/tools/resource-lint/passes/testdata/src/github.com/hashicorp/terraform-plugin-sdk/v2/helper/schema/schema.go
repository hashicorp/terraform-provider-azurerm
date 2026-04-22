// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

type ValueType int

const (
	TypeBool ValueType = iota
	TypeInt
	TypeFloat
	TypeString
	TypeList
	TypeMap
	TypeSet
)

type SchemaValidateFunc func(interface{}, string) ([]string, []error)

type Schema struct {
	Type          ValueType
	Required      bool
	Optional      bool
	Computed      bool
	ForceNew      bool
	ValidateFunc  SchemaValidateFunc
	MaxItems      int
	Elem          interface{}
	AtLeastOneOf  []string
}

type Resource struct {
	Schema map[string]*Schema
}

type ResourceData struct{}

func (d *ResourceData) HasChange(key string) bool { return false }

func (d *ResourceData) HasChanges(keys ...string) bool { return false }

func (d *ResourceData) GetOk(key string) (interface{}, bool) { return nil, false }

func (d *ResourceData) Get(key string) interface{} { return nil }
