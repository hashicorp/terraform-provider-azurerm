package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMAutoscaleSetting_basic(t *testing.T) {

	ri := acctest.RandInt()
	config := testAccAzureRMAutoscaleSetting_basic(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAvailabilitySetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutoscaleSettingExistsfunc,
					resource.TestCheckResourceAttr(
						"azurerm_autoscale_setting.test", "enabled", "true"),
					resource.TestCheckResourceAttr(
						"azurerm_autoscale_setting.test", "profile.#", "1"),
					resource.TestCheckResourceAttr(
						"azurerm_autoscale_setting.test", "profile.0.name", "defaultProfile"),
					resource.TestCheckResourceAttr(
						"azurerm_autoscale_setting.test", "profile.0.rule.#", "2"),
					resource.TestCheckResourceAttr(
						"azurerm_autoscale_setting.test", "notification.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMAutoscaleSetting_recurrence(t *testing.T) {

	ri := acctest.RandInt()
	config := testAccAzureRMAutoscaleSetting_recurrence(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAvailabilitySetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutoscaleSettingExistsfunc,
					resource.TestCheckResourceAttr(
						"azurerm_autoscale_setting.test", "enabled", "true"),
					resource.TestCheckResourceAttr(
						"azurerm_autoscale_setting.test", "profile.#", "2"),
					resource.TestCheckResourceAttr(
						"azurerm_autoscale_setting.test", "profile.0.name", "defaultProfile"),
					resource.TestCheckResourceAttr(
						"azurerm_autoscale_setting.test", "profile.0.rule.#", "2"),
					resource.TestCheckResourceAttr(
						"azurerm_autoscale_setting.test", "profile.1.recurrence.#", "1"),
					resource.TestCheckResourceAttr(
						"azurerm_autoscale_setting.test", "notification.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMAutoscaleSetting_fixedDate(t *testing.T) {

	ri := acctest.RandInt()
	config := testAccAzureRMAutoscaleSetting_fixedDate(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAvailabilitySetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutoscaleSettingExistsfunc,
					resource.TestCheckResourceAttr(
						"azurerm_autoscale_setting.test", "enabled", "true"),
					resource.TestCheckResourceAttr(
						"azurerm_autoscale_setting.test", "profile.#", "2"),
					resource.TestCheckResourceAttr(
						"azurerm_autoscale_setting.test", "profile.0.name", "defaultProfile"),
					resource.TestCheckResourceAttr(
						"azurerm_autoscale_setting.test", "profile.0.rule.#", "2"),
					resource.TestCheckResourceAttr(
						"azurerm_autoscale_setting.test", "profile.1.fixed_date.#", "1"),
					resource.TestCheckResourceAttr(
						"azurerm_autoscale_setting.test", "notification.#", "1"),
				),
			},
		},
	})
}

func testCheckAzureRMAutoscaleSettingExistsfunc(s *terraform.State) error {
	resourceName := "azurerm_autoscale_setting.test"
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return fmt.Errorf("Not found: %s", resourceName)
	}

	autoscaleSettingName := rs.Primary.Attributes["name"]
	resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
	if !hasResourceGroup {
		return fmt.Errorf("Bad: no resource group found in state for Autoscale Setting: %s", autoscaleSettingName)
	}

	asClient := testAccProvider.Meta().(*ArmClient).autoscaleSettingsClient

	resp, err := asClient.Get(resourceGroup, autoscaleSettingName)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Autoscale Setting %q (resource group: %q) does not exist", autoscaleSettingName, resourceGroup)
		}

		return fmt.Errorf("Bad: Get on autoscaleSettingsClient: %s", err)
	}

	return nil
}

func testCheckAzureRMAutoscaleSettingDestroy(s *terraform.State) error {
	asClient := testAccProvider.Meta().(*ArmClient).autoscaleSettingsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_autoscale_setting" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := asClient.Get(resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Autoscale Setting still exists:\n%#v", resp.AutoscaleSetting)
		}
	}

	return nil
}

