// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2025-04-15/cosmosdb"
)

func ExpandCosmosDbConflicResolutionPolicy(inputs []interface{}) *cosmosdb.ConflictResolutionPolicy {
	if len(inputs) == 0 || inputs[0] == nil {
		return nil
	}

	input := inputs[0].(map[string]interface{})
	conflictResolutionMode := cosmosdb.ConflictResolutionMode(input["mode"].(string))
	conflict := &cosmosdb.ConflictResolutionPolicy{
		Mode: &conflictResolutionMode,
	}

	if conflictResolutionPath, ok := input["conflict_resolution_path"].(string); ok {
		conflict.ConflictResolutionPath = pointer.To(conflictResolutionPath)
	}

	if conflictResolutionProcedure, ok := input["conflict_resolution_procedure"].(string); ok {
		conflict.ConflictResolutionProcedure = pointer.To(conflictResolutionProcedure)
	}

	return conflict
}

func FlattenCosmosDbConflictResolutionPolicy(input *cosmosdb.ConflictResolutionPolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	conflictResolutionPolicy := make(map[string]interface{})

	conflictResolutionPolicy["mode"] = input.Mode
	var path, procedure string
	if input.ConflictResolutionPath != nil {
		path = *input.ConflictResolutionPath
	}
	if input.ConflictResolutionProcedure != nil {
		procedure = *input.ConflictResolutionProcedure
	}

	return []interface{}{
		map[string]interface{}{
			"mode":                          input.Mode,
			"conflict_resolution_path":      path,
			"conflict_resolution_procedure": procedure,
		},
	}
}
