// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_batch_rule_set

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.ResourceWithCustomImporter = Resource{}
	_ sdk.ResourceWithCustomizeDiff  = Resource{}
	_ sdk.ResourceWithIdentity       = Resource{}
	_ sdk.ResourceWithUpdate         = Resource{}
)

type Resource struct{}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := metadata.Client.Cdn.FrontDoorRuleSetsClient

		id, err := rulesets.ParseRuleSetID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.Model == nil {
			return fmt.Errorf("retrieving %s: `model` was nil`", id)
		}

		if resp.Model.Properties == nil {
			return fmt.Errorf("retrieving %s: `properties` was nil`", id)
		}

		if !pointer.From(resp.Model.Properties.BatchMode) {
			return fmt.Errorf("%s was not provisioned using batch mode and cannot be managed by this resource, use `azurerm_cdn_frontdoor_rule_set` instead", id)
		}

		return nil
	}
}

func (Resource) Identity() resourceids.ResourceId {
	return &rulesets.RuleSetId{}
}

func (Resource) ResourceType() string {
	return "azurerm_cdn_frontdoor_batch_rule_set"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return rulesets.ValidateRuleSetID
}
