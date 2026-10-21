// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state ManagerModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.NetworkManager.NetworkManagers
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := networkmanagers.NewNetworkManagerID(subscriptionId, state.ResourceGroupName, state.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			input := networkmanagers.NetworkManager{
				Location: pointer.To(location.Normalize(state.Location)),
				Name:     pointer.To(state.Name),
				Properties: &networkmanagers.NetworkManagerProperties{
					Description:                 pointer.To(state.Description),
					NetworkManagerScopes:        expandNetworkManagerScope(state.Scope),
					NetworkManagerScopeAccesses: expandNetworkManagerScopeAccesses(state.ScopeAccesses),
				},
				Tags: pluginsdk.ExpandPtrMapStringString(state.Tags),
			}

			if _, err := client.CreateOrUpdate(ctx, id, input); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func expandNetworkManagerScope(input []ManagerScopeModel) networkmanagers.NetworkManagerPropertiesNetworkManagerScopes {
	return networkmanagers.NetworkManagerPropertiesNetworkManagerScopes{
		ManagementGroups: pointer.To(input[0].ManagementGroups),
		Subscriptions:    pointer.To(input[0].Subscriptions),
	}
}

func expandNetworkManagerScopeAccesses(input []string) *[]networkmanagers.ConfigurationType {
	result := make([]networkmanagers.ConfigurationType, 0)
	for _, v := range input {
		result = append(result, networkmanagers.ConfigurationType(v))
	}
	return &result
}
