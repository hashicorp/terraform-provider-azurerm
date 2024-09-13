// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

const autoTFVarsJson = "terraform-plugin-testing.auto.tfvars.json"

// Variable interface is an alias to json.Marshaler.
type Variable interface {
	json.Marshaler
}

// Variables is a type holding a key-value map of variable names
// to types implementing the Variable interface.
type Variables map[string]Variable

// Write creates a file in the destination supplied
// containing JSON encoded Variables.
func (v Variables) Write(dest string) error {
	if len(v) == 0 {
		return nil
	}

	b, err := json.Marshal(v)

	if err != nil {
		return fmt.Errorf("cannot marshal variables: %s", err)
	}

	outFilename := filepath.Join(dest, autoTFVarsJson)

	err = os.WriteFile(outFilename, b, 0600)

	if err != nil {
		return fmt.Errorf("cannot write variables file: %s", err)
	}

	return nil
}

var _ Variable = boolVariable{}

// boolVariable supports JSON encoding of a bool.
type boolVariable struct {
	value bool
}

// MarshalJSON returns the JSON encoding of boolVariable.
func (v boolVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

// BoolVariable returns boolVariable which implements Variable.
func BoolVariable(value bool) boolVariable {
	return boolVariable{
		value: value,
	}
}

var _ Variable = floatVariable{}

// floatVariable supports JSON encoding of any floating-point type.
type floatVariable struct {
	value any
}

// MarshalJSON returns the JSON encoding of floatVariable.
func (v floatVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

// FloatVariable returns floatVariable which implements Variable.
func FloatVariable[T anyFloat](value T) floatVariable {
	return floatVariable{
		value: value,
	}
}

var _ Variable = integerVariable{}

// integerVariable supports JSON encoding of any integer type.
type integerVariable struct {
	value any
}

// MarshalJSON returns the JSON encoding of integerVariable.
func (v integerVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

// IntegerVariable returns integerVariable which implements Variable.
func IntegerVariable[T anyInteger](value T) integerVariable {
	return integerVariable{
		value: value,
	}
}

var _ Variable = listVariable{}

// listVariable supports JSON encoding of slice of Variable.
type listVariable struct {
	value []Variable
}

// MarshalJSON returns the JSON encoding of listVariable.
// Every Variable within a listVariable must be the same
// underlying type.
func (v listVariable) MarshalJSON() ([]byte, error) {
	if !typesEq(v.value) {
		return nil, errors.New("lists must contain the same type")
	}

	return json.Marshal(v.value)
}

// ListVariable returns listVariable which implements Variable.
func ListVariable(value ...Variable) listVariable {
	return listVariable{
		value: value,
	}
}

var _ Variable = mapVariable{}

// mapVariable supports JSON encoding of a key-value map of
// string to Variable.
type mapVariable struct {
	value map[string]Variable
}

// MarshalJSON returns the JSON encoding of mapVariable.
// Every Variable in a mapVariable must be the same
// underlying type.
func (v mapVariable) MarshalJSON() ([]byte, error) {
	var variables []Variable

	for _, variable := range v.value {
		variables = append(variables, variable)
	}

	if !typesEq(variables) {
		return nil, errors.New("maps must contain the same type")
	}

	return json.Marshal(v.value)
}

// MapVariable returns mapVariable which implements Variable.
func MapVariable(value map[string]Variable) mapVariable {
	return mapVariable{
		value: value,
	}
}

var _ Variable = objectVariable{}

// objectVariable supports JSON encoding of a key-value
// map of string to Variable in which each Variable
// can be a different underlying type.
type objectVariable struct {
	value map[string]Variable
}

// MarshalJSON returns the JSON encoding of objectVariable.
func (v objectVariable) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(v.value)

	if err != nil {
		innerErr := err

		// Unwrap is used here to expose the initial error, for example
		// "maps must contain the same type" whilst removing any errors
		// related to the implementation (i.e., the usage of
		// encoding/json in this instance.
		for errors.Unwrap(innerErr) != nil {
			innerErr = errors.Unwrap(err)
		}

		return nil, innerErr
	}

	return b, nil
}

// ObjectVariable returns objectVariable which implements Variable.
func ObjectVariable(value map[string]Variable) objectVariable {
	return objectVariable{
		value: value,
	}
}

var _ Variable = setVariable{}

// setVariable supports JSON encoding of a slice of Variable.
type setVariable struct {
	value []Variable
}

// MarshalJSON returns the JSON encoding of setVariable.
// Every Variable in a setVariable must be the same
// underlying type.
func (v setVariable) MarshalJSON() ([]byte, error) {
	for kx, x := range v.value {
		for ky := kx + 1; ky < len(v.value); ky++ {
			y := v.value[ky]

			if _, ok := x.(setVariable); !ok {
				continue
			}

			if _, ok := y.(setVariable); !ok {
				continue
			}

			if reflect.DeepEqual(x, y) {
				return nil, errors.New("sets must contain unique elements")
			}
		}
	}

	if !typesEq(v.value) {
		return nil, errors.New("sets must contain the same type")
	}

	return json.Marshal(v.value)
}

// SetVariable returns setVariable which implements Variable.
func SetVariable(value ...Variable) setVariable {
	return setVariable{
		value: value,
	}
}

var _ Variable = stringVariable{}

// stringVariable supports JSON encoding of a string.
type stringVariable struct {
	value string
}

// MarshalJSON returns the JSON encoding of stringVariable.
func (v stringVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

// StringVariable returns stringVariable which implements Variable.
func StringVariable(value string) stringVariable {
	return stringVariable{
		value: value,
	}
}

var _ Variable = tupleVariable{}

// tupleVariable supports JSON encoding of a slice of Variable
// in which each element in the slice can be a different
// underlying type.
type tupleVariable struct {
	value []Variable
}

// MarshalJSON returns the JSON encoding of tupleVariable.
func (v tupleVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

// TupleVariable returns tupleVariable which implements Variable.
func TupleVariable(value ...Variable) tupleVariable {
	return tupleVariable{
		value: value,
	}
}

// typesEq verifies that every element in the supplied slice of Variable
// is the same underlying type.
func typesEq(variables []Variable) bool {
	var t reflect.Type

	for _, variable := range variables {
		switch x := variable.(type) {
		case listVariable:
			if !typesEq(x.value) {
				return false
			}
		case mapVariable:
			var vars []Variable

			for _, v := range x.value {
				vars = append(vars, v)
			}

			if !typesEq(vars) {
				return false
			}
		case setVariable:
			if !typesEq(x.value) {
				return false
			}
		}

		typeOfVariable := reflect.TypeOf(variable)

		if t == nil {
			t = typeOfVariable
			continue
		}

		if t != typeOfVariable {
			return false
		}
	}

	return true
}
