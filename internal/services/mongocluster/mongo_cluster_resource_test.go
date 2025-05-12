// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongocluster_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2024-07-01/mongoclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MongoClusterResource struct{}

func TestAccMongoClusterFreeTier(t *testing.T) {
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"freeTier": { // Run tests in sequence since each subscription is limited to one free tier cluster per region and free tier is currently only available in South India.
			"basic":  testAccMongoCluster_basic,
			"update": testAccMongoCluster_update,
			"import": testAccMongoCluster_requiresImport,
		},
	})
}

func testAccMongoCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mongo_cluster", "test")
	r := MongoClusterResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("connection_strings.0.value").HasValue(
					fmt.Sprintf(`mongodb+srv://adminTerraform:QAZwsx123basic@acctest-mc%d.mongocluster.cosmos.azure.com/?tls=true&authMechanism=SCRAM-SHA-256&retrywrites=false&maxIdleTimeMS=120000`,
						data.RandomInteger)),
			),
		},
		data.ImportStep("administrator_password", "create_mode", "connection_strings.0.value"),
	})
}

func testAccMongoCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mongo_cluster", "test")
	r := MongoClusterResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode", "connection_strings.0.value"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode", "connection_strings.0.value"),
	})
}

func testAccMongoCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mongo_cluster", "test")
	r := MongoClusterResource{}

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

func TestAccMongoCluster_previewFeature(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mongo_cluster", "test")
	r := MongoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.previewFeature(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode", "connection_strings.0.value", "connection_strings.1.value"),
		{
			Config: r.geoReplica(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode", "source_location", "connection_strings.0.value", "connection_strings.1.value"),
	})
}

func TestAccMongoCluster_geoReplica(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mongo_cluster", "test")
	r := MongoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.geoReplica(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode", "source_location", "connection_strings.0.value", "connection_strings.1.value"),
	})
}

func (r MongoClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := mongoclusters.ParseMongoClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MongoCluster.MongoClustersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r MongoClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mongo_cluster" "test" {
  name                   = "acctest-mc%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_username = "adminTerraform"
  administrator_password = "QAZwsx123basic"
  shard_count            = "1"
  compute_tier           = "Free"
  high_availability_mode = "Disabled"
  storage_size_in_gb     = "32"
  version                = "6.0"
}
`, r.template(data, data.Locations.Ternary), data.RandomInteger)
}

func (r MongoClusterResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mongo_cluster" "test" {
  name                   = "acctest-mc%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_username = "adminTerraform"
  administrator_password = "QAZwsx123update"
  shard_count            = "1"
  compute_tier           = "M30"
  high_availability_mode = "ZoneRedundantPreferred"
  public_network_access  = "Disabled"
  storage_size_in_gb     = "64"
  version                = "7.0"

  tags = {
    environment = "test"
  }
}
`, r.template(data, data.Locations.Ternary), data.RandomInteger)
}

func (r MongoClusterResource) source(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mongo_cluster" "test" {
  name                   = "acctest-mc%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_username = "adminTerraform"
  administrator_password = "QAZwsx123update"
  high_availability_mode = "ZoneRedundantPreferred"
  shard_count            = "1"
  compute_tier           = "M30"
  storage_size_in_gb     = "64"
  version                = "7.0"
}
`, r.template(data, data.Locations.Ternary), data.RandomInteger)
}

func (r MongoClusterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mongo_cluster" "import" {
  name                   = azurerm_mongo_cluster.test.name
  resource_group_name    = azurerm_mongo_cluster.test.resource_group_name
  location               = azurerm_mongo_cluster.test.location
  administrator_username = azurerm_mongo_cluster.test.administrator_username
  administrator_password = azurerm_mongo_cluster.test.administrator_password
  shard_count            = azurerm_mongo_cluster.test.shard_count
  compute_tier           = azurerm_mongo_cluster.test.compute_tier
  high_availability_mode = azurerm_mongo_cluster.test.high_availability_mode
  storage_size_in_gb     = azurerm_mongo_cluster.test.storage_size_in_gb
  version                = azurerm_mongo_cluster.test.version
}
`, r.basic(data))
}

func (r MongoClusterResource) previewFeature(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mongo_cluster" "test" {
  name                   = "acctest-mc%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_username = "adminTerraform"
  administrator_password = "testQAZwsx123"
  shard_count            = "1"
  compute_tier           = "M30"
  high_availability_mode = "ZoneRedundantPreferred"
  storage_size_in_gb     = "64"
  preview_features       = ["GeoReplicas"]
  version                = "7.0"
}
`, r.template(data, data.Locations.Primary), data.RandomInteger)
}

func (r MongoClusterResource) geoReplica(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mongo_cluster" "geo_replica" {
  name                = "acctest-mc-replica%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  source_server_id    = azurerm_mongo_cluster.test.id
  source_location     = azurerm_mongo_cluster.test.location
  create_mode         = "GeoReplica"

  lifecycle {
    ignore_changes = ["administrator_username", "high_availability_mode", "preview_features", "shard_count", "storage_size_in_gb", "compute_tier", "version"]
  }
}
`, r.source(data), data.RandomInteger, data.Locations.Secondary)
}

func (r MongoClusterResource) template(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
`, data.RandomInteger, location)
}
