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

func resourceSentinelDataConnectorOffice365() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentinelDataConnectorOffice365CreateUpdate,
		Read:   resourceSentinelDataConnectorOffice365Read,
		Update: resourceSentinelDataConnectorOffice365CreateUpdate,
		Delete: resourceSentinelDataConnectorOffice365Delete,

		Importer: azSchema.ValidateResourceIDPriorToImportThen(func(id string) error {
			_, err := parse.DataConnectorID(id)
			return err
		}, importSentinelDataConnector(securityinsight.DataConnectorKindOffice365)),

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
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"exchange_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"share_point_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"teams_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceSentinelDataConnectorOffice365CreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
				return fmt.Errorf("checking for existing Sentinel Data Connector Office 365 %q: %+v", id, err)
			}
		}

		id := dataConnectorID(resp.Value)
		if id != nil && *id != "" {
			return tf.ImportAsExistsError("azurerm_sentinel_data_connector_office_365", *id)
		}
	}

	tenantId := d.Get("tenant_id").(string)
	if tenantId == "" {
		tenantId = meta.(*clients.Client).Account.TenantId
	}

	exchangeEnabled := d.Get("exchange_enabled").(bool)
	sharePointEnabled := d.Get("share_point_enabled").(bool)
	teamsEnabled := d.Get("teams_enabled").(bool)

	// Service will not create the DC in case non of the toggle is enabled.
	if !exchangeEnabled && !sharePointEnabled && !teamsEnabled {
		return fmt.Errorf("one of `exchange_enabled`, `share_point_enabled` and `teams_enabled` should be `true`")
	}

	exchangeState := securityinsight.Enabled
	if !exchangeEnabled {
		exchangeState = securityinsight.Disabled
	}

	sharePointState := securityinsight.Enabled
	if !sharePointEnabled {
		sharePointState = securityinsight.Disabled
	}

	teamsState := securityinsight.Enabled
	if !teamsEnabled {
		teamsState = securityinsight.Disabled
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
		Kind: securityinsight.KindOffice365,
	}

	// Service avoid concurrent updates of this resource via checking the "etag" to guarantee it is the same value as last Read.
	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, operationalInsightsResourceProvider, id.WorkspaceName, name)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Data Connector Office 365 %q: %+v", id, err)
		}

		if err := assertDataConnectorKind(resp.Value, securityinsight.DataConnectorKindOffice365); err != nil {
			return fmt.Errorf("asserting Sentinel Data Connector of %q: %+v", id, err)
		}
		param.Etag = resp.Value.(securityinsight.OfficeDataConnector).Etag
	}

	_, err = client.CreateOrUpdate(ctx, id.ResourceGroup, operationalInsightsResourceProvider, id.WorkspaceName, id.Name, param)
	if err != nil {
		return fmt.Errorf("creating Sentinel Data Connector Office 365 %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelDataConnectorOffice365Read(d, meta)
}

func resourceSentinelDataConnectorOffice365Read(d *schema.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] Sentinel Data Connector Office 365 %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Sentinel Data Connector Office 365 %q: %+v", id, err)
	}

	if err := assertDataConnectorKind(resp.Value, securityinsight.DataConnectorKindOffice365); err != nil {
		return fmt.Errorf("asserting Sentinel Data Connector Office 365 of %q: %+v", id, err)
	}
	dc := resp.Value.(securityinsight.OfficeDataConnector)

	d.Set("name", id.Name)
	d.Set("log_analytics_workspace_id", workspaceId.ID())
	d.Set("tenant_id", dc.TenantID)

	if dt := dc.DataTypes; dt != nil {
		if exchange := dt.Exchange; exchange != nil {
			d.Set("exchange_enabled", strings.EqualFold(string(exchange.State), string(securityinsight.Enabled)))
		}
		if sharePoint := dt.SharePoint; sharePoint != nil {
			d.Set("share_point_enabled", strings.EqualFold(string(sharePoint.State), string(securityinsight.Enabled)))
		}
		if teams := dt.Teams; teams != nil {
			d.Set("teams_enabled", strings.EqualFold(string(teams.State), string(securityinsight.Enabled)))
		}
	}

	return nil
}

func resourceSentinelDataConnectorOffice365Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataConnectorID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, id.ResourceGroup, operationalInsightsResourceProvider, id.WorkspaceName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Sentinel Data Connector Office 365 %q: %+v", id, err)
	}

	return nil
}
