package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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
			ErrCount: 0,
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
	resourceName := "azurerm_app_service_plan.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_basicWindows(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "per_site_scaling", "false"),
					resource.TestCheckResourceAttr(resourceName, "reserved", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_basicLinux(t *testing.T) {
	resourceName := "azurerm_app_service_plan.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_basicLinux(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMAppServicePlan_basicLinuxNew(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "per_site_scaling", "false"),
					resource.TestCheckResourceAttr(resourceName, "reserved", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_app_service_plan.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_basicLinux(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMAppServicePlan_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_app_service_plan"),
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_standardWindows(t *testing.T) {
	resourceName := "azurerm_app_service_plan.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_standardWindows(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_premiumWindows(t *testing.T) {
	resourceName := "azurerm_app_service_plan.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_premiumWindows(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_premiumWindowsUpdated(t *testing.T) {
	resourceName := "azurerm_app_service_plan.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_premiumWindows(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
				),
			},
			{
				Config: testAccAzureRMAppServicePlan_premiumWindowsUpdated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_completeWindows(t *testing.T) {
	resourceName := "azurerm_app_service_plan.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_completeWindows(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "properties.0.per_site_scaling", "true"),
					resource.TestCheckResourceAttr(resourceName, "properties.0.reserved", "false"),
				),
			},
			{
				Config: testAccAzureRMAppServicePlan_completeWindowsNew(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "per_site_scaling", "true"),
					resource.TestCheckResourceAttr(resourceName, "reserved", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_consumptionPlan(t *testing.T) {
	resourceName := "azurerm_app_service_plan.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_consumptionPlan(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "Dynamic"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.size", "Y1"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_premiumConsumptionPlan(t *testing.T) {
	resourceName := "azurerm_app_service_plan.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_premiumConsumptionPlan(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "ElasticPremium"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.size", "EP1"),
					resource.TestCheckResourceAttr(resourceName, "maximum_elastic_worker_count", "20"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_basicWindowsContainer(t *testing.T) {
	resourceName := "azurerm_app_service_plan.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_basicWindowsContainer(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "kind", "xenon"),
					resource.TestCheckResourceAttr(resourceName, "is_xenon", "true"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "PremiumContainer"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.size", "PC2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMAppServicePlanDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicePlansClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_plan" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := conn.Get(ctx, resourceGroup, name)

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

func testCheckAzureRMAppServicePlanExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		appServicePlanName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Service Plan: %s", appServicePlanName)
		}

		conn := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicePlansClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := conn.Get(ctx, resourceGroup, appServicePlanName)
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

  reserved = true
}
`, rInt, location, rInt)
}

func testAccAzureRMAppServicePlan_requiresImport(rInt int, location string) string {
	template := testAccAzureRMAppServicePlan_basicLinux(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_plan" "import" {
  name                = "${azurerm_app_service_plan.test.name}"
  location            = "${azurerm_app_service_plan.test.location}"
  resource_group_name = "${azurerm_app_service_plan.test.resource_group_name}"
  kind                = "${azurerm_app_service_plan.test.kind}"

  sku {
    tier = "Basic"
    size = "B1"
  }

  reserved = true
}
`, template)
}

func testAccAzureRMAppServicePlan_basicLinuxNew(rInt int, location string) string {
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

  reserved = true
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

  per_site_scaling = true
  reserved         = false

  tags = {
    environment = "Test"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMAppServicePlan_completeWindowsNew(rInt int, location string) string {
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

  per_site_scaling = true
  reserved         = false

  tags = {
    environment = "Test"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMAppServicePlan_consumptionPlan(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  kind                = "FunctionApp"

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMAppServicePlan_premiumConsumptionPlan(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  kind                = "elastic"

  maximum_elastic_worker_count = 20

  sku {
    tier = "ElasticPremium"
    size = "EP1"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMAppServicePlan_basicWindowsContainer(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "xenon"
  is_xenon            = true

  sku {
    tier = "PremiumContainer"
    size = "PC2"
  }
}
`, rInt, location, rInt)
}
