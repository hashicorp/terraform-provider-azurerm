// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_dapr_component

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/daprcomponents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.DaprComponentsClient
			id, err := daprcomponents.ParseDaprComponentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ContainerAppEnvironmentDaprComponentModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil || existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s for update: %+v", *id, err)
			}

			// Populate the secrets from the List API to prevent accidental removal.
			secretsResp, err := client.ListSecrets(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving secrets for %s: %+v", *id, err)
			}

			existing.Model.Properties.Secrets = helpers.UnpackContainerDaprSecretsCollection(secretsResp.Model)

			if metadata.ResourceData.HasChange("version") {
				existing.Model.Properties.Version = pointer.To(state.Version)
			}

			if metadata.ResourceData.HasChange("init_timeout") {
				existing.Model.Properties.InitTimeout = pointer.To(state.InitTimeout)
			}

			if metadata.ResourceData.HasChange("ignore_errors") {
				existing.Model.Properties.IgnoreErrors = pointer.To(state.IgnoreErrors)
			}

			if metadata.ResourceData.HasChange("secret") {
				existing.Model.Properties.Secrets = helpers.ExpandDaprSecrets(state.Secrets)
			}

			if metadata.ResourceData.HasChange("metadata") {
				existing.Model.Properties.Metadata = expandDaprComponentPropertiesMetadata(state.Metadata)
			}

			if metadata.ResourceData.HasChange("scopes") {
				existing.Model.Properties.Scopes = pointer.To(state.Scopes)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
