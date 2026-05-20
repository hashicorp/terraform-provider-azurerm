// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/internal/tfplugin6"
)

func ValidateActionConfigRequest(in *tfplugin6.ValidateActionConfig_Request) *tfprotov6.ValidateActionConfigRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ValidateActionConfigRequest{
		ActionType: in.ActionType,
		Config:     DynamicValue(in.Config),
	}
}

func PlanActionRequest(in *tfplugin6.PlanAction_Request) *tfprotov6.PlanActionRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.PlanActionRequest{
		ActionType:         in.ActionType,
		Config:             DynamicValue(in.Config),
		ClientCapabilities: PlanActionClientCapabilities(in.ClientCapabilities),
	}

	return resp
}

func InvokeActionRequest(in *tfplugin6.InvokeAction_Request) *tfprotov6.InvokeActionRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.InvokeActionRequest{
		ActionType:         in.ActionType,
		Config:             DynamicValue(in.Config),
		ClientCapabilities: InvokeActionClientCapabilities(in.ClientCapabilities),
	}

	return resp
}
