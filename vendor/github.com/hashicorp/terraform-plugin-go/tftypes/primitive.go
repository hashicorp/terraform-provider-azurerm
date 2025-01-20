// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tftypes

import (
	"fmt"
	"math/big"
)

var (
	// DynamicPseudoType is a pseudo-type in Terraform's type system that
	// is used as a wildcard type. It indicates that any Terraform type can
	// be used.
	DynamicPseudoType = primitive{name: "DynamicPseudoType"}

	// String is a primitive type in Terraform that represents a UTF-8
	// string of bytes.
	String = primitive{name: "String"}

	// Number is a primitive type in Terraform that represents a real
	// number.
	Number = primitive{name: "Number"}

	// Bool is a primitive type in Terraform that represents a true or
	// false boolean value.
	Bool = primitive{name: "Bool"}
)

var (
	_ Type = primitive{name: "test"}
)

type primitive struct {
	name string

	// used to make this type uncomparable
	// see https://golang.org/ref/spec#Comparison_operators
	// this enforces the use of Is, instead
	_ []struct{}
}

// ApplyTerraform5AttributePathStep always returns an ErrInvalidStep error
// as it is invalid to step into a primitive.
func (p primitive) ApplyTerraform5AttributePathStep(step AttributePathStep) (interface{}, error) {
	return nil, ErrInvalidStep
}

func (p primitive) Equal(o Type) bool {
	v, ok := o.(primitive)
	if !ok {
		return false
	}
	return p.name == v.name
}

func (p primitive) Is(t Type) bool {
	return p.Equal(t)
}

func (p primitive) UsableAs(t Type) bool {
	v, ok := t.(primitive)
	if !ok {
		return false
	}
	if v.name == DynamicPseudoType.name {
		return true
	}
	return v.name == p.name
}

func (p primitive) String() string {
	return "tftypes." + p.name
}

func (p primitive) private() {}

func (p primitive) MarshalJSON() ([]byte, error) {
	switch p.name {
	case String.name:
		return []byte(`"string"`), nil
	case Number.name:
		return []byte(`"number"`), nil
	case Bool.name:
		return []byte(`"bool"`), nil
	case DynamicPseudoType.name:
		return []byte(`"dynamic"`), nil
	}

	// MarshalJSON should always be error safe and reaching this panic implies
	// a new primitive type was added that needs to be handled above.
	panic(fmt.Sprintf("unimplemented tftypes.primitive type: %+v", p))
}

func (p primitive) supportedGoTypes() []string {
	switch p.name {
	case String.name:
		return []string{"string", "*string"}
	case Number.name:
		return []string{
			"*big.Float",
			"uint", "*uint",
			"uint8", "*uint8",
			"uint16", "*uint16",
			"uint32", "*uint32",
			"uint64", "*uint64",
			"int", "*int",
			"int8", "*int8",
			"int16", "*int16",
			"int32", "*int32",
			"int64", "*int64",
			"float64", "*float64",
		}
	case Bool.name:
		return []string{"bool", "*bool"}
	case DynamicPseudoType.name:
		// List/Set is covered by Tuple, Map is covered by Object
		possibleTypes := []Type{
			String, Bool, Number,
			Tuple{}, Object{},
		}
		results := []string{}
		for _, t := range possibleTypes {
			results = append(results, t.supportedGoTypes()...)
		}
		return results
	}
	panic(fmt.Sprintf("unknown primitive type %q", p.name))
}

func valueFromString(in interface{}) (Value, error) {
	switch value := in.(type) {
	case *string:
		if value == nil {
			return Value{
				typ:   String,
				value: nil,
			}, nil
		}
		return Value{
			typ:   String,
			value: *value,
		}, nil
	case string:
		return Value{
			typ:   String,
			value: value,
		}, nil
	default:
		return Value{}, fmt.Errorf("tftypes.NewValue can't use %T as a tftypes.String; expected types are: %s", in, formattedSupportedGoTypes(String))
	}
}

func valueFromBool(in interface{}) (Value, error) {
	switch value := in.(type) {
	case *bool:
		if value == nil {
			return Value{
				typ:   Bool,
				value: nil,
			}, nil
		}
		return Value{
			typ:   Bool,
			value: *value,
		}, nil
	case bool:
		return Value{
			typ:   Bool,
			value: value,
		}, nil
	default:
		return Value{}, fmt.Errorf("tftypes.NewValue can't use %T as a tftypes.Bool; expected types are: %s", in, formattedSupportedGoTypes(Bool))
	}
}

