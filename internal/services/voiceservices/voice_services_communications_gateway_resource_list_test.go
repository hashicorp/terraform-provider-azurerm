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

func TestAccVoiceServicesCommunicationsGateway_list_basic(t *testing.T) {
	r := VoiceServicesCommunicationsGatewayResource{}
	listResourceAddress := "azurerm_voice_services_communications_gateway.list"

	data := acceptance.BuildTestData(t, "azurerm_voice_services_communications_gateway", "test")

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
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 3),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryByResourceGroupName(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r VoiceServicesCommunicationsGatewayResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_voice_services_communications_gateway" "test" {
  count = 3

  name                = "acctest-vscg-${count.index}-%[3]s"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r VoiceServicesCommunicationsGatewayResource) basicQuery() string {
	return `
list "azurerm_voice_services_communications_gateway" "list" {
  provider = azurerm
  config {}
}
`
}

func (r VoiceServicesCommunicationsGatewayResource) basicQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_voice_services_communications_gateway" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctest-rg-%[1]d"
  }
}
`, data.RandomInteger)
}
