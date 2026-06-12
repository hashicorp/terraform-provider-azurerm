// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/internal/tfplugin5"
)

func ResourceIdentityData(in *tfprotov5.ResourceIdentityData) *tfplugin5.ResourceIdentityData {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ResourceIdentityData{
		IdentityData: DynamicValue(in.IdentityData),
	}

	return resp
}
