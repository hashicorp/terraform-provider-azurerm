// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package elasticsan_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/elasticsans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ElasticSANTestResource struct{}

func TestAccElasticSAN_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_san", "test")
	r := ElasticSANTestResource{}

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

func TestAccElasticSAN_zoneWithZRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_san", "test")
	r := ElasticSANTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.zoneWithZRS(data),
			ExpectError: regexp.MustCompile("zones are not supported"),
		},
	})
}

func TestAccElasticSAN_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_san", "test")
	r := ElasticSANTestResource{}

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

func TestAccElasticSAN_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_san", "test")
	r := ElasticSANTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("total_iops").Exists(),
				check.That(data.ResourceName).Key("total_mbps").Exists(),
				check.That(data.ResourceName).Key("total_size_in_tib").Exists(),
				check.That(data.ResourceName).Key("total_volume_size_in_gib").Exists(),
				check.That(data.ResourceName).Key("volume_group_count").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccElasticSAN_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_san", "test")
	r := ElasticSANTestResource{}

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
	})
}

func TestAccElasticSAN_updateTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_san", "test")
	r := ElasticSANTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccElasticSAN_reduceBaseSize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_san", "test")
	r := ElasticSANTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config:      r.basic(data),
			ExpectError: regexp.MustCompile("new base_size_in_tib should be greater than the existing one"),
		},
		data.ImportStep(),
	})
}

func (r ElasticSANTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := elasticsans.ParseElasticSanID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ElasticSan.ElasticSans.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r ElasticSANTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
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
`, r.template(data))
}

func (r ElasticSANTestResource) basicWithTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_elastic_san" "test" {
  name                = "acctestes-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  base_size_in_tib    = 1
  sku {
    name = "Premium_LRS"
  }
  tags = {
    foo = "bar"
  }
}
`, r.template(data))
}

func (r ElasticSANTestResource) zoneWithZRS(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_elastic_san" "test" {
  name                 = "acctestes-${var.random_string}"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  base_size_in_tib     = 1
  extended_size_in_tib = 1
  zones                = ["1"]
  sku {
    name = "Premium_ZRS"
  }
}
`, r.template(data))
}

func (r ElasticSANTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_elastic_san" "import" {
  name                = azurerm_elastic_san.test.name
  resource_group_name = azurerm_elastic_san.test.resource_group_name
  location            = azurerm_elastic_san.test.location
  base_size_in_tib    = azurerm_elastic_san.test.base_size_in_tib
  zones               = azurerm_elastic_san.test.zones
  sku {
    name = azurerm_elastic_san.test.sku.0.name
  }
}
`, r.basic(data))
}

func (r ElasticSANTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_elastic_san" "test" {
  name                 = "acctestes-${var.random_string}"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  base_size_in_tib     = 2
  extended_size_in_tib = 4
  sku {
    name = "Premium_LRS"
    tier = "Premium"
  }
  tags = {
    some_key = "some-value"
  }
}
`, r.template(data))
}

func (r ElasticSANTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_elastic_san" "test" {
  name                 = "acctestes-${var.random_string}"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  base_size_in_tib     = 2
  extended_size_in_tib = 4
  zones                = ["1", "2"]
  sku {
    name = "Premium_LRS"
    tier = "Premium"
  }
  tags = {
    environment = "terraform-acctests"
    some_key    = "some-value"
  }
}
`, r.template(data))
}

func (r ElasticSANTestResource) template(data acceptance.TestData) string {
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
  name     = "acctestrg-${var.random_integer}"
  location = var.primary_location
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
