// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/internal/tfplugin5"
)

func GenerateResourceConfig_Response(in *tfprotov5.GenerateResourceConfigResponse) *tfplugin5.GenerateResourceConfig_Response {
	if in == nil {
		return nil
	}

	return &tfplugin5.GenerateResourceConfig_Response{
		Config:      DynamicValue(in.Config),
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}
