// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	legacyrulesets "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CdnFrontDoorBatchRuleResource struct{}

func TestAccCdnFrontDoorBatchRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule", "test")
	r := CdnFrontDoorBatchRuleResource{}

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

func TestAccCdnFrontDoorBatchRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule", "test")
	r := CdnFrontDoorBatchRuleResource{}

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

func TestAccCdnFrontDoorBatchRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule", "test")
	r := CdnFrontDoorBatchRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule", "test")
	r := CdnFrontDoorBatchRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRule_collectionReorderUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule", "test")
	r := CdnFrontDoorBatchRuleResource{}

	primaryRuleName := fmt.Sprintf("accTestBatchRule%d", data.RandomInteger)
	extraRuleName := fmt.Sprintf("accTestBatchRuleExtra%d", data.RandomInteger)
	insertedRuleName := fmt.Sprintf("accTestBatchRuleInserted%d", data.RandomInteger)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.insertAndReorder(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rules.#").HasValue("3"),
				check.That(data.ResourceName).Key("rules.0.name").HasValue(primaryRuleName),
				check.That(data.ResourceName).Key("rules.0.order").HasValue("1"),
				check.That(data.ResourceName).Key("rules.1.name").HasValue(insertedRuleName),
				check.That(data.ResourceName).Key("rules.1.order").HasValue("2"),
				check.That(data.ResourceName).Key("rules.2.name").HasValue(extraRuleName),
				check.That(data.ResourceName).Key("rules.2.order").HasValue("3"),
			),
		},
		{
			Config: r.deleteAndReorder(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rules.#").HasValue("2"),
				check.That(data.ResourceName).Key("rules.0.name").HasValue(insertedRuleName),
				check.That(data.ResourceName).Key("rules.0.order").HasValue("1"),
				check.That(data.ResourceName).Key("rules.1.name").HasValue(extraRuleName),
				check.That(data.ResourceName).Key("rules.1.order").HasValue("2"),
			),
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRule_routeConfigurationOverrideValidation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule", "test")
	r := CdnFrontDoorBatchRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.invalidRouteConfigurationOverrideMissingQueryStringCachingBehavior(data),
			ExpectError: regexp.MustCompile("the `route_configuration_override_action` block is not valid, the `query_string_caching_behavior` field must be set"),
		},
	})
}

func (r CdnFrontDoorBatchRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := rules.ParseRuleSetID(state.ID)
	if err != nil {
		return nil, err
	}

	legacyRuleSetID := legacyrulesets.NewRuleSetID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName)
	client := clients.Cdn.FrontDoorRuleSetsClient_v2025_12_01
	resp, err := client.Get(ctx, legacyRuleSetID)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Rules == nil || len(*resp.Model.Properties.Rules) == 0 {
		return pointer.To(false), nil
	}

	return pointer.To(true), nil
}

func (r CdnFrontDoorBatchRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-afdx-%[1]d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "accTestProfile-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_origin_group" "test" {
  name                     = "accTestOriginGroup-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  load_balancing {
    additional_latency_in_milliseconds = 0
    sample_size                        = 16
    successful_samples_required        = 3
  }
}

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "accTestOrigin-%[1]d"
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

resource "azurerm_cdn_frontdoor_rule_set" "test" {
  name                     = "accTestBatchRuleSet%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  batch_mode_enabled       = true
}

