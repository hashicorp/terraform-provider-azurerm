// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_deployment_setting

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/deploymentsettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.DeploymentSettings

			id, err := deploymentsettings.ParseDeploymentSettingID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			schema := StackHCIDeploymentSettingModel{
				StackHCIClusterId: deploymentsettings.NewClusterID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					schema.ArcResourceIds = props.ArcNodeResourceIds
					schema.Version = pointer.From(props.DeploymentConfiguration.Version)
					schema.ScaleUnit = flattenDeploymentSettingScaleUnits(props.DeploymentConfiguration.ScaleUnits)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}
