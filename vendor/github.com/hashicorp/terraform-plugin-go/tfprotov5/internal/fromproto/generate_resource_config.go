// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/internal/tfplugin5"
)

func GenerateResourceConfigRequest(in *tfplugin5.GenerateResourceConfig_Request) *tfprotov5.GenerateResourceConfigRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.GenerateResourceConfigRequest{
		TypeName: in.TypeName,
		State:    DynamicValue(in.State),
	}
}
