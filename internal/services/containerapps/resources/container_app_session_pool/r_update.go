// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_session_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/containerappssessionpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.SessionPoolClient

			id, err := containerappssessionpools.ParseSessionPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config ContainerAppSessionPoolModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			payload := *existing.Model
			props := payload.Properties

			if metadata.ResourceData.HasChange("custom_container_template") {
				customContainerTemplate, err := expandContainerAppSessionPoolCustomContainerTemplate(config.CustomContainerTemplate)
				if err != nil {
					return fmt.Errorf("expanding `custom_container_template`: %+v", err)
				}
				props.CustomContainerTemplate = customContainerTemplate
			}

			if metadata.ResourceData.HasChange("lifecycle_configuration") {
				props.DynamicPoolConfiguration = expandContainerAppSessionPoolDynamicPoolConfiguration(config.LifecycleConfiguration)
			}

			if metadata.ResourceData.HasChanges("max_concurrent_sessions", "ready_session_instances") {
				props.ScaleConfiguration = expandContainerAppSessionPoolScaleConfiguration(config)
			}

			// clearing this sends a null, which the API resets to the `None` lifecycle stage
			if metadata.ResourceData.HasChange("session_managed_identities") {
				props.ManagedIdentitySettings = expandContainerAppSessionPoolManagedIdentitySettings(config.SessionManagedIdentities)
			}

			if metadata.ResourceData.HasChange("network_egress_enabled") {
				networkStatus := containerappssessionpools.SessionNetworkStatusEgressDisabled
				if config.NetworkEgressEnabled {
					networkStatus = containerappssessionpools.SessionNetworkStatusEgressEnabled
				}

				props.SessionNetworkConfiguration = &containerappssessionpools.SessionNetworkConfiguration{
					Status: pointer.To(networkStatus),
				}
			}

			// The API doesn't return secret values, so these are always sent from the config to avoid clearing them.
			props.Secrets = expandContainerAppSessionPoolSecrets(config.Secrets)

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(config.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(config.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
