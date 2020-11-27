package automation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AutomationConnectionServicePrincipalResource struct {
}

func TestAccAutomationConnectionServicePrincipal_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_service_principal", "test")
	r := AutomationConnectionServicePrincipalResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationConnectionServicePrincipal_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_service_principal", "test")
	r := AutomationConnectionServicePrincipalResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAutomationConnectionServicePrincipal_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_service_principal", "test")
	r := AutomationConnectionServicePrincipalResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationConnectionServicePrincipal_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_service_principal", "test")
	r := AutomationConnectionServicePrincipalResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t AutomationConnectionServicePrincipalResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Automation.ConnectionClient.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Automation Connection (Service Pricipal) %q (resource group: %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ConnectionProperties != nil), nil
}

func (AutomationConnectionServicePrincipalResource) basic(data acceptance.TestData) string {
	template := AutomationConnectionServicePrincipalResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection_service_principal" "test" {
  name                    = "acctestACSP-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  application_id          = "00000000-0000-0000-0000-000000000000"
  tenant_id               = data.azurerm_client_config.test.tenant_id
  subscription_id         = data.azurerm_client_config.test.subscription_id
  certificate_thumbprint  = file("testdata/automation_certificate_test.thumb")
}
`, template, data.RandomInteger)
}

func (AutomationConnectionServicePrincipalResource) requiresImport(data acceptance.TestData) string {
	template := AutomationConnectionServicePrincipalResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection_service_principal" "import" {
  name                    = azurerm_automation_connection_service_principal.test.name
  resource_group_name     = azurerm_automation_connection_service_principal.test.resource_group_name
  automation_account_name = azurerm_automation_connection_service_principal.test.automation_account_name
  application_id          = azurerm_automation_connection_service_principal.test.application_id
  tenant_id               = azurerm_automation_connection_service_principal.test.tenant_id
  subscription_id         = azurerm_automation_connection_service_principal.test.subscription_id
  certificate_thumbprint  = azurerm_automation_connection_service_principal.test.certificate_thumbprint
}
`, template)
}

func (AutomationConnectionServicePrincipalResource) complete(data acceptance.TestData) string {
	template := AutomationConnectionServicePrincipalResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection_service_principal" "test" {
  name                    = "acctestACSP-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  application_id          = "00000000-0000-0000-0000-000000000000"
  tenant_id               = data.azurerm_client_config.test.tenant_id
  subscription_id         = data.azurerm_client_config.test.subscription_id
  certificate_thumbprint  = file("testdata/automation_certificate_test.thumb")
  description             = "acceptance test for automation connection"
}
`, template, data.RandomInteger)
}

func (AutomationConnectionServicePrincipalResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
  location = "%s"
}

data "azurerm_client_config" "test" {}

resource "azurerm_automation_account" "test" {
  name                = "acctestAA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
