package managementgroup_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ManagementGroupResource struct {
}

func TestAcc_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group", "test")
	r := ManagementGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagementGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group", "test")
	r := ManagementGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(),
			ExpectError: acceptance.RequiresImportError("azurerm_management_group"),
		},
	})
}

func TestAccManagementGroup_nested(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group", "parent")
	r := ManagementGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.nested(),
			Check: resource.ComposeTestCheckFunc(
				check.That("azurerm_management_group.parent").ExistsInAzure(r),
				check.That("azurerm_management_group.child").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		data.ImportStepFor("azurerm_management_group.child"),
	})
}

func TestAccManagementGroup_multiLevel(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group", "parent")
	r := ManagementGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multiLevel(),
			Check: resource.ComposeTestCheckFunc(
				check.That("azurerm_management_group.grandparent").ExistsInAzure(r),
				check.That("azurerm_management_group.parent").ExistsInAzure(r),
				check.That("azurerm_management_group.child").ExistsInAzure(r),
			),
		},
		data.ImportStepFor("azurerm_management_group.grandparent"),
		data.ImportStep(),
		data.ImportStepFor("azurerm_management_group.child"),
	})
}

func TestAccManagementGroup_multiLevelUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group", "parent")
	r := ManagementGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.nested(),
			Check: resource.ComposeTestCheckFunc(
				check.That("azurerm_management_group.parent").ExistsInAzure(r),
				check.That("azurerm_management_group.child").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		data.ImportStepFor("azurerm_management_group.child"),
		{
			Config: r.multiLevel(),
			Check: resource.ComposeTestCheckFunc(
				check.That("azurerm_management_group.grandparent").ExistsInAzure(r),
				check.That("azurerm_management_group.parent").ExistsInAzure(r),
				check.That("azurerm_management_group.child").ExistsInAzure(r),
			),
		},
		data.ImportStepFor("azurerm_management_group.grandparent"),
		data.ImportStep(),
		data.ImportStepFor("azurerm_management_group.child"),
	})
}

func TestAccManagementGroup_withName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group", "test")
	r := ManagementGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withName(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagementGroup_updateName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group", "test")
	r := ManagementGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.withName(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestmg-%d", data.RandomInteger)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagementGroup_withSubscriptions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group", "test")
	r := ManagementGroupResource{}
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subscription_ids.#").HasValue("0"),
			),
		},
		{
			Config: r.withSubscriptions(subscriptionID),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subscription_ids.#").HasValue("1"),
			),
		},
		{
			Config: r.basic(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subscription_ids.#").HasValue("0"),
			),
		},
	})
}

func (ManagementGroupResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ManagementGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ManagementGroups.GroupsClient.Get(ctx, id.Name, "children", utils.Bool(true), "", "no-cache")
	if err != nil {
		return nil, fmt.Errorf("retrieving Management Group %s: %v", id.Name, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (r ManagementGroupResource) basic() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
}
`
}

func (r ManagementGroupResource) requiresImport() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
}

resource "azurerm_management_group" "import" {
  name = azurerm_management_group.test.name
}
`
}

func (r ManagementGroupResource) nested() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "parent" {
}

resource "azurerm_management_group" "child" {
  parent_management_group_id = azurerm_management_group.parent.id
}
`
}

func (r ManagementGroupResource) multiLevel() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "grandparent" {
}

resource "azurerm_management_group" "parent" {
  parent_management_group_id = azurerm_management_group.grandparent.id
}

resource "azurerm_management_group" "child" {
  parent_management_group_id = azurerm_management_group.parent.id
}
`
}

func (ManagementGroupResource) withName(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  name         = "acctestmg-%d"
  display_name = "accTestMG-%d"
}
`, data.RandomInteger, data.RandomInteger)
}

// TODO: switch this out for dynamically creating a subscription once that's supported in the future
func (r ManagementGroupResource) withSubscriptions(subscriptionID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  subscription_ids = [
    "%s",
  ]
}
`, subscriptionID)
}
