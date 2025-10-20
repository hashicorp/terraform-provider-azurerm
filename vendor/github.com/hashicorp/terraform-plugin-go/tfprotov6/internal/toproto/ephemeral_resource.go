// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/internal/tfplugin6"
)

func GetMetadata_EphemeralResourceMetadata(in *tfprotov6.EphemeralResourceMetadata) *tfplugin6.GetMetadata_EphemeralResourceMetadata {
	if in == nil {
		return nil
	}

	return &tfplugin6.GetMetadata_EphemeralResourceMetadata{
		TypeName: in.TypeName,
	}
}

func ValidateEphemeralResourceConfig_Response(in *tfprotov6.ValidateEphemeralResourceConfigResponse) *tfplugin6.ValidateEphemeralResourceConfig_Response {
	if in == nil {
		return nil
	}

	return &tfplugin6.ValidateEphemeralResourceConfig_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}

func OpenEphemeralResource_Response(in *tfprotov6.OpenEphemeralResourceResponse) *tfplugin6.OpenEphemeralResource_Response {
	if in == nil {
		return nil
	}

	return &tfplugin6.OpenEphemeralResource_Response{
		Result:      DynamicValue(in.Result),
		Diagnostics: Diagnostics(in.Diagnostics),
		Private:     in.Private,
		RenewAt:     Timestamp(in.RenewAt),
		Deferred:    Deferred(in.Deferred),
	}
}

func RenewEphemeralResource_Response(in *tfprotov6.RenewEphemeralResourceResponse) *tfplugin6.RenewEphemeralResource_Response {
	if in == nil {
		return nil
	}

	return &tfplugin6.RenewEphemeralResource_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
		Private:     in.Private,
		RenewAt:     Timestamp(in.RenewAt),
	}
}

func CloseEphemeralResource_Response(in *tfprotov6.CloseEphemeralResourceResponse) *tfplugin6.CloseEphemeralResource_Response {
	if in == nil {
		return nil
	}

	return &tfplugin6.CloseEphemeralResource_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}
