package postgresqlhsc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PostgreSQLHyperScaleCoordinatorConfigurationResource struct{}

func TestPostgreSQLHyperScaleCoordinatorConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_hyperscale_coordinator_configuration", "test")
	r := PostgreSQLHyperScaleCoordinatorConfigurationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "on"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestPostgreSQLHyperScaleCoordinatorConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_hyperscale_coordinator_configuration", "test")
	r := PostgreSQLHyperScaleCoordinatorConfigurationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "on"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(func(data acceptance.TestData) string {
			return r.requiresImport(data, "on")
		}),
	})
}

func TestPostgreSQLHyperScaleCoordinatorConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_hyperscale_coordinator_configuration", "test")
	r := PostgreSQLHyperScaleCoordinatorConfigurationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "on"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "off"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r PostgreSQLHyperScaleCoordinatorConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := configurations.ParseCoordinatorConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.PostgreSQLHSC.ConfigurationsClient
	resp, err := client.GetCoordinator(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r PostgreSQLHyperScaleCoordinatorConfigurationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-postgresqlhsc-%d"
  location = "%s"
}

resource "azurerm_postgresql_hyperscale_cluster" "test" {
  name                = "acctest-postgresqlhscsg-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r PostgreSQLHyperScaleCoordinatorConfigurationResource) basic(data acceptance.TestData, value string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_hyperscale_coordinator_configuration" "test" {
  name       = "acctest-postgresqlnc-%d"
  cluster_id = azurerm_postgresql_hyperscale_cluster.test.id
  value      = "%s"
}
`, r.template(data), data.RandomInteger, value)
}

func (r PostgreSQLHyperScaleCoordinatorConfigurationResource) requiresImport(data acceptance.TestData, value string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_hyperscale_coordinator_configuration" "import" {
  name       = azurerm_postgresql_hyperscale_coordinator_configuration.test.name
  cluster_id = azurerm_postgresql_hyperscale_coordinator_configuration.test.cluster_id
  value      = azurerm_postgresql_hyperscale_coordinator_configuration.test.value
}
`, r.basic(data, value))
}
