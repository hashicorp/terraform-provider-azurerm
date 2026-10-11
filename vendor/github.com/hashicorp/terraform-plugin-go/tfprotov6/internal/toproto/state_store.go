// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/internal/tfplugin6"
)

func GetMetadata_StateStoreMetadata(in *tfprotov6.StateStoreMetadata) *tfplugin6.GetMetadata_StateStoreMetadata {
	if in == nil {
		return nil
	}

	return &tfplugin6.GetMetadata_StateStoreMetadata{
		TypeName: in.TypeName,
	}
}

func ValidateStateStoreConfig_Response(in *tfprotov6.ValidateStateStoreConfigResponse) *tfplugin6.ValidateStateStoreConfig_Response {
	if in == nil {
		return nil
	}

	return &tfplugin6.ValidateStateStoreConfig_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}

func ConfigureStateStore_Response(in *tfprotov6.ConfigureStateStoreResponse) *tfplugin6.ConfigureStateStore_Response {
	if in == nil {
		return nil
	}

	return &tfplugin6.ConfigureStateStore_Response{
		Diagnostics:  Diagnostics(in.Diagnostics),
		Capabilities: StateStoreServerCapabilities(in.Capabilities),
	}
}

func ReadStateBytes_ResponseChunk(in *tfprotov6.ReadStateByteChunk) *tfplugin6.ReadStateBytes_ResponseChunk {
	if in == nil {
		return nil
	}

	return &tfplugin6.ReadStateBytes_ResponseChunk{
		Diagnostics: Diagnostics(in.Diagnostics),
		Bytes:       in.Bytes,
		TotalLength: in.TotalLength,
		Range:       StateByteRange(in.Range),
	}
}

func StateByteRange(in tfprotov6.StateByteRange) *tfplugin6.StateByteRange {
	return &tfplugin6.StateByteRange{
		Start: in.Start,
		End:   in.End,
	}
}

func GetStates_Response(in *tfprotov6.GetStatesResponse) *tfplugin6.GetStates_Response {
	if in == nil {
		return nil
	}

	return &tfplugin6.GetStates_Response{
		StateIds:    in.StateIDs,
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}

func DeleteState_Response(in *tfprotov6.DeleteStateResponse) *tfplugin6.DeleteState_Response {
	if in == nil {
		return nil
	}

	return &tfplugin6.DeleteState_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}

func LockState_Response(in *tfprotov6.LockStateResponse) *tfplugin6.LockState_Response {
	if in == nil {
		return nil
	}

	return &tfplugin6.LockState_Response{
		LockId:      in.LockID,
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}

func UnlockState_Response(in *tfprotov6.UnlockStateResponse) *tfplugin6.UnlockState_Response {
	if in == nil {
		return nil
	}

	return &tfplugin6.UnlockState_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}
