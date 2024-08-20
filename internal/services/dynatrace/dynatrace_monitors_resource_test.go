package dynatrace_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type dynatraceInfo struct {
	UserCountry     string
	UserEmail       string
	UserFirstName   string
	UserLastName    string
	UserPhoneNumber string
}

type MonitorsResource struct {
	dynatraceInfo dynatraceInfo
}

func NewDynatraceMonitorResource() MonitorsResource {
	return MonitorsResource{
		dynatraceInfo: dynatraceInfo{
			UserCountry:     os.Getenv("DYNATRACE_USER_COUNTRY"),
			UserEmail:       os.Getenv("DYNATRACE_USER_EMAIL"),
			UserFirstName:   os.Getenv("DYNATRACE_USER_FIRST_NAME"),
			UserLastName:    os.Getenv("DYNATRACE_USER_LAST_NAME"),
			UserPhoneNumber: os.Getenv("DYNATRACE_USER_PHONE_NUMBER"),
		},
	}
}

func TestAccDynatraceMonitor_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dynatrace_monitor", "test")
	r := MonitorsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},

		data.ImportStep("user"),
	})
}

func TestAccDynatraceMonitor_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dynatrace_monitor", "test")
	r := MonitorsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user"),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user"),
	})
}

func TestAccDynatraceMonitor_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dynatrace_monitor", "test")
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

func (r MonitorsResource) preCheck(t *testing.T) {
	if r.dynatraceInfo.UserCountry == "" {
		t.Fatal("DYNATRACE_USER_COUNTRY must be set for acceptance tests")
	}
	if r.dynatraceInfo.UserEmail == "" {
		t.Fatal("DYNATRACE_USER_EMAIL must be set for acceptance tests")
	}
	if r.dynatraceInfo.UserFirstName == "" {
		t.Fatal("DYNATRACE_USER_FIRST_NAME must be set for acceptance tests")
	}
	if r.dynatraceInfo.UserLastName == "" {
		t.Fatal("DYNATRACE_USER_LAST_NAME must be set for acceptance tests")
	}
	if r.dynatraceInfo.UserPhoneNumber == "" {
		t.Fatal("DYNATRACE_USER_PHONE_NUMBER must be set for acceptance tests")
	}
}

func (r MonitorsResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := monitors.ParseMonitorID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Dynatrace.MonitorsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(true), nil
}

func (r MonitorsResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_dynatrace_monitor" "test" {
  name                     = "acctestacc%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  marketplace_subscription = "Active"

  identity {
    type = "SystemAssigned"
  }

  user {
    first_name   = "Alice"
    last_name    = "Bobab"
    email        = "alice@microsoft.com"
    phone_number = "12345"
    country      = "westus"
  }

  plan {
    usage_type     = "COMMITTED"
    billing_cycle  = "MONTHLY"
    plan           = "azureportalintegration_privatepreview@TIDgmz7xq9ge3py"
  }

  tags = {
    environment = "Dev"
  }
}
`, template, data.RandomInteger)
}

func (r MonitorsResource) updated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_dynatrace_monitor" "test" {
  name                     = "acctestacc%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  marketplace_subscription = "Active"

  identity {
    type = "SystemAssigned"
  }
  user {
    first_name   = "Alice"
    last_name    = "Bobab"
    email        = "alice@microsoft.com"
    phone_number = "12345"
    country      = "westus"
  }

  plan {
    usage_type     = "COMMITTED"
    billing_cycle  = "MONTHLY"
    plan           = "azureportalintegration_privatepreview@TIDgmz7xq9ge3py"
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

resource "azurerm_dynatrace_monitor" "import" {
  name                     = azurerm_dynatrace_monitor.test.name
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  identity = azurerm_dynatrace_monitor.test.identity
  monitoring_enabled       = azurerm_dynatrace_monitor.test.monitoring_enabled
  marketplace_subscription = azurerm_dynatrace_monitor.test.marketplace_subscription
  plan                     = azurerm_dynatrace_monitor.test.plan
  user                     = azurerm_dynatrace_monitor.test.user
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
`, data.RandomInteger, "eastus2euap", data.RandomString)
}
