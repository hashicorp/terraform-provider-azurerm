// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/virtualnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resourcegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ResourceGroupResource struct{}

func TestAccResourceGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")
	testResource := ResourceGroupResource{}
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		data.ApplyStep(testResource.basic, testResource),
		data.ImportStep(),
	})
}

func TestAccResourceGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")
	testResource := ResourceGroupResource{}
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		data.ApplyStep(testResource.basic, testResource),
		data.RequiresImportErrorStep(testResource.requiresImportConfig),
	})
}

func TestAccResourceGroup_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")
	testResource := ResourceGroupResource{}
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       testResource.basic,
			TestResource: testResource,
		}),
	})
}

func TestAccResourceGroup_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")
	testResource := ResourceGroupResource{}
	assert := check.That(data.ResourceName)
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		{
			Config: testResource.withTagsConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				assert.ExistsInAzure(testResource),
				assert.Key("tags.%").HasValue("2"),
				assert.Key("tags.cost_center").HasValue("MSFT"),
				assert.Key("tags.environment").HasValue("Production"),
			),
		},
		data.ImportStep(),
		{
			Config: testResource.withTagsUpdatedConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				assert.ExistsInAzure(testResource),
				assert.Key("tags.%").HasValue("1"),
				assert.Key("tags.environment").HasValue("staging"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceGroup_withManagedBy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")
	testResource := ResourceGroupResource{}
	assert := check.That(data.ResourceName)
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		{
			Config: testResource.withManagedByConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				assert.ExistsInAzure(testResource),
				assert.Key("managed_by").HasValue("test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceGroup_withNestedItemsAndFeatureFlag(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")
	r := ResourceGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withFeatureFlag(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// since we don't want to track/destroy this resource for test purposes, we can create this here
				// it'll be cleaned up in the final step with the feature flag disabled, so this should be fine.
				data.CheckWithClient(r.createNetworkOutsideTerraform(fmt.Sprintf("acctestvnet-%d", data.RandomInteger))),
			),
		},
		data.ImportStep(),
		{
			// attempting to delete this with the vnet should error
			Config:      r.withFeatureFlag(data, true),
			Destroy:     true,
			ExpectError: regexp.MustCompile("This feature is intended to avoid the unintentional destruction"),
		},
		{
			// with the feature disabled we should delete the RG and the Network
			Config:  r.withFeatureFlag(data, false),
			Destroy: true,
		},
	})
}

func (r ResourceGroupResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseResourceGroupIDInsensitively(state.ID)
	if err != nil {
		return nil, err
	}

	opts := resourcegroups.DefaultDeleteOperationOptions()
	opts.ForceDeletionTypes = pointer.To("Microsoft.Compute/virtualMachines,Microsoft.Compute/virtualMachineScaleSets")
	if resp, err := client.Resource.ResourceGroupsClient.Delete(ctx, *id, opts); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("deleting test %s: %+v", *id, err)
		}
	} else {
		if err := resp.Poller.PollUntilDone(ctx); err != nil {
			return nil, fmt.Errorf("polling deleting %s: %+v", *id, err)
		}
	}

	return pointer.To(true), nil
}

func (r ResourceGroupResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseResourceGroupIDInsensitively(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Resource.ResourceGroupsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ResourceGroupResource) createNetworkOutsideTerraform(name string) func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		client := clients.Network.VirtualNetworks

		id, err := commonids.ParseResourceGroupID(state.ID)
		if err != nil {
			return err
		}

		params := virtualnetworks.VirtualNetwork{
			Location: pointer.To(state.Attributes["location"]),
			Properties: &virtualnetworks.VirtualNetworkPropertiesFormat{
				AddressSpace: &virtualnetworks.AddressSpace{
					AddressPrefixes: &[]string{
						"10.0.0.0/16",
					},
				},
			},
		}
		vnetId := commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, name)

		ctx2, cancel := context.WithTimeout(ctx, 30*time.Minute)
		defer cancel()
		if err := client.CreateOrUpdateThenPoll(ctx2, vnetId, params); err != nil {
			return fmt.Errorf("creating nested virtual network: %+v", err)
		}

		return nil
	}
}

func (r ResourceGroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ResourceGroupResource) requiresImportConfig(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group" "import" {
  name     = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
}
`, template)
}

func (r ResourceGroupResource) withFeatureFlag(data acceptance.TestData, featureFlagEnabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = %t
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
`, featureFlagEnabled, data.RandomInteger, data.Locations.Primary)
}

func (r ResourceGroupResource) withTagsConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ResourceGroupResource) withTagsUpdatedConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ResourceGroupResource) withManagedByConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"

  managed_by = "test"
}
`, data.RandomInteger, data.Locations.Primary)
}
