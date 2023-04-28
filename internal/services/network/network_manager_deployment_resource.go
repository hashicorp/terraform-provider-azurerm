package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

var _ sdk.ResourceWithUpdate = ManagerDeploymentResource{}

type ManagerDeploymentModel struct {
	NetworkManagerId string   `tfschema:"network_manager_id"`
	ScopeAccess      string   `tfschema:"scope_access"`
	Location         string   `tfschema:"location"`
	ConfigurationIds []string `tfschema:"configuration_ids"`
	DeploymentStatus string   `tfschema:"deployment_status"`
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
			ValidateFunc: validate.NetworkManagerID,
		},

		"location": commonschema.Location(),

		"scope_access": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(network.ConfigurationTypeConnectivity),
				string(network.ConfigurationTypeSecurityAdmin),
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
	return map[string]*pluginsdk.Schema{
		"deployment_status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ManagerDeploymentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			metadata.Logger.Info("Decoding state..")
			var state ManagerDeploymentModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Network.ManagerDeploymentsClient
			statusClient := metadata.Client.Network.ManagerDeploymentStatusClient
			networkManagerId, err := parse.NetworkManagerID(state.NetworkManagerId)
			if err != nil {
				return err
			}

			normalizedLocation := azure.NormalizeLocation(state.Location)
			id := parse.NewNetworkManagerDeploymentID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroup, networkManagerId.Name, normalizedLocation, state.ScopeAccess)

			metadata.Logger.Infof("creating %s", id)

			listParam := network.ManagerDeploymentStatusParameter{
				Regions:         &[]string{normalizedLocation},
				DeploymentTypes: &[]network.ConfigurationType{network.ConfigurationType(state.ScopeAccess)},
			}
			resp, err := statusClient.List(ctx, listParam, id.ResourceGroup, id.NetworkManagerName, nil)

			if err != nil && !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(resp.Response) && resp.Value != nil && len(*resp.Value) != 0 && *(*resp.Value)[0].ConfigurationIds != nil && len(*(*resp.Value)[0].ConfigurationIds) != 0 {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			input := network.ManagerCommit{
				ConfigurationIds: &state.ConfigurationIds,
				TargetLocations:  &[]string{state.Location},
				CommitType:       network.ConfigurationType(state.ScopeAccess),
			}

			if _, err := client.Post(ctx, input, id.ResourceGroup, id.NetworkManagerName); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err = resourceManagerDeploymentWaitForFinished(ctx, statusClient, id, metadata.ResourceData); err != nil {
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
			client := metadata.Client.Network.ManagerDeploymentStatusClient
			id, err := parse.NetworkManagerDeploymentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("retrieving %s", *id)

			listParam := network.ManagerDeploymentStatusParameter{
				Regions:         &[]string{id.Location},
				DeploymentTypes: &[]network.ConfigurationType{network.ConfigurationType(id.ScopeAccess)},
			}
			resp, err := client.List(ctx, listParam, id.ResourceGroup, id.NetworkManagerName, nil)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					metadata.Logger.Infof("%s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Value == nil || len(*resp.Value) == 0 || (*resp.Value)[0].ConfigurationIds == nil || len(*(*resp.Value)[0].ConfigurationIds) == 0 {
				metadata.Logger.Infof("%s was not found - removing from state!", *id)
				return metadata.MarkAsGone(id)
			}

			deployment := (*resp.Value)[0]
			return metadata.Encode(&ManagerDeploymentModel{
				NetworkManagerId: parse.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName).ID(),
				Location:         location.NormalizeNilable(deployment.Region),
				ScopeAccess:      string(deployment.DeploymentType),
				ConfigurationIds: *deployment.ConfigurationIds,
				DeploymentStatus: string(deployment.DeploymentStatus),
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

			metadata.Logger.Infof("updating %s..", *id)
			client := metadata.Client.Network.ManagerDeploymentsClient
			statusClient := metadata.Client.Network.ManagerDeploymentStatusClient

			listParam := network.ManagerDeploymentStatusParameter{
				Regions:         &[]string{id.Location},
				DeploymentTypes: &[]network.ConfigurationType{network.ConfigurationType(id.ScopeAccess)},
			}
			resp, err := statusClient.List(ctx, listParam, id.ResourceGroup, id.NetworkManagerName, nil)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					metadata.Logger.Infof("%s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Value == nil || len(*resp.Value) == 0 || *(*resp.Value)[0].ConfigurationIds == nil || len(*(*resp.Value)[0].ConfigurationIds) == 0 {
				metadata.Logger.Infof("%s was not found - removing from state!", *id)
				return metadata.MarkAsGone(id)
			}

			deployment := (*resp.Value)[0]
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

			input := network.ManagerCommit{
				ConfigurationIds: deployment.ConfigurationIds,
				TargetLocations:  &[]string{state.Location},
				CommitType:       network.ConfigurationType(state.ScopeAccess),
			}

			if _, err := client.Post(ctx, input, id.ResourceGroup, id.NetworkManagerName); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err = resourceManagerDeploymentWaitForFinished(ctx, statusClient, id, metadata.ResourceData); err != nil {
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
			client := metadata.Client.Network.ManagerDeploymentsClient
			id, err := parse.NetworkManagerDeploymentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s..", *id)
			input := network.ManagerCommit{
				ConfigurationIds: &[]string{},
				TargetLocations:  &[]string{id.Location},
				CommitType:       network.ConfigurationType(id.ScopeAccess),
			}

			future, err := client.Post(ctx, input, id.ResourceGroup, id.NetworkManagerName)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
			}

			statusClient := metadata.Client.Network.ManagerDeploymentStatusClient
			if _, err = resourceManagerDeploymentWaitForDeleted(ctx, statusClient, id, metadata.ResourceData); err != nil {
				return err
			}

			return nil
		},
		Timeout: 24 * time.Hour,
	}
}

