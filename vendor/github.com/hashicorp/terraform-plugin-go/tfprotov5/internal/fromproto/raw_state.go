// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/internal/tfplugin5"
)

func RawState(in *tfplugin5.RawState) *tfprotov5.RawState {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.RawState{
		JSON:    in.Json,
		Flatmap: in.Flatmap,
	}

	return resp
}
