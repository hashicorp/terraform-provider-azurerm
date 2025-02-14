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

type WebAppSlotHybridConnectionResource struct{}

func TestAccWebAppSlotHybridConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_app_slot_hybrid_connection", "test")
	r := WebAppSlotHybridConnectionResource{}

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

func TestAccWebAppSlotHybridConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_app_slot_hybrid_connection", "test")
	r := WebAppSlotHybridConnectionResource{}

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

func TestAccWebAppSlotHybridConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_app_slot_hybrid_connection", "test")
	r := WebAppSlotHybridConnectionResource{}

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

func TestAccWebAppSlotHybridConnection_sendRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_app_slot_hybrid_connection", "test")
	r := WebAppSlotHybridConnectionResource{}

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

func TestAccWebAppSlotHybridConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_app_slot_hybrid_connection", "test")
	r := WebAppSlotHybridConnectionResource{}

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

func (r WebAppSlotHybridConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r WebAppSlotHybridConnectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_web_app_hybrid_connection" "test" {
  web_app_id = azurerm_windows_web_app.test.id
  relay_id   = azurerm_relay_hybrid_connection.test.id
  hostname   = "acctest%[2]s.hostname"
  port       = 8081
}

resource "azurerm_web_app_slot_hybrid_connection" "test" {
  name       = azurerm_windows_web_app_slot.test.name
  web_app_id = azurerm_windows_web_app.test.id
  relay_id   = azurerm_relay_hybrid_connection.test.id
  hostname   = "acctest%[2]s.hostname"
  port       = 8081
}
`, r.template(data), data.RandomStringOfLength(8))
}

func (r WebAppSlotHybridConnectionResource) basicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_web_app_slot_hybrid_connection" "test" {
  name       = azurerm_windows_web_app_slot.test.name
  web_app_id = azurerm_windows_web_app.test.id
  relay_id   = azurerm_relay_hybrid_connection.test.id
  hostname   = "acctest%[2]s.anothername"
  port       = 8888
}
`, r.template(data), data.RandomStringOfLength(8))
}

func (r WebAppSlotHybridConnectionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_web_app_slot_hybrid_connection" "import" {
  name       = azurerm_windows_web_app_slot.test.name
  web_app_id = azurerm_web_app_slot_hybrid_connection.test.web_app_id
  relay_id   = azurerm_web_app_slot_hybrid_connection.test.relay_id
  hostname   = azurerm_web_app_slot_hybrid_connection.test.hostname
  port       = azurerm_web_app_slot_hybrid_connection.test.port
}
`, r.basic(data))
}

func (r WebAppSlotHybridConnectionResource) sendRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_web_app_slot_hybrid_connection" "test" {
  name       = azurerm_windows_web_app_slot.test.name
  web_app_id = azurerm_windows_web_app.test.id
  relay_id   = azurerm_relay_hybrid_connection.test.id
  hostname   = "acctest%[2]s.hostname"
  port       = 8081

  send_key_name = azurerm_relay_hybrid_connection_authorization_rule.test.name
}
`, r.authRuleTemplate(data), data.RandomStringOfLength(8))
}

func (r WebAppSlotHybridConnectionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_web_app_slot_hybrid_connection" "test" {
  name       = azurerm_windows_web_app_slot.test.name
  web_app_id = azurerm_windows_web_app.test.id
  relay_id   = azurerm_relay_hybrid_connection.test.id
  hostname   = "acctest%[2]s.hostname"
  port       = 8081

  send_key_name = azurerm_relay_hybrid_connection_authorization_rule.test.name
}
`, r.authRuleInRemoteResourceGroupTemplate(data), data.RandomStringOfLength(8))
}

func (r WebAppSlotHybridConnectionResource) template(data acceptance.TestData) string {
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

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}
}

resource "azurerm_windows_web_app_slot" "test" {
  name           = "staging"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}
}
`, data.RandomInteger, data.Locations.Primary, SkuStandardPlan)
}

func (r WebAppSlotHybridConnectionResource) templateRelayInOtherResourceGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg-test-relay" {
  name     = "acctestRG-%[1]d-relay"
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

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}
}

resource "azurerm_windows_web_app_slot" "test" {
  name           = "staging"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}
}
`, data.RandomInteger, data.Locations.Primary, SkuStandardPlan)
}

func (r WebAppSlotHybridConnectionResource) authRuleTemplate(data acceptance.TestData) string {
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

func (r WebAppSlotHybridConnectionResource) authRuleInRemoteResourceGroupTemplate(data acceptance.TestData) string {
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
