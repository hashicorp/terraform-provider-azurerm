// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_aws_cloud_trail

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelDataConnectorAwsCloudTrailCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
