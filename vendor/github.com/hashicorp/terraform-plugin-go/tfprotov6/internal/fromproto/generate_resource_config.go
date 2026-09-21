// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/internal/tfplugin6"
)

func GenerateResourceConfigRequest(in *tfplugin6.GenerateResourceConfig_Request) *tfprotov6.GenerateResourceConfigRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.GenerateResourceConfigRequest{
		TypeName: in.TypeName,
		State:    DynamicValue(in.State),
	}
}
