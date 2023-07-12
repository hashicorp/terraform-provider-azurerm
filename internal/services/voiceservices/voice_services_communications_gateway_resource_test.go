// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package voiceservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/voiceservices/2023-04-03/communicationsgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VoiceServicesCommunicationsGatewayTestResource struct{}

func TestAccVoiceServicesCommunicationsGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_voice_services_communications_gateway", "test")
	r := VoiceServicesCommunicationsGatewayTestResource{}
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

func TestAccVoiceServicesCommunicationsGateway_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_voice_services_communications_gateway", "test")
	r := VoiceServicesCommunicationsGatewayTestResource{}
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

func TestAccVoiceServicesCommunicationsGateway_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_voice_services_communications_gateway", "test")
	r := VoiceServicesCommunicationsGatewayTestResource{}
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

func TestAccVoiceServicesCommunicationsGateway_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_voice_services_communications_gateway", "test")
	r := VoiceServicesCommunicationsGatewayTestResource{}
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
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r VoiceServicesCommunicationsGatewayTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := communicationsgateways.ParseCommunicationsGatewayID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.VoiceServices.CommunicationsGatewaysClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r VoiceServicesCommunicationsGatewayTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r VoiceServicesCommunicationsGatewayTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

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
`, template, data.RandomString, data.Locations.Primary)
}

func (r VoiceServicesCommunicationsGatewayTestResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_voice_services_communications_gateway" "import" {
  name                = azurerm_voice_services_communications_gateway.test.name
  resource_group_name = azurerm_voice_services_communications_gateway.test.resource_group_name
  location            = azurerm_voice_services_communications_gateway.test.location
  connectivity        = azurerm_voice_services_communications_gateway.test.connectivity
  e911_type           = azurerm_voice_services_communications_gateway.test.e911_type
  codecs              = azurerm_voice_services_communications_gateway.test.codecs
  platforms           = azurerm_voice_services_communications_gateway.test.platforms
  on_prem_mcp_enabled = azurerm_voice_services_communications_gateway.test.on_prem_mcp_enabled

  service_location {
    location                                  = "eastus"
    allowed_media_source_address_prefixes     = ["10.1.2.0/24"]
    allowed_signaling_source_address_prefixes = ["10.1.1.0/24"]
    operator_addresses                        = ["198.51.100.1"]
  }

  service_location {
    location                                  = "eastus2"
    allowed_media_source_address_prefixes     = ["10.2.2.0/24"]
    allowed_signaling_source_address_prefixes = ["10.2.1.0/24"]
    operator_addresses                        = ["198.51.100.2"]
  }
}
`, config)
}

func (r VoiceServicesCommunicationsGatewayTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`

%s

resource "azurerm_voice_services_communications_gateway" "test" {
  name                = "acctest-vscg-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"

  connectivity = "PublicAddress"
  e911_type    = "DirectToEsrp"
  codecs       = "PCMA"
  platforms    = ["OperatorConnect", "TeamsPhoneMobile"]

  service_location {
    location                                  = "eastus"
    allowed_media_source_address_prefixes     = ["10.1.2.0/24"]
    allowed_signaling_source_address_prefixes = ["10.1.1.0/24"]
    esrp_addresses                            = ["198.51.100.3"]
    operator_addresses                        = ["198.51.100.1"]
  }

  service_location {
    location                                  = "eastus2"
    allowed_media_source_address_prefixes     = ["10.2.2.0/24"]
    allowed_signaling_source_address_prefixes = ["10.2.1.0/24"]
    esrp_addresses                            = ["198.51.100.4"]
    operator_addresses                        = ["198.51.100.2"]
  }

  api_bridge                             = jsonencode({})
  auto_generated_domain_name_label_scope = "SubscriptionReuse"
  on_prem_mcp_enabled                    = true
  microsoft_teams_voicemail_pilot_number = "1"
  emergency_dial_strings                 = ["911", "933"]

  tags = {
    Environment = "Test"
  }
}
`, template, data.RandomString, data.Locations.Primary)
}

func (r VoiceServicesCommunicationsGatewayTestResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`


%s

resource "azurerm_voice_services_communications_gateway" "test" {
  name                = "acctest-vscg-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  connectivity        = "PublicAddress"
  e911_type           = "Standard"
  codecs              = "PCMU"
  platforms           = ["OperatorConnect"]

  service_location {
    location                                  = "eastus2"
    allowed_media_source_address_prefixes     = ["10.1.2.0/24"]
    allowed_signaling_source_address_prefixes = ["10.1.1.0/24"]
    operator_addresses                        = ["198.51.100.1"]
  }

  service_location {
    location                                  = "eastus"
    allowed_media_source_address_prefixes     = ["10.2.2.0/24"]
    allowed_signaling_source_address_prefixes = ["10.2.1.0/24"]
    operator_addresses                        = ["198.51.100.2"]
  }

  auto_generated_domain_name_label_scope = "SubscriptionReuse"
  emergency_dial_strings                 = ["911"]
  microsoft_teams_voicemail_pilot_number = "2"

  tags = {
    Environment = "dev"
  }
}
`, template, data.RandomString, data.Locations.Primary)
}
