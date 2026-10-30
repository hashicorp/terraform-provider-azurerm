// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type mysqlFlexibleServerRenameServer struct {
	tfprotov5.ProviderServer
}

// NewRenameAwareProviderServer preserves unknown rename IDs without the Plugin SDK's implicit ID replacement.
func NewRenameAwareProviderServer(server tfprotov5.ProviderServer) tfprotov5.ProviderServer {
	return &mysqlFlexibleServerRenameServer{ProviderServer: server}
}

func (s *mysqlFlexibleServerRenameServer) PlanResourceChange(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	resp, err := s.ProviderServer.PlanResourceChange(ctx, req)
	if err != nil || req.TypeName != mysqlFlexibleServerResourceName {
		return resp, err
	}
	for _, diagnostic := range resp.Diagnostics {
		if diagnostic.Severity == tfprotov5.DiagnosticSeverityError {
			return resp, nil
		}
	}
	if len(resp.RequiresReplace) != 1 || !resp.RequiresReplace[0].Equal(tftypes.NewAttributePath().WithAttributeName("id")) {
		return resp, nil
	}
	priorNull, err := req.PriorState.IsNull()
	if err != nil {
		return resp, fmt.Errorf("checking MySQL rename prior state: %w", err)
	}
	if priorNull {
		return resp, nil
	}

	schemas, err := s.ProviderServer.GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})
	if err != nil {
		return resp, fmt.Errorf("retrieving MySQL rename schema: %w", err)
	}
	resp.Diagnostics = append(resp.Diagnostics, schemas.Diagnostics...)
	for _, diagnostic := range schemas.Diagnostics {
		if diagnostic.Severity == tfprotov5.DiagnosticSeverityError {
			return resp, nil
		}
	}
	resourceSchema := schemas.ResourceSchemas[req.TypeName]
	if resourceSchema == nil {
		return resp, fmt.Errorf("missing schema for %s", req.TypeName)
	}
	prior, err := req.PriorState.Unmarshal(resourceSchema.ValueType())
	if err != nil {
		return resp, fmt.Errorf("decoding MySQL rename prior state: %w", err)
	}
	planned, err := resp.PlannedState.Unmarshal(resourceSchema.ValueType())
	if err != nil {
		return resp, fmt.Errorf("decoding MySQL rename planned state: %w", err)
	}
	if planned.IsNull() || !planned.IsKnown() {
		return resp, nil
	}
	var oldValues, newValues map[string]tftypes.Value
	if err := prior.As(&oldValues); err != nil {
		return resp, fmt.Errorf("reading MySQL rename prior state: %w", err)
	}
	if err := planned.As(&newValues); err != nil {
		return resp, fmt.Errorf("reading MySQL rename planned state: %w", err)
	}
	if newValues["id"].IsKnown() || newValues["fqdn"].IsKnown() ||
		!newValues["rename_enabled"].IsKnown() || newValues["rename_enabled"].IsNull() ||
		!newValues["name"].IsKnown() || newValues["name"].IsNull() ||
		!oldValues["id"].IsKnown() || oldValues["id"].IsNull() ||
		!oldValues["name"].IsKnown() || oldValues["name"].IsNull() {
		return resp, nil
	}
	var enabled bool
	var oldID, oldName, newName string
	if err := newValues["rename_enabled"].As(&enabled); err != nil {
		return resp, fmt.Errorf("reading MySQL rename opt-in: %w", err)
	}
	if err := oldValues["id"].As(&oldID); err != nil {
		return resp, fmt.Errorf("reading MySQL rename source ID: %w", err)
	}
	if err := oldValues["name"].As(&oldName); err != nil {
		return resp, fmt.Errorf("reading MySQL rename source name: %w", err)
	}
	if err := newValues["name"].As(&newName); err != nil {
		return resp, fmt.Errorf("reading MySQL rename destination name: %w", err)
	}
	if enabled && oldID != "" && newName != "" && !strings.EqualFold(oldName, newName) {
		// CustomizeDiff has accepted an opted-in rename. Preserve every real ForceNew path;
		// only the SDK's synthetic replacement for its unknown ID is removed.
		resp.RequiresReplace = nil
	}
	return resp, nil
}
