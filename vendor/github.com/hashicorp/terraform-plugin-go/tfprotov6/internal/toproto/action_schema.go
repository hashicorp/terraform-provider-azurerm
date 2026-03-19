// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/internal/tfplugin6"
)

func ActionSchema(in *tfprotov6.ActionSchema) *tfplugin6.ActionSchema {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.ActionSchema{
		Schema: Schema(in.Schema),
	}
	return resp
}
