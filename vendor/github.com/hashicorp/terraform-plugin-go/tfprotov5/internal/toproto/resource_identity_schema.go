// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/internal/tfplugin5"
)

func ResourceIdentitySchema(in *tfprotov5.ResourceIdentitySchema) *tfplugin5.ResourceIdentitySchema {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ResourceIdentitySchema{
		Version:            in.Version,
		IdentityAttributes: ResourceIdentitySchema_IdentityAttributes(in.IdentityAttributes),
	}

	return resp
}

func ResourceIdentitySchema_IdentityAttribute(in *tfprotov5.ResourceIdentitySchemaAttribute) *tfplugin5.ResourceIdentitySchema_IdentityAttribute {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ResourceIdentitySchema_IdentityAttribute{
		Name:              in.Name,
		Type:              CtyType(in.Type),
		RequiredForImport: in.RequiredForImport,
		OptionalForImport: in.OptionalForImport,
		Description:       in.Description,
	}

	return resp
}

func ResourceIdentitySchema_IdentityAttributes(in []*tfprotov5.ResourceIdentitySchemaAttribute) []*tfplugin5.ResourceIdentitySchema_IdentityAttribute {
	if in == nil {
		return nil
	}

	resp := make([]*tfplugin5.ResourceIdentitySchema_IdentityAttribute, 0, len(in))

	for _, a := range in {
		resp = append(resp, ResourceIdentitySchema_IdentityAttribute(a))
	}

	return resp
}
