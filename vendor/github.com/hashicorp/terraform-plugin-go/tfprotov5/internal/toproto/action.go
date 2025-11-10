// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/internal/tfplugin5"
)

func GetMetadata_ActionMetadata(in *tfprotov5.ActionMetadata) *tfplugin5.GetMetadata_ActionMetadata {
	if in == nil {
		return nil
	}

	return &tfplugin5.GetMetadata_ActionMetadata{
		TypeName: in.TypeName,
	}
}

func ValidateActionConfig_Response(in *tfprotov5.ValidateActionConfigResponse) *tfplugin5.ValidateActionConfig_Response {
	if in == nil {
		return nil
	}

	return &tfplugin5.ValidateActionConfig_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}

func PlanAction_Response(in *tfprotov5.PlanActionResponse) *tfplugin5.PlanAction_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.PlanAction_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
		Deferred:    Deferred(in.Deferred),
	}

	return resp
}

func InvokeAction_InvokeActionEvent(in *tfprotov5.InvokeActionEvent) *tfplugin5.InvokeAction_Event {
	if in == nil {
		return nil
	}

	switch event := (in.Type).(type) {
	case tfprotov5.ProgressInvokeActionEventType:
		return &tfplugin5.InvokeAction_Event{
			Type: &tfplugin5.InvokeAction_Event_Progress_{
				Progress: &tfplugin5.InvokeAction_Event_Progress{
					Message: event.Message,
				},
			},
		}
	case tfprotov5.CompletedInvokeActionEventType:
		return &tfplugin5.InvokeAction_Event{
			Type: &tfplugin5.InvokeAction_Event_Completed_{
				Completed: &tfplugin5.InvokeAction_Event_Completed{
					Diagnostics: Diagnostics(event.Diagnostics),
				},
			},
		}
	}

	// It is not currently possible to create tfprotov5.InvokeActionEventType
	// implementations outside the tfprotov5 package. If this panic was reached,
	// it implies that a new event type was introduced and needs to be implemented
	// as a new case above.
	panic(fmt.Sprintf("unimplemented tfprotov5.InvokeActionEventType type: %T", in.Type))
}
