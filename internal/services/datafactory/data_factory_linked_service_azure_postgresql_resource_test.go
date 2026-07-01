// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LinkedServiceAzurePostgreSQLResource struct{}

func TestAccDataFactoryLinkedServiceAzurePostgreSQL_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_postgresql", "test")
	r := LinkedServiceAzurePostgreSQLResource{}

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

func (t LinkedServiceAzurePostgreSQLResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := linkedservices.ParseLinkedServiceID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := clients.DataFactory.LinkedServicesClient.Get(ctx, *id, linkedservices.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("reading Data Factory Linked Service PostgreSQL (%s): %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (LinkedServiceAzurePostgreSQLResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r LinkedServiceAzurePostgreSQLResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_data_factory_linked_service_azure_postgresql" "test" {
  name                = "acctestadfpostgresql%d"
  data_factory_id     = azurerm_data_factory.test.id
  authentication_type = "SystemAssignedManagedIdentity"
  server              = "acctest-server.postgres.database.azure.com"
  port                = 5432
  database_name       = "acctestdb"
}
`, r.template(data), data.RandomInteger)
}
