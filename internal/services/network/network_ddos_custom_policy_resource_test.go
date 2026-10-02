// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ddoscustompolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkDdosCustomPolicyResource struct{}

func (NetworkDdosCustomPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := ddoscustompolicies.ParseDdosCustomPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.DdosCustomPoliciesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func TestAccNetworkDDoSCustomPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_ddos_custom_policy", "test")
	r := NetworkDdosCustomPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_ip_address_ids.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkDDoSCustomPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_ddos_custom_policy", "test")
	r := NetworkDdosCustomPolicyResource{}

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

func TestAccNetworkDDoSCustomPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_ddos_custom_policy", "test")
	r := NetworkDdosCustomPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("detection_rule.#").HasValue("3"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkDDoSCustomPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_ddos_custom_policy", "test")
	r := NetworkDdosCustomPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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
		{
			// back to basic to verify the additional detection rules and tags can be removed
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("detection_rule.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func (r NetworkDdosCustomPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_network_ddos_custom_policy" "test" {
  name                = "acctest-ddoscp-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  detection_rule {
    name               = "detectionRuleTcp"
    packets_per_second = 1000000
    traffic_type       = "Tcp"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkDdosCustomPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_network_ddos_custom_policy" "import" {
  name                = azurerm_network_ddos_custom_policy.test.name
  resource_group_name = azurerm_network_ddos_custom_policy.test.resource_group_name
  location            = azurerm_network_ddos_custom_policy.test.location

  detection_rule {
    name               = "detectionRuleTcp"
    packets_per_second = 1000000
    traffic_type       = "Tcp"
  }
}
`, r.basic(data))
}

func (r NetworkDdosCustomPolicyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_network_ddos_custom_policy" "test" {
  name                = "acctest-ddoscp-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  detection_rule {
    name               = "detectionRuleTcp"
    packets_per_second = 2000000
    traffic_type       = "Tcp"
  }

  detection_rule {
    name               = "detectionRuleUdp"
    packets_per_second = 200000
    traffic_type       = "Udp"
  }

  detection_rule {
    name               = "detectionRuleTcpSyn"
    packets_per_second = 1000
    traffic_type       = "TcpSyn"
  }

  tags = {
    environment = "acctest"
  }
}
`, r.template(data), data.RandomInteger)
}

func (NetworkDdosCustomPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ddoscp-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}
