// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package newrelic_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/monitors"
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
	const AccountIdEnv = "ARM_ACCTEST_NEW_RELIC_ACCOUNT_ID"
	const OrgIdEnv = "ARM_ACCTEST_NEW_RELIC_ORG_ID"

	accountId := os.Getenv(AccountIdEnv)
	orgId := os.Getenv(OrgIdEnv)

	if accountId == "" || orgId == "" {
		t.Skipf("Acceptance test skipped unless env '%s' and '%s' set", AccountIdEnv, OrgIdEnv)
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_new_relic_monitor", "test")
	r := NewRelicMonitorResource{}
	effectiveDate := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	email := "b9ba4f77-5e63-4f1e-9445-b982d35f635b@example.com"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, effectiveDate, email, accountId, orgId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("ingestion_key"),
	})
}

func TestAccNewRelicMonitor_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_new_relic_monitor", "test")
	r := NewRelicMonitorResource{}
	effectiveDate := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	email := "fdfc9282-8817-442f-9f32-605ab174b610@example.com"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identity(data, effectiveDate, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
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
  identity {
    type = "SystemAssigned"
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
  identity {
    type = "SystemAssigned"
  }
}
`, config)
}

func (r NewRelicMonitorResource) complete(data acceptance.TestData, effectiveDate string, email string, accountId string, orgId string) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%[1]s

resource "azurerm_new_relic_monitor" "test" {
  name                = "acctest-nrm-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  plan {
    billing_cycle  = "MONTHLY"
    effective_date = "%[3]s"
    plan_id        = "newrelic-pay-as-you-go-free-live"
    usage_type     = "PAYG"
  }
  user {
    email        = "%[4]s"
    first_name   = "first"
    last_name    = "last"
    phone_number = "123456"
  }
  identity {
    type = "SystemAssigned"
  }
  account_creation_source = "LIFTR"
  account_id              = "%[5]s"
  ingestion_key           = "wltnimmhqt"
  organization_id         = "%[6]s"
  org_creation_source     = "LIFTR"
  user_id                 = "123456"
}
`, template, data.RandomInteger, effectiveDate, email, accountId, orgId)
}

func (r NewRelicMonitorResource) identity(data acceptance.TestData, effectiveDate string, email string) string {
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
  identity {
    type = "SystemAssigned"
  }
}

data "azurerm_subscription" "primary" {}

data "azurerm_role_definition" "test" {
  name = "Monitoring Reader"
}

resource "azurerm_role_assignment" "test" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = azurerm_new_relic_monitor.test.identity[0].principal_id
}
`, template, data.RandomInteger, data.Locations.Primary, effectiveDate, email)
}
