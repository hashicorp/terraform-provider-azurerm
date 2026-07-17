// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package timeouts

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/internal/validators"
)

const (
	attributeNameRead = "read"
)

// Opts is used as an argument to BlockWithOpts and AttributesWithOpts to indicate
// whether supplied descriptions should override default descriptions.
type Opts struct {
	ReadDescription string
}

// BlockWithOpts returns a schema.Block containing attributes for `Read`, which is
// defined as types.StringType and optional. A validator is used to verify
// that the value assigned to `Read` can be parsed as time.Duration. The supplied
// Opts are used to override defaults.
func BlockWithOpts(ctx context.Context, opts Opts) schema.Block {
	return schema.SingleNestedBlock{
		Attributes: attributesMap(opts),
		CustomType: Type{
			ObjectType: types.ObjectType{
				AttrTypes: attrTypesMap(),
			},
		},
	}
}

// Block returns a schema.Block containing attributes for `Read`, which is
// defined as types.StringType and optional. A validator is used to verify
// that the value assigned to `Read` can be parsed as time.Duration.
func Block(ctx context.Context) schema.Block {
	return schema.SingleNestedBlock{
		Attributes: attributesMap(Opts{}),
		CustomType: Type{
			ObjectType: types.ObjectType{
				AttrTypes: attrTypesMap(),
			},
		},
	}
}

// AttributesWithOpts returns a schema.SingleNestedAttribute which contains an
// attribute for `Read`, which is defined as types.StringType and optional.
// A validator is used to verify that the value assigned to an attribute
// can be parsed as time.Duration. The supplied Opts are used to override defaults.
func AttributesWithOpts(ctx context.Context, opts Opts) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes: attributesMap(opts),
		CustomType: Type{
			ObjectType: types.ObjectType{
				AttrTypes: attrTypesMap(),
			},
		},
		Optional: true,
	}
}

// Attributes returns a schema.SingleNestedAttribute which contains an
// attribute for `Read`, which is defined as types.StringType and optional.
// A validator is used to verify that the value assigned to an attribute
// can be parsed as time.Duration.
func Attributes(ctx context.Context) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes: attributesMap(Opts{}),
		CustomType: Type{
			ObjectType: types.ObjectType{
				AttrTypes: attrTypesMap(),
			},
		},
		Optional: true,
	}
}

func attributesMap(opts Opts) map[string]schema.Attribute {
	attribute := schema.StringAttribute{
		Optional: true,
		Description: `A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) ` +
			`consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are ` +
			`"s" (seconds), "m" (minutes), "h" (hours).`,
		Validators: []validator.String{
			validators.TimeDuration(),
		},
	}

	if opts.ReadDescription != "" {
		attribute.Description = opts.ReadDescription
	}

	return map[string]schema.Attribute{
		attributeNameRead: attribute,
	}
}

func attrTypesMap() map[string]attr.Type {
	return map[string]attr.Type{
		attributeNameRead: types.StringType,
	}
}
