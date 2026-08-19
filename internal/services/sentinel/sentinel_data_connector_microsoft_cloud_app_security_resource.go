// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelDataConnectorMicrosoftCloudAppSecurity() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSentinelDataConnectorMicrosoftCloudAppSecurityCreateUpdate,
		Read:   resourceSentinelDataConnectorMicrosoftCloudAppSecurityRead,
		Update: resourceSentinelDataConnectorMicrosoftCloudAppSecurityCreateUpdate,
		Delete: resourceSentinelDataConnectorMicrosoftCloudAppSecurityDelete,

		Importer: importDataConnectorUntyped(dataconnectors.DataConnectorKindMicrosoftCloudAppSecurity),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"log_analytics_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
			},

			"alerts_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"discovery_logs_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceSentinelDataConnectorMicrosoftCloudAppSecurityCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := workspaces.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	id := dataconnectors.NewDataConnectorID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, name)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
		tenantId = meta.(*clients.Client).Account.TenantId
	}

	alertsEnabled := d.Get("alerts_enabled").(bool)
	discoveryLogsEnabled := d.Get("discovery_logs_enabled").(bool)

	// Service will not create the DC in case non of the toggle is enabled.
	if !alertsEnabled && !discoveryLogsEnabled {
		return fmt.Errorf("either `alerts_enabled` or `discovery_logs_enabled` should be `true`")
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

func resourceSentinelDataConnectorMicrosoftCloudAppSecurityRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceSentinelDataConnectorMicrosoftCloudAppSecurityDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataconnectors.ParseDataConnectorID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
