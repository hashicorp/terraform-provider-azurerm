package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CustomIpPrefixResource struct {
}

func TestAccCustomIpPrefix_withIpv4(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	if location.NormalizeNilable(utils.String(data.Locations.Primary)) != "eastus2euap" {
		t.Skip("skip as the test ip range is only available in `eastus2euap` region")
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},

		data.ImportStep("action"),
	})
}

func TestAccCustomIpPrefix_ipv4Update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	if location.NormalizeNilable(utils.String(data.Locations.Primary)) != "eastus2euap" {
		t.Skip("skip as the test ip range is only available in `eastus2euap` region")
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("action"),
		{
			Config: r.ipv4Update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("action"),
		{
			Config: r.ipv4(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("action"),
	})
}

func TestAccCustomIpPrefix_ipv4Complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	if location.NormalizeNilable(utils.String(data.Locations.Primary)) != "eastus2euap" {
		t.Skip("skip as the test ip range is only available in `eastus2euap` region")
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4Complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("action"),
	})
}

func TestAccCustomIpPrefix_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (CustomIpPrefixResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.CustomIpPrefixID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.CustomIPPrefixesClient.Get(ctx, id.ResourceGroup, id.CustomIpPrefixeName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (r CustomIpPrefixResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_custom_ip_prefix" "import" {
  name                = azurerm_custom_ip_prefix.test.name
  location            = azurerm_custom_ip_prefix.test.location
  resource_group_name = azurerm_custom_ip_prefix.test.resource_group_name
  cidr                = azurerm_custom_ip_prefix.test.cidr
  zones               = azurerm_custom_ip_prefix.test.zones
}
`, r.ipv4(data))
}

func (r CustomIpPrefixResource) ipv4(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%[1]d"
  location = "%[2]s"
}

resource "azurerm_custom_ip_prefix" "test" {
  name                = "acctest-CustomIpPrefix-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidr                = "194.41.19.0/24"
  zones               = ["1"]
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CustomIpPrefixResource) ipv4Update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%[1]d"
  location = "%[2]s"
}

resource "azurerm_custom_ip_prefix" "test" {
  name                = "acctest-CustomIpPrefix-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidr                = "194.41.19.0/24"
  action              = "Deprovision"
  zones               = ["1"]

  tags = {
    env = "prod"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CustomIpPrefixResource) ipv4Complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%[1]d"
  location = "%[2]s"
}

resource "azurerm_custom_ip_prefix" "test" {
  name                  = "acctest-CustomIpPrefix-%[1]d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  cidr                  = "194.41.19.0/24"
  zones                 = ["1","2","3"]
  action                = "Provision"
  authorization_message = "00000000-0000-0000-0000-000000000000|194.41.19.0/24|20991212"
  signed_message        = "singed message for WAN validation"

  tags = {
    env = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
