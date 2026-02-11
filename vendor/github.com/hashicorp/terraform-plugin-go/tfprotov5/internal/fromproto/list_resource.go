// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/internal/tfplugin5"
)

func ListResourceRequest(in *tfplugin5.ListResource_Request) *tfprotov5.ListResourceRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.ListResourceRequest{
		TypeName:        in.TypeName,
		Config:          DynamicValue(in.Config),
		IncludeResource: in.IncludeResourceObject,
		Limit:           in.Limit,
	}
}

func ValidateListResourceConfigRequest(in *tfplugin5.ValidateListResourceConfig_Request) *tfprotov5.ValidateListResourceConfigRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.ValidateListResourceConfigRequest{
		TypeName:              in.TypeName,
		Config:                DynamicValue(in.Config),
		IncludeResourceObject: DynamicValue(in.IncludeResourceObject),
		Limit:                 DynamicValue(in.Limit),
	}
}
