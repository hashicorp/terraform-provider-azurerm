// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package computefleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccComputeFleet_virtualMachineProfileNetwork_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkProfileBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileNetwork_completeForBaseVirtualMachineProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.netWorkProfileCompleteForBaseVirtualMachineProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileNetwork_completeForAdditionalLocation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.netWorkProfileCompleteForAdditionalLocation(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileNetwork_updateForBaseVirtualMachineProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.netWorkProfileCompleteForBaseVirtualMachineProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.netWorkProfileCompleteForBaseVirtualMachineProfileUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			// Test with `networkProfileBasicWithZones` as `netWorkProfileComplete` test requires `zones`
			Config: r.networkProfileBasicWithZones(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.netWorkProfileCompleteForBaseVirtualMachineProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileNetwork_updateForAdditionalLocation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.netWorkProfileCompleteForAdditionalLocation(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.netWorkProfileCompleteForAdditionalLocationUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.networkProfileBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.netWorkProfileCompleteForAdditionalLocation(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func (r ComputeFleetTestResource) networkProfileBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 1

  spot_priority_profile {
    min_capacity     = 0
    maintain_enabled = false
    capacity         = 1
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
    os_profile {
      linux_configuration {
        computer_name_prefix            = "testvm"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name                              = "nic-test"
      primary_network_interface_enabled = true

      ip_configuration {
        name                                   = "primary"
        primary_ip_configuration_enabled       = true
        subnet_id                              = azurerm_subnet.test.id
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
      }
    }

    os_disk {
      storage_account_type = "Standard_LRS"
      caching              = "ReadWrite"
    }

    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }
  }

  additional_location_profile {
    location = "%[4]s"
    virtual_machine_profile_override {
      network_api_version = "2020-11-01"
      os_profile {
        linux_configuration {
          computer_name_prefix            = "testvm"
          admin_username                  = local.admin_username
          admin_password                  = local.admin_password
          password_authentication_enabled = true
        }
      }

      network_interface {
        name                              = "nic-test"
        primary_network_interface_enabled = true

        ip_configuration {
          name                                   = "primary"
          primary_ip_configuration_enabled       = true
          subnet_id                              = azurerm_subnet.linux_test.id
          load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.linux_test.id]
        }
      }

      os_disk {
        storage_account_type = "Standard_LRS"
        caching              = "ReadWrite"
      }

      source_image_reference {
        publisher = "Canonical"
        offer     = "0001-com-ubuntu-server-jammy"
        sku       = "22_04-lts"
        version   = "latest"
      }
    }
  }
}
`, r.baseAndAdditionalLocationLinuxTemplate(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) networkProfileBasicWithZones(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 1
  zones                       = ["1", "2"]

  spot_priority_profile {
    min_capacity     = 0
    maintain_enabled = false
    capacity         = 1
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  compute_api_version = "2024-03-01"

  virtual_machine_profile {
    network_api_version = "2020-11-01"
    os_profile {
      linux_configuration {
        computer_name_prefix            = "testvm"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name                              = "nic-test"
      primary_network_interface_enabled = true

      ip_configuration {
        name                                   = "primary"
        primary_ip_configuration_enabled       = true
        subnet_id                              = azurerm_subnet.test.id
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
      }
    }

    os_disk {
      storage_account_type = "Standard_LRS"
      caching              = "ReadWrite"
    }

    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetTestResource) netWorkProfileCompleteForBaseVirtualMachineProfile(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[3]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[4]s"
  platform_fault_domain_count = 1

  zones = ["1", "2"]

  spot_priority_profile {
    min_capacity     = 0
    maintain_enabled = false
    capacity         = 0
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  compute_api_version = "2024-03-01"

  virtual_machine_profile {
    network_api_version = "2022-11-01"
    os_profile {
      linux_configuration {
        computer_name_prefix            = "testvm"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name                              = "nic-test"
      primary_network_interface_enabled = true
      accelerated_networking_enabled    = true
      ip_forwarding_enabled             = true
      auxiliary_mode                    = "AcceleratedConnections"
      auxiliary_sku                     = "A2"
      delete_option                     = "Delete"
      network_security_group_id         = azurerm_network_security_group.test.id
      dns_servers                       = ["8.8.8.8", "8.8.4.4"]
      ip_configuration {
        name                                         = "first"
        primary_ip_configuration_enabled             = true
        subnet_id                                    = azurerm_subnet.test.id
        application_gateway_backend_address_pool_ids = [tolist(azurerm_application_gateway.test.backend_address_pool)[0].id]
        application_security_group_ids               = [azurerm_application_security_group.test.id]
        load_balancer_backend_address_pool_ids       = [azurerm_lb_backend_address_pool.test.id]
        public_ip_address {
          name                    = "nic-pip-first"
          delete_option           = "Delete"
          domain_name_label       = "test-domain-label"
          domain_name_label_scope = "ResourceGroupReuse"
          idle_timeout_in_minutes = 4
          sku_name                = "Standard_Regional"
          version                 = "IPv4"
          ip_tag {
            type = "RoutingPreference"
            tag  = "Internet"
          }
        }
        version = "IPv4"
      }

      ip_configuration {
        name                           = "second"
        subnet_id                      = azurerm_subnet.test.id
        application_security_group_ids = [azurerm_application_security_group.test.id]
        public_ip_address {
          name                    = "nic-pip-second"
          idle_timeout_in_minutes = 15
        }
      }
    }

    os_disk {
      storage_account_type = "Standard_LRS"
      caching              = "ReadWrite"
    }

    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }
  }
}
`, r.template(data), r.netWorkProfileBaseVirtualMachineProfileResourceDependencies(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetTestResource) netWorkProfileCompleteForAdditionalLocation(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[3]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[4]s"
  platform_fault_domain_count = 1

  spot_priority_profile {
    min_capacity     = 0
    maintain_enabled = false
    capacity         = 0
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  compute_api_version = "2024-03-01"

  virtual_machine_profile {
    network_api_version = "2022-11-01"
    os_profile {
      linux_configuration {
        computer_name_prefix            = "testvm"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name                              = "nic-test"
      primary_network_interface_enabled = true
      accelerated_networking_enabled    = true
      ip_forwarding_enabled             = true
      auxiliary_mode                    = "AcceleratedConnections"
      auxiliary_sku                     = "A2"
      delete_option                     = "Delete"
      network_security_group_id         = azurerm_network_security_group.test.id
      dns_servers                       = ["8.8.8.8", "8.8.4.4"]
      ip_configuration {
        name                                         = "first"
        primary_ip_configuration_enabled             = true
        subnet_id                                    = azurerm_subnet.test.id
        application_gateway_backend_address_pool_ids = [tolist(azurerm_application_gateway.test.backend_address_pool)[0].id]
        application_security_group_ids               = [azurerm_application_security_group.test.id]
        load_balancer_backend_address_pool_ids       = [azurerm_lb_backend_address_pool.test.id]
        public_ip_address {
          name                    = "nic-pip-first"
          delete_option           = "Delete"
          domain_name_label       = "test-domain-label"
          domain_name_label_scope = "ResourceGroupReuse"
          idle_timeout_in_minutes = 4
          sku_name                = "Standard_Regional"
          version                 = "IPv4"
          ip_tag {
            type = "RoutingPreference"
            tag  = "Internet"
          }
        }
        version = "IPv4"
      }

      ip_configuration {
        name                           = "second"
        subnet_id                      = azurerm_subnet.test.id
        application_security_group_ids = [azurerm_application_security_group.test.id]
        public_ip_address {
          name                    = "nic-pip-second"
          idle_timeout_in_minutes = 15
        }
      }
    }

    os_disk {
      storage_account_type = "Standard_LRS"
      caching              = "ReadWrite"
    }

    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }
  }

  additional_location_profile {
    location = "%[5]s"

    virtual_machine_profile_override {
      network_api_version = "2022-11-01"
      os_profile {
        linux_configuration {
          computer_name_prefix            = "testvm"
          admin_username                  = local.admin_username
          admin_password                  = local.admin_password
          password_authentication_enabled = true
        }
      }

      network_interface {
        name                              = "nic-test"
        primary_network_interface_enabled = true
        accelerated_networking_enabled    = true
        ip_forwarding_enabled             = true
        auxiliary_mode                    = "AcceleratedConnections"
        auxiliary_sku                     = "A2"
        delete_option                     = "Delete"
        network_security_group_id         = azurerm_network_security_group.linux_test.id
        dns_servers                       = ["8.8.8.8", "8.8.4.4"]
        ip_configuration {
          name                                         = "first"
          primary_ip_configuration_enabled             = true
          subnet_id                                    = azurerm_subnet.linux_test.id
          application_gateway_backend_address_pool_ids = [tolist(azurerm_application_gateway.linux_test.backend_address_pool)[0].id]
          application_security_group_ids               = [azurerm_application_security_group.linux_test.id]
          load_balancer_backend_address_pool_ids       = [azurerm_lb_backend_address_pool.linux_test.id]
          public_ip_address {
            name                    = "nic-pip-first"
            delete_option           = "Delete"
            domain_name_label       = "test-domain-label"
            domain_name_label_scope = "ResourceGroupReuse"
            idle_timeout_in_minutes = 4
            sku_name                = "Standard_Regional"
            version                 = "IPv4"
            ip_tag {
              type = "RoutingPreference"
              tag  = "Internet"
            }
          }
          version = "IPv4"
        }

        ip_configuration {
          name                           = "second"
          subnet_id                      = azurerm_subnet.linux_test.id
          application_security_group_ids = [azurerm_application_security_group.linux_test.id]
          public_ip_address {
            name                    = "nic-pip-second"
            idle_timeout_in_minutes = 15
          }
        }
      }

      os_disk {
        storage_account_type = "Standard_LRS"
        caching              = "ReadWrite"
      }

      source_image_reference {
        publisher = "Canonical"
        offer     = "0001-com-ubuntu-server-jammy"
        sku       = "22_04-lts"
        version   = "latest"
      }
    }
  }
}
`, r.baseAndAdditionalLocationLinuxTemplate(data), r.netWorkProfileAdditionalLocationResourceDependencies(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) netWorkProfileCompleteForBaseVirtualMachineProfileUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[3]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[4]s"
  platform_fault_domain_count = 1

  zones = ["1", "2"]

  spot_priority_profile {
    min_capacity     = 0
    maintain_enabled = false
    capacity         = 0
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  compute_api_version = "2024-03-01"

  virtual_machine_profile {
    network_api_version = "2022-11-01"
    os_profile {
      linux_configuration {
        computer_name_prefix            = "testvm"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name                              = "nic-test"
      primary_network_interface_enabled = false
      accelerated_networking_enabled    = false
      ip_forwarding_enabled             = false
      delete_option                     = "Detach"
      network_security_group_id         = azurerm_network_security_group.other.id
      dns_servers                       = ["8.8.8.8"]
      ip_configuration {
        name                                         = "first"
        primary_ip_configuration_enabled             = true
        subnet_id                                    = azurerm_subnet.test.id
        application_gateway_backend_address_pool_ids = [tolist(azurerm_application_gateway.test.backend_address_pool)[0].id]
        application_security_group_ids               = [azurerm_application_security_group.test.id, azurerm_application_security_group.other.id]
        load_balancer_backend_address_pool_ids       = [azurerm_lb_backend_address_pool.test.id]
        public_ip_address {
          name                    = "nic-pip-first"
          delete_option           = "Detach"
          domain_name_label       = "test-domain-label-update"
          domain_name_label_scope = "SubscriptionReuse"
          idle_timeout_in_minutes = 14
          sku_name                = "Standard_Regional"
          version                 = "IPv4"
          ip_tag {
            type = "RoutingPreference"
            tag  = "Internet"
          }
        }
        version = "IPv4"
      }

      ip_configuration {
        name                           = "second"
        subnet_id                      = azurerm_subnet.test.id
        application_security_group_ids = [azurerm_application_security_group.test.id, azurerm_application_security_group.other.id]
        public_ip_address {
          name                    = "nic-pip-second"
          idle_timeout_in_minutes = 10
        }
      }
    }

    network_interface {
      name                              = "nic-test-multiple"
      primary_network_interface_enabled = true

      ip_configuration {
        name                                   = "primary"
        primary_ip_configuration_enabled       = true
        subnet_id                              = azurerm_subnet.test.id
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
      }
    }

    os_disk {
      storage_account_type = "Standard_LRS"
      caching              = "ReadWrite"
    }

    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }
  }
}
`, r.template(data), r.netWorkProfileBaseVirtualMachineProfileResourceDependencies(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetTestResource) netWorkProfileCompleteForAdditionalLocationUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[3]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[4]s"
  platform_fault_domain_count = 1

  spot_priority_profile {
    min_capacity     = 0
    maintain_enabled = false
    capacity         = 0
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  compute_api_version = "2024-03-01"

  virtual_machine_profile {
    network_api_version = "2022-11-01"
    os_profile {
      linux_configuration {
        computer_name_prefix            = "testvm"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name                              = "nic-test"
      primary_network_interface_enabled = false
      accelerated_networking_enabled    = false
      ip_forwarding_enabled             = false
      delete_option                     = "Detach"
      network_security_group_id         = azurerm_network_security_group.other.id
      dns_servers                       = ["8.8.8.8"]
      ip_configuration {
        name                                         = "first"
        primary_ip_configuration_enabled             = true
        subnet_id                                    = azurerm_subnet.test.id
        application_gateway_backend_address_pool_ids = [tolist(azurerm_application_gateway.test.backend_address_pool)[0].id]
        application_security_group_ids               = [azurerm_application_security_group.test.id, azurerm_application_security_group.other.id]
        load_balancer_backend_address_pool_ids       = [azurerm_lb_backend_address_pool.test.id]
        public_ip_address {
          name                    = "nic-pip-first"
          delete_option           = "Detach"
          domain_name_label       = "test-domain-label-update"
          domain_name_label_scope = "SubscriptionReuse"
          idle_timeout_in_minutes = 14
          sku_name                = "Standard_Regional"
          version                 = "IPv4"
          ip_tag {
            type = "RoutingPreference"
            tag  = "Internet"
          }
        }
        version = "IPv4"
      }

      ip_configuration {
        name                           = "second"
        subnet_id                      = azurerm_subnet.test.id
        application_security_group_ids = [azurerm_application_security_group.test.id, azurerm_application_security_group.other.id]
        public_ip_address {
          name                    = "nic-pip-second"
          idle_timeout_in_minutes = 10
        }
      }
    }

    network_interface {
      name                              = "nic-test-multiple"
      primary_network_interface_enabled = true

      ip_configuration {
        name                                   = "primary"
        primary_ip_configuration_enabled       = true
        subnet_id                              = azurerm_subnet.test.id
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
      }
    }

    os_disk {
      storage_account_type = "Standard_LRS"
      caching              = "ReadWrite"
    }

    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }
  }

  additional_location_profile {
    location = "%[5]s"

    virtual_machine_profile_override {
      network_api_version = "2022-11-01"
      os_profile {
        linux_configuration {
          computer_name_prefix            = "testvm"
          admin_username                  = local.admin_username
          admin_password                  = local.admin_password
          password_authentication_enabled = true
        }
      }

      network_interface {
        name                              = "nic-test"
        primary_network_interface_enabled = false
        accelerated_networking_enabled    = false
        ip_forwarding_enabled             = false
        delete_option                     = "Detach"
        network_security_group_id         = azurerm_network_security_group.linux_test_other.id
        dns_servers                       = ["8.8.8.8"]
        ip_configuration {
          name                                         = "first"
          primary_ip_configuration_enabled             = true
          subnet_id                                    = azurerm_subnet.linux_test.id
          application_gateway_backend_address_pool_ids = [tolist(azurerm_application_gateway.linux_test.backend_address_pool)[0].id]
          application_security_group_ids               = [azurerm_application_security_group.linux_test.id, azurerm_application_security_group.linux_test_other.id]
          load_balancer_backend_address_pool_ids       = [azurerm_lb_backend_address_pool.linux_test.id]
          public_ip_address {
            name                    = "nic-pip-first"
            delete_option           = "Detach"
            domain_name_label       = "test-domain-label-update"
            domain_name_label_scope = "SubscriptionReuse"
            idle_timeout_in_minutes = 14
            sku_name                = "Standard_Regional"
            version                 = "IPv4"
            ip_tag {
              type = "RoutingPreference"
              tag  = "Internet"
            }
          }
          version = "IPv4"
        }

        ip_configuration {
          name                           = "second"
          subnet_id                      = azurerm_subnet.linux_test.id
          application_security_group_ids = [azurerm_application_security_group.linux_test.id, azurerm_application_security_group.linux_test_other.id]
          public_ip_address {
            name                    = "nic-pip-second"
            idle_timeout_in_minutes = 10
          }
        }
      }

      network_interface {
        name                              = "nic-test-multiple"
        primary_network_interface_enabled = true

        ip_configuration {
          name                                   = "primary"
          primary_ip_configuration_enabled       = true
          subnet_id                              = azurerm_subnet.linux_test.id
          load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.linux_test.id]
        }
      }

      os_disk {
        storage_account_type = "Standard_LRS"
        caching              = "ReadWrite"
      }

      source_image_reference {
        publisher = "Canonical"
        offer     = "0001-com-ubuntu-server-jammy"
        sku       = "22_04-lts"
        version   = "latest"
      }
    }
  }
}
`, r.baseAndAdditionalLocationLinuxTemplate(data), r.netWorkProfileAdditionalLocationResourceDependencies(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) netWorkProfileAdditionalLocationResourceDependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_application_security_group" "test" {
  location            = azurerm_resource_group.test.location
  name                = "TestApplicationSecurityGroup"
  resource_group_name = azurerm_resource_group.test.name

  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_network_security_group" "other" {
  name                = "acceptanceTestSecurityGroup-other-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_application_security_group" "other" {
  location            = azurerm_resource_group.test.location
  name                = "TestApplicationSecurityGroupOther"
  resource_group_name = azurerm_resource_group.test.name

  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_subnet" "gwtest" {
  name                 = "gw-subnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.3.0/24"]
}

resource "azurerm_public_ip" "gwtest" {
  name                = "acctest-pubip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestgw-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "gw-ip-config1"
    subnet_id = azurerm_subnet.gwtest.id
  }

  frontend_ip_configuration {
    name                 = "ip-config-public"
    public_ip_address_id = azurerm_public_ip.gwtest.id
  }

  frontend_port {
    name = "port-8080"
    port = 8080
  }

  backend_address_pool {
    name = "pool-1"
  }

  backend_http_settings {
    name                  = "backend-http-1"
    port                  = 8010
    protocol              = "Http"
    cookie_based_affinity = "Enabled"
    request_timeout       = 30

    probe_name = "probe-1"
  }

  http_listener {
    name                           = "listener-1"
    frontend_ip_configuration_name = "ip-config-public"
    frontend_port_name             = "port-8080"
    protocol                       = "Http"
  }

  probe {
    name                = "probe-1"
    protocol            = "Http"
    path                = "/test"
    host                = "azure.com"
    timeout             = 120
    interval            = 300
    unhealthy_threshold = 8
    match {
      status_code = ["200-399"]
    }
  }

  request_routing_rule {
    name                       = "rule-basic-1"
    rule_type                  = "Basic"
    http_listener_name         = "listener-1"
    backend_address_pool_name  = "pool-1"
    backend_http_settings_name = "backend-http-1"

    priority = 10
  }

  tags = {
    environment = "tf01"
  }
}

resource "azurerm_network_security_group" "linux_test" {
  name                = "acceptanceTestSecurityGroup-%[1]d"
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name

  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_application_security_group" "linux_test" {
  location            = azurerm_resource_group.linux_test.location
  name                = "TestApplicationSecurityGroup"
  resource_group_name = azurerm_resource_group.linux_test.name

  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_network_security_group" "linux_test_other" {
  name                = "acceptanceTestSecurityGroup-other-%[1]d"
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name

  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_application_security_group" "linux_test_other" {
  location            = azurerm_resource_group.linux_test.location
  name                = "TestApplicationSecurityGroupOther"
  resource_group_name = azurerm_resource_group.linux_test.name

  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_subnet" "linux_test_gwtest" {
  name                 = "gw-subnet-%[1]d"
  resource_group_name  = azurerm_resource_group.linux_test.name
  virtual_network_name = azurerm_virtual_network.linux_test.name
  address_prefixes     = ["10.0.3.0/24"]
}

resource "azurerm_public_ip" "linux_test_gwtest" {
  name                = "acctest-pubip-%[1]d"
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_application_gateway" "linux_test" {
  name                = "acctestgw-%[1]d"
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "gw-ip-config1"
    subnet_id = azurerm_subnet.linux_test_gwtest.id
  }

  frontend_ip_configuration {
    name                 = "ip-config-public"
    public_ip_address_id = azurerm_public_ip.linux_test_gwtest.id
  }

  frontend_port {
    name = "port-8080"
    port = 8080
  }

  backend_address_pool {
    name = "pool-1"
  }

  backend_http_settings {
    name                  = "backend-http-1"
    port                  = 8010
    protocol              = "Http"
    cookie_based_affinity = "Enabled"
    request_timeout       = 30

    probe_name = "probe-1"
  }

  http_listener {
    name                           = "listener-1"
    frontend_ip_configuration_name = "ip-config-public"
    frontend_port_name             = "port-8080"
    protocol                       = "Http"
  }

  probe {
    name                = "probe-1"
    protocol            = "Http"
    path                = "/test"
    host                = "azure.com"
    timeout             = 120
    interval            = 300
    unhealthy_threshold = 8
    match {
      status_code = ["200-399"]
    }
  }

  request_routing_rule {
    name                       = "rule-basic-1"
    rule_type                  = "Basic"
    http_listener_name         = "listener-1"
    backend_address_pool_name  = "pool-1"
    backend_http_settings_name = "backend-http-1"

    priority = 10
  }

  tags = {
    environment = "tf01"
  }
}
`, data.RandomInteger)
}

