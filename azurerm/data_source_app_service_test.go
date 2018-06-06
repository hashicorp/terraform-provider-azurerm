package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMAppService_basic(t *testing.T) {
	dataSourceName := "data.azurerm_app_service.test"
	rInt := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "app_service_plan_id"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_tags(t *testing.T) {
	dataSourceName := "data.azurerm_app_service.test"
	rInt := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_tags(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.Hello", "World"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_clientAppAffinityDisabled(t *testing.T) {
	dataSourceName := "data.azurerm_app_service.test"
	rInt := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_clientAffinityDisabled(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "client_affinity_enabled", "false"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_32Bit(t *testing.T) {
	dataSourceName := "data.azurerm_app_service.test"
	rInt := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_32Bit(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "site_config.0.use_32_bit_worker_process", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_appSettings(t *testing.T) {
	dataSourceName := "data.azurerm_app_service.test"
	rInt := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_appSettings(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "app_settings.foo", "bar"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_connectionString(t *testing.T) {
	dataSourceName := "data.azurerm_app_service.test"
	rInt := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_connectionStrings(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "connection_string.0.name", "Example"),
					resource.TestCheckResourceAttr(dataSourceName, "connection_string.0.value", "some-postgresql-connection-string"),
					resource.TestCheckResourceAttr(dataSourceName, "connection_string.0.type", "PostgreSQL"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_ipRestriction(t *testing.T) {
	dataSourceName := "data.azurerm_app_service.test"
	rInt := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_ipRestriction(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "ip_restriction.0.ip_address", "10.10.10.10"),
					resource.TestCheckResourceAttr(dataSourceName, "ip_restriction.0.subnet_mask", "255.255.255.255"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_http2Enabled(t *testing.T) {
	dataSourceName := "data.azurerm_app_service.test"
	rInt := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_http2Enabled(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "site_config.0.http2_enabled", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAppService_basic(rInt int, location string) string {
	config := testAccAzureRMAppService_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = "${azurerm_app_service.test.name}"
  resource_group_name = "${azurerm_app_service.test.resource_group_name}"
}
`, config)
}

func testAccDataSourceAppService_tags(rInt int, location string) string {
	config := testAccAzureRMAppService_tags(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = "${azurerm_app_service.test.name}"
  resource_group_name = "${azurerm_app_service.test.resource_group_name}"
}
`, config)
}

func testAccDataSourceAppService_clientAffinityDisabled(rInt int, location string) string {
	config := testAccAzureRMAppService_clientAffinityDisabled(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = "${azurerm_app_service.test.name}"
  resource_group_name = "${azurerm_app_service.test.resource_group_name}"
}
`, config)
}

func testAccDataSourceAppService_32Bit(rInt int, location string) string {
	config := testAccAzureRMAppService_32Bit(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = "${azurerm_app_service.test.name}"
  resource_group_name = "${azurerm_app_service.test.resource_group_name}"
}
`, config)
}

func testAccDataSourceAppService_appSettings(rInt int, location string) string {
	config := testAccAzureRMAppService_appSettings(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = "${azurerm_app_service.test.name}"
  resource_group_name = "${azurerm_app_service.test.resource_group_name}"
}
`, config)
}

func testAccDataSourceAppService_connectionStrings(rInt int, location string) string {
	config := testAccAzureRMAppService_connectionStrings(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = "${azurerm_app_service.test.name}"
  resource_group_name = "${azurerm_app_service.test.resource_group_name}"
}
`, config)
}

func testAccDataSourceAppService_ipRestriction(rInt int, location string) string {
	config := testAccAzureRMAppService_oneIpRestriction(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = "${azurerm_app_service.test.name}"
  resource_group_name = "${azurerm_app_service.test.resource_group_name}"
}
`, config)
}

func testAccDataSourceAppService_http2Enabled(rInt int, location string) string {
	config := testAccAzureRMAppService_http2Enabled(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = "${azurerm_app_service.test.name}"
  resource_group_name = "${azurerm_app_service.test.resource_group_name}"
}
`, config)
}
