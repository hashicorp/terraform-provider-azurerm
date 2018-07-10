package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAppServiceSlot_basic(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_32Bit(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_32Bit(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.use_32_bit_worker_process", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_alwaysOn(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_alwaysOn(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.always_on", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_appSettings(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_appSettings(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app_settings.foo", "bar"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_clientAffinityEnabled(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_clientAffinityEnabled(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "client_affinity_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_connectionStrings(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_connectionStrings(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "connection_string.0.name", "Example"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.0.value", "some-postgresql-connection-string"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.0.type", "PostgreSQL"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_defaultDocuments(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_defaultDocuments(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.default_documents.0", "first.html"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.default_documents.1", "second.jsp"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.default_documents.2", "third.aspx"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_enabled(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_enabled(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_httpsOnly(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_httpsOnly(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "https_only", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_http2Enabled(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_http2Enabled(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.http2_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_oneIpRestriction(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_oneIpRestriction(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.0.ip_address", "10.10.10.10"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.0.subnet_mask", "255.255.255.255"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_manyIpRestrictions(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_manyIpRestrictions(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.0.ip_address", "10.10.10.10"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.0.subnet_mask", "255.255.255.255"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.1.ip_address", "20.20.20.0"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.1.subnet_mask", "255.255.255.0"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.2.ip_address", "30.30.0.0"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.2.subnet_mask", "255.255.0.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_localMySql(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_localMySql(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.local_mysql_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_managedPipelineMode(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_managedPipelineMode(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.managed_pipeline_mode", "Classic"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_tagsUpdate(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_tags(ri, testLocation())
	updatedConfig := testAccAzureRMAppServiceSlot_tagsUpdated(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Hello", "World"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Hello", "World"),
					resource.TestCheckResourceAttr(resourceName, "tags.Terraform", "AcceptanceTests"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_remoteDebugging(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_remoteDebugging(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.remote_debugging_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.remote_debugging_version", "VS2015"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsDotNet2(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_windowsDotNet(ri, testLocation(), "v2.0")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.dotnet_framework_version", "v2.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsDotNet4(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_windowsDotNet(ri, testLocation(), "v4.0")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.dotnet_framework_version", "v4.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsDotNetUpdate(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_windowsDotNet(ri, testLocation(), "v2.0")
	updatedConfig := testAccAzureRMAppServiceSlot_windowsDotNet(ri, testLocation(), "v4.0")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.dotnet_framework_version", "v2.0"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.dotnet_framework_version", "v4.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsJava7Jetty(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_windowsJava(ri, testLocation(), "1.7", "JETTY", "9.3")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_version", "1.7"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container", "JETTY"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container_version", "9.3"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsJava8Jetty(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_windowsJava(ri, testLocation(), "1.8", "JETTY", "9.3")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_version", "1.8"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container", "JETTY"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container_version", "9.3"),
				),
			},
		},
	})
}
func TestAccAzureRMAppServiceSlot_windowsJava7Tomcat(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_windowsJava(ri, testLocation(), "1.7", "TOMCAT", "9.0")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_version", "1.7"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container", "TOMCAT"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container_version", "9.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsJava8Tomcat(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_windowsJava(ri, testLocation(), "1.8", "TOMCAT", "9.0")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_version", "1.8"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container", "TOMCAT"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container_version", "9.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsPHP7(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_windowsPHP(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.php_version", "7.1"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsPython(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_windowsPython(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.python_version", "3.4"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_webSockets(t *testing.T) {
	resourceName := "azurerm_app_service_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceSlot_webSockets(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.websockets_enabled", "true"),
				),
			},
		},
	})
}

func testCheckAzureRMAppServiceSlotDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).appServicesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_slot" {
			continue
		}

		slot := rs.Primary.Attributes["name"]
		appServiceName := rs.Primary.Attributes["app_service_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.GetSlot(ctx, resourceGroup, appServiceName, slot)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMAppServiceSlotExists(slot string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[slot]
		if !ok {
			return fmt.Errorf("Slot Not found: %q", slot)
		}

		appServiceName := rs.Primary.Attributes["app_service_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Service Slot: %q/%q", appServiceName, slot)
		}

		client := testAccProvider.Meta().(*ArmClient).appServicesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.GetSlot(ctx, resourceGroup, appServiceName, slot)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: App Service slot %q/%q (resource group: %q) does not exist", appServiceName, slot, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on appServicesClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMAppServiceSlot_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceSlot_32Bit(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"

  site_config {
    use_32_bit_worker_process = true
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceSlot_alwaysOn(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"

  site_config {
    always_on = true
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceSlot_appSettings(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"

  app_settings {
    "foo" = "bar"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceSlot_clientAffinityEnabled(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                    = "acctestASSlot-%d"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  app_service_plan_id     = "${azurerm_app_service_plan.test.id}"
  app_service_name        = "${azurerm_app_service.test.name}"
  client_affinity_enabled = true
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceSlot_connectionStrings(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"

  connection_string {
    name  = "Example"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceSlot_defaultDocuments(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"

  site_config {
    default_documents = [
      "first.html",
      "second.jsp",
      "third.aspx",
    ]
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceSlot_enabled(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"
  enabled = false
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceSlot_httpsOnly(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"
  https_only          = true
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceSlot_http2Enabled(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"

  site_config {
    http2_enabled = true
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceSlot_oneIpRestriction(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10"
    }
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceSlot_manyIpRestrictions(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10"
    }

    ip_restriction {
      ip_address  = "20.20.20.0"
      subnet_mask = "255.255.255.0"
    }

    ip_restriction {
      ip_address  = "30.30.0.0"
      subnet_mask = "255.255.0.0"
    }
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceSlot_localMySql(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"

  site_config {
    local_mysql_enabled = true
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceSlot_managedPipelineMode(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"

  site_config {
    managed_pipeline_mode = "Classic"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceSlot_tags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"

  tags {
    "Hello" = "World"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceSlot_tagsUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"

  tags {
    "Hello"     = "World"
    "Terraform" = "AcceptanceTests"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceSlot_remoteDebugging(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"

  site_config {
    remote_debugging_enabled = true
    remote_debugging_version = "VS2015"
  }

  tags {
    "Hello" = "World"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceSlot_windowsDotNet(rInt int, location, version string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"

  site_config {
    dotnet_framework_version = "%s"
  }
}
`, rInt, location, rInt, rInt, rInt, version)
}

func testAccAzureRMAppServiceSlot_windowsJava(rInt int, location, javaVersion, container, containerVersion string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"

  site_config {
    java_version           = "%s"
    java_container         = "%s"
    java_container_version = "%s"
  }
}
`, rInt, location, rInt, rInt, rInt, javaVersion, container, containerVersion)
}

func testAccAzureRMAppServiceSlot_windowsPHP(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"

  site_config {
    php_version = "7.1"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceSlot_windowsPython(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"

  site_config {
    python_version = "3.4"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceSlot_webSockets(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  app_service_name    = "${azurerm_app_service.test.name}"

  site_config {
    websockets_enabled = true
  }
}
`, rInt, location, rInt, rInt, rInt)
}
