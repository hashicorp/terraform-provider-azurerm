// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SecurityInsightsIndicatorResource struct{}

func TestAccSecurityInsightsIndicator_basicDomainName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_threat_intelligence_indicator", "test")
	r := SecurityInsightsIndicatorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "domain-name", "http://test.com"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityInsightsIndicator_basicFile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_threat_intelligence_indicator", "test")
	r := SecurityInsightsIndicatorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "file", "MD5:78ecc5c05cd8b79af480df2f8fba0b9d"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityInsightsIndicator_basicIpV4(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_threat_intelligence_indicator", "test")
	r := SecurityInsightsIndicatorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "ipv4-addr", "1.1.1.1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityInsightsIndicator_basicIpV6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_threat_intelligence_indicator", "test")
	r := SecurityInsightsIndicatorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "ipv6-addr", "::1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityInsightsIndicator_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_threat_intelligence_indicator", "test")
	r := SecurityInsightsIndicatorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "domain-name", "http://test.com"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSecurityInsightsIndicator_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_threat_intelligence_indicator", "test")
	r := SecurityInsightsIndicatorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "domain-name", "http://test.com"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityInsightsIndicator_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_threat_intelligence_indicator", "test")
	r := SecurityInsightsIndicatorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "domain-name", "http://example.com"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, "domain-name", "http://example.com"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SecurityInsightsIndicatorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ThreatIntelligenceIndicatorID(state.ID)
	if err != nil {
		return nil, err
	}

	client := azuresdkhacks.ThreatIntelligenceIndicatorClient{
		BaseClient: clients.Sentinel.ThreatIntelligenceClient.BaseClient,
	}
	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.IndicatorName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Value != nil), nil
}

func (r SecurityInsightsIndicatorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "test" {
  workspace_id = azurerm_log_analytics_workspace.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r SecurityInsightsIndicatorResource) basic(data acceptance.TestData, patternType string, pattern string) string {
	return fmt.Sprintf(`
				%s

resource "azurerm_sentinel_threat_intelligence_indicator" "test" {
  workspace_id      = azurerm_log_analytics_workspace.test.id
  pattern_type      = "%s"
  pattern           = "%s"
  source            = "Microsoft Sentinel"
  validate_from_utc = "2022-12-14T16:00:00Z"
  display_name      = "test"

  depends_on = [azurerm_sentinel_log_analytics_workspace_onboarding.test]
}
`, r.template(data), patternType, pattern)
}

func (r SecurityInsightsIndicatorResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data, "domain-name", "http://test.com")
	return fmt.Sprintf(`
			%s

resource "azurerm_sentinel_threat_intelligence_indicator" "import" {
  workspace_id      = azurerm_sentinel_threat_intelligence_indicator.test.workspace_id
  pattern_type      = azurerm_sentinel_threat_intelligence_indicator.test.pattern_type
  pattern           = azurerm_sentinel_threat_intelligence_indicator.test.pattern
  source            = azurerm_sentinel_threat_intelligence_indicator.test.source
  validate_from_utc = azurerm_sentinel_threat_intelligence_indicator.test.validate_from_utc
  display_name      = azurerm_sentinel_threat_intelligence_indicator.test.display_name
}
`, config)
}

func (r SecurityInsightsIndicatorResource) complete(data acceptance.TestData, patternType string, pattern string) string {
	return fmt.Sprintf(`
			%s

resource "azurerm_sentinel_threat_intelligence_indicator" "test" {
  workspace_id    = azurerm_log_analytics_workspace.test.id
  pattern_type    = "%s"
  pattern         = "%s"
  confidence      = 5
  created_by      = "testcraeted@microsoft.com"
  description     = "test indicator"
  display_name    = "test"
  language        = "en"
  pattern_version = 1
  revoked         = true
  tags            = ["test-tags"]
  kill_chain_phase {
    name = "testtest"
  }
  external_reference {
    description = "test-external"
    source_name = "test-sourcename"
  }
  granular_marking {
    language = "en"
  }
  source            = "test Sentinel"
  validate_from_utc = "2022-12-14T16:00:00Z"

  depends_on = [azurerm_sentinel_log_analytics_workspace_onboarding.test]
}
`, r.template(data), patternType, pattern)
}

func (r SecurityInsightsIndicatorResource) update(data acceptance.TestData, patternType string, pattern string) string {
	return fmt.Sprintf(`
			%s

resource "azurerm_sentinel_threat_intelligence_indicator" "test" {
  workspace_id    = azurerm_log_analytics_workspace.test.id
  pattern_type    = "%s"
  pattern         = "%s"
  confidence      = 5
  created_by      = "testcraeted@microsoft.com"
  description     = "updated indicator"
  display_name    = "updated"
  language        = "en"
  pattern_version = 1
  revoked         = true
  tags            = ["updated-tags"]
  kill_chain_phase {
    name = "testtest"
  }
  external_reference {
    description = "test-external"
    source_name = "test-sourcename"
  }
  granular_marking {
    language = "en"
  }
  source            = "updated Sentinel"
  validate_from_utc = "2022-12-15T16:00:00Z"

  depends_on = [azurerm_sentinel_log_analytics_workspace_onboarding.test]
}
`, r.template(data), patternType, pattern)
}
