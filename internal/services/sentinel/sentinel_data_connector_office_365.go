package sentinel

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelDataConnectorOffice365() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSentinelDataConnectorOffice365CreateUpdate,
		Read:   resourceSentinelDataConnectorOffice365Read,
		Update: resourceSentinelDataConnectorOffice365CreateUpdate,
		Delete: resourceSentinelDataConnectorOffice365Delete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := dataconnectors.ParseDataConnectorID(id)
			return err
		}, importSentinelDataConnector(dataconnectors.DataConnectorKindOfficeThreeSixFive)),

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
				ValidateFunc: dataconnectors.ValidateWorkspaceID,
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

	workspaceId, err := dataconnectors.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	id := dataconnectors.NewDataConnectorID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, name)

	if d.IsNewResource() {
		resp, err := client.DataConnectorsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
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
	}

	if !d.IsNewResource() {
		resp, err := client.DataConnectorsGet(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Data Connector Office 365 %q: %+v", id, err)
		}

		if resp.Model == nil {
			return fmt.Errorf("model was nil for %s", id)
		}

		modelPtr := *resp.Model
		if _, ok := modelPtr.(dataconnectors.OfficeDataConnector); !ok {
			return fmt.Errorf("%s was not an Office 365 Data Connector", id)
		}
	}

	if _, err = client.DataConnectorsCreateOrUpdate(ctx, id, param); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelDataConnectorOffice365Read(d, meta)
}

func resourceSentinelDataConnectorOffice365Read(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataconnectors.ParseDataConnectorID(d.Id())
	if err != nil {
		return err
	}
	workspaceId := dataconnectors.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)

	resp, err := client.DataConnectorsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("model was nil for %s", id)
	}

	modelPtr := *resp.Model
	dc, ok := modelPtr.(dataconnectors.OfficeDataConnector)
	if !ok {
		return fmt.Errorf("%s was not an Office 365 Data Connector", id)
	}

	d.Set("name", id.DataConnectorId)
	d.Set("log_analytics_workspace_id", workspaceId.ID())

	if props := dc.Properties; props != nil {
		d.Set("tenant_id", props.TenantId)

		dt := props.DataTypes
		d.Set("exchange_enabled", strings.EqualFold(string(dt.Exchange.State), string(dataconnectors.DataTypeStateEnabled)))
		d.Set("sharepoint_enabled", strings.EqualFold(string(dt.SharePoint.State), string(dataconnectors.DataTypeStateEnabled)))
		d.Set("teams_enabled", strings.EqualFold(string(dt.Teams.State), string(dataconnectors.DataTypeStateEnabled)))
	}

	return nil
}

func resourceSentinelDataConnectorOffice365Delete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataconnectors.ParseDataConnectorID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.DataConnectorsDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
