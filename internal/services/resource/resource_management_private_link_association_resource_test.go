// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/privatelinkassociation"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ResourceManagementPrivateLinkAssociationTestResource struct {
	uuid string
}

func TestAccResourceManagementPrivateLinkAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_management_private_link_association", "test")
	randomUUID, _ := uuid.GenerateUUID()
	r := ResourceManagementPrivateLinkAssociationTestResource{
		uuid: randomUUID,
	}

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

func TestAccResourceManagementPrivateLinkAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_management_private_link_association", "test")
	randomUUID, _ := uuid.GenerateUUID()
	r := ResourceManagementPrivateLinkAssociationTestResource{
		uuid: randomUUID,
	}

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

func TestAccResourceManagementPrivateLinkAssociation_generateName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_management_private_link_association", "test")
	r := ResourceManagementPrivateLinkAssociationTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.generateName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ResourceManagementPrivateLinkAssociationTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := privatelinkassociation.ParsePrivateLinkAssociationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Resource.PrivateLinkAssociationClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r ResourceManagementPrivateLinkAssociationTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_resource_management_private_link_association" "test" {
  name                                = "%s"
  management_group_id                 = data.azurerm_management_group.test.id
  resource_management_private_link_id = azurerm_resource_management_private_link.test.id
  public_network_access_enabled       = true
}
`, r.template(data), r.uuid)
}

func (r ResourceManagementPrivateLinkAssociationTestResource) generateName(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_resource_management_private_link_association" "test" {
  management_group_id                 = data.azurerm_management_group.test.id
  resource_management_private_link_id = azurerm_resource_management_private_link.test.id
  public_network_access_enabled       = true
  lifecycle {
    ignore_changes = [name]
  }
}
`, r.template(data))
}

func (r ResourceManagementPrivateLinkAssociationTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_resource_management_private_link_association" "import" {
  name                                = azurerm_resource_management_private_link_association.test.name
  management_group_id                 = azurerm_resource_management_private_link_association.test.management_group_id
  resource_management_private_link_id = azurerm_resource_management_private_link_association.test.resource_management_private_link_id
  public_network_access_enabled       = azurerm_resource_management_private_link_association.test.public_network_access_enabled
}
`, r.basic(data))
}

func (r ResourceManagementPrivateLinkAssociationTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "test" {}

data "azurerm_management_group" "test" {
  name = data.azurerm_client_config.test.tenant_id
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_resource_management_private_link" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctestrmpl-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}
