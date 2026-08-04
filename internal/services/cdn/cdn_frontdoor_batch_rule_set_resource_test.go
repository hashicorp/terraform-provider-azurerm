// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CdnFrontdoorBatchRuleSetResource struct{}

func TestAccCdnFrontDoorBatchRuleSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAttachedRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAttachedRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_importNonBatchRuleSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: CdnFrontDoorRuleSetResource{}.basic(data, false),
		},
		{
			Config:      r.nonBatchRuleSetImport(data),
			ExpectError: regexp.MustCompile("was not provisioned using batch mode and cannot be managed by this resource"),
		},
	})
}

func TestAccCdnFrontDoorBatchRuleSet_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

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
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAttachedRoute(data),
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
		{
			Config: r.basicAttachedRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_disableCache_attachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAttachedRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.disableCacheAttachedRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicAttachedRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.disableCacheAndNoOriginGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_disableCache_unattachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disableCacheUnattachedRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_cacheDurationZero_attachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cacheDurationZeroAttachedRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_cacheDurationZero_unattachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cacheDurationZeroUnattachedRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_originGroupIdOptional_attachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.originGroupIdOptionalAttachedRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_originGroupIdOptional_unattachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.originGroupIdOptionalWithoutOrigin(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_originGroupIdOptionalUpdate_attachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.originGroupIdOptionalAttachedRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicAttachedRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.originGroupIdOptionalAttachedRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_originGroupIdOptionalUpdate_unattachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.originGroupIdOptionalWithoutOrigin(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicUnattachedRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicUnattachedRouteWithoutOriginGroupOverride(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.originGroupIdOptionalWithoutOrigin(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_honorOrigin_attachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.honorOriginAttachedRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_honorOrigin_unattachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.honorOriginWithoutOrigin(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_collectionReorderUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

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
			),
		},
		{
			Config: r.deleteAndReorder(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
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
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.invalidRouteConfigurationOverrideCaching(data, `
        caching {
          behaviour = "Disabled"
          duration  = "23:59:59"
        }`),
			ExpectError: regexp.MustCompile("when `route_configuration_override.caching.behaviour` is set to `Disabled`, you cannot define `route_configuration_override.caching.duration`"),
		},
		{
			Config: r.invalidRouteConfigurationOverrideCaching(data, `
        caching {
          behaviour               = "Disabled"
          query_string_parameters = ["foo"]
        }`),
			ExpectError: regexp.MustCompile("when `route_configuration_override.caching.behaviour` is set to `Disabled`, you cannot define any `route_configuration_override.caching.query_string_parameters`"),
		},
		{
			Config: r.invalidRouteConfigurationOverrideCaching(data, `
        caching {
          behaviour              = "Disabled"
          query_string_behaviour = "UseQueryString"
        }`),
			ExpectError: regexp.MustCompile("when `route_configuration_override.caching.behaviour` is set to `Disabled`, you cannot define `route_configuration_override.caching.query_string_behaviour`"),
		},
		{
			Config: r.invalidRouteConfigurationOverrideCaching(data, `
        caching {
          behaviour           = "Disabled"
          compression_enabled = true
        }`),
			ExpectError: regexp.MustCompile("when `route_configuration_override.caching.behaviour` is set to `Disabled`, you cannot define `route_configuration_override.caching.compression_enabled`"),
		},
		{
			Config: r.invalidRouteConfigurationOverrideCaching(data, `
        caching {
          behaviour = "OverrideIfOriginMissing"
          duration  = "23:59:59"
        }`),
			ExpectError: regexp.MustCompile("when `route_configuration_override.caching.behaviour` is set to `OverrideIfOriginMissing`, you must also define `route_configuration_override.caching.query_string_behaviour`"),
		},
		{
			Config: r.invalidRouteConfigurationOverrideCaching(data, `
        caching {
          behaviour              = "HonorOrigin"
          query_string_behaviour = "UseQueryString"
          duration               = "23:59:59"
        }`),
			ExpectError: regexp.MustCompile("when `route_configuration_override.caching.behaviour` is set to `HonorOrigin`, you cannot define `route_configuration_override.caching.duration`"),
		},
		{
			Config: r.invalidRouteConfigurationOverrideCaching(data, `
        caching {
          behaviour              = "OverrideAlways"
          query_string_behaviour = "UseQueryString"
        }`),
			ExpectError: regexp.MustCompile("when `route_configuration_override.caching.behaviour` is set to `OverrideAlways`, you must also define `route_configuration_override.caching.duration`"),
		},
		{
			Config: r.invalidRouteConfigurationOverrideCaching(data, `
        caching {
          behaviour              = "OverrideAlways"
          duration               = "23:59:59"
          query_string_behaviour = "IncludeSpecifiedQueryStrings"
        }`),
			ExpectError: regexp.MustCompile("when `route_configuration_override.caching.query_string_behaviour` is set to `IncludeSpecifiedQueryStrings`, you must also define one or more `route_configuration_override.caching.query_string_parameters`"),
		},
		{
			Config: r.invalidRouteConfigurationOverrideCaching(data, `
        caching {
          behaviour               = "OverrideAlways"
          duration                = "23:59:59"
          query_string_behaviour  = "UseQueryString"
          query_string_parameters = ["foo"]
        }`),
			ExpectError: regexp.MustCompile("when `route_configuration_override.caching.query_string_behaviour` is set to `UseQueryString`, you cannot define `route_configuration_override.caching.query_string_parameters`"),
		},
	})
}

func TestAccCdnFrontDoorBatchRuleSet_urlFilenameConditionOperatorAny(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.urlFilenameConditionOperator(data, "Any"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_conditionValidation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.urlFilenameConditionOperator(data, "Contains"),
			ExpectError: regexp.MustCompile("when `conditions.request_filename.operator` is set to `Contains`, `conditions.request_filename.values` must set one or more values"),
		},
		{
			Config:      r.conditionAnyOperatorWithValues(data),
			ExpectError: regexp.MustCompile("when `conditions.request_path.operator` is set to `Any`, `conditions.request_path.values` cannot be defined"),
		},
		{
			Config:      r.remoteAddressGeoMatchInvalid(data),
			ExpectError: regexp.MustCompile(`remote_address.*valid country code`),
		},
	})
}

