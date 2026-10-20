// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/partnerconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.PartnerConfigurations

			id, err := commonids.ParseResourceGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config EventGridPartnerConfigurationResourceModel
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
				return fmt.Errorf("retrieving %s: `model.Properties` was nil", *id)
			}

			payload := existing.Model

			if metadata.ResourceData.HasChange("default_maximum_expiration_time_in_days") {
				payload.Properties.PartnerAuthorization.DefaultMaximumExpirationTimeInDays = pointer.To(config.DefaultMaximumExpirationTimeInDays)
			}

			if metadata.ResourceData.HasChange("partner_authorization") {
				if payload.Properties.PartnerAuthorization == nil {
					payload.Properties.PartnerAuthorization = &partnerconfigurations.PartnerAuthorization{}
				}
				payload.Properties.PartnerAuthorization.AuthorizedPartnersList = expandAuthorizedPartnersList(config.PartnerAuthorizations)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(config.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
	}
}
