// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package reflect

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
)

// trueReflectValue returns the reflect.Value for `in` after derefencing all
// the pointers and unwrapping all the interfaces. It's the concrete value
// beneath it all.
func trueReflectValue(val reflect.Value) reflect.Value {
	kind := val.Type().Kind()
	for kind == reflect.Interface || kind == reflect.Ptr {
		innerVal := val.Elem()
		if !innerVal.IsValid() {
			break
		}
		val = innerVal
		kind = val.Type().Kind()
	}
	return val
}

// commaSeparatedString returns an English joining of the strings in `in`,
// using "and" and commas as appropriate.
func commaSeparatedString(in []string) string {
	switch len(in) {
	case 0:
		return ""
	case 1:
		return in[0]
	case 2:
		return strings.Join(in, " and ")
	default:
		in[len(in)-1] = "and " + in[len(in)-1]
		return strings.Join(in, ", ")
	}
}

// getStructTags returns a map of Terraform field names to their position in
// the fields of the struct `in`. `in` must be a struct.
//
// The position of the field in a struct is represented as an index sequence to support type embedding
// in structs. This index sequence can be used to retrieve the field with the Go "reflect" package FieldByIndex methods:
//   - https://pkg.go.dev/reflect#Type.FieldByIndex
//   - https://pkg.go.dev/reflect#Value.FieldByIndex
//   - https://pkg.go.dev/reflect#Value.FieldByIndexErr
//
// The following are not supported and will return an error if detected in a struct (including embedded structs):
//   - Duplicate "tfsdk" tags
//   - Exported fields without a "tfsdk" tag
//   - Exported fields with an invalid "tfsdk" tag (must be a valid Terraform identifier)
func getStructTags(ctx context.Context, typ reflect.Type, path path.Path) (map[string][]int, error) { //nolint:unparam // False positive, ctx is used below.
	tags := make(map[string][]int, 0)

	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%s: can't get struct tags of %s, is not a struct", path, typ)
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() && !field.Anonymous {
			// Skip unexported fields. Unexported embedded structs (anonymous fields) are allowed because they may
			// contain exported fields that are promoted; which means they can be read/set.
			continue
		}

		// This index sequence is the location of the field within the struct.
		// For promoted fields from an embedded struct, the length of this sequence will be > 1
		fieldIndexSequence := []int{i}
		tag, tagExists := field.Tag.Lookup(`tfsdk`)

		// "tfsdk" tags with "-" are being explicitly excluded
		if tag == "-" {
			continue
		}

		// Handle embedded structs
		if field.Anonymous {
			if tagExists {
				return nil, fmt.Errorf(`%s: embedded struct field %s cannot have tfsdk tag`, path.AtName(tag), field.Name)
			}

			embeddedTags, err := getStructTags(ctx, field.Type, path)
			if err != nil {
				return nil, fmt.Errorf(`error retrieving embedded struct %q field tags: %w`, field.Name, err)
			}
			for k, v := range embeddedTags {
				if other, ok := tags[k]; ok {
					otherField := typ.FieldByIndex(other)
					return nil, fmt.Errorf("embedded struct %q promotes a field with a duplicate tfsdk tag %q, conflicts with %q tfsdk tag", field.Name, k, otherField.Name)
				}

				tags[k] = append(fieldIndexSequence, v...)
			}
			continue
		}

		// All non-embedded fields must have a tfsdk tag
		if !tagExists {
			return nil, fmt.Errorf(`%s: need a struct tag for "tfsdk" on %s`, path, field.Name)
		}

		// Ensure the tfsdk tag has a valid name
		path := path.AtName(tag)
		if !isValidFieldName(tag) {
			return nil, fmt.Errorf("%s: invalid tfsdk tag, must only use lowercase letters, underscores, and numbers, and must start with a letter", path)
		}

		// Ensure there are no duplicate tfsdk tags
		if other, ok := tags[tag]; ok {
			otherField := typ.FieldByIndex(other)
			return nil, fmt.Errorf("%s: can't use tfsdk tag %q for both %s and %s fields", path, tag, otherField.Name, field.Name)
		}

		tags[tag] = fieldIndexSequence
	}

	return tags, nil
}

// isValidFieldName returns true if `name` can be used as a field name in a
// Terraform resource or data source.
func isValidFieldName(name string) bool {
	re := regexp.MustCompile("^[a-z][a-z0-9_]*$")
	return re.MatchString(name)
}

// canBeNil returns true if `target`'s type can hold a nil value
func canBeNil(target reflect.Value) bool {
	switch target.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Interface:
		// these types can all hold nils
		return true
	default:
		// nothing else can be set to nil
		return false
	}
}