func TestAccCdnFrontDoorBatchRuleSet_rulesValidation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.duplicateRuleName(data),
			ExpectError: regexp.MustCompile("the `rule` blocks must have unique `name` values, got duplicate"),
		},
		{
			Config:      r.duplicateRuleOrder(data),
			ExpectError: regexp.MustCompile("the `rule` blocks must have unique `order` values, got duplicate `1`"),
		},
		{
			Config:      r.unsortedRules(data),
			ExpectError: regexp.MustCompile("the `rule` blocks must be declared in ascending `order`, got `2` before `1`"),
		},
	})
}

func TestAccCdnFrontDoorBatchRuleSet_modifyHeaderActionValidation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.modifyHeaderActionMissingValue(data),
			ExpectError: regexp.MustCompile("the `modify_request_header` block is not valid, `header_value` cannot be empty if the `operator` is set to `Append` or `Overwrite`"),
		},
		{
			Config:      r.modifyHeaderActionUnexpectedValue(data),
			ExpectError: regexp.MustCompile("the `modify_response_header` block is not valid, `header_value` must be empty if the `operator` is set to `Delete`"),
		},
	})
}

func TestAccCdnFrontDoorBatchRuleSet_actionCountValidation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.actionCountValidation(data, ``),
			ExpectError: regexp.MustCompile("the `actions` block must define at least one action"),
		},
		{
			Config: r.actionCountValidation(data, `
      url_redirect {
        redirect_type = "Found"
      }

      url_rewrite {
        source_pattern   = "/"
        destination_path = "/rewritten"
      }`),
			ExpectError: regexp.MustCompile("cannot specify both `url_redirect` and the `url_rewrite` in the `actions` block"),
		},
		{
			Config: r.actionCountValidation(data, `
      url_redirect {
        redirect_type = "Found"
      }

      route_configuration_override {
        origin_group {
          cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
          forwarding_protocol           = "HttpsOnly"
        }

        caching {
          behaviour = "Disabled"
        }
      }`),
			ExpectError: regexp.MustCompile("cannot specify both `url_redirect` and the `route_configuration_override` in the `actions` block"),
		},
		{
			Config: r.actionCountValidation(data, `
      modify_request_header {
        operator     = "Overwrite"
        header_name  = "X-1"
        header_value = "1"
      }

      modify_request_header {
        operator     = "Overwrite"
        header_name  = "X-2"
        header_value = "2"
      }

      modify_request_header {
        operator     = "Overwrite"
        header_name  = "X-3"
        header_value = "3"
      }

      modify_response_header {
        operator     = "Overwrite"
        header_name  = "Y-1"
        header_value = "1"
      }

      modify_response_header {
        operator     = "Overwrite"
        header_name  = "Y-2"
        header_value = "2"
      }

      modify_response_header {
        operator     = "Overwrite"
        header_name  = "Y-3"
        header_value = "3"
      }`),
			ExpectError: regexp.MustCompile("the `actions` block may only contain up to 5 actions"),
		},
	})
}

