package devcenter_test

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/galleries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DevCenterGalleryTestResource struct{}

func TestAccDevCenterGallery_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_gallery", "test")
	r := DevCenterGalleryTestResource{}

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

func TestAccDevCenterGallery_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_gallery", "test")
	r := DevCenterGalleryTestResource{}

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

func TestAccDevCenterGallery_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_gallery", "test")
	r := DevCenterGalleryTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDevCenterGallery_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_gallery", "test")
	r := DevCenterGalleryTestResource{}

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

func (r DevCenterGalleryTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := galleries.ParseGalleryID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DevCenter.V20250201.Galleries.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r DevCenterGalleryTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center_gallery" "test" {
  dev_center_id     = azurerm_dev_center.test.id
  shared_gallery_id = azurerm_shared_image_gallery.test.id
  name              = "acctestdcg${var.random_string}"
}
`, r.template(data))
}

func (r DevCenterGalleryTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dev_center_gallery" "import" {
  dev_center_id     = azurerm_dev_center_gallery.test.dev_center_id
  shared_gallery_id = azurerm_dev_center_gallery.test.shared_gallery_id
  name              = azurerm_dev_center_gallery.test.name
}
`, r.basic(data))
}

func (r DevCenterGalleryTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center_gallery" "test" {
  dev_center_id     = azurerm_dev_center.test.id
  shared_gallery_id = azurerm_shared_image_gallery.test.id
  name              = "acctestdcg${var.random_string}"
}
`, r.template(data))
}

func (r DevCenterGalleryTestResource) template(data acceptance.TestData) string {
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

resource "azurerm_dev_center" "test" {
  name                = "acctestdc-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}


resource "azurerm_resource_group" "test" {
  name     = "acctestrg-${var.random_integer}"
  location = var.primary_location
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuami-${var.random_string}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}


resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig${var.random_string}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}


resource "azurerm_role_assignment" "test" {
  scope                = azurerm_shared_image_gallery.test.id
  role_definition_name = "Owner"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
