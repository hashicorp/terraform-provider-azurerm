// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/internal/tfplugin5"
)

func GetMetadata_ListResourceMetadata(in *tfprotov5.ListResourceMetadata) *tfplugin5.GetMetadata_ListResourceMetadata {
	if in == nil {
		return nil
	}

	return &tfplugin5.GetMetadata_ListResourceMetadata{
		TypeName: in.TypeName,
	}
}

func ListResource_ListResourceEvent(in *tfprotov5.ListResourceResult) *tfplugin5.ListResource_Event {
	return &tfplugin5.ListResource_Event{
		DisplayName:    in.DisplayName,
		ResourceObject: DynamicValue(in.Resource),
		Identity:       ResourceIdentityData(in.Identity),
		Diagnostic:     Diagnostics(in.Diagnostics),
	}
}

func ValidateListResourceConfig_Response(in *tfprotov5.ValidateListResourceConfigResponse) *tfplugin5.ValidateListResourceConfig_Response {
	if in == nil {
		return nil
	}

	return &tfplugin5.ValidateListResourceConfig_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}
