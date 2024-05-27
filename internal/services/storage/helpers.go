package storage

import (
	"sort"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
)

func expandStringSlice(input []interface{}) []string {
	result := make([]string, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, item.(string))
		}
	}
	return result
}

func sortedKeysFromSlice(input map[storageaccounts.Kind]struct{}) []string {
	keys := make([]string, 0)
	for key := range input {
		keys = append(keys, string(key))
	}
	sort.Strings(keys)
	return keys
}
