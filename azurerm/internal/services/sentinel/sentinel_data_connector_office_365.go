package sentinel

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2019-01-01-preview/securityinsight"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	loganalyticsParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	loganalyticsValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sentinel/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSentinelDataConnectorOffice365() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSentinelDataConnectorOffice365CreateUpdate,
		Read:   resourceSentinelDataConnectorOffice365Read,
		Update: resourceSentinelDataConnectorOffice365CreateUpdate,
		Delete: resourceSentinelDataConnectorOffice365Delete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.DataConnectorID(id)
			return err
		}, importSentinelDataConnector(securityinsight.DataConnectorKindOffice365)),

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
				ValidateFunc: loganalyticsValidate.LogAnalyticsWorkspaceID,
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

	workspaceId, err := loganalyticsParse.LogAnalyticsWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	id := parse.NewDataConnectorID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.WorkspaceName, name)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, OperationalInsightsResourceProvider, id.WorkspaceName, name)
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

	// Service avoid concurrent updates of this resource via checking the "etag" to guarantee it is the same value as last Read.
	// TODO: following code can be removed once the issue below is fixed:
	// https://github.com/Azure/azure-rest-api-specs/issues/13203
	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, OperationalInsightsResourceProvider, id.WorkspaceName, name)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Data Connector Office 365 %q: %+v", id, err)
		}

		dc, ok := resp.Value.(securityinsight.OfficeDataConnector)
		if !ok {
			return fmt.Errorf("%s was not an Office 365 Data Connector", id)
		}

		param.Etag = dc.Etag
	}

	if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, OperationalInsightsResourceProvider, id.WorkspaceName, id.Name, param); err != nil {
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
	workspaceId := loganalyticsParse.NewLogAnalyticsWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)

	resp, err := client.Get(ctx, id.ResourceGroup, OperationalInsightsResourceProvider, id.WorkspaceName, id.Name)
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

	if _, err = client.Delete(ctx, id.ResourceGroup, OperationalInsightsResourceProvider, id.WorkspaceName, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
