package datalake_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datalake/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datalake/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DataLakeStoreVirtualNetworkRuleResource struct{}

/*
	---Testing for Success---
	Test a basic Data Lake Store virtual network rule configuration setup and update scenario, and
	validate that new property is set correctly.
*/
func TestAccDataLakeStoreVirtualNetworkRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_store_virtual_network_rule", "test")
	r := DataLakeStoreVirtualNetworkRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withUpdates(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataLakeStoreVirtualNetworkRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_store_virtual_network_rule", "test")
	r := DataLakeStoreVirtualNetworkRuleResource{}

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

/*
	---Testing for Success---
	Test an update to the Data Lake Store Virtual Network Rule to connect to a different subnet, and
	validate that new subnet is set correctly.
*/
func TestAccDataLakeStoreVirtualNetworkRule_switchSubnets(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_store_virtual_network_rule", "test")
	r := DataLakeStoreVirtualNetworkRuleResource{}

	// Create regex strings that will ensure that one subnet name exists, but not the other
	preConfigRegex := regexp.MustCompile(fmt.Sprintf("(subnet1%d)$|(subnet[^2]%d)$", data.RandomInteger, data.RandomInteger))  // subnet 1 but not 2
	postConfigRegex := regexp.MustCompile(fmt.Sprintf("(subnet2%d)$|(subnet[^1]%d)$", data.RandomInteger, data.RandomInteger)) // subnet 2 but not 1

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.subnetSwitchPre(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").MatchesRegex(preConfigRegex),
			),
		},
		{
			Config: r.subnetSwitchPost(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").MatchesRegex(postConfigRegex),
			),
		},
	})
}

/*
	---Testing for Success---
*/
func TestAccDataLakeStoreVirtualNetworkRule_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_store_virtual_network_rule", "test")
	r := DataLakeStoreVirtualNetworkRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

/*
	--Testing for Success--
	Test if we are able to create multiple subnets and connect multiple subnets to the
	Data Lake Store.
*/
func TestAccDataLakeStoreVirtualNetworkRule_multipleSubnets(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_store_virtual_network_rule", "test")
	resourceName2 := "azurerm_data_lake_store_virtual_network_rule.rule2"
	resourceName3 := "azurerm_data_lake_store_virtual_network_rule.rule3"
	r := DataLakeStoreVirtualNetworkRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multipleSubnets(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(resourceName2).ExistsInAzure(r),
				check.That(resourceName3).ExistsInAzure(r),
			),
		},
	})
}

/*
	--Testing for Failure--
	Validation Function Tests - Invalid Name Validations
*/
func TestResourceAzureRMDataLakeStoreVirtualNetworkRule_invalidNameValidation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		// Must only contain alphanumeric characters, periods, underscores or hyphens (4 cases)
		{
			Value:    "test!Rule",
			ErrCount: 1,
		},
		{
			Value:    "test&Rule",
			ErrCount: 1,
		},
		{
			Value:    "test:Rule",
			ErrCount: 1,
		},
		{
			Value:    "test'Rule",
			ErrCount: 1,
		},
		// Cannot be more than 64 characters (1 case - ensure starts with a letter)
		{
			Value:    fmt.Sprintf("v%s", acctest.RandString(64)),
			ErrCount: 1,
		},
		// Cannot be empty (1 case)
		{
			Value:    "",
			ErrCount: 1,
		},
		// Cannot be single character (1 case)
		{
			Value:    "a",
			ErrCount: 1,
		},
		// Cannot end in a hyphen (1 case)
		{
			Value:    "testRule-",
			ErrCount: 1,
		},
		// Cannot end in a period (1 case)
		{
			Value:    "testRule.",
			ErrCount: 1,
		},
		// Cannot start with a period, underscore or hyphen (3 cases)
		{
			Value:    ".testRule",
			ErrCount: 1,
		},
		{
			Value:    "_testRule",
			ErrCount: 1,
		},
		{
			Value:    "-testRule",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validate.ValidateDataLakeStoreVirtualNetworkRuleName(tc.Value, "azurerm_data_lake_store_virtual_network_rule")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Bad: Expected the Azure RM Data Lake Store Virtual Network Rule Name to trigger a validation error.")
		}
	}
}

