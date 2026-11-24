// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_database_extended_auditing_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/blobauditing"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlDatabaseExtendedAuditingPolicyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.BlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabaseExtendedAuditingPolicyID(d.Id())
	if err != nil {
		return err
	}

	dbId := commonids.NewSqlDatabaseID(id.SubscriptionId, id.ResourceGroup, id.ServerName, id.DatabaseName)

	resp, err := client.ExtendedDatabaseBlobAuditingPoliciesGet(ctx, dbId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading MsSql Database %s Extended Auditing Policy (MsSql Server Name %q / Resource Group %q): %s", id.DatabaseName, id.ServerName, id.ResourceGroup, err)
	}

	databaseId := commonids.NewSqlDatabaseID(id.SubscriptionId, id.ResourceGroup, id.ServerName, id.DatabaseName)
	d.Set("database_id", databaseId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("blob_storage_endpoint", props.StorageEndpoint)
			d.Set("storage_account_access_key_is_secondary", props.IsStorageSecondaryKeyInUse)
			d.Set("retention_in_days", props.RetentionDays)
			d.Set("log_monitoring_enabled", props.IsAzureMonitorTargetEnabled)
			d.Set("enabled", props.State == blobauditing.BlobAuditingPolicyStateEnabled)
		}
	}

	return nil
}
