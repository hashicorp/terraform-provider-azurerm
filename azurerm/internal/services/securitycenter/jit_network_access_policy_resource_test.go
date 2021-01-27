package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type JitNetworkAccessPolicyResource struct {
}

func TestAccJitNetworkAccessPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_jit_network_access_policy", "test")
	r := JitNetworkAccessPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccJitNetworkAccessPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_jit_network_access_policy", "test")
	r := JitNetworkAccessPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccJitNetworkAccessPolicy_multiPorts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_jit_network_access_policy", "test")
	r := JitNetworkAccessPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multiPorts(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r JitNetworkAccessPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.JitNetworkAccessPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.SecurityCenter.NewJitNetworkAccessPoliciesClient(id.LocationName)
	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.JitNetworkAccessPolicyProperties != nil), nil
}

func (r JitNetworkAccessPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_jit_network_access_policy" "test" {
  name                = "policy%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  virtual_machines {
    id = azurerm_windows_virtual_machine.test.id
    ports {
      port                            = 22
      protocol                        = "*"
      allowed_source_address_prefixes = ["*"]
      max_request_access_duration     = "PT3H"
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r JitNetworkAccessPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_jit_network_access_policy" "import" {
  name                = azurerm_jit_network_access_policy.test.name
  location            = azurerm_jit_network_access_policy.test.location
  resource_group_name = azurerm_jit_network_access_policy.test.resource_group_name

  dynamic "virtual_machines" {
    for_each = azurerm_jit_network_access_policy.test.virtual_machines
    content {
      id = virtual_machines.value.id
      dynamic "ports" {
        for_each = lookup(virtual_machines.value, "ports", [])
        content {
          port                            = ports.value.port
          protocol                        = ports.value.protocol
          allowed_source_address_prefixes = ports.value.allowed_source_address_prefixes
          max_request_access_duration     = ports.value.max_request_access_duration
        }
      }
    }
  }
}
`, r.basic(data))
}

func (r JitNetworkAccessPolicyResource) multiPorts(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_jit_network_access_policy" "test" {
  name                = "policy%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  virtual_machines {
    id = azurerm_windows_virtual_machine.test.id
    ports {
      port                            = 22
      protocol                        = "*"
      allowed_source_address_prefixes = ["*"]
      max_request_access_duration     = "PT3H"
    }

    ports {
      port                            = 443
      protocol                        = "TCP"
      allowed_source_address_prefixes = ["192.168.0.3", "192.168.0.0/16"]
      max_request_access_duration     = "PT5M"
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r JitNetworkAccessPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctest_vm_ip"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Dynamic"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_interface_security_group_association" "test" {
  network_interface_id      = azurerm_network_interface.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_windows_virtual_machine" "test" {
  name                = "acctestvm%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  depends_on = [azurerm_network_interface_security_group_association.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString)
}
