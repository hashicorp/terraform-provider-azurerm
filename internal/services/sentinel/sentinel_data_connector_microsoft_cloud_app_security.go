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
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

func resourceSentinelDataConnectorMicrosoftCloudAppSecurity() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSentinelDataConnectorMicrosoftCloudAppSecurityCreateUpdate,
		Read:   resourceSentinelDataConnectorMicrosoftCloudAppSecurityRead,
		Update: resourceSentinelDataConnectorMicrosoftCloudAppSecurityCreateUpdate,
		Delete: resourceSentinelDataConnectorMicrosoftCloudAppSecurityDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := dataconnectors.ParseDataConnectorID(id)
			return err
		}, importSentinelDataConnector(dataconnectors.DataConnectorKindMicrosoftCloudAppSecurity)),

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
	}

	if !d.IsNewResource() {
		resp, err := client.DataConnectorsGet(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.Model == nil {
			return fmt.Errorf("model was nil for %s", id)
		}

		modelPtr := *resp.Model
		if _, ok := modelPtr.(securityinsight.MCASDataConnector); !ok {
			return fmt.Errorf("%s was not a Microsoft Cloud App Security Data Connector", id)
		}
	}

	if _, err = client.DataConnectorsCreateOrUpdate(ctx, id, param); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

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
	dc, ok := modelPtr.(dataconnectors.MCASDataConnector)
	if !ok {
		return fmt.Errorf("%s was not a Microsoft Cloud App Security Data Connector", id)
	}

	d.Set("name", id.DataConnectorId)
	d.Set("log_analytics_workspace_id", workspaceId.ID())

	if props := dc.Properties; props != nil {
		d.Set("tenant_id", props.TenantId)

		var (
			alertsEnabled        bool
			discoveryLogsEnabled bool
		)
		dt := props.DataTypes

		alertsEnabled = strings.EqualFold(string(dt.Alerts.State), string(securityinsight.DataTypeStateEnabled))

		if discoveryLogs := dt.DiscoveryLogs; discoveryLogs != nil {
			discoveryLogsEnabled = strings.EqualFold(string(discoveryLogs.State), string(securityinsight.DataTypeStateEnabled))

		}
		d.Set("discovery_logs_enabled", discoveryLogsEnabled)
		d.Set("alerts_enabled", alertsEnabled)
	}

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

	if _, err = client.DataConnectorsDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
