package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CustomIpPrefixResource struct {
}

func TestAccCustomIpPrefixIpv4(t *testing.T) {
	t.Error("@tombuildsstuff: disabling these tests so we can track down an issue where the specified CIDR isn't released")

	//// Only one test IPv4 range "194.41.20.0/24" could be provided to run tests, and the IP range could only create one resource at a time, so run the tests sequentially.
	//acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
	//	"ipv4": {
	//		"basic":                             testAccCustomIpPrefix_withIpv4,
	//		"update":                            testAccCustomIpPrefix_ipv4Update,
	//		"fromCommissionedTodoProvision":     testAccCustomIpPrefix_ipv4Update_from_commissioned_todo_provision,
	//		"fromDeprovisionedTodoCommission":   testAccCustomIpPrefix_ipv4Update_from_deprovisioned_todo_commission,
	//		"fromDeprovisionedTodoDecommission": testAccCustomIpPrefix_ipv4Update_from_deprovisioned_todo_decommission,
	//		"fromCommissionedTodoDeprovision":   testAccCustomIpPrefix_ipv4Update_from_commissioned_todo_deprovision,
	//		"complete":                          testAccCustomIpPrefix_ipv4Complete,
	//		"requiresImport":                    testAccCustomIpPrefix_requiresImport,
	//	},
	//})
}

func testAccCustomIpPrefix_withIpv4(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},

		data.ImportStep(),
	})
}

func testAccCustomIpPrefix_ipv4Update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ipv4UpdateCommissionedState(data, "Commission"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ipv4UpdateCommissionedState(data, "Decommission"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ipv4UpdateCommissionedState(data, "Deprovision"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccCustomIpPrefix_ipv4Update_from_commissioned_todo_provision(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ipv4UpdateCommissionedState(data, "Commission"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ipv4UpdateCommissionedState(data, "Provision"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccCustomIpPrefix_ipv4Update_from_deprovisioned_todo_commission(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ipv4UpdateCommissionedState(data, "Deprovision"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ipv4UpdateCommissionedState(data, "Commission"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccCustomIpPrefix_ipv4Update_from_deprovisioned_todo_decommission(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ipv4UpdateCommissionedState(data, "Deprovision"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ipv4UpdateCommissionedState(data, "Decommission"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccCustomIpPrefix_ipv4Update_from_commissioned_todo_deprovision(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ipv4UpdateCommissionedState(data, "Commission"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ipv4UpdateCommissionedState(data, "Deprovision"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccCustomIpPrefix_ipv4Complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4Complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccCustomIpPrefix_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
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
  cidr                = "194.41.20.0/24"
  zones               = ["1"]
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CustomIpPrefixResource) ipv4UpdateCommissionedState(data acceptance.TestData, state string) string {
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
  cidr                = "194.41.20.0/24"
  action              = "%[3]s"
  zones               = ["1"]

  tags = {
    env = "prod"
  }
}
`, data.RandomInteger, data.Locations.Primary, state)
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
  name                          = "acctest-CustomIpPrefix-%[1]d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  cidr                          = "194.41.20.0/24"
  zones                         = ["1", "2", "3"]
  action                        = "Provision"
  roa_expiration_date           = "20991212"
  wan_validation_signed_message = "signed message for WAN validation"

  tags = {
    env = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
