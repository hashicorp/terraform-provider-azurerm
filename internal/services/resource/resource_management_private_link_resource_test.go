// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/resourcemanagementprivatelink"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ResourceManagementPrivateLinkTestResource struct{}

func TestAccResourceManagementPrivateLink_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_management_private_link", "test")
	r := ResourceManagementPrivateLinkTestResource{}

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

func TestAccResourceManagementPrivateLink_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_management_private_link", "test")
	r := ResourceManagementPrivateLinkTestResource{}

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

func (r ResourceManagementPrivateLinkTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := resourcemanagementprivatelink.ParseResourceManagementPrivateLinkID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Resource.ResourceManagementPrivateLinkClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r ResourceManagementPrivateLinkTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_resource_management_private_link" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctestrmpl-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data))
}

func (r ResourceManagementPrivateLinkTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_resource_management_private_link" "import" {
  location            = azurerm_resource_management_private_link.test.location
  name                = azurerm_resource_management_private_link.test.name
  resource_group_name = azurerm_resource_management_private_link.test.resource_group_name
}
`, r.basic(data))
}

func (r ResourceManagementPrivateLinkTestResource) template(data acceptance.TestData) string {
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
