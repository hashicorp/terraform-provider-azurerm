// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/internal/tfplugin5"
)

func ValidateEphemeralResourceConfigRequest(in *tfplugin5.ValidateEphemeralResourceConfig_Request) *tfprotov5.ValidateEphemeralResourceConfigRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.ValidateEphemeralResourceConfigRequest{
		TypeName: in.TypeName,
		Config:   DynamicValue(in.Config),
	}
}

func OpenEphemeralResourceRequest(in *tfplugin5.OpenEphemeralResource_Request) *tfprotov5.OpenEphemeralResourceRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.OpenEphemeralResourceRequest{
		TypeName:           in.TypeName,
		Config:             DynamicValue(in.Config),
		ClientCapabilities: OpenEphemeralResourceClientCapabilities(in.ClientCapabilities),
	}
}

func RenewEphemeralResourceRequest(in *tfplugin5.RenewEphemeralResource_Request) *tfprotov5.RenewEphemeralResourceRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.RenewEphemeralResourceRequest{
		TypeName: in.TypeName,
		Private:  in.Private,
	}
}

func CloseEphemeralResourceRequest(in *tfplugin5.CloseEphemeralResource_Request) *tfprotov5.CloseEphemeralResourceRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.CloseEphemeralResourceRequest{
		TypeName: in.TypeName,
		Private:  in.Private,
	}
}
