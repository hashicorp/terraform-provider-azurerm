// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FunctionAppSlotHybridConnectionResource struct{}

func TestAccFunctionAppSlotHybridConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot_hybrid_connection", "test")
	r := FunctionAppSlotHybridConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionAppSlotHybridConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot_hybrid_connection", "test")
	r := FunctionAppSlotHybridConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionAppSlotHybridConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot_hybrid_connection", "test")
	r := FunctionAppSlotHybridConnectionResource{}

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

func TestAccFunctionAppSlotHybridConnection_sendRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot_hybrid_connection", "test")
	r := FunctionAppSlotHybridConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sendRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}
func TestAccFunctionAppSlotHybridConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot_hybrid_connection", "test")
	r := FunctionAppSlotHybridConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r FunctionAppSlotHybridConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := webapps.ParseSlotHybridConnectionNamespaceRelayID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppService.WebAppsClient.GetHybridConnectionSlot(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Windows %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (r FunctionAppSlotHybridConnectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_slot_hybrid_connection" "test" {
  name            = azurerm_windows_function_app_slot.test.name
  function_app_id = azurerm_windows_function_app.test.id
  relay_id        = azurerm_relay_hybrid_connection.test.id
  hostname        = "acctest%[2]s.hostname"
  port            = 8081
}
`, r.template(data), data.RandomStringOfLength(8))
}

func (r FunctionAppSlotHybridConnectionResource) basicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_slot_hybrid_connection" "test" {
  name            = azurerm_windows_function_app_slot.test.name
  function_app_id = azurerm_windows_function_app.test.id
  relay_id        = azurerm_relay_hybrid_connection.test.id
  hostname        = "acctest%[2]s.anothername"
  port            = 8888
}
`, r.template(data), data.RandomStringOfLength(8))
}

func (r FunctionAppSlotHybridConnectionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_function_app_slot_hybrid_connection" "import" {
  name            = azurerm_function_app_slot_hybrid_connection.test.name
  function_app_id = azurerm_function_app_slot_hybrid_connection.test.function_app_id
  relay_id        = azurerm_function_app_slot_hybrid_connection.test.relay_id
  hostname        = azurerm_function_app_slot_hybrid_connection.test.hostname
  port            = azurerm_function_app_slot_hybrid_connection.test.port
}
`, r.basic(data))
}

func (r FunctionAppSlotHybridConnectionResource) sendRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_slot_hybrid_connection" "test" {
  name            = azurerm_windows_function_app_slot.test.name
  function_app_id = azurerm_windows_function_app.test.id
  relay_id        = azurerm_relay_hybrid_connection.test.id
  hostname        = "acctest%[2]s.hostname"
  port            = 8081

  send_key_name = azurerm_relay_hybrid_connection_authorization_rule.test.name
}
`, r.authRuleTemplate(data), data.RandomStringOfLength(8))
}

func (r FunctionAppSlotHybridConnectionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_function_app_slot_hybrid_connection" "test" {
  name            = azurerm_windows_function_app_slot.test.name
  function_app_id = azurerm_windows_function_app.test.id
  relay_id        = azurerm_relay_hybrid_connection.test.id
  hostname        = "acctest%[2]s.hostname"
  port            = 8081

  send_key_name = azurerm_relay_hybrid_connection_authorization_rule.test.name
}
`, r.authRuleInRemoteResourceGroupTemplate(data), data.RandomStringOfLength(8))
}

func (r FunctionAppSlotHybridConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Windows"
  sku_name            = "%[3]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[4]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_relay_namespace" "test" {
  name                = "acctest-RN-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "Standard"
}

resource "azurerm_relay_hybrid_connection" "test" {
  name                 = "acctest-RHC-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  relay_namespace_name = azurerm_relay_namespace.test.name
  user_metadata        = "metadatatest"
}

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}
}

resource "azurerm_windows_function_app_slot" "test" {
  name            = "slot"
  function_app_id = azurerm_windows_function_app.test.id

  storage_account_name = azurerm_storage_account.test.name

  site_config {}
}
`, data.RandomInteger, data.Locations.Primary, SkuStandardPlan, data.RandomString)
}

func (r FunctionAppSlotHybridConnectionResource) templateRelayInOtherResourceGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Windows"
  sku_name            = "%[3]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[4]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_resource_group" "rg-test-relay" {
  name     = "acctestRG-relay-%[1]d"
  location = "%[2]s"
}

resource "azurerm_relay_namespace" "test" {
  name                = "acctest-RN-%[1]d"
  location            = azurerm_resource_group.rg-test-relay.location
  resource_group_name = azurerm_resource_group.rg-test-relay.name

  sku_name = "Standard"
}

resource "azurerm_relay_hybrid_connection" "test" {
  name                 = "acctest-RHC-%[1]d"
  resource_group_name  = azurerm_resource_group.rg-test-relay.name
  relay_namespace_name = azurerm_relay_namespace.test.name
  user_metadata        = "metadatatest"
}

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}
}

resource "azurerm_windows_function_app_slot" "test" {
  name            = "slot"
  function_app_id = azurerm_windows_function_app.test.id

  storage_account_name = azurerm_storage_account.test.name

  site_config {}
}
`, data.RandomInteger, data.Locations.Primary, SkuStandardPlan, data.RandomString)
}

func (r FunctionAppSlotHybridConnectionResource) authRuleTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_relay_hybrid_connection_authorization_rule" "test" {
  name                   = "sendKey"
  resource_group_name    = azurerm_resource_group.test.name
  hybrid_connection_name = azurerm_relay_hybrid_connection.test.name
  namespace_name         = azurerm_relay_namespace.test.name

  listen = true
  send   = true
  manage = false
}

`, r.template(data))
}

func (r FunctionAppSlotHybridConnectionResource) authRuleInRemoteResourceGroupTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_relay_hybrid_connection_authorization_rule" "test" {
  name                   = "sendKey"
  resource_group_name    = azurerm_resource_group.rg-test-relay.name
  hybrid_connection_name = azurerm_relay_hybrid_connection.test.name
  namespace_name         = azurerm_relay_namespace.test.name

  listen = true
  send   = true
  manage = false
}

`, r.templateRelayInOtherResourceGroup(data))
}
