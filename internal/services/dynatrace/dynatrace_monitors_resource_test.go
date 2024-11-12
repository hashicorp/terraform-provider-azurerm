// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	r.preCheck(t)

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
	r.preCheck(t)

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
	r.preCheck(t)

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
		t.Skipf("DYNATRACE_USER_COUNTRY must be set for acceptance tests")
	}
	if r.dynatraceInfo.UserEmail == "" {
		t.Skipf("DYNATRACE_USER_EMAIL must be set for acceptance tests")
	}
	if r.dynatraceInfo.UserFirstName == "" {
		t.Skipf("DYNATRACE_USER_FIRST_NAME must be set for acceptance tests")
	}
	if r.dynatraceInfo.UserLastName == "" {
		t.Skipf("DYNATRACE_USER_LAST_NAME must be set for acceptance tests")
	}
	if r.dynatraceInfo.UserPhoneNumber == "" {
		t.Skipf("DYNATRACE_USER_PHONE_NUMBER must be set for acceptance tests")
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
    first_name   = "%s"
    last_name    = "%s"
    email        = "%s"
    phone_number = "%s"
    country      = "%s"
  }

  plan {
    usage_type    = "COMMITTED"
    billing_cycle = "MONTHLY"
    plan          = "azureportalintegration_privatepreview@TIDgmz7xq9ge3py"
  }

  tags = {
    environment = "Dev"
  }
}
`, template, data.RandomInteger, r.dynatraceInfo.UserFirstName, r.dynatraceInfo.UserLastName, r.dynatraceInfo.UserEmail, r.dynatraceInfo.UserPhoneNumber, r.dynatraceInfo.UserCountry)
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
    first_name   = "%s"
    last_name    = "%s"
    email        = "%s"
    phone_number = "%s"
    country      = "%s"
  }

  plan {
    usage_type    = "COMMITTED"
    billing_cycle = "MONTHLY"
    plan          = "azureportalintegration_privatepreview@TIDgmz7xq9ge3py"
  }

  tags = {
    environment = "Prod"
    test        = "Patch"
  }
}
`, template, data.RandomInteger, r.dynatraceInfo.UserFirstName, r.dynatraceInfo.UserLastName, r.dynatraceInfo.UserEmail, r.dynatraceInfo.UserPhoneNumber, r.dynatraceInfo.UserCountry)
}

func (r MonitorsResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dynatrace_monitor" "import" {
  name                     = azurerm_dynatrace_monitor.test.name
  resource_group_name      = azurerm_dynatrace_monitor.test.resource_group_name
  location                 = azurerm_dynatrace_monitor.test.location
  marketplace_subscription = azurerm_dynatrace_monitor.test.marketplace_subscription

  identity {
    type = azurerm_dynatrace_monitor.test.identity.0.type
  }
  user {
    first_name   = azurerm_dynatrace_monitor.test.user.0.first_name
    last_name    = azurerm_dynatrace_monitor.test.user.0.last_name
    email        = azurerm_dynatrace_monitor.test.user.0.email
    phone_number = azurerm_dynatrace_monitor.test.user.0.phone_number
    country      = azurerm_dynatrace_monitor.test.user.0.country
  }

  plan {
    usage_type    = azurerm_dynatrace_monitor.test.plan.0.usage_type
    billing_cycle = azurerm_dynatrace_monitor.test.plan.0.billing_cycle
    plan          = azurerm_dynatrace_monitor.test.plan.0.plan
  }
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
