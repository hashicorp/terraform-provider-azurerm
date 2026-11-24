// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_cosmosdb_account

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name data_protection_backup_policy_cosmosdb_account -service-package-name dataprotection -properties "name" -compare-values "subscription_id:data_protection_backup_vault_id,resource_group_name:data_protection_backup_vault_id,backup_vault_name:data_protection_backup_vault_id"

type Resource struct{}

var (
	_ sdk.Resource                   = Resource{}
	_ sdk.ResourceWithIdentity       = Resource{}
	_ sdk.ResourceWithCustomImporter = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &basebackuppolicyresources.BackupPolicyId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_data_protection_backup_policy_cosmosdb_account"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return basebackuppolicyresources.ValidateBackupPolicyID
}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := basebackuppolicyresources.ParseBackupPolicyID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		client := metadata.Client.DataProtection.BackupPolicyClient20260601
		resp, err := client.BackupPoliciesGet(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}
		if resp.Model == nil {
			return fmt.Errorf("retrieving %s: model was nil", *id)
		}

		properties, ok := resp.Model.Properties.(basebackuppolicyresources.BackupPolicy)
		if !ok || !backupPolicySupportsCosmosdbAccounts(properties) {
			return fmt.Errorf("%s is not a backup policy for `Microsoft.DocumentDB/databaseAccounts`", *id)
		}

		return nil
	}
}
