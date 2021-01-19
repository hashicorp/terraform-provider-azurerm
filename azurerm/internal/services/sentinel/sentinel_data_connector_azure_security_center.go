package sentinel

import (
	"fmt"
	"log"
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

func resourceSentinelDataConnectorAzureSecurityCenter() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentinelDataConnectorAzureSecurityCenterCreateUpdate,
		Read:   resourceSentinelDataConnectorAzureSecurityCenterRead,
		Update: resourceSentinelDataConnectorAzureSecurityCenterCreateUpdate,
		Delete: resourceSentinelDataConnectorAzureSecurityCenterDelete,

		Importer: azSchema.ValidateResourceIDPriorToImportThen(func(id string) error {
			_, err := parse.DataConnectorID(id)
			return err
		}, importSentinelDataConnector(securityinsight.DataConnectorKindAzureSecurityCenter)),

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

			"subscription_id": {
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
		},
	}
}

func resourceSentinelDataConnectorAzureSecurityCenterCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
				return fmt.Errorf("checking for existing Sentinel Data Connector Azure Security Center %q: %+v", id, err)
			}
		}

		id := dataConnectorID(resp.Value)
		if id != nil && *id != "" {
			return tf.ImportAsExistsError("azurerm_sentinel_data_connector_azure_security_center", *id)
		}
	}

	subscriptionId := d.Get("subscription_id").(string)
	if subscriptionId == "" {
		subscriptionId = id.SubscriptionId
	}

	alertState := securityinsight.Enabled
	if !d.Get("alerts_enabled").(bool) {
		alertState = securityinsight.Disabled
	}

	param := securityinsight.ASCDataConnector{
		Name: &name,
		ASCDataConnectorProperties: &securityinsight.ASCDataConnectorProperties{
			SubscriptionID: &subscriptionId,
			DataTypes: &securityinsight.AlertsDataTypeOfDataConnector{
				Alerts: &securityinsight.AlertsDataTypeOfDataConnectorAlerts{
					State: alertState,
				},
			},
		},
		Kind: securityinsight.KindAzureSecurityCenter,
	}

	_, err = client.CreateOrUpdate(ctx, id.ResourceGroup, operationalInsightsResourceProvider, id.WorkspaceName, id.Name, param)
	if err != nil {
		return fmt.Errorf("creating Sentinel Data Connector Azure Security Center %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelDataConnectorAzureSecurityCenterRead(d, meta)
}

func resourceSentinelDataConnectorAzureSecurityCenterRead(d *schema.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] Sentinel Data Connector Azure Security Center %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Sentinel Data Connector Azure Security Center %q: %+v", id, err)
	}

	if err := assertDataConnectorKind(resp.Value, securityinsight.DataConnectorKindAzureSecurityCenter); err != nil {
		return fmt.Errorf("asserting Sentinel Data Connector Azure Security Center of %q: %+v", id, err)
	}
	dc := resp.Value.(securityinsight.ASCDataConnector)

	d.Set("name", id.Name)
	d.Set("log_analytics_workspace_id", workspaceId.ID())
	d.Set("subscription_id", dc.SubscriptionID)
	if dt := dc.DataTypes; dt != nil {
		if alert := dt.Alerts; alert != nil {
			d.Set("alerts_enabled", alert.State == securityinsight.Enabled)
		}
	}

	return nil
}

func resourceSentinelDataConnectorAzureSecurityCenterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataConnectorID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, id.ResourceGroup, operationalInsightsResourceProvider, id.WorkspaceName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Sentinel Data Connector Azure Security Center %q: %+v", id, err)
	}

	return nil
}
