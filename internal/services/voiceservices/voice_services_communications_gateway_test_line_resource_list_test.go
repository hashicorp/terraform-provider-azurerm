// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package voiceservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccVoiceServicesTestLine_list_basic(t *testing.T) {
	r := VoiceServicesCommunicationsGatewayTestLineResource{}
	listResourceAddress := "azurerm_voice_services_communications_gateway_test_line.list"

	data := acceptance.BuildTestData(t, "azurerm_voice_services_communications_gateway_test_line", "test")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicQueryByCommunicationsGatewayId(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r VoiceServicesCommunicationsGatewayTestLineResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-vscgtl-%[1]d"
  location = "%[2]s"
}

resource "azurerm_voice_services_communications_gateway" "test" {
  name                = "acctest-vscg-%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  connectivity        = "PublicAddress"
  e911_type           = "Standard"
  codecs              = "PCMA"
  platforms           = ["OperatorConnect"]
  on_prem_mcp_enabled = false

  service_location {
    location           = "eastus"
    operator_addresses = ["198.51.100.1"]
  }

  service_location {
    location           = "eastus2"
    operator_addresses = ["198.51.100.2"]
  }
}

resource "azurerm_voice_services_communications_gateway_test_line" "test" {
  count = 3

  name                                     = "acctest-tl-${count.index}-%[3]s"
  location                                 = azurerm_resource_group.test.location
  voice_services_communications_gateway_id = azurerm_voice_services_communications_gateway.test.id
  phone_number                             = "12345678${count.index}"
  purpose                                  = "Automated"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r VoiceServicesCommunicationsGatewayTestLineResource) basicQueryByCommunicationsGatewayId(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_voice_services_communications_gateway_test_line" "list" {
  provider = azurerm
  config {
    voice_services_communications_gateway_id = "/subscriptions/%[1]s/resourceGroups/acctest-vscgtl-%[2]d/providers/Microsoft.VoiceServices/communicationsGateways/acctest-vscg-%[3]s"
  }
}
`, data.Subscriptions.Primary, data.RandomInteger, data.RandomString)
}