/*
	--Testing for Success--
	Validation Function Tests - (Barely) Valid Name Validations
*/
func TestResourceAzureRMDataLakeStoreVirtualNetworkRule_validNameValidation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		// Test all lowercase
		{
			Value:    "thisisarule",
			ErrCount: 0,
		},
		// Test all uppercase
		{
			Value:    "THISISARULE",
			ErrCount: 0,
		},
		// Test alternating cases
		{
			Value:    "tHiSiSaRuLe",
			ErrCount: 0,
		},
		// Test hyphens
		{
			Value:    "this-is-a-rule",
			ErrCount: 0,
		},
		// Test multiple hyphens in a row
		{
			Value:    "this----1s----a----ru1e",
			ErrCount: 0,
		},
		// Test underscores
		{
			Value:    "this_is_a_rule",
			ErrCount: 0,
		},
		// Test ending with underscore
		{
			Value:    "this_is_a_rule_",
			ErrCount: 0,
		},
		// Test multiple underscoress in a row
		{
			Value:    "this____1s____a____ru1e",
			ErrCount: 0,
		},
		// Test periods
		{
			Value:    "this.is.a.rule",
			ErrCount: 0,
		},
		// Test multiple periods in a row
		{
			Value:    "this....1s....a....ru1e",
			ErrCount: 0,
		},
		// Test numbers
		{
			Value:    "1108501298509850810258091285091820-5",
			ErrCount: 0,
		},
		// Test a lot of hyphens and numbers
		{
			Value:    "x-5-4-1-2-5-2-6-1-5-2-5-1-2-5-6-2-2",
			ErrCount: 0,
		},
		// Test a lot of underscores and numbers
		{
			Value:    "x_5_4_1_2_5_2_6_1_5_2_5_1_2_5_6_2_2",
			ErrCount: 0,
		},
		// Test a lot of periods and numbers
		{
			Value:    "x.5.4.1.2.5.2.6.1.5.2.5.1.2.5.6.2.2",
			ErrCount: 0,
		},
		// Test exactly 64 characters
		{
			Value:    fmt.Sprintf("v%s", acctest.RandString(63)),
			ErrCount: 0,
		},
	}

	for _, tc := range cases {
		_, errors := validate.ValidateDataLakeStoreVirtualNetworkRuleName(tc.Value, "azurerm_data_lake_store_virtual_network_rule")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Bad: Expected the Azure RM Data Lake Store Virtual Network Rule Name pass name validation successfully but triggered a validation error.")
		}
	}
}

func (r DataLakeStoreVirtualNetworkRuleResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.VirtualNetworkRuleID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Datalake.VirtualNetworkRulesClient.Get(ctx, id.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Data Lake Store Virtual Network Rule %q (Account %q / Resource Group %q): %+v", id.Name, id.AccountName, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r DataLakeStoreVirtualNetworkRuleResource) Destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.VirtualNetworkRuleID(state.ID)
	if err != nil {
		return nil, err
	}
	rulesClient := client.Datalake.VirtualNetworkRulesClient

	if _, err = rulesClient.Delete(ctx, id.ResourceGroup, id.AccountName, id.Name); err != nil {
		return nil, fmt.Errorf("deleting Data Lake Store Virtual Network Rule %q (Account %q / Resource Group %q): %+v", id.Name, id.AccountName, id.ResourceGroup, err)
	}

	return utils.Bool(true), nil
}

/*
	(This test configuration is intended to succeed.)
	Basic Provisioning Configuration
*/
func (r DataLakeStoreVirtualNetworkRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/29"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.7.29.0/29"]
  service_endpoints    = ["Microsoft.AzureActiveDirectory"]
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctestadls%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_data_lake_store_virtual_network_rule" "test" {
  name                = "acctestadlsvnetrule%d"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_data_lake_store.test.name
  subnet_id           = azurerm_subnet.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, strconv.Itoa(data.RandomInteger)[2:13], data.RandomInteger)
}

func (r DataLakeStoreVirtualNetworkRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_store_virtual_network_rule" "import" {
  name                = azurerm_data_lake_store_virtual_network_rule.test.name
  resource_group_name = azurerm_data_lake_store_virtual_network_rule.test.resource_group_name
  account_name        = azurerm_data_lake_store_virtual_network_rule.test.account_name
  subnet_id           = azurerm_data_lake_store_virtual_network_rule.test.subnet_id
}
`, r.basic(data))
}

/*
	(This test configuration is intended to succeed.)
	Basic Provisioning Update Configuration (all other properties would recreate the rule)
*/
func (r DataLakeStoreVirtualNetworkRuleResource) withUpdates(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/29"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.7.29.0/29"]
  service_endpoints    = ["Microsoft.AzureActiveDirectory"]
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctestadls%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_data_lake_store_virtual_network_rule" "test" {
  name                = "acctestadlsvnetrule%d"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_data_lake_store.test.name
  subnet_id           = azurerm_subnet.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, strconv.Itoa(data.RandomInteger)[2:13], data.RandomInteger)
}

/*
	(This test configuration is intended to succeed.)
	This test is designed to set up a scenario where a user would want to update the subnet
	on a given Data Lake Store virtual network rule. This configuration sets up the resources initially.
*/
func (r DataLakeStoreVirtualNetworkRuleResource) subnetSwitchPre(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test1" {
  name                 = "subnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.7.29.0/25"]
  service_endpoints    = ["Microsoft.AzureActiveDirectory"]
}

resource "azurerm_subnet" "test2" {
  name                 = "subnet2%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.7.29.128/25"]
  service_endpoints    = ["Microsoft.AzureActiveDirectory"]
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctestadls%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_data_lake_store_virtual_network_rule" "test" {
  name                = "acctestadlsvnetrule%d"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_data_lake_store.test.name
  subnet_id           = azurerm_subnet.test1.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, strconv.Itoa(data.RandomInteger)[2:13], data.RandomInteger)
}

/*
	(This test configuration is intended to succeed.)
	This test is designed to set up a scenario where a user would want to update the subnet
	on a given Data Lake Store virtual network rule. This configuration contains the update from
	azurerm_subnet.test1 to azurerm_subnet.test2.
*/
func (r DataLakeStoreVirtualNetworkRuleResource) subnetSwitchPost(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test1" {
  name                 = "subnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.7.29.0/25"]
  service_endpoints    = ["Microsoft.AzureActiveDirectory"]
}

resource "azurerm_subnet" "test2" {
  name                 = "subnet2%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.7.29.128/25"]
  service_endpoints    = ["Microsoft.AzureActiveDirectory"]
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctestadls%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_data_lake_store_virtual_network_rule" "test" {
  name                = "acctestadlsvnetrule%d"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_data_lake_store.test.name
  subnet_id           = azurerm_subnet.test2.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, strconv.Itoa(data.RandomInteger)[2:13], data.RandomInteger)
}

/*
	(This test configuration is intended to succeed.)
	This configuration sets up 3 subnets in 2 different virtual networks, and adds
	Data Lake Store virtual network rules for all 3 subnets to the Data Lake Store.
*/
func (r DataLakeStoreVirtualNetworkRuleResource) multipleSubnets(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "vnet1" {
  name                = "acctestvnet1%d"
  address_space       = ["10.7.29.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "vnet2" {
  name                = "acctestvnet2%d"
  address_space       = ["10.1.29.0/29"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "vnet1_subnet1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.vnet1.name
  address_prefixes     = ["10.7.29.0/29"]
  service_endpoints    = ["Microsoft.AzureActiveDirectory"]
}

resource "azurerm_subnet" "vnet1_subnet2" {
  name                 = "acctestsubnet2%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.vnet1.name
  address_prefixes     = ["10.7.29.128/29"]
  service_endpoints    = ["Microsoft.AzureActiveDirectory"]
}

resource "azurerm_subnet" "vnet2_subnet1" {
  name                 = "acctestsubnet3%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.vnet2.name
  address_prefixes     = ["10.1.29.0/29"]
  service_endpoints    = ["Microsoft.AzureActiveDirectory"]
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctestadls%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_data_lake_store_virtual_network_rule" "test" {
  name                = "acctestsqlvnetrule1%d"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_data_lake_store.test.name
  subnet_id           = azurerm_subnet.vnet1_subnet1.id
}

resource "azurerm_data_lake_store_virtual_network_rule" "rule2" {
  name                = "acctestsqlvnetrule2%d"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_data_lake_store.test.name
  subnet_id           = azurerm_subnet.vnet1_subnet2.id
}

resource "azurerm_data_lake_store_virtual_network_rule" "rule3" {
  name                = "acctestsqlvnetrule3%d"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_data_lake_store.test.name
  subnet_id           = azurerm_subnet.vnet2_subnet1.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, strconv.Itoa(data.RandomInteger)[2:13], data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
