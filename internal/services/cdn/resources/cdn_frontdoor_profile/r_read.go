// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_profile

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorProfileRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.FrontDoorProfilesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := profiles.ParseProfileID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ProfileName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if skuName := model.Sku.Name; skuName != nil {
			d.Set("sku_name", pointer.FromEnum(skuName))
		}

		identity, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}

		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("response_timeout_seconds", int(pointer.From(props.OriginResponseTimeoutSeconds)))

			// whilst this is returned in the API as FrontDoorID other resources refer to
			// this as the Resource GUID, so we will for consistency
			d.Set("resource_guid", pointer.From(props.FrontDoorId))

			if err := d.Set("log_scrubbing_rule", flattenCdnFrontDoorProfileLogScrubbingRules(props.LogScrubbing)); err != nil {
				return fmt.Errorf("setting `log_scrubbing_rule`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func flattenCdnFrontDoorProfileLogScrubbingRules(input *profiles.ProfileLogScrubbing) []any {
	result := make([]any, 0)

	if input == nil || pointer.From(input.State) == profiles.ProfileScrubbingStateDisabled {
		return result
	}

	if input.ScrubbingRules == nil {
		return result
	}

	for _, scrubbingRule := range *input.ScrubbingRules {
		if scrubbingRule.State != nil && pointer.From(scrubbingRule.State) == profiles.ScrubbingRuleEntryStateEnabled {
			result = append(result, map[string]any{
				"match_variable": scrubbingRule.MatchVariable,
			})
		}
	}

	return result
}
