// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package paloalto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/firewalls"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NextGenerationFirewallVnetResource struct{}

func TestAccPaloAltoNextGenerationFirewallLocalRulestackVNet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack", "test")

	r := NextGenerationFirewallVnetResource{}

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

func TestAccPaloAltoNextGenerationFirewallLocalRulestackVNet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack", "test")

	r := NextGenerationFirewallVnetResource{}

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

func TestAccPaloAltoNextGenerationFirewallLocalRulestackVNet_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack", "test")

	r := NextGenerationFirewallVnetResource{}

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

func TestAccPaloAltoNextGenerationFirewallLocalRulestackVNet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack", "test")

	r := NextGenerationFirewallVnetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_profile.0.egress_nat_ip_address_ids.%").HasValue("0"),
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
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r NextGenerationFirewallVnetResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := firewalls.ParseFirewallID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.PaloAlto.PaloAltoClient_v2023_09_01.Firewalls.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r NextGenerationFirewallVnetResource) basic(data acceptance.TestData) string {
	if !features.FivePointOh() {
		return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack" "test" {
  name                = "acctest-ngfwvn-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  rulestack_id        = azurerm_palo_alto_local_rulestack.test.id
  plan_id             = "panw-cngfw-payg"

  network_profile {
    public_ip_address_ids = [azurerm_public_ip.test.id]

    vnet_configuration {
      virtual_network_id  = azurerm_virtual_network.test.id
      trusted_subnet_id   = azurerm_subnet.test1.id
      untrusted_subnet_id = azurerm_subnet.test2.id
    }
  }

  depends_on = [azurerm_palo_alto_local_rulestack_rule.test]
}
`, r.template(data), data.RandomInteger)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack" "test" {
  name                = "acctest-ngfwvn-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  rulestack_id        = azurerm_palo_alto_local_rulestack.test.id

  network_profile {
    public_ip_address_ids = [azurerm_public_ip.test.id]

    vnet_configuration {
      virtual_network_id  = azurerm_virtual_network.test.id
      trusted_subnet_id   = azurerm_subnet.test1.id
      untrusted_subnet_id = azurerm_subnet.test2.id
    }
  }

  depends_on = [azurerm_palo_alto_local_rulestack_rule.test]
}
`, r.template(data), data.RandomInteger)
}

func (r NextGenerationFirewallVnetResource) requiresImport(data acceptance.TestData) string {
	if !features.FivePointOh() {
		return fmt.Sprintf(`
%[1]s

resource "azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack" "import" {
  name                = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.name
  resource_group_name = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.resource_group_name
  rulestack_id        = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.rulestack_id
  plan_id             = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.plan_id

  network_profile {
    public_ip_address_ids = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.network_profile.0.public_ip_address_ids

    vnet_configuration {
      virtual_network_id  = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.network_profile.0.vnet_configuration.0.virtual_network_id
      trusted_subnet_id   = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.network_profile.0.vnet_configuration.0.trusted_subnet_id
      untrusted_subnet_id = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.network_profile.0.vnet_configuration.0.untrusted_subnet_id
    }
  }
}
`, r.basic(data))
	}
	return fmt.Sprintf(`
%[1]s

resource "azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack" "import" {
  name                = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.name
  resource_group_name = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.resource_group_name
  rulestack_id        = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.rulestack_id

  network_profile {
    public_ip_address_ids = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.network_profile.0.public_ip_address_ids

    vnet_configuration {
      virtual_network_id  = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.network_profile.0.vnet_configuration.0.virtual_network_id
      trusted_subnet_id   = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.network_profile.0.vnet_configuration.0.trusted_subnet_id
      untrusted_subnet_id = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.network_profile.0.vnet_configuration.0.untrusted_subnet_id
    }
  }
}
`, r.basic(data))
}

func (r NextGenerationFirewallVnetResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack" "test" {
  name                 = "acctest-ngfwvn-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  rulestack_id         = azurerm_palo_alto_local_rulestack.test.id
  marketplace_offer_id = "pan_swfw_cloud_ngfw"
  plan_id              = "panw-cngfw-payg"

  network_profile {
    public_ip_address_ids     = [azurerm_public_ip.test.id]
    egress_nat_ip_address_ids = [azurerm_public_ip.egress.id]
    trusted_address_ranges    = ["20.22.92.11", "20.23.92.11"]

    vnet_configuration {
      virtual_network_id  = azurerm_virtual_network.test.id
      trusted_subnet_id   = azurerm_subnet.test1.id
      untrusted_subnet_id = azurerm_subnet.test2.id
    }
  }

  dns_settings {
    dns_servers = ["8.8.8.8", "8.8.4.4"]
  }

  destination_nat {
    name     = "testDNAT-1"
    protocol = "TCP"
    frontend_config {
      public_ip_address_id = azurerm_public_ip.test.id
      port                 = 8081
    }
    backend_config {
      public_ip_address = "10.0.1.101"
      port              = 18081
    }
  }

  destination_nat {
    name     = "testDNAT-2"
    protocol = "UDP"
    frontend_config {
      public_ip_address_id = azurerm_public_ip.test.id
      port                 = 8082
    }
    backend_config {
      public_ip_address = "10.0.1.102"
      port              = 18082
    }
  }

  depends_on = [azurerm_palo_alto_local_rulestack_rule.test]
}
`, r.template(data), data.RandomInteger)
}

func (r NextGenerationFirewallVnetResource) completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack" "test" {
  name                 = "acctest-ngfwvn-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  rulestack_id         = azurerm_palo_alto_local_rulestack.test.id
  marketplace_offer_id = "pan_swfw_cloud_ngfw"
  plan_id              = "panw-cngfw-payg"

  network_profile {
    public_ip_address_ids     = [azurerm_public_ip.test.id]
    egress_nat_ip_address_ids = [azurerm_public_ip.egress.id]
    trusted_address_ranges    = ["20.22.92.11", "20.23.92.11"]

    vnet_configuration {
      virtual_network_id  = azurerm_virtual_network.test.id
      trusted_subnet_id   = azurerm_subnet.test1.id
      untrusted_subnet_id = azurerm_subnet.test2.id
    }
  }

  dns_settings {
    use_azure_dns = true
  }

  destination_nat {
    name     = "testDNAT-2"
    protocol = "UDP"
    frontend_config {
      public_ip_address_id = azurerm_public_ip.test.id
      port                 = 8082
    }
    backend_config {
      public_ip_address = "10.0.1.102"
      port              = 18082
    }
  }

  depends_on = [azurerm_palo_alto_local_rulestack_rule.test]
}
`, r.template(data), data.RandomInteger)
}

func (r NextGenerationFirewallVnetResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-PANGFWVN-%[1]d"
  location = "%[2]s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"

  depends_on = [azurerm_public_ip.egress]
}

resource "azurerm_public_ip" "egress" {
  name                = "acctestpublicip-%[1]d-e"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "Production"
  }
}

resource "azurerm_subnet" "test1" {
  name                 = "acctest-pangfw-%[1]d-1"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "trusted"

    service_delegation {
      name = "PaloAltoNetworks.Cloudngfw/firewalls"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "test1" {
  subnet_id                 = azurerm_subnet.test1.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_subnet" "test2" {
  name                 = "acctest-pangfw-%[1]d-2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "untrusted"

    service_delegation {
      name = "PaloAltoNetworks.Cloudngfw/firewalls"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "test2" {
  subnet_id                 = azurerm_subnet.test2.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_palo_alto_local_rulestack" "test" {
  name                = "testAcc-palrs-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"

  depends_on = [azurerm_subnet_network_security_group_association.test1, azurerm_subnet_network_security_group_association.test2]
}

resource "azurerm_palo_alto_local_rulestack_rule" "test" {
  name         = "testacc-palr-%[1]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id
  priority     = 1001
  action       = "Allow"
  protocol     = "application-default"
  applications = ["any"]

  destination {
    cidrs = ["any"]
  }

  source {
    cidrs = ["any"]
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
