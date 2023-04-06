package network_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"testing"
)

type CustomIpPrefixResource struct{}

const ipv4TestCidr = "194.41.20.0/24"

func TestAccCustomIpPrefixIpv4(t *testing.T) {
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"ipv4": {
			//"basic":        testAccCustomIpPrefix_withIpv4,
			"commissioned": testAccCustomIpPrefix_withIpv4_commissioned,
			//"update":                            testAccCustomIpPrefix_ipv4Update,
			//"fromCommissionedTodoProvision":     testAccCustomIpPrefix_ipv4Update_from_commissioned_todo_provision,
			//"fromDeprovisionedTodoCommission":   testAccCustomIpPrefix_ipv4Update_from_deprovisioned_todo_commission,
			//"fromDeprovisionedTodoDecommission": testAccCustomIpPrefix_ipv4Update_from_deprovisioned_todo_decommission,
			//"fromCommissionedTodoDeprovision":   testAccCustomIpPrefix_ipv4Update_from_commissioned_todo_deprovision,
			//"complete":                          testAccCustomIpPrefix_ipv4Complete,
			//"requiresImport":                    testAccCustomIpPrefix_requiresImport,
		},
	})
}

func testAccCustomIpPrefix_withIpv4(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},

		data.ImportStep(),
	})
}

func testAccCustomIpPrefix_withIpv4_commissioned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4Commissioned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},

		data.ImportStep(),
	})
}

func (CustomIpPrefixResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.CustomIpPrefixID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.CustomIPPrefixesClient.Get(ctx, id.ResourceGroup, id.CustomIpPrefixeName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (r CustomIpPrefixResource) ipv4(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_custom_ip_prefix" "test" {
  name                          = "acctest-%[1]d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  cidr                          = "%[3]s"
  roa_validity_end_date         = "2099-12-12"
  wan_validation_signed_message = "signed message for WAN validation"
  zones                         = ["1"]
}
`, data.RandomInteger, data.Locations.Primary, ipv4TestCidr)
}

func (r CustomIpPrefixResource) ipv4Commissioned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_custom_ip_prefix" "test" {
  name                          = "acctest-%[1]d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  commissioning_enabled         = true
  cidr                          = "%[3]s"
  roa_validity_end_date         = "2099-12-12"
  wan_validation_signed_message = "signed message for WAN validation"
  zones                         = ["1"]
}
`, data.RandomInteger, data.Locations.Primary, ipv4TestCidr)
}
