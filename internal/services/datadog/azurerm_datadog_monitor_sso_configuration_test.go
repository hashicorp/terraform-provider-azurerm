// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datadog_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/singlesignon"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SSODatadogMonitorResource struct {
	datadogApiKey         string
	datadogApplicationKey string
	enterpriseAppId       string
}

// Because the test resource interface requires an Exists() function with a value receiver, I refactored this function
// to take a resource as arg rather than receiver in order to comply with the Go recommendation to not make methods on
// both value and pointer receivers.
func NewSSODatadogMonitorResource(t *testing.T) SSODatadogMonitorResource {
	if os.Getenv("ARM_TEST_DATADOG_API_KEY") == "" {
		t.Skip("Skipping as ARM_TEST_DATADOG_API_KEY is not specified")
	}
	if os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY") == "" {
		t.Skip("Skipping as ARM_TEST_DATADOG_APPLICATION_KEY is not specified")
	}
	if os.Getenv("ARM_TEST_ENTERPRISE_APP_ID") == "" {
		t.Skip("Skipping as ARM_TEST_ENTERPRISE_APP_ID is not specified")
	}
	return SSODatadogMonitorResource{
		datadogApiKey:         os.Getenv("ARM_TEST_DATADOG_API_KEY"),
		datadogApplicationKey: os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY"),
		enterpriseAppId:       os.Getenv("ARM_TEST_ENTERPRISE_APP_ID"),
	}
}

func TestAccDatadogMonitorSSO_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor_sso_configuration", "test")
	r := NewSSODatadogMonitorResource(t)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("single_sign_on").HasValue("Enable"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatadogMonitorSSO_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor_sso_configuration", "test")
	r := NewSSODatadogMonitorResource(t)
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

func TestAccDatadogMonitorSSO_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor_sso_configuration", "test")
	r := NewSSODatadogMonitorResource(t)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("single_sign_on").HasValue(string(singlesignon.SingleSignOnStatesEnable)),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, singlesignon.SingleSignOnStatesDisable),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("single_sign_on").HasValue(string(singlesignon.SingleSignOnStatesDisable)),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("single_sign_on").HasValue(string(singlesignon.SingleSignOnStatesEnable)),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, singlesignon.SingleSignOnStatesExisting),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("single_sign_on").HasValue(string(singlesignon.SingleSignOnStatesExisting)),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("single_sign_on").HasValue(string(singlesignon.SingleSignOnStatesEnable)),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, singlesignon.SingleSignOnStatesInitial),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("single_sign_on").HasValue(string(singlesignon.SingleSignOnStatesInitial)),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("single_sign_on").HasValue(string(singlesignon.SingleSignOnStatesEnable)),
			),
		},
		data.ImportStep(),
	})
}

func (SSODatadogMonitorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := singlesignon.ParseSingleSignOnConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Datadog.SingleSignOn.ConfigurationsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r *SSODatadogMonitorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-datadogrg-%[1]d"
  location = %[2]q
}

resource "azurerm_datadog_monitor" "test" {
  name                = "acctest-datadog-%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datadog_organization {
    api_key         = %[4]q
    application_key = %[5]q
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, r.datadogApiKey, r.datadogApplicationKey)
}

func (r SSODatadogMonitorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_datadog_monitor_sso_configuration" "test" {
  datadog_monitor_id        = azurerm_datadog_monitor.test.id
  single_sign_on            = "Enable"
  enterprise_application_id = %q
}
`, r.template(data), r.enterpriseAppId)
}

func (r SSODatadogMonitorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_datadog_monitor_sso_configuration" "import" {
  datadog_monitor_id        = azurerm_datadog_monitor_sso_configuration.test.datadog_monitor_id
  single_sign_on            = azurerm_datadog_monitor_sso_configuration.test.single_sign_on
  enterprise_application_id = azurerm_datadog_monitor_sso_configuration.test.enterprise_application_id
}
`, r.basic(data))
}

func (r SSODatadogMonitorResource) update(data acceptance.TestData, ssoState singlesignon.SingleSignOnStates) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_datadog_monitor_sso_configuration" "test" {
  datadog_monitor_id        = azurerm_datadog_monitor.test.id
  single_sign_on            = %s
  enterprise_application_id = %q
}
`, r.template(data), ssoState, r.enterpriseAppId)
}