resource "azurerm_cdn_frontdoor_endpoint" "test" {
  name                     = "accTestEndpoint-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}

resource "azurerm_cdn_frontdoor_route" "test" {
  name                          = "accTestRoute-%[1]d"
  cdn_frontdoor_endpoint_id     = azurerm_cdn_frontdoor_endpoint.test.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  cdn_frontdoor_origin_ids      = [azurerm_cdn_frontdoor_origin.test.id]
  cdn_frontdoor_rule_set_ids    = [azurerm_cdn_frontdoor_rule_set.test.id]
  patterns_to_match             = ["/*"]
  supported_protocols           = ["Http", "Https"]
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CdnFrontDoorBatchRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_batch_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  rules {
    name  = "accTestBatchRule%d"
    order = 0

    actions {
      route_configuration_override_action {
        cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
        forwarding_protocol           = "HttpsOnly"
        query_string_caching_behavior = "IncludeSpecifiedQueryStrings"
        query_string_parameters       = ["foo", "clientIp={client_ip}"]
        compression_enabled           = true
        cache_behavior                = "OverrideIfOriginMissing"
        cache_duration                = "365.23:59:59"
      }
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorBatchRuleResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_batch_rule" "import" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  rules {
    name  = "accTestBatchRule%d"
    order = 0

    actions {
      route_configuration_override_action {
        cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
        forwarding_protocol           = "HttpsOnly"
        query_string_caching_behavior = "IncludeSpecifiedQueryStrings"
        query_string_parameters       = ["foo", "clientIp={client_ip}"]
        compression_enabled           = true
        cache_behavior                = "OverrideIfOriginMissing"
        cache_duration                = "365.23:59:59"
      }
    }
  }
}
`, config, data.RandomInteger)
}

func (r CdnFrontDoorBatchRuleResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_batch_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  rules {
    name              = "accTestBatchRule%d"
    behavior_on_match = "Continue"
    order             = 1

    actions {
      route_configuration_override_action {
        cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
        forwarding_protocol           = "HttpsOnly"
        query_string_caching_behavior = "IncludeSpecifiedQueryStrings"
        query_string_parameters       = ["foo", "clientIp={client_ip}"]
        compression_enabled           = true
        cache_behavior                = "OverrideIfOriginMissing"
        cache_duration                = "365.23:59:59"
      }

      response_header_action {
        header_action = "Append"
        header_name   = "Set-Cookie"
        value         = "sessionId=12345678"
      }
    }

    conditions {
      host_name_condition {
        operator         = "Equal"
        negate_condition = false
        match_values     = ["www.contoso.com", "images.contoso.com", "video.contoso.com"]
        transforms       = ["Lowercase", "Trim"]
      }

      is_device_condition {
        operator         = "Equal"
        negate_condition = false
        match_values     = ["Mobile"]
      }

      request_method_condition {
        operator         = "Equal"
        negate_condition = false
        match_values     = ["DELETE"]
      }
    }
  }

  rules {
    name  = "accTestBatchRuleExtra%d"
    order = 2

    actions {
      request_header_action {
        header_action = "Overwrite"
        header_name   = "X-Test"
        value         = "second-rule"
      }
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r CdnFrontDoorBatchRuleResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_batch_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  rules {
    name              = "accTestBatchRuleExtra%d"
    behavior_on_match = "Stop"
    order             = 1

    actions {
      route_configuration_override_action {
        cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
        forwarding_protocol           = "HttpsOnly"
        query_string_caching_behavior = "IgnoreSpecifiedQueryStrings"
        query_string_parameters       = ["clientIp={client_ip}"]
        compression_enabled           = false
        cache_behavior                = "OverrideIfOriginMissing"
        cache_duration                = "23:59:59"
      }
    }

    conditions {
      host_name_condition {
        operator         = "Equal"
        negate_condition = true
        match_values     = ["www.contoso.com", "images.contoso.com", "video.contoso.com"]
        transforms       = ["Lowercase", "Trim"]
      }

      is_device_condition {
        operator         = "Equal"
        negate_condition = true
        match_values     = ["Mobile"]
      }

      request_method_condition {
        operator         = "Equal"
        negate_condition = false
        match_values     = ["DELETE"]
      }
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorBatchRuleResource) insertAndReorder(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_batch_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  rules {
    name              = "accTestBatchRule%d"
    behavior_on_match = "Continue"
    order             = 1

    actions {
      route_configuration_override_action {
        cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
        forwarding_protocol           = "HttpsOnly"
        query_string_caching_behavior = "IncludeSpecifiedQueryStrings"
        query_string_parameters       = ["foo", "clientIp={client_ip}"]
        compression_enabled           = true
        cache_behavior                = "OverrideIfOriginMissing"
        cache_duration                = "365.23:59:59"
      }
    }
  }

  rules {
    name              = "accTestBatchRuleInserted%d"
    behavior_on_match = "Continue"
    order             = 2

    actions {
      request_header_action {
        header_action = "Overwrite"
        header_name   = "X-Test-Inserted"
        value         = "inserted-rule"
      }
    }
  }

  rules {
    name  = "accTestBatchRuleExtra%d"
    order = 3

    actions {
      request_header_action {
        header_action = "Overwrite"
        header_name   = "X-Test"
        value         = "second-rule"
      }
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r CdnFrontDoorBatchRuleResource) deleteAndReorder(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_batch_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  rules {
    name              = "accTestBatchRuleInserted%d"
    behavior_on_match = "Stop"
    order             = 1

    actions {
      route_configuration_override_action {
        cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
        forwarding_protocol           = "HttpsOnly"
        query_string_caching_behavior = "IgnoreSpecifiedQueryStrings"
        query_string_parameters       = ["clientIp={client_ip}"]
        compression_enabled           = false
        cache_behavior                = "OverrideIfOriginMissing"
        cache_duration                = "23:59:59"
      }
    }
  }

  rules {
    name  = "accTestBatchRuleExtra%d"
    order = 2

    actions {
      request_header_action {
        header_action = "Overwrite"
        header_name   = "X-Test"
        value         = "second-rule"
      }
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r CdnFrontDoorBatchRuleResource) invalidRouteConfigurationOverrideMissingQueryStringCachingBehavior(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_batch_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  rules {
    name  = "accTestBatchRule%d"
    order = 1

    actions {
      route_configuration_override_action {
        cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
        forwarding_protocol           = "HttpsOnly"
        cache_behavior                = "OverrideIfOriginMissing"
        cache_duration                = "23:59:59"
      }
    }
  }
}
`, r.template(data), data.RandomInteger)
}
