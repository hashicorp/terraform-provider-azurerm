// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_discovery_workspace

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagediscovery/2025-09-01/storagediscoveryworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	storageClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := storageClient.From(metadata.Client).StorageDiscoveryWorkspacesClient

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

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := StorageDiscoveryWorkspaceModel{
				Name:              id.StorageDiscoveryWorkspaceName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(resp.Model.Location),
				Tags:              pointer.From(resp.Model.Tags),
			}

			if props := resp.Model.Properties; props != nil {
				state.Description = pointer.From(props.Description)
				state.WorkspaceRoots = props.WorkspaceRoots
				state.Sku = pointer.FromEnum(props.Sku)
				state.Scope = flattenStorageDiscoveryScopes(props.Scopes)
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}
			return metadata.Encode(&state)
		},
	}
}

func flattenStorageDiscoveryScopes(input []storagediscoveryworkspaces.StorageDiscoveryScope) []StorageDiscoveryScopeModel {
	result := make([]StorageDiscoveryScopeModel, 0)

	for _, scope := range input {
		model := StorageDiscoveryScopeModel{
			DisplayName:   scope.DisplayName,
			ResourceTypes: pointer.FromEnumSlice(&scope.ResourceTypes),
		}

		model.TagKeysOnly = pointer.From(scope.TagKeysOnly)
		model.Tags = pointer.From(scope.Tags)

		result = append(result, model)
	}

	return result
}

func storageDiscoveryScopeReplacementPath(oldRaw, newRaw any) string {
	oldScopes := oldRaw.([]any)
	newScopes := newRaw.([]any)

	oldScopesByName := make(map[string]map[string]any, len(oldScopes))
	for _, item := range oldScopes {
		if item == nil {
			continue
		}

		scope := item.(map[string]any)
		oldScopesByName[scope["display_name"].(string)] = scope
	}

	for index, item := range newScopes {
		if item == nil {
			continue
		}

		newScope := item.(map[string]any)
		oldScope, exists := oldScopesByName[newScope["display_name"].(string)]
		if !exists {
			continue
		}

		replacementField := ""
		switch {
		case !oldScope["resource_types"].(*pluginsdk.Set).Equal(newScope["resource_types"]):
			replacementField = "resource_types"
		case !oldScope["tag_keys_only"].(*pluginsdk.Set).Equal(newScope["tag_keys_only"]):
			replacementField = "tag_keys_only"
		case !reflect.DeepEqual(oldScope["tags"], newScope["tags"]):
			replacementField = "tags"
		}

		if replacementField != "" {
			// ForceNew requires a positional change. When a scope moves, its filters
			// can match those previously at its new index, but its display name changes.
			if index >= len(oldScopes) || oldScopes[index] == nil || oldScopes[index].(map[string]any)["display_name"] != newScope["display_name"] {
				replacementField = "display_name"
			}
			return fmt.Sprintf("scope.%d.%s", index, replacementField)
		}
	}

	return ""
}

func validateStorageDiscoveryScopes(raw any) error {
	scopes := raw.([]any)
	displayNames := make(map[string]struct{}, len(scopes))

	for _, item := range scopes {
		if item == nil {
			continue
		}

		scope := item.(map[string]any)
		displayName := scope["display_name"].(string)
		if displayName == "" {
			continue
		}

		if _, exists := displayNames[displayName]; exists {
			return fmt.Errorf("the `display_name` %q must be unique across all `scope` blocks", displayName)
		}

		displayNames[displayName] = struct{}{}
	}

	return nil
}
