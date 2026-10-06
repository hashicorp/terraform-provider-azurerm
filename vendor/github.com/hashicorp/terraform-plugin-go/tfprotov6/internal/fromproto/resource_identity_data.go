// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/internal/tfplugin6"
)

func ResourceIdentityData(in *tfplugin6.ResourceIdentityData) *tfprotov6.ResourceIdentityData {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.ResourceIdentityData{
		IdentityData: DynamicValue(in.IdentityData),
	}

	return resp
}
