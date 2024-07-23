// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = MachineLearningWorkspaceProjectResource{}

type machineLearningWorkspaceProjectModel struct {
	Name                string                                     `tfschema:"name"`
	Location            string                                     `tfschema:"location"`
	ResourceGroupName   string                                     `tfschema:"resource_group_name"`
	WorkspaceHubID      string                                     `tfschema:"workspace_hub_id"`
	PublicNetworkAccess string                                     `tfschema:"public_network_access"`
	WorkspaceId         string                                     `tfschema:"workspace_id"`
	FriendlyName        string                                     `tfschema:"friendly_name"`
	Identity            []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Tags                map[string]string                          `tfschema:"tags"`
}

type MachineLearningWorkspaceProjectResource struct{}

var _ sdk.Resource = MachineLearningWorkspaceProjectResource{}

func (r MachineLearningWorkspaceProjectResource) ResourceType() string {
	return "azurerm_machine_learning_workspace_project"
}

func (r MachineLearningWorkspaceProjectResource) ModelObject() interface{} {
	return &machineLearningWorkspaceProjectModel{}
}

func (r MachineLearningWorkspaceProjectResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return workspaces.ValidateWorkspaceID
}

func (r MachineLearningWorkspaceProjectResource) Arguments() map[string]*pluginsdk.Schema {
	arguments := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"workspace_hub_id": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateFunc:     workspaces.ValidateWorkspaceID,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"identity": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(identity.TypeSystemAssigned),
							string(identity.TypeSystemAssignedUserAssigned),
						}, false),
					},
					"identity_ids": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: commonids.ValidateUserAssignedIdentityID,
						},
					},
					"principal_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"tenant_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},

		"friendly_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"tags": commonschema.Tags(),
	}
	return arguments
}

func (r MachineLearningWorkspaceProjectResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"public_network_access": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"workspace_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r MachineLearningWorkspaceProjectResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model machineLearningWorkspaceProjectModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MachineLearning.Workspaces
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := workspaces.NewWorkspaceID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_machine_learning_workspace_project", id.ID())
			}

			projectIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			project := workspaces.Workspace{
				Name:     &model.Name,
				Location: pointer.To(location.Normalize(model.Location)),
				Tags:     &model.Tags,
				Kind:     pointer.To("Project"),
				Identity: projectIdentity,
				Properties: &workspaces.WorkspaceProperties{
					HubResourceId: &model.WorkspaceHubID,
					FriendlyName:  &model.FriendlyName,
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, project); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r MachineLearningWorkspaceProjectResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model machineLearningWorkspaceProjectModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MachineLearning.Workspaces
			id, err := workspaces.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			payload := existing.Model
			payload.Properties.StorageAccount = nil
			payload.Properties.KeyVault = nil
			payload.Properties.ApplicationInsights = nil
			payload.Properties.ContainerRegistry = nil
			payload.Properties.Encryption = nil
			payload.Properties.ManagedNetwork = nil

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}

				payload.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("friendly_name") {
				payload.Properties.FriendlyName = &model.FriendlyName
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r MachineLearningWorkspaceProjectResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Workspaces

			id, err := workspaces.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := machineLearningWorkspaceProjectModel{
				Name:              id.WorkspaceName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          *model.Location,
				Tags:              *model.Tags,
			}

			if props := model.Properties; props != nil {
				if model.Properties.PublicNetworkAccess != nil {
					state.PublicNetworkAccess = string(*model.Properties.PublicNetworkAccess)
				}

				if model.Properties.FriendlyName != nil {
					state.FriendlyName = *model.Properties.FriendlyName
				}

				if model.Properties.WorkspaceId != nil {
					state.WorkspaceId = *model.Properties.WorkspaceId
				}

				if model.Properties.HubResourceId != nil {
					state.WorkspaceHubID = *model.Properties.HubResourceId
				}

				flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = flattenedIdentity
			}

			return metadata.Encode(&state)
		},
	}
}

func (r MachineLearningWorkspaceProjectResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Workspaces

			id, err := workspaces.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			options := workspaces.DefaultDeleteOperationOptions()
			if metadata.Client.Features.MachineLearning.PurgeSoftDeletedWorkspaceOnDestroy {
				options = workspaces.DeleteOperationOptions{
					ForceToPurge: pointer.To(true),
				}
			}

			future, err := client.Delete(ctx, *id, options)
			if err != nil {
				return fmt.Errorf("deleting Machine Learning Workspace Project %q (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroupName, err)
			}

			if err := future.Poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for deletion of Machine Learning Workspace Project %q (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroupName, err)
			}

			return nil
		},
	}
}
