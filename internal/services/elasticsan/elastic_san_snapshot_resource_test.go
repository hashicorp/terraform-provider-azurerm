package elasticsan_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/snapshots"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ElasticSANSnapshotTestResource struct{}

func TestAccElasticSANSnapshot_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_san_snapshot", "test")
	r := ElasticSANSnapshotTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_volume_size_gib").HasValue("1"),
				check.That(data.ResourceName).Key("volume_name").IsNotEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccElasticSANSnapshot_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_san_snapshot", "test")
	r := ElasticSANSnapshotTestResource{}

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

func (r ElasticSANSnapshotTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := snapshots.ParseSnapshotID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.ElasticSan.Snapshots
	resp, err := client.VolumeSnapshotsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r ElasticSANSnapshotTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_elastic_san_snapshot" "test" {
  name            = "acctestess-${var.random_string}"
  volume_group_id = azurerm_elastic_san_volume_group.test.id
  creation_source {
    source_id = azurerm_elastic_san_volume.test.id
  }
}
`, r.template(data))
}

func (r ElasticSANSnapshotTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_elastic_san_snapshot" "import" {
  name            = azurerm_elastic_san_snapshot.test.name
  volume_group_id = azurerm_elastic_san_snapshot.test.volume_group_id
  creation_source {
    source_id = azurerm_elastic_san_snapshot.test.creation_source.0.source_id
  }
}
`, r.basic(data))
}

func (r ElasticSANSnapshotTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "primary_location" {
  default = %q
}
variable "random_integer" {
  default = %d
}
variable "random_string" {
  default = %q
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-esvg-${var.random_integer}"
  location = var.primary_location
}

resource "azurerm_elastic_san" "test" {
  name                = "acctestes-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  base_size_in_tib    = 1
  sku {
    name = "Premium_LRS"
  }
}

resource "azurerm_elastic_san_volume_group" "test" {
  name           = "acctestesvg-${var.random_string}"
  elastic_san_id = azurerm_elastic_san.test.id
}

resource "azurerm_elastic_san_volume" "test" {
  name            = "acctestesv-${var.random_string}"
  volume_group_id = azurerm_elastic_san_volume_group.test.id
  size_in_gib     = 1
}


`, data.Locations.Primary, data.RandomInteger, data.RandomString)

}
