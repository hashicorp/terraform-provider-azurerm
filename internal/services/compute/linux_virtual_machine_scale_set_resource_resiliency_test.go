// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccLinuxVirtualMachineScaleSet_resiliency_Basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.resiliencyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// Basic test includes VM policies and zones (no health probe)
				// Azure may automatically enable zone rebalancing when zones are present
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("false"),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_resiliency_Minimal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.resiliencyMinimal(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// Minimal test: resiliency with all policies disabled
				// When all policies are disabled, resiliency block should be removed from state ("None" pattern)
				check.That(data.ResourceName).Key("resiliency.#").HasValue("0"),
				check.That(data.ResourceName).Key("zones.#").HasValue("0"),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_resiliency_Complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.resiliencyComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resiliency.0.automatic_zone_rebalancing.0.rebalance_behavior").HasValue("CreateBeforeDelete"),
				check.That(data.ResourceName).Key("resiliency.0.automatic_zone_rebalancing.0.rebalance_strategy").HasValue("Recreate"),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("true"),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_resiliency_Update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Start with VM policies only (no zones)
			Config: r.resiliencyVMPoliciesOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("true"),
			),
		},
		data.ImportStep("admin_password"),
		{
			// Add zones and health probe with zone rebalancing (ForceNew because zones added)
			Config: r.resiliencyComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resiliency.0.automatic_zone_rebalancing.0.rebalance_behavior").HasValue("CreateBeforeDelete"),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("true"),
			),
		},
		data.ImportStep("admin_password"),
		{
			// Update resiliency settings (remove zone rebalancing, keep VM policies and zones)
			Config: r.resiliencyVMPoliciesOnlyMultipleZonesWithHealthProbe(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("true"),
				// Following "None" pattern - automatic_zone_rebalancing removed from state when disabled
				check.That(data.ResourceName).Key("resiliency.#").HasValue("0"),
			),
		},
		data.ImportStep("admin_password"),
		{
			// Finally test minimal config (disable VM policies but keep zones and health probe)
			Config: r.resiliencyMinimalWithHealthProbeConsistent(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("false"),
			),
		},
		data.ImportStep("admin_password"),
		{
			// Scale down to 0 instances to avoid post-test cleanup issues
			Config: r.resiliencyMinimalWithHealthProbeConsistentScaledDown(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("instances").HasValue("0"),
			),
		},
	})
}

func TestAccLinuxVirtualMachineScaleSet_resiliency_RequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.resiliencyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.resiliencyRequiresImport),
	})
}

