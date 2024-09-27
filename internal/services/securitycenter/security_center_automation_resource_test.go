// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2019-01-01-preview/automations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SecurityCenterAutomationResource struct{}

func TestAccSecurityCenterAutomation_logicApp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_automation", "test")
	r := SecurityCenterAutomationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logicApp(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("action.0.trigger_url"), // trigger_url needs to be ignored
	})
}

func TestAccSecurityCenterAutomation_logAnalytics(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_automation", "test")
	r := SecurityCenterAutomationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logAnalytics(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityCenterAutomation_eventHub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_automation", "test")
	r := SecurityCenterAutomationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.eventHub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("action.0.connection_string"), // trigger_url needs to be ignored
	})
}

func TestAccSecurityCenterAutomation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_automation", "test")
	r := SecurityCenterAutomationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logicApp(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSecurityCenterAutomation_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_automation", "test")
	r := SecurityCenterAutomationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logicApp(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("action.0.trigger_url"), // trigger_url needs to be ignored
		{
			Config: r.ruleSingle(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("action.0.trigger_url"), // trigger_url needs to be ignored
		{
			Config: r.logicApp(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("action.0.trigger_url"), // trigger_url needs to be ignored
	})
}

func TestAccSecurityCenterAutomation_ruleSingle(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_automation", "test")
	r := SecurityCenterAutomationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ruleSingle(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("source.#").HasValue("1"),
				check.That(data.ResourceName).Key("source.0.rule_set.#").HasValue("1"),
			),
		},
		data.ImportStep("action.0.trigger_url"), // trigger_url needs to be ignored
	})
}

func TestAccSecurityCenterAutomation_ruleMulti(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_automation", "test")
	r := SecurityCenterAutomationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ruleMulti(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("source.#").HasValue("1"),
				check.That(data.ResourceName).Key("source.0.rule_set.#").HasValue("1"),
				check.That(data.ResourceName).Key("source.0.rule_set.0.rule.#").HasValue("3"),
			),
		},
		data.ImportStep("action.0.trigger_url"), // trigger_url needs to be ignored
	})
}

func TestAccSecurityCenterAutomation_ruleSetMulti(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_automation", "test")
	r := SecurityCenterAutomationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ruleSetMulti(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("source.#").HasValue("1"),
				check.That(data.ResourceName).Key("source.0.rule_set.#").HasValue("2"),
				check.That(data.ResourceName).Key("source.0.rule_set.0.rule.#").HasValue("2"),
				check.That(data.ResourceName).Key("source.0.rule_set.1.rule.#").HasValue("2"),
			),
		},
		data.ImportStep("action.0.trigger_url"), // trigger_url needs to be ignored
	})
}

func TestAccSecurityCenterAutomation_scopeMulti(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_automation", "test")
	r := SecurityCenterAutomationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scopeMulti(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("scopes.#").HasValue("3"),
			),
		},
		data.ImportStep("action.0.trigger_url"), // trigger_url needs to be ignored
	})
}

func TestAccSecurityCenterAutomation_actionMulti(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_automation", "test")
	r := SecurityCenterAutomationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.actionMulti(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("action.#").HasValue("2"),
			),
		},
		data.ImportStep("action.0.trigger_url", "action.1.trigger_url"), // trigger_url needs to be ignored
	})
}

func TestAccSecurityCenterAutomation_sourceMulti(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_automation", "test")
	r := SecurityCenterAutomationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sourceMulti(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("source.#").HasValue("9"),
				check.That(data.ResourceName).Key("source.0.rule_set.#").HasValue("1"),
				check.That(data.ResourceName).Key("source.1.rule_set.#").HasValue("1"),
				check.That(data.ResourceName).Key("source.2.rule_set.#").HasValue("1"),
				check.That(data.ResourceName).Key("source.3.rule_set.#").HasValue("0"),
				check.That(data.ResourceName).Key("source.4.rule_set.#").HasValue("0"),
				check.That(data.ResourceName).Key("source.5.rule_set.#").HasValue("0"),
				check.That(data.ResourceName).Key("source.6.rule_set.#").HasValue("0"),
				check.That(data.ResourceName).Key("source.7.rule_set.#").HasValue("0"),
				check.That(data.ResourceName).Key("source.8.rule_set.#").HasValue("0"),
			),
		},
		data.ImportStep("action.0.trigger_url", "action.1.trigger_url"), // trigger_url needs to be ignored
	})
}

