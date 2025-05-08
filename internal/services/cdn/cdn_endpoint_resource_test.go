// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnEndpointResource struct{}

func TestAccCdnEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccCdnEndpoint_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccCdnEndpoint_updateHostHeader(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.hostHeader(data, "www.contoso.com"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.hostHeader(data, "www.example2.com"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpoint_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withTagsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpoint_withoutCompression(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withoutCompression(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpoint_optimized(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.optimized(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpoint_isHttpAndHttpsAllowedUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.isHttpAndHttpsAllowed(data, "true", "false"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.isHttpAndHttpsAllowed(data, "false", "true"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpoint_globalDeliveryRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.globalDeliveryRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.globalDeliveryRuleUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.globalDeliveryRuleRemove(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpoint_deliveryRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.deliveryRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.deliveryRuleUpdate1(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.deliveryRuleUpdate2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.deliveryRuleRemove(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpoint_dnsAlias(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.testAccAzureRMCdnEndpoint_dnsAlias(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpoint_deliveryRuleOptionalMatchValue(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.deliveryRuleOptionalMatchValue(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// Covers https://github.com/hashicorp/terraform-provider-azurerm/issue/21450
func TestAccCdnEndpoint_longQueryString(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.longQueryString(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// Covers https://github.com/hashicorp/terraform-provider-azurerm/issues/22326
func TestAccCdnEndpoint_compressionUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	r := CdnEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.compressionUpdate(data, "true", "sandbox"), // PUT
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("is_compression_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("sandbox"),
			),
		},
		data.ImportStep(),
		{
			Config: r.compressionUpdate(data, "false", "sandbox"), // PUT
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("is_compression_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("sandbox"),
			),
		},
		data.ImportStep(),
		{
			Config: r.compressionUpdate(data, "true", "sandbox"), // PUT
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("is_compression_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("sandbox"),
			),
		},
		data.ImportStep(),
		{
			Config: r.compressionUpdate(data, "true", "production"), // PATCH
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("is_compression_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("production"),
			),
		},
		data.ImportStep(),
	})
}

func (r CdnEndpointResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r CdnEndpointResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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
  sku                 = "Standard_Microsoft"
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
  sku                 = "Standard_Microsoft"
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
  sku                 = "Standard_Microsoft"
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
  sku                 = "Standard_Microsoft"
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
  sku                 = "Standard_Microsoft"
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
  sku                 = "Standard_Microsoft"
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
  sku                 = "Standard_Microsoft"
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
  name     = "acctestRG-%[1]d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%[1]d"
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
    cache_key_query_string_action {
      behavior   = "IncludeAll"
      parameters = "test"
    }
    modify_request_header_action {
      action = "Append"
      name   = "www.contoso1.com"
      value  = "test value"
    }
    url_redirect_action {
      redirect_type = "Found"
      protocol      = "Https"
      hostname      = "www.contoso.com"
      fragment      = "5fgdfg"
      path          = "/article.aspx"
      query_string  = "id={var_uri_path_1}&title={var_uri_path_2}"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
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
    url_rewrite_action {
      source_pattern          = "/test_source_pattern"
      destination             = "/test_destination"
      preserve_unmatched_path = false
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
    cache_expiration_action {
      behavior = "Override"
      duration = "5.04:44:23"
    }
    cache_key_query_string_action {
      behavior   = "IncludeAll"
      parameters = "test"
    }
    cookies_condition {
      operator         = "Contains"
      selector         = "abc"
      negate_condition = false
      match_values     = ["windows"]
      transforms       = ["Lowercase"]

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
  sku                 = "Standard_Microsoft"
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

func (r CdnEndpointResource) deliveryRuleOptionalMatchValue(data acceptance.TestData) string {
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
    name  = "cookieCondition"
    order = 1

    cookies_condition {
      selector = "abc"
      operator = "Any"
    }

    modify_response_header_action {
      action = "Delete"
      name   = "Content-Language"
    }
  }

  delivery_rule {
    name  = "postArg"
    order = 2

    post_arg_condition {
      selector = "abc"
      operator = "Any"
    }

    modify_response_header_action {
      action = "Delete"
      name   = "Content-Language"
    }
  }

  delivery_rule {
    name  = "queryString"
    order = 3

    query_string_condition {
      operator = "Any"
    }

    modify_response_header_action {
      action = "Delete"
      name   = "Content-Language"
    }
  }

  delivery_rule {
    name  = "remoteAddress"
    order = 4

    remote_address_condition {
      operator = "Any"
    }

    modify_response_header_action {
      action = "Delete"
      name   = "Content-Language"
    }
  }

  delivery_rule {
    name  = "requestBody"
    order = 5

    request_body_condition {
      operator = "Any"
    }

    modify_response_header_action {
      action = "Delete"
      name   = "Content-Language"
    }
  }

  delivery_rule {
    name  = "requestHeader"
    order = 6

    request_header_condition {
      selector = "abc"
      operator = "Any"
    }

    modify_response_header_action {
      action = "Delete"
      name   = "Content-Language"
    }
  }

  delivery_rule {
    name  = "requestUri"
    order = 7

    request_uri_condition {
      operator = "Any"
    }

    modify_response_header_action {
      action = "Delete"
      name   = "Content-Language"
    }
  }

  delivery_rule {
    name  = "uriFileExtension"
    order = 8

    url_file_extension_condition {
      operator = "Any"
    }

    modify_response_header_action {
      action = "Delete"
      name   = "Content-Language"
    }
  }

  delivery_rule {
    name  = "uriFileName"
    order = 9

    url_file_name_condition {
      operator = "Any"
    }

    modify_response_header_action {
      action = "Delete"
      name   = "Content-Language"
    }
  }

  delivery_rule {
    name  = "uriPath"
    order = 10

    url_path_condition {
      operator = "Any"
    }

    modify_response_header_action {
      action = "Delete"
      name   = "Content-Language"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CdnEndpointResource) longQueryString(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctesa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

data "azurerm_storage_account_sas" "test" {
  connection_string = azurerm_storage_account.test.primary_connection_string
  https_only        = true

  resource_types {
    service   = false
    container = false
    object    = true
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  start  = "2023-04-01T00:00:00Z"
  expiry = "2123-04-01T00:00:00Z"

  permissions {
    read    = true
    write   = false
    delete  = false
    list    = false
    add     = false
    create  = false
    update  = false
    process = false
    tag     = false
    filter  = false
  }
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%[1]d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin {
    name      = "acceptanceTestCdnOrigin1"
    host_name = "www.contoso.com"
  }

  delivery_rule {
    name  = "TokenSAS"
    order = 1
    query_string_condition {
      operator         = "Contains"
      negate_condition = true
      match_values     = ["sig"]
    }
    url_redirect_action {
      redirect_type = "PermanentRedirect"
      query_string  = trimprefix(data.azurerm_storage_account_sas.test.sas, "?")
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(8))
}

func (r CdnEndpointResource) compressionUpdate(data acceptance.TestData, compressionEnabled string, tag string) string {
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
  name                      = "acctestcdnend%d"
  profile_name              = azurerm_cdn_profile.test.name
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  is_http_allowed           = true
  is_https_allowed          = true
  content_types_to_compress = ["text/html"]
  is_compression_enabled    = %s

  origin {
    name      = "acceptanceTestCdnOrigin1"
    host_name = "www.contoso.com"
  }

  tags = {
    environment = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, compressionEnabled, tag)
}
