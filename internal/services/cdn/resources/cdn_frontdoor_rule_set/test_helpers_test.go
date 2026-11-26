// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_rule_set_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type CdnFrontdoorBatchRuleSetResource struct{}

func (r CdnFrontdoorBatchRuleSetResource) basicUnattachedRoute(data acceptance.TestData) string {
	return r.basicWithTemplate(data, r.templateUnattachedRoute(data))
}

func (r CdnFrontdoorBatchRuleSetResource) basicWithTemplate(data acceptance.TestData, template string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "acctestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rule {
    name  = "acctestBatchRule%[2]d"
    order = 0

    actions {
      route_configuration_override {
        origin_group {
          cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
          forwarding_protocol           = "HttpsOnly"
        }
        caching {
          query_string_behaviour  = "IncludeSpecifiedQueryStrings"
          query_string_parameters = ["foo", "clientIp={client_ip}"]
          compression_enabled     = true
          behaviour               = "OverrideIfOriginMissing"
          duration                = "365.23:59:59"
        }
      }
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "acctestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rule {
    name               = "acctestBatchRule%[2]d"
    behaviour_on_match = "Continue"
    order              = 1

    actions {
      route_configuration_override {
        origin_group {
          cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
          forwarding_protocol           = "HttpsOnly"
        }
        caching {
          query_string_behaviour  = "IncludeSpecifiedQueryStrings"
          query_string_parameters = ["foo", "clientIp={client_ip}"]
          compression_enabled     = true
          behaviour               = "OverrideIfOriginMissing"
          duration                = "365.23:59:59"
        }
      }

      modify_response_header {
        operator     = "Append"
        header_name  = "Set-Cookie"
        header_value = "sessionId=12345678"
      }
    }

    conditions {
      host_name {
        operator   = "Equal"
        values     = ["www.contoso.com", "images.contoso.com", "video.contoso.com"]
        transforms = ["Lowercase", "Trim"]
      }

      device_type {
        operator = "Equal"
        values   = ["Mobile"]
      }

      request_method {
        operator = "Equal"
        values   = ["DELETE"]
      }
    }
  }

  rule {
    name  = "acctestBatchRuleExtra%[3]d"
    order = 2

    actions {
      modify_request_header {
        operator     = "Overwrite"
        header_name  = "X-Test"
        header_value = "second-rule"
      }
    }
  }
}
`, r.templateAttachedRoute(data), data.RandomInteger, data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) disableCacheAndNoOriginGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "acctestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rule {
    name  = "acctestBatchRule%[2]d"
    order = 0

    actions {
      route_configuration_override {
        caching {
          behaviour = "Disabled"
        }
      }
    }
  }
}
`, r.templateAttachedRoute(data), data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) routeTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_cdn_frontdoor_endpoint" "test" {
  name                     = "acctestEndpoint-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}

resource "azurerm_cdn_frontdoor_route" "test" {
  name                          = "acctestRoute-%[1]d"
  cdn_frontdoor_endpoint_id     = azurerm_cdn_frontdoor_endpoint.test.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  cdn_frontdoor_origin_ids      = [azurerm_cdn_frontdoor_origin.test.id]
  cdn_frontdoor_rule_set_ids    = [azurerm_cdn_frontdoor_batch_rule_set.test.id]
  patterns_to_match             = ["/*"]
  supported_protocols           = ["Http", "Https"]
}
`, data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) templateAttachedRoute(data acceptance.TestData) string {
	return fmt.Sprintf(`%[1]s

%[2]s`, r.templateUnattachedRoute(data), r.routeTemplate(data))
}

func (r CdnFrontdoorBatchRuleSetResource) templateUnattachedRoute(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-afdx-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestProfile-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_origin_group" "test" {
  name                     = "acctestOriginGroup-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  load_balancing {
    additional_latency_in_milliseconds = 0
    sample_size                        = 16
    successful_samples_required        = 3
  }
}

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "acctestOrigin-%[1]d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  enabled                       = true

  certificate_name_check_enabled = false
  host_name                      = "contoso.com"
  http_port                      = 80
  https_port                     = 443
  origin_host_header             = "www.contoso.com"
  priority                       = 1
  weight                         = 1
}
`, data.RandomInteger, data.Locations.Primary)
}
