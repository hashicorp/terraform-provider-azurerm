// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package customproviders_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/customproviders/2018-09-01-preview/customresourceprovider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CustomProviderResource struct{}

func TestAccCustomProvider_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_provider", "test")
	r := CustomProviderResource{}

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

func TestAccCustomProvider_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_provider", "test")
	r := CustomProviderResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
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

func TestAccCustomProvider_action(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_provider", "test")
	r := CustomProviderResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.action(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.actionUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.action(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CustomProviderResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := customresourceprovider.ParseResourceProviderID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.CustomProviders.CustomProviderClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r CustomProviderResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cp-%d"
  location = "%s"
}
resource "azurerm_custom_provider" "test" {
  name                = "accTEst_saa%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  resource_type {
    name     = "dEf1"
    endpoint = "https://example.com/"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CustomProviderResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cp-%d"
  location = "%s"
}
resource "azurerm_custom_provider" "test" {
  name                = "accTEst_saa%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  resource_type {
    name     = "dEf1"
    endpoint = "https://example.com/"
  }

  action {
    name     = "dEf2"
    endpoint = "https://example.com/"
  }

  validation {
    specification = "https://raw.githubusercontent.com/Azure/azure-custom-providers/master/CustomRPWithSwagger/Artifacts/Swagger/pingaction.json"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CustomProviderResource) action(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cp-%d"
  location = "%s"
}
resource "azurerm_custom_provider" "test" {
  name                = "accTEst_saa%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  action {
    name     = "dEf1"
    endpoint = "https://example.com/"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CustomProviderResource) actionUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cp-%d"
  location = "%s"
}
resource "azurerm_custom_provider" "test" {
  name                = "accTEst_saa%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  action {
    name     = "dEf2"
    endpoint = "https://example.com/"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
