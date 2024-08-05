// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package chaosstudio_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ChaosStudioCapabilityTestResource struct{}

func TestAccChaosStudioCapability_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_chaos_studio_capability", "test")
	r := ChaosStudioCapabilityTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccChaosStudioCapability_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_chaos_studio_capability", "test")
	r := ChaosStudioCapabilityTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccChaosStudioCapability_multipleCapabilities(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_chaos_studio_capability", "test")
	r := ChaosStudioCapabilityTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ChaosStudioCapabilityTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseChaosStudioCapabilityID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ChaosStudio.V20231101.Capabilities.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}
func (r ChaosStudioCapabilityTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_chaos_studio_capability" "test" {
  chaos_studio_target_id = azurerm_chaos_studio_target.test.id
  capability_type        = "NetworkChaos-2.0"
}
`, r.template(data))
}

func (r ChaosStudioCapabilityTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_chaos_studio_capability" "import" {
  chaos_studio_target_id = azurerm_chaos_studio_capability.test.chaos_studio_target_id
  capability_type        = azurerm_chaos_studio_capability.test.capability_type
}
`, r.basic(data))
}

func (r ChaosStudioCapabilityTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_chaos_studio_capability" "another" {
  chaos_studio_target_id = azurerm_storage_account.test.id
  capability_type        = "NetworkChaos-2.0"
}

resource "azurerm_chaos_studio_capability" "test" {
  chaos_studio_target_id = azurerm_storage_account.test.id
  capability_type        = "PodChaos-2.1"
}
`, r.template(data))
}

func (r ChaosStudioCapabilityTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "primary_location" {
  default = %q
}
variable "random_integer" {
  default = %d
}
variable "random_string" {
  default = %q
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-${var.random_integer}"
  location = var.primary_location
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks${var.random_string}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks${var.random_string}"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_chaos_studio_target" "test" {
  location           = azurerm_resource_group.test.location
  target_resource_id = azurerm_kubernetes_cluster.test.id
  target_type        = "Microsoft-AzureKubernetesServiceChaosMesh"
}

`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
