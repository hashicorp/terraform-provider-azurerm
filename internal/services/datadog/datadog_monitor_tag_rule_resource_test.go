// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datadog_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type TagRulesDatadogMonitorResource struct {
	datadogApiKey         string
	datadogApplicationKey string
}

func (r *TagRulesDatadogMonitorResource) populateFromEnvironment(t *testing.T) {
	if os.Getenv("ARM_TEST_DATADOG_API_KEY") == "" {
		t.Skip("Skipping as ARM_TEST_DATADOG_API_KEY is not specified")
	}
	if os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY") == "" {
		t.Skip("Skipping as ARM_TEST_DATADOG_APPLICATION_KEY is not specified")
	}
	r.datadogApiKey = os.Getenv("ARM_TEST_DATADOG_API_KEY")
	r.datadogApplicationKey = os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY")
}

func TestAccDatadogMonitorTagRules_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor_tag_rule", "test")
	r := TagRulesDatadogMonitorResource{}
	r.populateFromEnvironment(t)
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

func TestAccDatadogMonitorTagRules_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor_tag_rule", "test")
	r := TagRulesDatadogMonitorResource{}
	r.populateFromEnvironment(t)
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

func TestAccDatadogMonitorTagRules_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor_tag_rule", "test")
	r := TagRulesDatadogMonitorResource{}
	r.populateFromEnvironment(t)
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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r TagRulesDatadogMonitorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := rules.ParseTagRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Datadog.Rules.TagRulesGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r TagRulesDatadogMonitorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-datadogrg-%[1]d"
  location = %[2]q
}

resource "azurerm_datadog_monitor" "test" {
  name                = "acctest-datadog-%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datadog_organization {
    api_key         = %[4]q
    application_key = %[5]q
  }
  user {
    name  = "Test Datadog"
    email = "abc@xyz.com"
  }
  sku_name = "Linked"
  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, r.datadogApiKey, r.datadogApplicationKey)
}

func (r TagRulesDatadogMonitorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_datadog_monitor_tag_rule" "test" {
  datadog_monitor_id = azurerm_datadog_monitor.test.id
  log {
    subscription_log_enabled = true
  }
  metric {
    filter {
      name   = "Test"
      value  = "Testing-Logs"
      action = "Include"
    }
  }
}
`, r.template(data))
}

func (r TagRulesDatadogMonitorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_datadog_monitor_tag_rule" "import" {
  datadog_monitor_id = azurerm_datadog_monitor_tag_rule.test.datadog_monitor_id
  name               = azurerm_datadog_monitor_tag_rule.test.name
  log {
    subscription_log_enabled = true
  }
  metric {
    filter {
      name   = "Test"
      value  = "Testing-Logs"
      action = "Include"
    }
  }
}
`, r.basic(data))
}

func (r TagRulesDatadogMonitorResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_datadog_monitor_tag_rule" "test" {
  datadog_monitor_id = azurerm_datadog_monitor.test.id
  log {
    subscription_log_enabled = false
    resource_log_enabled     = true
    filter {
      name   = "Test"
      value  = "Testing-Logs"
      action = "Include"
    }
  }
  metric {
    filter {
      name   = "Test"
      value  = "Testing-Logs"
      action = "Include"
    }
  }
}
`, r.template(data))
}
