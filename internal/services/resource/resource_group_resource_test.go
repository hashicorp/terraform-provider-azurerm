// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

type ResourceGroupResource struct{}

func TestAccResourceGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")
	testResource := ResourceGroupResource{}
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		data.ApplyStep(testResource.basicConfig, testResource),
		data.ImportStep(),
	})
}

func TestAccResourceGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")
	testResource := ResourceGroupResource{}
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		data.ApplyStep(testResource.basicConfig, testResource),
		data.RequiresImportErrorStep(testResource.requiresImportConfig),
	})
}

func TestAccResourceGroup_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")
	testResource := ResourceGroupResource{}
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       testResource.basicConfig,
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

func (t ResourceGroupResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	resourceGroup := state.Attributes["name"]

	groupsClient := client.Resource.GroupsClient
	deleteFuture, err := groupsClient.Delete(ctx, resourceGroup, "Microsoft.Compute/virtualMachines,Microsoft.Compute/virtualMachineScaleSets")
	if err != nil {
		return nil, fmt.Errorf("deleting Resource Group %q: %+v", resourceGroup, err)
	}

	err = deleteFuture.WaitForCompletionRef(ctx, groupsClient.Client)
	if err != nil {
		return nil, fmt.Errorf("waiting for deletion of Resource Group %q: %+v", resourceGroup, err)
	}

	return utils.Bool(true), nil
}

func (t ResourceGroupResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	name := state.Attributes["name"]

	resp, err := client.Resource.GroupsClient.Get(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Resource Group %q: %+v", name, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (t ResourceGroupResource) createNetworkOutsideTerraform(name string) func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		client := clients.Network.VnetClient
		resourceGroup := state.Attributes["name"]
		location := state.Attributes["location"]
		params := network.VirtualNetwork{
			Location: utils.String(location),
			VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
				AddressSpace: &network.AddressSpace{
					AddressPrefixes: &[]string{
						"10.0.0.0/16",
					},
				},
			},
		}
		future, err := client.CreateOrUpdate(ctx, resourceGroup, name, params)
		if err != nil {
			return fmt.Errorf("creating nested virtual network: %+v", err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for the creation of nested virtual network: %+v", err)
		}

		return nil
	}
}

func (t ResourceGroupResource) basicConfig(data acceptance.TestData) string {
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

func (t ResourceGroupResource) requiresImportConfig(data acceptance.TestData) string {
	template := t.basicConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group" "import" {
  name     = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
}
`, template)
}

func (t ResourceGroupResource) withFeatureFlag(data acceptance.TestData, featureFlagEnabled bool) string {
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

func (t ResourceGroupResource) withTagsConfig(data acceptance.TestData) string {
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

func (t ResourceGroupResource) withTagsUpdatedConfig(data acceptance.TestData) string {
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

func (t ResourceGroupResource) withManagedByConfig(data acceptance.TestData) string {
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
