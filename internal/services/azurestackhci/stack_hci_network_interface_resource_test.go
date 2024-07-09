package azurestackhci_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StackHCINetworkInterfaceResource struct{}

// https://learn.microsoft.com/en-us/azure-stack/hci/manage/create-logical-networks?tabs=azurecli#prerequisites
// The resource can only be created on the customlocation generated after HCI deployment
const (
	customLocationIdEnv = "ARM_TEST_STACK_HCI_CUSTOM_LOCATION_ID"
)

func TestAccStackHCINetworkInterface_dynamic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_network_interface", "test")
	r := StackHCINetworkInterfaceResource{}

	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skip the test as one or more of below environment variables are not specified: %q", customLocationIdEnv)
	}

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

func TestAccStackHCINetworkInterface_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_network_interface", "test")
	r := StackHCINetworkInterfaceResource{}

	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skip the test as one or more of below environment variables are not specified: %q", customLocationIdEnv)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
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

func TestAccStackHCINetworkInterface_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_network_interface", "test")
	r := StackHCINetworkInterfaceResource{}

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

func TestAccStackHCINetworkInterface_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_network_interface", "test")
	r := StackHCINetworkInterfaceResource{}

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

func (r StackHCINetworkInterfaceResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	clusterClient := client.AzureStackHCI.NetworkInterfaces
	id, err := networkinterfaces.ParseNetworkInterfaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clusterClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r StackHCINetworkInterfaceResource) dynamic(data acceptance.TestData) string {
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

resource "azurerm_stack_hci_network_interface" "test" {
  name                = "acctest-ni-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %q

  ip_configuration {
    subnet_id = azurerm_stack_hci_logical_network.test.id
  }
}
`, template, os.Getenv(customLocationIdEnv))
}

func (r StackHCINetworkInterfaceResource) requiresImport(data acceptance.TestData) string {
	config := r.dynamic(data)

	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_network_interface" "import" {
  name                = azurerm_stack_hci_network_interface.test.name
  resource_group_name = azurerm_stack_hci_network_interface.test.resource_group_name
  location            = azurerm_stack_hci_network_interface.test.location
  custom_location_id  = azurerm_stack_hci_network_interface.test.custom_location_id

  ip_configuration {
    subnet_id = azurerm_stack_hci_network_interface.test.ip_configuration.0.subnet_id
  }
}
`, config)
}

func (r StackHCINetworkInterfaceResource) basic(data acceptance.TestData) string {
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
  dns_servers         = ["10.0.0.7"]
  subnet {
    ip_allocation_method = "Static"
    address_prefix       = "10.0.0.0/24"
    vlan_id              = 123

    route {
      name                = "test-route"
      address_prefix      = "0.0.0.0/0"
      next_hop_ip_address = "10.0.20.1"
    }
  }
}

resource "azurerm_stack_hci_network_interface" "test" {
  name                = "acctest-ni-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %q
  dns_servers         = ["10.0.0.8"]
  mac_address         = "02:ec:01:0c:00:09"

  ip_configuration {
    private_ip_address = "10.0.0.2"
    subnet_id          = azurerm_stack_hci_logical_network.test.id
  }
}
`, template, os.Getenv(customLocationIdEnv))
}

func (r StackHCINetworkInterfaceResource) update(data acceptance.TestData) string {
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
  custom_location_id  = %[1]q
  virtual_switch_name = "ConvergedSwitch(managementcompute)"
  dns_servers         = ["10.0.0.7"]
  subnet {
    ip_allocation_method = "Static"
    address_prefix       = "10.0.0.0/24"
    vlan_id              = 123

    route {
      name                = "test-route"
      address_prefix      = "0.0.0.0/0"
      next_hop_ip_address = "10.0.20.1"
    }
  }
}

resource "azurerm_stack_hci_network_interface" "test" {
  name                = "acctest-ni-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[1]q
  dns_servers         = ["10.0.0.8"]
  mac_address         = "02:ec:01:0c:00:09"

  ip_configuration {
    private_ip_address = "10.0.0.2"
    subnet_id          = azurerm_stack_hci_logical_network.test.id
  }

  tags = {
    foo = "bar"
  }
}
`, template, os.Getenv(customLocationIdEnv))
}

func (r StackHCINetworkInterfaceResource) complete(data acceptance.TestData) string {
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
  dns_servers         = ["10.0.0.7"]
  subnet {
    ip_allocation_method = "Static"
    address_prefix       = "10.0.0.0/24"
    vlan_id              = 123

    route {
      name                = "test-route"
      address_prefix      = "0.0.0.0/0"
      next_hop_ip_address = "10.0.20.1"
    }
  }
}

resource "azurerm_stack_hci_network_interface" "test" {
  name                = "acctest-ni-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[1]q
  dns_servers         = ["10.0.0.8"]
  mac_address         = "02:ec:01:0c:00:09"

  ip_configuration {
    private_ip_address = "10.0.0.2"
    subnet_id          = azurerm_stack_hci_logical_network.test.id
  }

  tags = {
    foo = "bar"
    env = "test"
  }
}
`, template, os.Getenv(customLocationIdEnv))
}

func (r StackHCINetworkInterfaceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "primary_location" {
  default = %q
}

variable "random_string" {
  default = %q
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-hci-ni-${var.random_string}"
  location = var.primary_location
}
`, data.Locations.Primary, data.RandomString)
}
