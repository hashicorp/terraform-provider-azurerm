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
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	eventgridClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/client"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := eventgridClient.From(metadata.Client).PartnerConfigurations

			subscriptionId := account.From(metadata.Client).SubscriptionId

			var config EventGridPartnerConfigurationResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Uses resource group ID as it has 1-1 mapping with resource group
			// See the SDK sample usage: https://github.com/hashicorp/go-azure-sdk/tree/main/resource-manager/eventgrid/2025-02-15/partnerconfigurations
			id := commonids.NewResourceGroupID(subscriptionId, config.ResourceGroup)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
		partner := partnerconfigurations.Partner{
			AuthorizationExpirationTimeInUtc: pointer.To(partnerAuth.AuthorizationExpirationTimeInUtc),
		}
		if partnerAuth.PartnerRegistrationId != "" {
			partner.PartnerRegistrationImmutableId = pointer.To(partnerAuth.PartnerRegistrationId)
		}
		if partnerAuth.PartnerName != "" {
			partner.PartnerName = pointer.To(partnerAuth.PartnerName)
		}
		partners = append(partners, partner)
	}

	return &partners
}
