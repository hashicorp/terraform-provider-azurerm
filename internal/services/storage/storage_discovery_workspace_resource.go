// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagediscovery/2025-09-01/storagediscoveryworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = StorageDiscoveryWorkspaceResource{}

type StorageDiscoveryWorkspaceResource struct{}

type StorageDiscoveryWorkspaceModel struct {
	Name              string                       `tfschema:"name"`
	ResourceGroupName string                       `tfschema:"resource_group_name"`
	Location          string                       `tfschema:"location"`
	WorkspaceRoot     []string                     `tfschema:"workspace_root"`
	Scopes            []StorageDiscoveryScopeModel `tfschema:"scopes"`
	Description       string                       `tfschema:"description"`
	Sku               string                       `tfschema:"sku"`
	Tags              map[string]string            `tfschema:"tags"`
}

type StorageDiscoveryScopeModel struct {
	DisplayName   string            `tfschema:"display_name"`
	ResourceTypes []string          `tfschema:"resource_types"`
	TagKeysOnly   []string          `tfschema:"tag_keys_only"`
	Tags          map[string]string `tfschema:"tags"`
}

func (r StorageDiscoveryWorkspaceResource) ResourceType() string {
	return "azurerm_storage_discovery_workspace"
}

func (r StorageDiscoveryWorkspaceResource) ModelObject() interface{} {
	return &StorageDiscoveryWorkspaceModel{}
}

func (r StorageDiscoveryWorkspaceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return storagediscoveryworkspaces.ValidateProviderStorageDiscoveryWorkspaceID
}

func (r StorageDiscoveryWorkspaceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.StorageDiscoveryWorkspaceName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"workspace_root": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			MaxItems: 100,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.StorageDiscoveryWorkspaceRoot,
			},
		},

		"scopes": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"display_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"resource_types": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"tag_keys_only": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"tags": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(storagediscoveryworkspaces.StorageDiscoverySkuStandard),
			ValidateFunc: validation.StringInSlice(
				storagediscoveryworkspaces.PossibleValuesForStorageDiscoverySku(),
				false,
			),
		},

		"tags": commonschema.Tags(),
	}
}

func (r StorageDiscoveryWorkspaceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StorageDiscoveryWorkspaceResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diff := metadata.ResourceDiff

			workspaceRootsRaw := diff.Get("workspace_root")
			if workspaceRootsRaw == nil {
				return nil
			}

			workspaceRootsSet := workspaceRootsRaw.(*pluginsdk.Set)
			workspaceRoots := make([]string, 0)
			for _, item := range workspaceRootsSet.List() {
				workspaceRoots = append(workspaceRoots, item.(string))
			}

			subscriptionIDs := make(map[string]bool)
			resourceGroupIDs := make([]commonids.ResourceGroupId, 0)

			// First pass: collect all subscription IDs and resource group IDs
			for _, rootID := range workspaceRoots {
				// Try to parse as subscription ID
				if subscriptionID, err := commonids.ParseSubscriptionID(rootID); err == nil {
					subscriptionIDs[subscriptionID.SubscriptionId] = true
					continue
				}

				// Try to parse as resource group ID
				if resourceGroupID, err := commonids.ParseResourceGroupID(rootID); err == nil {
					resourceGroupIDs = append(resourceGroupIDs, *resourceGroupID)
				}
			}

			// Second pass: check if any resource group belongs to a subscription in the list
			for _, rgID := range resourceGroupIDs {
				if subscriptionIDs[rgID.SubscriptionId] {
					return fmt.Errorf("cannot specify both subscription ID `/subscriptions/%s` and its child resource group ID `%s` in `workspace_root`", rgID.SubscriptionId, rgID.ID())
				}
			}

			return nil
		},
	}
}

func (r StorageDiscoveryWorkspaceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.StorageDiscoveryWorkspacesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model StorageDiscoveryWorkspaceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := storagediscoveryworkspaces.NewProviderStorageDiscoveryWorkspaceID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			sku := storagediscoveryworkspaces.StorageDiscoverySkuStandard
			if model.Sku != "" {
				sku = storagediscoveryworkspaces.StorageDiscoverySku(model.Sku)
			}

			payload := storagediscoveryworkspaces.StorageDiscoveryWorkspace{
				Location: location.Normalize(model.Location),
				Properties: &storagediscoveryworkspaces.StorageDiscoveryWorkspaceProperties{
					WorkspaceRoots: model.WorkspaceRoot,
					Scopes:         expandStorageDiscoveryScopes(model.Scopes),
					Sku:            &sku,
				},
				Tags: pointer.To(model.Tags),
			}

			if model.Description != "" {
				payload.Properties.Description = pointer.To(model.Description)
			}

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r StorageDiscoveryWorkspaceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.StorageDiscoveryWorkspacesClient

			id, err := storagediscoveryworkspaces.ParseProviderStorageDiscoveryWorkspaceID(metadata.ResourceData.Id())
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

			state := StorageDiscoveryWorkspaceModel{
				Name:              id.StorageDiscoveryWorkspaceName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if resp.Model != nil {
				state.Location = location.Normalize(resp.Model.Location)
				state.Tags = pointer.From(resp.Model.Tags)

				if props := resp.Model.Properties; props != nil {
					state.Description = pointer.From(props.Description)
					state.WorkspaceRoot = props.WorkspaceRoots

					sku := string(storagediscoveryworkspaces.StorageDiscoverySkuStandard)
					if props.Sku != nil {
						sku = string(*props.Sku)
					}
					state.Sku = sku

					state.Scopes = flattenStorageDiscoveryScopes(props.Scopes)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r StorageDiscoveryWorkspaceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.StorageDiscoveryWorkspacesClient

			id, err := storagediscoveryworkspaces.ParseProviderStorageDiscoveryWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model StorageDiscoveryWorkspaceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload := storagediscoveryworkspaces.StorageDiscoveryWorkspaceUpdate{
				Properties: &storagediscoveryworkspaces.StorageDiscoveryWorkspacePropertiesUpdate{},
			}

			if metadata.ResourceData.HasChange("description") {
				payload.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("sku") {
				sku := storagediscoveryworkspaces.StorageDiscoverySku(model.Sku)
				payload.Properties.Sku = &sku
			}

			if metadata.ResourceData.HasChange("workspace_root") {
				payload.Properties.WorkspaceRoots = &model.WorkspaceRoot
			}

			if metadata.ResourceData.HasChange("scopes") {
				scopes := expandStorageDiscoveryScopes(model.Scopes)
				payload.Properties.Scopes = &scopes
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if _, err := client.Update(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r StorageDiscoveryWorkspaceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.StorageDiscoveryWorkspacesClient

			id, err := storagediscoveryworkspaces.ParseProviderStorageDiscoveryWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandStorageDiscoveryScopes(input []StorageDiscoveryScopeModel) []storagediscoveryworkspaces.StorageDiscoveryScope {
	result := make([]storagediscoveryworkspaces.StorageDiscoveryScope, 0)

	for _, scope := range input {
		apiScope := storagediscoveryworkspaces.StorageDiscoveryScope{
			DisplayName:   scope.DisplayName,
			ResourceTypes: expandStorageDiscoveryResourceTypes(scope.ResourceTypes),
		}

		if len(scope.TagKeysOnly) > 0 {
			apiScope.TagKeysOnly = &scope.TagKeysOnly
		}

		if len(scope.Tags) > 0 {
			apiScope.Tags = &scope.Tags
		}

		result = append(result, apiScope)
	}

	return result
}

func expandStorageDiscoveryResourceTypes(input []string) []storagediscoveryworkspaces.StorageDiscoveryResourceType {
	result := make([]storagediscoveryworkspaces.StorageDiscoveryResourceType, 0)
	for _, item := range input {
		if item != "" {
			result = append(result, storagediscoveryworkspaces.StorageDiscoveryResourceType(item))
		}
	}
	return result
}

func flattenStorageDiscoveryScopes(input []storagediscoveryworkspaces.StorageDiscoveryScope) []StorageDiscoveryScopeModel {
	result := make([]StorageDiscoveryScopeModel, 0)

	for _, scope := range input {
		model := StorageDiscoveryScopeModel{
			DisplayName:   scope.DisplayName,
			ResourceTypes: flattenStorageDiscoveryResourceTypes(scope.ResourceTypes),
		}

		if scope.TagKeysOnly != nil {
			model.TagKeysOnly = *scope.TagKeysOnly
		}

		if scope.Tags != nil {
			model.Tags = *scope.Tags
		}

		result = append(result, model)
	}

	return result
}

func flattenStorageDiscoveryResourceTypes(input []storagediscoveryworkspaces.StorageDiscoveryResourceType) []string {
	result := make([]string, 0)
	for _, item := range input {
		result = append(result, string(item))
	}
	return result
}
