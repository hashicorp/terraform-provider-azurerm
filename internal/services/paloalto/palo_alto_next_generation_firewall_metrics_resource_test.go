// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package paloalto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	firewalls "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallresources"
	metricsobjectfirewall "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/metricsobjectfirewallresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NextGenerationFirewallMetricsResourceTest struct{}

func (r NextGenerationFirewallMetricsResourceTest) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := firewalls.ParseFirewallID(state.ID)
	if err != nil {
		return nil, err
	}

	metricsFirewallId := metricsobjectfirewall.NewFirewallID(id.SubscriptionId, id.ResourceGroupName, id.FirewallName)

	resp, err := client.PaloAlto.MetricsObjectFirewallResources.MetricsObjectFirewallGet(ctx, metricsFirewallId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving Metrics for %s: %+v", metricsFirewallId, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func TestAccPaloAltoNextGenerationFirewallMetrics_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_next_generation_firewall_metrics", "test")
	r := NextGenerationFirewallMetricsResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("application_insights_connection_string"),
	})
}

func TestAccPaloAltoNextGenerationFirewallMetrics_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_next_generation_firewall_metrics", "test")
	r := NextGenerationFirewallMetricsResourceTest{}

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

func TestAccPaloAltoNextGenerationFirewallMetrics_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_next_generation_firewall_metrics", "test")
	r := NextGenerationFirewallMetricsResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("application_insights_connection_string"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("application_insights_connection_string"),
	})
}

func (r NextGenerationFirewallMetricsResourceTest) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_palo_alto_next_generation_firewall_metrics" "test" {
  firewall_id                            = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.id
  application_insights_connection_string = azurerm_application_insights.test.connection_string
  application_insights_resource_id       = azurerm_application_insights.test.id
}
`, r.template(data))
}

func (r NextGenerationFirewallMetricsResourceTest) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_palo_alto_next_generation_firewall_metrics" "import" {
  firewall_id                            = azurerm_palo_alto_next_generation_firewall_metrics.test.firewall_id
  application_insights_connection_string = azurerm_palo_alto_next_generation_firewall_metrics.test.application_insights_connection_string
  application_insights_resource_id       = azurerm_palo_alto_next_generation_firewall_metrics.test.application_insights_resource_id
}
`, r.basic(data))
}

func (r NextGenerationFirewallMetricsResourceTest) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_palo_alto_next_generation_firewall_metrics" "test" {
  firewall_id                            = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.id
  application_insights_connection_string = azurerm_application_insights.test2.connection_string
  application_insights_resource_id       = azurerm_application_insights.test2.id
}
`, r.template(data))
}

func (r NextGenerationFirewallMetricsResourceTest) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-PANGFWMETRICS-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights" "test2" {
  name                = "acctestappinsights2-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
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
  name                = "acctest-palrs-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"

  depends_on = [azurerm_subnet_network_security_group_association.test1, azurerm_subnet_network_security_group_association.test2]
}

resource "azurerm_palo_alto_local_rulestack_rule" "test" {
  name         = "acctest-palr-%[1]d"
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

resource "azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack" "test" {
  name                = "acctest-ngfwvn-%[1]d"
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
`, data.RandomInteger, data.Locations.Primary)
}
