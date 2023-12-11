// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SentinelAlertRuleFusionResource struct{}

func TestAccSentinelAlertRuleFusion_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_fusion", "test")
	r := SentinelAlertRuleFusionResource{}

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

func TestAccSentinelAlertRuleFusion_disable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_fusion", "test")
	r := SentinelAlertRuleFusionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.disabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSentinelAlertRuleFusion_sourceSetting(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_fusion", "test")
	r := SentinelAlertRuleFusionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sourceSetting(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.sourceSetting(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSentinelAlertRuleFusion_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_fusion", "test")
	r := SentinelAlertRuleFusionResource{}

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

func (r SentinelAlertRuleFusionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	alertRuleClient := client.Sentinel.AlertRulesClient
	id, err := alertrules.ParseAlertRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := alertRuleClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving Sentinel Alert Rule Fusion (%q): %+v", state.String(), err)
	}

	if model := resp.Model; model != nil {
		modelPtr := *model
		rule, ok := modelPtr.(alertrules.FusionAlertRule)
		if !ok {
			return nil, fmt.Errorf("the Alert Rule %q is not a Fusion Alert Rule", id)
		}
		return utils.Bool(rule.Id != nil), nil
	}

	return utils.Bool(false), nil
}

func (r SentinelAlertRuleFusionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_sentinel_alert_rule_template" "test" {
  display_name               = "Advanced Multistage Attack Detection"
  log_analytics_workspace_id = azurerm_log_analytics_solution.test.workspace_resource_id
}

resource "azurerm_sentinel_alert_rule_fusion" "test" {
  name                       = "acctest-SentinelAlertRule-Fusion-%d"
  log_analytics_workspace_id = azurerm_log_analytics_solution.test.workspace_resource_id
  alert_rule_template_guid   = data.azurerm_sentinel_alert_rule_template.test.name
}
`, r.template(data), data.RandomInteger)
}

func (r SentinelAlertRuleFusionResource) disabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_sentinel_alert_rule_template" "test" {
  display_name               = "Advanced Multistage Attack Detection"
  log_analytics_workspace_id = azurerm_log_analytics_solution.test.workspace_resource_id
}

resource "azurerm_sentinel_alert_rule_fusion" "test" {
  name                       = "acctest-SentinelAlertRule-Fusion-%d"
  log_analytics_workspace_id = azurerm_log_analytics_solution.test.workspace_resource_id
  alert_rule_template_guid   = data.azurerm_sentinel_alert_rule_template.test.name
  enabled                    = false
}
`, r.template(data), data.RandomInteger)
}

func (r SentinelAlertRuleFusionResource) sourceSetting(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_sentinel_alert_rule_template" "test" {
  display_name               = "Advanced Multistage Attack Detection"
  log_analytics_workspace_id = azurerm_log_analytics_solution.test.workspace_resource_id
}

resource "azurerm_sentinel_alert_rule_fusion" "test" {
  name                       = "acctest-SentinelAlertRule-Fusion-%[2]d"
  log_analytics_workspace_id = azurerm_log_analytics_solution.test.workspace_resource_id
  alert_rule_template_guid   = data.azurerm_sentinel_alert_rule_template.test.name
  source {
    name    = "Anomalies"
    enabled = %[3]t
  }
  source {
    name    = "Alert providers"
    enabled = %[3]t
    sub_type {
      severities_allowed = ["High", "Informational", "Low", "Medium"]
      name               = "Azure Active Directory Identity Protection"
      enabled            = %[3]t
    }
    sub_type {
      severities_allowed = ["High", "Informational", "Low", "Medium"]
      name               = "Microsoft 365 Defender"
      enabled            = %[3]t
    }
    sub_type {
      severities_allowed = ["High", "Informational", "Low", "Medium"]
      name               = "Microsoft Cloud App Security"
      enabled            = %[3]t
    }
    sub_type {
      severities_allowed = ["High", "Informational", "Low", "Medium"]
      name               = "Azure Defender"
      enabled            = %[3]t
    }
    sub_type {
      severities_allowed = ["High", "Informational", "Low", "Medium"]
      name               = "Microsoft Defender for Endpoint"
      enabled            = %[3]t
    }
    sub_type {
      severities_allowed = ["High", "Informational", "Low", "Medium"]
      name               = "Microsoft Defender for Identity"
      enabled            = %[3]t
    }
    sub_type {
      severities_allowed = ["High", "Informational", "Low", "Medium"]
      name               = "Azure Defender for IoT"
      enabled            = %[3]t
    }
    sub_type {
      severities_allowed = ["High", "Informational", "Low", "Medium"]
      name               = "Microsoft Defender for Office 365"
      enabled            = %[3]t
    }
    sub_type {
      severities_allowed = ["High", "Informational", "Low", "Medium"]
      name               = "Azure Sentinel scheduled analytics rules"
      enabled            = %[3]t
    }
    sub_type {
      severities_allowed = ["High", "Informational", "Low", "Medium"]
      name               = "Azure Sentinel NRT analytic rules"
      enabled            = %[3]t
    }
  }
}
`, r.template(data), data.RandomInteger, enabled)
}

func (r SentinelAlertRuleFusionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_fusion" "import" {
  name                       = azurerm_sentinel_alert_rule_fusion.test.name
  log_analytics_workspace_id = azurerm_sentinel_alert_rule_fusion.test.log_analytics_workspace_id
  alert_rule_template_guid   = azurerm_sentinel_alert_rule_fusion.test.alert_rule_template_guid
}
`, r.basic(data))
}

func (r SentinelAlertRuleFusionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sentinel-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_log_analytics_solution" "test" {
  solution_name         = "SecurityInsights"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  workspace_resource_id = azurerm_log_analytics_workspace.test.id
  workspace_name        = azurerm_log_analytics_workspace.test.name

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/SecurityInsights"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
