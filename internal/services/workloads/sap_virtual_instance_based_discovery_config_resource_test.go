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

type SAPVirtualInstanceBasedDiscoveryConfigResource struct{}

func TestAccSAPVirtualInstanceBasedDiscoveryConfig_basic(t *testing.T) {
	if os.Getenv("ARM_CENTRAL_SERVER_VM_ID") == "" || os.Getenv("ARM_MANAGED_STORAGE_ACCOUNT_NAME") == "" || os.Getenv("ARM_IDENTITY_ID") == "" {
		t.Skip("Skipping as ARM_CENTRAL_SERVER_VM_ID or ARM_MANAGED_STORAGE_ACCOUNT_NAME or ARM_IDENTITY_ID is not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_sap_virtual_instance_based_discovery_config", "test")
	r := SAPVirtualInstanceBasedDiscoveryConfigResource{}

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

func TestAccSAPVirtualInstanceBasedDiscoveryConfig_requiresImport(t *testing.T) {
	if os.Getenv("ARM_CENTRAL_SERVER_VM_ID") == "" || os.Getenv("ARM_MANAGED_STORAGE_ACCOUNT_NAME") == "" || os.Getenv("ARM_IDENTITY_ID") == "" {
		t.Skip("Skipping as ARM_CENTRAL_SERVER_VM_ID or ARM_MANAGED_STORAGE_ACCOUNT_NAME or ARM_IDENTITY_ID is not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_sap_virtual_instance_based_discovery_config", "test")
	r := SAPVirtualInstanceBasedDiscoveryConfigResource{}

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

func TestAccSAPVirtualInstanceBasedDiscoveryConfig_complete(t *testing.T) {
	if os.Getenv("ARM_CENTRAL_SERVER_VM_ID") == "" || os.Getenv("ARM_MANAGED_STORAGE_ACCOUNT_NAME") == "" || os.Getenv("ARM_IDENTITY_ID") == "" {
		t.Skip("Skipping as ARM_CENTRAL_SERVER_VM_ID or ARM_MANAGED_STORAGE_ACCOUNT_NAME or ARM_IDENTITY_ID is not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_sap_virtual_instance_based_discovery_config", "test")
	r := SAPVirtualInstanceBasedDiscoveryConfigResource{}

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

func TestAccSAPVirtualInstanceBasedDiscoveryConfig_update(t *testing.T) {
	if os.Getenv("ARM_CENTRAL_SERVER_VM_ID") == "" || os.Getenv("ARM_MANAGED_STORAGE_ACCOUNT_NAME") == "" || os.Getenv("ARM_IDENTITY_ID") == "" {
		t.Skip("Skipping as ARM_CENTRAL_SERVER_VM_ID or ARM_MANAGED_STORAGE_ACCOUNT_NAME or ARM_IDENTITY_ID is not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_sap_virtual_instance_based_discovery_config", "test")
	r := SAPVirtualInstanceBasedDiscoveryConfigResource{}

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

func (r SAPVirtualInstanceBasedDiscoveryConfigResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r SAPVirtualInstanceBasedDiscoveryConfigResource) template(data acceptance.TestData) string {
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

func (r SAPVirtualInstanceBasedDiscoveryConfigResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sap_virtual_instance_based_discovery_config" "test" {
  name                = "X%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  environment         = "NonProd"
  sap_product         = "S4HANA"

  configuration {
    central_server_vm_id         = "%s"
    managed_storage_account_name = "%s"
  }
}
`, r.template(data), data.RandomIntOfLength(2), os.Getenv("ARM_CENTRAL_SERVER_VM_ID"), os.Getenv("ARM_MANAGED_STORAGE_ACCOUNT_NAME"))
}

func (r SAPVirtualInstanceBasedDiscoveryConfigResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sap_virtual_instance_based_discovery_config" "import" {
  name                = azurerm_sap_virtual_instance_based_discovery_config.test.name
  resource_group_name = azurerm_sap_virtual_instance_based_discovery_config.test.name
  location            = azurerm_sap_virtual_instance_based_discovery_config.test.name
  environment         = azurerm_sap_virtual_instance_based_discovery_config.test.name
  sap_product         = azurerm_sap_virtual_instance_based_discovery_config.test.name

  configuration {
    central_server_vm_id         = azurerm_sap_virtual_instance_based_discovery_config.test.configuration.0.central_server_vm_id
    managed_storage_account_name = azurerm_sap_virtual_instance_based_discovery_config.test.configuration.0.managed_storage_account_name
  }
}
`, r.basic(data))
}

func (r SAPVirtualInstanceBasedDiscoveryConfigResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sap_virtual_instance_based_discovery_config" "test" {
  name                = "X%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  environment                 = "NonProd"
  sap_product                 = "S4HANA"
  managed_resource_group_name = "managedrg%s"

  configuration {
    central_server_vm_id         = "%s"
    managed_storage_account_name = "%s"
  }

  identity {
    type = "UserAssigned"
    
    identity_ids = ["%s"]
  }

  tags = {
    Env = "Test"
  }
}
`, r.template(data), data.RandomIntOfLength(2), data.RandomString, os.Getenv("ARM_CENTRAL_SERVER_VM_ID"), os.Getenv("ARM_MANAGED_STORAGE_ACCOUNT_NAME"), os.Getenv("ARM_IDENTITY_ID"))
}

func (r SAPVirtualInstanceBasedDiscoveryConfigResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sap_virtual_instance_based_discovery_config" "test" {
  name                = "X%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  environment                 = "NonProd"
  sap_product                 = "S4HANA"
  managed_resource_group_name = "managedrg%s"

  configuration {
    central_server_vm_id         = "%s"
    managed_storage_account_name = "%s"
  }

  identity {
    type = "UserAssigned"
    
    identity_ids = ["%s"]
  }

  tags = {
    Env = "Test2"
  }
}
`, r.template(data), data.RandomIntOfLength(2), data.RandomString, os.Getenv("ARM_CENTRAL_SERVER_VM_ID"), os.Getenv("ARM_MANAGED_STORAGE_ACCOUNT_NAME"), os.Getenv("ARM_IDENTITY_ID"))
}
