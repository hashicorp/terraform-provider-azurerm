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

type LogzSubAccountResource struct{}

func TestAccLogzSubAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logz_sub_account", "test")
	r := LogzSubAccountResource{}
	effectiveDate := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	email := "0bc0fe71-6e2f-4552-bc48-6ca0c22f4db0@example.com"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, effectiveDate, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user"),
	})
}

func TestAccLogzSubAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logz_sub_account", "test")
	r := LogzSubAccountResource{}
	effectiveDate := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	email := "429420d4-0abf-4cd8-b149-fe63fa141fc6@example.com"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, effectiveDate, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, effectiveDate, email),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
	})
}

func TestAccLogzSubAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logz_sub_account", "test")
	r := LogzSubAccountResource{}
	effectiveDate := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	email := "253e466c-3a78-4a79-9260-98854eef2b5c@example.com"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, effectiveDate, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user"),
	})
}

func TestAccLogzSubAccount_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logz_sub_account", "test")
	r := LogzSubAccountResource{}
	effectiveDate := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	email := "5c3b9a35-06c5-4c75-928e-6505a10541a5@example.com"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, effectiveDate, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user"),
		{
			Config: r.update(data, effectiveDate, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user"),
		{
			Config: r.basic(data, effectiveDate, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user"),
	})
}

func (r LogzSubAccountResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.LogzSubAccountID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Logz.SubAccountClient.Get(ctx, id.ResourceGroup, id.MonitorName, id.AccountName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r LogzSubAccountResource) template(data acceptance.TestData, effectiveDate string, email string) string {
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
  plan {
    billing_cycle  = "MONTHLY"
    effective_date = "%s"
    plan_id        = "100gb14days"
    usage_type     = "COMMITTED"
  }

  user {
    email        = "%s"
    first_name   = "first"
    last_name    = "last"
    phone_number = "123456"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, effectiveDate, email)
}

func (r LogzSubAccountResource) basic(data acceptance.TestData, effectiveDate string, email string) string {
	template := r.template(data, effectiveDate, email)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_sub_account" "test" {
  name            = "acctest-lsa-%d"
  logz_monitor_id = azurerm_logz_monitor.test.id
  user {
    email        = azurerm_logz_monitor.test.user[0].email
    first_name   = azurerm_logz_monitor.test.user[0].first_name
    last_name    = azurerm_logz_monitor.test.user[0].last_name
    phone_number = azurerm_logz_monitor.test.user[0].phone_number
  }
}
`, template, data.RandomInteger)
}

func (r LogzSubAccountResource) update(data acceptance.TestData, effectiveDate string, email string) string {
	template := r.template(data, effectiveDate, email)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_sub_account" "test" {
  name            = "acctest-lsa-%d"
  logz_monitor_id = azurerm_logz_monitor.test.id
  user {
    email        = azurerm_logz_monitor.test.user[0].email
    first_name   = azurerm_logz_monitor.test.user[0].first_name
    last_name    = azurerm_logz_monitor.test.user[0].last_name
    phone_number = azurerm_logz_monitor.test.user[0].phone_number
  }
  enabled = false
}
`, template, data.RandomInteger)
}

func (r LogzSubAccountResource) requiresImport(data acceptance.TestData, effectiveDate string, email string) string {
	config := r.basic(data, effectiveDate, email)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_sub_account" "import" {
  name            = azurerm_logz_sub_account.test.name
  logz_monitor_id = azurerm_logz_sub_account.test.logz_monitor_id
  user {
    email        = azurerm_logz_monitor.test.user[0].email
    first_name   = azurerm_logz_monitor.test.user[0].first_name
    last_name    = azurerm_logz_monitor.test.user[0].last_name
    phone_number = azurerm_logz_monitor.test.user[0].phone_number
  }
}
`, config)
}

func (r LogzSubAccountResource) complete(data acceptance.TestData, effectiveDate string, email string) string {
	template := r.template(data, effectiveDate, email)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_sub_account" "test" {
  name            = "acctest-lsa-%d"
  logz_monitor_id = azurerm_logz_monitor.test.id
  user {
    email        = azurerm_logz_monitor.test.user[0].email
    first_name   = azurerm_logz_monitor.test.user[0].first_name
    last_name    = azurerm_logz_monitor.test.user[0].last_name
    phone_number = azurerm_logz_monitor.test.user[0].phone_number
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}
