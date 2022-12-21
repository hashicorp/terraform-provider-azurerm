package labservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/lab"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LabServiceLabResource struct{}

func TestAccLabServicesLab_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_lab", "test")
	r := LabServiceLabResource{}

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

func TestAccLabServiceLab_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_lab", "test")
	r := LabServiceLabResource{}

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

func TestAccLabServiceLab_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_lab", "test")
	r := LabServiceLabResource{}

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

func TestAccLabServiceLab_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_lab", "test")
	r := LabServiceLabResource{}

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

func (r LabServiceLabResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := lab.ParseLabID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.LabService.LabClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r LabServiceLabResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-lab-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r LabServiceLabResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

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
      username = "testAdmin"
      password = "Password1234!"
    }

    image_reference {
      offer     = "WindowsServer"
      publisher = "MicrosoftWindowsServer"
      sku       = "2022-datacenter-g2"
      version   = "latest"
    }

    sku {
      name     = "Classic_Fsv2_2_4GB_128_S_SSD"
      capacity = 1
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LabServiceLabResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lab_service_lab" "import" {
  name                = azurerm_lab_service_lab.test.name
  resource_group_name = azurerm_lab_service_lab.test.resource_group_name
  location            = azurerm_lab_service_lab.test.location
  title               = azurerm_lab_service_lab.test.title

  security {
    open_access_enabled = false
  }

  virtual_machine {
    admin_user {
      username = "testAdmin"
      password = "Password1234!"
    }

    image_reference {
      offer     = "WindowsServer"
      publisher = "MicrosoftWindowsServer"
      sku       = "2022-datacenter-g2"
      version   = "latest"
    }

    sku {
      name     = "Classic_Fsv2_2_4GB_128_S_SSD"
      capacity = 1
    }
  }
}
`, r.basic(data))
}

func (r LabServiceLabResource) complete(data acceptance.TestData) string {
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

resource "azurerm_lab_service_plan" "test" {
  name                = "acctest-lslp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allowed_regions     = [azurerm_resource_group.test.location]
}

resource "azurerm_lab_service_lab" "test" {
  name                = "acctest-lab-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  title               = "Test Title"
  description         = "Test Description"
  lab_plan_id         = azurerm_lab_service_plan.test.id

  security {
    open_access_enabled = false
  }

  virtual_machine {
    usage_quota                                 = "PT10H"
    additional_capability_gpu_drivers_installed = false
    create_option                               = "Image"
    shared_password_enabled                     = true

    admin_user {
      username = "testAdmin"
      password = "Password1234!"
    }

    non_admin_user {
      username = "testNonAdmin"
      password = "Password1234!"
    }

    image_reference {
      offer     = "WindowsServer"
      publisher = "MicrosoftWindowsServer"
      sku       = "2022-datacenter-g2"
      version   = "latest"
    }

    sku {
      name     = "Classic_Fsv2_2_4GB_128_S_SSD"
      capacity = 1
    }
  }

  auto_shutdown {
    disconnect_delay = "PT17M"
    idle_delay       = "PT17M"
    no_connect_delay = "PT17M"
    shutdown_on_idle = "UserAbsence"
  }

  connection_setting {
    client_rdp_access = "Public"
  }

  network {
    subnet_id = azurerm_subnet.test.id
  }

  roster {
    lms_instance        = "https://terraform.io/"
    lti_context_id      = "b0a538ec-1b9d-4e00-b3c3-f689cb34a30c"
    lti_roster_endpoint = "https://terraform.io/"
  }

  tags = {
    Env = "Test"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r LabServiceLabResource) update(data acceptance.TestData) string {
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

resource "azurerm_lab_service_plan" "test" {
  name                = "acctest-lslp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allowed_regions     = [azurerm_resource_group.test.location]
}

resource "azurerm_lab_service_lab" "test" {
  name                = "acctest-lab-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  title               = "Test Title2"
  description         = "Test Description2"

  security {
    open_access_enabled = true
  }

  virtual_machine {
    usage_quota                                 = "PT11H"
    additional_capability_gpu_drivers_installed = false
    create_option                               = "Image"
    shared_password_enabled                     = true

    admin_user {
      username = "testAdmin"
      password = "Password1234!"
    }

    image_reference {
      offer     = "WindowsServer"
      publisher = "MicrosoftWindowsServer"
      sku       = "2022-datacenter-g2"
      version   = "latest"
    }

    sku {
      name     = "Classic_Fsv2_2_4GB_128_S_SSD"
      capacity = 0
    }
  }

  auto_shutdown {
    disconnect_delay = "PT16M"
    idle_delay       = "PT16M"
    no_connect_delay = "PT16M"
    shutdown_on_idle = "LowUsage"
  }

  roster {
    lms_instance        = "https://registry.terraform.io/"
    lti_context_id      = "72fa88e9-e65a-44f3-91ac-008226ae91dc"
    lti_roster_endpoint = "https://registry.terraform.io/"
  }

  tags = {
    Env = "Test2"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
