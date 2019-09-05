package azure

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2018-12-01/batch"
)

// ExpandBatchAccountKeyVaultReference expands Batch account KeyVault reference
func ExpandBatchAccountKeyVaultReference(list []interface{}) (*batch.KeyVaultReference, error) {
	if len(list) == 0 {
		return nil, fmt.Errorf("Error: key vault reference should be defined")
	}

	keyVaultRef := list[0].(map[string]interface{})

	keyVaultRefID := keyVaultRef["id"].(string)
	keyVaultRefURL := keyVaultRef["url"].(string)

	ref := &batch.KeyVaultReference{
		ID:  &keyVaultRefID,
		URL: &keyVaultRefURL,
	}

	return ref, nil
}

// FlattenBatchAccountKeyvaultReference flattens a Batch account keyvault reference
func FlattenBatchAccountKeyvaultReference(keyVaultReference *batch.KeyVaultReference) interface{} {
	result := make(map[string]interface{})

	if keyVaultReference == nil {
		return []interface{}{}
	}

	if keyVaultReference.ID != nil {
		result["id"] = *keyVaultReference.ID
	}

	if keyVaultReference.URL != nil {
		result["url"] = *keyVaultReference.URL
	}

	return []interface{}{result}
}
