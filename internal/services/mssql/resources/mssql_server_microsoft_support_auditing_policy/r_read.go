// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server_microsoft_support_auditing_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/serverdevopsaudit"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlServerMicrosoftSupportAuditingPolicyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.ServerDevOpsAuditSettingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerMicrosoftSupportAuditingPolicyID(d.Id())
	if err != nil {
		return err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	resp, err := client.SettingsGet(ctx, serverId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %s", *id, err)
	}

	d.Set("server_id", serverId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("blob_storage_endpoint", props.StorageEndpoint)
			d.Set("log_monitoring_enabled", props.IsAzureMonitorTargetEnabled)
			d.Set("enabled", props.State == serverdevopsaudit.BlobAuditingPolicyStateEnabled)

			if pointer.From(props.StorageAccountSubscriptionId) != "00000000-0000-0000-0000-000000000000" {
				d.Set("storage_account_subscription_id", pointer.From(props.StorageAccountSubscriptionId))
			}
		}
	}

	return nil
}
