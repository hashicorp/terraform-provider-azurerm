// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/backuppolicy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetAppBackupPolicyDataSource struct{}

var _ sdk.DataSource = NetAppBackupPolicyDataSource{}

func (r NetAppBackupPolicyDataSource) ResourceType() string {
	return "azurerm_netapp_backup_policy"
}

func (r NetAppBackupPolicyDataSource) ModelObject() interface{} {
	return &netAppModels.NetAppBackupVaultModel{}
}

func (r NetAppBackupPolicyDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backuppolicy.ValidateBackupPolicyID
}

func (r NetAppBackupPolicyDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"account_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (r NetAppBackupPolicyDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"tags": commonschema.TagsDataSource(),

		"daily_backups_to_keep": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"weekly_backups_to_keep": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"monthly_backups_to_keep": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
	}
}

func (r NetAppBackupPolicyDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BackupPolicyClient

			var state netAppModels.NetAppBackupPolicyModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			backupPolicyID := backuppolicy.NewBackupPolicyID(metadata.Client.Account.SubscriptionId, state.ResourceGroupName, state.AccountName, state.Name)

			existing, err := client.BackupPoliciesGet(ctx, backupPolicyID)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", backupPolicyID)
				}
				return fmt.Errorf("retrieving %s: %v", backupPolicyID, err)
			}

			if model := existing.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
				state.DailyBackupsToKeep = pointer.From(model.Properties.DailyBackupsToKeep)
				state.WeeklyBackupsToKeep = pointer.From(model.Properties.WeeklyBackupsToKeep)
				state.MonthlyBackupsToKeep = pointer.From(model.Properties.MonthlyBackupsToKeep)
				state.Enabled = pointer.From(model.Properties.Enabled)
			}

			metadata.SetID(backupPolicyID)

			return metadata.Encode(&state)
		},
	}
}