func testAccAzureRMAutoscaleSetting_basic(rInt int) string {
	return fmt.Sprintf(`
    variable "time_zone" {
      type    = "string"
      default = "Pacific Standard Time"
    }

    variable "adminPwd" {
      type    = "string"
      default = "^l&f_ZwX2L4C"
    }

    resource "azurerm_resource_group" "test" {
      name     = "acctas%d"
      location = "West US"
    }

    resource "azurerm_virtual_network" "test" { 
      name                = "acctvn" 
      address_space       = ["10.0.0.0/16"] 
      location            = "${azurerm_resource_group.test.location}" 
      resource_group_name = "${azurerm_resource_group.test.name}" 
    } 
    
    resource "azurerm_subnet" "test" { 
      name                 = "acctsub" 
      resource_group_name  = "${azurerm_resource_group.test.name}" 
      virtual_network_name = "${azurerm_virtual_network.test.name}" 
      address_prefix       = "10.0.2.0/24" 
    } 
    
    resource "azurerm_public_ip" "test" { 
      name                         = "acctPublicIp" 
      location                     = "${azurerm_resource_group.test.location}" 
      resource_group_name          = "${azurerm_resource_group.test.name}" 
      public_ip_address_allocation = "static" 
    } 
    
    resource "azurerm_lb" "test" { 
      name                = "acctlb" 
      location            = "${azurerm_resource_group.test.location}" 
      resource_group_name = "${azurerm_resource_group.test.name}" 

      frontend_ip_configuration { 
        name                 = "PublicIPAddress" 
        public_ip_address_id = "${azurerm_public_ip.test.id}" 
      } 
    } 
    
    resource "azurerm_lb_backend_address_pool" "bpepool" { 
      name                = "BackEndAddressPool" 
      resource_group_name = "${azurerm_resource_group.test.name}" 
      loadbalancer_id     = "${azurerm_lb.test.id}" 
    } 
    
    resource "azurerm_lb_nat_pool" "lbnatpool" { 
      count                          = 3 
      resource_group_name            = "${azurerm_resource_group.test.name}" 
      name                           = "ssh" 
      loadbalancer_id                = "${azurerm_lb.test.id}" 
      protocol                       = "Tcp" 
      frontend_port_start            = 50000 
      frontend_port_end              = 50119 
      backend_port                   = 22 
      frontend_ip_configuration_name = "PublicIPAddress" 
    } 
    
    resource "azurerm_virtual_machine_scale_set" "test" { 
      name                = "myTestScaleset-1" 
      location            = "${azurerm_resource_group.test.location}" 
      resource_group_name = "${azurerm_resource_group.test.name}" 
      upgrade_policy_mode = "Manual" 
    
      sku { 
        name     = "Standard_A0" 
        tier     = "Standard" 
        capacity = 1 
      } 
    
      storage_profile_image_reference { 
        publisher = "Canonical" 
        offer     = "UbuntuServer" 
        sku       = "14.04.2-LTS" 
        version   = "latest" 
      } 
    
      storage_profile_os_disk { 
        name              = "" 
        caching           = "ReadWrite" 
        create_option     = "FromImage" 
        managed_disk_type = "Standard_LRS" 
      } 
    
      storage_profile_data_disk { 
        lun            = 0 
        caching        = "ReadWrite" 
        create_option  = "Empty" 
        disk_size_gb   = 10  
      } 
    
      os_profile { 
        computer_name_prefix = "testvm" 
        admin_username       = "myadmin" 
        admin_password       = "${var.adminPwd}" 
      } 
    
      os_profile_linux_config { 
        disable_password_authentication = false 
      } 
    
      network_profile { 
        name    = "terraformnetworkprofile" 
        primary = true 
    
        ip_configuration { 
          name                                   = "TestIPConfiguration" 
          subnet_id                              = "${azurerm_subnet.test.id}" 
          load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.bpepool.id}"] 
          load_balancer_inbound_nat_rules_ids    = ["${element(azurerm_lb_nat_pool.lbnatpool.*.id, count.index)}"] 
        } 
      } 
    
      tags { 
        environment = "staging" 
      } 
    } 

    resource "azurerm_autoscale_setting" "test" {
      name                = "autoScale%[1]d"
      enabled             = true
      resource_group_name = "${azurerm_resource_group.test.name}"
      location            = "${azurerm_resource_group.test.location}"
      target_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
      profile {
        name = "defaultProfile"
        capacity {
          default = 1
          minimum = 1
          maximum = 10
        }
        rule {
          metric_trigger {
            metric_name         = "Percentage CPU"
            metric_resource_id   = "${azurerm_virtual_machine_scale_set.test.id}"
            time_grain           = "PT1M"
            statistic            = "Average"
            time_window          = "PT5M"
            time_aggregation     = "Average"
            operator             = "GreaterThan"
            threshold            = 75
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
            metric_name         = "Percentage CPU"
            metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
            time_grain          = "PT1M"
            statistic           = "Average"
            time_window         = "PT5M"
            time_aggregation    = "Average"
            operator            = "LessThan"
            threshold           = 25
          }
          scale_action {
            direction = "Decrease"
            type      = "ChangeCount"
            value     = "1"
            cooldown  = "PT1M"
          }
        }
      }

      notification {
        operation = "Scale"
        email {
          send_to_subscription_administrator    = false
          send_to_subscription_co_administrator = false
        }
      }
    }`, rInt)
}

