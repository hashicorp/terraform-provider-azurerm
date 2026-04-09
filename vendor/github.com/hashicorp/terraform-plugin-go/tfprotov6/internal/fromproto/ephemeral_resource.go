// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/internal/tfplugin6"
)

func ValidateEphemeralResourceConfigRequest(in *tfplugin6.ValidateEphemeralResourceConfig_Request) *tfprotov6.ValidateEphemeralResourceConfigRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ValidateEphemeralResourceConfigRequest{
		TypeName: in.TypeName,
		Config:   DynamicValue(in.Config),
	}
}

func OpenEphemeralResourceRequest(in *tfplugin6.OpenEphemeralResource_Request) *tfprotov6.OpenEphemeralResourceRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.OpenEphemeralResourceRequest{
		TypeName:           in.TypeName,
		Config:             DynamicValue(in.Config),
		ClientCapabilities: OpenEphemeralResourceClientCapabilities(in.ClientCapabilities),
	}
}

func RenewEphemeralResourceRequest(in *tfplugin6.RenewEphemeralResource_Request) *tfprotov6.RenewEphemeralResourceRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.RenewEphemeralResourceRequest{
		TypeName: in.TypeName,
		Private:  in.Private,
	}
}

func CloseEphemeralResourceRequest(in *tfplugin6.CloseEphemeralResource_Request) *tfprotov6.CloseEphemeralResourceRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.CloseEphemeralResourceRequest{
		TypeName: in.TypeName,
		Private:  in.Private,
	}
}
