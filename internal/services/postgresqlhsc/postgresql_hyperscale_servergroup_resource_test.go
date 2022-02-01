package postgresqlhsc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgresqlhsc/sdk/2020-10-05-privatepreview/servergroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PostgreSQLHyperScaleServerGroupResource struct{}

func TestPostgreSQLHyperScaleServerGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_hyperscale_servergroup", "test")
	r := PostgreSQLHyperScaleServerGroupResource{}

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

// Exists func

func (r PostgreSQLHyperScaleServerGroupResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := servergroups.ParseServerGroupsv2ID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.PostgreSQLHSC.ServerGroupsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Server Group %s: %+v", id, err)
	}
	if response.WasNotFound(resp.HttpResponse) {
		return utils.Bool(false), nil
	}
	return utils.Bool(true), nil
}

// Configs

func (r PostgreSQLHyperScaleServerGroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_postgresql_hyperscale_servergroup" "test" {
  name                = "acctestcitus-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  administrator_login_password = "CMU7kVMB0cl2SsXk236iC0O71AoAqsSm"
  create_mode = "Default"
  citus_version = "8.3"
  postgresql_version = "11"
  server_role_group {
    role = "Coordinator"
    server_count = 1
    vcores = 16
    storage_quota_in_mb = 524288
  }
  server_role_group {
    role = "Worker"
    server_count = 1
    vcores = 16
    storage_quota_in_mb = 524288
  }

  tags = {
    Environment = "hyperscale"
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (PostgreSQLHyperScaleServerGroupResource) baseTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

`, data.RandomInteger, data.Locations.Primary)
}
