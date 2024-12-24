package resourcegraph_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

// import (
// 	"context"
// 	"fmt"
// 	"testing"

// 	"github.com/hashicorp/go-azure-helpers/lang/response"
// )

type ResourceGraphQueryResource struct{}

func TestResourceGraphQuery_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_alert_v2", "test")
	r := ResourceGraphQueryResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ResourceGraphQueryResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_graph_query" "test" {
  name                = "acctestrgq-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  
  query = <<QUERY
		resources 
		| limit 1
	QUERY
}
`, r.template(data), data.RandomInteger)
}

func (r ResourceGraphQueryResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_resource_graph_query" "import" {
  name          		  = azurerm_resource_graph_query.test.name
  resource_group_name = azurerm_resource_graph_query.test.resource_group_name
  location 						= azurerm_resource_graph_query.test.location

	query         			= azurerm_resource_graph_query.test.query

  
}
`, r.basic(data))
}

func (r ResourceGraphQueryResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_graph_query" "test" {
  name                = "acctestrgq-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  
  query = <<QUERY
		resources
		| limit 2
	QUERY
`, r.template(data), data.RandomInteger)
}

func (r ResourceGraphQueryResource) template(data acceptance.TestData) string {
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
