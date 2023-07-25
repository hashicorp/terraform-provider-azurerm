// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package voiceservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/voiceservices/2023-04-03/testlines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CommunicationsGatewayTestLineTestResource struct{}

func TestAccVoiceServicesTestLine_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_voice_services_communications_gateway_test_line", "test")
	r := CommunicationsGatewayTestLineTestResource{}
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

func TestAccVoiceServicesTestLine_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_voice_services_communications_gateway_test_line", "test")
	r := CommunicationsGatewayTestLineTestResource{}
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

func TestAccVoiceServicesTestLine_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_voice_services_communications_gateway_test_line", "test")
	r := CommunicationsGatewayTestLineTestResource{}
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

func TestAccVoiceServicesTestLine_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_voice_services_communications_gateway_test_line", "test")
	r := CommunicationsGatewayTestLineTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CommunicationsGatewayTestLineTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := testlines.ParseTestLineID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.VoiceServices.TestLinesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r CommunicationsGatewayTestLineTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-vscgtl-%d"
  location = "%s"
}

resource "azurerm_voice_services_communications_gateway" "test" {
  name                = "acctest-vscg-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.Locations.Primary)
}

func (r CommunicationsGatewayTestLineTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_voice_services_communications_gateway_test_line" "test" {
  name                                     = "acctest-tl-%s"
  location                                 = "%s"
  voice_services_communications_gateway_id = azurerm_voice_services_communications_gateway.test.id
  phone_number                             = "123456789"
  purpose                                  = "Automated"
}
`, template, data.RandomString, data.Locations.Primary)
}

func (r CommunicationsGatewayTestLineTestResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_voice_services_communications_gateway_test_line" "import" {
  name                                     = azurerm_voice_services_communications_gateway_test_line.test.name
  location                                 = azurerm_voice_services_communications_gateway_test_line.test.location
  voice_services_communications_gateway_id = azurerm_voice_services_communications_gateway_test_line.test.voice_services_communications_gateway_id
  phone_number                             = azurerm_voice_services_communications_gateway_test_line.test.phone_number
  purpose                                  = azurerm_voice_services_communications_gateway_test_line.test.purpose
}
`, config)
}

func (r CommunicationsGatewayTestLineTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_voice_services_communications_gateway_test_line" "test" {
  name                                     = "acctest-tl-%s"
  location                                 = "%s"
  voice_services_communications_gateway_id = azurerm_voice_services_communications_gateway.test.id
  phone_number                             = "123456789"
  purpose                                  = "Automated"
  tags = {
    key = "value"
  }
}
`, template, data.RandomString, data.Locations.Primary)
}

func (r CommunicationsGatewayTestLineTestResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_voice_services_communications_gateway_test_line" "test" {
  name                                     = "acctest-tl-%s"
  location                                 = "%s"
  voice_services_communications_gateway_id = azurerm_voice_services_communications_gateway.test.id
  phone_number                             = "987654321"
  purpose                                  = "Manual"
  tags = {
    key2 = "value2"
  }
}
`, template, data.RandomString, data.Locations.Primary)
}
