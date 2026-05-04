// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/internal/tfplugin5"
)

func Deferred(in *tfprotov5.Deferred) *tfplugin5.Deferred {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.Deferred{
		Reason: tfplugin5.Deferred_Reason(in.Reason),
	}

	return resp
}
