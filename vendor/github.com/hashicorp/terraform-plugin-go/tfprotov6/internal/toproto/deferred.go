// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/internal/tfplugin6"
)

func Deferred(in *tfprotov6.Deferred) *tfplugin6.Deferred {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.Deferred{
		Reason: tfplugin6.Deferred_Reason(in.Reason),
	}

	return resp
}