func valueFromNumber(in interface{}) (Value, error) {
	switch value := in.(type) {
	case *big.Float:
		if value == nil {
			return Value{
				typ:   Number,
				value: nil,
			}, nil
		}
		return Value{
			typ:   Number,
			value: value,
		}, nil
	case uint:
		return Value{
			typ:   Number,
			value: new(big.Float).SetUint64(uint64(value)),
		}, nil
	case *uint:
		if value == nil {
			return Value{
				typ:   Number,
				value: nil,
			}, nil
		}
		return Value{
			typ:   Number,
			value: new(big.Float).SetUint64(uint64(*value)),
		}, nil
	case uint8:
		return Value{
			typ:   Number,
			value: new(big.Float).SetUint64(uint64(value)),
		}, nil
	case *uint8:
		if value == nil {
			return Value{
				typ:   Number,
				value: nil,
			}, nil
		}
		return Value{
			typ:   Number,
			value: new(big.Float).SetUint64(uint64(*value)),
		}, nil
	case uint16:
		return Value{
			typ:   Number,
			value: new(big.Float).SetUint64(uint64(value)),
		}, nil
	case *uint16:
		if value == nil {
			return Value{
				typ:   Number,
				value: nil,
			}, nil
		}
		return Value{
			typ:   Number,
			value: new(big.Float).SetUint64(uint64(*value)),
		}, nil
	case uint32:
		return Value{
			typ:   Number,
			value: new(big.Float).SetUint64(uint64(value)),
		}, nil
	case *uint32:
		if value == nil {
			return Value{
				typ:   Number,
				value: nil,
			}, nil
		}
		return Value{
			typ:   Number,
			value: new(big.Float).SetUint64(uint64(*value)),
		}, nil
	case uint64:
		return Value{
			typ:   Number,
			value: new(big.Float).SetUint64(value),
		}, nil
	case *uint64:
		if value == nil {
			return Value{
				typ:   Number,
				value: nil,
			}, nil
		}
		return Value{
			typ:   Number,
			value: new(big.Float).SetUint64(*value),
		}, nil
	case int:
		return Value{
			typ:   Number,
			value: new(big.Float).SetInt64(int64(value)),
		}, nil
	case *int:
		if value == nil {
			return Value{
				typ:   Number,
				value: nil,
			}, nil
		}
		return Value{
			typ:   Number,
			value: new(big.Float).SetInt64(int64(*value)),
		}, nil
	case int8:
		return Value{
			typ:   Number,
			value: new(big.Float).SetInt64(int64(value)),
		}, nil
	case *int8:
		if value == nil {
			return Value{
				typ:   Number,
				value: nil,
			}, nil
		}
		return Value{
			typ:   Number,
			value: new(big.Float).SetInt64(int64(*value)),
		}, nil
	case int16:
		return Value{
			typ:   Number,
			value: new(big.Float).SetInt64(int64(value)),
		}, nil
	case *int16:
		if value == nil {
			return Value{
				typ:   Number,
				value: nil,
			}, nil
		}
		return Value{
			typ:   Number,
			value: new(big.Float).SetInt64(int64(*value)),
		}, nil
	case int32:
		return Value{
			typ:   Number,
			value: new(big.Float).SetInt64(int64(value)),
		}, nil
	case *int32:
		if value == nil {
			return Value{
				typ:   Number,
				value: nil,
			}, nil
		}
		return Value{
			typ:   Number,
			value: new(big.Float).SetInt64(int64(*value)),
		}, nil
	case int64:
		return Value{
			typ:   Number,
			value: new(big.Float).SetInt64(value),
		}, nil
	case *int64:
		if value == nil {
			return Value{
				typ:   Number,
				value: nil,
			}, nil
		}
		return Value{
			typ:   Number,
			value: new(big.Float).SetInt64(*value),
		}, nil
	case float64:
		return Value{
			typ:   Number,
			value: big.NewFloat(value),
		}, nil
	case *float64:
		if value == nil {
			return Value{
				typ:   Number,
				value: nil,
			}, nil
		}
		return Value{
			typ:   Number,
			value: big.NewFloat(*value),
		}, nil
	default:
		return Value{}, fmt.Errorf("tftypes.NewValue can't use %T as a tftypes.Number; expected types are: %s", in, formattedSupportedGoTypes(Number))
	}
}

func valueFromDynamicPseudoType(val interface{}) (Value, error) {
	switch val := val.(type) {
	case string, *string:
		v, err := valueFromString(val)
		if err != nil {
			return Value{}, err
		}
		v.typ = DynamicPseudoType
		return v, nil
	case *big.Float, float64, *float64, int, *int, int8, *int8, int16, *int16, int32, *int32, int64, *int64, uint, *uint, uint8, *uint8, uint16, *uint16, uint32, *uint32, uint64, *uint64:
		v, err := valueFromNumber(val)
		if err != nil {
			return Value{}, err
		}
		v.typ = DynamicPseudoType
		return v, nil
	case bool, *bool:
		v, err := valueFromBool(val)
		if err != nil {
			return Value{}, err
		}
		v.typ = DynamicPseudoType
		return v, nil
	case map[string]Value:
		v, err := valueFromObject(nil, nil, val)
		if err != nil {
			return Value{}, err
		}
		v.typ = DynamicPseudoType
		return v, nil
	case []Value:
		v, err := valueFromTuple(nil, val)
		if err != nil {
			return Value{}, err
		}
		v.typ = DynamicPseudoType
		return v, nil
	default:
		return Value{}, fmt.Errorf("tftypes.NewValue can't use %T as a tftypes.DynamicPseudoType; expected types are: %s", val, formattedSupportedGoTypes(DynamicPseudoType))
	}
}
