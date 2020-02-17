package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMLinuxVirtualMachineScaleSet_otherBootDiagnostics(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				// Enabled
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_otherBootDiagnostics(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// Removed
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_otherBootDiagnosticsDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// Enabled
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_otherBootDiagnostics(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_otherComputerNamePrefix(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_otherComputerNamePrefix(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_otherCustomData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_otherCustomData(data, "/bin/bash"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
				"custom_data",
			),
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_otherCustomData(data, "/bin/zsh"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
				"custom_data",
			),
			{
				// removed
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_authPassword(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
				"custom_data",
			),
		},
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_otherDoNotRunExtensionsOnOverProvisionedMachines(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_otherDoNotRunExtensionsOnOverProvisionedMachines(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_otherPrioritySpotDeallocate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_otherPrioritySpot(data, "Deallocate"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_otherPrioritySpotDelete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_otherPrioritySpot(data, "Delete"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_otherPrioritySpotMaxBidPrice(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				// expensive, but guarantees this test will pass
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_otherPrioritySpotMaxBidPrice(data, "0.5000"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_otherPrioritySpotMaxBidPrice(data, "-1"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_otherPriorityRegular(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_otherPriorityRegular(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_otherRequiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_authPassword(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMLinuxVirtualMachineScaleSet_otherRequiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_linux_virtual_machine_scale_set"),
			},
		},
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_otherSecret(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_otherSecret(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// update
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_otherSecretUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),

			{
				// removed
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_otherSecretRemoved(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_otherTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_otherTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// add one
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_otherTagsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// remove all
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_authPassword(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_otherVMAgent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_otherVMAgent(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_otherVMAgentDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_otherVMAgent(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_updateHealthProbe(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_updateLoadBalancerHealthProbeSKUBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_updateLoadBalancerHealthProbeSKUStandard(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
		},
	})
}

func testAccAzureRMLinuxVirtualMachineScaleSet_otherBootDiagnostics(data acceptance.TestData) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "accsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  boot_diagnostics {
    storage_account_uri = azurerm_storage_account.test.primary_blob_endpoint
  }
}
`, template, data.RandomString, data.RandomInteger)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_otherBootDiagnosticsDisabled(data acceptance.TestData) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "accsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, template, data.RandomString, data.RandomInteger)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_otherComputerNamePrefix(data acceptance.TestData) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                 = "acctestvmss-%d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  sku                  = "Standard_F2"
  instances            = 1
  admin_username       = "adminuser"
  admin_password       = "P@ssword1234!"
  computer_name_prefix = "morty"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_otherCustomData(data acceptance.TestData, customData string) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  custom_data         = base64encode(%q)

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, template, data.RandomInteger, customData)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_otherDoNotRunExtensionsOnOverProvisionedMachines(data acceptance.TestData) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  overprovision       = true

  disable_password_authentication                   = false
  do_not_run_extensions_on_overprovisioned_machines = true

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_otherPrioritySpot(data acceptance.TestData, evictionPolicy string) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  eviction_policy     = %q
  priority            = "Spot"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, template, data.RandomInteger, evictionPolicy)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_otherPrioritySpotMaxBidPrice(data acceptance.TestData, maxBid string) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  eviction_policy     = "Delete"
  priority            = "Spot"
  max_bid_price       = %q

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, template, data.RandomInteger, maxBid)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_otherPriorityRegular(data acceptance.TestData) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  priority            = "Regular"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_otherRequiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_authPassword(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "import" {
  name                = azurerm_linux_virtual_machine_scale_set.test.name
  resource_group_name = azurerm_linux_virtual_machine_scale_set.test.resource_group_name
  location            = azurerm_linux_virtual_machine_scale_set.test.location
  sku                 = azurerm_linux_virtual_machine_scale_set.test.sku
  instances           = azurerm_linux_virtual_machine_scale_set.test.instances
  admin_username      = azurerm_linux_virtual_machine_scale_set.test.admin_username
  admin_password      = azurerm_linux_virtual_machine_scale_set.test.admin_password

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, template)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_otherSecret(data acceptance.TestData) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_otherSecretTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  secret {
    key_vault_id = azurerm_key_vault.test.id

    certificate {
      url = azurerm_key_vault_certificate.first.secret_id
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_otherSecretRemoved(data acceptance.TestData) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_otherSecretTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_otherSecretUpdated(data acceptance.TestData) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_otherSecretTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  secret {
    key_vault_id = azurerm_key_vault.test.id

    certificate {
      url = azurerm_key_vault_certificate.first.secret_id
    }

    certificate {
      url = azurerm_key_vault_certificate.second.secret_id
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_otherSecretTemplate(data acceptance.TestData) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name                        = "standard"
  enabled_for_template_deployment = true
  enabled_for_deployment          = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.service_principal_object_id

    certificate_permissions = [
      "create",
      "delete",
      "get",
      "update",
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]

    storage_permissions = [
      "set",
    ]
  }
}

resource "azurerm_key_vault_certificate" "first" {
  name         = "first"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=hello-world-first"
      validity_in_months = 12
    }
  }
}

resource "azurerm_key_vault_certificate" "second" {
  name         = "second"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=hello-world-second"
      validity_in_months = 12
    }
  }
}
`, template, data.RandomString)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_otherTags(data acceptance.TestData) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  tags = {
    artist = "Billy"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_otherTagsUpdated(data acceptance.TestData) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  tags = {
    artist = "Billy"
    when   = "we all fall asleep"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_otherVMAgent(data acceptance.TestData, enabled bool) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false
  provision_vm_agent              = %t

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, template, data.RandomInteger, enabled)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_updateLoadBalancerHealthProbeSKUBasic(data acceptance.TestData) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_public_ip" "test" {
  name                    = "acctestpip-%[2]d"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  allocation_method       = "Dynamic"
  idle_timeout_in_minutes = 4
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[2]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Basic"

  frontend_ip_configuration {
    name                 = "internal"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "test"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
}

resource "azurerm_lb_nat_pool" "test" {
  name                           = "test"
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  frontend_ip_configuration_name = "internal"
  protocol                       = "Tcp"
  frontend_port_start            = 80
  frontend_port_end              = 81
  backend_port                   = 8080
}

resource "azurerm_lb_probe" "test" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "acctest-lb-probe"
  port                = 22
  protocol            = "Tcp"
}

resource "azurerm_lb_rule" "test" {
  name                           = "AccTestLBRule"
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  probe_id                       = azurerm_lb_probe.test.id
  backend_address_pool_id        = azurerm_lb_backend_address_pool.test.id
  frontend_ip_configuration_name = "internal"
  protocol                       = "Tcp"
  frontend_port                  = 22
  backend_port                   = 22
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  health_probe_id     = azurerm_lb_probe.test.id

  depends_on = ["azurerm_lb_rule.test"]

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
    disk_size_gb         = 10
    lun                  = 10
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name                                   = "internal"
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
      load_balancer_inbound_nat_rules_ids    = [azurerm_lb_nat_pool.test.id]
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_updateLoadBalancerHealthProbeSKUStandard(data acceptance.TestData) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_public_ip" "test" {
  name                    = "acctestpip-%[2]d"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  allocation_method       = "Static"
  idle_timeout_in_minutes = 4
  sku                     = "Standard"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[2]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "internal"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "test"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
}

resource "azurerm_lb_nat_pool" "test" {
  name                           = "test"
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  frontend_ip_configuration_name = "internal"
  protocol                       = "Tcp"
  frontend_port_start            = 80
  frontend_port_end              = 81
  backend_port                   = 8080
}

resource "azurerm_lb_probe" "test" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "acctest-lb-probe"
  port                = 22
  protocol            = "Tcp"
}

resource "azurerm_lb_rule" "test" {
  name                           = "AccTestLBRule"
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  probe_id                       = azurerm_lb_probe.test.id
  backend_address_pool_id        = azurerm_lb_backend_address_pool.test.id
  frontend_ip_configuration_name = "internal"
  protocol                       = "Tcp"
  frontend_port                  = 22
  backend_port                   = 22
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  health_probe_id     = azurerm_lb_probe.test.id

  depends_on = ["azurerm_lb_rule.test"]

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
    disk_size_gb         = 10
    lun                  = 10
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name                                   = "internal"
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
      load_balancer_inbound_nat_rules_ids    = [azurerm_lb_nat_pool.test.id]
    }
  }
}
`, template, data.RandomInteger)
}
