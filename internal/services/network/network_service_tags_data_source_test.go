// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type NetworkServiceTagsDataSource struct{}

func TestAccDataSourceAzureRMServiceTags_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_service_tags", "test")
	r := NetworkServiceTagsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("address_prefixes.#").Exists(),
				check.That(data.ResourceName).Key("ipv4_cidrs.#").Exists(),
				check.That(data.ResourceName).Key("ipv6_cidrs.#").Exists(),
			),
		},
	})
}

func TestAccDataSourceAzureRMServiceTags_tagName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_service_tags", "test")
	r := NetworkServiceTagsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.tagName(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue("westus2-Storage"),
				check.That(data.ResourceName).Key("name").HasValue("Storage.WestUS2"),
				check.That(data.ResourceName).Key("address_prefixes.#").Exists(),
				check.That(data.ResourceName).Key("ipv4_cidrs.#").Exists(),
				check.That(data.ResourceName).Key("ipv6_cidrs.#").Exists(),
			),
		},
	})
}

func TestAccDataSourceAzureRMServiceTags_region(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_service_tags", "test")
	r := NetworkServiceTagsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.region(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("address_prefixes.#").Exists(),
				check.That(data.ResourceName).Key("ipv4_cidrs.#").Exists(),
				check.That(data.ResourceName).Key("ipv6_cidrs.#").Exists(),
			),
		},
	})
}

func TestAccDataSourceAzureRMServiceTags_AzureFrontDoor(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_service_tags", "test")
	r := NetworkServiceTagsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.azureFrontDoor(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("address_prefixes.#").Exists(),
				check.That(data.ResourceName).Key("ipv4_cidrs.#").Exists(),
				check.That(data.ResourceName).Key("ipv6_cidrs.#").Exists(),
			),
		},
	})
}

func TestAccDataSourceAzureRMServiceTags_AzureFrontDoorBackend(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_service_tags", "test")
	r := NetworkServiceTagsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.azureFrontDoorBackend(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("address_prefixes.#").Exists(),
				check.That(data.ResourceName).Key("ipv4_cidrs.#").Exists(),
				check.That(data.ResourceName).Key("ipv6_cidrs.#").Exists(),
			),
		},
	})
}

func TestAccDataSourceAzureRMServiceTags_AzureFrontDoorFrontend(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_service_tags", "test")
	r := NetworkServiceTagsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.azureFrontDoorFrontend(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("address_prefixes.#").Exists(),
				check.That(data.ResourceName).Key("ipv4_cidrs.#").Exists(),
				check.That(data.ResourceName).Key("ipv6_cidrs.#").Exists(),
			),
		},
	})
}

func TestAccDataSourceAzureRMServiceTags_AzureFrontDoorFirstParty(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_service_tags", "test")
	r := NetworkServiceTagsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.azureFrontDoorFirstParty(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("address_prefixes.#").Exists(),
				check.That(data.ResourceName).Key("ipv4_cidrs.#").Exists(),
				check.That(data.ResourceName).Key("ipv6_cidrs.#").Exists(),
			),
		},
	})
}

func (NetworkServiceTagsDataSource) basic() string {
	return `data "azurerm_network_service_tags" "test" {
  location = "westcentralus"
  service  = "AzureKeyVault"
}`
}

func (NetworkServiceTagsDataSource) region() string {
	return `data "azurerm_network_service_tags" "test" {
  location        = "westcentralus"
  service         = "AzureKeyVault"
  location_filter = "australiacentral"
}`
}

func (NetworkServiceTagsDataSource) tagName() string {
	return `data "azurerm_network_service_tags" "test" {
  location        = "westus2"
  service         = "Storage"
  location_filter = "westus2"
}`
}

func (NetworkServiceTagsDataSource) azureFrontDoor() string {
	return `data "azurerm_network_service_tags" "test" {
  location = "northeurope"
  service  = "AzureFrontDoor"
}`
}

func (NetworkServiceTagsDataSource) azureFrontDoorBackend() string {
	return `data "azurerm_network_service_tags" "test" {
  location = "northeurope"
  service  = "AzureFrontDoor.Backend"
}`
}

func (NetworkServiceTagsDataSource) azureFrontDoorFrontend() string {
	return `data "azurerm_network_service_tags" "test" {
  location = "northeurope"
  service  = "AzureFrontDoor.Frontend"
}`
}

func (NetworkServiceTagsDataSource) azureFrontDoorFirstParty() string {
	return `data "azurerm_network_service_tags" "test" {
  location = "northeurope"
  service  = "AzureFrontDoor.FirstParty"
}`
}
