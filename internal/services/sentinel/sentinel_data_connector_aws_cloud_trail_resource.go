// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelDataConnectorAwsCloudTrail() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSentinelDataConnectorAwsCloudTrailCreateUpdate,
		Read:   resourceSentinelDataConnectorAwsCloudTrailRead,
		Update: resourceSentinelDataConnectorAwsCloudTrailCreateUpdate,
		Delete: resourceSentinelDataConnectorAwsCloudTrailDelete,

		Importer: importDataConnectorUntyped(dataconnectors.DataConnectorKindAmazonWebServicesCloudTrail),

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

			"aws_role_arn": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.IsARN,
			},
		},
	}
}

func resourceSentinelDataConnectorAwsCloudTrailCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := workspaces.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	id := dataconnectors.NewDataConnectorID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, name)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			resp, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(resp.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_sentinel_data_connector_aws_cloud_trail", id.ID())
			}
		}
	}

	param := dataconnectors.AwsCloudTrailDataConnector{
		Name: &name,
		Properties: &dataconnectors.AwsCloudTrailDataConnectorProperties{
			AwsRoleArn: pointer.To(d.Get("aws_role_arn").(string)),
			DataTypes: dataconnectors.AwsCloudTrailDataConnectorDataTypes{
				Logs: dataconnectors.DataConnectorDataTypeCommon{
					State: dataconnectors.DataTypeStateEnabled,
				},
			},
		},
		Kind: dataconnectors.DataConnectorKindAmazonWebServicesCloudTrail,
	}

	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.Model == nil {
			return fmt.Errorf("retrieving %s: `model` was nil", id)
		}

		if _, ok := resp.Model.(dataconnectors.AwsCloudTrailDataConnector); !ok {
			return fmt.Errorf("%s was not an AWS Cloud Trail Data Connector", id)
		}
	}

	if _, err = client.CreateOrUpdate(ctx, id, param); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelDataConnectorAwsCloudTrailRead(d, meta)
}

func resourceSentinelDataConnectorAwsCloudTrailRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataconnectors.ParseDataConnectorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	dc, ok := resp.Model.(dataconnectors.AwsCloudTrailDataConnector)
	if !ok {
		return fmt.Errorf("%s was not an AWS Cloud Trail Data Connector", id)
	}

	d.Set("name", id.DataConnectorId)
	d.Set("log_analytics_workspace_id", workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID())
	if props := dc.Properties; props != nil {
		d.Set("aws_role_arn", props.AwsRoleArn)
	}

	return nil
}

func resourceSentinelDataConnectorAwsCloudTrailDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataconnectors.ParseDataConnectorID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
