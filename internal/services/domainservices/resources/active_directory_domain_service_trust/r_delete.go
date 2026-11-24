// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package active_directory_domain_service_trust

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aad/2021-05-01/domainservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/parse"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DomainServices.DomainServicesClient

			id, err := parse.DomainServiceTrustID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			idsdk := domainservices.NewDomainServiceID(id.SubscriptionId, id.ResourceGroup, id.DomainServiceName)

			locks.ByName(id.DomainServiceName, DomainServiceResourceName)
			defer locks.UnlockByName(id.DomainServiceName, DomainServiceResourceName)

			existing, err := client.Get(ctx, idsdk)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", idsdk, err)
			}

			model := existing.Model
			if model == nil {
				return fmt.Errorf("reading %s: returned with null model", idsdk)
			}
			existingTrusts := []domainservices.ForestTrust{}
			if props := model.Properties; props != nil {
				if fsettings := props.ResourceForestSettings; fsettings != nil {
					if settings := fsettings.Settings; settings != nil {
						existingTrusts = *settings
					}
				}
			}
			var found bool
			newTrusts := []domainservices.ForestTrust{}
			for _, trust := range existingTrusts {
				if trust.FriendlyName != nil && *trust.FriendlyName == id.TrustName {
					found = true
					continue
				}
				newTrusts = append(newTrusts, trust)
			}

			if !found {
				return metadata.MarkAsGone(id)
			}

			params := domainservices.DomainService{
				Properties: &domainservices.DomainServiceProperties{
					ResourceForestSettings: &domainservices.ResourceForestSettings{
						Settings: &newTrusts,
					},
				},
			}

			if err := client.UpdateThenPoll(ctx, idsdk, params); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
