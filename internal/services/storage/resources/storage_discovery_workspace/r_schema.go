// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_discovery_workspace

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagediscovery/2025-09-01/storagediscoveryworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.StorageDiscoveryWorkspaceName(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"scope": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 10,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"display_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.StorageDiscoveryScopeDisplayName(),
					},

					"resource_types": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(
								storagediscoveryworkspaces.PossibleValuesForStorageDiscoveryResourceType(),
								false,
							),
						},
					},

					"tag_keys_only": {
						Type:     pluginsdk.TypeSet,
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

		"workspace_roots": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			MaxItems: 100,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.Any(commonids.ValidateSubscriptionID, commonids.ValidateResourceGroupID),
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(0, 1024),
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

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diff := metadata.ResourceDiff

			if err := validateStorageDiscoveryScopes(diff.Get("scope")); err != nil {
				return err
			}

			if diff.HasChange("scope") {
				oldScopes, newScopes := diff.GetChange("scope")
				if replacementPath := storageDiscoveryScopeReplacementPath(oldScopes, newScopes); replacementPath != "" {
					if err := diff.ForceNew(replacementPath); err != nil {
						return err
					}
				}
			}

			workspaceRootsRaw := diff.Get("workspace_roots")
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

			for _, rootID := range workspaceRoots {
				if subscriptionID, err := commonids.ParseSubscriptionID(rootID); err == nil {
					subscriptionIDs[strings.ToLower(subscriptionID.SubscriptionId)] = true
					continue
				}

				if resourceGroupID, err := commonids.ParseResourceGroupID(rootID); err == nil {
					resourceGroupIDs = append(resourceGroupIDs, *resourceGroupID)
				}
			}

			for _, rgID := range resourceGroupIDs {
				if subscriptionIDs[strings.ToLower(rgID.SubscriptionId)] {
					return fmt.Errorf("cannot specify both subscription ID `/subscriptions/%s` and its child resource group ID `%s` in `workspace_roots`", rgID.SubscriptionId, rgID.ID())
				}
			}

			return nil
		},
	}
}
