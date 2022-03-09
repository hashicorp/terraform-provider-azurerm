package mssql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlManagedDatabase struct{}

func TestAccMsSqlManagedDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_database", "test")
	r := MsSqlManagedDatabase{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(""),
	})
}

func (r MsSqlManagedDatabase) Exists(ctx context.Context, client *clients.Client, state *acceptance.InstanceState) (*bool, error) {
	id, err := parse.ManagedDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.ManagedDatabasesClient.Get(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving SQL Managed Database %q: %+v", id.ID(), err)
	}

	return utils.Bool(true), nil
}

func (r MsSqlManagedDatabase) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_managed_database" "test" {
  managed_instance_id = azurerm_mssql_managed_instance.test.id
  name                = "acctest-%[2]d"
}
`, MsSqlManagedInstanceResource{}.basic(data), data.RandomInteger)
}
