package logz_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logz/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogzMonitorResource struct{}

func TestAccLogzMonitor_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logz_monitor", "test")
	r := LogzMonitorResource{}
	effectiveDate := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, effectiveDate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user"),
	})
}

func TestAccLogzMonitor_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logz_monitor", "test")
	r := LogzMonitorResource{}
	effectiveDate := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, effectiveDate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccLogzMonitor_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logz_monitor", "test")
	r := LogzMonitorResource{}
	effectiveDate := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, effectiveDate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user"),
	})
}

func TestAccLogzMonitor_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logz_monitor", "test")
	r := LogzMonitorResource{}
	effectiveDate := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, effectiveDate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user"),
		{
			Config: r.update(data, effectiveDate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user"),
		{
			Config: r.basic(data, effectiveDate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user"),
	})
}

func (r LogzMonitorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.LogzMonitorID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := clients.Logz.MonitorClient.Get(ctx, id.ResourceGroup, id.MonitorName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r LogzMonitorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-logz-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r LogzMonitorResource) basic(data acceptance.TestData, effectiveDate string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_monitor" "test" {
  name                = "acctest-lm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  plan {
    billing_cycle  = "Monthly"
    effective_date = "%s"
    plan_id        = "100gb14days"
    usage_type     = "Committed"
  }

  user {
    email        = "e081a27c-bc01-4159-bc06-7f9f711e3b3a@example.com"
    first_name   = "first"
    last_name    = "last"
    phone_number = "123456"
  }
}
`, template, data.RandomInteger, effectiveDate)
}

func (r LogzMonitorResource) update(data acceptance.TestData, effectiveDate string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_monitor" "test" {
  name                = "acctest-lm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  plan {
    billing_cycle  = "Monthly"
    effective_date = "%s"
    plan_id        = "100gb14days"
    usage_type     = "Committed"
  }

  user {
    email        = "e081a27c-bc01-4159-bc06-7f9f711e3b3a@example.com"
    first_name   = "first"
    last_name    = "last"
    phone_number = "123456"
  }
  enabled = false
}
`, template, data.RandomInteger, effectiveDate)
}

func (r LogzMonitorResource) requiresImport(data acceptance.TestData) string {
	effectiveDate := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	config := r.basic(data, effectiveDate)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_monitor" "import" {
  name                = azurerm_logz_monitor.test.name
  resource_group_name = azurerm_logz_monitor.test.resource_group_name
  location            = azurerm_logz_monitor.test.location
  plan {
    billing_cycle  = "Monthly"
    effective_date = "%s"
    plan_id        = "100gb14days"
    usage_type     = "Committed"
  }

  user {
    email        = "e081a27c-bc01-4159-bc06-7f9f711e3b3a@example.com"
    first_name   = "first"
    last_name    = "last"
    phone_number = "123456"
  }
}
`, config, effectiveDate)
}

func (r LogzMonitorResource) complete(data acceptance.TestData, effectiveDate string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_monitor" "test" {
  name                = "acctest-lm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  company_name      = "company-name-1"
  enterprise_app_id = "e081a27c-bc01-4159-bc06-7f9f711e3b3a"
  plan {
    billing_cycle  = "Monthly"
    effective_date = "%s"
    plan_id        = "100gb14days"
    usage_type     = "Committed"
  }

  user {
    email        = "e081a27c-bc01-4159-bc06-7f9f711e3b3a@example.com"
    first_name   = "first"
    last_name    = "last"
    phone_number = "123456"
  }
  enabled = false
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger, effectiveDate)
}