func TestAccLinuxVirtualMachineScaleSet_resiliency_Removal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Start with VM policies only (no zones to avoid ForceNew issues)
			Config: r.resiliencyVMPoliciesOnlyWithoutRollingUpgrade(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("true"),
			),
		},
		data.ImportStep("admin_password"),
		{
			// Remove resiliency block entirely (no zones change)
			Config: r.basicWithoutRollingUpgrade(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		{
			// Re-add resiliency block (VM policies only, no zone rebalancing since VMSS has no health probe)
			Config: r.resiliencyVMPoliciesOnlyWithoutRollingUpgrade(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("true"),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_resiliency_ZoneRebalancingOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.resiliencyZoneRebalancingOnlyWithHealthProbe(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resiliency.0.automatic_zone_rebalancing.0.rebalance_behavior").HasValue("CreateBeforeDelete"),
				check.That(data.ResourceName).Key("resiliency.0.automatic_zone_rebalancing.0.rebalance_strategy").HasValue("Recreate"),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("false"),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_resiliency_VMPoliciesOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.resiliencyVMPoliciesOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("true"),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_resiliency_ZoneRebalancingNonePattern(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Start with VM policies only (NO automatic_zone_rebalancing block) but with multiple zones
			Config: r.resiliencyVMPoliciesOnlyMultipleZones(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("true"),
				// Following "None" pattern - field omitted when automatic zone rebalancing not explicitly enabled
				check.That(data.ResourceName).Key("resiliency.0.automatic_zone_rebalancing.#").HasValue("0"),
			),
		},
		data.ImportStep("admin_password"),
		{
			// Add automatic_zone_rebalancing block
			Config: r.resiliencyComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resiliency.0.automatic_zone_rebalancing.0.rebalance_behavior").HasValue("CreateBeforeDelete"),
				check.That(data.ResourceName).Key("resiliency.0.automatic_zone_rebalancing.0.rebalance_strategy").HasValue("Recreate"),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("true"),
			),
		},
		data.ImportStep("admin_password"),
		{
			// Remove automatic_zone_rebalancing block again (back to VM policies only with zones)
			Config: r.resiliencyVMPoliciesOnlyMultipleZones(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("true"),
				// Following "None" pattern - field omitted when automatic zone rebalancing not explicitly enabled
				check.That(data.ResourceName).Key("resiliency.0.automatic_zone_rebalancing.#").HasValue("0"),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_resiliency_ValidationRebalanceBehavior(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.resiliencyInvalidRebalanceBehavior(data),
			ExpectError: regexp.MustCompile(`expected resiliency\.0\.automatic_zone_rebalancing\.0\.rebalance_behavior to be one of \["CreateBeforeDelete"\], got InvalidBehavior`),
		},
	})
}

func TestAccLinuxVirtualMachineScaleSet_resiliency_ValidationRebalanceStrategy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.resiliencyInvalidRebalanceStrategy(data),
			ExpectError: regexp.MustCompile(`expected resiliency\.0\.automatic_zone_rebalancing\.0\.rebalance_strategy to be one of \["Recreate"\], got InvalidStrategy`),
		},
	})
}

func TestAccLinuxVirtualMachineScaleSet_resiliency_customizeDiffResiliencyZoneValidation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.resiliencyZoneRebalancingSingleZone(data),
			ExpectError: regexp.MustCompile("`automatic_zone_rebalancing` can only be configured when the Virtual Machine Scale Set is deployed across multiple `zones`"),
		},
		{
			Config:      r.resiliencyZoneRebalancingNoZones(data),
			ExpectError: regexp.MustCompile("`automatic_zone_rebalancing` can only be configured when the Virtual Machine Scale Set is deployed across multiple `zones`"),
		},
		{
			Config: r.resiliencyZoneRebalancingMultipleZones(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_resiliency_CustomizeDiffValidation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Test zone rebalancing without health probe should fail
			Config:      r.resiliencyZoneRebalancingWithoutHealthProbe(data),
			ExpectError: regexp.MustCompile("`automatic_zone_rebalancing` requires a `health_probe_id` to be configured. Azure requires health probes to be applied to all instances when automatic zone rebalancing is `enabled`"),
		},
		{
			// Test zone rebalancing with single zone should fail
			Config:      r.resiliencyZoneRebalancingSingleZone(data),
			ExpectError: regexp.MustCompile("`automatic_zone_rebalancing` can only be configured when the Virtual Machine Scale Set is deployed across multiple `zones`"),
		},
		{
			// Test zone rebalancing with no zones should fail
			Config:      r.resiliencyZoneRebalancingNoZones(data),
			ExpectError: regexp.MustCompile("`automatic_zone_rebalancing` can only be configured when the Virtual Machine Scale Set is deployed across multiple `zones`"),
		},
		{
			// Test valid zone rebalancing with health probe should succeed
			Config: r.resiliencyZoneRebalancingMultipleZones(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resiliency.0.automatic_zone_rebalancing.0.rebalance_behavior").HasValue("CreateBeforeDelete"),
			),
		},
		data.ImportStep("admin_password"),
	})
}

// Test configurations

// loadBalancerInfrastructureTemplate generates the common load balancer infrastructure
func (LinuxVirtualMachineScaleSetResource) loadBalancerInfrastructureTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "BackEndAddressPool"
}

resource "azurerm_lb_probe" "test" {
  loadbalancer_id     = azurerm_lb.test.id
  name                = "ssh-running-probe"
  port                = 22
  protocol            = "Tcp"
  interval_in_seconds = 15
  number_of_probes    = 2
}

