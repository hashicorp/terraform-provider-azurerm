package common

import (
	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-01-15/documentdb"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func ExpandCosmosDbConflicResolutionPolicy(inputs []interface{}) *documentdb.ConflictResolutionPolicy {
	if len(inputs) == 0 || inputs[0] == nil {
		return nil
	}

	input := inputs[0].(map[string]interface{})
	conflictResolutionMode := input["mode"].(string)
	conflict := &documentdb.ConflictResolutionPolicy{
		Mode: documentdb.ConflictResolutionMode(conflictResolutionMode),
	}

	if conflictResolutionPath, ok := input["conflict_resolution_path"].(string); ok {
		conflict.ConflictResolutionPath = utils.String(conflictResolutionPath)
	}

	if conflictResolutionProcedure, ok := input["conflict_resolution_procedure"].(string); ok {
		conflict.ConflictResolutionProcedure = utils.String(conflictResolutionProcedure)
	}

	return conflict
}

func FlattenCosmosDbConflictResolutionPolicy(input *documentdb.ConflictResolutionPolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	conflictResolutionPolicy := make(map[string]interface{})

	conflictResolutionPolicy["mode"] = string(input.Mode)
	var path, procedure string
	if input.ConflictResolutionPath != nil {
		path = *input.ConflictResolutionPath
	}
	if input.ConflictResolutionProcedure != nil {
		procedure = *input.ConflictResolutionProcedure
	}

	return []interface{}{
		map[string]interface{}{
			"mode":                          string(input.Mode),
			"conflict_resolution_path":      path,
			"conflict_resolution_procedure": procedure,
		},
	}
}
