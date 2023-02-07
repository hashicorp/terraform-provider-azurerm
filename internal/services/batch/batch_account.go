package batch

import (
	"fmt"
<<<<<<< HEAD
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2022-01-01/batchaccount"
)

// expandBatchAccountKeyVaultReference expands Batch account KeyVault reference
func expandBatchAccountKeyVaultReference(list []interface{}) (*batchaccount.KeyVaultReference, error) {
	if len(list) == 0 {
		return nil, fmt.Errorf("key vault reference should be defined")
=======

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2022-01-01/batch" // nolint: staticcheck
)

// expandBatchAccountKeyVaultReference expands Batch account KeyVault reference
func expandBatchAccountKeyVaultReference(list []interface{}) (*batch.KeyVaultReference, error) {
	if len(list) == 0 || list[0] == nil {
		return nil, fmt.Errorf("Error: key vault reference should be defined")
>>>>>>> main
	}

	keyVaultRef := list[0].(map[string]interface{})

	ref := &batchaccount.KeyVaultReference{
		Id:  keyVaultRef["id"].(string),
		Url: keyVaultRef["url"].(string),
	}

	return ref, nil
}

// flattenBatchAccountKeyvaultReference flattens a Batch account keyvault reference
func flattenBatchAccountKeyvaultReference(keyVaultReference *batchaccount.KeyVaultReference) interface{} {
	result := make(map[string]interface{})

	if keyVaultReference == nil {
		return []interface{}{}
	}

	if keyVaultReference.Id != "" {
		result["id"] = keyVaultReference.Id
	}

	if keyVaultReference.Url != "" {
		result["url"] = keyVaultReference.Url
	}

	return []interface{}{result}
}
