// Copyright IBM Corp. 2014, 2025
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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.PartnerConfigurations

			subscriptionId := metadata.Client.Account.SubscriptionId

			var config EventGridPartnerConfigurationResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Uses resource group ID as it has 1-1 mapping with resource group
			// See the SDK sample usage: https://github.com/hashicorp/go-azure-sdk/tree/main/resource-manager/eventgrid/2025-02-15/partnerconfigurations
			id := commonids.NewResourceGroupID(subscriptionId, config.ResourceGroup)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			param := partnerconfigurations.PartnerConfiguration{
				Location: pointer.To("global"),
				Properties: &partnerconfigurations.PartnerConfigurationProperties{
					PartnerAuthorization: &partnerconfigurations.PartnerAuthorization{
						DefaultMaximumExpirationTimeInDays: pointer.To(config.DefaultMaximumExpirationTimeInDays),
						AuthorizedPartnersList:             expandAuthorizedPartnersList(config.PartnerAuthorizations),
					},
				},
				Tags: pointer.To(config.Tags),
			}
			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, param, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandAuthorizedPartnersList(partnerAuthorization []PartnerAuthorization) *[]partnerconfigurations.Partner {
	if len(partnerAuthorization) == 0 {
		return nil
	}

	partners := []partnerconfigurations.Partner{}

	for _, partnerAuth := range partnerAuthorization {
		partners = append(partners, partnerconfigurations.Partner{
			PartnerName:                      pointer.To(partnerAuth.PartnerName),
			PartnerRegistrationImmutableId:   pointer.To(partnerAuth.PartnerRegistrationId),
			AuthorizationExpirationTimeInUtc: pointer.To(partnerAuth.AuthorizationExpirationTimeInUtc),
		})
	}

	return &partners
}
