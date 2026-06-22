// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resourcegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// TestAccResourceGroup_providerIgnoreTags verifies that the provider-level
// `ignore_tags` block prevents tags applied out-of-band (here, simulated via a
// direct API call) from producing a perpetual plan diff, when their key matches
// `keys` exactly or starts with one of `key_prefixes`.
func TestAccResourceGroup_providerIgnoreTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")
	r := ResourceGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.providerIgnoreTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
				// apply tags out-of-band whose keys are ignored by the provider config
				data.CheckWithClient(r.setTagsOutOfBand(map[string]string{
					"ignoreThisKey":      "set-outside-terraform",
					"ignore-prefix-team": "platform",
				})),
			),
		},
		{
			// because the out-of-band tags match the provider `ignore_tags` block they
			// are scrubbed from state on read, so the plan must be empty (no drift)
			Config:             r.providerIgnoreTags(data),
			PlanOnly:           true,
			ExpectNonEmptyPlan: false,
		},
	})
}

func (r ResourceGroupResource) providerIgnoreTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}

  ignore_tags {
    keys         = ["ignoreThisKey"]
    key_prefixes = ["ignore-prefix-"]
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

// setTagsOutOfBand merges the supplied tags into the resource group via a direct
// API call, simulating tags applied by an external system such as Azure Policy.
func (r ResourceGroupResource) setTagsOutOfBand(extra map[string]string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
		id, err := commonids.ParseResourceGroupIDInsensitively(state.ID)
		if err != nil {
			return err
		}

		existing, err := client.Resource.ResourceGroupsClient.Get(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}
		if existing.Model == nil {
			return fmt.Errorf("retrieving %s: `model` was nil", *id)
		}

		newTags := map[string]string{}
		if existing.Model.Tags != nil {
			for k, v := range *existing.Model.Tags {
				newTags[k] = v
			}
		}
		for k, v := range extra {
			newTags[k] = v
		}

		patch := resourcegroups.ResourceGroupPatchable{
			Tags: &newTags,
		}
		if _, err := client.Resource.ResourceGroupsClient.Update(ctx, *id, patch); err != nil {
			return fmt.Errorf("setting tags out-of-band on %s: %+v", *id, err)
		}

		return nil
	}
}
