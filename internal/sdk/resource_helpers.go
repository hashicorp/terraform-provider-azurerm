// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"fmt"
	"reflect"
	"strings"
)

type decodedStructTags struct {
	// hclPath defines the path to this field used for this in the Schema for this Resource
	hclPath string

	// addedInNextMajorVersion specifies whether this field should only be introduced in a next major
	// version of the Provider
	addedInNextMajorVersion bool

	// removedInNextMajorVersion specifies whether this field is deprecated and should not
	// be set into the state in the next major version of the Provider
	removedInNextMajorVersion bool
}

// parseStructTags parses the struct tags defined in input into a decodedStructTags object
// which allows for the consistent parsing of struct tags across the Typed SDK.
func parseStructTags(input reflect.StructTag) (*decodedStructTags, error) {
	tag, ok := input.Lookup("tfschema")
	if !ok {
		// doesn't exist - ignore it?
		return nil, nil
	}
	if tag == "" {
		return nil, fmt.Errorf("the `tfschema` struct tag was defined but empty")
	}

	components := strings.Split(tag, ",")
	output := &decodedStructTags{
		// NOTE: `hclPath` has to be the first item in the struct tag
		hclPath:                   strings.TrimSpace(components[0]),
		addedInNextMajorVersion:   false,
		removedInNextMajorVersion: false,
	}
	if output.hclPath == "" {
		return nil, fmt.Errorf("hclPath was empty")
	}

	if len(components) > 1 {
		// remove the hcl field name since it's been parsed
		components = components[1:]
		for _, item := range components {
			item = strings.TrimSpace(item) // allowing for both `foo,bar` and `foo, bar` in struct tags
			if strings.EqualFold(item, "removedInNextMajorVersion") {
				if output.addedInNextMajorVersion {
					return nil, fmt.Errorf("the struct-tags `removedInNextMajorVersion` and `addedInNextMajorVersion` cannot be set together")
				}
				output.removedInNextMajorVersion = true
				continue
			}
			if strings.EqualFold(item, "addedInNextMajorVersion") {
				if output.removedInNextMajorVersion {
					return nil, fmt.Errorf("the struct-tags `removedInNextMajorVersion` and `addedInNextMajorVersion` cannot be set together")
				}
				output.addedInNextMajorVersion = true
				continue
			}

			return nil, fmt.Errorf("internal-error: the struct-tag %q is not implemented - struct tags are %q", item, tag)
		}
	}

	return output, nil
}

// PluginSDKFlattenStructSliceToInterface converts a slice of structs with tfschema tags to []interface{}.
// This is specifically needed for list resources that use pluginsdk.ResourceData.Set() instead of
// metadata.Encode() from the typed SDK. The pluginsdk API cannot handle custom struct types directly
// and requires conversion to primitive types (maps, slices of primitives, etc.).
//
// WHEN TO USE:
//   - List resources that use pluginsdk.ResourceData.Set() for setting state
//   - Any resource still using the untyped/plugin SDK instead of the typed SDK
//
// WHEN NOT TO USE:
//   - Resources using the typed SDK with metadata.Encode() - those handle encoding automatically
//   - The typed SDK's internal encoding (see resource_encode.go) is preferred for typed resources
//
// Parameters:
//   - input: A slice of structs with tfschema tags
//   - debugLogger: Optional function for debug logging. Pass nil to disable logging.
//
// Returns a slice of map[string]interface{} suitable for pluginsdk.ResourceData.Set()
func PluginSDKFlattenStructSliceToInterface(input interface{}, debugLogger func(string, ...interface{})) []interface{} {
	val := reflect.ValueOf(input)
	if val.Kind() != reflect.Slice {
		if debugLogger != nil {
			debugLogger("PluginSDKFlattenStructSliceToInterface: input is not a slice, returning empty slice")
		}
		return []interface{}{}
	}

	result := make([]interface{}, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		item := pluginSDKFlattenStructToMap(val.Index(i), debugLogger)
		if item != nil {
			result = append(result, item)
		}
	}

	if debugLogger != nil {
		debugLogger("PluginSDKFlattenStructSliceToInterface: converted %d items", len(result))
	}
	return result
}

// pluginSDKFlattenStructToMap converts a struct with tfschema tags to map[string]interface{}.
// This is an internal helper for PluginSDKFlattenStructSliceToInterface.
func pluginSDKFlattenStructToMap(val reflect.Value, debugLogger func(string, ...interface{})) map[string]interface{} {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			if debugLogger != nil {
				debugLogger("pluginSDKFlattenStructToMap: nil pointer encountered")
			}
			return nil
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		if debugLogger != nil {
			debugLogger("pluginSDKFlattenStructToMap: value is not a struct, kind=%v", val.Kind())
		}
		return nil
	}

	result := make(map[string]interface{})
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		tag := field.Tag.Get("tfschema")
		if tag == "" {
			continue
		}

		result[tag] = pluginSDKFlattenValue(fieldVal, debugLogger)
	}

	if debugLogger != nil {
		debugLogger("pluginSDKFlattenStructToMap: converted struct with %d fields", len(result))
	}
	return result
}

// pluginSDKFlattenValue recursively flattens a reflect.Value to a type compatible with pluginsdk.ResourceData.Set().
// This is an internal helper for the PluginSDK flatten functions.
func pluginSDKFlattenValue(val reflect.Value, debugLogger func(string, ...interface{})) interface{} {
	switch val.Kind() {
	case reflect.String:
		return val.String()
	case reflect.Int, reflect.Int64:
		return val.Int()
	case reflect.Float64:
		return val.Float()
	case reflect.Bool:
		return val.Bool()
	case reflect.Slice:
		// Primitive slices can be returned as-is
		if val.Type().Elem().Kind() == reflect.String {
			return val.Interface()
		}
		if val.Type().Elem().Kind() == reflect.Int64 {
			return val.Interface()
		}
		// Struct slices need recursive conversion
		return PluginSDKFlattenStructSliceToInterface(val.Interface(), debugLogger)
	case reflect.Map:
		return val.Interface()
	case reflect.Ptr:
		if val.IsNil() {
			return nil
		}
		return pluginSDKFlattenValue(val.Elem(), debugLogger)
	case reflect.Struct:
		return pluginSDKFlattenStructToMap(val, debugLogger)
	default:
		if debugLogger != nil {
			debugLogger("pluginSDKFlattenValue: unhandled kind=%v, returning as-is", val.Kind())
		}
		return val.Interface()
	}
}
