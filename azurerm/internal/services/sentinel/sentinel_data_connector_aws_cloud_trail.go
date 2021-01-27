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

func resourceSentinelDataConnectorAwsCloudTrail() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentinelDataConnectorAwsCloudTrailCreateUpdate,
		Read:   resourceSentinelDataConnectorAwsCloudTrailRead,
		Update: resourceSentinelDataConnectorAwsCloudTrailCreateUpdate,
		Delete: resourceSentinelDataConnectorAwsCloudTrailDelete,

		Importer: azSchema.ValidateResourceIDPriorToImportThen(func(id string) error {
			_, err := parse.DataConnectorID(id)
			return err
		}, importSentinelDataConnector(securityinsight.DataConnectorKindAmazonWebServicesCloudTrail)),

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

			"aws_role_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceSentinelDataConnectorAwsCloudTrailCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
				return fmt.Errorf("checking for existing Sentinel Data Connector AWS Cloud Trail %q: %+v", id, err)
			}
		}

		id := dataConnectorID(resp.Value)
		if id != nil && *id != "" {
			return tf.ImportAsExistsError("azurerm_sentinel_data_connector_aws_cloud_trail", *id)
		}
	}

	param := securityinsight.AwsCloudTrailDataConnector{
		Name: &name,
		AwsCloudTrailDataConnectorProperties: &securityinsight.AwsCloudTrailDataConnectorProperties{
			AwsRoleArn: utils.String(d.Get("aws_role_arn").(string)),
			DataTypes: &securityinsight.AwsCloudTrailDataConnectorDataTypes{
				Logs: &securityinsight.AwsCloudTrailDataConnectorDataTypesLogs{
					State: securityinsight.Enabled,
				},
			},
		},
		Kind: securityinsight.KindAmazonWebServicesCloudTrail,
	}

	// Service avoid concurrent updates of this resource via checking the "etag" to guarantee it is the same value as last Read.
	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, operationalInsightsResourceProvider, id.WorkspaceName, name)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Data Connector AWS Cloud Trail %q: %+v", id, err)
		}

		if err := assertDataConnectorKind(resp.Value, securityinsight.DataConnectorKindAmazonWebServicesCloudTrail); err != nil {
			return fmt.Errorf("asserting Sentinel Data Connector of %q: %+v", id, err)
		}
		param.Etag = resp.Value.(securityinsight.AwsCloudTrailDataConnector).Etag
	}

	_, err = client.CreateOrUpdate(ctx, id.ResourceGroup, operationalInsightsResourceProvider, id.WorkspaceName, id.Name, param)
	if err != nil {
		return fmt.Errorf("creating Sentinel Data Connector AWS Cloud Trail %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelDataConnectorAwsCloudTrailRead(d, meta)
}

func resourceSentinelDataConnectorAwsCloudTrailRead(d *schema.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] Sentinel Data Connector AWS Cloud Trail %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Sentinel Data Connector AWS Cloud Trail %q: %+v", id, err)
	}

	if err := assertDataConnectorKind(resp.Value, securityinsight.DataConnectorKindAmazonWebServicesCloudTrail); err != nil {
		return fmt.Errorf("asserting Sentinel Data Connector AWS Cloud Trail of %q: %+v", id, err)
	}
	dc := resp.Value.(securityinsight.AwsCloudTrailDataConnector)

	d.Set("name", id.Name)
	d.Set("log_analytics_workspace_id", workspaceId.ID())
	if prop := dc.AwsCloudTrailDataConnectorProperties; prop != nil {
		d.Set("aws_role_arn", prop.AwsRoleArn)
	}

	return nil
}

func resourceSentinelDataConnectorAwsCloudTrailDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataConnectorID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, id.ResourceGroup, operationalInsightsResourceProvider, id.WorkspaceName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Sentinel Data Connector AWS Cloud Trail %q: %+v", id, err)
	}

	return nil
}