resource "azurerm_lb_rule" "test" {
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "AccTestLBRule"
  protocol                       = "Tcp"
  frontend_port                  = 22
  backend_port                   = 22
  frontend_ip_configuration_name = "PublicIPAddress"
  probe_id                       = azurerm_lb_probe.test.id
  backend_address_pool_ids       = [azurerm_lb_backend_address_pool.test.id]
}`, data.RandomInteger, data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) resiliencyBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestVMSS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  zones               = ["1", "2", "3"]

  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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

  resilient_vm_creation_enabled = true
  resilient_vm_deletion_enabled = false
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) resiliencyComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  depends_on = [azurerm_lb_rule.test]

  name                = "acctestVMSS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 3
  zones               = ["1", "2", "3"]
  health_probe_id     = azurerm_lb_probe.test.id
  upgrade_mode        = "Rolling"
  overprovision       = false

  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false

  rolling_upgrade_policy {
    cross_zone_upgrades_enabled             = true
    max_batch_instance_percent              = 20
    max_unhealthy_instance_percent          = 20
    max_unhealthy_upgraded_instance_percent = 20
    pause_time_between_batches              = "PT0S"
    prioritize_unhealthy_instances_enabled  = false
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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
      name                                   = "internal"
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
    }
  }

  resiliency {
    automatic_zone_rebalancing {
      rebalance_behavior = "CreateBeforeDelete"
      rebalance_strategy = "Recreate"
    }
  }

  resilient_vm_creation_enabled = true
  resilient_vm_deletion_enabled = true
}
`, r.template(data), r.loadBalancerInfrastructureTemplate(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) resiliencyZoneRebalancingOnlyWithHealthProbe(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = "internal"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "test"
  loadbalancer_id = azurerm_lb.test.id
}

resource "azurerm_lb_probe" "test" {
  loadbalancer_id     = azurerm_lb.test.id
  name                = "acctest-lb-probe"
  port                = 22
  protocol            = "Tcp"
  interval_in_seconds = 15
  number_of_probes    = 2
}

