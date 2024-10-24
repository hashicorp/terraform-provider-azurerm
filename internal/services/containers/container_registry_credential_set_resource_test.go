// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/credentialsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerRegistryCredentialSetResource struct{}

func TestAccContainerRegistryCredentialSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_credential_set", "test")
	r := ContainerRegistryCredentialSetResource{}

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

func TestAccContainerRegistryCredentialSet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_credential_set", "test")
	r := ContainerRegistryCredentialSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_container_registry_credential_set"),
		},
	})
}

func TestAccContainerRegistryCredentialSet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_credential_set", "test")
	r := ContainerRegistryCredentialSetResource{}

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

func (ContainerRegistryCredentialSetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := credentialsets.ParseCredentialSetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Containers.CredentialSetsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (ContainerRegistryCredentialSetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "accTestRG-acr-credetial-set-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_container_registry_credential_set" "test" {
  name                  = "testacc-acr-credential-set-%d"
  container_registry_id = azurerm_container_registry.test.id
  login_server          = "docker.io"
  auth_credentials {
    username_secret_identifier = "https://example-keyvault.vault.azure.net/secrets/acr-cs-user-name"
    password_secret_identifier = "https://example-keyvault.vault.azure.net/secrets/acr-cs-user-password"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ContainerRegistryCredentialSetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_credential_set" "import" {
  name                  = azurerm_container_registry_credential_set.test.name
  container_registry_id = azurerm_container_registry_credential_set.test.container_registry_id
  login_server          = azurerm_container_registry_credential_set.test.login_server
  auth_credentials {
    username_secret_identifier = azurerm_container_registry_credential_set.test.auth_credentials[0].username_secret_identifier
    password_secret_identifier = azurerm_container_registry_credential_set.test.auth_credentials[0].password_secret_identifier
  }
}
`, r.basic(data))
}

func (ContainerRegistryCredentialSetResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "accTestRG-acr-credetial-set-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_container_registry_credential_set" "test" {
  name                  = "testacc-acr-credential-set-%d"
  container_registry_id = azurerm_container_registry.test.id
  login_server          = "docker.io"
  auth_credentials {
    username_secret_identifier = "https://example-keyvault.vault.azure.net/secrets/acr-cs-user-name-changed"
    password_secret_identifier = "https://example-keyvault.vault.azure.net/secrets/acr-cs-user-password"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
