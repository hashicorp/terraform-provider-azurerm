// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_microsoft_cloud_app_security

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelDataConnectorMicrosoftCloudAppSecurityRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataconnectors.ParseDataConnectorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	dc, ok := resp.Model.(dataconnectors.MCASDataConnector)
	if !ok {
		return fmt.Errorf("%s was not a Microsoft Cloud App Security Data Connector", id)
	}

	d.Set("name", id.DataConnectorId)
	d.Set("log_analytics_workspace_id", workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID())

	var (
		alertsEnabled        bool
		discoveryLogsEnabled bool
	)
	if props := dc.Properties; props != nil {
		d.Set("tenant_id", props.TenantId)

		dt := props.DataTypes
		alertsEnabled = strings.EqualFold(string(dt.Alerts.State), string(dataconnectors.DataTypeStateEnabled))

		if discoveryLogs := dt.DiscoveryLogs; discoveryLogs != nil {
			discoveryLogsEnabled = strings.EqualFold(string(discoveryLogs.State), string(dataconnectors.DataTypeStateEnabled))
		}
	}
	d.Set("discovery_logs_enabled", discoveryLogsEnabled)
	d.Set("alerts_enabled", alertsEnabled)

	return nil
}
