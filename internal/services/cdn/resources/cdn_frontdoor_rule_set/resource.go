// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_rule_set

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorRuleSetCreate,
		Read:   resourceCdnFrontDoorRuleSetRead,
		Delete: resourceCdnFrontDoorRuleSetDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(4 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := rulesets.ParseRuleSetID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta any) ([]*pluginsdk.ResourceData, error) {
			client := meta.(*clients.Client).Cdn.FrontDoorRuleSetsClient

			id, _ := rulesets.ParseRuleSetID(d.Id())
			resp, err := client.Get(ctx, *id)
			if err != nil {
				return nil, fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return nil, fmt.Errorf("retrieving %s: `model` was nil`", id)
			}

			if resp.Model.Properties == nil {
				return nil, fmt.Errorf("retrieving %s: `properties` was nil`", id)
			}

			if pointer.From(resp.Model.Properties.BatchMode) {
				return nil, fmt.Errorf("%s was provisioned using batch mode and cannot be managed by this resource, use `azurerm_cdn_frontdoor_batch_rule_set` instead", id)
			}

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Schema: cdnFrontdoorRuleSetSchema(),
	}
}