func TestAccCdnFrontDoorBatchRuleSet_gapInRuleOrder(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gapInRuleOrder(data, [2]int{0, 2}),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.gapInRuleOrder(data, [2]int{1, 2}),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_diffQuotaValidation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.allConditions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config:      r.diffQuotaExceeded(data),
			ExpectError: regexp.MustCompile("the number of changed `rule` blocks exceeds the service-side quota"),
		},
	})
}

func TestAccCdnFrontDoorBatchRuleSet_urlRedirect(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.urlRedirect(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_allConditions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.allConditions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_allConditionsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.allConditions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.allConditionsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.allConditions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_allActions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.allActions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicUnattachedRouteWithoutOriginGroupOverride(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.originGroupIdOptionalWithoutOrigin(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorBatchRuleSet_allActionsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "test")
	r := CdnFrontdoorBatchRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.allActions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.allActionsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.allActions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicUnattachedRouteWithoutOriginGroupOverride(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.originGroupIdOptionalWithoutOrigin(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CdnFrontdoorBatchRuleSetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := rulesets.ParseRuleSetID(state.ID)
	if err != nil {
		return nil, err
	}

	batchModeRuleSetClient := clients.Cdn.FrontDoorRuleSetsClient
	resp, err := batchModeRuleSetClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r CdnFrontdoorBatchRuleSetResource) nonBatchRuleSetImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

import {
  id = azurerm_cdn_frontdoor_rule_set.test.id
  to = azurerm_cdn_frontdoor_batch_rule_set.test
}

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  name                     = "acctestBatchRuleSet0706"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rule {
    name  = "acctestBatchRule%[2]d"
    order = 0

    actions {
      url_rewrite {
        source_pattern   = "/test"
        destination_path = "/new"
      }
    }
  }
}
`, CdnFrontDoorRuleSetResource{}.basic(data, false), data.RandomInteger)
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

func (r CdnFrontdoorBatchRuleSetResource) templateWithoutOrigin(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary)
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

func (r CdnFrontdoorBatchRuleSetResource) basicAttachedRoute(data acceptance.TestData) string {
	return r.basicWithTemplate(data, r.templateAttachedRoute(data))
}

func (r CdnFrontdoorBatchRuleSetResource) basicUnattachedRoute(data acceptance.TestData) string {
	return r.basicWithTemplate(data, r.templateUnattachedRoute(data))
}

func (r CdnFrontdoorBatchRuleSetResource) basicUnattachedRouteWithoutOriginGroupOverride(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  name                     = "acctestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rule {
    name  = "acctestBatchRule%[2]d"
    order = 0

    actions {
      route_configuration_override {
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
`, r.templateUnattachedRoute(data), data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) allActions(data acceptance.TestData) string {
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
    order = 1

    actions {
      url_redirect {
        redirect_type         = "Found"
        redirect_protocol     = "Https"
        destination_host_name = "contoso.com"
        destination_path      = "/redirected"
        destination_fragment  = "fragment"
        query_string          = "a=b"
      }

      modify_request_header {
        operator     = "Overwrite"
        header_name  = "X-Request"
        header_value = "request-value"
      }

      modify_response_header {
        operator     = "Append"
        header_name  = "X-Response"
        header_value = "response-value"
      }
    }
  }

  rule {
    name  = "acctestBatchRule2%[2]d"
    order = 2

    actions {
      url_rewrite {
        source_pattern                  = "/"
        destination_path                = "/rewritten"
        preserve_unmatched_path_enabled = true
      }
    }
  }

  rule {
    name  = "acctestBatchRule3%[2]d"
    order = 3

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
`, r.templateUnattachedRoute(data), data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) allActionsUpdate(data acceptance.TestData) string {
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
    order = 10

    actions {
      url_redirect {
        redirect_type         = "PermanentRedirect"
        redirect_protocol     = "Http"
        destination_host_name = "www.contoso.com"
        destination_path      = "/moved"
        query_string          = "c=d"
      }

      modify_request_header {
        operator     = "Append"
        header_name  = "X-RequestUpdate"
        header_value = "request-value-update"
      }

      modify_response_header {
        operator    = "Delete"
        header_name = "X-Response"
      }
    }
  }

  rule {
    name  = "acctestBatchRule2%[2]d"
    order = 20

    actions {
      url_rewrite {
        source_pattern                  = "/updated"
        destination_path                = "/rewritten-update"
        preserve_unmatched_path_enabled = false
      }
    }
  }

  rule {
    name  = "acctestBatchRule3%[2]d"
    order = 30

    actions {
      route_configuration_override {
        origin_group {
          cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
          forwarding_protocol           = "MatchRequest"
        }
        caching {
          query_string_behaviour = "IgnoreQueryString"
          compression_enabled    = false
          behaviour              = "HonorOrigin"
        }
      }
    }
  }
}
`, r.templateUnattachedRoute(data), data.RandomInteger)
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

func (r CdnFrontdoorBatchRuleSetResource) disableCacheAttachedRoute(data acceptance.TestData) string {
	return r.disableCacheWithTemplate(data, r.templateAttachedRoute(data))
}

func (r CdnFrontdoorBatchRuleSetResource) disableCacheUnattachedRoute(data acceptance.TestData) string {
	return r.disableCacheWithTemplate(data, r.templateUnattachedRoute(data))
}

func (r CdnFrontdoorBatchRuleSetResource) disableCacheWithTemplate(data acceptance.TestData, template string) string {
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
          behaviour = "Disabled"
        }
      }
    }
  }
}
`, template, data.RandomInteger)
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

func (r CdnFrontdoorBatchRuleSetResource) cacheDurationZeroAttachedRoute(data acceptance.TestData) string {
	return r.cacheDurationZeroWithTemplate(data, r.templateAttachedRoute(data))
}

func (r CdnFrontdoorBatchRuleSetResource) cacheDurationZeroUnattachedRoute(data acceptance.TestData) string {
	return r.cacheDurationZeroWithoutOrigin(data)
}

func (r CdnFrontdoorBatchRuleSetResource) cacheDurationZeroWithoutOrigin(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  name                     = "acctestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rule {
    name  = "acctestBatchRule%[2]d"
    order = 0

    actions {
      route_configuration_override {
        caching {
          query_string_behaviour  = "IncludeSpecifiedQueryStrings"
          query_string_parameters = ["foo", "clientIp={client_ip}"]
          compression_enabled     = true
          behaviour               = "OverrideIfOriginMissing"
          duration                = "00:00:00"
        }
      }
    }
  }
}
`, r.templateWithoutOrigin(data), data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) cacheDurationZeroWithTemplate(data acceptance.TestData, template string) string {
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
          duration                = "00:00:00"
        }
      }
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) originGroupIdOptionalAttachedRoute(data acceptance.TestData) string {
	return r.originGroupIdOptionalWithTemplate(data, r.templateAttachedRoute(data))
}

