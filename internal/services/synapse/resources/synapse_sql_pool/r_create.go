// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_sql_pool

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseSqlPoolCreate(d *pluginsdk.ResourceData, meta any) error {
	sqlClient := client.From(meta).SqlPoolClient
	sqlPTDEClient := client.From(meta).SqlPoolTransparentDataEncryptionClient
	workspaceClient := client.From(meta).WorkspaceClient
	geoBackUpClient := client.From(meta).SqlPoolGeoBackupPoliciesClient
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	// todo 6.0 - move to the case-sensitive parser when validation.AsGeneratedID is removed: this parses a config
	// value which the paired AsGeneratedID validator accepts with legacy casing, and configs cannot be migrated.
	workspaceId, err := workspaces.ParseWorkspaceIDInsensitively(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewSqlPoolID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := sqlClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.Response.Response) {
			return tf.ImportAsExistsError("azurerm_synapse_sql_pool", id.ID())
		}
	}
	workspace, err := workspaceClient.Get(ctx, workspaceId.ResourceGroupName, workspaceId.WorkspaceName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", workspaceId, err)
	}

	geoBackupEnabled := d.Get("geo_backup_policy_enabled").(bool)
	storageAccountType := synapse.StorageAccountType(d.Get("storage_account_type").(string))

	mode := d.Get("create_mode").(string)
	sqlPoolInfo := synapse.SQLPool{
		Location: workspace.Location,
		SQLPoolResourceProperties: &synapse.SQLPoolResourceProperties{
			CreateMode:         synapse.CreateMode(*pointer.To(mode)),
			StorageAccountType: storageAccountType,
		},
		Sku: &synapse.Sku{
			Name: pointer.To(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	switch mode {
	case DefaultCreateMode:
		sqlPoolInfo.Collation = pointer.To(d.Get("collation").(string))
	case RecoveryCreateMode:
		recoveryDatabaseId := constructSourceDatabaseId(d.Get("recovery_database_id").(string))

		if recoveryDatabaseId == "" {
			return fmt.Errorf("`recovery_database_id` must be set when `create_mode` is %q", RecoveryCreateMode)
		}

		sqlPoolInfo.RecoverableDatabaseID = pointer.To(recoveryDatabaseId)
	case PointInTimeRestoreCreateMode:
		restore := d.Get("restore").([]any)
		if len(restore) == 0 || restore[0] == nil {
			return fmt.Errorf("`restore` block must be set when `create_mode` is %q", PointInTimeRestoreCreateMode)
		}

		v := restore[0].(map[string]any)
		sourceDatabaseId := constructSourceDatabaseId(v["source_database_id"].(string))
		vTime, parseErr := date.ParseTime(time.RFC3339, v["point_in_time"].(string))

		if parseErr != nil {
			return fmt.Errorf("parsing time format: %+v", parseErr)
		}

		sqlPoolInfo.RestorePointInTime = &date.Time{Time: vTime}
		sqlPoolInfo.SourceDatabaseID = pointer.To(sourceDatabaseId)
	}

	future, err := sqlClient.Create(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, sqlPoolInfo)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if err = future.WaitForCompletionRef(ctx, sqlClient.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	if d.Get("data_encrypted").(bool) {
		parameter := synapse.TransparentDataEncryption{
			TransparentDataEncryptionProperties: &synapse.TransparentDataEncryptionProperties{
				Status: synapse.TransparentDataEncryptionStatusEnabled,
			},
		}

		if _, err := sqlPTDEClient.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, parameter); err != nil {
			return fmt.Errorf("setting `data_encrypted`: %+v", err)
		}
	}

	// Only update the Geo Backup Policy if it has been disabled since it is
	// already enabled by default...
	if !geoBackupEnabled {
		geoBackupParams := synapse.GeoBackupPolicy{
			GeoBackupPolicyProperties: &synapse.GeoBackupPolicyProperties{
				State: synapse.GeoBackupPolicyStateDisabled,
			},
		}

		if _, err := geoBackUpClient.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, geoBackupParams); err != nil {
			return fmt.Errorf("setting `geo_backup_policy_enabled`: %+v", err)
		}
	}

	return resourceSynapseSqlPoolRead(d, meta)
}
