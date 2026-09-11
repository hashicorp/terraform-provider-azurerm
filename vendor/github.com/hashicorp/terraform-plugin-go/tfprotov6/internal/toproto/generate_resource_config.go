// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/internal/tfplugin6"
)

func GenerateResourceConfig_Response(in *tfprotov6.GenerateResourceConfigResponse) *tfplugin6.GenerateResourceConfig_Response {
	if in == nil {
		return nil
	}

	return &tfplugin6.GenerateResourceConfig_Response{
		Config:      DynamicValue(in.Config),
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}
