// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/identityprovider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementIdentityProviderAADResource struct{}

func TestAccApiManagementIdentityProviderAAD_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_aad", "test")
	r := ApiManagementIdentityProviderAADResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret"),
	})
}

func TestAccApiManagementIdentityProviderAAD_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_aad", "test")
	r := ApiManagementIdentityProviderAADResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_id").HasValue("00000000-0000-0000-0000-000000000000"),
				check.That(data.ResourceName).Key("client_secret").HasValue("00000000000000000000000000000000"),
				check.That(data.ResourceName).Key("allowed_tenants.#").HasValue("1"),
				check.That(data.ResourceName).Key("allowed_tenants.0").HasValue(data.Client().TenantID),
			),
		},
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_id").HasValue("11111111-1111-1111-1111-111111111111"),
				check.That(data.ResourceName).Key("client_secret").HasValue("11111111111111111111111111111111"),
				check.That(data.ResourceName).Key("allowed_tenants.#").HasValue("2"),
				check.That(data.ResourceName).Key("allowed_tenants.0").HasValue(data.Client().TenantID),
				check.That(data.ResourceName).Key("allowed_tenants.1").HasValue(data.Client().TenantID),
			),
		},
		data.ImportStep("client_secret"),
	})
}

func TestAccApiManagementIdentityProviderAAD_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_aad", "test")
	r := ApiManagementIdentityProviderAADResource{}

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

func (ApiManagementIdentityProviderAADResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := identityprovider.ParseIdentityProviderID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.IdentityProviderClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (ApiManagementIdentityProviderAADResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-api-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_identity_provider_aad" "test" {
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  client_id           = "00000000-0000-0000-0000-000000000000"
  client_secret       = "00000000000000000000000000000000"
  signin_tenant       = "00000000-0000-0000-0000-000000000000"
  allowed_tenants     = ["%s"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Client().TenantID)
}

func (ApiManagementIdentityProviderAADResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-api-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_identity_provider_aad" "test" {
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  client_id           = "11111111-1111-1111-1111-111111111111"
  client_secret       = "11111111111111111111111111111111"
  allowed_tenants     = ["%s", "%s"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Client().TenantID, data.Client().TenantID)
}

func (r ApiManagementIdentityProviderAADResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_identity_provider_aad" "import" {
  resource_group_name = azurerm_api_management_identity_provider_aad.test.resource_group_name
  api_management_name = azurerm_api_management_identity_provider_aad.test.api_management_name
  client_id           = azurerm_api_management_identity_provider_aad.test.client_id
  client_secret       = azurerm_api_management_identity_provider_aad.test.client_secret
  allowed_tenants     = azurerm_api_management_identity_provider_aad.test.allowed_tenants
}
`, r.basic(data))
}
