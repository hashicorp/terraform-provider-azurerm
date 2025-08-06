// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package elasticsan_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumes"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

type ElasticSANVolumeSnapshotDataSource struct{}

// https://github.com/hashicorp/terraform-provider-azurerm/pull/25372#issuecomment-2022105240
// Elastic SAN Volume Snapshot is context-based and should not be regarded as the infrastructure managed by Terraform
// so we only onboard this as a data source instead of a resource. The acctest creates the snapshot as a test step
func TestAccElasticSANVolumeSnapshotDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_elastic_san_volume_snapshot", "test")
	d := ElasticSANVolumeSnapshotDataSource{}

	data.DataSourceTestInSequence(t, []acceptance.TestStep{
		{
			Config: d.snapshotSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
					if _, ok := ctx.Deadline(); !ok {
						var cancel context.CancelFunc
						ctx, cancel = context.WithTimeout(ctx, 30*time.Minute)
						defer cancel()
					}

					volumeId, err := volumes.ParseVolumeID(state.ID)
					if err != nil {
						return err
					}

					id := snapshots.NewSnapshotID(volumeId.SubscriptionId, volumeId.ResourceGroupName, volumeId.ElasticSanName, volumeId.VolumeGroupName, data.RandomString)

					snapshot := snapshots.Snapshot{
						Properties: snapshots.SnapshotProperties{
							CreationData: snapshots.SnapshotCreationData{
								SourceId: volumeId.ID(),
							},
						},
					}

					client := clients.ElasticSan.Snapshots
					if err = client.VolumeSnapshotsCreateThenPoll(ctx, id, snapshot); err != nil {
						return fmt.Errorf("creating %s: %+v", id, err)
					}

					return nil
				}, "azurerm_elastic_san_volume.test"),
			),
		},
		{
			Config: d.snapshotRestore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("source_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("source_volume_size_in_gib").IsNotEmpty(),
				check.That(data.ResourceName).Key("volume_name").IsNotEmpty(),
			),
		},
		{
			Config: d.snapshotSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
					if _, ok := ctx.Deadline(); !ok {
						var cancel context.CancelFunc
						ctx, cancel = context.WithTimeout(ctx, 30*time.Minute)
						defer cancel()
					}

					volumeId, err := volumes.ParseVolumeID(state.ID)
					if err != nil {
						return err
					}

					id := snapshots.NewSnapshotID(volumeId.SubscriptionId, volumeId.ResourceGroupName, volumeId.ElasticSanName, volumeId.VolumeGroupName, data.RandomString)

					client := clients.ElasticSan.Snapshots
					if err = client.VolumeSnapshotsDeleteThenPoll(ctx, id); err != nil {
						return fmt.Errorf("creating %s: %+v", id, err)
					}

					return nil
				}, "azurerm_elastic_san_volume.test"),
			),
		},
	})
}

func (d ElasticSANVolumeSnapshotDataSource) snapshotSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-esvg-%[2]d"
  location = "%[1]s"
}

resource "azurerm_elastic_san" "test" {
  name                = "acctestes-%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  base_size_in_tib    = 1
  sku {
    name = "Premium_LRS"
  }
}

resource "azurerm_elastic_san_volume_group" "test" {
  name           = "acctestesvg-%[3]s"
  elastic_san_id = azurerm_elastic_san.test.id
}

resource "azurerm_elastic_san_volume" "test" {
  name            = "acctestesv-%[3]s"
  volume_group_id = azurerm_elastic_san_volume_group.test.id
  size_in_gib     = 1
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (d ElasticSANVolumeSnapshotDataSource) snapshotRestore(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-esvg-%[2]d"
  location = "%[1]s"
}

resource "azurerm_elastic_san" "test" {
  name                = "acctestes-%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  base_size_in_tib    = 1
  sku {
    name = "Premium_LRS"
  }
}

resource "azurerm_elastic_san_volume_group" "test" {
  name           = "acctestesvg-%[3]s"
  elastic_san_id = azurerm_elastic_san.test.id
}

resource "azurerm_elastic_san_volume" "test" {
  name            = "acctestesv-%[3]s"
  volume_group_id = azurerm_elastic_san_volume_group.test.id
  size_in_gib     = 1
}

data "azurerm_elastic_san_volume_snapshot" "test" {
  name            = "%[3]s"
  volume_group_id = azurerm_elastic_san_volume_group.test.id
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
