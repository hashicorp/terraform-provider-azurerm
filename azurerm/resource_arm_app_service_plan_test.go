package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAzureRMAppServicePlanName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "ab",
			ErrCount: 0,
		},
		{
			Value:    "abc",
			ErrCount: 0,
		},
		{
			Value:    "webapp1",
			ErrCount: 0,
		},
		{
			Value:    "hello-world",
			ErrCount: 0,
		},
		{
			Value:    "hello_world",
			ErrCount: 1,
		},
		{
			Value:    "helloworld21!",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateAppServicePlanName(tc.Value, "azurerm_app_service_plan")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the App Service Plan Name to trigger a validation error for '%s'", tc.Value)
		}
	}
}

func TestAccAzureRMAppServicePlan_basicWindows(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMAppServicePlan_basicWindows(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists("azurerm_app_service_plan.test"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_basicLinux(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMAppServicePlan_basicLinux(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists("azurerm_app_service_plan.test"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_standardWindows(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMAppServicePlan_standardWindows(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists("azurerm_app_service_plan.test"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_premiumWindows(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMAppServicePlan_premiumWindows(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists("azurerm_app_service_plan.test"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_premiumWindowsUpdated(t *testing.T) {
	resourceName := "azurerm_app_service_plan.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMAppServicePlan_premiumWindows(ri, location)
	updatedConfig := testAccAzureRMAppServicePlan_premiumWindowsUpdated(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_completeWindows(t *testing.T) {
	resourceName := "azurerm_app_service_plan.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServicePlan_completeWindows(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "properties.0.per_site_scaling", "true"),
					resource.TestCheckResourceAttr(resourceName, "properties.0.reserved", "false"),
				),
			},
		},
	})
}

func testCheckAzureRMAppServicePlanDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).appServicePlansClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_plan" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("App Service Plan still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMAppServicePlanExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		appServicePlanName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Service Plan: %s", appServicePlanName)
		}

		conn := testAccProvider.Meta().(*ArmClient).appServicePlansClient

		resp, err := conn.Get(resourceGroup, appServicePlanName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: App Service Plan %q (resource group: %q) does not exist", appServicePlanName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on appServicePlansClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMAppServicePlan_basicWindows(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Basic"
    size = "B1"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMAppServicePlan_basicLinux(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "Linux"

  sku {
    tier = "Basic"
    size = "B1"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMAppServicePlan_standardWindows(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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
`, rInt, location, rInt)
}

func testAccAzureRMAppServicePlan_premiumWindows(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Premium"
    size = "P1"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMAppServicePlan_premiumWindowsUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier     = "Premium"
    size     = "P1"
    capacity = 2
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMAppServicePlan_completeWindows(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "Windows"

  sku {
    tier = "Standard"
    size = "S1"
  }

  properties {
    per_site_scaling          = true
    reserved                  = false
  }

  tags {
    environment = "Test"
  }
}
`, rInt, location, rInt)
}
