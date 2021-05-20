package web_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AppServiceEnvironmentResource struct{}

func TestAccAppServiceEnvironment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")
	r := AppServiceEnvironmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("pricing_tier").HasValue("I1"),
				check.That(data.ResourceName).Key("front_end_scale_factor").HasValue("15"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServiceEnvironment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")
	r := AppServiceEnvironmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAppServiceEnvironment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")
	r := AppServiceEnvironmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("pricing_tier").HasValue("I1"),
				check.That(data.ResourceName).Key("front_end_scale_factor").HasValue("15"),
			),
		},
		data.ImportStep(),
		{
			Config: r.tierAndScaleFactor(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("pricing_tier").HasValue("I2"),
				check.That(data.ResourceName).Key("front_end_scale_factor").HasValue("10"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServiceEnvironment_tierAndScaleFactor(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")
	r := AppServiceEnvironmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.tierAndScaleFactor(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("pricing_tier").HasValue("I2"),
				check.That(data.ResourceName).Key("front_end_scale_factor").HasValue("10"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServiceEnvironment_withAppServicePlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")
	aspData := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")
	r := AppServiceEnvironmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withAppServicePlan(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").MatchesOtherKey(
					check.That(aspData.ResourceName).Key("app_service_environment_id"),
				),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServiceEnvironment_dedicatedResourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")
	r := AppServiceEnvironmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dedicatedResourceGroup(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServiceEnvironment_withCertificatePfx(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")
	certData := acceptance.BuildTestData(t, "azurerm_app_service_certificate", "test")
	r := AppServiceEnvironmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withCertificatePfx(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").MatchesOtherKey(
					check.That(certData.ResourceName).Key("hosting_environment_profile_id"),
				),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServiceEnvironment_internalLoadBalancer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")
	r := AppServiceEnvironmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.internalLoadBalancerAndWhitelistedIpRanges(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("internal_load_balancing_mode").HasValue("Web, Publishing"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServiceEnvironment_clusterSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")
	r := AppServiceEnvironmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.clusterSettings(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cluster_setting.#").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.clusterSettingsUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cluster_setting.#").HasValue("3"),
			),
		},
		data.ImportStep(),
		{
			Config: r.clusterSettings(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cluster_setting.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func (r AppServiceEnvironmentResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.AppServiceEnvironmentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Web.AppServiceEnvironmentsClient.Get(ctx, id.ResourceGroup, id.HostingEnvironmentName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving App Service Environment %q (Resource Group %q): %+v", id.HostingEnvironmentName, id.ResourceGroup, err)
	}

	return utils.Bool(true), nil
}

func (r AppServiceEnvironmentResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_environment" "test" {
  name      = "acctest-ase-%d"
  subnet_id = azurerm_subnet.ase.id
}
`, template, data.RandomInteger)
}

func (r AppServiceEnvironmentResource) clusterSettings(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_environment" "test" {
  name      = "acctest-ase-%d"
  subnet_id = azurerm_subnet.ase.id

  cluster_setting {
    name  = "InternalEncryption"
    value = "true"
  }

  cluster_setting {
    name  = "DisableTls1.0"
    value = "1"
  }
}
`, template, data.RandomInteger)
}

func (r AppServiceEnvironmentResource) clusterSettingsUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_environment" "test" {
  name      = "acctest-ase-%d"
  subnet_id = azurerm_subnet.ase.id

  cluster_setting {
    name  = "InternalEncryption"
    value = "true"
  }

  cluster_setting {
    name  = "DisableTls1.0"
    value = "1"
  }

  cluster_setting {
    name  = "FrontEndSSLCipherSuiteOrder"
    value = "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384_P256,TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256_P256,TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384_P256,TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256_P256,TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA_P256,TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA_P256"
  }
}
`, template, data.RandomInteger)
}

func (r AppServiceEnvironmentResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_environment" "import" {
  name      = azurerm_app_service_environment.test.name
  subnet_id = azurerm_app_service_environment.test.subnet_id
}
`, template)
}

func (r AppServiceEnvironmentResource) tierAndScaleFactor(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_environment" "test" {
  name                   = "acctest-ase-%d"
  subnet_id              = azurerm_subnet.ase.id
  pricing_tier           = "I2"
  front_end_scale_factor = 10
}
`, template, data.RandomInteger)
}

func (r AppServiceEnvironmentResource) withAppServicePlan(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_plan" "test" {
  name                       = "acctest-ASP-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_environment_id = azurerm_app_service_environment.test.id

  sku {
    tier     = "Isolated"
    size     = "I1"
    capacity = 1
  }
}
`, template, data.RandomInteger)
}

func (r AppServiceEnvironmentResource) dedicatedResourceGroup(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-%[2]d-2"
  location = "%s"
}

resource "azurerm_app_service_environment" "test" {
  name                = "acctest-ase-%[2]d"
  resource_group_name = azurerm_resource_group.test2.name
  subnet_id           = azurerm_subnet.ase.id
}
`, template, data.RandomInteger, data.Locations.Secondary)
}

func (r AppServiceEnvironmentResource) withCertificatePfx(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_certificate" "test" {
  name                           = "acctest-cert-%d"
  resource_group_name            = azurerm_app_service_environment.test.resource_group_name
  location                       = azurerm_resource_group.test.location
  pfx_blob                       = filebase64("testdata/app_service_certificate.pfx")
  password                       = "terraform"
  hosting_environment_profile_id = azurerm_app_service_environment.test.id
}
`, template, data.RandomInteger)
}

func (r AppServiceEnvironmentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "ase" {
  name                 = "asesubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_subnet" "gateway" {
  name                 = "gatewaysubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r AppServiceEnvironmentResource) internalLoadBalancerAndWhitelistedIpRanges(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_environment" "test" {
  name                         = "acctest-ase-%d"
  subnet_id                    = azurerm_subnet.ase.id
  pricing_tier                 = "I1"
  front_end_scale_factor       = 5
  internal_load_balancing_mode = "Web, Publishing"
  allowed_user_ip_cidrs        = ["11.22.33.44/32", "55.66.77.0/24"]
}
`, template, data.RandomInteger)
}
