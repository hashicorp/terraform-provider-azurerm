// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontDoorRuleResource struct{}

func TestAccCdnFrontDoorRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
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

func TestAccCdnFrontDoorRule_cacheDuration(t *testing.T) {
	// NOTE: Regression test case for issue #22668
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cacheDuration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_urlRedirectAction(t *testing.T) {
	// NOTE: Regression test case for issue #18249
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.urlRedirectAction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_originGroupIdOptional(t *testing.T) {
	// NOTE: Regression test case for issue #18889
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.originGroupIdOptional(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_originGroupIdOptionalUpdate(t *testing.T) {
	// NOTE: Regression test case for issue #18889
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.originGroupIdOptional(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.originGroupIdOptionalUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.originGroupIdOptional(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_originGroupIdOptionalError(t *testing.T) {
	// NOTE: Regression test case for issue #18889
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.originGroupIdOptionalError(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExpectError: regexp.MustCompile("the 'route_configuration_override_action' block is not valid, if the 'cdn_frontdoor_origin_group_id' is not set you cannot define the 'forwarding_protocol'"),
		},
	})
}

func TestAccCdnFrontDoorRule_disableCache(t *testing.T) {
	// NOTE: Regression test case for issue #19008
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disableCache(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_disableCacheOriginGroupId(t *testing.T) {
	// NOTE: Regression test case for issue #19008
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disableCacheOriginGroupId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_disableCacheOriginGroupIdUpdate(t *testing.T) {
	// NOTE: Regression test case for issue #19008
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disableCacheOriginGroupId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.enableCacheOriginGroupId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.disableCacheOriginGroupId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_disableCacheUpdate(t *testing.T) {
	// NOTE: Regression test case for issue #19008
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disableCache(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.enableCache(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.disableCache(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_disableCacheError(t *testing.T) {
	// NOTE: Regression test case for issue #19008
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disableCacheError(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExpectError: regexp.MustCompile("the 'route_configuration_override_action' block is not valid, if the 'cache_behavior' is set to 'Disabled' you cannot define the 'cache_duration'"),
		},
	})
}

func TestAccCdnFrontDoorRule_actionOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.actionOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
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

func TestAccCdnFrontDoorRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
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

func TestAccCdnFrontDoorRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
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
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_invalidCacheDuration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.invalidCacheDuration(data),
			ExpectError: regexp.MustCompile(`if the duration is less than 1`),
		},
	})
}

func TestAccCdnFrontDoorRule_multipleQueryStringParameters(t *testing.T) {
	// NOTE: Regression test case for issue #19097
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleQueryStringParameters(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_multipleQueryStringParametersError(t *testing.T) {
	// NOTE: Regression test case for issue #19097
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.multipleQueryStringParametersError(data),
			ExpectError: regexp.MustCompile(`cannot be longer than 2048 characters in length`),
		},
	})
}

func TestAccCdnFrontDoorRule_honorOrigin(t *testing.T) {
	// NOTE: Regression test case for issue #19311
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.honorOrigin(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_allowEmptyQueryString(t *testing.T) {
	// NOTE: Regression test case for issue #19682
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.allowEmptyQueryString(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_allowForwardSlashUrlConditionMatchValue(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.allowForwardSlashUrlConditionMatchValue(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CdnFrontDoorRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FrontDoorRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontDoorRulesClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.RuleSetName, id.RuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r CdnFrontDoorRuleResource) template(data acceptance.TestData) string {
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
  name                     = "accTestRuleSet%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CdnFrontDoorRuleResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

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
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) cacheDuration(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  order = 0

  actions {
    route_configuration_override_action {
      cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
      forwarding_protocol           = "HttpsOnly"
      query_string_caching_behavior = "IncludeSpecifiedQueryStrings"
      query_string_parameters       = ["foo", "clientIp={client_ip}"]
      compression_enabled           = true
      cache_behavior                = "OverrideIfOriginMissing"
      cache_duration                = "167.23:59:59"
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) urlRedirectAction(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  order = 1

  conditions {
    request_scheme_condition {
      match_values     = ["HTTP"]
      negate_condition = false
      operator         = "Equal"
    }
  }

  actions {
    url_redirect_action {
      redirect_type        = "PermanentRedirect"
      redirect_protocol    = "Https"
      destination_hostname = ""
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) originGroupIdOptional(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  order = 0

  conditions {
    url_path_condition {
      operator         = "RegEx"
      negate_condition = false
      match_values     = ["api/?(.*)"]
      transforms       = ["Lowercase", "Trim"]
    }
  }

  actions {
    route_configuration_override_action {
      query_string_caching_behavior = "IncludeSpecifiedQueryStrings"
      query_string_parameters       = ["foo", "clientIp={client_ip}"]
      compression_enabled           = true
      cache_behavior                = "OverrideIfOriginMissing"
      cache_duration                = "365.23:59:59"
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) originGroupIdOptionalUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  order = 0

  conditions {
    url_path_condition {
      operator         = "RegEx"
      negate_condition = false
      match_values     = ["api/?(.*)"]
      transforms       = ["Lowercase", "Trim"]
    }
  }

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
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) originGroupIdOptionalError(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  order = 0

  conditions {
    url_path_condition {
      operator         = "RegEx"
      negate_condition = false
      match_values     = ["api/?(.*)"]
      transforms       = ["Lowercase", "Trim"]
    }
  }

  actions {
    route_configuration_override_action {
      forwarding_protocol           = "HttpsOnly"
      query_string_caching_behavior = "IncludeSpecifiedQueryStrings"
      query_string_parameters       = ["foo", "clientIp={client_ip}"]
      compression_enabled           = true
      cache_behavior                = "OverrideIfOriginMissing"
      cache_duration                = "365.23:59:59"
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) disableCache(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  order = 0

  conditions {
    url_path_condition {
      operator         = "RegEx"
      negate_condition = false
      match_values     = ["api/?(.*)"]
      transforms       = ["Lowercase", "Trim"]
    }
  }

  actions {
    route_configuration_override_action {
      cache_behavior = "Disabled"
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) disableCacheOriginGroupId(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  order = 0

  conditions {
    url_path_condition {
      operator         = "RegEx"
      negate_condition = false
      match_values     = ["api/?(.*)"]
      transforms       = ["Lowercase", "Trim"]
    }
  }

  actions {
    route_configuration_override_action {
      cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
      forwarding_protocol           = "HttpsOnly"
      cache_behavior                = "Disabled"
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) enableCacheOriginGroupId(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  order = 0

  conditions {
    url_path_condition {
      operator         = "RegEx"
      negate_condition = false
      match_values     = ["api/?(.*)"]
      transforms       = ["Lowercase", "Trim"]
    }
  }

  actions {
    route_configuration_override_action {
      cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
      forwarding_protocol           = "HttpsOnly"
      query_string_caching_behavior = "IncludeSpecifiedQueryStrings"
      query_string_parameters       = ["RUSH", "clientIp={client_ip}"]
      compression_enabled           = true
      cache_behavior                = "OverrideIfOriginMissing"
      cache_duration                = "21.12:04:01"
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) enableCache(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  order = 0

  conditions {
    url_path_condition {
      operator         = "RegEx"
      negate_condition = false
      match_values     = ["api/?(.*)"]
      transforms       = ["Lowercase", "Trim"]
    }
  }

  actions {
    route_configuration_override_action {
      query_string_caching_behavior = "IncludeSpecifiedQueryStrings"
      query_string_parameters       = ["RUSH", "clientIp={client_ip}"]
      compression_enabled           = true
      cache_behavior                = "OverrideIfOriginMissing"
      cache_duration                = "21.12:04:01"
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) disableCacheError(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  order = 0

  conditions {
    url_path_condition {
      operator         = "RegEx"
      negate_condition = false
      match_values     = ["api/?(.*)"]
      transforms       = ["Lowercase", "Trim"]
    }
  }

  actions {
    route_configuration_override_action {
      cache_behavior = "Disabled"
      cache_duration = "365.23:59:59"
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cdn_frontdoor_rule" "import" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = azurerm_cdn_frontdoor_rule.test.name
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

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
`, config)
}

func (r CdnFrontDoorRuleResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id
  behavior_on_match         = "Continue"
  order                     = 1

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

    url_redirect_action {
      redirect_type        = "PermanentRedirect"
      redirect_protocol    = "MatchRequest"
      query_string         = "clientIp={client_ip}"
      destination_path     = "/exampleredirection"
      destination_hostname = "contoso.com"
      destination_fragment = "UrlRedirect"
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

    post_args_condition {
      post_args_name = "customerName"
      operator       = "BeginsWith"
      match_values   = ["J", "K"]
      transforms     = ["Uppercase"]
    }

    request_method_condition {
      operator         = "Equal"
      negate_condition = false
      match_values     = ["DELETE"]
    }

    url_filename_condition {
      operator         = "Equal"
      negate_condition = false
      match_values     = ["media.mp4"]
      transforms       = ["Lowercase", "RemoveNulls", "Trim"]
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id
  behavior_on_match         = "Stop"
  order                     = 2

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

    post_args_condition {
      post_args_name = "customerName"
      operator       = "BeginsWith"
      match_values   = ["J", "K"]
      transforms     = ["Uppercase"]
    }

    request_method_condition {
      operator         = "Equal"
      negate_condition = false
      match_values     = ["DELETE"]
    }

    url_filename_condition {
      operator         = "Equal"
      negate_condition = false
      match_values     = ["media.mp4"]
      transforms       = ["Lowercase"]
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) actionOnly(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id
  order                     = 1

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
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) invalidCacheDuration(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id
  order                     = 1

  actions {
    route_configuration_override_action {
      cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
      forwarding_protocol           = "HttpsOnly"
      query_string_caching_behavior = "IgnoreSpecifiedQueryStrings"
      query_string_parameters       = ["clientIp={client_ip}"]
      compression_enabled           = false
      cache_behavior                = "OverrideIfOriginMissing"
      cache_duration                = "0.23:59:59"
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) multipleQueryStringParameters(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  order = 0

  conditions {
    request_uri_condition {
      match_values     = ["https://contoso.com/test"]
      negate_condition = false
      operator         = "Equal"
    }
  }

  actions {
    url_redirect_action {
      redirect_type        = "PermanentRedirect"
      redirect_protocol    = "Https"
      query_string         = "TE=&PFalse=&source=TE&medium=mai&campaign=y10"
      destination_hostname = "contoso.com"
      destination_path     = "/test/page"
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) multipleQueryStringParametersError(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  order = 0

  conditions {
    request_uri_condition {
      match_values     = ["https://contoso.com/test"]
      negate_condition = false
      operator         = "Equal"
    }
  }

  actions {
    url_redirect_action {
      redirect_type        = "PermanentRedirect"
      redirect_protocol    = "Https"
      query_string         = "origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.c"
      destination_hostname = "fabrikam.com"
      destination_path     = "/test/page"
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) honorOrigin(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  order = 0

  actions {
    route_configuration_override_action {
      cache_behavior                = "HonorOrigin"
      compression_enabled           = true
      query_string_caching_behavior = "IgnoreQueryString"
    }
  }

  conditions {
    url_path_condition {
      match_values     = ["data/", ]
      negate_condition = false
      operator         = "BeginsWith"
    }

    url_path_condition {
      match_values     = [".html", ".htm"]
      negate_condition = false
      operator         = "EndsWith"
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) allowEmptyQueryString(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

  %s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  order = 0

  conditions {
    request_uri_condition {
      match_values     = ["contoso"]
      negate_condition = false
      operator         = "Contains"
    }
  }

  actions {
    url_redirect_action {
      redirect_type        = "PermanentRedirect"
      redirect_protocol    = "MatchRequest"
      query_string         = ""
      destination_hostname = "contoso.com"
      destination_path     = "/test/page"
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) allowForwardSlashUrlConditionMatchValue(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  order = 0

  actions {
    route_configuration_override_action {
      cache_behavior                = "HonorOrigin"
      compression_enabled           = true
      query_string_caching_behavior = "IgnoreQueryString"
    }
  }

  conditions {
    url_path_condition {
      match_values     = ["/"]
      negate_condition = false
      operator         = "EndsWith"
    }
  }
}
`, template, data.RandomInteger)
}
