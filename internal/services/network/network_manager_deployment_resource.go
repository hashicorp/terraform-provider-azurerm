// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/networkmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = ManagerDeploymentResource{}

type ManagerDeploymentModel struct {
	NetworkManagerId string   `tfschema:"network_manager_id"`
	ScopeAccess      string   `tfschema:"scope_access"`
	Location         string   `tfschema:"location"`
	ConfigurationIds []string `tfschema:"configuration_ids"`
}

type ManagerDeploymentResource struct{}

func (r ManagerDeploymentResource) ResourceType() string {
	return "azurerm_network_manager_deployment"
}

func (r ManagerDeploymentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.NetworkManagerDeploymentID
}

func (r ManagerDeploymentResource) ModelObject() interface{} {
	return &ManagerDeploymentModel{}
}

func (r ManagerDeploymentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"network_manager_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: networkmanagers.ValidateNetworkManagerID,
		},

		"location": commonschema.Location(),

		"scope_access": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(networkmanagers.ConfigurationTypeConnectivity),
				string(networkmanagers.ConfigurationTypeSecurityAdmin),
				string(networkmanagers.ConfigurationTypeRouting),
			}, false),
		},

		"configuration_ids": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azure.ValidateResourceID,
			},
		},

		// TODO: look at removing this workaround in v4.0, see https://github.com/hashicorp/terraform-provider-azurerm/pull/20451#discussion_r1179646861 (manicminer)
		"triggers": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r ManagerDeploymentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagerDeploymentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			metadata.Logger.Info("Decoding state..")
			var state ManagerDeploymentModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Network.NetworkManagers

			networkManagerId, err := networkmanagers.ParseNetworkManagerID(state.NetworkManagerId)
			if err != nil {
				return err
			}

			normalizedLocation := azure.NormalizeLocation(state.Location)
			id := parse.NewNetworkManagerDeploymentID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroupName, networkManagerId.NetworkManagerName, normalizedLocation, state.ScopeAccess)

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			metadata.Logger.Infof("creating %s", *id)

			listParam := networkmanagers.NetworkManagerDeploymentStatusParameter{
				Regions:         &[]string{normalizedLocation},
				DeploymentTypes: &[]networkmanagers.ConfigurationType{networkmanagers.ConfigurationType(state.ScopeAccess)},
			}
			resp, err := client.NetworkManagerDeploymentStatusList(ctx, *networkManagerId, listParam)

			if err != nil && !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("unexpected null model of %s", *id)
			}

			if !response.WasNotFound(resp.HttpResponse) && resp.Model.Value != nil && len(*resp.Model.Value) != 0 && *(*resp.Model.Value)[0].ConfigurationIds != nil && len(*(*resp.Model.Value)[0].ConfigurationIds) != 0 {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			input := networkmanagers.NetworkManagerCommit{
				ConfigurationIds: &state.ConfigurationIds,
				TargetLocations:  []string{state.Location},
				CommitType:       networkmanagers.ConfigurationType(state.ScopeAccess),
			}

			if _, err := client.NetworkManagerCommitsPost(ctx, *networkManagerId, input); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			if err = resourceManagerDeploymentWaitForFinished(ctx, client, id, time.Until(deadline)); err != nil {
				return err
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 24 * time.Hour,
	}
}

