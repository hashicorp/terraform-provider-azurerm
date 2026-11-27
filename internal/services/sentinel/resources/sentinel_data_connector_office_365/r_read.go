// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_office_365

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	sentinelClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelDataConnectorOffice365Read(d *pluginsdk.ResourceData, meta any) error {
	client := sentinelClient.From(meta).DataConnectorsClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
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

	dc, ok := resp.Model.(dataconnectors.OfficeDataConnector)
	if !ok {
		return fmt.Errorf("%s was not an Office 365 Data Connector", id)
	}

	d.Set("name", id.DataConnectorId)
	d.Set("log_analytics_workspace_id", workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID())

	if props := dc.Properties; props != nil {
		d.Set("tenant_id", props.TenantId)

		dt := props.DataTypes
		d.Set("exchange_enabled", strings.EqualFold(string(dt.Exchange.State), string(dataconnectors.DataTypeStateEnabled)))

		d.Set("sharepoint_enabled", strings.EqualFold(string(dt.SharePoint.State), string(dataconnectors.DataTypeStateEnabled)))

		d.Set("teams_enabled", strings.EqualFold(string(dt.Teams.State), string(dataconnectors.DataTypeStateEnabled)))
	}

	return nil
}
