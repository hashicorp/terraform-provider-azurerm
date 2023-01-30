package sentinel_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
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
			Config: r.basic(data, "domain-name", "test.com"),
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

//
//func TestAccSecurityInsightsIndicator_requiresImport(t *testing.T) {
//	data := acceptance.BuildTestData(t, "azurerm_sentinel_threat_intelligence_indicator", "test")
//	r := SecurityInsightsIndicatorResource{}
//	data.ResourceTest(t, r, []acceptance.TestStep{
//		{
//			Config: r.basic(data, "domain-name", "test.com"),
//			Check: acceptance.ComposeTestCheckFunc(
//				check.That(data.ResourceName).ExistsInAzure(r),
//			),
//		},
//		data.RequiresImportErrorStep(r.requiresImport),
//	})
//}

//
//func TestAccSecurityInsightsIndicator_complete(t *testing.T) {
//	data := acceptance.BuildTestData(t, "azurerm_sentinel_threat_intelligence_indicator", "test")
//	r := SecurityInsightsIndicatorResource{}
//	data.ResourceTest(t, r, []acceptance.TestStep{
//		{
//			Config: r.complete(data),
//			Check: acceptance.ComposeTestCheckFunc(
//				check.That(data.ResourceName).ExistsInAzure(r),
//			),
//		},
//		data.ImportStep(),
//	})
//}
//
//func TestAccSecurityInsightsIndicator_update(t *testing.T) {
//	data := acceptance.BuildTestData(t, "azurerm_sentinel_threat_intelligence_indicator", "test")
//	r := SecurityInsightsIndicatorResource{}
//	data.ResourceTest(t, r, []acceptance.TestStep{
//		{
//			Config: r.complete(data),
//			Check: acceptance.ComposeTestCheckFunc(
//				check.That(data.ResourceName).ExistsInAzure(r),
//			),
//		},
//		data.ImportStep(),
//		{
//			Config: r.update(data),
//			Check: acceptance.ComposeTestCheckFunc(
//				check.That(data.ResourceName).ExistsInAzure(r),
//			),
//		},
//		data.ImportStep(),
//	})
//}

func (r SecurityInsightsIndicatorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ThreatIntelligenceIndicatorID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Sentinel.ThreatIntelligenceClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.IndicatorName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	info, ok := resp.Value.AsThreatIntelligenceIndicatorModel()
	if !ok {
		return nil, fmt.Errorf("converting %s: %+v", id, err)
	}

	return utils.Bool(info.ID != nil), nil
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
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = azurerm_log_analytics_workspace.test.name
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
	config := r.basic(data, "domain-name", "test.com")
	return fmt.Sprintf(`
			%s

resource "azurerm_sentinel_threat_intelligence_indicator" "import" {
  name                = azurerm_sentinel_threat_intelligence_indicator.test.name
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = ""
  kind                = ""

  depends_on = [azurerm_sentinel_log_analytics_workspace_onboarding.test]

}
`, config)
}

func (r SecurityInsightsIndicatorResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
			%s

resource "azurerm_sentinel_threat_intelligence_indicator" "test" {
  name                           = "acctest-sii-%d"
  resource_group_name            = azurerm_resource_group.test.name
  workspace_name                 = ""
  confidence                     = 0
  created                        = ""
  created_by_ref                 = ""
  defanged                       = false
  description                    = ""
  display_name                   = ""
  external_id                    = ""
  external_last_updated_time_utc = ""
  kind                           = ""
  language                       = ""
  last_updated_time_utc          = ""
  modified                       = ""
  pattern                        = ""
  pattern_type                   = ""
  pattern_version                = ""
  revoked                        = false
  source                         = ""
  valid_from                     = ""
  valid_until                    = ""
  indicator_types                = []
  labels                         = []
  object_marking_refs            = []
  threat_intelligence_tags       = []
  threat_types                   = []
  external_references {
    description = ""
    external_id = ""
    source_name = ""
    url         = ""
    hashes = {
      key = ""
    }
  }
  granular_markings {
    language    = ""
    marking_ref = 0
    selectors   = []
  }
  kill_chain_phases {
    kill_chain_name = ""
    phase_name      = ""
  }
  parsed_pattern {
    pattern_type_key = ""
    pattern_type_values {
      value      = ""
      value_type = ""
    }
  }
  extensions = jsonencode({
    "key" : "value"
  })

  depends_on = [azurerm_sentinel_log_analytics_workspace_onboarding.test]

}
`, r.template(data), data.RandomInteger)
}

func (r SecurityInsightsIndicatorResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
			%s

resource "azurerm_sentinel_threat_intelligence_indicator" "test" {
  name                           = "acctest-sii-%d"
  resource_group_name            = azurerm_resource_group.test.name
  workspace_name                 = ""
  confidence                     = 0
  created                        = ""
  created_by_ref                 = ""
  defanged                       = false
  description                    = ""
  display_name                   = ""
  external_id                    = ""
  external_last_updated_time_utc = ""
  kind                           = ""
  language                       = ""
  last_updated_time_utc          = ""
  modified                       = ""
  pattern                        = ""
  pattern_type                   = ""
  pattern_version                = ""
  revoked                        = false
  source                         = ""
  valid_from                     = ""
  valid_until                    = ""
  indicator_types                = []
  labels                         = []
  object_marking_refs            = []
  threat_intelligence_tags       = []
  threat_types                   = []
  external_references {
    description = ""
    external_id = ""
    source_name = ""
    url         = ""
    hashes = {
      key = ""
    }
  }
  granular_markings {
    language    = ""
    marking_ref = 0
    selectors   = []
  }
  kill_chain_phases {
    kill_chain_name = ""
    phase_name      = ""
  }
  parsed_pattern {
    pattern_type_key = ""
    pattern_type_values {
      value      = ""
      value_type = ""
    }
  }
  extensions = jsonencode({
    "key" : "value"
  })

  depends_on = [azurerm_sentinel_log_analytics_workspace_onboarding.test]

}
`, r.template(data), data.RandomInteger)
}
