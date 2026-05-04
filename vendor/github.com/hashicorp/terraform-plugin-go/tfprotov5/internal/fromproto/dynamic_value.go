// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/internal/tfplugin5"
)

func DynamicValue(in *tfplugin5.DynamicValue) *tfprotov5.DynamicValue {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.DynamicValue{
		MsgPack: in.Msgpack,
		JSON:    in.Json,
	}

	return resp
}
