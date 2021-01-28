package sentinel

import (
	"fmt"
	"log"
	"strconv"
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

func resourceSentinelDataConnectorThreatIntelligenceTaxii() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentinelDataConnectorThreatIntelligenceTaxiiCreateUpdate,
		Read:   resourceSentinelDataConnectorThreatIntelligenceTaxiiRead,
		Update: resourceSentinelDataConnectorThreatIntelligenceTaxiiCreateUpdate,
		Delete: resourceSentinelDataConnectorThreatIntelligenceTaxiiDelete,

		Importer: azSchema.ValidateResourceIDPriorToImportThen(func(id string) error {
			_, err := parse.DataConnectorID(id)
			return err
		}, importSentinelDataConnector(securityinsight.DataConnectorKindThreatIntelligenceTaxii)),

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

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"taxii_server_api_root": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"taxii_server_collection_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"taxii_server_username": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"taxii_server_password": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"tenant_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourceSentinelDataConnectorThreatIntelligenceTaxiiCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	wspClient := meta.(*clients.Client).LogAnalytics.WorkspacesClient
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
				return fmt.Errorf("checking for existing Sentinel Data Connector Threat Intelligence Taxii %q: %+v", id, err)
			}
		}

		id := dataConnectorID(resp.Value)
		if id != nil && *id != "" {
			return tf.ImportAsExistsError("azurerm_sentinel_data_connector_threat_intelligence_taxii", *id)
		}
	}

	workspace, err := wspClient.Get(ctx, workspaceId.ResourceGroup, workspaceId.WorkspaceName)
	if err != nil {
		return fmt.Errorf("retrieving Log Analytics Workspaces %q: %+v", workspaceId, err)
	}
	wspProp := workspace.WorkspaceProperties
	if wspProp == nil {
		return fmt.Errorf("unexpected nil properties of Log Analytics Workspace %q", workspaceId)
	}
	workspaceCustomerId := wspProp.CustomerID

	tenantId := d.Get("tenant_id").(string)
	if tenantId == "" {
		tenantId = meta.(*clients.Client).Account.TenantId
	}

	param := securityinsight.TiTaxiiDataConnector{
		Name: &name,
		TiTaxiiDataConnectorProperties: &securityinsight.TiTaxiiDataConnectorProperties{
			WorkspaceID:  workspaceCustomerId,
			FriendlyName: utils.String(d.Get("display_name").(string)),
			TaxiiServer:  utils.String(d.Get("taxii_server_api_root").(string)),
			CollectionID: utils.String(strconv.Itoa(d.Get("taxii_server_collection_id").(int))),
			TenantID:     &tenantId,
			DataTypes: &securityinsight.TiTaxiiDataConnectorDataTypes{
				TaxiiClient: &securityinsight.TiTaxiiDataConnectorDataTypesTaxiiClient{
					State: securityinsight.Enabled,
				},
			},
		},
		Kind: securityinsight.KindThreatIntelligenceTaxii,
	}

	if username := d.Get("taxii_server_username").(string); username != "" {
		param.TiTaxiiDataConnectorProperties.UserName = &username
	}

	if password := d.Get("taxii_server_password").(string); password != "" {
		param.TiTaxiiDataConnectorProperties.Password = &password
	}

	_, err = client.CreateOrUpdate(ctx, id.ResourceGroup, operationalInsightsResourceProvider, id.WorkspaceName, id.Name, param)
	if err != nil {
		return fmt.Errorf("creating Sentinel Data Connector Threat Intelligence Taxii %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelDataConnectorThreatIntelligenceTaxiiRead(d, meta)
}

func resourceSentinelDataConnectorThreatIntelligenceTaxiiRead(d *schema.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] Sentinel Data Connector Threat Intelligence Taxii %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Sentinel Data Connector Threat Intelligence Taxii %q: %+v", id, err)
	}

	if err := assertDataConnectorKind(resp.Value, securityinsight.DataConnectorKindThreatIntelligenceTaxii); err != nil {
		return fmt.Errorf("asserting Sentinel Data Connector Threat Intelligence Taxii of %q: %+v", id, err)
	}
	dc := resp.Value.(securityinsight.TiTaxiiDataConnector)

	d.Set("name", id.Name)
	d.Set("log_analytics_workspace_id", workspaceId.ID())

	if prop := dc.TiTaxiiDataConnectorProperties; prop != nil {
		d.Set("display_name", prop.FriendlyName)
		d.Set("taxii_server_api_root", prop.TaxiiServer)

		collectionId := 0
		if prop.CollectionID != nil {
			var err error
			collectionId, err = strconv.Atoi(*prop.CollectionID)
			if err != nil {
				return fmt.Errorf("converting `collectionId` of Sentinel Data Connector Threat Intelligence Taxii %q to int: %v", id, err)
			}
		}
		d.Set("taxii_server_collection_id", collectionId)
		d.Set("tenant_id", dc.TenantID)
	}

	return nil
}

func resourceSentinelDataConnectorThreatIntelligenceTaxiiDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataConnectorID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, id.ResourceGroup, operationalInsightsResourceProvider, id.WorkspaceName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Sentinel Data Connector Threat Intelligence Taxii %q: %+v", id, err)
	}

	return nil
}
