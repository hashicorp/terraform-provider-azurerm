// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"sort"

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
