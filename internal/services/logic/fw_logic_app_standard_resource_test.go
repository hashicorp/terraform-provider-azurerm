package logic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type FWLogicAppStandardResource struct{}

func TestAccFWLogicAppStandard_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fw_logic_app_standard", "test")
	r := FWLogicAppStandardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "kind", "functionapp,workflowapp"),
				check.That(data.ResourceName).Key("version").HasValue("~4"),
				check.That(data.ResourceName).Key("outbound_ip_addresses").Exists(),
				check.That(data.ResourceName).Key("possible_outbound_ip_addresses").Exists(),
				check.That(data.ResourceName).Key("custom_domain_verification_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFWLogicAppStandard_systemIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fw_logic_app_standard", "test")
	r := FWLogicAppStandardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withSystemIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "kind", "functionapp,workflowapp"),
				check.That(data.ResourceName).Key("version").HasValue("~4"),
				check.That(data.ResourceName).Key("outbound_ip_addresses").Exists(),
				check.That(data.ResourceName).Key("possible_outbound_ip_addresses").Exists(),
				check.That(data.ResourceName).Key("custom_domain_verification_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (r FWLogicAppStandardResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseLogicAppId(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppService.WebAppsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r FWLogicAppStandardResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_fw_logic_app_standard" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FWLogicAppStandardResource) withSystemIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_fw_logic_app_standard" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  identity {
    type = "SystemAssigned"
  }

}
`, r.template(data), data.RandomInteger)
}

func (FWLogicAppStandardResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestRGFW-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "WS1"
  os_type             = "Windows"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
