// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package active_directory_domain_service_trust

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aad/2021-05-01/domainservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/parse"
)

func (r Resource) Update() sdk.ResourceFunc {
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

			var plan DomainServiceTrustModel
			if err := metadata.Decode(&plan); err != nil {
				return err
			}

			var found bool
			newTrusts := []domainservices.ForestTrust{}
			for _, trust := range existingTrusts {
				if trust.FriendlyName != nil && *trust.FriendlyName == id.TrustName {
					found = true
					if metadata.ResourceData.HasChange("trusted_domain_fqdn") {
						trust.TrustedDomainFqdn = pointer.To(plan.TrustedDomainFqdn)
					}
					if metadata.ResourceData.HasChange("trusted_domain_dns_ips") {
						trust.RemoteDnsIPs = pointer.To(strings.Join(plan.TrustedDomainDnsIPs, ","))
					}
					trust.TrustPassword = pointer.To(plan.Password)
				}
				newTrusts = append(newTrusts, trust)
			}
			if !found {
				return fmt.Errorf("%s not exists: %+v", id, err)
			}

			params := domainservices.DomainService{
				Properties: &domainservices.DomainServiceProperties{
					ResourceForestSettings: &domainservices.ResourceForestSettings{
						Settings: &newTrusts,
					},
				},
			}

			if err := client.UpdateThenPoll(ctx, idsdk, params); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}
