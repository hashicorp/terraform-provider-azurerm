// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tftypes

import (
	"bytes"
	"fmt"
)

// Set is a Terraform type representing an unordered collection of unique
// elements, all of the same type.
type Set struct {
	ElementType Type

	// used to make this type uncomparable
	// see https://golang.org/ref/spec#Comparison_operators
	// this enforces the use of Is, instead
	_ []struct{}
}

// ApplyTerraform5AttributePathStep applies an AttributePathStep to a Set,
// returning the Type found at that AttributePath within the Set. If the
// AttributePathStep cannot be applied to the Set, an ErrInvalidStep error
// will be returned.
func (s Set) ApplyTerraform5AttributePathStep(step AttributePathStep) (interface{}, error) {
	switch step.(type) {
	case ElementKeyValue:
		return s.ElementType, nil
	default:
		return nil, ErrInvalidStep
	}
}

// Equal returns true if the two Sets are exactly equal. Unlike Is, passing in
// a Set with no ElementType will always return false.
func (s Set) Equal(o Type) bool {
	v, ok := o.(Set)
	if !ok {
		return false
	}
	if v.ElementType == nil || s.ElementType == nil {
		// when doing exact comparisons, we can't compare types that
		// don't have element types set, so we just consider them not
		// equal
		return false
	}
	return s.ElementType.Equal(v.ElementType)
}

// UsableAs returns whether the two Sets are type compatible.
//
// If the other type is DynamicPseudoType, it will return true.
// If the other type is not a Set, it will return false.
// If the other Set does not have a type compatible ElementType, it will
// return false.
func (s Set) UsableAs(o Type) bool {
	if o.Is(DynamicPseudoType) {
		return true
	}
	v, ok := o.(Set)
	if !ok {
		return false
	}
	return s.ElementType.UsableAs(v.ElementType)
}

// Is returns whether `t` is a Set type or not. It does not perform any
// ElementType checks.
func (s Set) Is(t Type) bool {
	_, ok := t.(Set)
	return ok
}

func (s Set) String() string {
	return "tftypes.Set[" + s.ElementType.String() + "]"
}

func (s Set) private() {}

func (s Set) supportedGoTypes() []string {
	return []string{"[]tftypes.Value"}
}

func valueFromSet(typ Type, in interface{}) (Value, error) {
	switch value := in.(type) {
	case []Value:
		var elType Type
		for _, v := range value {
			if !v.Type().UsableAs(typ) {
				return Value{}, NewAttributePath().WithElementKeyValue(v).NewErrorf("can't use %s as %s", v.Type(), typ)
			}
			if elType == nil {
				elType = v.Type()
			}
			if !elType.Equal(v.Type()) {
				return Value{}, fmt.Errorf("sets must only contain one type of element, saw %s and %s", elType, v.Type())
			}
		}
		return Value{
			typ:   Set{ElementType: typ},
			value: value,
		}, nil
	default:
		return Value{}, fmt.Errorf("tftypes.NewValue can't use %T as a tftypes.Set; expected types are: %s", in, formattedSupportedGoTypes(Set{}))
	}
}

// MarshalJSON returns a JSON representation of the full type signature of `s`,
// including its ElementType.
//
// Deprecated: this is not meant to be called by third-party code.
func (s Set) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(`["set",`)

	// MarshalJSON is always error safe
	elementTypeBytes, _ := s.ElementType.MarshalJSON()

	buf.Write(elementTypeBytes)
	buf.WriteString(`]`)

	return buf.Bytes(), nil
}
