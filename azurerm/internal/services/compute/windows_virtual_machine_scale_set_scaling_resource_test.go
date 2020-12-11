package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMWindowsVirtualMachineScaleSet_scalingAutoScale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_scalingAutoScale(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMWindowsVirtualMachineScaleSet_scalingInstanceCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_scalingInstanceCount(data, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_scalingInstanceCount(data, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_scalingInstanceCount(data, 5),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// update the count but the `sku` should be ignored
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_scalingInstanceCountIgnoreUpdatedSku(data, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// confirm that the `sku` hasn't been changed
				Config:   testAccAzureRMWindowsVirtualMachineScaleSet_scalingInstanceCount(data, 3),
				PlanOnly: true,
			},
		},
	})
}

func TestAccAzureRMWindowsVirtualMachineScaleSet_scalingOverProvisionDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_scalingOverProvisionDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMWindowsVirtualMachineScaleSet_scalingProximityPlacementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_scalingProximityPlacementGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMWindowsVirtualMachineScaleSet_scalingSinglePlacementGroupDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_scalingSinglePlacementGroupDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMWindowsVirtualMachineScaleSet_scalingSinglePlacementGroupDisabledUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_authPassword(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_scalingSinglePlacementGroupDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMWindowsVirtualMachineScaleSet_scalingUpdateSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_scalingUpdateSku(data, "Standard_F2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_scalingUpdateSku(data, "Standard_F4"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_scalingUpdateSku(data, "Standard_F2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// confirms that the `instances` count comes from the API
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_scalingUpdateSkuIgnoredUpdatedCount(data, "Standard_F2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			{
				Config:   testAccAzureRMWindowsVirtualMachineScaleSet_scalingUpdateSku(data, "Standard_F2"),
				PlanOnly: true,
			},
		},
	})
}

func TestAccAzureRMWindowsVirtualMachineScaleSet_scalingZonesSingle(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_scalingZonesSingle(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMWindowsVirtualMachineScaleSet_scalingZonesMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_scalingZonesMultiple(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMWindowsVirtualMachineScaleSet_scalingZonesBalance(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_scalingZonesBalance(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func testAccAzureRMWindowsVirtualMachineScaleSet_scalingAutoScale(data acceptance.TestData) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 2
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  overprovision       = true

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
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

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "autoscale-config"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  target_resource_id  = azurerm_windows_virtual_machine_scale_set.test.id

  profile {
    name = "AutoScale"

    capacity {
      default = 2
      minimum = 1
      maximum = 5
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_windows_virtual_machine_scale_set.test.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = "1"
        cooldown  = "PT1M"
      }
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_windows_virtual_machine_scale_set.test.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "LessThan"
        threshold          = 25
      }

      scale_action {
        direction = "Decrease"
        type      = "ChangeCount"
        value     = "1"
        cooldown  = "PT1M"
      }
    }
  }
}
`, template)
}

func testAccAzureRMWindowsVirtualMachineScaleSet_scalingInstanceCount(data acceptance.TestData, instanceCount int) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = %d
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
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
`, template, instanceCount)
}

func testAccAzureRMWindowsVirtualMachineScaleSet_scalingInstanceCountIgnoreUpdatedSku(data acceptance.TestData, instanceCount int) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F4"
  instances           = %d
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
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

  lifecycle {
    ignore_changes = ["sku"]
  }
}
`, template, instanceCount)
}

func testAccAzureRMWindowsVirtualMachineScaleSet_scalingOverProvisionDisabled(data acceptance.TestData) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 3
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  overprovision       = false

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
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

func testAccAzureRMWindowsVirtualMachineScaleSet_scalingProximityPlacementGroup(data acceptance.TestData) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_proximity_placement_group" "test" {
  name                = "acctestPPG-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                         = local.vm_name
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  sku                          = "Standard_F2"
  instances                    = 1
  admin_username               = "adminuser"
  admin_password               = "P@ssword1234!"
  proximity_placement_group_id = azurerm_proximity_placement_group.test.id

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
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

func testAccAzureRMWindowsVirtualMachineScaleSet_scalingSinglePlacementGroupDisabled(data acceptance.TestData) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                   = local.vm_name
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  sku                    = "Standard_F2"
  instances              = 1
  admin_username         = "adminuser"
  admin_password         = "P@ssword1234!"
  single_placement_group = false

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
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

func testAccAzureRMWindowsVirtualMachineScaleSet_scalingUpdateSku(data acceptance.TestData, skuName string) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = %q
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
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
`, template, skuName)
}

func testAccAzureRMWindowsVirtualMachineScaleSet_scalingUpdateSkuIgnoredUpdatedCount(data acceptance.TestData, skuName string) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = %q
  instances           = 5
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
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

  lifecycle {
    ignore_changes = ["instances"]
  }
}
`, template, skuName)
}

func testAccAzureRMWindowsVirtualMachineScaleSet_scalingZonesSingle(data acceptance.TestData) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  zones               = ["1"]

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
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

func testAccAzureRMWindowsVirtualMachineScaleSet_scalingZonesMultiple(data acceptance.TestData) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 2
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  zones               = ["1", "2"]

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
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

func testAccAzureRMWindowsVirtualMachineScaleSet_scalingZonesBalance(data acceptance.TestData) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 4
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  zones               = ["1", "2"]
  zone_balance        = true

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
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
