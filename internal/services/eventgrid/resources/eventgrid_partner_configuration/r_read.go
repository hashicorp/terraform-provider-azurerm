// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/partnerconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.PartnerConfigurations

			id, err := commonids.ParseResourceGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := EventGridPartnerConfigurationResourceModel{
				ResourceGroup: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil && props.PartnerAuthorization != nil {
					state.DefaultMaximumExpirationTimeInDays = pointer.From(props.PartnerAuthorization.DefaultMaximumExpirationTimeInDays)
					state.PartnerAuthorizations = flattenAuthorizedPartnersList(props.PartnerAuthorization.AuthorizedPartnersList)
				}

				if model.Tags != nil {
					state.Tags = pointer.From(model.Tags)
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func flattenAuthorizedPartnersList(partner *[]partnerconfigurations.Partner) []PartnerAuthorization {
	partnerAuthorizations := make([]PartnerAuthorization, 0)

	if partner == nil {
		return partnerAuthorizations
	}

	for _, partnerAuth := range *partner {
		partnerAuthorizations = append(partnerAuthorizations, PartnerAuthorization{
			PartnerName:                      pointer.From(partnerAuth.PartnerName),
			PartnerRegistrationId:            pointer.From(partnerAuth.PartnerRegistrationImmutableId),
			AuthorizationExpirationTimeInUtc: pointer.From(partnerAuth.AuthorizationExpirationTimeInUtc),
		})
	}

	return partnerAuthorizations
}

type PartnerAuthorization struct {
	PartnerRegistrationId            string `tfschema:"partner_registration_id"`
	PartnerName                      string `tfschema:"partner_name"`
	AuthorizationExpirationTimeInUtc string `tfschema:"authorization_expiration_time_in_utc"`
}