func resourceManagerDeploymentWaitForDeleted(ctx context.Context, client *network.ManagerDeploymentStatusClient, ManagerDeploymentId *parse.ManagerDeploymentId, d *pluginsdk.ResourceData) (network.ManagerDeploymentStatusListResult, error) {
	state := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"NotStarted", "Deploying", "Deployed", "Failed"},
		Target:     []string{"NotFound"},
		Refresh:    resourceManagerDeploymentResultRefreshFunc(ctx, client, ManagerDeploymentId),
		Timeout:    d.Timeout(pluginsdk.TimeoutCreate),
	}

	resp, err := state.WaitForStateContext(ctx)
	if err != nil {
		return resp.(network.ManagerDeploymentStatusListResult), fmt.Errorf("waiting for the Deployment %s: %+v", *ManagerDeploymentId, err)
	}

	return resp.(network.ManagerDeploymentStatusListResult), nil
}

func resourceManagerDeploymentWaitForFinished(ctx context.Context, client *network.ManagerDeploymentStatusClient, ManagerDeploymentId *parse.ManagerDeploymentId, d *pluginsdk.ResourceData) error {
	state := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"NotStarted", "Deploying"},
		Target:     []string{"Deployed"},
		Refresh:    resourceManagerDeploymentResultRefreshFunc(ctx, client, ManagerDeploymentId),
		Timeout:    d.Timeout(pluginsdk.TimeoutCreate),
	}

	_, err := state.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for the Deployment %s: %+v", *ManagerDeploymentId, err)
	}

	return nil
}

func resourceManagerDeploymentResultRefreshFunc(ctx context.Context, client *network.ManagerDeploymentStatusClient, id *parse.ManagerDeploymentId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		listParam := network.ManagerDeploymentStatusParameter{
			Regions:         &[]string{azure.NormalizeLocation(id.Location)},
			DeploymentTypes: &[]network.ConfigurationType{network.ConfigurationType(id.ScopeAccess)},
		}
		resp, err := client.List(ctx, listParam, id.ResourceGroup, id.NetworkManagerName, nil)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return resp, "NotFound", nil
			}
			return resp, "Error", fmt.Errorf("retrieving Deployment: %+v", err)
		}

		if resp.Value == nil || len(*resp.Value) == 0 || *(*resp.Value)[0].ConfigurationIds == nil || len(*(*resp.Value)[0].ConfigurationIds) == 0 {
			return resp, "NotFound", nil
		}

		deploymentStatus := string((*resp.Value)[0].DeploymentStatus)
		return resp, deploymentStatus, nil
	}
}
