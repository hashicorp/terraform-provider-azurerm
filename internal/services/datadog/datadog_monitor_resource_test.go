// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datadog_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/monitorsresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DatadogMonitorResource struct {
	datadogApiKey         string
	datadogApplicationKey string
}

func (r *DatadogMonitorResource) populateFromEnvironment(t *testing.T) {
	if os.Getenv("ARM_TEST_DATADOG_API_KEY") == "" {
		t.Skip("Skipping as ARM_TEST_DATADOG_API_KEY is not specified")
	}
	if os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY") == "" {
		t.Skip("Skipping as ARM_TEST_DATADOG_APPLICATION_KEY is not specified")
	}
	r.datadogApiKey = os.Getenv("ARM_TEST_DATADOG_API_KEY")
	r.datadogApplicationKey = os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY")
}

func TestAccDatadogMonitor_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor", "test")
	r := DatadogMonitorResource{}
	r.populateFromEnvironment(t)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user",
			"user.0.name",
			"user.0.email",
			"datadog_organization",
			"datadog_organization.0",
			"datadog_organization.0.id",
			"datadog_organization.0.name",
			"datadog_organization.0.api_key",
			"datadog_organization.0.application_key",
			"datadog_organization.0.enterprise_app_id",
			"datadog_organization.0.linking_auth_code",
			"datadog_organization.0.linking_client_id",
			"datadog_organization.0.redirect_uri"),
	})
}

func TestAccDatadogMonitor_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor", "test")
	r := DatadogMonitorResource{}
	r.populateFromEnvironment(t)
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

func TestAccDatadogMonitor_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor", "test")
	r := DatadogMonitorResource{}
	r.populateFromEnvironment(t)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user",
			"user.0.name",
			"user.0.email",
			"datadog_organization",
			"datadog_organization.0",
			"datadog_organization.0.id",
			"datadog_organization.0.name",
			"datadog_organization.0.api_key",
			"datadog_organization.0.application_key",
			"datadog_organization.0.enterprise_app_id",
			"datadog_organization.0.linking_auth_code",
			"datadog_organization.0.linking_client_id",
			"datadog_organization.0.redirect_uri"),
	})
}

func TestAccDatadogMonitor_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor", "test")
	r := DatadogMonitorResource{}
	r.populateFromEnvironment(t)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user",
			"user.0.name",
			"user.0.email",
			"datadog_organization",
			"datadog_organization.0",
			"datadog_organization.0.id",
			"datadog_organization.0.name",
			"datadog_organization.0.api_key",
			"datadog_organization.0.application_key",
			"datadog_organization.0.enterprise_app_id",
			"datadog_organization.0.linking_auth_code",
			"datadog_organization.0.linking_client_id",
			"datadog_organization.0.redirect_uri"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user",
			"user.0.name",
			"user.0.email",
			"datadog_organization",
			"datadog_organization.0",
			"datadog_organization.0.id",
			"datadog_organization.0.name",
			"datadog_organization.0.api_key",
			"datadog_organization.0.application_key",
			"datadog_organization.0.enterprise_app_id",
			"datadog_organization.0.linking_auth_code",
			"datadog_organization.0.linking_client_id",
			"datadog_organization.0.redirect_uri"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user",
			"user.0.name",
			"user.0.email",
			"datadog_organization",
			"datadog_organization.0",
			"datadog_organization.0.id",
			"datadog_organization.0.name",
			"datadog_organization.0.api_key",
			"datadog_organization.0.application_key",
			"datadog_organization.0.enterprise_app_id",
			"datadog_organization.0.linking_auth_code",
			"datadog_organization.0.linking_client_id",
			"datadog_organization.0.redirect_uri"),
	})
}

func (r DatadogMonitorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := monitorsresource.ParseMonitorID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Datadog.MonitorsResource.MonitorsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r DatadogMonitorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-datadogrg-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r DatadogMonitorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_datadog_monitor" "test" {
  name                = "acctest-datadog-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
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
`, r.template(data), data.RandomString, os.Getenv("ARM_TEST_DATADOG_API_KEY"), os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY"))
}

func (r DatadogMonitorResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_datadog_monitor" "test" {
  name                = "acctest-datadog-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
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
  monitoring_enabled = false
  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomString, r.datadogApiKey, r.datadogApplicationKey)
}

func (r DatadogMonitorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s
resource "azurerm_datadog_monitor" "import" {
  name                = azurerm_datadog_monitor.test.name
  resource_group_name = azurerm_datadog_monitor.test.resource_group_name
  location            = azurerm_datadog_monitor.test.location
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
`, r.basic(data), r.datadogApiKey, r.datadogApplicationKey)
}

func (r DatadogMonitorResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s
resource "azurerm_datadog_monitor" "test" {
  name                = "acctest-datadog-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datadog_organization {
    api_key           = %q
    application_key   = %q
    enterprise_app_id = ""
    linking_auth_code = ""
    linking_client_id = ""
    redirect_uri      = ""
  }
  identity {
    type = "SystemAssigned"
  }
  sku_name = "Linked"
  user {
    name         = "Test Datadog"
    email        = "abc@xyz.com"
    phone_number = ""
  }
  monitoring_enabled = true
  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomString, r.datadogApiKey, r.datadogApplicationKey)
}
