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
		LinkedResources:    ProposedLinkedResources(in.LinkedResources),
		Config:             DynamicValue(in.Config),
		ClientCapabilities: PlanActionClientCapabilities(in.ClientCapabilities),
	}

	return resp
}

func ProposedLinkedResources(in []*tfplugin6.PlanAction_Request_LinkedResource) []*tfprotov6.ProposedLinkedResource {
	resp := make([]*tfprotov6.ProposedLinkedResource, 0, len(in))

	for _, inLinkedResource := range in {
		resp = append(resp, &tfprotov6.ProposedLinkedResource{
			PriorState:    DynamicValue(inLinkedResource.PriorState),
			PlannedState:  DynamicValue(inLinkedResource.PlannedState),
			Config:        DynamicValue(inLinkedResource.Config),
			PriorIdentity: ResourceIdentityData(inLinkedResource.PriorIdentity),
		})
	}

	return resp
}

func InvokeActionRequest(in *tfplugin6.InvokeAction_Request) *tfprotov6.InvokeActionRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.InvokeActionRequest{
		ActionType:      in.ActionType,
		LinkedResources: InvokeLinkedResources(in.LinkedResources),
		Config:          DynamicValue(in.Config),
	}

	return resp
}

func InvokeLinkedResources(in []*tfplugin6.InvokeAction_Request_LinkedResource) []*tfprotov6.InvokeLinkedResource {
	resp := make([]*tfprotov6.InvokeLinkedResource, 0, len(in))

	for _, inLinkedResource := range in {
		resp = append(resp, &tfprotov6.InvokeLinkedResource{
			PriorState:      DynamicValue(inLinkedResource.PriorState),
			PlannedState:    DynamicValue(inLinkedResource.PlannedState),
			Config:          DynamicValue(inLinkedResource.Config),
			PlannedIdentity: ResourceIdentityData(inLinkedResource.PlannedIdentity),
		})
	}

	return resp
}
