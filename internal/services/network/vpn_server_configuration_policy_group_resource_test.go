package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VPNServerConfigurationPolicyGroupResource struct{}

func TestAccVPNServerConfigurationPolicyGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_server_configuration_policy_group", "test")
	r := VPNServerConfigurationPolicyGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVPNServerConfigurationPolicyGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_server_configuration_policy_group", "test")
	r := VPNServerConfigurationPolicyGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r VPNServerConfigurationPolicyGroupResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.VpnServerConfigurationPolicyGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.ConfigurationPolicyGroupClient.Get(ctx, id.ResourceGroup, id.VpnServerConfigurationName, id.ConfigurationPolicyGroupName)
	if err != nil {
		return nil, fmt.Errorf("reading Vpn Server Configuration Policy Group (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
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

func (r VPNServerConfigurationPolicyGroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_server_configuration_policy_group" "test" {
  name = "acctest-ncpg-%d"
  resource_group_name = azurerm_resource_group.test.name
  vpn_server_configuration_id = azurerm_network_vpn_server_configuration.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r VPNServerConfigurationPolicyGroupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_configuration_policy_group" "import" {
  name = azurerm_network_configuration_policy_group.test.name
  resource_group_name = azurerm_network_configuration_policy_group.test.resource_group_name
  vpn_server_configuration_id = azurerm_network_configuration_policy_group.test.id
}
`, r.basic(data))
}
