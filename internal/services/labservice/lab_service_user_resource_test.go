// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package labservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/user"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LabServiceUserResource struct{}

func TestAccLabServiceUser_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_user", "test")
	r := LabServiceUserResource{}

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

func TestAccLabServiceUser_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_user", "test")
	r := LabServiceUserResource{}

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

func TestAccLabServiceUser_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_user", "test")
	r := LabServiceUserResource{}

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

func TestAccLabServiceUser_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_user", "test")
	r := LabServiceUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
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

func (r LabServiceUserResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := user.ParseUserID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.LabService.UserClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r LabServiceUserResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-labuser-%d"
  location = "%s"
}

resource "azurerm_lab_service_lab" "test" {
  name                = "acctest-lab-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  title               = "Test Title"

  security {
    open_access_enabled = false
  }

  virtual_machine {
    admin_user {
      username = "testadmin"
      password = "Password1234!"
    }

    image_reference {
      offer     = "0001-com-ubuntu-server-focal"
      publisher = "canonical"
      sku       = "20_04-lts"
      version   = "latest"
    }

    sku {
      name     = "Classic_Fsv2_2_4GB_128_S_SSD"
      capacity = 1
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(17))
}

func (r LabServiceUserResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lab_service_user" "test" {
  name   = "acctest-labuser-%d"
  lab_id = azurerm_lab_service_lab.test.id
  email  = "terraform-acctest@hashicorp.com"
}
`, r.template(data), data.RandomInteger)
}

func (r LabServiceUserResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lab_service_user" "import" {
  name   = azurerm_lab_service_user.test.name
  lab_id = azurerm_lab_service_user.test.lab_id
  email  = azurerm_lab_service_user.test.email
}
`, r.basic(data))
}

func (r LabServiceUserResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lab_service_user" "test" {
  name                   = "acctest-lsu-%d"
  lab_id                 = azurerm_lab_service_lab.test.id
  email                  = "terraform-acctest@hashicorp.com"
  additional_usage_quota = "PT10H"
}
`, r.template(data), data.RandomInteger)
}

func (r LabServiceUserResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lab_service_user" "test" {
  name                   = "acctest-lsu-%d"
  lab_id                 = azurerm_lab_service_lab.test.id
  email                  = "terraform-acctest@hashicorp.com"
  additional_usage_quota = "PT11H"
}
`, r.template(data), data.RandomInteger)
}
