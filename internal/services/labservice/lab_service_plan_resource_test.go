// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package labservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/labplan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LabServicePlanResource struct{}

func TestAccLabServicePlan_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_plan", "test")
	r := LabServicePlanResource{}

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

func TestAccLabServicePlan_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_plan", "test")
	r := LabServicePlanResource{}

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

func TestAccLabServicePlan_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_plan", "test")
	r := LabServicePlanResource{}

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

func TestAccLabServicePlan_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_plan", "test")
	r := LabServicePlanResource{}

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

func (r LabServicePlanResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := labplan.ParseLabPlanID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.LabService.LabPlanClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r LabServicePlanResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-lslp-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r LabServicePlanResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lab_service_plan" "test" {
  name                = "acctest-lslp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allowed_regions     = [azurerm_resource_group.test.location]
}
`, r.template(data), data.RandomInteger)
}

func (r LabServicePlanResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lab_service_plan" "import" {
  name                = azurerm_lab_service_plan.test.name
  resource_group_name = azurerm_lab_service_plan.test.resource_group_name
  location            = azurerm_lab_service_plan.test.location
  allowed_regions     = azurerm_lab_service_plan.test.allowed_regions
}
`, r.basic(data))
}

func (r LabServicePlanResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

data "azuread_service_principal" "test" {
  application_id = "c7bb12bf-0b39-4f7f-9171-f418ff39b76a"
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_shared_image_gallery.test.id
  role_definition_name = "Contributor"
  principal_id         = data.azuread_service_principal.test.object_id
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-sn-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "Microsoft.LabServices.labplans"

    service_delegation {
      name = "Microsoft.LabServices/labplans"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_lab_service_plan" "test" {
  name                      = "acctest-lslp-%d"
  resource_group_name       = azurerm_resource_group.test.name
  location                  = azurerm_resource_group.test.location
  allowed_regions           = [azurerm_resource_group.test.location]
  default_network_subnet_id = azurerm_subnet.test.id
  shared_gallery_id         = azurerm_shared_image_gallery.test.id

  default_auto_shutdown {
    disconnect_delay = "PT15M"
    idle_delay       = "PT15M"
    no_connect_delay = "PT15M"
    shutdown_on_idle = "LowUsage"
  }

  default_connection {
    client_ssh_access = "Public"
    web_ssh_access    = "Public"
  }

  support {
    email        = "company@terraform.io"
    instructions = "Contact support for help"
    phone        = "+1-555-555-5555"
    url          = "https://www.terraform.io/"
  }

  tags = {
    Env = "Test"
  }

  depends_on = [azurerm_role_assignment.test]
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r LabServicePlanResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-sn-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "Microsoft.LabServices.labplans"

    service_delegation {
      name = "Microsoft.LabServices/labplans"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_subnet" "test2" {
  name                 = "acctest-sn2-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.3.0/24"]

  delegation {
    name = "Microsoft.LabServices.labplans"

    service_delegation {
      name = "Microsoft.LabServices/labplans"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_lab_service_plan" "test" {
  name                      = "acctest-lslp-%d"
  resource_group_name       = azurerm_resource_group.test.name
  location                  = azurerm_resource_group.test.location
  allowed_regions           = [azurerm_resource_group.test.location, "%s"]
  default_network_subnet_id = azurerm_subnet.test2.id

  default_auto_shutdown {
    disconnect_delay = "PT16M"
    idle_delay       = "PT16M"
    no_connect_delay = "PT16M"
    shutdown_on_idle = "UserAbsence"
  }

  default_connection {
    client_rdp_access = "Public"
    web_rdp_access    = "Public"
  }

  support {
    email        = "company2@terraform.io"
    instructions = "Contacting support for help"
    phone        = "+1-555-555-6666"
    url          = "https://www.terraform2.io/"
  }

  tags = {
    Env = "Test2"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.Locations.Secondary)
}
