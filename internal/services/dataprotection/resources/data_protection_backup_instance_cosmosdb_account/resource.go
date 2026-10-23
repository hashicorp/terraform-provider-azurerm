// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_cosmosdb_account

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstanceresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name data_protection_backup_instance_cosmosdb_account -service-package-name dataprotection -properties "name" -compare-values "subscription_id:data_protection_backup_vault_id,resource_group_name:data_protection_backup_vault_id,backup_vault_name:data_protection_backup_vault_id"

type Resource struct{}

var (
	_ sdk.Resource                   = Resource{}
	_ sdk.ResourceWithIdentity       = Resource{}
	_ sdk.ResourceWithCustomImporter = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &backupinstanceresources.BackupInstanceId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_data_protection_backup_instance_cosmosdb_account"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backupinstanceresources.ValidateBackupInstanceID
}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := backupinstanceresources.ParseBackupInstanceID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		client := metadata.Client.DataProtection.BackupInstanceClient20260601
		resp, err := client.BackupInstancesGet(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if resp.Model == nil {
			return fmt.Errorf("retrieving %s: `model` was nil", *id)
		}
		if resp.Model.Properties == nil || !backupInstanceSupportsCosmosdbAccounts(*resp.Model.Properties) {
			return fmt.Errorf("%s is not a backup instance for `%s`", *id, cosmosDBAccountDataSourceType)
		}

		return nil
	}
}

func backupInstanceSupportsCosmosdbAccounts(input backupinstanceresources.BackupInstance) bool {
	return strings.EqualFold(pointer.From(input.DataSourceInfo.ResourceType), cosmosDBAccountDataSourceType)
}
