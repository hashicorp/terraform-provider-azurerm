// Copyright (c) HashiCorp, Inc.
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
