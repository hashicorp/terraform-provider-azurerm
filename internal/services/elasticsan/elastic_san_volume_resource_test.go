// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package elasticsan_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ElasticSANVolumeTestResource struct{}

func TestAccElasticSANVolume_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_san_volume", "test")
	r := ElasticSANVolumeTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_iqn").IsNotEmpty(),
				check.That(data.ResourceName).Key("target_portal_hostname").IsNotEmpty(),
				check.That(data.ResourceName).Key("target_portal_port").IsNotEmpty(),
				check.That(data.ResourceName).Key("volume_id").IsNotEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccElasticSANVolume_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_san_volume", "test")
	r := ElasticSANVolumeTestResource{}

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

func TestAccElasticSANVolume_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_san_volume", "test")
	r := ElasticSANVolumeTestResource{}

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

func TestAccElasticSANVolume_fromDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_san_volume", "test")
	r := ElasticSANVolumeTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fromDisk(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccElasticSANVolume_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_san_volume", "test")
	r := ElasticSANVolumeTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config:      r.basic(data),
			ExpectError: regexp.MustCompile("new size_in_gib should be greater than the existing one"),
		},
	})
}

func (r ElasticSANVolumeTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := volumes.ParseVolumeID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ElasticSan.Volumes.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r ElasticSANVolumeTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_elastic_san_volume" "test" {
  name            = "acctestesv-${var.random_string}"
  volume_group_id = azurerm_elastic_san_volume_group.test.id
  size_in_gib     = 1
}
`, r.template(data))
}

func (r ElasticSANVolumeTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_elastic_san_volume" "import" {
  name            = azurerm_elastic_san_volume.test.name
  volume_group_id = azurerm_elastic_san_volume.test.volume_group_id
  size_in_gib     = azurerm_elastic_san_volume.test.size_in_gib
}
`, r.basic(data))
}

func (r ElasticSANVolumeTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_elastic_san_volume" "test" {
  name            = "acctestesv-${var.random_string}"
  volume_group_id = azurerm_elastic_san_volume_group.test.id
  size_in_gib     = 2
}
`, r.template(data))
}

func (r ElasticSANVolumeTestResource) fromDisk(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-${var.random_integer}"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  create_option        = "Empty"
  storage_account_type = "Standard_LRS"
  disk_size_gb         = 2
}

resource "azurerm_elastic_san_volume" "test" {
  name            = "acctestesv-${var.random_string}"
  size_in_gib     = 2
  volume_group_id = azurerm_elastic_san_volume_group.test.id
  create_source {
    source_type = "Disk"
    source_id   = azurerm_managed_disk.test.id
  }
}
`, r.template(data))
}

func (r ElasticSANVolumeTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-${var.random_integer}"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  create_option        = "Empty"
  storage_account_type = "Standard_LRS"
  disk_size_gb         = 2
}

resource "azurerm_snapshot" "test" {
  name                = "acctestss_${var.random_string}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Copy"
  source_uri          = azurerm_managed_disk.test.id
}

resource "azurerm_elastic_san_volume" "test" {
  name            = "acctestesv-${var.random_string}"
  size_in_gib     = 2
  volume_group_id = azurerm_elastic_san_volume_group.test.id
  create_source {
    source_type = "DiskSnapshot"
    source_id   = azurerm_snapshot.test.id
  }
}
`, r.template(data))
}

func (r ElasticSANVolumeTestResource) template(data acceptance.TestData) string {
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
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
