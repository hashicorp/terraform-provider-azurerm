// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccount"
)

// expandBatchAccountKeyVaultReference expands Batch account KeyVault reference
func expandBatchAccountKeyVaultReference(list []any) (*batchaccount.KeyVaultReference, error) {
	if len(list) == 0 || list[0] == nil {
		return nil, fmt.Errorf("key vault reference should be defined")
	}

	keyVaultRef := list[0].(map[string]any)

	ref := &batchaccount.KeyVaultReference{
		Id:  keyVaultRef["id"].(string),
		Url: keyVaultRef["url"].(string),
	}

	return ref, nil
}

// flattenBatchAccountKeyvaultReference flattens a Batch account keyvault reference
func flattenBatchAccountKeyvaultReference(keyVaultReference *batchaccount.KeyVaultReference) any {
	result := make(map[string]any)

	if keyVaultReference == nil {
		return []any{}
	}

	if keyVaultReference.Id != "" {
		result["id"] = keyVaultReference.Id
	}

	if keyVaultReference.Url != "" {
		result["url"] = keyVaultReference.Url
	}

	return []any{result}
}
