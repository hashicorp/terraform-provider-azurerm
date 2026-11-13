// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulestackresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestackResources

			model := LocalRuleStackModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := localrulestackresources.NewLocalRulestackID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)
			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.LocalRulestacksGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			secServices := localrulestackresources.SecurityServices{
				AntiSpywareProfile:   pointer.To(RuleStackSecurityServicesNone),
				AntiVirusProfile:     pointer.To(RuleStackSecurityServicesNone),
				DnsSubscription:      pointer.To(RuleStackSecurityServicesNone),
				FileBlockingProfile:  pointer.To(RuleStackSecurityServicesNone),
				UrlFilteringProfile:  pointer.To(RuleStackSecurityServicesNone),
				VulnerabilityProfile: pointer.To(RuleStackSecurityServicesNone),
			}

			if model.AntiSpywareProfile != "" {
				secServices.AntiSpywareProfile = pointer.To(model.AntiSpywareProfile)
			}
			if model.AntiVirusProfile != "" {
				secServices.AntiVirusProfile = pointer.To(model.AntiVirusProfile)
			}
			if model.DNSSubscription != "" {
				secServices.DnsSubscription = pointer.To(model.DNSSubscription)
			}
			if model.FileBlockingProfile != "" {
				secServices.FileBlockingProfile = pointer.To(model.FileBlockingProfile)
			}
			if model.URLFilteringProfile != "" {
				secServices.UrlFilteringProfile = pointer.To(model.URLFilteringProfile)
			}
			if model.VulnerabilityProfile != "" {
				secServices.VulnerabilityProfile = pointer.To(model.VulnerabilityProfile)
			}

			localRuleStack := localrulestackresources.LocalRulestackResource{
				Location: location.Normalize(model.Location),
				Properties: localrulestackresources.RulestackProperties{
					DefaultMode:      pointer.To(localrulestackresources.DefaultModeNONE),
					Description:      pointer.To(model.Description),
					Scope:            pointer.To(localrulestackresources.ScopeTypeLOCAL),
					SecurityServices: pointer.To(secServices),
				},
			}

			if err := client.LocalRulestacksCreateOrUpdateCallbackThenPoll(ctx, id, localRuleStack, metadata.SetIDCallback(&id)); err != nil {
				return err
			}
			metadata.SetID(id)

			return nil
		},
	}
}