func (t SecurityCenterAutomationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AutomationID(state.ID)
	if err != nil {
		return nil, err
	}
	automationId := automations.NewAutomationID(id.SubscriptionId, id.ResourceGroup, id.Name)

	resp, err := clients.SecurityCenter.AutomationsClient.Get(ctx, automationId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model.Properties != nil), nil
}

func (SecurityCenterAutomationResource) logicApp(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlogicapp-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_client_config" "current" {
}

resource "azurerm_security_center_automation" "test" {
  name                = "acctestautomation-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name

  scopes = [
    "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  ]

  action {
    type        = "logicapp"
    resource_id = azurerm_logic_app_workflow.test.id
    trigger_url = "https://example.net/this_is_never_validated_by_azure"
  }

  source {
    event_source = "Alerts"
  }

  tags = {
    Env2 = "Test2"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (SecurityCenterAutomationResource) logAnalytics(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestlogs-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

data "azurerm_client_config" "current" {
}

resource "azurerm_security_center_automation" "test" {
  name                = "acctestautomation-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name

  scopes = [
    "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  ]

  action {
    type        = "loganalytics"
    resource_id = azurerm_log_analytics_workspace.test.id
  }

  source {
    event_source = "Alerts"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (SecurityCenterAutomationResource) eventHub(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhub-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
  capacity            = 1
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  partition_count   = 2
  message_retention = 1
}

resource "azurerm_eventhub_authorization_rule" "test" {
  name                = "acctest-eventhub-auth-rule-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name

  listen = true
  send   = false
  manage = false
}

data "azurerm_client_config" "current" {
}

resource "azurerm_security_center_automation" "test" {
  name                = "acctestautomation-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name

  scopes = [
    "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  ]

  action {
    type              = "eventhub"
    resource_id       = azurerm_eventhub.test.id
    connection_string = azurerm_eventhub_authorization_rule.test.primary_connection_string
  }

  source {
    event_source = "Alerts"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.Locations.Primary)
}

func (r SecurityCenterAutomationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_security_center_automation" "import" {
  name                = azurerm_security_center_automation.test.name
  location            = azurerm_security_center_automation.test.location
  resource_group_name = azurerm_security_center_automation.test.resource_group_name

  scopes = [
    "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  ]

  action {
    type        = "logicapp"
    resource_id = azurerm_logic_app_workflow.test.id
    trigger_url = "https://example.net/this_is_never_validated_by_azure"
  }

  source {
    event_source = "Alerts"
  }
}
`, r.logicApp(data))
}

func (SecurityCenterAutomationResource) ruleSingle(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlogicapp-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_client_config" "current" {
}

resource "azurerm_security_center_automation" "test" {
  name                = "acctestautomation-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name

  scopes = [
    "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  ]

  action {
    type        = "logicapp"
    resource_id = azurerm_logic_app_workflow.test.id
    trigger_url = "https://example.net/this_is_never_validated_by_azure"
  }

  source {
    event_source = "Alerts"
    rule_set {
      rule {
        property_path  = "properties.metadata.severity"
        operator       = "Equals"
        expected_value = "High"
        property_type  = "String"
      }
    }
  }

  description = "Security Center Automation Acc test"
  tags = {
    Env = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (SecurityCenterAutomationResource) scopeMulti(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlogicapp-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_client_config" "current" {
}

resource "azurerm_security_center_automation" "test" {
  name                = "acctestautomation-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name

  scopes = [
    "/subscriptions/${data.azurerm_client_config.current.subscription_id}",
    "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/test",
    "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/test2"
  ]

  action {
    type        = "logicapp"
    resource_id = azurerm_logic_app_workflow.test.id
    trigger_url = "https://example.net/this_is_never_validated_by_azure"
  }

  source {
    event_source = "Alerts"
    rule_set {
      rule {
        property_path  = "properties.metadata.severity"
        operator       = "Equals"
        expected_value = "High"
        property_type  = "String"
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (SecurityCenterAutomationResource) ruleMulti(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlogicapp-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_client_config" "current" {
}

resource "azurerm_security_center_automation" "test" {
  name                = "acctestautomation-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name

  scopes = [
    "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  ]

  action {
    type        = "logicapp"
    resource_id = azurerm_logic_app_workflow.test.id
    trigger_url = "https://example.net/this_is_never_validated_by_azure"
  }

  source {
    event_source = "Alerts"
    rule_set {
      rule {
        property_path  = "properties.metadata.severity"
        operator       = "Equals"
        expected_value = "Low"
        property_type  = "String"
      }
      rule {
        property_path  = "properties.metadata.severity"
        operator       = "Equals"
        expected_value = "Medium"
        property_type  = "String"
      }
      rule {
        property_path  = "properties.metadata.severity"
        operator       = "Equals"
        expected_value = "High"
        property_type  = "String"
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (SecurityCenterAutomationResource) ruleSetMulti(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlogicapp-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_client_config" "current" {
}

resource "azurerm_security_center_automation" "test" {
  name                = "acctestautomation-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name

  scopes = [
    "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  ]

  action {
    type        = "logicapp"
    resource_id = azurerm_logic_app_workflow.test.id
    trigger_url = "https://example.net/this_is_never_validated_by_azure"
  }

  source {
    event_source = "Alerts"
    rule_set {
      rule {
        property_path  = "properties.metadata.title"
        operator       = "Equals"
        expected_value = "Tony Iommi"
        property_type  = "String"
      }
      rule {
        property_path  = "properties.metadata.title"
        operator       = "Equals"
        expected_value = "Ozzy Osbourne"
        property_type  = "String"
      }
    }
    rule_set {
      rule {
        property_path  = "properties.metadata.title"
        operator       = "Equals"
        expected_value = "Bill Ward"
        property_type  = "String"
      }
      rule {
        property_path  = "properties.metadata.title"
        operator       = "Equals"
        expected_value = "Geezer Butler"
        property_type  = "String"
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (SecurityCenterAutomationResource) actionMulti(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlogicapp-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_logic_app_workflow" "test2" {
  name                = "acctestlogicapp2-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_client_config" "current" {
}

resource "azurerm_security_center_automation" "test" {
  name                = "acctestautomation-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name

  scopes = [
    "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  ]

  action {
    type        = "logicapp"
    resource_id = azurerm_logic_app_workflow.test.id
    trigger_url = "https://example.net/this_is_never_validated_by_azure"
  }

  action {
    type        = "logicapp"
    resource_id = azurerm_logic_app_workflow.test2.id
    trigger_url = "https://example.net/this_is_also_never_validated_by_azure"
  }

  source {
    event_source = "Alerts"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (SecurityCenterAutomationResource) sourceMulti(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlogicapp-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_client_config" "current" {
}

resource "azurerm_security_center_automation" "test" {
  name                = "acctestautomation-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name

  scopes = [
    "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  ]

  action {
    type        = "logicapp"
    resource_id = azurerm_logic_app_workflow.test.id
    trigger_url = "https://example.net/this_is_never_validated_by_azure"
  }

  source {
    event_source = "Alerts"
    rule_set {
      rule {
        property_path  = "properties.metadata.severity"
        operator       = "Equals"
        expected_value = "Low"
        property_type  = "String"
      }
    }
  }

  source {
    event_source = "Assessments"
    rule_set {
      rule {
        property_path  = "properties.metadata.severity"
        operator       = "Equals"
        expected_value = "Low"
        property_type  = "String"
      }
    }
  }

  source {
    event_source = "SubAssessments"
    rule_set {
      rule {
        property_path  = "properties.metadata.severity"
        operator       = "Equals"
        expected_value = "Low"
        property_type  = "String"
      }
    }
  }

  source {
    event_source = "SecureScores"
  }

  source {
    event_source = "SecureScoreControls"
  }

  source {
    event_source = "RegulatoryComplianceAssessment"
  }

  source {
    event_source = "RegulatoryComplianceAssessmentSnapshot"
  }

  source {
    event_source = "SecureScoreControlsSnapshot"
  }

  source {
    event_source = "SecureScoresSnapshot"
  }

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}
