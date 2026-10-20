// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_profile

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorProfileCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.FrontDoorProfilesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := profiles.NewProfileID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_profile", id.ID())
		}
	}

	props := profiles.Profile{
		Location: location.Normalize("global"),
		Properties: &profiles.ProfileProperties{
			OriginResponseTimeoutSeconds: pointer.To(int64(d.Get("response_timeout_seconds").(int))),
		},
		Sku: profiles.Sku{
			Name: pointer.ToEnum[profiles.SkuName](d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	props.Properties.LogScrubbing = expandCdnFrontDoorProfileLogScrubbing(d.Get("log_scrubbing_rule").(*pluginsdk.Set).List())

	if v, ok := d.GetOk("identity"); ok {
		i, err := identity.ExpandLegacySystemAndUserAssignedMap(v.([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		props.Identity = i
	}

	if err := client.CreateCallbackThenPoll(ctx, id, props, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnFrontDoorProfileRead(d, meta)
}

func expandCdnFrontDoorProfileLogScrubbing(input []any) *profiles.ProfileLogScrubbing {
	if len(input) == 0 {
		return &profiles.ProfileLogScrubbing{
			State: pointer.To(profiles.ProfileScrubbingStateDisabled),
		}
	}

	return &profiles.ProfileLogScrubbing{
		State:          pointer.To(profiles.ProfileScrubbingStateEnabled),
		ScrubbingRules: expandCdnFrontDoorProfileLogScrubbingRules(input),
	}
}

func expandCdnFrontDoorProfileLogScrubbingRules(input []any) *[]profiles.ProfileScrubbingRules {
	if len(input) == 0 {
		return nil
	}

	scrubbingRules := make([]profiles.ProfileScrubbingRules, 0)

	for _, rule := range input {
		v := rule.(map[string]any)

		item := profiles.ProfileScrubbingRules{
			MatchVariable:         profiles.ScrubbingRuleEntryMatchVariable(v["match_variable"].(string)),
			SelectorMatchOperator: profiles.ScrubbingRuleEntryMatchOperatorEqualsAny, // EqualsAny is the only valid SelectorMatchOperator for log scrubbing in the Profile API
			State:                 pointer.To(profiles.ScrubbingRuleEntryStateEnabled),
		}

		scrubbingRules = append(scrubbingRules, item)
	}

	return &scrubbingRules
}
