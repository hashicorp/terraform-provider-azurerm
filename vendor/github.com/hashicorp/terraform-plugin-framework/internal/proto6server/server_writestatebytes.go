// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package proto6server

import (
	"bytes"
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/toproto6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-log/tfsdklog"
)

func (s *Server) WriteStateBytes(ctx context.Context, proto6Req *tfprotov6.WriteStateBytesStream) (*tfprotov6.WriteStateBytesResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.WriteStateBytesResponse{}

	var typeName string
	var stateID string
	var stateBuffer bytes.Buffer

	// Chunk size is set on the server during the ConfigureStateStore RPC
	configuredChunkSize := s.FrameworkServer.StateStoreConfigureData.ServerCapabilities.ChunkSize
	if configuredChunkSize <= 0 {
		fwResp.Diagnostics.AddError(
			"Error Writing State",
			"The provider server does not have a chunk size configured. This is a bug in either Terraform core or terraform-plugin-framework and should be reported.",
		)
		return toproto6.WriteStateBytesResponse(ctx, fwResp), nil
	}

	// This field will be collected from one of the chunks
	var expectedTotalLength int64 = 0

	for chunk, diags := range proto6Req.Chunks {
		// Any diagnostics here are either a GRPC communication error or invalid data from the client.
		if len(diags) > 0 {
			return &tfprotov6.WriteStateBytesResponse{Diagnostics: diags}, nil
		}

		// Only the first chunk will have meta set
		if chunk.Meta != nil {
			typeName = chunk.Meta.TypeName
			stateID = chunk.Meta.StateID
		}

		if chunk.Range.End < chunk.TotalLength-1 {
			// Ensure each chunk (except the last) exactly match the configured size.
			if int64(len(chunk.Bytes)) != configuredChunkSize {
				fwResp.Diagnostics.AddError(
					"Error Writing State",
					fmt.Sprintf("Unexpected chunk of size %d was received from Terraform, expected chunk size was %d. This is a bug in Terraform core and should be reported.", len(chunk.Bytes), configuredChunkSize),
				)
				return toproto6.WriteStateBytesResponse(ctx, fwResp), nil
			}
		} else {
			// Ensure the last chunk is within the configured size.
			if int64(len(chunk.Bytes)) > configuredChunkSize {
				fwResp.Diagnostics.AddError(
					"Error Writing State",
					fmt.Sprintf("Unexpected final chunk of size %d was received from Terraform, which exceeds the configured chunk size of %d. This is a bug in Terraform core and should be reported.", len(chunk.Bytes), configuredChunkSize),
				)
				return toproto6.WriteStateBytesResponse(ctx, fwResp), nil
			}
		}

		if expectedTotalLength == 0 {
			expectedTotalLength = chunk.TotalLength
		}

		stateBuffer.Write(chunk.Bytes)
	}

	// MAINTAINER NOTE: Typically these fields are set in terraform-plugin-go (see link below), however because
	// the type name is not extracted until the stream is consumed, it's easier to set the logger fields here.
	//
	// https://github.com/hashicorp/terraform-plugin-go/blob/14fe65ea923b5e306dbb4f67f2bb861f74b9e3ec/internal/logging/context.go#L112-L119
	ctx = tfsdklog.SetField(ctx, "tf_state_store_type", typeName)
	ctx = tfsdklog.SubsystemSetField(ctx, "proto", "tf_state_store_type", typeName)
	ctx = tflog.SetField(ctx, "tf_state_store_type", typeName)

	if stateBuffer.Len() == 0 {
		fwResp.Diagnostics.AddError(
			"Error Writing State",
			"No state data was received from Terraform. This is a bug in Terraform core and should be reported.",
		)
		return toproto6.WriteStateBytesResponse(ctx, fwResp), nil
	}

	if int64(stateBuffer.Len()) != expectedTotalLength {
		fwResp.Diagnostics.AddError(
			"Error Writing State",
			fmt.Sprintf("Unexpected size of state data received from Terraform, got: %d, expected: %d. This is a bug in Terraform core and should be reported.", stateBuffer.Len(), expectedTotalLength),
		)
		return toproto6.WriteStateBytesResponse(ctx, fwResp), nil
	}

	if stateID == "" {
		fwResp.Diagnostics.AddError(
			"Error Writing State",
			"No state ID was received from Terraform. This is a bug in Terraform core and should be reported.",
		)
		return toproto6.WriteStateBytesResponse(ctx, fwResp), nil
	}

	stateStore, diags := s.FrameworkServer.StateStore(ctx, typeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.WriteStateBytesResponse(ctx, fwResp), nil
	}

	fwReq := &fwserver.WriteStateBytesRequest{
		StateStore: stateStore,
		StateID:    stateID,
		StateBytes: stateBuffer.Bytes(),
	}

	s.FrameworkServer.WriteStateBytes(ctx, fwReq, fwResp)

	return toproto6.WriteStateBytesResponse(ctx, fwResp), nil
}
