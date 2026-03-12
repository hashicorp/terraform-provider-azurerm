// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package timeouts

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/internal/validators"
)

const (
	attributeNameCreate = "create"
	attributeNameRead   = "read"
	attributeNameUpdate = "update"
	attributeNameDelete = "delete"
)

// Opts is used as an argument to Block and Attributes to indicate which attributes
// should be created and whether supplied descriptions should override default
// descriptions.
type Opts struct {
	Create            bool
	Read              bool
	Update            bool
	Delete            bool
	CreateDescription string
	ReadDescription   string
	UpdateDescription string
	DeleteDescription string
}

// Block returns a schema.Block containing attributes for each of the fields
// in Opts which are set to true. Each attribute is defined as types.StringType
// and optional. A validator is used to verify that the value assigned to an
// attribute can be parsed as time.Duration.
func Block(ctx context.Context, opts Opts) schema.Block {
	return schema.SingleNestedBlock{
		Attributes: attributesMap(opts),
		CustomType: Type{
			ObjectType: types.ObjectType{
				AttrTypes: attrTypesMap(opts),
			},
		},
	}
}

// BlockAll returns a schema.Block containing attributes for each of create, read,
// update and delete. Each attribute is defined as types.StringType and optional.
// A validator is used to verify that the value assigned to an attribute can be
// parsed as time.Duration.
func BlockAll(ctx context.Context) schema.Block {
	return Block(ctx, Opts{
		Create: true,
		Read:   true,
		Update: true,
		Delete: true,
	})
}

// Attributes returns a schema.SingleNestedAttribute which contains attributes for
// each of the fields in Opts which are set to true. Each attribute is defined as
// types.StringType and optional. A validator is used to verify that the value
// assigned to an attribute can be parsed as time.Duration.
func Attributes(ctx context.Context, opts Opts) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes: attributesMap(opts),
		CustomType: Type{
			ObjectType: types.ObjectType{
				AttrTypes: attrTypesMap(opts),
			},
		},
		Optional: true,
	}
}

// AttributesAll returns a schema.SingleNestedAttribute which contains attributes
// for each of create, read, update and delete. Each attribute is defined as
// types.StringType and optional. A validator is used to verify that the value
// assigned to an attribute can be parsed as time.Duration.
func AttributesAll(ctx context.Context) schema.Attribute {
	return Attributes(ctx, Opts{
		Create: true,
		Read:   true,
		Update: true,
		Delete: true,
	})
}

func attributesMap(opts Opts) map[string]schema.Attribute {
	description := `A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) ` +
		`consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are ` +
		`"s" (seconds), "m" (minutes), "h" (hours).`
	attributes := map[string]schema.Attribute{}
	attribute := schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			validators.TimeDuration(),
		},
	}

	if opts.Create {
		attribute.Description = description

		if opts.CreateDescription != "" {
			attribute.Description = opts.CreateDescription
		}

		attributes[attributeNameCreate] = attribute
	}

	if opts.Read {
		attribute.Description = description + ` Read operations occur during any refresh or planning operation ` +
			`when refresh is enabled.`

		if opts.ReadDescription != "" {
			attribute.Description = opts.ReadDescription
		}

		attributes[attributeNameRead] = attribute
	}

	if opts.Update {
		attribute.Description = description

		if opts.UpdateDescription != "" {
			attribute.Description = opts.UpdateDescription
		}

		attributes[attributeNameUpdate] = attribute
	}

	if opts.Delete {
		attribute.Description = description + ` Setting a timeout for a Delete operation is only applicable if ` +
			`changes are saved into state before the destroy operation occurs.`

		if opts.DeleteDescription != "" {
			attribute.Description = opts.DeleteDescription
		}

		attributes[attributeNameDelete] = attribute
	}

	return attributes
}

func attrTypesMap(opts Opts) map[string]attr.Type {
	attrTypes := map[string]attr.Type{}

	if opts.Create {
		attrTypes[attributeNameCreate] = types.StringType
	}

	if opts.Read {
		attrTypes[attributeNameRead] = types.StringType
	}

	if opts.Update {
		attrTypes[attributeNameUpdate] = types.StringType
	}

	if opts.Delete {
		attrTypes[attributeNameDelete] = types.StringType
	}

	return attrTypes
}
