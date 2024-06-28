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
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2023-05-01/deployments"
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
	Name                  string                                     `tfschema:"name"`
	Location              string                                     `tfschema:"location"`
	ResourceGroupName     string                                     `tfschema:"resource_group_name"`
	KeyVaultID            string                                     `tfschema:"key_vault_id"`
	StorageAccountID      string                                     `tfschema:"storage_account_id"`
	ContainerRegistryID   string                                     `tfschema:"container_registry_id"`
	ApplicationInsightsID string                                     `tfschema:"application_insights_id"`
	PublicNetworkAccess   string                                     `tfschema:"public_network_access"`
	Identity              []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Encryption            encryptionModel                            `tfschema:"encryption"`
	Tags                  map[string]string                          `tfschema:"tags"`
}

type encryptionModel struct {
	KeyVaultID             string `tfschema:"key_vault_key_id"`
	KeyID                  string `tfschema:"key_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

type MachineLearningWorkspaceHubResource struct{}

func (r MachineLearningWorkspaceHubResource) CustomImporter() sdk.ResourceRunFunc {
	panic("implement me")
}

var _ sdk.Resource = MachineLearningWorkspaceHubResource{}

func (r MachineLearningWorkspaceHubResource) ResourceType() string {
	return "azurerm_machine_learning_workspace_hub"
}

func (r MachineLearningWorkspaceHubResource) ModelObject() interface{} {
	return &machineLearningWorkspaceHubModel{}
}

func (r MachineLearningWorkspaceHubResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return deployments.ValidateDeploymentID
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

		"container_registry_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: registries.ValidateRegistryID,
			// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"application_insights_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: components.ValidateComponentID,
			// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityRequired(),

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
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
						// TODO: remove this
						DiffSuppressFunc: suppress.CaseDifference,
					},
				},
			},
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
				//Sku: &workspaces.Sku{
				//	Name: d.Get("sku_name").(string),
				//	Tier: pointer.To(workspaces.SkuTier(d.Get("sku_name").(string))),
				//},
				Kind:     pointer.To("Hub"),
				Identity: hubIdentity,
				Properties: &workspaces.WorkspaceProperties{
					ApplicationInsights: &model.ApplicationInsightsID,
					Encryption:          expandMachineLearningWorkspaceHubEncryption(model.Encryption),
					KeyVault:            &model.KeyVaultID,
					//ManagedNetwork:      expandMachineLearningWorkspaceManagedNetwork(d.Get("managed_network").([]interface{})),
					PublicNetworkAccess: pointer.To(workspaces.PublicNetworkAccess(model.PublicNetworkAccess)),
					StorageAccount:      &model.StorageAccountID,
				},
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

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = &model.Tags
			}

			if metadata.ResourceData.HasChange("public_network_access") {
				payload.Properties.PublicNetworkAccess = pointer.To(workspaces.PublicNetworkAccess(model.PublicNetworkAccess))
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
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

				state.KeyVaultID = *model.Properties.KeyVault
				state.StorageAccountID = *model.Properties.StorageAccount
				state.ContainerRegistryID = *model.Properties.ContainerRegistry
				state.PublicNetworkAccess = string(*model.Properties.PublicNetworkAccess)

				flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = flattenedIdentity

				flattenedEncryption, err := flattenMachineLearningWorkspaceHubEncryption(props.Encryption)
				if err != nil {
					return fmt.Errorf("flattening `encryption`: %+v", err)
				}
				state.Encryption = *flattenedEncryption
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
				return fmt.Errorf("deleting Machine Learning Hub Workspace %q (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroupName, err)
			}

			if err := future.Poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for deletion of Machine Learning Hub Workspace %q (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroupName, err)
			}

			return nil
		},
	}
}

func expandMachineLearningWorkspaceHubEncryption(encryption encryptionModel) *workspaces.EncryptionProperty {
	if encryption.KeyID == "" || encryption.KeyVaultID == "" {
		return nil
	}

	out := workspaces.EncryptionProperty{
		Identity: &workspaces.IdentityForCmk{
			UserAssignedIdentity: nil,
		},
		KeyVaultProperties: workspaces.EncryptionKeyVaultProperties{
			KeyVaultArmId: encryption.KeyVaultID,
			KeyIdentifier: encryption.KeyID,
		},
		Status: workspaces.EncryptionStatusEnabled,
	}

	if encryption.UserAssignedIdentityID != "" {
		out.Identity.UserAssignedIdentity = &encryption.UserAssignedIdentityID
	}

	return &out
}

func flattenMachineLearningWorkspaceHubEncryption(input *workspaces.EncryptionProperty) (*encryptionModel, error) {
	if input == nil || input.Status != workspaces.EncryptionStatusEnabled {
		return nil, nil
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
			return nil, fmt.Errorf("parsing userAssignedIdentityId %q: %+v", *input.Identity.UserAssignedIdentity, err)
		}
		userAssignedIdentityId = id.ID()
	}

	return &encryptionModel{
		UserAssignedIdentityID: userAssignedIdentityId,
		KeyVaultID:             keyVaultId,
		KeyID:                  keyVaultKeyId,
	}, nil
}