func (r ManagerDeploymentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagers
			id, err := parse.NetworkManagerDeploymentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("retrieving %s", *id)

			listParam := networkmanagers.NetworkManagerDeploymentStatusParameter{
				Regions:         &[]string{id.Location},
				DeploymentTypes: &[]networkmanagers.ConfigurationType{networkmanagers.ConfigurationType(id.ScopeAccess)},
			}

			networkManagerId := networkmanagers.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName)

			resp, err := client.NetworkManagerDeploymentStatusList(ctx, networkManagerId, listParam)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					metadata.Logger.Infof("%s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("unexpected null model of %s", *id)
			}

			if resp.Model.Value == nil || len(*resp.Model.Value) == 0 || (*resp.Model.Value)[0].ConfigurationIds == nil || len(*(*resp.Model.Value)[0].ConfigurationIds) == 0 {
				metadata.Logger.Infof("%s was not found - removing from state!", *id)
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

func (r ManagerDeploymentResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.NetworkManagerDeploymentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			metadata.Logger.Infof("updating %s..", *id)
			client := metadata.Client.Network.NetworkManagers

			listParam := networkmanagers.NetworkManagerDeploymentStatusParameter{
				Regions:         &[]string{id.Location},
				DeploymentTypes: &[]networkmanagers.ConfigurationType{networkmanagers.ConfigurationType(id.ScopeAccess)},
			}

			networkManagerId := networkmanagers.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName)

			resp, err := client.NetworkManagerDeploymentStatusList(ctx, networkManagerId, listParam)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					metadata.Logger.Infof("%s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("unexpected null model of %s", *id)
			}

			if resp.Model.Value == nil || len(*resp.Model.Value) == 0 || *(*resp.Model.Value)[0].ConfigurationIds == nil || len(*(*resp.Model.Value)[0].ConfigurationIds) == 0 {
				metadata.Logger.Infof("%s was not found - removing from state!", *id)
				return metadata.MarkAsGone(id)
			}

			if len(*resp.Model.Value) > 1 {
				return fmt.Errorf("found more than one deployment with id %s", *id)
			}

			deployment := (*resp.Model.Value)[0]
			if deployment.ConfigurationIds == nil {
				return fmt.Errorf("unexpected null configuration ID of %s", *id)
			}

			var state ManagerDeploymentModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("configuration_ids") {
				deployment.ConfigurationIds = &state.ConfigurationIds
			}

			input := networkmanagers.NetworkManagerCommit{
				ConfigurationIds: deployment.ConfigurationIds,
				TargetLocations:  []string{state.Location},
				CommitType:       networkmanagers.ConfigurationType(state.ScopeAccess),
			}

			if _, err := client.NetworkManagerCommitsPost(ctx, networkManagerId, input); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			if err = resourceManagerDeploymentWaitForFinished(ctx, client, id, time.Until(deadline)); err != nil {
				return err
			}

			return nil
		},
		Timeout: 24 * time.Hour,
	}
}

func (r ManagerDeploymentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagers
			id, err := parse.NetworkManagerDeploymentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			metadata.Logger.Infof("deleting %s..", *id)
			input := networkmanagers.NetworkManagerCommit{
				ConfigurationIds: &[]string{},
				TargetLocations:  []string{id.Location},
				CommitType:       networkmanagers.ConfigurationType(id.ScopeAccess),
			}

			networkManagerId := networkmanagers.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName)

			if _, err := client.NetworkManagerCommitsPost(ctx, networkManagerId, input); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			statusClient := metadata.Client.Network.NetworkManagers
			if err = resourceManagerDeploymentWaitForDeleted(ctx, statusClient, id, time.Until(deadline)); err != nil {
				return err
			}

			return nil
		},
		Timeout: 24 * time.Hour,
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

	_, err := state.WaitForStateContext(ctx)
	if err != nil {
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
		Refresh: func() (interface{}, string, error) {
			result, state, err := resourceManagerDeploymentResultRefreshFunc(ctx, client, managerDeploymentId)()
			if state == "NotFound" {
				// the deployment might not found after initial commit, https://github.com/Azure/azure-rest-api-specs/issues/27327
				// to serve NotFoundChecks, return nil result
				return nil, state, err
			}
			return result, state, err
		},
	}

	_, err := state.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for the %s: %+v", *managerDeploymentId, err)
	}

	return nil
}

func resourceManagerDeploymentResultRefreshFunc(ctx context.Context, client *networkmanagers.NetworkManagersClient, id *parse.ManagerDeploymentId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		listParam := networkmanagers.NetworkManagerDeploymentStatusParameter{
			Regions:         &[]string{azure.NormalizeLocation(id.Location)},
			DeploymentTypes: &[]networkmanagers.ConfigurationType{networkmanagers.ConfigurationType(id.ScopeAccess)},
		}

		networkManagerId := networkmanagers.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName)

		resp, err := client.NetworkManagerDeploymentStatusList(ctx, networkManagerId, listParam)
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
