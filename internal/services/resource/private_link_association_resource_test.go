package resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/privatelinkassociation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PrivateLinkAssociationTestResource struct{}

func TestAccPrivateLinkAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_link_association", "test")
	r := PrivateLinkAssociationTestResource{}

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

func TestAccPrivateLinkAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_link_association", "test")
	r := PrivateLinkAssociationTestResource{}

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

func (r PrivateLinkAssociationTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r PrivateLinkAssociationTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_link_association" "test" {
 management_group_id           = data.azurerm_management_group.test.id
 private_link_id               = azurerm_resource_management_private_link.test.id
 public_network_access_enabled = true
}
`, r.template(data))
}

func (r PrivateLinkAssociationTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_link_association" "import" {
 name                          = azurerm_private_link_association.test.name
 management_group_id           = azurerm_private_link_association.test.management_group_id
 private_link_id               = azurerm_private_link_association.test.private_link_id
 public_network_access_enabled = azurerm_private_link_association.test.public_network_access_enabled
}
`, r.basic(data))
}

func (r PrivateLinkAssociationTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "primary_location" {
 default = %q
}
variable "random_string" {
 default = %q
}
variable "random_integer" {
 default = %d
}

provider "azurerm" {
 features {}
}

data "azurerm_client_config" "test" {}

data "azurerm_management_group" "test" {
 name = data.azurerm_client_config.test.tenant_id
}

resource "azurerm_resource_group" "test" {
 name     = "acctestrg-${var.random_integer}"
 location = var.primary_location
}

resource "azurerm_resource_management_private_link" "test" {
 location            = azurerm_resource_group.test.location
 name                = "acctestrmpl-${var.random_string}"
 resource_group_name = azurerm_resource_group.test.name
}
`, data.Locations.Primary, data.RandomString, data.RandomInteger)
}
