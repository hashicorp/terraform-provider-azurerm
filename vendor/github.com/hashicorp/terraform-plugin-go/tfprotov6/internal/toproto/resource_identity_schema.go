// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/internal/tfplugin6"
)

func ResourceIdentitySchema(in *tfprotov6.ResourceIdentitySchema) *tfplugin6.ResourceIdentitySchema {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.ResourceIdentitySchema{
		Version:            in.Version,
		IdentityAttributes: ResourceIdentitySchema_IdentityAttributes(in.IdentityAttributes),
	}

	return resp
}

func ResourceIdentitySchema_IdentityAttribute(in *tfprotov6.ResourceIdentitySchemaAttribute) *tfplugin6.ResourceIdentitySchema_IdentityAttribute {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.ResourceIdentitySchema_IdentityAttribute{
		Name:              in.Name,
		Type:              CtyType(in.Type),
		RequiredForImport: in.RequiredForImport,
		OptionalForImport: in.OptionalForImport,
		Description:       in.Description,
	}

	return resp
}

func ResourceIdentitySchema_IdentityAttributes(in []*tfprotov6.ResourceIdentitySchemaAttribute) []*tfplugin6.ResourceIdentitySchema_IdentityAttribute {
	if in == nil {
		return nil
	}

	resp := make([]*tfplugin6.ResourceIdentitySchema_IdentityAttribute, 0, len(in))

	for _, a := range in {
		resp = append(resp, ResourceIdentitySchema_IdentityAttribute(a))
	}

	return resp
}
