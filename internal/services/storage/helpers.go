// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"sort"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
)

func sortedKeysFromSlice(input map[storageaccounts.Kind]struct{}) []string {
	keys := make([]string, 0)
	for key := range input {
		keys = append(keys, string(key))
	}
	sort.Strings(keys)
	return keys
}

func validateExistingModel(input *storageaccounts.StorageAccount, id *commonids.StorageAccountId) error {
	if input == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if input.Kind == nil {
		return fmt.Errorf("retrieving %s: `model.Kind` was nil", id)
	}

	if input.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", id)
	}

	if input.Sku == nil {
		return fmt.Errorf("retrieving %s: `model.Sku` was nil", id)
	}

	return nil
}
