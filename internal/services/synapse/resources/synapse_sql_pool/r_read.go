// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_sql_pool

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseSqlPoolRead(d *pluginsdk.ResourceData, meta any) error {
	sqlClient := meta.(*clients.Client).Synapse.SqlPoolClient
	sqlPTDEClient := meta.(*clients.Client).Synapse.SqlPoolTransparentDataEncryptionClient
	geoBackUpClient := meta.(*clients.Client).Synapse.SqlPoolGeoBackupPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlPoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := sqlClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	transparentDataEncryption, err := sqlPTDEClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving Transparent Data Encryption settings of %s: %+v", *id, err)
	}

	geoBackupPolicy, err := geoBackUpClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving Geo Backup Policy of %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("synapse_workspace_id", workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID())
	if resp.Sku != nil {
		d.Set("sku_name", resp.Sku.Name)
	}

	if props := resp.SQLPoolResourceProperties; props != nil {
		d.Set("collation", props.Collation)
		d.Set("storage_account_type", props.StorageAccountType)
	}

	geoBackupEnabled := true
	if geoBackupProps := geoBackupPolicy.GeoBackupPolicyProperties; geoBackupProps != nil {
		geoBackupEnabled = geoBackupProps.State == synapse.GeoBackupPolicyStateEnabled
	}
	d.Set("geo_backup_policy_enabled", geoBackupEnabled)

	if tdeProps := transparentDataEncryption.TransparentDataEncryptionProperties; tdeProps != nil {
		d.Set("data_encrypted", tdeProps.Status == synapse.TransparentDataEncryptionStatusEnabled)
	}

	// whole "restore" block is not returned. to avoid conflict, so set it from the old state
	d.Set("restore", d.Get("restore").([]any))

	return tags.FlattenAndSet(d, resp.Tags)
}

const (
	DefaultCreateMode            = "Default"
	RecoveryCreateMode           = "Recovery"
	PointInTimeRestoreCreateMode = "PointInTimeRestore"
)

func synapseSqlPoolCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
	_, value := d.GetChange("geo_backup_policy_enabled")
	geoBackupEnabled := value.(bool)

	_, value = d.GetChange("storage_account_type")
	storageAccountType := synapse.StorageAccountType(value.(string))

	if storageAccountType == synapse.StorageAccountTypeLRS && geoBackupEnabled {
		return fmt.Errorf("`geo_backup_policy_enabled` cannot be `true` if the `storage_account_type` is `LRS`")
	}

	return nil
}

func synapseSqlPoolScaleStateRefreshFunc(ctx context.Context, client *synapse.SQLPoolsClient, resourceGroup, workspaceName, name string) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		resp, err := client.Get(ctx, resourceGroup, workspaceName, name)
		if err != nil {
			return resp, "failed", err
		}

		if resp.SQLPoolResourceProperties == nil || resp.Status == nil {
			return resp, "failed", nil
		}

		return resp, *resp.Status, nil
	}
}

// sqlPool backend service is a proxy to sql database
// backend service restore and backup only accept id format of sql database
// so if the id is sqlPool, we need to construct the corresponding sql database id
func constructSourceDatabaseId(id string) string {
	sqlPoolId, err := parse.SqlPoolID(id)
	if err != nil {
		return id
	}

	return commonids.NewSqlDatabaseID(sqlPoolId.SubscriptionId, sqlPoolId.ResourceGroup, sqlPoolId.WorkspaceName, sqlPoolId.Name).ID()
}
