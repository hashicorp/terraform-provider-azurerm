package postgresqlhsc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PostgreSQLHyperScaleClusterResource struct{}

func TestAccPostgreSQLHyperScaleCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_hyperscale_cluster", "test")
	r := PostgreSQLHyperScaleClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccPostgreSQLHyperScaleCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_hyperscale_cluster", "test")
	r := PostgreSQLHyperScaleClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccPostgreSQLHyperScaleCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_hyperscale_cluster", "test")
	r := PostgreSQLHyperScaleClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccPostgreSQLHyperScaleCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_hyperscale_cluster", "test")
	r := PostgreSQLHyperScaleClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func (r PostgreSQLHyperScaleClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := clusters.ParseServerGroupsv2ID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.PostgreSQLHSC.ClustersClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r PostgreSQLHyperScaleClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-pshsc-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r PostgreSQLHyperScaleClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_hyperscale_cluster" "test" {
  name                            = "acctestcluster%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  administrator_login_password    = "H@Sh1CoR3!"
  coordinator_storage_quota_in_mb = 131072
  coordinator_vcores              = 2
  node_count                      = 0
}
`, r.template(data), data.RandomInteger)
}

func (r PostgreSQLHyperScaleClusterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_hyperscale_cluster" "import" {
  name                            = azurerm_postgresql_hyperscale_cluster.test.name
  resource_group_name             = azurerm_postgresql_hyperscale_cluster.test.resource_group_name
  location                        = azurerm_postgresql_hyperscale_cluster.test.location
  administrator_login_password    = azurerm_postgresql_hyperscale_cluster.test.administrator_login_password
  coordinator_storage_quota_in_mb = azurerm_postgresql_hyperscale_cluster.test.coordinator_storage_quota_in_mb
  coordinator_vcores              = azurerm_postgresql_hyperscale_cluster.test.coordinator_vcores
  node_count                      = azurerm_postgresql_hyperscale_cluster.test.node_count
}
`, r.basic(data))
}

func (r PostgreSQLHyperScaleClusterResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_hyperscale_cluster" "test" {
  name                            = "acctestcluster%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  administrator_login_password    = "H@Sh1CoR3!"
  coordinator_storage_quota_in_mb = 131072
  coordinator_vcores              = 2
  node_count                      = 0
}
`, r.template(data), data.RandomInteger)
}

func (r PostgreSQLHyperScaleClusterResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_hyperscale_cluster" "test" {
  name                            = "acctestcluster%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  administrator_login_password    = "H@Sh1CoR4!"
  coordinator_storage_quota_in_mb = 262144
  coordinator_vcores              = 2
  node_count                      = 0
}
`, r.template(data), data.RandomInteger)
}
