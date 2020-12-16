package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type CdnEndpointResource struct{}

func TestAccCdnEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

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

func TestAccCdnEndpoint_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccCdnEndpoint_updateHostHeader(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.hostHeader(data, "www.contoso.com"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.hostHeader(data, ""),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.hostHeader(data, "www.example2.com"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpoint_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withTagsUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpoint_withoutCompression(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withoutCompression(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpoint_optimized(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.optimized(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpoint_withGeoFilters(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.geoFilters(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpoint_fullFields(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.fullFields(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("is_http_allowed").HasValue("true"),
				check.That(data.ResourceName).Key("is_https_allowed").HasValue("true"),
				check.That(data.ResourceName).Key("origin_path").HasValue("/origin-path"),
				check.That(data.ResourceName).Key("probe_path").HasValue("/origin-path/probe"),
				check.That(data.ResourceName).Key("origin_host_header").HasValue("www.contoso.com"),
				check.That(data.ResourceName).Key("optimization_type").HasValue("GeneralWebDelivery"),
				check.That(data.ResourceName).Key("querystring_caching_behaviour").HasValue("UseQueryString"),
				check.That(data.ResourceName).Key("content_types_to_compress.#").HasValue("1"),
				check.That(data.ResourceName).Key("is_compression_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("geo_filter.#").HasValue("1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
			),
		},
		// TODO -- add import step. Import step now gives us an error complaining that `is_compression_enabled` is not imported
	})
}

func TestAccCdnEndpoint_isHttpAndHttpsAllowedUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.isHttpAndHttpsAllowed(data, "true", "false"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.isHttpAndHttpsAllowed(data, "false", "true"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpoint_globalDeliveryRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.globalDeliveryRule(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.globalDeliveryRuleUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.globalDeliveryRuleRemove(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpoint_deliveryRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.deliveryRule(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.deliveryRuleUpdate1(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.deliveryRuleUpdate2(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.deliveryRuleRemove(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpoint_dnsAlias(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.testAccAzureRMCdnEndpoint_dnsAlias(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpoint_PremiumVerizon(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.PremiumVerizon(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CdnEndpointResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.EndpointID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Cdn.EndpointsClient.Get(ctx, id.ResourceGroup, id.ProfileName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving CDN Endpoint %q (Resource Group %q / Profile Name %q): %+v", id.Name, id.ResourceGroup, id.ProfileName, err)
	}
	return utils.Bool(true), nil
}

func (r CdnEndpointResource) Destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.EndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	endpointsClient := client.Cdn.EndpointsClient
	future, err := endpointsClient.Delete(ctx, id.ResourceGroup, id.ProfileName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("deleting CDN Endpoint %q (Resource Group %q / Profile %q): %+v", id.Name, id.ResourceGroup, id.ProfileName, err)
	}
	if err := future.WaitForCompletionRef(ctx, endpointsClient.Client); err != nil {
		return nil, fmt.Errorf("waiting for deletion of CDN Endpoint %q (Resource Group %q / Profile %q): %+v", id.Name, id.ResourceGroup, id.ProfileName, err)
	}

	return utils.Bool(true), nil
}

func (r CdnEndpointResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CdnEndpointResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_endpoint" "import" {
  name                = azurerm_cdn_endpoint.test.name
  profile_name        = azurerm_cdn_endpoint.test.profile_name
  location            = azurerm_cdn_endpoint.test.location
  resource_group_name = azurerm_cdn_endpoint.test.resource_group_name

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }
}
`, template)
}

func (r CdnEndpointResource) hostHeader(data acceptance.TestData, domain string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  origin_host_header  = "%s"

  origin {
    name       = "acceptanceTestCdnOrigin2"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, domain)
}

func (r CdnEndpointResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin {
    name       = "acceptanceTestCdnOrigin2"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CdnEndpointResource) withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin {
    name       = "acceptanceTestCdnOrigin2"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CdnEndpointResource) geoFilters(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  is_http_allowed     = false
  is_https_allowed    = true
  origin_path         = "/origin-path"
  probe_path          = "/origin-path/probe"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  geo_filter {
    relative_path = "/some-example-endpoint"
    action        = "Allow"
    country_codes = ["GB"]
  }

  geo_filter {
    relative_path = "/some-other-endpoint"
    action        = "Block"
    country_codes = ["US"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CdnEndpointResource) optimized(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  is_http_allowed     = false
  is_https_allowed    = true
  optimization_type   = "GeneralWebDelivery"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CdnEndpointResource) withoutCompression(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  is_http_allowed     = false
  is_https_allowed    = true
  optimization_type   = "GeneralWebDelivery"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CdnEndpointResource) fullFields(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                          = "acctestcdnend%d"
  profile_name                  = azurerm_cdn_profile.test.name
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  is_http_allowed               = true
  is_https_allowed              = true
  content_types_to_compress     = ["text/html"]
  is_compression_enabled        = true
  querystring_caching_behaviour = "UseQueryString"
  origin_host_header            = "www.contoso.com"
  optimization_type             = "GeneralWebDelivery"
  origin_path                   = "/origin-path"
  probe_path                    = "/origin-path/probe"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  geo_filter {
    relative_path = "/some-example-endpoint"
    action        = "Allow"
    country_codes = ["GB"]
  }

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CdnEndpointResource) isHttpAndHttpsAllowed(data acceptance.TestData, isHttpAllowed string, isHttpsAllowed string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  is_http_allowed     = %s
  is_https_allowed    = %s

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, isHttpAllowed, isHttpsAllowed)
}

func (r CdnEndpointResource) globalDeliveryRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin_host_header = "www.contoso.com"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  global_delivery_rule {
    cache_expiration_action {
      behavior = "Override"
      duration = "5.04:44:23"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CdnEndpointResource) globalDeliveryRuleUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin_host_header = "www.contoso.com"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  global_delivery_rule {
    cache_expiration_action {
      behavior = "SetIfMissing"
      duration = "12.04:11:22"
    }

    modify_response_header_action {
      action = "Overwrite"
      name   = "Content-Type"
      value  = "application/json"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CdnEndpointResource) globalDeliveryRuleRemove(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin_host_header = "www.contoso.com"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CdnEndpointResource) deliveryRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin_host_header = "www.contoso.com"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  delivery_rule {
    name  = "http2https"
    order = 1

    request_scheme_condition {
      match_values = ["HTTP"]
    }

    url_redirect_action {
      redirect_type = "Found"
      protocol      = "Https"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CdnEndpointResource) deliveryRuleUpdate1(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin_host_header = "www.contoso.com"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  delivery_rule {
    name  = "http2https"
    order = 1

    request_scheme_condition {
      negate_condition = true
      match_values     = ["HTTPS"]
    }

    url_redirect_action {
      redirect_type = "Found"
      protocol      = "Https"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CdnEndpointResource) deliveryRuleUpdate2(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin_host_header = "www.contoso.com"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  delivery_rule {
    name  = "http2https"
    order = 1

    request_scheme_condition {
      negate_condition = true
      match_values     = ["HTTPS"]
    }

    url_redirect_action {
      redirect_type = "Found"
      protocol      = "Https"
    }
  }

  delivery_rule {
    name  = "test"
    order = 2

    device_condition {
      match_values = ["Mobile"]
    }

    modify_response_header_action {
      action = "Delete"
      name   = "Content-Language"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CdnEndpointResource) deliveryRuleRemove(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin_host_header = "www.contoso.com"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CdnEndpointResource) testAccAzureRMCdnEndpoint_dnsAlias(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestcdnep%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnep%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnep%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }
}

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  target_resource_id  = azurerm_cdn_endpoint.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r CdnEndpointResource) PremiumVerizon(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Premium_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  is_http_allowed     = false
  is_https_allowed    = true
  querystring_caching_behaviour = "NotSet"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
