package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAppConfigurationName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "four",
			ErrCount: 1,
		},
		{
			Value:    "5five",
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
			Value:    "helloWorld",
			ErrCount: 0,
		},
		{
			Value:    "helloworld12",
			ErrCount: 0,
		},
		{
			Value:    "hello@world",
			ErrCount: 1,
		},
		{
			Value:    "qfvbdsbvipqdbwsbddbdcwqffewsqwcdw21ddwqwd3324120",
			ErrCount: 0,
		},
		{
			Value:    "qfvbdsbvipqdbwsbddbdcwqffewsqwcdw21ddwqwd332412020",
			ErrCount: 0,
		},
		{
			Value:    "qfvbdsbvipqdbwsbddbdcwqfjjfewsqwcdw21ddwqwd33241201",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateAppConfigurationName(tc.Value, "azurerm_app_configuration")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure App Configuration Name to trigger a validation error: %v", tc)
		}
	}
}

func TestAccAzureAppConfiguration_free(t *testing.T) {
	rn := "azurerm_app_configuration.test"
	ri := tf.AccRandTimeInt()
	l := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureAppConfiguration_free(ri, l),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureAppConfigurationExists(rn),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureAppConfiguration_standard(t *testing.T) {
	rn := "azurerm_app_configuration.test"
	ri := tf.AccRandTimeInt()
	l := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureAppConfiguration_standard(ri, l),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureAppConfigurationExists(rn),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureAppConfiguration_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	rn := "azurerm_app_configuration.test"
	ri := tf.AccRandTimeInt()
	l := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureAppConfiguration_free(ri, l),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureAppConfigurationExists(rn),
				),
			},
			{
				Config:      testAccAzureAppConfiguration_requiresImport(ri, l),
				ExpectError: testRequiresImportError("azurerm_app_configuration"),
			},
		},
	})
}

func TestAccAzureAppConfiguration_complete(t *testing.T) {
	rn := "azurerm_app_configuration.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureAppConfiguration_complete(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureAppConfigurationExists(rn),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureAppConfiguration_update(t *testing.T) {
	rn := "azurerm_app_configuration.test"
	ri := tf.AccRandTimeInt()
	l := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureAppConfiguration_complete(ri, l),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureAppConfigurationExists(rn),
				),
			},
			{
				Config: testAccAzureAppConfiguration_completeUpdated(ri, l),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureAppConfigurationExists(rn),
				),
			},
		},
	})
}

func testCheckAzureAppConfigurationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).AppConfiguration.AppConfigurationsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_configuration" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}

func testCheckAzureAppConfigurationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Configuration: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).AppConfiguration.AppConfigurationsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on appConfigurationsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: App Configuration %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testAccAzureAppConfiguration_free(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_app_configuration" "test" {
  name                = "testaccappconf%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "free"
}
`, rInt, location, rInt)
}

func testAccAzureAppConfiguration_standard(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_app_configuration" "test" {
  name                = "testaccappconf%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "standard"
}
`, rInt, location, rInt)
}

func testAccAzureAppConfiguration_requiresImport(rInt int, location string) string {
	template := testAccAzureAppConfiguration_free(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration" "import" {
  name                = "${azurerm_app_configuration.test.name}"
  resource_group_name = "${azurerm_app_configuration.test.resource_group_name}"
  location            = "${azurerm_app_configuration.test.location}"
  sku                 = "${azurerm_app_configuration.test.sku}"

}
`, template)
}

func testAccAzureAppConfiguration_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_app_configuration" "test" {
  name                = "testaccappconf%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "free"

  tags = {
    environment = "development"
  }
}
`, rInt, location, rInt)
}

func testAccAzureAppConfiguration_completeUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_app_configuration" "test" {
  name                = "testaccappconf%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "free"

  tags = {
    environment = "production"
  }
}
`, rInt, location, rInt)
}
