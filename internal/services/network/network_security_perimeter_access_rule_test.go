package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeteraccessrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkSecurityPerimeterAccessRuleTestResource struct{}

func TestAccNetworkSecurityPerimeterAccessRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_perimeter_access_rule", "test")
	r := NetworkSecurityPerimeterAccessRuleTestResource{}

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

func TestAccNetworkSecurityPerimeterAccessRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_perimeter_access_rule", "test")
	r := NetworkSecurityPerimeterAccessRuleTestResource{}

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

func TestAccNetworkSecurityPerimeterAccessRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_perimeter_access_rule", "test")
	r := NetworkSecurityPerimeterAccessRuleTestResource{}

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

func TestAccNetworkSecurityPerimeterAccessRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_perimeter_access_rule", "test")
	r := NetworkSecurityPerimeterAccessRuleTestResource{}

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
				check.That(data.ResourceName).Key("address_prefixes.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkSecurityPerimeterAccessRule_subscriptionIds(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_perimeter_access_rule", "test")
	r := NetworkSecurityPerimeterAccessRuleTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.subscriptionIds(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subscription_ids.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkSecurityPerimeterAccessRule_fdqns(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_perimeter_access_rule", "test")
	r := NetworkSecurityPerimeterAccessRuleTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fqdns(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdns.#").HasValue("3"),
			),
		},
		data.ImportStep(),
	})
}

func (NetworkSecurityPerimeterAccessRuleTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := networksecurityperimeteraccessrules.ParseAccessRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.NetworkSecurityPerimeterAccessRulesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (NetworkSecurityPerimeterAccessRuleTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_perimeter" "test" {
  name     = "acctestNsp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = "%s"
}

resource "azurerm_network_security_perimeter_profile" "test" {
	name = "acctestProfile-%d"
	perimeter_id = azurerm_network_security_perimeter.test.id
}

resource "azurerm_network_security_perimeter_access_rule" "test" {
	name = "acctestRule-%d"
	profile_id = azurerm_network_security_perimeter_profile.test.id

	direction = "Inbound"
	address_prefixes = [
		"8.8.8.8/32"
	]	
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r NetworkSecurityPerimeterAccessRuleTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_security_perimeter_access_rule" "import" {
	name = azurerm_network_security_perimeter_access_rule.test.name
	profile_id =azurerm_network_security_perimeter_access_rule.test.profile_id

	direction = "Inbound"
	address_prefixes = [
		"8.8.8.8/32"
	]	
}
`, r.basic(data))
}

func (NetworkSecurityPerimeterAccessRuleTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_perimeter" "test" {
  name     = "acctestNsp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = "%s"
}

resource "azurerm_network_security_perimeter_profile" "test" {
	name = "acctestProfile-%d"
	perimeter_id = azurerm_network_security_perimeter.test.id
}

resource "azurerm_network_security_perimeter_access_rule" "test" {
	name = "acctestRule-%d"
	profile_id = azurerm_network_security_perimeter_profile.test.id

	direction = "Inbound"
	address_prefixes = [
		"8.8.8.8/32",
		"8.8.4.4/32"
	]	
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (NetworkSecurityPerimeterAccessRuleTestResource) subscriptionIds(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_perimeter" "test" {
  name     = "acctestNsp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = "%s"
}

resource "azurerm_network_security_perimeter_profile" "test" {
	name = "acctestProfile-%d"
	perimeter_id = azurerm_network_security_perimeter.test.id
}

resource "azurerm_network_security_perimeter_access_rule" "test" {
	name = "acctestRule-%d"
	profile_id = azurerm_network_security_perimeter_profile.test.id

	direction = "Inbound"
	subscription_ids = [
		data.azurerm_subscription.current.id
	]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (NetworkSecurityPerimeterAccessRuleTestResource) fqdns(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_perimeter" "test" {
  name     = "acctestNsp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = "%s"
}

resource "azurerm_network_security_perimeter_profile" "test" {
	name = "acctestProfile-%d"
	perimeter_id = azurerm_network_security_perimeter.test.id
}

resource "azurerm_network_security_perimeter_access_rule" "test" {
	name = "acctestRule-%d"
	profile_id = azurerm_network_security_perimeter_profile.test.id

	direction = "Outbound"
	fqdns = [
		"hashicorp.com",
		"google.com",
		"azure.com"
	]

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