func testAccAzureRMAutoscaleSetting_recurrence(rInt int) string {
	return fmt.Sprintf(`
    variable "time_zone" {
      type    = "string"
      default = "Pacific Standard Time"
    }

    variable "adminPwd" {
      type    = "string"
      default = "^l&f_ZwX2L4C"
    }

    resource "azurerm_resource_group" "test" {
      name     = "acctas%d"
      location = "West US"
    }

    resource "azurerm_virtual_network" "test" { 
      name                = "acctvn" 
      address_space       = ["10.0.0.0/16"] 
      location            = "${azurerm_resource_group.test.location}" 
      resource_group_name = "${azurerm_resource_group.test.name}" 
    } 
    
    resource "azurerm_subnet" "test" { 
      name                 = "acctsub" 
      resource_group_name  = "${azurerm_resource_group.test.name}" 
      virtual_network_name = "${azurerm_virtual_network.test.name}" 
      address_prefix       = "10.0.2.0/24" 
    } 
    
    resource "azurerm_public_ip" "test" { 
      name                         = "acctPublicIp" 
      location                     = "${azurerm_resource_group.test.location}" 
      resource_group_name          = "${azurerm_resource_group.test.name}" 
      public_ip_address_allocation = "static" 
    } 
    
    resource "azurerm_lb" "test" { 
      name                = "acctlb" 
      location            = "${azurerm_resource_group.test.location}" 
      resource_group_name = "${azurerm_resource_group.test.name}" 

      frontend_ip_configuration { 
        name                 = "PublicIPAddress" 
        public_ip_address_id = "${azurerm_public_ip.test.id}" 
      } 
    } 
    
    resource "azurerm_lb_backend_address_pool" "bpepool" { 
      name                = "BackEndAddressPool" 
      resource_group_name = "${azurerm_resource_group.test.name}" 
      loadbalancer_id     = "${azurerm_lb.test.id}" 
    } 
    
    resource "azurerm_lb_nat_pool" "lbnatpool" { 
      count                          = 3 
      resource_group_name            = "${azurerm_resource_group.test.name}" 
      name                           = "ssh" 
      loadbalancer_id                = "${azurerm_lb.test.id}" 
      protocol                       = "Tcp" 
      frontend_port_start            = 50000 
      frontend_port_end              = 50119 
      backend_port                   = 22 
      frontend_ip_configuration_name = "PublicIPAddress" 
    } 
    
    resource "azurerm_virtual_machine_scale_set" "test" { 
      name                = "myTestScaleset-1" 
      location            = "${azurerm_resource_group.test.location}" 
      resource_group_name = "${azurerm_resource_group.test.name}" 
      upgrade_policy_mode = "Manual" 
    
      sku { 
        name     = "Standard_A0" 
        tier     = "Standard" 
        capacity = 1 
      } 
    
      storage_profile_image_reference { 
        publisher = "Canonical" 
        offer     = "UbuntuServer" 
        sku       = "14.04.2-LTS" 
        version   = "latest" 
      } 
    
      storage_profile_os_disk { 
        name              = "" 
        caching           = "ReadWrite" 
        create_option     = "FromImage" 
        managed_disk_type = "Standard_LRS" 
      } 
    
      storage_profile_data_disk { 
        lun            = 0 
        caching        = "ReadWrite" 
        create_option  = "Empty" 
        disk_size_gb   = 10  
      } 
    
      os_profile { 
        computer_name_prefix = "testvm" 
        admin_username       = "myadmin" 
        admin_password       = "${var.adminPwd}" 
      } 
    
      os_profile_linux_config { 
        disable_password_authentication = false 
      } 
    
      network_profile { 
        name    = "terraformnetworkprofile" 
        primary = true 
    
        ip_configuration { 
          name                                   = "TestIPConfiguration" 
          subnet_id                              = "${azurerm_subnet.test.id}" 
          load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.bpepool.id}"] 
          load_balancer_inbound_nat_rules_ids    = ["${element(azurerm_lb_nat_pool.lbnatpool.*.id, count.index)}"] 
        } 
      } 
    
      tags { 
        environment = "staging" 
      } 
    } 

    resource "azurerm_autoscale_setting" "test" {
      name                = "autoScale%[1]d"
      enabled             = true
      resource_group_name = "${azurerm_resource_group.test.name}"
      location            = "${azurerm_resource_group.test.location}"
      target_resource_id = "${azurerm_virtual_machine_scale_set.test.id}"
      profile {
        name = "defaultProfile"
        capacity {
          default = 1
          minimum = 1
          maximum = 10
        }
        rule {
          metric_trigger {
            metric_name         = "Percentage CPU"
            metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
            time_grain          = "PT1M"
            statistic           = "Average"
            time_window         = "PT5M"
            time_aggregation    = "Average"
            operator            = "GreaterThan"
            threshold           = 75
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
            metric_name         = "Percentage CPU"
            metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
            time_grain          = "PT1M"
            statistic           = "Average"
            time_window         = "PT5M"
            time_aggregation    = "Average"
            operator            = "LessThan"
            threshold           = 25
          }
          scale_action {
            direction = "Decrease"
            type      = "ChangeCount"
            value     = "1"
            cooldown  = "PT1M"
          }
        }
      }

      profile {
        name = "recurrence"
        capacity {
          default = 1
          minimum = 1
          maximum = 10
        }
        rule {
          metric_trigger {
            metric_name         = "Percentage CPU"
            metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
            time_grain          = "PT1M"
            statistic           = "Average"
            time_window         = "PT5M"
            time_aggregation    = "Average"
            operator            = "GreaterThan"
            threshold           = 75
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
            metric_name         = "Percentage CPU"
            metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
            time_grain          = "PT1M"
            statistic           = "Average"
            time_window         = "PT5M"
            time_aggregation    = "Average"
            operator            = "LessThan"
            threshold           = 25
          }
          scale_action {
            direction = "Decrease"
            type      = "ChangeCount"
            value     = "1"
            cooldown  = "PT1M"
          }
        }

        recurrence {
          frequency = "Week"
          schedule {
            time_zone = "${var.time_zone}"
            days      = [
              "Monday",
              "Wednesday",
              "Friday"
            ]
            hours     = [ 18 ]
            minutes   = [ 0 ]
          }
        }
      }

      notification {
        operation = "Scale"
        email {
          send_to_subscription_administrator    = false
          send_to_subscription_co_administrator = false
        }
      }
    }`, rInt)
}