resource "azurerm_lb_rule" "test" {
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "AccTestLBRule"
  protocol                       = "Tcp"
  frontend_port                  = 22
  backend_port                   = 22
  frontend_ip_configuration_name = "internal"
  backend_address_pool_ids       = [azurerm_lb_backend_address_pool.test.id]
  probe_id                       = azurerm_lb_probe.test.id
  enable_tcp_reset               = false
  enable_floating_ip             = false
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  depends_on = [azurerm_lb_rule.test]

  name                = "acctestVMSS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 3
  zones               = ["1", "2", "3"]
  health_probe_id     = azurerm_lb_probe.test.id
  upgrade_mode        = "Rolling"
  overprovision       = false

  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false

  rolling_upgrade_policy {
    cross_zone_upgrades_enabled             = true
    max_batch_instance_percent              = 20
    max_unhealthy_instance_percent          = 20
    max_unhealthy_upgraded_instance_percent = 20
    pause_time_between_batches              = "PT0S"
    prioritize_unhealthy_instances_enabled  = false
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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
      name                                   = "internal"
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
    }
  }

  resiliency {
    automatic_zone_rebalancing {
      rebalance_behavior = "CreateBeforeDelete"
      rebalance_strategy = "Recreate"
    }
  }

  resilient_vm_creation_enabled = false
  resilient_vm_deletion_enabled = false
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) resiliencyMinimal(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestVMSS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1

  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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

  resilient_vm_creation_enabled = false
  resilient_vm_deletion_enabled = false
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) resiliencyRequiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "import" {
  name                = azurerm_linux_virtual_machine_scale_set.test.name
  resource_group_name = azurerm_linux_virtual_machine_scale_set.test.resource_group_name
  location            = azurerm_linux_virtual_machine_scale_set.test.location
  sku                 = azurerm_linux_virtual_machine_scale_set.test.sku
  instances           = azurerm_linux_virtual_machine_scale_set.test.instances

  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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

  resilient_vm_creation_enabled = true
  resilient_vm_deletion_enabled = false
}
`, r.resiliencyBasic(data))
}

func (r LinuxVirtualMachineScaleSetResource) resiliencyVMPoliciesOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  depends_on = [azurerm_lb_rule.test]

  name                = "acctestVMSS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  health_probe_id     = azurerm_lb_probe.test.id
  upgrade_mode        = "Rolling"

  rolling_upgrade_policy {
    max_batch_instance_percent              = 20
    max_unhealthy_instance_percent          = 20
    max_unhealthy_upgraded_instance_percent = 20
    pause_time_between_batches              = "PT0S"
    prioritize_unhealthy_instances_enabled  = false
  }

  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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
      name                                   = "internal"
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
    }
  }

  resilient_vm_creation_enabled = true
  resilient_vm_deletion_enabled = true
}
`, r.template(data), r.loadBalancerInfrastructureTemplate(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) resiliencyVMPoliciesOnlyMultipleZones(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  depends_on = [azurerm_lb_rule.test]

  name                = "acctestVMSS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  zones               = ["1", "2", "3"]
  health_probe_id     = azurerm_lb_probe.test.id
  upgrade_mode        = "Rolling"

  rolling_upgrade_policy {
    max_batch_instance_percent              = 20
    max_unhealthy_instance_percent          = 20
    max_unhealthy_upgraded_instance_percent = 20
    pause_time_between_batches              = "PT0S"
    prioritize_unhealthy_instances_enabled  = false
    cross_zone_upgrades_enabled             = true
  }

  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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
      name                                   = "internal"
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
    }
  }

  resilient_vm_creation_enabled = true
  resilient_vm_deletion_enabled = true
}
`, r.template(data), r.loadBalancerInfrastructureTemplate(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) resiliencyInvalidRebalanceBehavior(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestVMSS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1

  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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

  resilient_vm_creation_enabled = true
  resilient_vm_deletion_enabled = false

  resiliency {
    automatic_zone_rebalancing {
      rebalance_behavior = "InvalidBehavior"
      rebalance_strategy = "Recreate"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) resiliencyInvalidRebalanceStrategy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestVMSS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1

  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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

  resilient_vm_creation_enabled = true
  resilient_vm_deletion_enabled = false

  resiliency {
    automatic_zone_rebalancing {
      rebalance_behavior = "CreateBeforeDelete"
      rebalance_strategy = "InvalidStrategy"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) resiliencyZoneRebalancingSingleZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestVMSS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  zones               = ["1"]

  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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

  resilient_vm_creation_enabled = false
  resilient_vm_deletion_enabled = false

  resiliency {
    automatic_zone_rebalancing {
      rebalance_behavior = "CreateBeforeDelete"
      rebalance_strategy = "Recreate"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) resiliencyZoneRebalancingNoZones(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestVMSS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  # No zones specified - should trigger validation error

  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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

  resilient_vm_creation_enabled = false
  resilient_vm_deletion_enabled = false

  resiliency {
    automatic_zone_rebalancing {
      rebalance_behavior = "CreateBeforeDelete"
      rebalance_strategy = "Recreate"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) resiliencyZoneRebalancingMultipleZones(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = "internal"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "test"
  loadbalancer_id = azurerm_lb.test.id
}

resource "azurerm_lb_probe" "test" {
  loadbalancer_id     = azurerm_lb.test.id
  name                = "acctest-lb-probe"
  port                = 22
  protocol            = "Tcp"
  interval_in_seconds = 15
  number_of_probes    = 2
}

resource "azurerm_lb_rule" "test" {
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "AccTestLBRule"
  protocol                       = "Tcp"
  frontend_port                  = 22
  backend_port                   = 22
  frontend_ip_configuration_name = "internal"
  backend_address_pool_ids       = [azurerm_lb_backend_address_pool.test.id]
  probe_id                       = azurerm_lb_probe.test.id
  enable_tcp_reset               = false
  enable_floating_ip             = false
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  depends_on = [azurerm_lb_rule.test]

  name                = "acctestVMSS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 3
  zones               = ["1", "2", "3"]
  health_probe_id     = azurerm_lb_probe.test.id
  upgrade_mode        = "Rolling"
  overprovision       = false

  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false

  rolling_upgrade_policy {
    cross_zone_upgrades_enabled             = true
    max_batch_instance_percent              = 20
    max_unhealthy_instance_percent          = 20
    max_unhealthy_upgraded_instance_percent = 20
    pause_time_between_batches              = "PT0S"
    prioritize_unhealthy_instances_enabled  = false
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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
      name                                   = "internal"
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
    }
  }

  resilient_vm_creation_enabled = false
  resilient_vm_deletion_enabled = false

  resiliency {
    automatic_zone_rebalancing {
      rebalance_behavior = "CreateBeforeDelete"
      rebalance_strategy = "Recreate"
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) resiliencyZoneRebalancingWithoutHealthProbe(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestVMSS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  zones               = ["1", "2", "3"]
  # No health_probe_id configured - should trigger validation error

  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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

  resiliency {
    automatic_zone_rebalancing {
      rebalance_behavior = "CreateBeforeDelete"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) basicWithoutRollingUpgrade(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestVMSS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  upgrade_mode        = "Manual"

  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) resiliencyVMPoliciesOnlyWithoutRollingUpgrade(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestVMSS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  upgrade_mode        = "Manual"

  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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

  resilient_vm_creation_enabled = true
  resilient_vm_deletion_enabled = true
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) resiliencyVMPoliciesOnlyMultipleZonesWithHealthProbe(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  depends_on = [azurerm_lb_rule.test]

  name                = "acctestVMSS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 3
  zones               = ["1", "2", "3"]
  health_probe_id     = azurerm_lb_probe.test.id
  upgrade_mode        = "Rolling"
  overprovision       = false

  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false

  rolling_upgrade_policy {
    cross_zone_upgrades_enabled             = true
    max_batch_instance_percent              = 20
    max_unhealthy_instance_percent          = 20
    max_unhealthy_upgraded_instance_percent = 20
    pause_time_between_batches              = "PT0S"
    prioritize_unhealthy_instances_enabled  = false
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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
      name                                   = "internal"
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
    }
  }

  resilient_vm_creation_enabled = true
  resilient_vm_deletion_enabled = true
}
`, r.template(data), r.loadBalancerInfrastructureTemplate(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) resiliencyMinimalWithHealthProbeConsistent(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  depends_on = [azurerm_lb_rule.test]

  name                = "acctestVMSS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 3
  zones               = ["1", "2", "3"]
  health_probe_id     = azurerm_lb_probe.test.id
  upgrade_mode        = "Rolling"
  overprovision       = false

  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false

  rolling_upgrade_policy {
    cross_zone_upgrades_enabled             = true
    max_batch_instance_percent              = 20
    max_unhealthy_instance_percent          = 20
    max_unhealthy_upgraded_instance_percent = 20
    pause_time_between_batches              = "PT0S"
    prioritize_unhealthy_instances_enabled  = false
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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
      name                                   = "internal"
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
    }
  }

  resilient_vm_creation_enabled = false
  resilient_vm_deletion_enabled = false
}
`, r.template(data), r.loadBalancerInfrastructureTemplate(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) resiliencyMinimalWithHealthProbeConsistentScaledDown(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  depends_on = [azurerm_lb_rule.test]

  name                = "acctestVMSS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 0
  zones               = ["1", "2", "3"]
  health_probe_id     = azurerm_lb_probe.test.id
  upgrade_mode        = "Rolling"
  overprovision       = false

  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false

  rolling_upgrade_policy {
    cross_zone_upgrades_enabled             = true
    max_batch_instance_percent              = 20
    max_unhealthy_instance_percent          = 20
    max_unhealthy_upgraded_instance_percent = 20
    pause_time_between_batches              = "PT0S"
    prioritize_unhealthy_instances_enabled  = false
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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
      name                                   = "internal"
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
    }
  }

  resilient_vm_creation_enabled = false
  resilient_vm_deletion_enabled = false
}
`, r.templateWithForceDelete(data), r.loadBalancerInfrastructureTemplate(data), data.RandomInteger)
}
