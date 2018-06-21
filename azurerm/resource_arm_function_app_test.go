package azurerm

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMFunctionApp_basic(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := acctest.RandInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMFunctionApp_basic(ri, rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "version", "~1"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_tags(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := acctest.RandInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMFunctionApp_tags(ri, rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "production"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_tagsUpdate(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := acctest.RandInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMFunctionApp_tags(ri, rs, testLocation())
	updatedConfig := testAccAzureRMFunctionApp_tagsUpdated(ri, rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "production"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "production"),
					resource.TestCheckResourceAttr(resourceName, "tags.hello", "Berlin"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_appSettings(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := acctest.RandInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMFunctionApp_basic(ri, rs, testLocation())
	updatedConfig := testAccAzureRMFunctionApp_appSettings(ri, rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app_settings.%", "0"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app_settings.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_settings.hello", "world"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app_settings.%", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_siteConfig(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := acctest.RandInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMFunctionApp_alwaysOn(ri, rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.always_on", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_connectionStrings(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := acctest.RandInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMFunctionApp_connectionStrings(ri, rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "connection_string.0.name", "Example"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.0.value", "some-postgresql-connection-string"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.0.type", "PostgreSQL"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_siteConfigMulti(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := acctest.RandInt()
	rs := strings.ToLower(acctest.RandString(11))
	configBase := testAccAzureRMFunctionApp_basic(ri, rs, testLocation())
	configUpdate1 := testAccAzureRMFunctionApp_appSettings(ri, rs, testLocation())
	configUpdate2 := testAccAzureRMFunctionApp_appSettingsAlwaysOn(ri, rs, testLocation())
	configUpdate3 := testAccAzureRMFunctionApp_appSettingsAlwaysOnConnectionStrings(ri, rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBase,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app_settings.%", "0"),
				),
			},
			{
				Config: configUpdate1,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app_settings.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_settings.hello", "world"),
				),
			},
			{
				Config: configUpdate2,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app_settings.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_settings.hello", "world"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.always_on", "true"),
				),
			},
			{
				Config: configUpdate3,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app_settings.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_settings.hello", "world"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.always_on", "true"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.0.name", "Example"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.0.value", "some-postgresql-connection-string"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.0.type", "PostgreSQL"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_updateVersion(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := acctest.RandInt()
	rs := strings.ToLower(acctest.RandString(11))
	preConfig := testAccAzureRMFunctionApp_version(ri, rs, testLocation(), "beta")
	postConfig := testAccAzureRMFunctionApp_version(ri, rs, testLocation(), "~1")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "version", "beta"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "version", "~1"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_3264bit(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := acctest.RandInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	config := testAccAzureRMFunctionApp_basic(ri, rs, location)
	updatedConfig := testAccAzureRMFunctionApp_64bit(ri, rs, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.use_32_bit_worker_process", "true"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.use_32_bit_worker_process", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_httpsOnly(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := acctest.RandInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	config := testAccAzureRMFunctionApp_httpsOnly(ri, rs, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "https_only", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_consumptionPlan(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := acctest.RandInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	config := testAccAzureRMFunctionApp_consumptionPlan(ri, rs, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.use_32_bit_worker_process", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_createIdentity(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := acctest.RandInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMFunctionApp_basicIdentity(ri, rs, testLocation())

	uuidMatch := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(resourceName, "identity.0.principal_id", uuidMatch),
					resource.TestMatchResourceAttr(resourceName, "identity.0.tenant_id", uuidMatch),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_updateIdentity(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := acctest.RandInt()
	rs := strings.ToLower(acctest.RandString(11))

	preConfig := testAccAzureRMFunctionApp_basic(ri, rs, testLocation())
	postConfig := testAccAzureRMFunctionApp_basicIdentity(ri, rs, testLocation())

	uuidMatch := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.#", "0"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(resourceName, "identity.0.principal_id", uuidMatch),
					resource.TestMatchResourceAttr(resourceName, "identity.0.tenant_id", uuidMatch),
				),
			},
		},
	})
}

func testCheckAzureRMFunctionAppDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).appServicesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_function_app" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)

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

func testCheckAzureRMFunctionAppExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		functionAppName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Function App: %s", functionAppName)
		}

		client := testAccProvider.Meta().(*ArmClient).appServicesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, functionAppName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Function App %q (resource group: %q) does not exist", functionAppName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on appServicesClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMFunctionApp_basic(rInt int, storage string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%[1]d"
	location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
	name                     = "acctestsa%[3]s"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
	name                = "acctestASP-%[1]d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	sku {
		tier = "Standard"
		size = "S1"
	}
}

resource "azurerm_function_app" "test" {
	name                      = "acctest-%[1]d-func"
	location                  = "${azurerm_resource_group.test.location}"
	resource_group_name       = "${azurerm_resource_group.test.name}"
	app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
	storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
}`, rInt, location, storage)
}

func testAccAzureRMFunctionApp_tags(rInt int, storage string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%[1]d"
	location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
	name                     = "acctestsa%[3]s"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
	name                = "acctestASP-%[1]d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	sku {
		tier = "Standard"
		size = "S1"
	}
}

resource "azurerm_function_app" "test" {
	name                      = "acctest-%[1]d-func"
	location                  = "${azurerm_resource_group.test.location}"
	resource_group_name       = "${azurerm_resource_group.test.name}"
	app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
	storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
	tags {
		environment = "production"
	}
}
`, rInt, location, storage)
}

func testAccAzureRMFunctionApp_tagsUpdated(rInt int, storage string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%[1]d"
	location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
	name                     = "acctestsa%[3]s"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
	name                = "acctestASP-%[1]d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	sku {
		tier = "Standard"
		size = "S1"
	}
}

resource "azurerm_function_app" "test" {
	name                      = "acctest-%[1]d-func"
	location                  = "${azurerm_resource_group.test.location}"
	resource_group_name       = "${azurerm_resource_group.test.name}"
	app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
	storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
	tags {
		environment = "production"
		hello       = "Berlin"
	}
}
`, rInt, location, storage)
}

func testAccAzureRMFunctionApp_version(rInt int, storage string, location string, version string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%[1]d"
	location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
	name                     = "acctestsa%[3]s"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
	name                = "acctestASP-%[1]d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	sku {
		tier = "Standard"
		size = "S1"
	}
}

resource "azurerm_function_app" "test" {
	name                      = "acctest-%[1]d-func"
	location                  = "${azurerm_resource_group.test.location}"
	resource_group_name       = "${azurerm_resource_group.test.name}"
	app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
	version                   = "%[4]s"
	storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
}`, rInt, location, storage, version)
}

