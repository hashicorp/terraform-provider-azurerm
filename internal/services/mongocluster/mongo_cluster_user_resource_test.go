// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongocluster_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2025-09-01/users"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MongoClusterUserResource struct{}

func TestAccMongoClusterUser_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mongo_cluster_user", "test")
	r := MongoClusterUserResource{}

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

func TestAccMongoClusterUser_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mongo_cluster_user", "test")
	r := MongoClusterUserResource{}

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

func TestAccMongoClusterUser_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mongo_cluster_user", "test")
	r := MongoClusterUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r MongoClusterUserResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := users.ParseUserID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MongoCluster.UsersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r MongoClusterUserResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mongo_cluster_user" "test" {
  object_id              = data.azurerm_client_config.current.object_id
  mongo_cluster_id       = azurerm_mongo_cluster.test.id
  identity_provider_type = "MicrosoftEntraID"
  principal_type         = "servicePrincipal"

  role {
    database = "admin"
    role     = "root"
  }
}
`, r.template(data))
}

func (r MongoClusterUserResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mongo_cluster_user" "import" {
  object_id              = azurerm_mongo_cluster_user.test.object_id
  mongo_cluster_id       = azurerm_mongo_cluster_user.test.mongo_cluster_id
  identity_provider_type = azurerm_mongo_cluster_user.test.identity_provider_type
  principal_type         = azurerm_mongo_cluster_user.test.principal_type

  role {
    database = azurerm_mongo_cluster_user.test.role.0.database
    role     = azurerm_mongo_cluster_user.test.role.0.role
  }
}
`, r.basic(data))
}

func (r MongoClusterUserResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azuread_domains" "test" {
  only_initial = true
}

resource "azuread_user" "test" {
  user_principal_name = "acctestAadUser-%[2]d@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestAadUser-%[2]d"
  password            = "TerrAform321!"
}

resource "azurerm_mongo_cluster_user" "test" {
  object_id              = azuread_user.test.object_id
  mongo_cluster_id       = azurerm_mongo_cluster.test.id
  identity_provider_type = "MicrosoftEntraID"
  principal_type         = "user"

  role {
    database = "admin"
    role     = "root"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MongoClusterUserResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mc-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mongo_cluster" "test" {
  name                      = "acctest-mc%[1]d"
  resource_group_name       = azurerm_resource_group.test.name
  location                  = azurerm_resource_group.test.location
  administrator_username    = "adminTerraform"
  administrator_password    = "QAZwsx123"
  shard_count               = "1"
  compute_tier              = "M30"
  high_availability_mode    = "Disabled"
  storage_size_in_gb        = "32"
  version                   = "8.0"
  auth_config_allowed_modes = ["NativeAuth", "MicrosoftEntraID"]
}
`, data.RandomInteger, data.Locations.Primary)
}
