// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package chaos_studio_experiment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/experiments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ChaosStudio.V20231101.Experiments

			id, err := experiments.ParseExperimentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config ChaosStudioExperimentResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving existing %s: properties was nil", *id)
			}
			payload := *existing.Model

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandSystemOrUserAssignedMapFromModel(config.Identity)
				if err != nil {
					return fmt.Errorf("expanding SystemOrUserAssigned Identity: %+v", err)
				}
				payload.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("selectors") {
				selectors, err := expandSelectors(config.Selectors)
				if err != nil {
					return fmt.Errorf("expanding `selectors`: %+v", err)
				}
				payload.Properties.Selectors = *selectors
			}

			if metadata.ResourceData.HasChange("steps") {
				steps, err := expandSteps(config.Steps)
				if err != nil {
					return fmt.Errorf("expanding `steps`: %+v", err)
				}
				payload.Properties.Steps = *steps
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
