// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto5

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// IdentitySchema converts a *tfprotov5.ResourceIdentitySchema into a resource/identityschema Schema, used for
// converting protocol identity schemas (from another provider server, such as SDKv2 or terraform-plugin-go)
// into Framework identity schemas.
func IdentitySchema(ctx context.Context, s *tfprotov5.ResourceIdentitySchema) (*identityschema.Schema, error) {
	if s == nil {
		return nil, nil
	}

	attrs, err := IdentitySchemaAttributes(ctx, s.IdentityAttributes)
	if err != nil {
		return nil, err
	}

	return &identityschema.Schema{
		// MAINTAINER NOTE: At the moment, there isn't a need to copy all of the data from the protocol identity schema
		// to the resource identity schema, just enough data to allow provider developers to read and set data.
		Attributes: attrs,
	}, nil
}

func IdentitySchemaAttributes(ctx context.Context, protoAttrs []*tfprotov5.ResourceIdentitySchemaAttribute) (map[string]identityschema.Attribute, error) {
	attrs := make(map[string]identityschema.Attribute, len(protoAttrs))
	for _, protoAttr := range protoAttrs {
		// MAINTAINER NOTE: At the moment, there isn't a need to copy all of the data from the protocol identity schema
		// to the resource identity schema, just enough data to allow provider developers to read and set data.
		switch {
		case protoAttr.Type.Is(tftypes.Bool):
			attrs[protoAttr.Name] = identityschema.BoolAttribute{
				RequiredForImport: protoAttr.RequiredForImport,
				OptionalForImport: protoAttr.OptionalForImport,
			}
		case protoAttr.Type.Is(tftypes.Number):
			attrs[protoAttr.Name] = identityschema.NumberAttribute{
				RequiredForImport: protoAttr.RequiredForImport,
				OptionalForImport: protoAttr.OptionalForImport,
			}
		case protoAttr.Type.Is(tftypes.String):
			attrs[protoAttr.Name] = identityschema.StringAttribute{
				RequiredForImport: protoAttr.RequiredForImport,
				OptionalForImport: protoAttr.OptionalForImport,
			}
		case protoAttr.Type.Is(tftypes.List{}):
			//nolint:forcetypeassert // Type assertion is guaranteed by the above `(tftypes.Type).Is` function
			l := protoAttr.Type.(tftypes.List)

			elementType, err := basetypes.TerraformTypeToFrameworkType(l.ElementType)
			if err != nil {
				return nil, err
			}

			attrs[protoAttr.Name] = identityschema.ListAttribute{
				ElementType:       elementType,
				RequiredForImport: protoAttr.RequiredForImport,
				OptionalForImport: protoAttr.OptionalForImport,
			}
		default:
			// MAINTAINER NOTE: Not all terraform types are valid identity attribute types. Framework fully supports
			// all of the possible identity attribute types, so any errors here would be invalid protocol identities.
			return nil, fmt.Errorf("no supported identity attribute for %q, type: %T", protoAttr.Name, protoAttr.Type)
		}
	}

	return attrs, nil
}
