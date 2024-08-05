package azurestackhci_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StackHCINetworkInterfaceResource struct{}

func TestAccStackHCINetworkInterface_dynamic(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_network_interface", "test")
	r := StackHCINetworkInterfaceResource{}

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
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_network_interface", "test")
	r := StackHCINetworkInterfaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.updateNoTag(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateTag(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateTagAgain(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateNoTag(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStackHCINetworkInterface_requiresImport(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

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
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

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
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r StackHCINetworkInterfaceResource) dynamic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_logical_network" "test" {
  name                = "acctest-ln-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  virtual_switch_name = "ConvergedSwitch(managementcompute)"

  subnet {
    ip_allocation_method = "Dynamic"
  }
}

resource "azurerm_stack_hci_network_interface" "test" {
  name                = "acctest-ni-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q

  ip_configuration {
    subnet_id = azurerm_stack_hci_logical_network.test.id
  }
}
`, template, data.RandomString, os.Getenv(customLocationIdEnv))
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
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_logical_network" "test" {
  name                = "acctest-ln-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  virtual_switch_name = "ConvergedSwitch(managementcompute)"
  dns_servers         = ["10.0.0.7"]
  subnet {
    ip_allocation_method = "Static"
    address_prefix       = "10.0.11.0/24"
    vlan_id              = 123

    route {
      name                = "test-route"
      address_prefix      = "0.0.0.0/0"
      next_hop_ip_address = "10.0.0.1"
    }
  }
}

resource "azurerm_stack_hci_network_interface" "test" {
  name                = "acctest-ni-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  dns_servers         = ["10.0.0.8"]
  mac_address         = "02:ec:01:0c:00:07"

  ip_configuration {
    private_ip_address = "10.0.11.%[4]d"
    subnet_id          = azurerm_stack_hci_logical_network.test.id
  }
}
`, template, data.RandomString, os.Getenv(customLocationIdEnv), data.RandomInteger%100)
}

func (r StackHCINetworkInterfaceResource) updateNoTag(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s
provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_logical_network" "test" {
  name                = "acctest-ln-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  virtual_switch_name = "ConvergedSwitch(managementcompute)"
  dns_servers         = ["10.0.0.7"]
  subnet {
    ip_allocation_method = "Static"
    address_prefix       = "10.0.12.0/24"
    vlan_id              = 123

    route {
      name                = "test-route"
      address_prefix      = "0.0.0.0/0"
      next_hop_ip_address = "10.0.20.1"
    }
  }
}

resource "azurerm_stack_hci_network_interface" "test" {
  name                = "acctest-ni-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  dns_servers         = ["10.0.0.8"]
  mac_address         = "02:ec:01:0c:00:08"

  ip_configuration {
    private_ip_address = "10.0.12.%[4]d"
    subnet_id          = azurerm_stack_hci_logical_network.test.id
  }
}
`, template, data.RandomString, os.Getenv(customLocationIdEnv), data.RandomInteger%100)
}

func (r StackHCINetworkInterfaceResource) updateTag(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_logical_network" "test" {
  name                = "acctest-ln-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  virtual_switch_name = "ConvergedSwitch(managementcompute)"
  dns_servers         = ["10.0.0.7"]
  subnet {
    ip_allocation_method = "Static"
    address_prefix       = "10.0.12.0/24"
    vlan_id              = 123

    route {
      name                = "test-route"
      address_prefix      = "0.0.0.0/0"
      next_hop_ip_address = "10.0.20.1"
    }
  }
}

resource "azurerm_stack_hci_network_interface" "test" {
  name                = "acctest-ni-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  dns_servers         = ["10.0.0.8"]
  mac_address         = "02:ec:01:0c:00:08"

  ip_configuration {
    private_ip_address = "10.0.12.%[4]d"
    subnet_id          = azurerm_stack_hci_logical_network.test.id
  }

  tags = {
    foo = "bar"
  }
}
`, template, data.RandomString, os.Getenv(customLocationIdEnv), data.RandomInteger%100)
}

func (r StackHCINetworkInterfaceResource) updateTagAgain(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_logical_network" "test" {
  name                = "acctest-ln-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  virtual_switch_name = "ConvergedSwitch(managementcompute)"
  dns_servers         = ["10.0.0.7"]
  subnet {
    ip_allocation_method = "Static"
    address_prefix       = "10.0.12.0/24"
    vlan_id              = 123

    route {
      name                = "test-route"
      address_prefix      = "0.0.0.0/0"
      next_hop_ip_address = "10.0.20.1"
    }
  }
}

resource "azurerm_stack_hci_network_interface" "test" {
  name                = "acctest-ni-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  dns_servers         = ["10.0.0.8"]
  mac_address         = "02:ec:01:0c:00:08"

  ip_configuration {
    private_ip_address = "10.0.12.%[4]d"
    subnet_id          = azurerm_stack_hci_logical_network.test.id
  }

  tags = {
    foo = "bar"
    env = "test"
  }
}
`, template, data.RandomString, os.Getenv(customLocationIdEnv), data.RandomInteger%100)
}

func (r StackHCINetworkInterfaceResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_logical_network" "test" {
  name                = "acctest-ln-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  virtual_switch_name = "ConvergedSwitch(managementcompute)"
  dns_servers         = ["10.0.0.7"]
  subnet {
    ip_allocation_method = "Static"
    address_prefix       = "10.0.13.0/24"
    vlan_id              = 123

    route {
      name                = "test-route"
      address_prefix      = "0.0.0.0/0"
      next_hop_ip_address = "10.0.20.1"
    }
  }
}

resource "azurerm_stack_hci_network_interface" "test" {
  name                = "acctest-ni-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  dns_servers         = ["10.0.0.8"]
  mac_address         = "02:ec:01:0c:00:09"

  ip_configuration {
    private_ip_address = "10.0.13.%[4]d"
    subnet_id          = azurerm_stack_hci_logical_network.test.id
  }

  tags = {
    foo = "bar"
    env = "test"
  }
}
`, template, data.RandomString, os.Getenv(customLocationIdEnv), data.RandomInteger%100)
}

func (r StackHCINetworkInterfaceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-hci-ni-%[2]s"
  location = "%[1]s"
}
`, data.Locations.Primary, data.RandomString)
}
