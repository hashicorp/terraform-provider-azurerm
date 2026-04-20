// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/internal/tfplugin5"
)

func ValidateActionConfigRequest(in *tfplugin5.ValidateActionConfig_Request) *tfprotov5.ValidateActionConfigRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.ValidateActionConfigRequest{
		ActionType: in.ActionType,
		Config:     DynamicValue(in.Config),
	}
}

func PlanActionRequest(in *tfplugin5.PlanAction_Request) *tfprotov5.PlanActionRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.PlanActionRequest{
		ActionType:         in.ActionType,
		Config:             DynamicValue(in.Config),
		ClientCapabilities: PlanActionClientCapabilities(in.ClientCapabilities),
	}

	return resp
}

func InvokeActionRequest(in *tfplugin5.InvokeAction_Request) *tfprotov5.InvokeActionRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.InvokeActionRequest{
		ActionType:         in.ActionType,
		Config:             DynamicValue(in.Config),
		ClientCapabilities: InvokeActionClientCapabilities(in.ClientCapabilities),
	}

	return resp
}
