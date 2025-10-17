package resourcegraph_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resourcegraph/2022-10-01/graphquery"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ResourceGraphQuery struct{}

func TestResourceGraphQuery_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_graph_query", "test")
	r := ResourceGraphQuery{}

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

func TestResourceGraphQuery_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_graph_query", "test")
	r := ResourceGraphQuery{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestResourceGraphQuery_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_graph_query", "test")
	r := ResourceGraphQuery{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
		data.ImportStep(),
	})
}

func TestResourceGraphQuery_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_example", "test")
	r := ResourceGraphQuery{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func (r ResourceGraphQuery) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_graph_query" "test" {
  name                = "acctest-rgq-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  query = <<QUERY
		resources 
		| limit 1
	QUERY
}
`, r.template(data), data.RandomInteger)
}

func (r ResourceGraphQuery) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_resource_graph_query" "import" {
  name                = azurerm_resource_graph_query.test.name
  resource_group_name = azurerm_resource_graph_query.test.resource_group_name
  location            = azurerm_resource_graph_query.test.location

  query = azurerm_resource_graph_query.test.query


}
`, r.basic(data))
}

func (r ResourceGraphQuery) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_graph_query" "test" {
  name                = "acctest-rgq-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  query = <<QUERY
		resources
		| limit 2
	QUERY
}
`, r.template(data), data.RandomInteger)
}

func (r ResourceGraphQuery) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-RGQ-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (t ResourceGraphQuery) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := graphquery.ParseQueryID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ResourceGraph.ResourceGraphQueryClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r ResourceGraphQuery) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_graph_query" "test" {
  name                = "acctest-rgq-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  display_name        = "acctest-rgq"
  description         = "test rgq"

  query = <<QUERY
		resources 
		| limit 1
	QUERY

  tags = {
    key = "value"
  }
}
`, r.template(data), data.RandomInteger)
}
