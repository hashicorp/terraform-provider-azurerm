// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/internal/tfplugin6"
)

func DynamicValue(in *tfplugin6.DynamicValue) *tfprotov6.DynamicValue {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.DynamicValue{
		MsgPack: in.Msgpack,
		JSON:    in.Json,
	}

	return resp
}