func testAccAzureRMFunctionApp_appSettings(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%[1]d"
	location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
	name                     = "acctestsa%[3]s"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
	name                = "acctestASP-%[1]d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	sku {
		tier = "Standard"
		size = "S1"
	}
}

resource "azurerm_function_app" "test" {
	name                      = "acctest-%[1]d-func"
	location                  = "${azurerm_resource_group.test.location}"
	resource_group_name       = "${azurerm_resource_group.test.name}"
	app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
	storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
	app_settings {
		"hello" = "world"
	}
}
`, rInt, location, rString)
}

func testAccAzureRMFunctionApp_alwaysOn(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%[1]d"
	location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
	name                     = "acctestsa%[3]s"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
	name                = "acctestASP-%[1]d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	sku {
		tier = "Standard"
		size = "S1"
	}
}

resource "azurerm_function_app" "test" {
	name                      = "acctest-%[1]d-func"
	location                  = "${azurerm_resource_group.test.location}"
	resource_group_name       = "${azurerm_resource_group.test.name}"
	app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
	storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
	site_config {
		always_on = true
	}
}
`, rInt, location, rString)
}

func testAccAzureRMFunctionApp_connectionStrings(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%[1]d"
	location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
	name                     = "acctestsa%[3]s"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
	name                = "acctestASP-%[1]d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	sku {
		tier = "Standard"
		size = "S1"
	}
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
	location                  = "${azurerm_resource_group.test.location}"
	resource_group_name       = "${azurerm_resource_group.test.name}"
	app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
	storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  connection_string {
  	name  = "Example"
  	value = "some-postgresql-connection-string"
  	type  = "PostgreSQL"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMFunctionApp_appSettingsAlwaysOn(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%[1]d"
	location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
	name                     = "acctestsa%[3]s"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
	name                = "acctestASP-%[1]d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	sku {
		tier = "Standard"
		size = "S1"
	}
}

resource "azurerm_function_app" "test" {
	name                      = "acctest-%[1]d-func"
	location                  = "${azurerm_resource_group.test.location}"
	resource_group_name       = "${azurerm_resource_group.test.name}"
	app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
	storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
	app_settings {
		"hello" = "world"
	}
	site_config {
		always_on = true
	}
}
`, rInt, location, rString)
}

func testAccAzureRMFunctionApp_appSettingsAlwaysOnConnectionStrings(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%[1]d"
	location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
	name                     = "acctestsa%[3]s"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
	name                = "acctestASP-%[1]d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	sku {
		tier = "Standard"
		size = "S1"
	}
}

resource "azurerm_function_app" "test" {
	name                      = "acctest-%[1]d-func"
	location                  = "${azurerm_resource_group.test.location}"
	resource_group_name       = "${azurerm_resource_group.test.name}"
	app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
	storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
	app_settings {
		"hello" = "world"
	}
	site_config {
		always_on = true
	}
	connection_string {
  	name  = "Example"
  	value = "some-postgresql-connection-string"
  	type  = "PostgreSQL"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMFunctionApp_64bit(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
	name                     = "acctestsa%s"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"
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

resource "azurerm_function_app" "test" {
  name                      = "acctest-%d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  site_config {
    use_32_bit_worker_process = false
  }
}
`, rInt, location, rString, rInt, rInt)
}

func testAccAzureRMFunctionApp_httpsOnly(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
	name                     = "acctestsa%s"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"
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

resource "azurerm_function_app" "test" {
  name                      = "acctest-%d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
  https_only                = true
}
`, rInt, location, rString, rInt, rInt)
}

func testAccAzureRMFunctionApp_consumptionPlan(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
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

resource "azurerm_function_app" "test" {
  name                      = "acctest-%d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
}
`, rInt, location, rString, rInt, rInt)
}

func testAccAzureRMFunctionApp_basicIdentity(rInt int, storage string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  identity {
    type = "SystemAssigned"
  }
}`, rInt, location, storage)
}
