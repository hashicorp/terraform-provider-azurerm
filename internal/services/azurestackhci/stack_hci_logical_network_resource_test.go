// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurestackhci_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/logicalnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StackHCILogicalNetworkResource struct{}

// https://learn.microsoft.com/en-us/azure-stack/hci/manage/create-logical-networks?tabs=azurecli#prerequisites
// The resource can only be created on the customlocation generated after HCI deployment
const (
	customLocationIdEnv = "ARM_TEST_STACK_HCI_CUSTOM_LOCATION_ID"
)

func TestAccStackHCILogicalNetwork_dynamic(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_logical_network", "test")
	r := StackHCILogicalNetworkResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dynamic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStackHCILogicalNetwork_update(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_logical_network", "test")
	r := StackHCILogicalNetworkResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStackHCILogicalNetwork_requiresImport(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_logical_network", "test")
	r := StackHCILogicalNetworkResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dynamic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccStackHCILogicalNetwork_complete(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_logical_network", "test")
	r := StackHCILogicalNetworkResource{}

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

func (r StackHCILogicalNetworkResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	clusterClient := client.AzureStackHCI.LogicalNetworks
	id, err := logicalnetworks.ParseLogicalNetworkID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clusterClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r StackHCILogicalNetworkResource) dynamic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_logical_network" "test" {
  name                = "acctest-ln-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %q
  virtual_switch_name = "ConvergedSwitch(managementcompute)"

  subnet {
    ip_allocation_method = "Dynamic"
  }
}
`, template, os.Getenv(customLocationIdEnv))
}

func (r StackHCILogicalNetworkResource) requiresImport(data acceptance.TestData) string {
	config := r.dynamic(data)

	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_logical_network" "import" {
  name                = azurerm_stack_hci_logical_network.test.name
  resource_group_name = azurerm_stack_hci_logical_network.test.resource_group_name
  location            = azurerm_stack_hci_logical_network.test.location
  custom_location_id  = azurerm_stack_hci_logical_network.test.custom_location_id
  virtual_switch_name = azurerm_stack_hci_logical_network.test.virtual_switch_name

  subnet {
    ip_allocation_method = azurerm_stack_hci_logical_network.test.subnet.0.ip_allocation_method
  }
}
`, config)
}

func (r StackHCILogicalNetworkResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_logical_network" "test" {
  name                = "acctest-ln-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %q
  virtual_switch_name = "ConvergedSwitch(managementcompute)"
  dns_servers         = ["10.0.0.7", "10.0.0.8"]

  subnet {
    ip_allocation_method = "Static"
    address_prefix       = "10.0.0.0/24"
    vlan_id              = 123
    ip_pool {
      start = "10.0.0.218"
      end   = "10.0.0.230"
    }
    ip_pool {
      start = "10.0.0.234"
      end   = "10.0.0.239"
    }
    route {
      address_prefix      = "10.0.0.0/28"
      next_hop_ip_address = "10.0.0.1"
    }
  }
}
`, template, os.Getenv(customLocationIdEnv))
}

func (r StackHCILogicalNetworkResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_logical_network" "test" {
  name                = "acctest-ln-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %q
  virtual_switch_name = "ConvergedSwitch(managementcompute)"
  dns_servers         = ["10.0.0.7", "10.0.0.8"]

  subnet {
    ip_allocation_method = "Static"
    address_prefix       = "10.0.0.0/24"
    vlan_id              = 123
    ip_pool {
      start = "10.0.0.218"
      end   = "10.0.0.230"
    }
    ip_pool {
      start = "10.0.0.234"
      end   = "10.0.0.239"
    }
    route {
      address_prefix      = "10.0.0.0/28"
      next_hop_ip_address = "10.0.0.1"
    }
  }

  tags = {
    foo = "bar"
  }
}
`, template, os.Getenv(customLocationIdEnv))
}

func (r StackHCILogicalNetworkResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_logical_network" "test" {
  name                = "acctest-ln-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %q
  virtual_switch_name = "ConvergedSwitch(managementcompute)"
  dns_servers         = ["10.0.0.7", "10.0.0.8"]

  subnet {
    ip_allocation_method = "Static"
    address_prefix       = "10.0.0.0/24"
    vlan_id              = 123
    ip_pool {
      start = "10.0.0.208"
      end   = "10.0.0.210"
    }
    ip_pool {
      start = "10.0.0.244"
      end   = "10.0.0.249"
    }
    route {
      address_prefix      = "10.0.0.0/28"
      next_hop_ip_address = "10.0.0.1"
    }
  }

  tags = {
    foo = "bar"
    env = "test"
  }
}
`, template, os.Getenv(customLocationIdEnv))
}

func (r StackHCILogicalNetworkResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "primary_location" {
  default = %q
}

variable "random_string" {
  default = %q
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-hci-ln-${var.random_string}"
  location = var.primary_location
}
`, data.Locations.Primary, data.RandomString)
}
