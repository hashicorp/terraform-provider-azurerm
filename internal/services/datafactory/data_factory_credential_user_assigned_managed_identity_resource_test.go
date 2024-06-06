// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/credentials"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CredentialUserAssignedManagedIdentityResource struct{}

func TestAccDataFactoryCredentialUserAssignedManagedIdentity_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_credential_user_managed_identity", "test")
	r := CredentialUserAssignedManagedIdentityResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("annotations.#").HasValue("1"),
				check.That(data.ResourceName).Key("description").HasValue("ORIGINAL DESCRIPTION"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryCredentialUserAssignedManagedIdentity_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_credential_user_managed_identity", "test")
	r := CredentialUserAssignedManagedIdentityResource{}
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

func TestAccDataFactoryCredentialUserAssignedManagedIdentity_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_credential_user_managed_identity", "test")
	r := CredentialUserAssignedManagedIdentityResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("annotations.#").HasValue("1"),
				check.That(data.ResourceName).Key("description").HasValue("ORIGINAL DESCRIPTION"),
			),
		},
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("annotations.#").HasValue("2"),
				check.That(data.ResourceName).Key("description").HasValue("UPDATED DESCRIPTION"),
			),
		},
	})
}

func (t CredentialUserAssignedManagedIdentityResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := credentials.ParseCredentialID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataFactory.Credentials.CredentialOperationsGet(ctx, *id, credentials.DefaultCredentialOperationsGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func templateBase(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctestdf%d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CredentialUserAssignedManagedIdentityResource) basic(data acceptance.TestData) string {
	base := templateBase(data)

	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_credential_user_managed_identity" "test" {
  name            = "credential%d"
  description     = "ORIGINAL DESCRIPTION"
  data_factory_id = azurerm_data_factory.test.id
  identity_id     = azurerm_user_assigned_identity.test.id
  annotations     = ["1"]
}
`, base, data.RandomInteger)
}

func (r CredentialUserAssignedManagedIdentityResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_credential_user_managed_identity" "import" {
  name            = azurerm_data_factory_credential_user_managed_identity.test.name
  data_factory_id = azurerm_data_factory_credential_user_managed_identity.test.data_factory_id
  identity_id     = azurerm_data_factory_credential_user_managed_identity.test.identity_id
}
`, r.basic(data))
}

func (r CredentialUserAssignedManagedIdentityResource) update(data acceptance.TestData) string {
	base := templateBase(data)

	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_credential_user_managed_identity" "test" {
  name            = "credential%d"
  description     = "UPDATED DESCRIPTION"
  data_factory_id = azurerm_data_factory.test.id
  identity_id     = azurerm_user_assigned_identity.test.id
  annotations     = ["1", "2"]
}
`, base, data.RandomInteger)
}
