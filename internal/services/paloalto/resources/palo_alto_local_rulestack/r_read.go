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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestackResources

			id, err := localrulestackresources.ParseLocalRulestackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LocalRuleStackModel

			existing, err := client.LocalRulestacksGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state.Name = id.LocalRulestackName
			state.ResourceGroupName = id.ResourceGroupName
			if model := existing.Model; model != nil {
				props := model.Properties

				state.Description = pointer.From(props.Description)
				state.Location = location.Normalize(existing.Model.Location)

				if secServices := props.SecurityServices; secServices != nil {
					if v := pointer.From(secServices.VulnerabilityProfile); v != RuleStackSecurityServicesNone {
						state.VulnerabilityProfile = v
					}
					if v := pointer.From(secServices.AntiSpywareProfile); v != RuleStackSecurityServicesNone {
						state.AntiSpywareProfile = v
					}
					if v := pointer.From(secServices.AntiVirusProfile); v != RuleStackSecurityServicesNone {
						state.AntiVirusProfile = v
					}
					if v := pointer.From(secServices.FileBlockingProfile); v != RuleStackSecurityServicesNone {
						state.FileBlockingProfile = v
					}
					if v := pointer.From(secServices.UrlFilteringProfile); v != RuleStackSecurityServicesNone {
						state.URLFilteringProfile = v
					}
					if v := pointer.From(secServices.DnsSubscription); v != RuleStackSecurityServicesNone {
						state.DNSSubscription = v
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

const (
	RuleStackSecurityServicesCustom       string = "Custom"
	RuleStackSecurityServicesNone         string = "None"
	RuleStackSecurityServicesBestPractice string = "BestPractice"
)
