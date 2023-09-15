package graphservices_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/graphservices/2023-04-13/graphservicesprods"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func TestAccGraphAccount(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to
	// the account need a pre-existing AD application, here we use the service principal.
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		// remove in 4.0
		"legacyAccount": {
			"basic": testAccGraphAccount_legacy,
		},
		"account": {
			"basic":          testAccGraphAccount_basic,
			"update":         testAccGraphAccount_update,
			"complete":       testAccGraphAccount_complete,
			"requiresImport": testAccGraphAccount_requiresImport,
		},
	})
}

type AccountTestResource struct{}

func testAccGraphAccount_legacy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_graph_account", "test")
	r := AccountTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.legacy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccGraphAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_graph_services_account", "test")
	r := AccountTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccGraphAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_graph_services_account", "test")

	r := AccountTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccGraphAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_graph_services_account", "test")
	r := AccountTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccGraphAccount_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_graph_services_account", "test")
	r := AccountTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
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
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}
func (r AccountTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := graphservicesprods.ParseAccountID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Graph.V20230413.Graphservicesprods.AccountsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r AccountTestResource) legacy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_graph_account" "test" {
  name                = "acctesta-%[2]d"
  application_id      = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data), data.RandomInteger, os.Getenv("ARM_CLIENT_ID"))
}

func (r AccountTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_graph_services_account" "test" {
  name                = "acctesta-%[2]d"
  application_id      = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data), data.RandomInteger, os.Getenv("ARM_CLIENT_ID"))
}

func (r AccountTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_graph_services_account" "import" {
  application_id      = azurerm_graph_services_account.test.application_id
  name                = azurerm_graph_services_account.test.name
  resource_group_name = azurerm_graph_services_account.test.resource_group_name
}
`, r.basic(data))
}

func (r AccountTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_graph_services_account" "test" {
  name                = "acctesta-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  application_id      = "%[3]s"
  tags = {
    environment = "terraform-acctests"
    some_key    = "some-value"
  }
}
`, r.template(data), data.RandomInteger, os.Getenv("ARM_CLIENT_ID"))
}

func (r AccountTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-graph-%d"
  location = %q
}

resource "azuread_application" "test" {
  display_name = "acctestsap%[1]d"
}
`, data.RandomInteger, data.Locations.Primary)
}
