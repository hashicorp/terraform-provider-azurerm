package sql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type SqlManagedDatabase struct{}

func TestAccAzureRMSqlManagedDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_managed_database", "test")
	r := SqlManagedDatabase{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(""),
	})
}

func (r SqlManagedDatabase) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ManagedDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Sql.ManagedDatabasesClient.Get(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving SQL Managed Database %q: %+v", id.ID(), err)
	}
	return utils.Bool(true), nil
}

func (r SqlManagedDatabase) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_sql_managed_database" "test" {
  sql_managed_instance_id = azurerm_sql_managed_instance.test.id
  name                    = "acctest-%d"
  location                = azurerm_resource_group.test.location
}
`, SqlManagedInstanceResource{}.basic(data), data.RandomInteger)
}
