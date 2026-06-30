// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package schema

type ValueType int

const (
	TypeInvalid ValueType = iota
	TypeBool
	TypeInt
	TypeFloat
	TypeString
	TypeList
	TypeMap
	TypeSet
)

type SchemaValidateFunc func(interface{}, string) ([]string, []error)

type Resource struct {
	Schema map[string]*Schema
}

type ResourceData struct{}

func (d *ResourceData) Get(string) interface{} {
	return nil
}

func (d *ResourceData) HasChange(string) bool {
	return false
}

func (d *ResourceData) HasChanges(...string) bool {
	return false
}

type Schema struct {
	Type          ValueType
	Required      bool
	Optional      bool
	Computed      bool
	ForceNew      bool
	Sensitive     bool
	Default       interface{}
	Description   string
	MaxItems      int
	Elem          interface{}
	ValidateFunc  SchemaValidateFunc
	AtLeastOneOf  []string
	ExactlyOneOf  []string
	ConflictsWith []string
}
