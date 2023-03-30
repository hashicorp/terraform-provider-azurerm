package workloads_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkloadsSAPVirtualInstanceResource struct{}

func TestAccWorkloadsSAPVirtualInstance_basic(t *testing.T) {
	if os.Getenv("ARM_CENTRAL_SERVER_VM_ID") == "" {
		t.Skip("Skipping as ARM_CENTRAL_SERVER_VM_ID is not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_virtual_instance", "test")
	r := WorkloadsSAPVirtualInstanceResource{}

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

func TestAccWorkloadsSAPVirtualInstance_requiresImport(t *testing.T) {
	if os.Getenv("ARM_CENTRAL_SERVER_VM_ID") == "" {
		t.Skip("Skipping as ARM_CENTRAL_SERVER_VM_ID or ARM_IDENTITY_ID is not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_virtual_instance", "test")
	r := WorkloadsSAPVirtualInstanceResource{}

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

func TestAccWorkloadsSAPVirtualInstance_complete(t *testing.T) {
	if os.Getenv("ARM_CENTRAL_SERVER_VM_ID") == "" || os.Getenv("ARM_IDENTITY_ID") == "" {
		t.Skip("Skipping as ARM_CENTRAL_SERVER_VM_ID or ARM_IDENTITY_ID is not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_virtual_instance", "test")
	r := WorkloadsSAPVirtualInstanceResource{}

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

func TestAccWorkloadsSAPVirtualInstance_update(t *testing.T) {
	if os.Getenv("ARM_CENTRAL_SERVER_VM_ID") == "" || os.Getenv("ARM_IDENTITY_ID") == "" {
		t.Skip("Skipping as ARM_CENTRAL_SERVER_VM_ID or ARM_IDENTITY_ID is not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_virtual_instance", "test")
	r := WorkloadsSAPVirtualInstanceResource{}

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

func (r WorkloadsSAPVirtualInstanceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := sapvirtualinstances.ParseSapVirtualInstanceID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Workloads.SAPVirtualInstances
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r WorkloadsSAPVirtualInstanceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vis-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r WorkloadsSAPVirtualInstanceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_workloads_sap_virtual_instance" "test" {
  name                = "X%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  environment         = "NonProd"
  sap_product         = "S4HANA"

  discovery_configuration {
    central_server_vm_id = "%s"
  }
}
`, r.template(data), data.RandomIntOfLength(2), os.Getenv("ARM_CENTRAL_SERVER_VM_ID"))
}

func (r WorkloadsSAPVirtualInstanceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_workloads_sap_virtual_instance" "import" {
  name                = azurerm_workloads_sap_virtual_instance.test.name
  resource_group_name = azurerm_workloads_sap_virtual_instance.test.name
  location            = azurerm_workloads_sap_virtual_instance.test.name
  environment         = azurerm_workloads_sap_virtual_instance.test.name
  sap_product         = azurerm_workloads_sap_virtual_instance.test.name

  discovery_configuration {
    central_server_vm_id = azurerm_workloads_sap_virtual_instance.test.configuration.0.central_server_vm_id
  }
}
`, r.basic(data))
}

func (r WorkloadsSAPVirtualInstanceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_workloads_sap_virtual_instance" "test" {
  name                = "X%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  environment                 = "NonProd"
  sap_product                 = "S4HANA"
  managed_resource_group_name = "managedrg%s"

  discovery_configuration {
    central_server_vm_id         = "%s"
    managed_storage_account_name = "managedsa%s"
  }

  identity {
    type = "UserAssigned"
    
    identity_ids = ["%s"]
  }

  tags = {
    Env = "Test"
  }
}
`, r.template(data), data.RandomIntOfLength(2), data.RandomString, os.Getenv("ARM_CENTRAL_SERVER_VM_ID"), data.RandomString, os.Getenv("ARM_IDENTITY_ID"))
}

func (r WorkloadsSAPVirtualInstanceResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_workloads_sap_virtual_instance" "test" {
  name                = "X%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  environment                 = "NonProd"
  sap_product                 = "S4HANA"
  managed_resource_group_name = "managedrg%s"

  discovery_configuration {
    central_server_vm_id         = "%s"
    managed_storage_account_name = "managedsa%s"
  }

  identity {
    type = "UserAssigned"
    
    identity_ids = ["%s"]
  }

  tags = {
    Env = "Test2"
  }
}
`, r.template(data), data.RandomIntOfLength(2), data.RandomString, os.Getenv("ARM_CENTRAL_SERVER_VM_ID"), data.RandomString, os.Getenv("ARM_IDENTITY_ID"))
}
