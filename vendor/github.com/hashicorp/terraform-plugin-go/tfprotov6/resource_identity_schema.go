// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfprotov6

import "github.com/hashicorp/terraform-plugin-go/tftypes"

// ResourceIdentitySchema is the identity schema for a Resource.
type ResourceIdentitySchema struct {
	// Version indicates which version of the schema this is. Versions
	// should be monotonically incrementing numbers. When Terraform
	// encounters a resource identity stored in state with a schema version
	// lower that the identity schema version the provider advertises for
	// that resource, Terraform requests the provider upgrade the resource's
	// identity state.
	Version int64

	// IdentityAttributes is a list of attributes that uniquely identify a
	// resource. These attributes are used to identify a resource in the
	// state and to import existing resources into the state.
	IdentityAttributes []*ResourceIdentitySchemaAttribute
}

// ValueType returns the tftypes.Type for a ResourceIdentitySchema.
//
// If ResourceIdentitySchema is missing, an empty Object is returned.
func (s *ResourceIdentitySchema) ValueType() tftypes.Type {
	if s == nil {
		return tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{},
		}
	}

	attributeTypes := map[string]tftypes.Type{}

	for _, attribute := range s.IdentityAttributes {
		if attribute == nil {
			continue
		}

		attributeType := attribute.ValueType()

		if attributeType == nil {
			continue
		}

		attributeTypes[attribute.Name] = attributeType
	}

	return tftypes.Object{
		AttributeTypes: attributeTypes,
	}
}

// ResourceIdentitySchemaAttribute represents one value of data within
// resource identity.
// These are always used in resource identity comparisons.
type ResourceIdentitySchemaAttribute struct {
	// Name is the name of the attribute. This is what the user will put
	// before the equals sign to assign a value to this attribute during import.
	Name string

	// Type indicates the type of data the attribute expects. See the
	// documentation for the tftypes package for information on what types
	// are supported and their behaviors.
	// For resource identity Terraform core only supports the following types:
	// - bool
	// - number
	// - string
	// - list of bool
	// - list of number
	// - list of string
	Type tftypes.Type

	// RequiredForImport indicates whether this attribute is required to
	// import the resource. For example it might be false if the value
	// can be derived from provider configuration. Either this or OptionalForImport
	// needs to be true.
	RequiredForImport bool

	// OptionalForImport indicates whether this attribute is optional to
	// import the resource. For example it might be true if the value
	// can be derived from provider configuration. Either this or RequiredForImport
	// needs to be true.
	OptionalForImport bool

	// Description is a human-readable description of the attribute.
	Description string
}

// ValueType returns the tftypes.Type for a ResourceIdentitySchemaAttribute.
//
// If ResourceIdentitySchemaAttribute is missing, nil is returned.
func (s *ResourceIdentitySchemaAttribute) ValueType() tftypes.Type {
	if s == nil {
		return nil
	}

	return s.Type
}
