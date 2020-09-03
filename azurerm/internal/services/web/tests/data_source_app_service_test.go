package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMAppService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "app_service_plan_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "outbound_ip_addresses"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "possible_outbound_ip_addresses"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_tags(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Hello", "World"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_clientAppAffinityDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_clientAffinityDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "client_affinity_enabled", "false"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_32Bit(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_32Bit(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.use_32_bit_worker_process", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_appSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_appSettings(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.foo", "bar"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_connectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_connectionStrings(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.0.name", "First"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.0.value", "first-connection-string"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.0.type", "Custom"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.1.name", "Second"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.1.value", "some-postgresql-connection-string"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.1.type", "PostgreSQL"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_ipRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_ipRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.ip_address", "10.10.10.10/32"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.name", "test-restriction"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.priority", "123"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.action", "Allow"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_oneVNetSubnetIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_oneVNetSubnetIpRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_scmUseMainIPRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_scmUseMainIPRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.scm_use_main_ip_restriction", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_scmIPRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_scmIPRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.scm_ip_restriction.0.ip_address", "10.10.10.10/32"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.scm_ip_restriction.0.name", "test-restriction"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.scm_ip_restriction.0.priority", "123"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.scm_ip_restriction.0.action", "Allow"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_withSourceControl(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_withSourceControl(data, "main"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "source_control.0.branch", "main"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_http2Enabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_http2Enabled(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.http2_enabled", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_minTls(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_minTls(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.min_tls_version", "1.1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppService_basicWindowsContainer(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_basicWindowsContainer(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.windows_fx_version", "DOCKER|mcr.microsoft.com/azure-app-service/samples/aspnethelloworld:latest"),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.DOCKER_REGISTRY_SERVER_URL", "https://mcr.microsoft.com"),
				),
			},
		},
	})
}

func testAccDataSourceAppService_basic(data acceptance.TestData) string {
	config := testAccAzureRMAppService_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func testAccDataSourceAppService_tags(data acceptance.TestData) string {
	config := testAccAzureRMAppService_tags(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func testAccDataSourceAppService_clientAffinityDisabled(data acceptance.TestData) string {
	config := testAccAzureRMAppService_clientAffinityDisabled(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func testAccDataSourceAppService_32Bit(data acceptance.TestData) string {
	config := testAccAzureRMAppService_32Bit(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func testAccDataSourceAppService_appSettings(data acceptance.TestData) string {
	config := testAccAzureRMAppService_appSettings(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func testAccDataSourceAppService_connectionStrings(data acceptance.TestData) string {
	config := testAccAzureRMAppService_connectionStrings(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func testAccDataSourceAppService_ipRestriction(data acceptance.TestData) string {
	config := testAccAzureRMAppService_completeIpRestriction(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func testAccDataSourceAppService_oneVNetSubnetIpRestriction(data acceptance.TestData) string {
	config := testAccAzureRMAppService_oneVNetSubnetIpRestriction(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func testAccDataSourceAppService_scmUseMainIPRestriction(data acceptance.TestData) string {
	config := testAccAzureRMAppService_scmUseMainIPRestriction(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func testAccDataSourceAppService_scmIPRestriction(data acceptance.TestData) string {
	config := testAccAzureRMAppService_completeScmIpRestriction(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func testAccDataSourceAppService_withSourceControl(data acceptance.TestData, branch string) string {
	config := testAccAzureRMAppService_withSourceControl(data, branch)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func testAccDataSourceAppService_http2Enabled(data acceptance.TestData) string {
	config := testAccAzureRMAppService_http2Enabled(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func testAccDataSourceAppService_minTls(data acceptance.TestData) string {
	config := testAccAzureRMAppService_minTls(data, "1.1")
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func testAccDataSourceAppService_basicWindowsContainer(data acceptance.TestData) string {
	config := testAccAzureRMAppService_basicWindowsContainer(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}
