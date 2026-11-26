// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
)

func FlattenBackupPolicyBlobStorageDefaultRetentionRuleDuration(input []basebackuppolicyresources.BasePolicyRule, dsType basebackuppolicyresources.DataStoreTypes) any {
	if input == nil {
		return nil
	}

	for _, item := range input {
		if retentionRule, ok := item.(basebackuppolicyresources.AzureRetentionRule); ok && retentionRule.IsDefault != nil && *retentionRule.IsDefault {
			if len(retentionRule.Lifecycles) > 0 {
				if deleteOption, ok := retentionRule.Lifecycles[0].DeleteAfter.(basebackuppolicyresources.AbsoluteDeleteOption); ok {
					if retentionRule.Lifecycles[0].SourceDataStore.DataStoreType == dsType {
						return deleteOption.Duration
					}
				}
			}
		}
	}
	return nil
}
