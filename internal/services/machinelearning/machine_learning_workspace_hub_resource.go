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
	components "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = MachineLearningWorkspaceHubResource{}

type machineLearningWorkspaceHubModel struct {
	Name                        string                                     `tfschema:"name"`
	Location                    string                                     `tfschema:"location"`
	ResourceGroupName           string                                     `tfschema:"resource_group_name"`
	KeyVaultID                  string                                     `tfschema:"key_vault_id"`
	StorageAccountID            string                                     `tfschema:"storage_account_id"`
	ContainerRegistryID         string                                     `tfschema:"container_registry_id"`
	ApplicationInsightsID       string                                     `tfschema:"application_insights_id"`
	PublicNetworkAccess         string                                     `tfschema:"public_network_access"`
	WorkspaceId                 string                                     `tfschema:"workspace_id"`
	FriendlyName                string                                     `tfschema:"friendly_name"`
	PrimaryUserAssignedIdentity string                                     `tfschema:"primary_user_assigned_identity"`
	Identity                    []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Encryption                  []encryptionModel                          `tfschema:"encryption"`
	Tags                        map[string]string                          `tfschema:"tags"`
}

type encryptionModel struct {
	KeyVaultID             string `tfschema:"key_vault_id"`
	KeyID                  string `tfschema:"key_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

type MachineLearningWorkspaceHubResource struct{}

var _ sdk.Resource = MachineLearningWorkspaceHubResource{}

func (r MachineLearningWorkspaceHubResource) ResourceType() string {
	return "azurerm_machine_learning_workspace_hub"
}

func (r MachineLearningWorkspaceHubResource) ModelObject() interface{} {
	return &machineLearningWorkspaceHubModel{}
}

func (r MachineLearningWorkspaceHubResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return workspaces.ValidateWorkspaceID
}

func (r MachineLearningWorkspaceHubResource) Arguments() map[string]*pluginsdk.Schema {
	arguments := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"key_vault_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
			// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
			// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"application_insights_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: components.ValidateComponentID,
			// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityRequired(),

		"container_registry_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: registries.ValidateRegistryID,
			// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"encryption": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),
					"key_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
					"user_assigned_identity_id": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						ValidateFunc:     commonids.ValidateUserAssignedIdentityID,
						DiffSuppressFunc: suppress.CaseDifference,
					},
				},
			},
		},

		"friendly_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"primary_user_assigned_identity": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"public_network_access": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(workspaces.PublicNetworkAccessEnabled),
				string(workspaces.PublicNetworkAccessDisabled),
			}, false),
			Default: string(workspaces.PublicNetworkAccessEnabled),
		},

		"tags": commonschema.Tags(),
	}
	return arguments
}

func (r MachineLearningWorkspaceHubResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"workspace_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r MachineLearningWorkspaceHubResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model machineLearningWorkspaceHubModel
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
				return tf.ImportAsExistsError("azurerm_machine_learning_workspace_hub", id.ID())
			}

			hubIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			hub := workspaces.Workspace{
				Name:     &model.Name,
				Location: pointer.To(location.Normalize(model.Location)),
				Tags:     &model.Tags,
				Kind:     pointer.To("Hub"),
				Identity: hubIdentity,
				Properties: &workspaces.WorkspaceProperties{
					FriendlyName:                &model.FriendlyName,
					Encryption:                  expandMachineLearningWorkspaceHubEncryption(model.Encryption),
					KeyVault:                    &model.KeyVaultID,
					PublicNetworkAccess:         pointer.To(workspaces.PublicNetworkAccess(model.PublicNetworkAccess)),
					StorageAccount:              &model.StorageAccountID,
					PrimaryUserAssignedIdentity: &model.PrimaryUserAssignedIdentity,
				},
			}

			if model.ApplicationInsightsID != "" {
				hub.Properties.ApplicationInsights = &model.ApplicationInsightsID
			}

			if model.ContainerRegistryID != "" {
				hub.Properties.ContainerRegistry = &model.ContainerRegistryID
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, hub); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r MachineLearningWorkspaceHubResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model machineLearningWorkspaceHubModel
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

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}

				payload.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("primary_user_assigned_identity") {
				payload.Properties.PrimaryUserAssignedIdentity = &model.PrimaryUserAssignedIdentity
			}

			if metadata.ResourceData.HasChange("friendly_name") {
				payload.Properties.FriendlyName = &model.FriendlyName
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = &model.Tags
			}

			if metadata.ResourceData.HasChange("public_network_access") {
				payload.Properties.PublicNetworkAccess = pointer.To(workspaces.PublicNetworkAccess(model.PublicNetworkAccess))
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r MachineLearningWorkspaceHubResource) Read() sdk.ResourceFunc {
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

			state := machineLearningWorkspaceHubModel{
				Name:              id.WorkspaceName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          *model.Location,
				Tags:              *model.Tags,
			}

			if props := model.Properties; props != nil {
				appInsightsId := ""
				if props.ApplicationInsights != nil {
					applicationInsightsId, err := components.ParseComponentIDInsensitively(*props.ApplicationInsights)
					if err != nil {
						return err
					}
					appInsightsId = applicationInsightsId.ID()
				}
				state.ApplicationInsightsID = appInsightsId

				if model.Properties.KeyVault != nil {
					state.KeyVaultID = *model.Properties.KeyVault
				}

				if model.Properties.StorageAccount != nil {
					state.StorageAccountID = *model.Properties.StorageAccount
				}

				if model.Properties.PublicNetworkAccess != nil {
					state.PublicNetworkAccess = string(*model.Properties.PublicNetworkAccess)
				}

				if model.Properties.WorkspaceId != nil {
					state.WorkspaceId = *model.Properties.WorkspaceId
				}

				if model.Properties.PrimaryUserAssignedIdentity != nil {
					state.PrimaryUserAssignedIdentity = *model.Properties.PrimaryUserAssignedIdentity
				}

				if model.Properties.FriendlyName != nil {
					state.FriendlyName = *model.Properties.FriendlyName
				}

				if model.Properties.ContainerRegistry != nil {
					state.ContainerRegistryID = *model.Properties.ContainerRegistry
				}

				if model.Properties.WorkspaceId != nil {
					state.WorkspaceId = *model.Properties.WorkspaceId
				}

				flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = flattenedIdentity

				flattenedEncryption, err := flattenMachineLearningWorkspaceHubEncryption(props.Encryption)
				if err != nil {
					return fmt.Errorf("flattening `encryption`: %+v", err)
				}
				if flattenedEncryption != nil {
					state.Encryption = flattenedEncryption
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r MachineLearningWorkspaceHubResource) Delete() sdk.ResourceFunc {
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
				return fmt.Errorf("deleting Machine Learning Workspace Hub %q (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroupName, err)
			}

			if err := future.Poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for deletion of Machine Learning Workspace Hub %q (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroupName, err)
			}

			return nil
		},
	}
}

func expandMachineLearningWorkspaceHubEncryption(encryption []encryptionModel) *workspaces.EncryptionProperty {
	if len(encryption) == 0 {
		return &workspaces.EncryptionProperty{}
	}

	encrypt := encryption[0]
	out := workspaces.EncryptionProperty{
		Identity: &workspaces.IdentityForCmk{
			UserAssignedIdentity: nil,
		},
		KeyVaultProperties: workspaces.EncryptionKeyVaultProperties{
			KeyVaultArmId: encrypt.KeyVaultID,
			KeyIdentifier: encrypt.KeyID,
		},
		Status: workspaces.EncryptionStatusEnabled,
	}

	if encrypt.UserAssignedIdentityID != "" {
		out.Identity.UserAssignedIdentity = &encrypt.UserAssignedIdentityID
	}

	return &out
}

func flattenMachineLearningWorkspaceHubEncryption(input *workspaces.EncryptionProperty) ([]encryptionModel, error) {
	if input == nil || input.Status != workspaces.EncryptionStatusEnabled {
		return []encryptionModel{}, nil
	}

	keyVaultId := ""
	keyVaultKeyId := ""

	if input.KeyVaultProperties.KeyIdentifier != "" {
		keyVaultKeyId = input.KeyVaultProperties.KeyIdentifier
	}
	if input.KeyVaultProperties.KeyVaultArmId != "" {
		keyVaultId = input.KeyVaultProperties.KeyVaultArmId
	}

	userAssignedIdentityId := ""
	if input.Identity != nil && input.Identity.UserAssignedIdentity != nil {
		id, err := commonids.ParseUserAssignedIdentityIDInsensitively(*input.Identity.UserAssignedIdentity)
		if err != nil {
			return nil, fmt.Errorf("parsing userAssignedIdentity %q: %+v", *input.Identity.UserAssignedIdentity, err)
		}
		userAssignedIdentityId = id.ID()
	}

	return []encryptionModel{
		{
			UserAssignedIdentityID: userAssignedIdentityId,
			KeyVaultID:             keyVaultId,
			KeyID:                  keyVaultKeyId,
		},
	}, nil
}