func testAccAzureRMAutoscaleSetting_fixedDate(rInt int) string {
	return fmt.Sprintf(`
    variable "time_zone" {
      type    = "string"
      default = "Pacific Standard Time"
    }

    variable "adminPwd" {
      type    = "string"
      default = "^l&f_ZwX2L4C"
    }

    resource "azurerm_resource_group" "test" {
      name     = "acctas%d"
      location = "West US"
    }

    resource "azurerm_virtual_network" "test" { 
      name                = "acctvn" 
      address_space       = ["10.0.0.0/16"] 
      location            = "${azurerm_resource_group.test.location}" 
      resource_group_name = "${azurerm_resource_group.test.name}" 
    } 
    
    resource "azurerm_subnet" "test" { 
      name                 = "acctsub" 
      resource_group_name  = "${azurerm_resource_group.test.name}" 
      virtual_network_name = "${azurerm_virtual_network.test.name}" 
      address_prefix       = "10.0.2.0/24" 
    } 
    
    resource "azurerm_public_ip" "test" { 
      name                         = "acctPublicIp" 
      location                     = "${azurerm_resource_group.test.location}" 
      resource_group_name          = "${azurerm_resource_group.test.name}" 
      public_ip_address_allocation = "static" 
    } 
    
    resource "azurerm_lb" "test" { 
      name                = "acctlb" 
      location            = "${azurerm_resource_group.test.location}" 
      resource_group_name = "${azurerm_resource_group.test.name}" 

      frontend_ip_configuration { 
        name                 = "PublicIPAddress" 
        public_ip_address_id = "${azurerm_public_ip.test.id}" 
      } 
    } 
    
    resource "azurerm_lb_backend_address_pool" "bpepool" { 
      name                = "BackEndAddressPool" 
      resource_group_name = "${azurerm_resource_group.test.name}" 
      loadbalancer_id     = "${azurerm_lb.test.id}" 
    } 
    
    resource "azurerm_lb_nat_pool" "lbnatpool" { 
      count                          = 3 
      resource_group_name            = "${azurerm_resource_group.test.name}" 
      name                           = "ssh" 
      loadbalancer_id                = "${azurerm_lb.test.id}" 
      protocol                       = "Tcp" 
      frontend_port_start            = 50000 
      frontend_port_end              = 50119 
      backend_port                   = 22 
      frontend_ip_configuration_name = "PublicIPAddress" 
    } 
    
    resource "azurerm_virtual_machine_scale_set" "test" { 
      name                = "myTestScaleset-1" 
      location            = "${azurerm_resource_group.test.location}" 
      resource_group_name = "${azurerm_resource_group.test.name}" 
      upgrade_policy_mode = "Manual" 
    
      sku { 
        name     = "Standard_A0" 
        tier     = "Standard" 
        capacity = 1 
      } 
    
      storage_profile_image_reference { 
        publisher = "Canonical" 
        offer     = "UbuntuServer" 
        sku       = "14.04.2-LTS" 
        version   = "latest" 
      } 
    
      storage_profile_os_disk { 
        name              = "" 
        caching           = "ReadWrite" 
        create_option     = "FromImage" 
        managed_disk_type = "Standard_LRS" 
      } 
    
      storage_profile_data_disk { 
        lun            = 0 
        caching        = "ReadWrite" 
        create_option  = "Empty" 
        disk_size_gb   = 10  
      } 
    
      os_profile { 
        computer_name_prefix = "testvm" 
        admin_username       = "myadmin" 
        admin_password       = "${var.adminPwd}" 
      } 
    
      os_profile_linux_config { 
        disable_password_authentication = false 
      } 
    
      network_profile { 
        name    = "terraformnetworkprofile" 
        primary = true 
    
        ip_configuration { 
          name                                   = "TestIPConfiguration" 
          subnet_id                              = "${azurerm_subnet.test.id}" 
          load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.bpepool.id}"] 
          load_balancer_inbound_nat_rules_ids    = ["${element(azurerm_lb_nat_pool.lbnatpool.*.id, count.index)}"] 
        } 
      } 
    
      tags { 
        environment = "staging" 
      } 
    } 

    resource "azurerm_autoscale_setting" "test" {
      name                = "autoScale%[1]d"
      enabled             = true
      resource_group_name = "${azurerm_resource_group.test.name}"
      location            = "${azurerm_resource_group.test.location}"
      target_resource_id = "${azurerm_virtual_machine_scale_set.test.id}"
      profile {
        name = "defaultProfile"
        capacity {
          default = 1
          minimum = 1
          maximum = 10
        }
        rule {
          metric_trigger {
            metric_name         = "Percentage CPU"
            metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
            time_grain          = "PT1M"
            statistic           = "Average"
            time_window         = "PT5M"
            time_aggregation    = "Average"
            operator            = "GreaterThan"
            threshold           = 75
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
            metric_name         = "Percentage CPU"
            metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
            time_grain          = "PT1M"
            statistic           = "Average"
            time_window         = "PT5M"
            time_aggregation    = "Average"
            operator            = "LessThan"
            threshold           = 25
          }
          scale_action {
            direction = "Decrease"
            type      = "ChangeCount"
            value     = "1"
            cooldown  = "PT1M"
          }
        }
      }

      profile {
        name = "fixedDate"
        capacity {
          default = 1
          minimum = 1
          maximum = 10
        }
        rule {
          metric_trigger {
            metric_name         = "Percentage CPU"
            metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
            time_grain          = "PT1M"
            statistic           = "Average"
            time_window         = "PT5M"
            time_aggregation    = "Average"
            operator            = "GreaterThan"
            threshold           = 75
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
            metric_name         = "Percentage CPU"
            metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
            time_grain          = "PT1M"
            statistic           = "Average"
            time_window         = "PT5M"
            time_aggregation    = "Average"
            operator            = "LessThan"
            threshold           = 25
          }
          scale_action {
            direction = "Decrease"
            type      = "ChangeCount"
            value     = "1"
            cooldown  = "PT1M"
          }
        }

        fixed_date {
          time_zone = "${var.time_zone}"
          start     = "2020-06-18T00:00:00Z"
          end       = "2020-06-18T23:59:59Z"
        }
      }

      notification {
        operation = "Scale"
        email {
          send_to_subscription_administrator    = false
          send_to_subscription_co_administrator = false
        }
      }
    }`, rInt)
}
