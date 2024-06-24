package elasticsan_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ElasticSANVolumeSnapshotDataSource struct{}

// https://github.com/hashicorp/terraform-provider-azurerm/pull/25372#issuecomment-2022105240
// Elastic SAN Volume Snapshot is context-based and should not be regarded as the infrastructure managed by Terraform
// so we only onboard this as a data source instead of a resource. The acctest relys on Azure CLI to create/delete the snapshot.
func TestAccElasticSANVolumeSnapshotDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_elastic_san_volume_snapshot", "test")
	d := ElasticSANVolumeSnapshotDataSource{}

	const SnapshotTestRunEnv = "ARM_TEST_ELASTIC_SAN_VOLUME_SNAPSHOT_RUN"

	if os.Getenv(SnapshotTestRunEnv) == "" {
		t.Skipf("skip the test as one or more of below environment variables are not specified: %q", SnapshotTestRunEnv)
	}

	data.DataSourceTestInSequence(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("creation_source.#").HasValue("1"),
				check.That(data.ResourceName).Key("creation_source.0.source_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("source_volume_size_gib").IsNotEmpty(),
				check.That(data.ResourceName).Key("volume_name").IsNotEmpty(),
			),
		},
	})
}

func (d ElasticSANVolumeSnapshotDataSource) basic(data acceptance.TestData) string {
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

provider "azurerm" {
  features {}
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

locals {
  snapshot_name = "acctest-ess-${var.random_string}"
}

resource "terraform_data" "test" {
  input = {
    resource_group_name = azurerm_resource_group.test.name
    elastic_san_name    = azurerm_elastic_san.test.name
    volume_group_name   = azurerm_elastic_san_volume_group.test.name
    snapshot_name       = local.snapshot_name
  }

  provisioner "local-exec" {
    command = <<COMMAND
      az elastic-san volume snapshot create -g ${self.input.resource_group_name} -e ${self.input.elastic_san_name} -v ${self.input.volume_group_name} -n ${self.input.snapshot_name} --creation-data '{source-id:${azurerm_elastic_san_volume.test.id}}'
    COMMAND
  }

  provisioner "local-exec" {
    when    = destroy
    command = <<COMMAND
      az elastic-san volume snapshot delete -g ${self.input.resource_group_name} -e ${self.input.elastic_san_name} -v ${self.input.volume_group_name} -n ${self.input.snapshot_name} -y
    COMMAND
  }
}

data "azurerm_elastic_san_volume_snapshot" "test" {
  name            = local.snapshot_name
  volume_group_id = azurerm_elastic_san_volume_group.test.id
  depends_on      = [terraform_data.test]
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
