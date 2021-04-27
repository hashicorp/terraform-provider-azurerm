package firewall_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/firewall/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type FirewallPolicyResource struct {
}

func TestAccFirewallPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")
	r := FirewallPolicyResource{}

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

func TestAccFirewallPolicy_basicPremium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")
	r := FirewallPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicPremium(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")
	r := FirewallPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")
	r := FirewallPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")
	r := FirewallPolicyResource{}

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

func TestAccFirewallPolicy_inherit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")
	r := FirewallPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.inherit(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (FirewallPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	var id, err = parse.FirewallPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Firewall.FirewallPolicyClient.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.FirewallPolicyPropertiesFormat != nil), nil
}

func (FirewallPolicyResource) basic(data acceptance.TestData) string {
	template := FirewallPolicyResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_policy" "test" {
  name                = "acctest-networkfw-Policy-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func (FirewallPolicyResource) basicPremium(data acceptance.TestData) string {
	template := FirewallPolicyResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_policy" "test" {
  name                = "acctest-networkfw-Policy-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
}
`, template, data.RandomInteger)
}

func (FirewallPolicyResource) complete(data acceptance.TestData) string {
	template := FirewallPolicyResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_policy" "test" {
  name                     = "acctest-networkfw-Policy-%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  threat_intelligence_mode = "Off"
  threat_intelligence_allowlist {
    ip_addresses = ["1.1.1.1", "2.2.2.2"]
    fqdns        = ["foo.com", "bar.com"]
  }
  dns {
    servers       = ["1.1.1.1", "2.2.2.2"]
    proxy_enabled = true
  }
  tags = {
    env = "Test"
  }
}
`, template, data.RandomInteger)
}

func (FirewallPolicyResource) requiresImport(data acceptance.TestData) string {
	template := FirewallPolicyResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_policy" "import" {
  name                = azurerm_firewall_policy.test.name
  resource_group_name = azurerm_firewall_policy.test.resource_group_name
  location            = azurerm_firewall_policy.test.location
}
`, template)
}

func (FirewallPolicyResource) inherit(data acceptance.TestData) string {
	template := FirewallPolicyResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_policy" "test-parent" {
  name                = "acctest-networkfw-Policy-%d-parent"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_firewall_policy" "test" {
  name                = "acctest-networkfw-Policy-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  base_policy_id      = azurerm_firewall_policy.test-parent.id
  threat_intelligence_allowlist {
    ip_addresses = ["1.1.1.1", "2.2.2.2"]
    fqdns        = ["foo.com", "bar.com"]
  }
  dns {
    servers       = ["1.1.1.1", "2.2.2.2"]
    proxy_enabled = true
  }
  tags = {
    env = "Test"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (FirewallPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-networkfw-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
