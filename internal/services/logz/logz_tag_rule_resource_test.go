package logz_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logz/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogzTagRuleResource struct{}

func TestAccLogzTagRule_basic(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_LOGZ_MONITOR_EMAIL") == "" {
		t.Skip("Skipping as `ARM_RUN_TEST_LOGZ_MONITOR_EMAIL` was not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_logz_tag_rule", "test")
	r := LogzTagRuleResource{}
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

func TestAccLogzTagRule_requiresImport(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_LOGZ_MONITOR_EMAIL") == "" {
		t.Skip("Skipping as `ARM_RUN_TEST_LOGZ_MONITOR_EMAIL` was not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_logz_tag_rule", "test")
	r := LogzTagRuleResource{}
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

func TestAccLogzTagRule_complete(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_LOGZ_MONITOR_EMAIL") == "" {
		t.Skip("Skipping as `ARM_RUN_TEST_LOGZ_MONITOR_EMAIL` was not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_logz_tag_rule", "test")
	r := LogzTagRuleResource{}
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

func TestAccLogzTagRule_update(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_LOGZ_MONITOR_EMAIL") == "" {
		t.Skip("Skipping as `ARM_RUN_TEST_LOGZ_MONITOR_EMAIL` was not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_logz_tag_rule", "test")
	r := LogzTagRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LogzTagRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.LogzTagRuleID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := clients.Logz.TagRuleClient.Get(ctx, id.ResourceGroup, id.MonitorName, id.TagRuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r LogzTagRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-logz-%d"
  location = "%s"
}

resource "azurerm_logz_monitor" "test" {
  name                = "acctest-lm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  plan_data {
    billing_cycle  = "Monthly"
    effective_date = "%s"
    plan_id        = "100gb14days"
    usage_type     = "Committed"
  }

  user_info {
    email_address = "%s"
    first_name    = "first"
    last_name     = "last"
    phone_number  = "123456"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, getEffectiveDate(), os.Getenv("ARM_RUN_TEST_LOGZ_MONITOR_EMAIL"))
}

func (r LogzTagRuleResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_tag_rule" "test" {
  resource_group_name = azurerm_resource_group.test.name
  logz_monitor_id     = azurerm_logz_monitor.test.name
}
`, template)
}

func (r LogzTagRuleResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_tag_rule" "import" {
  resource_group_name = azurerm_logz_tag_rule.test.resource_group_name
  logz_monitor_id     = azurerm_logz_tag_rule.test.logz_monitor_id
}
`, config)
}

func (r LogzTagRuleResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_tag_rule" "test" {
  resource_group_name = azurerm_resource_group.test.name
  logz_monitor_id     = azurerm_logz_monitor.test.name
  filtering_tag {
    name   = "ccc"
    action = "Include"
    value  = "ccc"
  }

  filtering_tag {
    name   = "bbb"
    action = "Exclude"
    value  = ""
  }

  filtering_tag {
    name   = "ccc"
    action = "Include"
    value  = "ccc"
  }
  send_aad_logs          = true
  send_activity_logs     = true
  send_subscription_logs = true
}
`, template)
}

func getEffectiveDate() string {
	effectiveDate := time.Now().Add(time.Hour * 72)
	year, month, day := effectiveDate.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, effectiveDate.Location()).Format(time.RFC3339)
}