func (r ComputeFleetTestResource) netWorkProfileBaseVirtualMachineProfileResourceDependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_application_security_group" "test" {
  location            = azurerm_resource_group.test.location
  name                = "TestApplicationSecurityGroup"
  resource_group_name = azurerm_resource_group.test.name

  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_network_security_group" "other" {
  name                = "acceptanceTestSecurityGroup-other-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_application_security_group" "other" {
  location            = azurerm_resource_group.test.location
  name                = "TestApplicationSecurityGroupOther"
  resource_group_name = azurerm_resource_group.test.name

  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_subnet" "gwtest" {
  name                 = "gw-subnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.3.0/24"]
}

resource "azurerm_public_ip" "gwtest" {
  name                = "acctest-pubip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestgw-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "gw-ip-config1"
    subnet_id = azurerm_subnet.gwtest.id
  }

  frontend_ip_configuration {
    name                 = "ip-config-public"
    public_ip_address_id = azurerm_public_ip.gwtest.id
  }

  frontend_port {
    name = "port-8080"
    port = 8080
  }

  backend_address_pool {
    name = "pool-1"
  }

  backend_http_settings {
    name                  = "backend-http-1"
    port                  = 8010
    protocol              = "Http"
    cookie_based_affinity = "Enabled"
    request_timeout       = 30

    probe_name = "probe-1"
  }

  http_listener {
    name                           = "listener-1"
    frontend_ip_configuration_name = "ip-config-public"
    frontend_port_name             = "port-8080"
    protocol                       = "Http"
  }

  probe {
    name                = "probe-1"
    protocol            = "Http"
    path                = "/test"
    host                = "azure.com"
    timeout             = 120
    interval            = 300
    unhealthy_threshold = 8
    match {
      status_code = ["200-399"]
    }
  }

  request_routing_rule {
    name                       = "rule-basic-1"
    rule_type                  = "Basic"
    http_listener_name         = "listener-1"
    backend_address_pool_name  = "pool-1"
    backend_http_settings_name = "backend-http-1"

    priority = 10
  }

  tags = {
    environment = "tf01"
  }
}
`, data.RandomInteger)
}
