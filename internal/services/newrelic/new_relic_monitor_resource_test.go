// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package newrelic_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2022-07-01/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NewRelicMonitorResource struct{}

func TestAccNewRelicMonitor_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_new_relic_monitor", "test")
	r := NewRelicMonitorResource{}
	effectiveDate := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	email := "d327e362-8431-4df1-8d99-8dc1c383a4f3@example.com"
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

func TestAccNewRelicMonitor_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_new_relic_monitor", "test")
	r := NewRelicMonitorResource{}
	effectiveDate := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	email := "15f0c06e-0cda-4a46-8baa-f6ec19f0ff94@example.com"
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

func TestAccNewRelicMonitor_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_new_relic_monitor", "test")
	r := NewRelicMonitorResource{}
	effectiveDate := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	email := "b9ba4f77-5e63-4f1e-9445-b982d35f635b@example.com"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, effectiveDate, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("ingestion_key"),
	})
}

func (r NewRelicMonitorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := monitors.ParseMonitorID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.NewRelic.MonitorsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r NewRelicMonitorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r NewRelicMonitorResource) basic(data acceptance.TestData, effectiveDate string, email string) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_new_relic_monitor" "test" {
  name                = "acctest-nrm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  plan {
    effective_date = "%s"
  }
  user {
    email        = "%s"
    first_name   = "first"
    last_name    = "last"
    phone_number = "123456"
  }
}
`, template, data.RandomInteger, data.Locations.Primary, effectiveDate, email)
}

func (r NewRelicMonitorResource) requiresImport(data acceptance.TestData, effectiveDate string, email string) string {
	config := r.basic(data, effectiveDate, email)
	return fmt.Sprintf(`
			%s

resource "azurerm_new_relic_monitor" "import" {
  name                = azurerm_new_relic_monitor.test.name
  resource_group_name = azurerm_new_relic_monitor.test.resource_group_name
  location            = azurerm_new_relic_monitor.test.location
  plan {
    effective_date = azurerm_new_relic_monitor.test.plan[0].effective_date
  }
  user {
    email        = azurerm_new_relic_monitor.test.user[0].email
    first_name   = azurerm_new_relic_monitor.test.user[0].first_name
    last_name    = azurerm_new_relic_monitor.test.user[0].last_name
    phone_number = azurerm_new_relic_monitor.test.user[0].phone_number
  }
}
`, config)
}

func (r NewRelicMonitorResource) complete(data acceptance.TestData, effectiveDate string, email string) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%[1]s

resource "azurerm_new_relic_monitor" "org" {
  name                = "acctest-nrmo-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  plan {
    effective_date = "%[4]s"
  }
  user {
    email        = "%[5]s"
    first_name   = "first"
    last_name    = "last"
    phone_number = "123456"
  }
}

resource "azurerm_new_relic_monitor" "test" {
  name                = "acctest-nrm-%[2]d"
  resource_group_name = azurerm_new_relic_monitor.org.resource_group_name
  location            = azurerm_new_relic_monitor.org.location
  plan {
    billing_cycle  = azurerm_new_relic_monitor.org.plan[0].billing_cycle
    effective_date = "%[4]s"
    plan_id        = azurerm_new_relic_monitor.org.plan[0].plan_id
    usage_type     = azurerm_new_relic_monitor.org.plan[0].usage_type
  }
  user {
    email        = azurerm_new_relic_monitor.org.user[0].email
    first_name   = azurerm_new_relic_monitor.org.user[0].first_name
    last_name    = azurerm_new_relic_monitor.org.user[0].last_name
    phone_number = azurerm_new_relic_monitor.org.user[0].phone_number
  }
  account_creation_source = azurerm_new_relic_monitor.org.account_creation_source
  account_id              = azurerm_new_relic_monitor.org.account_id
  ingestion_key           = "wltnimmhqt"
  organization_id         = azurerm_new_relic_monitor.org.organization_id
  org_creation_source     = azurerm_new_relic_monitor.org.org_creation_source
  user_id                 = "123456"
}
`, template, data.RandomInteger, data.Locations.Primary, effectiveDate, email)
}
