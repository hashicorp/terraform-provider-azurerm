package dynatrace_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2021-09-01/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MonitorsResource struct{}

func TestAccDynatraceMonitors_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dynatrace_monitors", "test")
	r := MonitorsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},

		data.ImportStep("user_info"),
	})
}

func TestAccDynatraceMonitors_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dynatrace_monitors", "test")
	r := MonitorsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user_info"),
	})
}

func TestAccDynatraceMonitors_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dynatrace_monitors", "test")
	r := MonitorsResource{}

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

func (r MonitorsResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := monitors.ParseMonitorID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Dynatrace.MonitorClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r MonitorsResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_dynatrace_monitors" "test" {
  name                            = "acctestacc%[2]d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  identity_type                   = "SystemAssigned"
  monitoring_status               = "Enabled"
  marketplace_subscription_status = "Active"

  user_info {
    first_name    = "Alice"
    last_name     = "Bobab"
    email_address = "alice@microsoft.com"
    phone_number  = "123456"
    country       = "westus"
  }

  plan_data {
    usage_type     = "COMMITTED"
    billing_cycle  = "MONTHLY"
    plan_details   = "azureportalintegration_privatepreview@TIDhjdtn7tfnxcy"
    effective_date = "2019-08-30T15:14:33Z"
  }
}
`, template, data.RandomInteger)
}

func (r MonitorsResource) updated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_dynatrace_monitors" "test" {
  name                            = "acctestacc%[2]d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  identity_type                   = "SystemAssigned"
  monitoring_status               = "Enabled"
  marketplace_subscription_status = "Active"

  user_info {
    first_name    = "Alice"
    last_name     = "Bobab"
    email_address = "alice@microsoft.com"
    phone_number  = "123456"
    country       = "westus"
  }

  plan_data {
    usage_type     = "COMMITTED"
    billing_cycle  = "MONTHLY"
    plan_details   = "azureportalintegration_privatepreview@TIDhjdtn7tfnxcy"
    effective_date = "2019-08-30T15:14:33Z"
  }

  tags = {
    environment = "Prod"
    test        = "Patch"
  }
}
`, template, data.RandomInteger)
}

func (r MonitorsResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dynatrace_monitors" "import" {
  name                            = azurerm_dynatrace_monitors.test.name
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  identity_type                   = azurerm_dynatrace_monitors.test.identity_type
  monitoring_status               = azurerm_dynatrace_monitors.test.monitoring_status
  marketplace_subscription_status = azurerm_dynatrace_monitors.test.marketplace_subscription_status
}
`, template)
}

func (r MonitorsResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
