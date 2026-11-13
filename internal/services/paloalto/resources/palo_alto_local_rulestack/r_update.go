// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulestackresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestackResources

			id, err := localrulestackresources.ParseLocalRulestackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			model := LocalRuleStackModel{}

			if err = metadata.Decode(&model); err != nil {
				return err
			}

			existing, err := client.LocalRulestacksGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			localRuleStack := *existing.Model
			update := localRuleStack.Properties

			if metadata.ResourceData.HasChange("description") {
				update.Description = pointer.To(model.Description)
			}

			secServices := pointer.From(update.SecurityServices)

			if metadata.ResourceData.HasChange("dns_subscription") {
				if model.DNSSubscription != "" {
					secServices.DnsSubscription = pointer.To(model.DNSSubscription)
				} else {
					secServices.DnsSubscription = pointer.To(RuleStackSecurityServicesNone)
				}
			}

			if metadata.ResourceData.HasChange("vulnerability_profile") {
				if model.VulnerabilityProfile != "" {
					secServices.VulnerabilityProfile = pointer.To(model.VulnerabilityProfile)
				} else {
					secServices.VulnerabilityProfile = pointer.To(RuleStackSecurityServicesNone)
				}
			}

			if metadata.ResourceData.HasChange("anti_spyware_profile") {
				if model.AntiSpywareProfile != "" {
					secServices.AntiSpywareProfile = pointer.To(model.AntiSpywareProfile)
				} else {
					secServices.AntiSpywareProfile = pointer.To(RuleStackSecurityServicesNone)
				}
			}

			if metadata.ResourceData.HasChange("anti_virus_profile") {
				if model.AntiVirusProfile != "" {
					secServices.AntiVirusProfile = pointer.To(model.AntiVirusProfile)
				} else {
					secServices.AntiVirusProfile = pointer.To(RuleStackSecurityServicesNone)
				}
			}

			if metadata.ResourceData.HasChange("url_filtering_profile") {
				if model.URLFilteringProfile != "" {
					secServices.UrlFilteringProfile = pointer.To(model.URLFilteringProfile)
				} else {
					secServices.UrlFilteringProfile = pointer.To(RuleStackSecurityServicesNone)
				}
			}

			if metadata.ResourceData.HasChange("file_blocking_profile") {
				if model.FileBlockingProfile != "" {
					secServices.FileBlockingProfile = pointer.To(model.FileBlockingProfile)
				} else {
					secServices.FileBlockingProfile = pointer.To(RuleStackSecurityServicesNone)
				}
			}

			update.SecurityServices = pointer.To(secServices)

			localRuleStack.Properties = update

			if err = client.LocalRulestacksCreateOrUpdateThenPoll(ctx, *id, localRuleStack); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			if err = client.LocalRulestackscommitThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("committing config for %s: %+v", *id, err)
			}

			return nil
		},
	}
}
