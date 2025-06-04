// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	securityinsight "github.com/jackofallops/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

func resourceSentinelDataConnectorMicrosoftCloudAppSecurity() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSentinelDataConnectorMicrosoftCloudAppSecurityCreateUpdate,
		Read:   resourceSentinelDataConnectorMicrosoftCloudAppSecurityRead,
		Update: resourceSentinelDataConnectorMicrosoftCloudAppSecurityCreateUpdate,
		Delete: resourceSentinelDataConnectorMicrosoftCloudAppSecurityDelete,

		Importer: importDataConnectorUntyped(securityinsight.DataConnectorKindMicrosoftCloudAppSecurity),

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
	id := parse.NewDataConnectorID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, name)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_sentinel_data_connector_microsoft_cloud_app_security", id.ID())
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

	alertState := securityinsight.DataTypeStateEnabled
	if !alertsEnabled {
		alertState = securityinsight.DataTypeStateDisabled
	}

	discoveryLogsState := securityinsight.DataTypeStateEnabled
	if !discoveryLogsEnabled {
		discoveryLogsState = securityinsight.DataTypeStateDisabled
	}

	param := securityinsight.MCASDataConnector{
		Name: &name,
		MCASDataConnectorProperties: &securityinsight.MCASDataConnectorProperties{
			TenantID: &tenantId,
			DataTypes: &securityinsight.MCASDataConnectorDataTypes{
				Alerts: &securityinsight.DataConnectorDataTypeCommon{
					State: alertState,
				},
				DiscoveryLogs: &securityinsight.DataConnectorDataTypeCommon{
					State: discoveryLogsState,
				},
			},
		},
		Kind: securityinsight.KindBasicDataConnectorKindMicrosoftCloudAppSecurity,
	}

	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, name)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if _, ok := resp.Value.(securityinsight.MCASDataConnector); !ok {
			return fmt.Errorf("%s was not a Microsoft Cloud App Security Data Connector", id)
		}
	}

	if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, param); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelDataConnectorMicrosoftCloudAppSecurityRead(d, meta)
}

func resourceSentinelDataConnectorMicrosoftCloudAppSecurityRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataConnectorID(d.Id())
	if err != nil {
		return err
	}
	workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	dc, ok := resp.Value.(securityinsight.MCASDataConnector)
	if !ok {
		return fmt.Errorf("%s was not a Microsoft Cloud App Security Data Connector", id)
	}

	d.Set("name", id.Name)
	d.Set("log_analytics_workspace_id", workspaceId.ID())
	d.Set("tenant_id", dc.TenantID)

	var (
		alertsEnabled        bool
		discoveryLogsEnabled bool
	)
	if dt := dc.DataTypes; dt != nil {
		if alert := dt.Alerts; alert != nil {
			alertsEnabled = strings.EqualFold(string(alert.State), string(securityinsight.DataTypeStateEnabled))
		}

		if discoveryLogs := dt.DiscoveryLogs; discoveryLogs != nil {
			discoveryLogsEnabled = strings.EqualFold(string(discoveryLogs.State), string(securityinsight.DataTypeStateEnabled))
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

	id, err := parse.DataConnectorID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
