package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccOrchestratedVirtualMachineScaleSet_disksOSDiskCaching(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksOSDiskEphemeral(data, "CacheDisk"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.disksOSDiskEphemeral(data, "ResourceDisk"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func (r OrchestratedVirtualMachineScaleSetResource) disksOSDiskEphemeral(data acceptance.TestData, placement string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[3]d"
  location = "%[2]s"
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name  = "Standard_F4s_v2"
  instances = 1

  platform_fault_domain_count = 2

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm-%[3]d"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadOnly"

    diff_disk_settings {
      option    = "Local"
      placement = "%[4]s"
    }
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }
}
`, r.natgateway_template(data), data.Locations.Primary, data.RandomInteger, placement)
}
