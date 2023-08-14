// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logz_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LogzMonitorResource struct{}

func TestAccLogzMonitor_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logz_monitor", "test")
	r := LogzMonitorResource{}
	effectiveDate := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	email := "9d186100-1e0f-4b4a-bb10-753d2d52b750@example.com"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, effectiveDate, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogzMonitor_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logz_monitor", "test")
	r := LogzMonitorResource{}
	effectiveDate := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	email := "88841ffe-6376-487c-950c-1c3318f63dc5@example.com"
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

func TestAccLogzMonitor_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logz_monitor", "test")
	r := LogzMonitorResource{}
	effectiveDate := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	email := "37d395aa-4b30-4566-b141-72ea4bf84e11@example.com"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, effectiveDate, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogzMonitor_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logz_monitor", "test")
	r := LogzMonitorResource{}
	effectiveDate := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	email := "d5d750ce-94fd-475d-816f-48110e1ca04a@example.com"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, effectiveDate, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, effectiveDate, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, effectiveDate, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LogzMonitorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := monitors.ParseMonitorID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := clients.Logz.MonitorClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
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

func (r LogzMonitorResource) basic(data acceptance.TestData, effectiveDate string, email string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_monitor" "test" {
  name                = "acctest-lm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  plan {
    billing_cycle  = "MONTHLY"
    effective_date = "%s"
    usage_type     = "COMMITTED"
  }

  user {
    email        = "%s"
    first_name   = "first"
    last_name    = "last"
    phone_number = "123456"
  }
}
`, template, data.RandomInteger, effectiveDate, email)
}

func (r LogzMonitorResource) update(data acceptance.TestData, effectiveDate string, email string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_monitor" "test" {
  name                = "acctest-lm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  plan {
    billing_cycle  = "MONTHLY"
    effective_date = "%s"
    usage_type     = "COMMITTED"
  }

  user {
    email        = "%s"
    first_name   = "first"
    last_name    = "last"
    phone_number = "123456"
  }
  enabled = false
}
`, template, data.RandomInteger, effectiveDate, email)
}

func (r LogzMonitorResource) requiresImport(data acceptance.TestData, effectiveDate string, email string) string {
	config := r.basic(data, effectiveDate, email)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_monitor" "import" {
  name                = azurerm_logz_monitor.test.name
  resource_group_name = azurerm_logz_monitor.test.resource_group_name
  location            = azurerm_logz_monitor.test.location
  plan {
    billing_cycle  = "MONTHLY"
    effective_date = "%s"
    usage_type     = "COMMITTED"
  }

  user {
    email        = "%s"
    first_name   = "first"
    last_name    = "last"
    phone_number = "123456"
  }
}
`, config, effectiveDate, email)
}

func (r LogzMonitorResource) complete(data acceptance.TestData, effectiveDate string, email string) string {
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
    billing_cycle  = "MONTHLY"
    effective_date = "%s"
    usage_type     = "COMMITTED"
  }

  user {
    email        = "%s"
    first_name   = "first"
    last_name    = "last"
    phone_number = "123456"
  }
  enabled = true
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger, effectiveDate, email)
}
