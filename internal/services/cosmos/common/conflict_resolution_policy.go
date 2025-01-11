// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
		conflict.ConflictResolutionPath = utils.String(conflictResolutionPath)
	}

	if conflictResolutionProcedure, ok := input["conflict_resolution_procedure"].(string); ok {
		conflict.ConflictResolutionProcedure = utils.String(conflictResolutionProcedure)
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