func (r CdnFrontdoorBatchRuleSetResource) originGroupIdOptionalWithoutOrigin(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  name                     = "acctestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rule {
    name  = "acctestBatchRule%[2]d"
    order = 0

    actions {
      route_configuration_override {
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
`, r.templateWithoutOrigin(data), data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) originGroupIdOptionalWithTemplate(data acceptance.TestData, template string) string {
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

func (r CdnFrontdoorBatchRuleSetResource) honorOriginAttachedRoute(data acceptance.TestData) string {
	return r.honorOriginWithTemplate(data, r.templateAttachedRoute(data))
}

func (r CdnFrontdoorBatchRuleSetResource) honorOriginWithoutOrigin(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  name                     = "acctestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rule {
    name  = "acctestBatchRule%[2]d"
    order = 0

    actions {
      route_configuration_override {
        caching {
          behaviour              = "HonorOrigin"
          compression_enabled    = true
          query_string_behaviour = "IgnoreQueryString"
        }
      }
    }

    conditions {
      request_path {
        values   = ["data/"]
        operator = "BeginsWith"
      }

      request_path {
        values   = [".html", ".htm"]
        operator = "EndsWith"
      }
    }
  }
}
`, r.templateWithoutOrigin(data), data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) honorOriginWithTemplate(data acceptance.TestData, template string) string {
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
          behaviour              = "HonorOrigin"
          compression_enabled    = true
          query_string_behaviour = "IgnoreQueryString"
        }
      }
    }

    conditions {
      request_path {
        values   = ["data/"]
        operator = "BeginsWith"
      }

      request_path {
        values   = [".html", ".htm"]
        operator = "EndsWith"
      }
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "import" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = azurerm_cdn_frontdoor_batch_rule_set.test.name
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_batch_rule_set.test.cdn_frontdoor_profile_id

  rule {
    name  = azurerm_cdn_frontdoor_batch_rule_set.test.rule.0.name
    order = azurerm_cdn_frontdoor_batch_rule_set.test.rule.0.order

    actions {
      route_configuration_override {
        origin_group {
          cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_batch_rule_set.test.rule.0.actions.0.route_configuration_override.0.origin_group.0.cdn_frontdoor_origin_group_id
          forwarding_protocol           = azurerm_cdn_frontdoor_batch_rule_set.test.rule.0.actions.0.route_configuration_override.0.origin_group.0.forwarding_protocol
        }
        caching {
          query_string_behaviour  = azurerm_cdn_frontdoor_batch_rule_set.test.rule.0.actions.0.route_configuration_override.0.caching.0.query_string_behaviour
          query_string_parameters = azurerm_cdn_frontdoor_batch_rule_set.test.rule.0.actions.0.route_configuration_override.0.caching.0.query_string_parameters
          compression_enabled     = azurerm_cdn_frontdoor_batch_rule_set.test.rule.0.actions.0.route_configuration_override.0.caching.0.compression_enabled
          behaviour               = azurerm_cdn_frontdoor_batch_rule_set.test.rule.0.actions.0.route_configuration_override.0.caching.0.behaviour
          duration                = azurerm_cdn_frontdoor_batch_rule_set.test.rule.0.actions.0.route_configuration_override.0.caching.0.duration
        }
      }
    }
  }
}
`, r.basicAttachedRoute(data), data.RandomInteger)
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

func (r CdnFrontdoorBatchRuleSetResource) update(data acceptance.TestData) string {
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
    name               = "acctestBatchRuleExtra%[2]d"
    behaviour_on_match = "Stop"
    order              = 1

    actions {
      route_configuration_override {
        origin_group {
          cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
          forwarding_protocol           = "HttpsOnly"
        }
        caching {
          query_string_behaviour  = "IgnoreSpecifiedQueryStrings"
          query_string_parameters = ["clientIp={client_ip}"]
          compression_enabled     = false
          behaviour               = "OverrideIfOriginMissing"
          duration                = "23:59:59"
        }
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
}
`, r.templateAttachedRoute(data), data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) insertAndReorder(data acceptance.TestData) string {
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
    }
  }

  rule {
    name               = "acctestBatchRuleInserted%[3]d"
    behaviour_on_match = "Continue"
    order              = 2

    actions {
      modify_request_header {
        operator     = "Overwrite"
        header_name  = "X-Test-Inserted"
        header_value = "inserted-rule"
      }
    }
  }

  rule {
    name  = "acctestBatchRuleExtra%[4]d"
    order = 3

    actions {
      modify_request_header {
        operator     = "Overwrite"
        header_name  = "X-Test"
        header_value = "second-rule"
      }
    }
  }
}
`, r.templateAttachedRoute(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) deleteAndReorder(data acceptance.TestData) string {
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
    name               = "acctestBatchRuleInserted%[2]d"
    behaviour_on_match = "Stop"
    order              = 1

    actions {
      route_configuration_override {
        origin_group {
          cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
          forwarding_protocol           = "HttpsOnly"
        }
        caching {
          query_string_behaviour  = "IgnoreSpecifiedQueryStrings"
          query_string_parameters = ["clientIp={client_ip}"]
          compression_enabled     = false
          behaviour               = "OverrideIfOriginMissing"
          duration                = "23:59:59"
        }
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

func (r CdnFrontdoorBatchRuleSetResource) invalidRouteConfigurationOverrideCaching(data acceptance.TestData, cachingBlock string) string {
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
    order = 1

    actions {
      route_configuration_override {
        origin_group {
          cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
          forwarding_protocol           = "HttpsOnly"
        }
%[3]s
      }
    }
  }
}
`, r.templateAttachedRoute(data), data.RandomInteger, cachingBlock)
}

