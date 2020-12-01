package web_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type AppServiceDataSource struct{}

func TestAccDataSourceAzureRMAppService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppServiceDataSource{}.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("app_service_plan_id").Exists(),
				check.That(data.ResourceName).Key("outbound_ip_addresses").Exists(),
				check.That(data.ResourceName).Key("possible_outbound_ip_addresses").Exists(),
				check.That(data.ResourceName).Key("custom_domain_verification_id").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourceAzureRMAppService_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppServiceDataSource{}.tags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Hello").HasValue("World"),
			),
		},
	})
}

func TestAccDataSourceAzureRMAppService_clientAppAffinityDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppServiceDataSource{}.clientAffinityDisabled(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("client_affinity_enabled").HasValue("false"),
			),
		},
	})
}

func TestAccDataSourceAzureRMAppService_32Bit(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppServiceDataSource{}.service32Bit(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("site_config.0.use_32_bit_worker_process").HasValue("true"),
			),
		},
	})
}

func TestAccDataSourceAzureRMAppService_appSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppServiceDataSource{}.appSettings(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("app_settings.foo").HasValue("bar"),
			),
		},
	})
}

func TestAccDataSourceAzureRMAppService_connectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppServiceDataSource{}.connectionStrings(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("connection_string.0.name").HasValue("First"),
				check.That(data.ResourceName).Key("connection_string.0.value").HasValue("first-connection-string"),
				check.That(data.ResourceName).Key("connection_string.0.type").HasValue("Custom"),
				check.That(data.ResourceName).Key("connection_string.1.name").HasValue("Second"),
				check.That(data.ResourceName).Key("connection_string.1.value").HasValue("some-postgresql-connection-string"),
				check.That(data.ResourceName).Key("connection_string.1.type").HasValue("PostgreSQL"),
			),
		},
	})
}

func TestAccDataSourceAzureRMAppService_ipRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppServiceDataSource{}.ipRestriction(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.ip_address").HasValue("10.10.10.10/32"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.name").HasValue("test-restriction"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.priority").HasValue("123"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.action").HasValue("Allow"),
			),
		},
	})
}

func TestAccDataSourceAzureRMAppService_oneVNetSubnetIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppServiceDataSource{}.oneVNetSubnetIpRestriction(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.#").HasValue("1"),
			),
		},
	})
}

func TestAccDataSourceAzureRMAppService_scmUseMainIPRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppServiceDataSource{}.scmUseMainIPRestriction(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("site_config.0.scm_use_main_ip_restriction").HasValue("true"),
			),
		},
	})
}

func TestAccDataSourceAzureRMAppService_scmIPRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppServiceDataSource{}.scmIPRestriction(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("site_config.0.scm_ip_restriction.0.ip_address").HasValue("10.10.10.10/32"),
				check.That(data.ResourceName).Key("site_config.0.scm_ip_restriction.0.name").HasValue("test-restriction"),
				check.That(data.ResourceName).Key("site_config.0.scm_ip_restriction.0.priority").HasValue("123"),
				check.That(data.ResourceName).Key("site_config.0.scm_ip_restriction.0.action").HasValue("Allow"),
			),
		},
	})
}

func TestAccDataSourceAzureRMAppService_withSourceControl(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppServiceDataSource{}.withSourceControl(data, "main"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("source_control.0.branch").HasValue("main"),
			),
		},
	})
}

func TestAccDataSourceAzureRMAppService_http2Enabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppServiceDataSource{}.http2Enabled(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("site_config.0.http2_enabled").HasValue("true"),
			),
		},
	})
}

func TestAccDataSourceAzureRMAppService_minTls(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppServiceDataSource{}.minTls(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("site_config.0.min_tls_version").HasValue("1.1"),
			),
		},
	})
}

func TestAccDataSourceAzureRMAppService_basicWindowsContainer(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppServiceDataSource{}.basicWindowsContainer(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("site_config.0.windows_fx_version").HasValue("DOCKER|mcr.microsoft.com/azure-app-service/samples/aspnethelloworld:latest"),
				check.That(data.ResourceName).Key("app_settings.DOCKER_REGISTRY_SERVER_URL").HasValue("https://mcr.microsoft.com"),
			),
		},
	})
}

func (d AppServiceDataSource) basic(data acceptance.TestData) string {
	config := AppServiceResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func (d AppServiceDataSource) tags(data acceptance.TestData) string {
	config := AppServiceResource{}.tags(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func (d AppServiceDataSource) clientAffinityDisabled(data acceptance.TestData) string {
	config := AppServiceResource{}.clientAffinityDisabled(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func (d AppServiceDataSource) service32Bit(data acceptance.TestData) string {
	config := AppServiceResource{}.service32Bit(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func (d AppServiceDataSource) appSettings(data acceptance.TestData) string {
	config := AppServiceResource{}.appSettings(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func (d AppServiceDataSource) connectionStrings(data acceptance.TestData) string {
	config := AppServiceResource{}.connectionStrings(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func (d AppServiceDataSource) ipRestriction(data acceptance.TestData) string {
	config := AppServiceResource{}.completeIpRestriction(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func (d AppServiceDataSource) oneVNetSubnetIpRestriction(data acceptance.TestData) string {
	config := AppServiceResource{}.oneVNetSubnetIpRestriction(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func (d AppServiceDataSource) scmUseMainIPRestriction(data acceptance.TestData) string {
	config := AppServiceResource{}.scmUseMainIPRestriction(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func (d AppServiceDataSource) scmIPRestriction(data acceptance.TestData) string {
	config := AppServiceResource{}.completeScmIpRestriction(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func (d AppServiceDataSource) withSourceControl(data acceptance.TestData, branch string) string {
	config := AppServiceResource{}.withSourceControl(data, branch)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func (d AppServiceDataSource) http2Enabled(data acceptance.TestData) string {
	config := AppServiceResource{}.http2Enabled(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func (d AppServiceDataSource) minTls(data acceptance.TestData) string {
	config := AppServiceResource{}.minTls(data, "1.1")
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}

func (d AppServiceDataSource) basicWindowsContainer(data acceptance.TestData) string {
	config := AppServiceResource{}.basicWindowsContainer(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service" "test" {
  name                = azurerm_app_service.test.name
  resource_group_name = azurerm_app_service.test.resource_group_name
}
`, config)
}
