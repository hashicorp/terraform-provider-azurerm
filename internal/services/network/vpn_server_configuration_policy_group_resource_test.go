// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VPNServerConfigurationPolicyGroupResource struct{}

func TestAccVPNServerConfigurationPolicyGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_server_configuration_policy_group", "test")
	r := VPNServerConfigurationPolicyGroupResource{}

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

func TestAccVPNServerConfigurationPolicyGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_server_configuration_policy_group", "test")
	r := VPNServerConfigurationPolicyGroupResource{}

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

func TestAccVPNServerConfigurationPolicyGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_server_configuration_policy_group", "test")
	r := VPNServerConfigurationPolicyGroupResource{}

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

func TestAccVPNServerConfigurationPolicyGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_server_configuration_policy_group", "test")
	r := VPNServerConfigurationPolicyGroupResource{}

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

func TestAccVPNServerConfigurationPolicyGroup_multiplePolicyGroups(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_server_configuration_policy_group", "test")
	r := VPNServerConfigurationPolicyGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiplePolicyGroups(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r VPNServerConfigurationPolicyGroupResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := virtualwans.ParseConfigurationPolicyGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.VirtualWANs.ConfigurationPolicyGroupsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Vpn Server Configuration Policy Group (%s): %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r VPNServerConfigurationPolicyGroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_server_configuration_policy_group" "test" {
  name                        = "acctestVPNSCPG-%d"
  vpn_server_configuration_id = azurerm_vpn_server_configuration.test.id

  policy {
    name  = "policy1"
    type  = "RadiusAzureGroupId"
    value = "6ad1bd08"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VPNServerConfigurationPolicyGroupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_server_configuration_policy_group" "import" {
  name                        = azurerm_vpn_server_configuration_policy_group.test.name
  vpn_server_configuration_id = azurerm_vpn_server_configuration_policy_group.test.vpn_server_configuration_id

  policy {
    name  = "policy1"
    type  = "RadiusAzureGroupId"
    value = "6ad1bd08"
  }
}
`, r.basic(data))
}

func (r VPNServerConfigurationPolicyGroupResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_server_configuration_policy_group" "test" {
  name                        = "acctestVPNSCPG-%d"
  vpn_server_configuration_id = azurerm_vpn_server_configuration.test.id
  is_default                  = true
  priority                    = 1

  policy {
    name  = "policy1"
    type  = "RadiusAzureGroupId"
    value = "6ad1bd08"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VPNServerConfigurationPolicyGroupResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_server_configuration_policy_group" "test" {
  name                        = "acctestVPNSCPG-%d"
  vpn_server_configuration_id = azurerm_vpn_server_configuration.test.id
  is_default                  = true
  priority                    = 2

  policy {
    name  = "policy2"
    type  = "CertificateGroupId"
    value = "red.com"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VPNServerConfigurationPolicyGroupResource) multiplePolicyGroups(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_server_configuration_policy_group" "test" {
  name                        = "acctestVPNSCPG-%d"
  vpn_server_configuration_id = azurerm_vpn_server_configuration.test.id
  is_default                  = true
  priority                    = 1

  policy {
    name  = "policy1"
    type  = "RadiusAzureGroupId"
    value = "6ad1bd08"
  }
}

resource "azurerm_vpn_server_configuration_policy_group" "test2" {
  name                        = "acctestVPNSCPG2-%d"
  vpn_server_configuration_id = azurerm_vpn_server_configuration.test.id
  is_default                  = false
  priority                    = 2

  policy {
    name  = "policy2"
    type  = "CertificateGroupId"
    value = "red.com"
  }
}

resource "azurerm_vpn_server_configuration_policy_group" "test3" {
  name                        = "acctestVPNSCPG3-%d"
  vpn_server_configuration_id = azurerm_vpn_server_configuration.test.id
  is_default                  = false
  priority                    = 3

  policy {
    name  = "policy3"
    type  = "CertificateGroupId"
    value = "green.com"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r VPNServerConfigurationPolicyGroupResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_vpn_server_configuration" "test" {
  name                     = "acctestVPNSC-%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  vpn_authentication_types = ["Radius"]

  radius {
    server {
      address = "10.105.1.1"
      secret  = "vindicators-the-return-of-worldender"
      score   = 15
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
