package sentinel

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2019-01-01-preview/securityinsight"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	loganalyticsParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	loganalyticsValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sentinel/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSentinelDataConnectorMicrosoftCloudAppSecurity() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentinelDataConnectorMicrosoftCloudAppSecurityCreateUpdate,
		Read:   resourceSentinelDataConnectorMicrosoftCloudAppSecurityRead,
		Update: resourceSentinelDataConnectorMicrosoftCloudAppSecurityCreateUpdate,
		Delete: resourceSentinelDataConnectorMicrosoftCloudAppSecurityDelete,

		Importer: azSchema.ValidateResourceIDPriorToImportThen(func(id string) error {
			_, err := parse.DataConnectorID(id)
			return err
		}, importSentinelDataConnector(securityinsight.DataConnectorKindMicrosoftCloudAppSecurity)),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"log_analytics_workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: loganalyticsValidate.LogAnalyticsWorkspaceID,
			},

			"tenant_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
			},

			"alerts_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"discovery_logs_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceSentinelDataConnectorMicrosoftCloudAppSecurityCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := loganalyticsParse.LogAnalyticsWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	id := parse.NewDataConnectorID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.WorkspaceName, name)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, operationalInsightsResourceProvider, id.WorkspaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing Sentinel Data Connector Microsoft Cloud App Security %q: %+v", id, err)
			}
		}

		id := dataConnectorID(resp.Value)
		if id != nil && *id != "" {
			return tf.ImportAsExistsError("azurerm_sentinel_data_connector_microsoft_cloud_app_security", *id)
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

	alertState := securityinsight.Enabled
	if !alertsEnabled {
		alertState = securityinsight.Disabled
	}

	discoveryLogsState := securityinsight.Enabled
	if !discoveryLogsEnabled {
		discoveryLogsState = securityinsight.Disabled
	}

	param := securityinsight.MCASDataConnector{
		Name: &name,
		MCASDataConnectorProperties: &securityinsight.MCASDataConnectorProperties{
			TenantID: &tenantId,
			DataTypes: &securityinsight.MCASDataConnectorDataTypes{
				Alerts: &securityinsight.AlertsDataTypeOfDataConnectorAlerts{
					State: alertState,
				},
				DiscoveryLogs: &securityinsight.MCASDataConnectorDataTypesDiscoveryLogs{
					State: discoveryLogsState,
				},
			},
		},
		Kind: securityinsight.KindMicrosoftCloudAppSecurity,
	}

	// Service avoid concurrent updates of this resource via checking the "etag" to guarantee it is the same value as last Read.
	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, operationalInsightsResourceProvider, id.WorkspaceName, name)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Data Connector Microsoft Cloud App Security %q: %+v", id, err)
		}

		if err := assertDataConnectorKind(resp.Value, securityinsight.DataConnectorKindMicrosoftCloudAppSecurity); err != nil {
			return fmt.Errorf("asserting Sentinel Data Connector of %q: %+v", id, err)
		}
		param.Etag = resp.Value.(securityinsight.MCASDataConnector).Etag
	}

	_, err = client.CreateOrUpdate(ctx, id.ResourceGroup, operationalInsightsResourceProvider, id.WorkspaceName, id.Name, param)
	if err != nil {
		return fmt.Errorf("creating Sentinel Data Connector Microsoft Cloud App Security %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelDataConnectorMicrosoftCloudAppSecurityRead(d, meta)
}

func resourceSentinelDataConnectorMicrosoftCloudAppSecurityRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataConnectorID(d.Id())
	if err != nil {
		return err
	}
	workspaceId := loganalyticsParse.NewLogAnalyticsWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)

	resp, err := client.Get(ctx, id.ResourceGroup, operationalInsightsResourceProvider, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Sentinel Data Connector Microsoft Cloud App Security %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Sentinel Data Connector Microsoft Cloud App Security %q: %+v", id, err)
	}

	if err := assertDataConnectorKind(resp.Value, securityinsight.DataConnectorKindMicrosoftCloudAppSecurity); err != nil {
		return fmt.Errorf("asserting Sentinel Data Connector Microsoft Cloud App Security of %q: %+v", id, err)
	}
	dc := resp.Value.(securityinsight.MCASDataConnector)

	d.Set("name", id.Name)
	d.Set("log_analytics_workspace_id", workspaceId.ID())
	d.Set("tenant_id", dc.TenantID)
	if dt := dc.DataTypes; dt != nil {
		if alert := dt.Alerts; alert != nil {
			d.Set("alerts_enabled", strings.EqualFold(string(alert.State), string(securityinsight.Enabled)))
		}
		if discoveryLogs := dt.DiscoveryLogs; discoveryLogs != nil {
			d.Set("discovery_logs_enabled", strings.EqualFold(string(discoveryLogs.State), string(securityinsight.Enabled)))
		}
	}

	return nil
}

func resourceSentinelDataConnectorMicrosoftCloudAppSecurityDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataConnectorID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, id.ResourceGroup, operationalInsightsResourceProvider, id.WorkspaceName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Sentinel Data Connector Microsoft Cloud App Security %q: %+v", id, err)
	}

	return nil
}
