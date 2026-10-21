// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_deployment

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.NetworkManagers
			id, err := parse.NetworkManagerDeploymentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			listParam := networkmanagers.NetworkManagerDeploymentStatusParameter{
				Regions:         &[]string{id.Location},
				DeploymentTypes: &[]networkmanagers.ConfigurationType{networkmanagers.ConfigurationType(id.ScopeAccess)},
			}

			networkManagerId := networkmanagers.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName)

			resp, err := client.NetworkManagerDeploymentStatusList(ctx, networkManagerId, listParam, networkmanagers.DefaultNetworkManagerDeploymentStatusListOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("unexpected null model of %s", *id)
			}

			if resp.Model.Value == nil || len(*resp.Model.Value) == 0 || (*resp.Model.Value)[0].ConfigurationIds == nil || len(*(*resp.Model.Value)[0].ConfigurationIds) == 0 {
				return metadata.MarkAsGone(id)
			}

			if len(*resp.Model.Value) > 1 {
				return fmt.Errorf("found more than one deployment with id %s", *id)
			}

			deployment := (*resp.Model.Value)[0]
			return metadata.Encode(&ManagerDeploymentModel{
				NetworkManagerId: networkmanagers.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName).ID(),
				Location:         location.NormalizeNilable(deployment.Region),
				ScopeAccess:      string(*deployment.DeploymentType),
				ConfigurationIds: *deployment.ConfigurationIds,
			})
		},
		Timeout: 5 * time.Minute,
	}
}

func resourceManagerDeploymentWaitForDeleted(ctx context.Context, client *networkmanagers.NetworkManagersClient, managerDeploymentId *parse.ManagerDeploymentId, d time.Duration) error {
	state := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"NotStarted", "Deploying", "Deployed", "Failed"},
		Target:     []string{"NotFound"},
		Refresh:    resourceManagerDeploymentResultRefreshFunc(ctx, client, managerDeploymentId),
		Timeout:    d,
	}

	if _, err := state.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for the %s: %+v", *managerDeploymentId, err)
	}

	return nil
}

func resourceManagerDeploymentWaitForFinished(ctx context.Context, client *networkmanagers.NetworkManagersClient, managerDeploymentId *parse.ManagerDeploymentId, d time.Duration) error {
	state := &pluginsdk.StateChangeConf{
		MinTimeout:     30 * time.Second,
		Delay:          10 * time.Second,
		Pending:        []string{"NotStarted", "Deploying"},
		Target:         []string{"Deployed"},
		NotFoundChecks: 20,
		Timeout:        d,
		Refresh: func() (any, string, error) {
			result, state, err := resourceManagerDeploymentResultRefreshFunc(ctx, client, managerDeploymentId)()
			if state == "NotFound" {
				// the deployment might not found after initial commit, https://github.com/Azure/azure-rest-api-specs/issues/27327
				// to serve NotFoundChecks, return nil result
				return nil, state, err
			}
			return result, state, err
		},
	}

	if _, err := state.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for the %s: %+v", *managerDeploymentId, err)
	}

	return nil
}

func resourceManagerDeploymentResultRefreshFunc(ctx context.Context, client *networkmanagers.NetworkManagersClient, id *parse.ManagerDeploymentId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		listParam := networkmanagers.NetworkManagerDeploymentStatusParameter{
			Regions:         &[]string{location.Normalize(id.Location)},
			DeploymentTypes: &[]networkmanagers.ConfigurationType{networkmanagers.ConfigurationType(id.ScopeAccess)},
		}

		networkManagerId := networkmanagers.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName)

		resp, err := client.NetworkManagerDeploymentStatusList(ctx, networkManagerId, listParam, networkmanagers.DefaultNetworkManagerDeploymentStatusListOperationOptions())
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return resp, "NotFound", nil
			}
			return resp, "Error", fmt.Errorf("retrieving Deployment: %+v", err)
		}

		if resp.Model == nil {
			return resp, "Error", fmt.Errorf("unexpected null model of %s", *id)
		}

		if resp.Model.Value == nil || len(*resp.Model.Value) == 0 || *(*resp.Model.Value)[0].ConfigurationIds == nil || len(*(*resp.Model.Value)[0].ConfigurationIds) == 0 {
			log.Printf("[DEBUG] retrieving %s: listing deployments succeeds however the specific deployment was not found", *id)
			return resp, "NotFound", nil
		}

		if len(*resp.Model.Value) > 1 {
			return resp, "Error", fmt.Errorf("found more than one deployment with id %s", *id)
		}

		deploymentStatus := string(*(*resp.Model.Value)[0].DeploymentStatus)
		return resp, deploymentStatus, nil
	}
}
