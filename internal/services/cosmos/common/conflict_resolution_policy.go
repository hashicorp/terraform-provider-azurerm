// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2026-03-15/openapis"
)

func ExpandCosmosDbConflicResolutionPolicy(inputs []any) *openapis.ConflictResolutionPolicy {
	if len(inputs) == 0 || inputs[0] == nil {
		return nil
	}

	input := inputs[0].(map[string]any)
	conflict := &openapis.ConflictResolutionPolicy{
		Mode: pointer.ToEnum[openapis.ConflictResolutionMode](input["mode"].(string)),
	}

	if conflictResolutionPath, ok := input["conflict_resolution_path"].(string); ok {
		conflict.ConflictResolutionPath = pointer.To(conflictResolutionPath)
	}

	if conflictResolutionProcedure, ok := input["conflict_resolution_procedure"].(string); ok {
		conflict.ConflictResolutionProcedure = pointer.To(conflictResolutionProcedure)
	}

	return conflict
}

func FlattenCosmosDbConflictResolutionPolicy(input *openapis.ConflictResolutionPolicy) []any {
	if input == nil {
		return []any{}
	}
	conflictResolutionPolicy := make(map[string]any)

	conflictResolutionPolicy["mode"] = input.Mode
	var path, procedure string
	if input.ConflictResolutionPath != nil {
		path = *input.ConflictResolutionPath
	}
	if input.ConflictResolutionProcedure != nil {
		procedure = *input.ConflictResolutionProcedure
	}

	return []any{
		map[string]any{
			"mode":                          input.Mode,
			"conflict_resolution_path":      path,
			"conflict_resolution_procedure": procedure,
		},
	}
}
