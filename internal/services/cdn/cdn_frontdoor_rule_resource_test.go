// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CdnFrontDoorRuleResource struct{}

func TestAccCdnFrontDoorRule_basic_unattachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_basic_attachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_cacheDuration_unattachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #22668
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cacheDuration(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_cacheDuration_attachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #22668
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cacheDuration(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_cacheDurationZero_unattachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #23376
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cacheDurationZero(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_cacheDurationZero_attachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #23376
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cacheDurationZero(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_urlRedirectAction_unattachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #18249
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.urlRedirectAction(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_urlRedirectAction_attachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #18249
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.urlRedirectAction(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_originGroupIdOptional_unattachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #18889
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.originGroupIdOptional(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_originGroupIdOptional_attachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #18889
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.originGroupIdOptional(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_originGroupIdOptionalUpdate_unattachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #18889
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.originGroupIdOptional(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.originGroupIdOptionalUpdate(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.originGroupIdOptional(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_originGroupIdOptionalUpdate_attachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #18889
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.originGroupIdOptional(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.originGroupIdOptionalUpdate(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.originGroupIdOptional(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_disableCache_unattachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #19008
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disableCache(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_disableCache_attachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #19008
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disableCache(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_disableCacheOriginGroupId_unattachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #19008
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disableCacheOriginGroupId(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_disableCacheOriginGroupId_attachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #19008
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disableCacheOriginGroupId(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_disableCacheOriginGroupIdUpdate_unattachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #19008
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disableCacheOriginGroupId(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.enableCacheOriginGroupId(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.disableCacheOriginGroupId(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_disableCacheOriginGroupIdUpdate_attachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #19008
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disableCacheOriginGroupId(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.enableCacheOriginGroupId(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.disableCacheOriginGroupId(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_disableCacheUpdate_unattachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #19008
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disableCache(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.enableCache(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.disableCache(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_disableCacheUpdate_attachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #19008
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disableCache(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.enableCache(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.disableCache(data, true),
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
			ExpectError: regexp.MustCompile("when `route_configuration_override.caching.behaviour` is set to `Disabled`, you cannot define `route_configuration_override.caching.duration`"),
		},
	})
}

func TestAccCdnFrontDoorRule_actionOnly_unattachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.actionOnly(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_actionOnly_attachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.actionOnly(data, true),
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
			Config: r.basic(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccCdnFrontDoorRule_importBatchRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: CdnFrontdoorBatchRuleSetResource{}.disableCacheAndNoOriginGroup(data),
		},
		{
			Config:      r.batchRuleImport(data),
			ExpectError: regexp.MustCompile("was provisioned using batch mode, and individual rules for this cannot be managed by this resource"),
		},
	})
}

func TestAccCdnFrontDoorRule_complete_unattachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_complete_attachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_update_unattachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_update_attachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, true),
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

func TestAccCdnFrontDoorRule_multipleQueryStringParameters_unattachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #19097
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleQueryStringParameters(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_multipleQueryStringParameters_attachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #19097
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleQueryStringParameters(data, true),
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

func TestAccCdnFrontDoorRule_honorOrigin_unattachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #19311
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.honorOrigin(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_honorOrigin_attachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #19311
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.honorOrigin(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_allowForwardSlashUrlConditionMatchValue_unattachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.allowForwardSlashUrlConditionMatchValue(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_allowForwardSlashUrlConditionMatchValue_attachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.allowForwardSlashUrlConditionMatchValue(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_allowForwardSlashUrl2ConditionMatchValue_unattachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.allowForwardSlashUrl2ConditionMatchValue(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_allowForwardSlashUrl2ConditionMatchValue_attachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.allowForwardSlashUrl2ConditionMatchValue(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_urlFilenameConditionOperatorAny_unattachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #23504
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.urlFilenameConditionOperator(data, "Any", false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("conditions.0.request_filename.0.operator").HasValue("Any"),
				check.That(data.ResourceName).Key("conditions.0.request_filename.0.values").DoesNotExist(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_urlFilenameConditionOperatorAny_attachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #23504
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.urlFilenameConditionOperator(data, "Any", true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("conditions.0.request_filename.0.operator").HasValue("Any"),
				check.That(data.ResourceName).Key("conditions.0.request_filename.0.values").DoesNotExist(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_urlFilenameConditionOperatorError(t *testing.T) {
	// NOTE: Regression test case for issue #23504
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.urlFilenameConditionOperator(data, "Contains", false),
			ExpectError: regexp.MustCompile("when `conditions.request_filename.operator` is set to `Contains`, `conditions.request_filename.values` must set one or more values"),
		},
	})
}

func TestAccCdnFrontDoorRule_urlPathConditionOperatorWildcard_unattachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #29415
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.urlPathWildcard(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.urlPathWildcardNegate(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.urlPathWildcard(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRule_urlPathConditionOperatorWildcard_attachedRoute(t *testing.T) {
	// NOTE: Regression test case for issue #29415
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontDoorRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.urlPathWildcard(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.urlPathWildcardNegate(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.urlPathWildcard(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CdnFrontDoorRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := rules.ParseRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cdn.FrontDoorRulesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
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

func (r CdnFrontDoorRuleResource) templateWithAttachedRoute(data acceptance.TestData, attachRoute bool) string {
	template := r.template(data)
	if !attachRoute {
		return template
	}

	return fmt.Sprintf(`%s

%s`, template, r.routeTemplate(data))
}

func (r CdnFrontDoorRuleResource) routeTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
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
`, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) basic(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
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
    route_configuration_override {
      caching {
        query_string_behaviour  = "IncludeSpecifiedQueryStrings"
        query_string_parameters = ["foo", "clientIp={client_ip}"]
        compression_enabled     = true
        behaviour               = "OverrideIfOriginMissing"
        duration                = "365.23:59:59"
      }

      origin_group {
        cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
        forwarding_protocol           = "HttpsOnly"
      }
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) cacheDuration(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
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
    route_configuration_override {
      caching {
        query_string_behaviour  = "IncludeSpecifiedQueryStrings"
        query_string_parameters = ["foo", "clientIp={client_ip}"]
        compression_enabled     = true
        behaviour               = "OverrideIfOriginMissing"
        duration                = "167.23:59:59"
      }

      origin_group {
        cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
        forwarding_protocol           = "HttpsOnly"
      }
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) cacheDurationZero(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
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
    route_configuration_override {
      caching {
        query_string_behaviour  = "IncludeSpecifiedQueryStrings"
        query_string_parameters = ["foo", "clientIp={client_ip}"]
        compression_enabled     = true
        behaviour               = "OverrideIfOriginMissing"
        duration                = "00:00:00"
      }

      origin_group {
        cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
        forwarding_protocol           = "HttpsOnly"
      }
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) urlRedirectAction(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
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
    request_scheme {
      values   = ["HTTP"]
      operator = "Equal"
    }
  }

  actions {
    url_redirect {
      redirect_type     = "PermanentRedirect"
      redirect_protocol = "Https"
      #destination_host_name = ""
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) originGroupIdOptional(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
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
    request_path {
      operator   = "RegEx"
      values     = ["api/?(.*)"]
      transforms = ["Lowercase", "Trim"]
    }
  }

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
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) originGroupIdOptionalUpdate(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
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
    request_path {
      operator   = "RegEx"
      values     = ["api/?(.*)"]
      transforms = ["Lowercase", "Trim"]
    }
  }

  actions {
    route_configuration_override {
      caching {
        query_string_behaviour  = "IncludeSpecifiedQueryStrings"
        query_string_parameters = ["foo", "clientIp={client_ip}"]
        compression_enabled     = true
        behaviour               = "OverrideIfOriginMissing"
        duration                = "365.23:59:59"
      }

      origin_group {
        cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
        forwarding_protocol           = "HttpsOnly"
      }
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) disableCache(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
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
    request_path {
      operator   = "RegEx"
      values     = ["api/?(.*)"]
      transforms = ["Lowercase", "Trim"]
    }
  }

  actions {
    route_configuration_override {
      caching {
        behaviour = "Disabled"
      }
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) disableCacheOriginGroupId(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
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
    request_path {
      operator   = "RegEx"
      values     = ["api/?(.*)"]
      transforms = ["Lowercase", "Trim"]
    }
  }

  actions {
    route_configuration_override {
      caching {
        behaviour = "Disabled"
      }

      origin_group {
        cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
        forwarding_protocol           = "HttpsOnly"
      }
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) enableCacheOriginGroupId(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
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
    request_path {
      operator   = "RegEx"
      values     = ["api/?(.*)"]
      transforms = ["Lowercase", "Trim"]
    }
  }

  actions {
    route_configuration_override {
      caching {
        query_string_behaviour  = "IncludeSpecifiedQueryStrings"
        query_string_parameters = ["RUSH", "clientIp={client_ip}"]
        compression_enabled     = true
        behaviour               = "OverrideIfOriginMissing"
        duration                = "21.12:04:01"
      }

      origin_group {
        cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
        forwarding_protocol           = "HttpsOnly"
      }
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) enableCache(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
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
    request_path {
      operator   = "RegEx"
      values     = ["api/?(.*)"]
      transforms = ["Lowercase", "Trim"]
    }
  }

  actions {
    route_configuration_override {
      caching {
        query_string_behaviour  = "IncludeSpecifiedQueryStrings"
        query_string_parameters = ["RUSH", "clientIp={client_ip}"]
        compression_enabled     = true
        behaviour               = "OverrideIfOriginMissing"
        duration                = "21.12:04:01"
      }
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
    request_path {
      operator   = "RegEx"
      values     = ["api/?(.*)"]
      transforms = ["Lowercase", "Trim"]
    }
  }

  actions {
    route_configuration_override {
      caching {
        behaviour = "Disabled"
        duration  = "365.23:59:59"
      }

    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data, true)
	return fmt.Sprintf(`
			%s

resource "azurerm_cdn_frontdoor_rule" "import" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = azurerm_cdn_frontdoor_rule.test.name
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

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

      origin_group {
        cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
        forwarding_protocol           = "HttpsOnly"
      }
    }
  }

}
`, config)
}

func (r CdnFrontDoorRuleResource) batchRuleImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

import {
  id = "${azurerm_cdn_frontdoor_batch_rule_set.test.id}/rules/${azurerm_cdn_frontdoor_batch_rule_set.test.rule.0.name}"
  to = azurerm_cdn_frontdoor_rule.test
}

resource "azurerm_cdn_frontdoor_rule" "test" {
  name                      = "accTestRule%[2]d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_batch_rule_set.test.id

  order = 0

  actions {
    route_configuration_override {
      caching {
        behaviour = "Disabled"
      }
    }
  }
}
`, CdnFrontdoorBatchRuleSetResource{}.basicUnattachedRoute(data), data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) complete(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id
  behaviour_on_match        = "Continue"
  order                     = 1

  actions {
    route_configuration_override {
      caching {
        query_string_behaviour  = "IncludeSpecifiedQueryStrings"
        query_string_parameters = ["foo", "clientIp={client_ip}"]
        compression_enabled     = true
        behaviour               = "OverrideIfOriginMissing"
        duration                = "365.23:59:59"
      }

      origin_group {
        cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
        forwarding_protocol           = "HttpsOnly"
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

    post_argument {
      name       = "customerName"
      operator   = "BeginsWith"
      values     = ["J", "K"]
      transforms = ["Uppercase"]
    }

    request_method {
      operator = "Equal"
      values   = ["DELETE"]
    }

    request_filename {
      operator   = "Equal"
      values     = ["media.mp4"]
      transforms = ["Lowercase", "RemoveNulls", "Trim"]
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) update(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.test, azurerm_cdn_frontdoor_origin.test]

  name                      = "accTestRule%d"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id
  behaviour_on_match        = "Stop"
  order                     = 2

  actions {
    route_configuration_override {
      caching {
        query_string_behaviour  = "IgnoreSpecifiedQueryStrings"
        query_string_parameters = ["clientIp={client_ip}"]
        compression_enabled     = false
        behaviour               = "OverrideIfOriginMissing"
        duration                = "23:59:59"
      }

      origin_group {
        cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
        forwarding_protocol           = "HttpsOnly"
      }
    }
  }

  conditions {
    host_name {
      operator   = "NotEqual"
      values     = ["www.contoso.com", "images.contoso.com", "video.contoso.com"]
      transforms = ["Lowercase", "Trim"]
    }

    device_type {
      operator = "NotEqual"
      values   = ["Mobile"]
    }

    post_argument {
      name       = "customerName"
      operator   = "BeginsWith"
      values     = ["J", "K"]
      transforms = ["Uppercase"]
    }

    request_method {
      operator = "Equal"
      values   = ["DELETE"]
    }

    request_filename {
      operator   = "NotEqual"
      values     = ["mp3"]
      transforms = ["Lowercase"]
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) actionOnly(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
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
    route_configuration_override {
      caching {
        query_string_behaviour  = "IgnoreSpecifiedQueryStrings"
        query_string_parameters = ["clientIp={client_ip}"]
        compression_enabled     = false
        behaviour               = "OverrideIfOriginMissing"
        duration                = "23:59:59"
      }

      origin_group {
        cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
        forwarding_protocol           = "HttpsOnly"
      }
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
    route_configuration_override {
      caching {
        query_string_behaviour  = "IgnoreSpecifiedQueryStrings"
        query_string_parameters = ["clientIp={client_ip}"]
        compression_enabled     = false
        behaviour               = "OverrideIfOriginMissing"
        duration                = "0.23:59:59"
      }

      origin_group {
        cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
        forwarding_protocol           = "HttpsOnly"
      }
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) multipleQueryStringParameters(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
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
    request_url {
      values   = ["https://contoso.com/test"]
      operator = "Equal"
    }
  }

  actions {
    url_redirect {
      redirect_type         = "PermanentRedirect"
      redirect_protocol     = "Https"
      query_string          = "TE=&PFalse=&source=TE&medium=mai&campaign=y10"
      destination_host_name = "contoso.com"
      destination_path      = "/test/page"
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
    request_url {
      values   = ["https://contoso.com/test"]
      operator = "Equal"
    }
  }

  actions {
    url_redirect {
      redirect_type         = "PermanentRedirect"
      redirect_protocol     = "Https"
      query_string          = "origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.com&destination_host=fabrikam.com&redirect_from=frontdoor&origin_host=contoso.c"
      destination_host_name = "fabrikam.com"
      destination_path      = "/test/page"
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) honorOrigin(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
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
      values   = ["data/", ]
      operator = "BeginsWith"
    }

    request_path {
      values   = [".html", ".htm"]
      operator = "EndsWith"
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) allowForwardSlashUrlConditionMatchValue(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
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
      values   = ["/"]
      operator = "EndsWith"
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) allowForwardSlashUrl2ConditionMatchValue(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
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
      values   = ["/legacy-login"]
      operator = "EndsWith"
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) urlFilenameConditionOperator(data acceptance.TestData, operator string, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
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
  behaviour_on_match        = "Stop"

  actions {
    url_rewrite {
      source_pattern                  = "/"
      destination_path                = "/index.html"
      preserve_unmatched_path_enabled = false
    }
  }

  conditions {
    request_filename {
      operator = "%s"
    }
  }
}
`, template, data.RandomInteger, operator)
}

func (r CdnFrontDoorRuleResource) urlPathWildcard(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
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
    request_path {
      operator   = "Wildcard"
      values     = ["files/customer*/file.pdf"]
      transforms = ["Lowercase", "Trim"]
    }
  }

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
`, template, data.RandomInteger)
}

func (r CdnFrontDoorRuleResource) urlPathWildcardNegate(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
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
    request_path {
      operator   = "NotWildcard"
      values     = ["files/customer*/file.pdf"]
      transforms = ["Lowercase", "Trim"]
    }
  }

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
`, template, data.RandomInteger)
}
