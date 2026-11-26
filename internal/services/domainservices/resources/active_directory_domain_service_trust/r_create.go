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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DomainServices.DomainServicesClient

			var plan DomainServiceTrustModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			dsid, err := parse.DomainServiceID(plan.DomainServiceId)
			if err != nil {
				return err
			}

			id := parse.NewDomainServiceTrustID(dsid.SubscriptionId, dsid.ResourceGroup, dsid.Name, plan.Name)
			idsdk := domainservices.NewDomainServiceID(id.SubscriptionId, id.ResourceGroup, id.DomainServiceName)

			locks.ByName(id.DomainServiceName, DomainServiceResourceName)
			defer locks.UnlockByName(id.DomainServiceName, DomainServiceResourceName)

			existing, err := client.Get(ctx, idsdk)
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
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

			exists := false
			for _, setting := range existingTrusts {
				if setting.FriendlyName != nil && *setting.FriendlyName == id.TrustName {
					exists = true
					if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
						return metadata.ResourceRequiresImport(r.ResourceType(), id)
					}
				}
			}

			if !exists {
				existingTrusts = append(existingTrusts, domainservices.ForestTrust{
					TrustedDomainFqdn: pointer.To(plan.TrustedDomainFqdn),
					TrustDirection:    pointer.To("Inbound"),
					FriendlyName:      pointer.To(id.TrustName),
					RemoteDnsIPs:      pointer.To(strings.Join(plan.TrustedDomainDnsIPs, ",")),
					TrustPassword:     pointer.To(plan.Password),
				})
			}

			params := domainservices.DomainService{
				Properties: &domainservices.DomainServiceProperties{
					ResourceForestSettings: &domainservices.ResourceForestSettings{
						Settings: &existingTrusts,
					},
				},
			}

			// TODO: implement `CallbackThenPoll`, requires migrating to an ID that implements `resourceids.ResourceId`
			if err := client.UpdateThenPoll(ctx, idsdk, params); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