func (r CdnFrontdoorBatchRuleSetResource) urlFilenameConditionOperator(data acceptance.TestData, operator string) string {
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
    behaviour_on_match = "Stop"
    order              = 1

    actions {
      url_rewrite {
        source_pattern                  = "/"
        destination_path                = "/index.html"
        preserve_unmatched_path_enabled = false
      }
    }

    conditions {
      request_filename {
        operator = "%[3]s"
      }
    }
  }
}
`, r.templateAttachedRoute(data), data.RandomInteger, operator)
}

func (r CdnFrontdoorBatchRuleSetResource) conditionAnyOperatorWithValues(data acceptance.TestData) string {
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
    behaviour_on_match = "Stop"
    order              = 1

    actions {
      url_rewrite {
        source_pattern                  = "/"
        destination_path                = "/index.html"
        preserve_unmatched_path_enabled = false
      }
    }

    conditions {
      request_path {
        operator = "Any"
        values   = ["foo", "bar"]
      }
    }
  }
}
`, r.templateAttachedRoute(data), data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) remoteAddressGeoMatchInvalid(data acceptance.TestData) string {
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
    order = 1

    actions {
      url_rewrite {
        source_pattern                  = "/"
        destination_path                = "/index.html"
        preserve_unmatched_path_enabled = false
      }
    }

    conditions {
      remote_address {
        operator = "GeoMatch"
        values   = ["us"]
      }
    }
  }
}
`, r.templateAttachedRoute(data), data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) duplicateRuleName(data acceptance.TestData) string {
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
    order = 1

    actions {
      url_rewrite {
        source_pattern                  = "/"
        destination_path                = "/first.html"
        preserve_unmatched_path_enabled = false
      }
    }
  }

  rule {
    name  = "acctestBatchRule%[3]d"
    order = 2

    actions {
      url_rewrite {
        source_pattern                  = "/"
        destination_path                = "/second.html"
        preserve_unmatched_path_enabled = false
      }
    }
  }
}
`, r.templateAttachedRoute(data), data.RandomInteger, data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) duplicateRuleOrder(data acceptance.TestData) string {
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
    order = 1

    actions {
      url_rewrite {
        source_pattern                  = "/"
        destination_path                = "/first.html"
        preserve_unmatched_path_enabled = false
      }
    }
  }

  rule {
    name  = "acctestBatchRuleExtra%[3]d"
    order = 1

    actions {
      url_rewrite {
        source_pattern                  = "/"
        destination_path                = "/second.html"
        preserve_unmatched_path_enabled = false
      }
    }
  }
}
`, r.templateAttachedRoute(data), data.RandomInteger, data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) unsortedRules(data acceptance.TestData) string {
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
    order = 2

    actions {
      url_rewrite {
        source_pattern                  = "/"
        destination_path                = "/second.html"
        preserve_unmatched_path_enabled = false
      }
    }
  }

  rule {
    name  = "acctestBatchRulePrimary%[3]d"
    order = 1

    actions {
      url_rewrite {
        source_pattern                  = "/"
        destination_path                = "/first.html"
        preserve_unmatched_path_enabled = false
      }
    }
  }
}
`, r.templateAttachedRoute(data), data.RandomInteger, data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) modifyHeaderActionMissingValue(data acceptance.TestData) string {
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
    order = 1

    actions {
      modify_request_header {
        operator    = "Overwrite"
        header_name = "X-Test"
      }
    }
  }
}
`, r.templateAttachedRoute(data), data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) modifyHeaderActionUnexpectedValue(data acceptance.TestData) string {
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
    order = 1

    actions {
      modify_response_header {
        operator     = "Delete"
        header_name  = "X-Test"
        header_value = "should-be-empty"
      }
    }
  }
}
`, r.templateAttachedRoute(data), data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) actionCountValidation(data acceptance.TestData, actionsBlock string) string {
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
    order = 1

    actions {
      %[3]s
    }
  }
}
`, r.templateAttachedRoute(data), data.RandomInteger, actionsBlock)
}

func (r CdnFrontdoorBatchRuleSetResource) gapInRuleOrder(data acceptance.TestData, orders [2]int) string {
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
    name  = "acctestBatchRulePrimary%[2]d"
    order = %[4]d

    actions {
      url_rewrite {
        source_pattern                  = "/"
        destination_path                = "/first.html"
        preserve_unmatched_path_enabled = false
      }
    }
  }

  rule {
    name  = "acctestBatchRuleGap%[3]d"
    order = %[5]d

    actions {
      url_rewrite {
        source_pattern                  = "/"
        destination_path                = "/second.html"
        preserve_unmatched_path_enabled = false
      }
    }
  }
}
`, r.templateAttachedRoute(data), data.RandomInteger, data.RandomInteger, orders[0], orders[1])
}

