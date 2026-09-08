// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/internal/tfplugin6"
)

func ValidateStateStoreConfigRequest(in *tfplugin6.ValidateStateStoreConfig_Request) *tfprotov6.ValidateStateStoreConfigRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ValidateStateStoreConfigRequest{
		TypeName: in.TypeName,
		Config:   DynamicValue(in.Config),
	}
}

func ConfigureStateStoreRequest(in *tfplugin6.ConfigureStateStore_Request) *tfprotov6.ConfigureStateStoreRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ConfigureStateStoreRequest{
		TypeName:     in.TypeName,
		Config:       DynamicValue(in.Config),
		Capabilities: ConfigureStateStoreClientCapabilities(in.Capabilities),
	}
}

func ReadStateBytesRequest(in *tfplugin6.ReadStateBytes_Request) *tfprotov6.ReadStateBytesRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ReadStateBytesRequest{
		TypeName: in.TypeName,
		StateID:  in.StateId,
	}
}

func GetStatesRequest(in *tfplugin6.GetStates_Request) *tfprotov6.GetStatesRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.GetStatesRequest{
		TypeName: in.TypeName,
	}
}

func DeleteStateRequest(in *tfplugin6.DeleteState_Request) *tfprotov6.DeleteStateRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.DeleteStateRequest{
		TypeName: in.TypeName,
		StateID:  in.StateId,
	}
}

func LockStateRequest(in *tfplugin6.LockState_Request) *tfprotov6.LockStateRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.LockStateRequest{
		TypeName:  in.TypeName,
		StateID:   in.StateId,
		Operation: in.Operation,
	}
}

func UnlockStateRequest(in *tfplugin6.UnlockState_Request) *tfprotov6.UnlockStateRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.UnlockStateRequest{
		TypeName: in.TypeName,
		StateID:  in.StateId,
		LockID:   in.LockId,
	}
}

func WriteStateBytesChunk(in *tfplugin6.WriteStateBytes_RequestChunk) (*tfprotov6.WriteStateBytesChunk, *tfprotov6.Diagnostic) {
	if in == nil {
		return nil, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Unexpected empty state chunk in WriteStateBytes",
			Detail:   "An empty state byte chunk was received. This is a bug in Terraform that should be reported to the maintainers.",
		}
	}

	if in.Range == nil {
		return nil, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Unexpected state chunk data received in WriteStateBytes",
			Detail: "An invalid state byte chunk was received with no range start/end information. " +
				"This is a bug in Terraform that should be reported to the maintainers.",
		}
	}

	var meta *tfprotov6.WriteStateChunkMeta
	if in.Meta != nil {
		// Metadata is only attached to the first chunk
		meta = &tfprotov6.WriteStateChunkMeta{
			TypeName: in.Meta.TypeName,
			StateID:  in.Meta.StateId,
		}
	}

	stateChunk := &tfprotov6.WriteStateBytesChunk{
		Meta: meta,
		StateByteChunk: tfprotov6.StateByteChunk{
			Bytes:       in.Bytes,
			TotalLength: in.TotalLength,
			Range: tfprotov6.StateByteRange{
				Start: in.Range.Start,
				End:   in.Range.End,
			},
		},
	}

	return stateChunk, nil
}
