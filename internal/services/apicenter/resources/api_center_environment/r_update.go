// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_environment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/helpers"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiCenter.EnvironmentsClient
			id, err := environments.ParseEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ApiCenterEnvironmentResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			model := existing.Model

			if metadata.ResourceData.HasChange("custom_properties") {
				model.Properties.CustomProperties, err = helpers.ExpandCustomProperties(state.CustomProperties)
				if err != nil {
					return err
				}
			}

			if metadata.ResourceData.HasChange("description") {
				model.Properties.Description = pointer.To(state.Description)
			}

			if metadata.ResourceData.HasChange("development_portal_uri") {
				model.Properties.Onboarding.DeveloperPortalUri = pointer.To([]string{state.DevPortalUri})
			}

			if metadata.ResourceData.HasChange("server_type") {
				model.Properties.Server.Type = pointer.ToEnum[environments.EnvironmentServerType](state.ServerType)
			}

			if metadata.ResourceData.HasChange("management_portal_uri") {
				model.Properties.Server.ManagementPortalUri = pointer.To([]string{state.MgmtPortalUri})
			}

			if metadata.ResourceData.HasChange("instructions") {
				model.Properties.Onboarding.Instructions = pointer.To(state.Instructions)
			}

			if _, err = client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