func (r CdnFrontdoorBatchRuleSetResource) diffQuotaExceeded(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                     = "acctestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  dynamic "rule" {
    for_each = range(51)

    content {
      name  = "acctestRule${rule.key}"
      order = rule.key

      actions {
        route_configuration_override {
          caching {
            behaviour              = "OverrideIfOriginMissing"
            duration               = "365.23:59:59"
            query_string_behaviour = "UseQueryString"
          }
        }
      }
    }
  }

}
`, r.templateAttachedRoute(data), data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) urlRedirect(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  name                     = "acctestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rule {
    name  = "acctestBatchRule%[2]d"
    order = 1

    actions {
      url_redirect {
        redirect_type         = "Found"
        redirect_protocol     = "Https"
        destination_host_name = "contoso.com"
        destination_path      = "/redirected"
        query_string          = "a=b"
      }
    }
  }
}
`, r.templateWithoutOrigin(data), data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) allConditions(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  name                     = "acctestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rule {
    name  = "acctestBatchRule%[2]d"
    order = 1

    actions {
      url_rewrite {
        source_pattern                  = "/"
        destination_path                = "/rewritten"
        preserve_unmatched_path_enabled = true
      }
    }

    conditions {
      remote_address {
        operator = "GeoMatch"
        values   = ["US", "CA"]
      }

      request_scheme {
        operator = "Equal"
        values   = ["HTTP"]
      }

      socket_address {
        operator = "IPMatch"
        values   = ["10.0.0.0/24"]
      }

      query_string {
        operator   = "Contains"
        values     = ["a=b"]
        transforms = ["Lowercase"]
      }

      post_argument {
        name     = "arg"
        operator = "Equal"
        values   = ["value"]
      }

      request_url {
        operator = "BeginsWith"
        values   = ["https://contoso.com"]
      }

      request_header {
        name     = "X-Test"
        operator = "Equal"
        values   = ["value"]
      }

      request_body {
        operator = "Contains"
        values   = ["body"]
      }

      request_file_extension {
        operator = "Equal"
        values   = ["html"]
      }

      http_version {
        operator = "Equal"
        values   = ["2.0"]
      }
    }
  }

  rule {
    name  = "acctestBatchRuleExtra%[2]d"
    order = 2

    actions {
      modify_request_header {
        operator     = "Overwrite"
        header_name  = "X-Test"
        header_value = "value"
      }
    }

    conditions {
      request_cookies {
        name     = "cookie"
        operator = "Equal"
        values   = ["value"]
      }

      client_port {
        operator = "Equal"
        values   = ["8080"]
      }

      server_port {
        operator = "Equal"
        values   = ["443"]
      }

      ssl_protocol {
        operator = "Equal"
        values   = ["TLSv1.2"]
      }

      device_type {
        operator = "NotEqual"
        values   = ["Mobile"]
      }

      host_name {
        operator   = "NotEqual"
        values     = ["www.contoso.com"]
        transforms = ["Lowercase"]
      }

      request_path {
        operator = "NotBeginsWith"
        values   = ["data"]
      }

      request_method {
        operator = "NotEqual"
        values   = ["GET"]
      }
    }
  }
}
`, r.templateWithoutOrigin(data), data.RandomInteger)
}

func (r CdnFrontdoorBatchRuleSetResource) allConditionsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  name                     = "acctestBatchRuleSet%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  rule {
    name  = "acctestBatchRule%[2]d"
    order = 10

    actions {
      url_rewrite {
        source_pattern                  = "/"
        destination_path                = "/rewritten"
        preserve_unmatched_path_enabled = true
      }
    }

    conditions {
      remote_address {
        operator = "NotIPMatch"
        values   = ["10.0.0.0/24"]
      }

      request_scheme {
        operator = "NotEqual"
        values   = ["HTTPS"]
      }

      socket_address {
        operator = "NotIPMatch"
        values   = ["10.0.1.0/24"]
      }

      query_string {
        operator = "Any"
      }

      post_argument {
        name     = "arg"
        operator = "NotEqual"
        values   = ["foo", "bar"]
      }

      request_url {
        operator = "EndsWith"
        values   = [".ca"]
      }

      request_header {
        name       = "X-TestUpdate"
        operator   = "NotEqual"
        values     = ["value", "value2"]
        transforms = ["Lowercase", "RemoveNulls"]
      }

      request_body {
        operator   = "NotBeginsWith"
        values     = ["HELLO"]
        transforms = ["Uppercase"]
      }

      request_file_extension {
        operator = "NotEqual"
        values   = ["php"]
      }

      http_version {
        operator = "NotEqual"
        values   = ["1.0"]
      }
    }
  }

  rule {
    name  = "acctestBatchRuleExtra%[2]d"
    order = 20

    actions {
      modify_request_header {
        operator     = "Append"
        header_name  = "X-TestUpdate"
        header_value = "valueUpdate"
      }
    }

    conditions {
      request_cookies {
        name       = "cookie"
        operator   = "NotEqual"
        values     = ["hello", "world"]
        transforms = ["Lowercase"]
      }

      client_port {
        operator = "NotEqual"
        values   = ["8080", "8181"]
      }

      server_port {
        operator = "NotEqual"
        values   = ["80"]
      }

      ssl_protocol {
        operator = "NotEqual"
        values   = ["TLSv1"]
      }

      device_type {
        operator = "Equal"
        values   = ["Desktop"]
      }

      host_name {
        operator   = "Equal"
        values     = ["www.google.com"]
        transforms = ["RemoveNulls"]
      }

      request_path {
        operator   = "EndsWith"
        values     = ["foo", "bar", "hello", "world"]
        transforms = ["Lowercase", "UrlDecode", "RemoveNulls", "Trim"]
      }

      request_method {
        operator = "Equal"
        values   = ["DELETE", "GET", "HEAD", "OPTIONS", "POST", "PUT", "TRACE"]
      }
    }
  }
}
`, r.templateWithoutOrigin(data), data.RandomInteger)
}
