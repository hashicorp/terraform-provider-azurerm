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

type TagRulesDatadogMonitorResource struct{}

func TestAccDatadogMonitorTagRules_basic(t *testing.T) {
	if os.Getenv("ARM_TEST_DATADOG_API_KEY") == "" || os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY") == "" {
		t.Skip("Skipping as ARM_TEST_DATADOG_API_KEY and/or ARM_TEST_DATADOG_APPLICATION_KEY are not specified")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor_tag_rule", "test")
	r := TagRulesDatadogMonitorResource{}
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

func TestAccDatadogMonitorTagRules_update(t *testing.T) {
	if os.Getenv("ARM_TEST_DATADOG_API_KEY") == "" || os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY") == "" {
		t.Skip("Skipping as ARM_TEST_DATADOG_API_KEY and/or ARM_TEST_DATADOG_APPLICATION_KEY are not specified")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor_tag_rule", "test")
	r := TagRulesDatadogMonitorResource{}
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
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctest-datadog-%d"
  location = "%s"
}

resource "azurerm_datadog_monitor" "test" {
  name                = "acctest-datadog-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "WEST US 2"
  datadog_organization {
    api_key         = %q
    application_key = %q
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger%100, os.Getenv("ARM_TEST_DATADOG_API_KEY"), os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY"))
}

func (r TagRulesDatadogMonitorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`


	%s

resource "azurerm_datadog_monitor_tag_rule" "testbasic" {
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

func (r TagRulesDatadogMonitorResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

resource "azurerm_datadog_monitor_tag_rule" "testupdate" {
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
