// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_office_365

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

func resourceSentinelDataConnectorOffice365CreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
				return tf.ImportAsExistsError("azurerm_sentinel_data_connector_office_365", id.ID())
			}
		}
	}

	tenantId := d.Get("tenant_id").(string)
	if tenantId == "" {
		tenantId = account.From(meta).TenantId
	}

	exchangeEnabled := d.Get("exchange_enabled").(bool)
	sharePointEnabled := d.Get("sharepoint_enabled").(bool)
	teamsEnabled := d.Get("teams_enabled").(bool)

	// Service will not create the DC in case non of the toggle is enabled.
	if !exchangeEnabled && !sharePointEnabled && !teamsEnabled {
		return errors.New("one of `exchange_enabled`, `sharepoint_enabled` and `teams_enabled` should be `true`")
	}

	exchangeState := dataconnectors.DataTypeStateEnabled
	if !exchangeEnabled {
		exchangeState = dataconnectors.DataTypeStateDisabled
	}

	sharePointState := dataconnectors.DataTypeStateEnabled
	if !sharePointEnabled {
		sharePointState = dataconnectors.DataTypeStateDisabled
	}

	teamsState := dataconnectors.DataTypeStateEnabled
	if !teamsEnabled {
		teamsState = dataconnectors.DataTypeStateDisabled
	}

	param := dataconnectors.OfficeDataConnector{
		Name: &name,
		Properties: &dataconnectors.OfficeDataConnectorProperties{
			TenantId: tenantId,
			DataTypes: dataconnectors.OfficeDataConnectorDataTypes{
				Exchange: dataconnectors.DataConnectorDataTypeCommon{
					State: exchangeState,
				},
				SharePoint: dataconnectors.DataConnectorDataTypeCommon{
					State: sharePointState,
				},
				Teams: dataconnectors.DataConnectorDataTypeCommon{
					State: teamsState,
				},
			},
		},
		Kind: dataconnectors.DataConnectorKindOfficeThreeSixFive,
	}

	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Data Connector Office 365 %q: %+v", id, err)
		}

		if resp.Model == nil {
			return fmt.Errorf("retrieving %s: `model` was nil", id)
		}

		if _, ok := resp.Model.(dataconnectors.OfficeDataConnector); !ok {
			return fmt.Errorf("%s was not an Office 365 Data Connector", id)
		}
	}

	if _, err = client.CreateOrUpdate(ctx, id, param); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceSentinelDataConnectorOffice365Read(d, meta)
}
