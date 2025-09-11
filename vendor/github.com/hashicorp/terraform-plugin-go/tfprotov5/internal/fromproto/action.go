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
		LinkedResources:    ProposedLinkedResources(in.LinkedResources),
		Config:             DynamicValue(in.Config),
		ClientCapabilities: PlanActionClientCapabilities(in.ClientCapabilities),
	}

	return resp
}

func ProposedLinkedResources(in []*tfplugin5.PlanAction_Request_LinkedResource) []*tfprotov5.ProposedLinkedResource {
	resp := make([]*tfprotov5.ProposedLinkedResource, 0, len(in))

	for _, inLinkedResource := range in {
		resp = append(resp, &tfprotov5.ProposedLinkedResource{
			PriorState:    DynamicValue(inLinkedResource.PriorState),
			PlannedState:  DynamicValue(inLinkedResource.PlannedState),
			Config:        DynamicValue(inLinkedResource.Config),
			PriorIdentity: ResourceIdentityData(inLinkedResource.PriorIdentity),
		})
	}

	return resp
}

func InvokeActionRequest(in *tfplugin5.InvokeAction_Request) *tfprotov5.InvokeActionRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.InvokeActionRequest{
		ActionType:      in.ActionType,
		LinkedResources: InvokeLinkedResources(in.LinkedResources),
		Config:          DynamicValue(in.Config),
	}

	return resp
}

func InvokeLinkedResources(in []*tfplugin5.InvokeAction_Request_LinkedResource) []*tfprotov5.InvokeLinkedResource {
	resp := make([]*tfprotov5.InvokeLinkedResource, 0, len(in))

	for _, inLinkedResource := range in {
		resp = append(resp, &tfprotov5.InvokeLinkedResource{
			PriorState:      DynamicValue(inLinkedResource.PriorState),
			PlannedState:    DynamicValue(inLinkedResource.PlannedState),
			Config:          DynamicValue(inLinkedResource.Config),
			PlannedIdentity: ResourceIdentityData(inLinkedResource.PlannedIdentity),
		})
	}

	return resp
}
