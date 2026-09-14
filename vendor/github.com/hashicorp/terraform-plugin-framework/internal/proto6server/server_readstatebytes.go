// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package proto6server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/toproto6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func (s *Server) ReadStateBytes(ctx context.Context, proto6Req *tfprotov6.ReadStateBytesRequest) (*tfprotov6.ReadStateBytesStream, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &tfprotov6.ReadStateBytesStream{}

	statestore, diags := s.FrameworkServer.StateStore(ctx, proto6Req.TypeName)

	if diags.HasError() {
		return &tfprotov6.ReadStateBytesStream{
			Chunks: func(push func(chunk tfprotov6.ReadStateByteChunk) bool) {
				push(tfprotov6.ReadStateByteChunk{
					Diagnostics: toproto6.Diagnostics(ctx, diags),
				})
			},
		}, nil
	}

	readStateBytesReq := &fwserver.ReadStateBytesRequest{
		StateID:    proto6Req.StateID,
		StateStore: statestore,
	}
	readStateBytesResp := &fwserver.ReadStateBytesResponse{}

	s.FrameworkServer.ReadStateBytes(ctx, readStateBytesReq, readStateBytesResp)

	if readStateBytesResp.Diagnostics.HasError() {
		return &tfprotov6.ReadStateBytesStream{
			Chunks: func(push func(chunk tfprotov6.ReadStateByteChunk) bool) {
				push(tfprotov6.ReadStateByteChunk{
					Diagnostics: toproto6.Diagnostics(ctx, readStateBytesResp.Diagnostics),
				})
			},
		}, nil
	}

	// If ConfigureStateStore isn't called prior to ReadStateBytes
	if int(s.FrameworkServer.StateStoreConfigureData.ServerCapabilities.ChunkSize) <= 0 {
		return &tfprotov6.ReadStateBytesStream{
			Chunks: func(push func(chunk tfprotov6.ReadStateByteChunk) bool) {
				push(tfprotov6.ReadStateByteChunk{
					Diagnostics: []*tfprotov6.Diagnostic{
						{
							Severity: tfprotov6.DiagnosticSeverityError,
							Summary:  "Error reading state",
							Detail:   "The provider server does not have a chunk size configured. This is a bug in either Terraform core or terraform-plugin-framework and should be reported.",
						},
					},
				})
			},
		}, nil
	}

	chunkSize := int(s.FrameworkServer.StateStoreConfigureData.ServerCapabilities.ChunkSize)

	reader := bytes.NewReader(readStateBytesResp.StateBytes)
	totalLength := reader.Size()
	rangeStart := 0

	fwResp.Chunks = func(yield func(tfprotov6.ReadStateByteChunk) bool) {
		for {
			readBytes := make([]byte, chunkSize)
			byteCount, err := reader.Read(readBytes)
			if err != nil && !errors.Is(err, io.EOF) {
				chunkWithDiag := tfprotov6.ReadStateByteChunk{
					Diagnostics: []*tfprotov6.Diagnostic{
						{
							Severity: tfprotov6.DiagnosticSeverityError,
							Summary:  "Error reading state",
							Detail: fmt.Sprintf("An unexpected error occurred while reading state data for %s: %s",
								proto6Req.StateID,
								err,
							),
						},
					},
				}
				if !yield(chunkWithDiag) {
					return
				}
			}

			if byteCount == 0 {
				return
			}

			chunk := tfprotov6.ReadStateByteChunk{
				StateByteChunk: tfprotov6.StateByteChunk{
					Bytes:       readBytes[:byteCount],
					TotalLength: totalLength,
					Range: tfprotov6.StateByteRange{
						Start: int64(rangeStart),
						End:   int64(rangeStart + byteCount - 1),
					},
				},
			}
			if !yield(chunk) {
				return
			}

			rangeStart += byteCount
		}
	}

	return fwResp, nil
}
