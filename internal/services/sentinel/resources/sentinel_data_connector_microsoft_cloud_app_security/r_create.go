// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_microsoft_cloud_app_security

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	sentinelClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelDataConnectorMicrosoftCloudAppSecurityCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := sentinelClient.From(meta).DataConnectorsClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	workspaceId, err := workspaces.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	id := dataconnectors.NewDataConnectorID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, name)

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			resp, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(resp.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_sentinel_data_connector_microsoft_cloud_app_security", id.ID())
			}
		}
	}

	tenantId := d.Get("tenant_id").(string)
	if tenantId == "" {
		tenantId = account.From(meta).TenantId
	}

	alertsEnabled := d.Get("alerts_enabled").(bool)
	discoveryLogsEnabled := d.Get("discovery_logs_enabled").(bool)

	// Service will not create the DC in case non of the toggle is enabled.
	if !alertsEnabled && !discoveryLogsEnabled {
		return errors.New("either `alerts_enabled` or `discovery_logs_enabled` should be `true`")
	}

	alertState := dataconnectors.DataTypeStateEnabled
	if !alertsEnabled {
		alertState = dataconnectors.DataTypeStateDisabled
	}

	discoveryLogsState := dataconnectors.DataTypeStateEnabled
	if !discoveryLogsEnabled {
		discoveryLogsState = dataconnectors.DataTypeStateDisabled
	}

	param := dataconnectors.MCASDataConnector{
		Name: &name,
		Properties: &dataconnectors.MCASDataConnectorProperties{
			TenantId: tenantId,
			DataTypes: dataconnectors.MCASDataConnectorDataTypes{
				Alerts: dataconnectors.DataConnectorDataTypeCommon{
					State: alertState,
				},
				DiscoveryLogs: &dataconnectors.DataConnectorDataTypeCommon{
					State: discoveryLogsState,
				},
			},
		},
		Kind: dataconnectors.DataConnectorKindMicrosoftCloudAppSecurity,
	}

	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.Model == nil {
			return fmt.Errorf("retrieving %s: `model` was nil", id)
		}

		if _, ok := resp.Model.(dataconnectors.MCASDataConnector); !ok {
			return fmt.Errorf("%s was not a Microsoft Cloud App Security Data Connector", id)
		}
	}

	if _, err = client.CreateOrUpdate(ctx, id, param); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceSentinelDataConnectorMicrosoftCloudAppSecurityRead(d, meta)
}
