// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedidentity_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2023-01-31/managedidentities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FederatedIdentityCredentialTestResource struct{}

func TestAccFederatedIdentityCredential_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_federated_identity_credential", "test")
	r := FederatedIdentityCredentialTestResource{}

	rg := *regexp.MustCompile(`-updated`)

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
				check.That(data.ResourceName).Key("audience.0").MatchesRegex(&rg),
				check.That(data.ResourceName).Key("issuer").MatchesRegex(&rg),
				check.That(data.ResourceName).Key("subject").MatchesRegex(&rg),
			),
		},
	})
}

func TestAccFederatedIdentityCredential_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_federated_identity_credential", "test")
	r := FederatedIdentityCredentialTestResource{}

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

func (r FederatedIdentityCredentialTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := managedidentities.ParseFederatedIdentityCredentialID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ManagedIdentity.V20230131.ManagedIdentities.FederatedIdentityCredentialsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r FederatedIdentityCredentialTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_federated_identity_credential" "test" {
  audience            = ["foo"]
  issuer              = "https://foo"
  name                = "acctest-${local.random_integer}"
  resource_group_name = azurerm_resource_group.test.name
  parent_id           = azurerm_user_assigned_identity.test.id
  subject             = "foo"
}
`, r.template(data))
}

func (r FederatedIdentityCredentialTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_federated_identity_credential" "test" {
  audience            = ["foo-updated"]
  issuer              = "https://foo-updated"
  name                = "acctest-${local.random_integer}"
  resource_group_name = azurerm_resource_group.test.name
  parent_id           = azurerm_user_assigned_identity.test.id
  subject             = "foo-updated"
}
`, r.template(data))
}

func (r FederatedIdentityCredentialTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_federated_identity_credential" "import" {
  audience            = ["foo"]
  issuer              = "https://foo"
  name                = "acctest-${local.random_integer}"
  resource_group_name = azurerm_resource_group.test.name
  parent_id           = azurerm_user_assigned_identity.test.id
  subject             = "foo"
}
`, r.basic(data))
}

func (r FederatedIdentityCredentialTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
locals {
  random_integer   = %[1]d
  primary_location = %[2]q
}
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-${local.random_integer}"
  location = local.primary_location
}

resource "azurerm_user_assigned_identity" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctestuai-${local.random_integer}"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}
