// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	securityinsight "github.com/jackofallops/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

func resourceSentinelDataConnectorAwsCloudTrail() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSentinelDataConnectorAwsCloudTrailCreateUpdate,
		Read:   resourceSentinelDataConnectorAwsCloudTrailRead,
		Update: resourceSentinelDataConnectorAwsCloudTrailCreateUpdate,
		Delete: resourceSentinelDataConnectorAwsCloudTrailDelete,

		Importer: importDataConnectorUntyped(securityinsight.DataConnectorKindAmazonWebServicesCloudTrail),

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
	id := parse.NewDataConnectorID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, name)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_sentinel_data_connector_aws_cloud_trail", id.ID())
		}
	}

	param := securityinsight.AwsCloudTrailDataConnector{
		Name: &name,
		AwsCloudTrailDataConnectorProperties: &securityinsight.AwsCloudTrailDataConnectorProperties{
			AwsRoleArn: utils.String(d.Get("aws_role_arn").(string)),
			DataTypes: &securityinsight.AwsCloudTrailDataConnectorDataTypes{
				Logs: &securityinsight.AwsCloudTrailDataConnectorDataTypesLogs{
					State: securityinsight.DataTypeStateEnabled,
				},
			},
		},
		Kind: securityinsight.KindBasicDataConnectorKindAmazonWebServicesCloudTrail,
	}

	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, name)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if _, ok := resp.Value.(securityinsight.AwsCloudTrailDataConnector); !ok {
			return fmt.Errorf("%s was not an AWS Cloud Trail Data Connector", id)
		}
	}

	if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, param); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelDataConnectorAwsCloudTrailRead(d, meta)
}

func resourceSentinelDataConnectorAwsCloudTrailRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

	dc, ok := resp.Value.(securityinsight.AwsCloudTrailDataConnector)
	if !ok {
		return fmt.Errorf("%s was not an AWS Cloud Trail Data Connector", id)
	}

	d.Set("name", id.Name)
	d.Set("log_analytics_workspace_id", workspaceId.ID())
	if prop := dc.AwsCloudTrailDataConnectorProperties; prop != nil {
		d.Set("aws_role_arn", prop.AwsRoleArn)
	}

	return nil
}

func resourceSentinelDataConnectorAwsCloudTrailDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
