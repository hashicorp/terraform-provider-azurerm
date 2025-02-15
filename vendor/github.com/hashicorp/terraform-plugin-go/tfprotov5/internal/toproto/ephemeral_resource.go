// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/internal/tfplugin5"
)

func GetMetadata_EphemeralResourceMetadata(in *tfprotov5.EphemeralResourceMetadata) *tfplugin5.GetMetadata_EphemeralResourceMetadata {
	if in == nil {
		return nil
	}

	return &tfplugin5.GetMetadata_EphemeralResourceMetadata{
		TypeName: in.TypeName,
	}
}

func ValidateEphemeralResourceConfig_Response(in *tfprotov5.ValidateEphemeralResourceConfigResponse) *tfplugin5.ValidateEphemeralResourceConfig_Response {
	if in == nil {
		return nil
	}

	return &tfplugin5.ValidateEphemeralResourceConfig_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}

func OpenEphemeralResource_Response(in *tfprotov5.OpenEphemeralResourceResponse) *tfplugin5.OpenEphemeralResource_Response {
	if in == nil {
		return nil
	}

	return &tfplugin5.OpenEphemeralResource_Response{
		Result:      DynamicValue(in.Result),
		Diagnostics: Diagnostics(in.Diagnostics),
		Private:     in.Private,
		RenewAt:     Timestamp(in.RenewAt),
		Deferred:    Deferred(in.Deferred),
	}
}

func RenewEphemeralResource_Response(in *tfprotov5.RenewEphemeralResourceResponse) *tfplugin5.RenewEphemeralResource_Response {
	if in == nil {
		return nil
	}

	return &tfplugin5.RenewEphemeralResource_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
		Private:     in.Private,
		RenewAt:     Timestamp(in.RenewAt),
	}
}

func CloseEphemeralResource_Response(in *tfprotov5.CloseEphemeralResourceResponse) *tfplugin5.CloseEphemeralResource_Response {
	if in == nil {
		return nil
	}

	return &tfplugin5.CloseEphemeralResource_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}
