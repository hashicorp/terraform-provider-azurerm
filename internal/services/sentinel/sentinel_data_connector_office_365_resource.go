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

func resourceSentinelDataConnectorOffice365() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSentinelDataConnectorOffice365CreateUpdate,
		Read:   resourceSentinelDataConnectorOffice365Read,
		Update: resourceSentinelDataConnectorOffice365CreateUpdate,
		Delete: resourceSentinelDataConnectorOffice365Delete,

		Importer: importDataConnectorUntyped(securityinsight.DataConnectorKindOffice365),

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
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"exchange_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"sharepoint_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"teams_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceSentinelDataConnectorOffice365CreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
			return tf.ImportAsExistsError("azurerm_sentinel_data_connector_office_365", id.ID())
		}
	}

	tenantId := d.Get("tenant_id").(string)
	if tenantId == "" {
		tenantId = meta.(*clients.Client).Account.TenantId
	}

	exchangeEnabled := d.Get("exchange_enabled").(bool)
	sharePointEnabled := d.Get("sharepoint_enabled").(bool)
	teamsEnabled := d.Get("teams_enabled").(bool)

	// Service will not create the DC in case non of the toggle is enabled.
	if !exchangeEnabled && !sharePointEnabled && !teamsEnabled {
		return fmt.Errorf("one of `exchange_enabled`, `sharepoint_enabled` and `teams_enabled` should be `true`")
	}

	exchangeState := securityinsight.DataTypeStateEnabled
	if !exchangeEnabled {
		exchangeState = securityinsight.DataTypeStateDisabled
	}

	sharePointState := securityinsight.DataTypeStateEnabled
	if !sharePointEnabled {
		sharePointState = securityinsight.DataTypeStateDisabled
	}

	teamsState := securityinsight.DataTypeStateEnabled
	if !teamsEnabled {
		teamsState = securityinsight.DataTypeStateDisabled
	}

	param := securityinsight.OfficeDataConnector{
		Name: &name,
		OfficeDataConnectorProperties: &securityinsight.OfficeDataConnectorProperties{
			TenantID: &tenantId,
			DataTypes: &securityinsight.OfficeDataConnectorDataTypes{
				Exchange: &securityinsight.OfficeDataConnectorDataTypesExchange{
					State: exchangeState,
				},
				SharePoint: &securityinsight.OfficeDataConnectorDataTypesSharePoint{
					State: sharePointState,
				},
				Teams: &securityinsight.OfficeDataConnectorDataTypesTeams{
					State: teamsState,
				},
			},
		},
		Kind: securityinsight.KindBasicDataConnectorKindOffice365,
	}

	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, name)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Data Connector Office 365 %q: %+v", id, err)
		}

		if _, ok := resp.Value.(securityinsight.OfficeDataConnector); !ok {
			return fmt.Errorf("%s was not an Office 365 Data Connector", id)
		}
	}

	if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, param); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelDataConnectorOffice365Read(d, meta)
}

func resourceSentinelDataConnectorOffice365Read(d *pluginsdk.ResourceData, meta interface{}) error {
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

	dc, ok := resp.Value.(securityinsight.OfficeDataConnector)
	if !ok {
		return fmt.Errorf("%s was not an Office 365 Data Connector", id)
	}

	d.Set("name", id.Name)
	d.Set("log_analytics_workspace_id", workspaceId.ID())
	d.Set("tenant_id", dc.TenantID)

	if dt := dc.DataTypes; dt != nil {
		exchangeEnabled := false
		if exchange := dt.Exchange; exchange != nil {
			exchangeEnabled = strings.EqualFold(string(exchange.State), string(securityinsight.DataTypeStateEnabled))
		}
		d.Set("exchange_enabled", exchangeEnabled)

		sharePointEnabled := false
		if sharePoint := dt.SharePoint; sharePoint != nil {
			sharePointEnabled = strings.EqualFold(string(sharePoint.State), string(securityinsight.DataTypeStateEnabled))
		}
		d.Set("sharepoint_enabled", sharePointEnabled)

		teamsEnabled := false
		if teams := dt.Teams; teams != nil {
			teamsEnabled = strings.EqualFold(string(teams.State), string(securityinsight.DataTypeStateEnabled))
		}
		d.Set("teams_enabled", teamsEnabled)
	}

	return nil
}

func resourceSentinelDataConnectorOffice365Delete(d *pluginsdk.ResourceData, meta interface{}) error {
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
