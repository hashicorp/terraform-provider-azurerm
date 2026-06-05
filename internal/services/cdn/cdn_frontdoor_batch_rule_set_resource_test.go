// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	legacyrulesets "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CdnFrontDoorBatchRuleSetResource struct{}

func TestAccCdnFrontDoorBatchRuleSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontDoorBatchRuleSetResource{}

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

func TestAccCdnFrontDoorBatchRuleSet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontDoorBatchRuleSetResource{}

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

func TestAccCdnFrontDoorBatchRuleSet_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontDoorBatchRuleSetResource{}

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

func TestAccCdnFrontDoorBatchRuleSet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontDoorBatchRuleSetResource{}

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

func TestAccCdnFrontDoorBatchRuleSet_collectionReorderUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontDoorBatchRuleSetResource{}

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

func TestAccCdnFrontDoorBatchRuleSet_routeConfigurationOverrideValidation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontDoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.invalidRouteConfigurationOverrideMissingQueryStringCachingBehavior(data),
			ExpectError: regexp.MustCompile("the 'route_configuration_override_action' block is not valid, the 'query_string_caching_behavior' field must be set"),
		},
	})
}

func TestAccCdnFrontDoorBatchRuleSet_urlFilenameConditionOperatorAny(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontDoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.urlFilenameConditionOperator(data, "Any"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rules.0.conditions.0.url_filename_condition.0.operator").HasValue("Any"),
				check.That(data.ResourceName).Key("rules.0.conditions.0.url_filename_condition.0.match_values").DoesNotExist(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_conditionValidation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontDoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.urlFilenameConditionOperator(data, "Contains"),
			ExpectError: regexp.MustCompile(`"url_filename_condition" is invalid: the 'match_values' field must be set if the conditions 'operator' is not set to 'Any'`),
		},
		{
			Config:      r.requestSchemeConditionOperatorAny(data),
			ExpectError: regexp.MustCompile(`"request_scheme_condition" is invalid: the 'match_values' field must not be set if the conditions 'operator' is set to 'Any'`),
		},
		{
			Config:      r.requestSchemeConditionMissingMatchValues(data),
			ExpectError: regexp.MustCompile("the `request_scheme_condition` block requires `match_values'|the `request_scheme_condition` block requires `match_values`"),
		},
		{
			Config:      r.isDeviceConditionMissingMatchValues(data),
			ExpectError: regexp.MustCompile("the `is_device_condition` block requires `match_values'|the `is_device_condition` block requires `match_values`"),
		},
		{
			Config:      r.remoteAddressGeoMatchInvalid(data),
			ExpectError: regexp.MustCompile(`"remote_address_condition" is invalid: when the 'operator' is set to 'GeoMatch' the value must be a valid country code`),
		},
		{
			Config:      r.socketAddressConditionInvalidCIDR(data),
			ExpectError: regexp.MustCompile(`"socket_address_condition" is invalid: when the 'operator' is set to 'IPMatch' the 'match_values' must be a valid IPv4 or IPv6 CIDR`),
		},
	})
}

func TestAccCdnFrontDoorBatchRuleSet_rulesValidation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontDoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.duplicateRuleName(data),
			ExpectError: regexp.MustCompile("the `rules` blocks must have unique `name` values, got duplicate"),
		},
		{
			Config:      r.duplicateRuleOrder(data),
			ExpectError: regexp.MustCompile("the `rules` blocks must have unique `order` values, got duplicate `1`"),
		},
		{
			Config:      r.unsortedRules(data),
			ExpectError: regexp.MustCompile("the `rules` blocks must be declared in ascending `order`, got `2` before `1`"),
		},
	})
}

func TestAccCdnFrontDoorBatchRuleSet_gapInRuleOrder(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontDoorBatchRuleSetResource{}
	primaryRuleName := fmt.Sprintf("accTestBatchRulePrimary%d", data.RandomInteger)
	gapRuleName := fmt.Sprintf("accTestBatchRuleGap%d", data.RandomInteger)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gapInRuleOrder(data, [2]int{0, 2}),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rules.#").HasValue("2"),
				check.That(data.ResourceName).Key("rules.0.name").HasValue(primaryRuleName),
				check.That(data.ResourceName).Key("rules.0.order").HasValue("0"),
				check.That(data.ResourceName).Key("rules.1.name").HasValue(gapRuleName),
				check.That(data.ResourceName).Key("rules.1.order").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.gapInRuleOrder(data, [2]int{1, 2}),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rules.#").HasValue("2"),
				check.That(data.ResourceName).Key("rules.0.name").HasValue(primaryRuleName),
				check.That(data.ResourceName).Key("rules.0.order").HasValue("1"),
				check.That(data.ResourceName).Key("rules.1.name").HasValue(gapRuleName),
				check.That(data.ResourceName).Key("rules.1.order").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_diffQuotaValidation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontDoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.diffQuotaExceeded(data, 51),
			ExpectError: regexp.MustCompile("the effective diff for `rules` exceeds the service-side quota"),
		},
	})
}

func (r CdnFrontDoorBatchRuleSetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := rules.ParseRuleSetID(state.ID)
	if err != nil {
		return nil, err
	}

	ruleSetResourceId := legacyrulesets.NewRuleSetID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName)
	batchModeRuleSetClient := clients.Cdn.FrontDoorRuleSetsClient_v2025_12_01
	resp, err := batchModeRuleSetClient.Get(ctx, ruleSetResourceId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Rules == nil || len(*resp.Model.Properties.Rules) == 0 {
		return pointer.To(false), nil
	}

	return pointer.To(true), nil
}

func (r CdnFrontDoorBatchRuleSetResource) template(data acceptance.TestData) string {
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

resource "azurerm_cdn_frontdoor_endpoint" "test" {
  name                     = "accTestEndpoint-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}

resource "azurerm_cdn_frontdoor_route" "test" {
  name                          = "accTestRoute-%[1]d"
  cdn_frontdoor_endpoint_id     = azurerm_cdn_frontdoor_endpoint.test.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  cdn_frontdoor_origin_ids      = [azurerm_cdn_frontdoor_origin.test.id]
  cdn_frontdoor_rule_set_ids    = [azurerm_cdn_frontdoor_batch_rule_set.test.id]
  patterns_to_match             = ["/*"]
  supported_protocols           = ["Http", "Https"]
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CdnFrontDoorBatchRuleSetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "accTestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rules {
    name  = "accTestBatchRule%[2]d"
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

func (r CdnFrontDoorBatchRuleSetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "import" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "accTestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rules {
    name  = "accTestBatchRule%[2]d"
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
`, r.basic(data), data.RandomInteger)
}

func (r CdnFrontDoorBatchRuleSetResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "accTestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rules {
    name              = "accTestBatchRule%[2]d"
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
    name  = "accTestBatchRuleExtra%[3]d"
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

func (r CdnFrontDoorBatchRuleSetResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "accTestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rules {
    name              = "accTestBatchRuleExtra%[2]d"
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

func (r CdnFrontDoorBatchRuleSetResource) insertAndReorder(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "accTestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rules {
    name              = "accTestBatchRule%[2]d"
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
    name              = "accTestBatchRuleInserted%[3]d"
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
    name  = "accTestBatchRuleExtra%[4]d"
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

func (r CdnFrontDoorBatchRuleSetResource) deleteAndReorder(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "accTestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rules {
    name              = "accTestBatchRuleInserted%[2]d"
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
    name  = "accTestBatchRuleExtra%[3]d"
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

func (r CdnFrontDoorBatchRuleSetResource) invalidRouteConfigurationOverrideMissingQueryStringCachingBehavior(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "accTestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rules {
    name  = "accTestBatchRule%[2]d"
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

func (r CdnFrontDoorBatchRuleSetResource) urlFilenameConditionOperator(data acceptance.TestData, operator string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "accTestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rules {
    name              = "accTestBatchRule%[2]d"
    behavior_on_match = "Stop"
    order             = 1

    actions {
      url_rewrite_action {
        source_pattern          = "/"
        destination             = "/index.html"
        preserve_unmatched_path = false
      }
    }

    conditions {
      url_filename_condition {
        operator = "%[3]s"
      }
    }
  }
}
`, r.template(data), data.RandomInteger, operator)
}

func (r CdnFrontDoorBatchRuleSetResource) requestSchemeConditionOperatorAny(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "accTestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rules {
    name  = "accTestBatchRule%[2]d"
    order = 1

    actions {
      url_rewrite_action {
        source_pattern          = "/"
        destination             = "/index.html"
        preserve_unmatched_path = false
      }
    }

    conditions {
      request_scheme_condition {
        operator     = "Any"
        match_values = ["HTTP"]
      }
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorBatchRuleSetResource) requestSchemeConditionMissingMatchValues(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "accTestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rules {
    name  = "accTestBatchRule%[2]d"
    order = 1

    actions {
      url_rewrite_action {
        source_pattern          = "/"
        destination             = "/index.html"
        preserve_unmatched_path = false
      }
    }

    conditions {
      request_scheme_condition {
        negate_condition = false
        operator         = "Equal"
      }
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorBatchRuleSetResource) isDeviceConditionMissingMatchValues(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "accTestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rules {
    name  = "accTestBatchRule%[2]d"
    order = 1

    actions {
      url_rewrite_action {
        source_pattern          = "/"
        destination             = "/index.html"
        preserve_unmatched_path = false
      }
    }

    conditions {
      is_device_condition {
        negate_condition = false
        operator         = "Equal"
      }
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorBatchRuleSetResource) remoteAddressGeoMatchInvalid(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "accTestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rules {
    name  = "accTestBatchRule%[2]d"
    order = 1

    actions {
      url_rewrite_action {
        source_pattern          = "/"
        destination             = "/index.html"
        preserve_unmatched_path = false
      }
    }

    conditions {
      remote_address_condition {
        operator     = "GeoMatch"
        match_values = ["us"]
      }
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorBatchRuleSetResource) socketAddressConditionInvalidCIDR(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "accTestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rules {
    name  = "accTestBatchRule%[2]d"
    order = 1

    actions {
      url_rewrite_action {
        source_pattern          = "/"
        destination             = "/index.html"
        preserve_unmatched_path = false
      }
    }

    conditions {
      socket_address_condition {
        operator     = "IPMatch"
        match_values = ["not-a-cidr"]
      }
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorBatchRuleSetResource) duplicateRuleName(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "accTestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rules {
    name  = "accTestBatchRule%[2]d"
    order = 1

    actions {
      url_rewrite_action {
        source_pattern          = "/"
        destination             = "/first.html"
        preserve_unmatched_path = false
      }
    }
  }

  rules {
    name  = "accTestBatchRule%[3]d"
    order = 2

    actions {
      url_rewrite_action {
        source_pattern          = "/"
        destination             = "/second.html"
        preserve_unmatched_path = false
      }
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r CdnFrontDoorBatchRuleSetResource) duplicateRuleOrder(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "accTestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rules {
    name  = "accTestBatchRule%[2]d"
    order = 1

    actions {
      url_rewrite_action {
        source_pattern          = "/"
        destination             = "/first.html"
        preserve_unmatched_path = false
      }
    }
  }

  rules {
    name  = "accTestBatchRuleExtra%[3]d"
    order = 1

    actions {
      url_rewrite_action {
        source_pattern          = "/"
        destination             = "/second.html"
        preserve_unmatched_path = false
      }
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r CdnFrontDoorBatchRuleSetResource) unsortedRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "accTestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rules {
    name  = "accTestBatchRule%[2]d"
    order = 2

    actions {
      url_rewrite_action {
        source_pattern          = "/"
        destination             = "/second.html"
        preserve_unmatched_path = false
      }
    }
  }

  rules {
    name  = "accTestBatchRulePrimary%[3]d"
    order = 1

    actions {
      url_rewrite_action {
        source_pattern          = "/"
        destination             = "/first.html"
        preserve_unmatched_path = false
      }
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r CdnFrontDoorBatchRuleSetResource) gapInRuleOrder(data acceptance.TestData, orders [2]int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "accTestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rules {
    name  = "accTestBatchRulePrimary%[2]d"
    order = %[4]d

    actions {
      url_rewrite_action {
        source_pattern          = "/"
        destination             = "/first.html"
        preserve_unmatched_path = false
      }
    }
  }

  rules {
    name  = "accTestBatchRuleGap%[3]d"
    order = %[5]d

    actions {
      url_rewrite_action {
        source_pattern          = "/"
        destination             = "/second.html"
        preserve_unmatched_path = false
      }
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, orders[0], orders[1])
}

func (r CdnFrontDoorBatchRuleSetResource) diffQuotaExceeded(data acceptance.TestData, ruleCount int) string {
	var rulesBuilder strings.Builder
	for index := 0; index < ruleCount; index++ {
		rulesBuilder.WriteString(fmt.Sprintf(`
  rules {
    name  = "accTestBatchRuleQuota%[1]d"
    order = %[1]d

    actions {
      route_configuration_override_action {
        cache_behavior                = "OverrideIfOriginMissing"
        query_string_caching_behavior = "UseQueryString"
        cache_duration                = "365.23:59:59"
      }
    }
  }
`, index+1))
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "accTestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
%[3]s
}
`, r.template(data), data.RandomInteger, rulesBuilder.String())
}
