// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/internal/tfplugin6"
)

func GetMetadata_ActionMetadata(in *tfprotov6.ActionMetadata) *tfplugin6.GetMetadata_ActionMetadata {
	if in == nil {
		return nil
	}

	return &tfplugin6.GetMetadata_ActionMetadata{
		TypeName: in.TypeName,
	}
}

func ValidateActionConfig_Response(in *tfprotov6.ValidateActionConfigResponse) *tfplugin6.ValidateActionConfig_Response {
	if in == nil {
		return nil
	}

	return &tfplugin6.ValidateActionConfig_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}

func PlanAction_Response(in *tfprotov6.PlanActionResponse) *tfplugin6.PlanAction_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.PlanAction_Response{
		LinkedResources: PlannedLinkedResources(in.LinkedResources),
		Diagnostics:     Diagnostics(in.Diagnostics),
		Deferred:        Deferred(in.Deferred),
	}

	return resp
}

func PlannedLinkedResources(in []*tfprotov6.PlannedLinkedResource) []*tfplugin6.PlanAction_Response_LinkedResource {
	resp := make([]*tfplugin6.PlanAction_Response_LinkedResource, 0, len(in))

	for _, inLinkedResource := range in {
		resp = append(resp, &tfplugin6.PlanAction_Response_LinkedResource{
			PlannedState:    DynamicValue(inLinkedResource.PlannedState),
			PlannedIdentity: ResourceIdentityData(inLinkedResource.PlannedIdentity),
		})
	}

	return resp
}

func InvokeAction_InvokeActionEvent(in *tfprotov6.InvokeActionEvent) *tfplugin6.InvokeAction_Event {
	if in == nil {
		return nil
	}

	switch event := (in.Type).(type) {
	case tfprotov6.ProgressInvokeActionEventType:
		return &tfplugin6.InvokeAction_Event{
			Type: &tfplugin6.InvokeAction_Event_Progress_{
				Progress: &tfplugin6.InvokeAction_Event_Progress{
					Message: event.Message,
				},
			},
		}
	case tfprotov6.CompletedInvokeActionEventType:
		return &tfplugin6.InvokeAction_Event{
			Type: &tfplugin6.InvokeAction_Event_Completed_{
				Completed: &tfplugin6.InvokeAction_Event_Completed{
					LinkedResources: NewLinkedResources(event.LinkedResources),
					Diagnostics:     Diagnostics(event.Diagnostics),
				},
			},
		}
	}

	// It is not currently possible to create tfprotov6.InvokeActionEventType
	// implementations outside the tfprotov6 package. If this panic was reached,
	// it implies that a new event type was introduced and needs to be implemented
	// as a new case above.
	panic(fmt.Sprintf("unimplemented tfprotov6.InvokeActionEventType type: %T", in.Type))
}

func NewLinkedResources(in []*tfprotov6.NewLinkedResource) []*tfplugin6.InvokeAction_Event_Completed_LinkedResource {
	resp := make([]*tfplugin6.InvokeAction_Event_Completed_LinkedResource, 0, len(in))

	for _, inLinkedResource := range in {
		resp = append(resp, &tfplugin6.InvokeAction_Event_Completed_LinkedResource{
			NewState:        DynamicValue(inLinkedResource.NewState),
			NewIdentity:     ResourceIdentityData(inLinkedResource.NewIdentity),
			RequiresReplace: inLinkedResource.RequiresReplace,
		})
	}

	return resp
}
