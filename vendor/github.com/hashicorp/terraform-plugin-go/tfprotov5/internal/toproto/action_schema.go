// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/internal/tfplugin5"
)

func ActionSchema(in *tfprotov5.ActionSchema) *tfplugin5.ActionSchema {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ActionSchema{
		Schema: Schema(in.Schema),
	}

	return resp
}
