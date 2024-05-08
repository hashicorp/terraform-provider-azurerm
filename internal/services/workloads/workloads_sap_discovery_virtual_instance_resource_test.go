// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

type WorkloadsSAPDiscoveryVirtualInstanceResource struct{}

func TestAccWorkloadsSAPDiscoveryVirtualInstanceSequential(t *testing.T) {
	// The dependent central server VM requires many complicated manual configurations. So it has to test based on the resource provided by service team.
	if os.Getenv("ARM_TEST_SAP_VIRTUAL_INSTANCE_NAME") == "" || os.Getenv("ARM_TEST_CENTRAL_SERVER_VM_ID") == "" || os.Getenv("ARM_TEST_IDENTITY_ID") == "" {
		t.Skip("Skipping as `ARM_TEST_SAP_VIRTUAL_INSTANCE_NAME`, `ARM_TEST_CENTRAL_SERVER_VM_ID` and `ARM_TEST_IDENTITY_ID` are not specified")
	}

	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"sapVirtualInstance": {
			"basic":          testAccWorkloadsSAPDiscoveryVirtualInstance_basic,
			"requiresImport": testAccWorkloadsSAPDiscoveryVirtualInstance_requiresImport,
			"complete":       testAccWorkloadsSAPDiscoveryVirtualInstance_complete,
			"update":         testAccWorkloadsSAPDiscoveryVirtualInstance_update,
		},
	})
}

func testAccWorkloadsSAPDiscoveryVirtualInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_discovery_virtual_instance", "test")
	r := WorkloadsSAPDiscoveryVirtualInstanceResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccWorkloadsSAPDiscoveryVirtualInstance_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_discovery_virtual_instance", "test")
	r := WorkloadsSAPDiscoveryVirtualInstanceResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccWorkloadsSAPDiscoveryVirtualInstance_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_discovery_virtual_instance", "test")
	r := WorkloadsSAPDiscoveryVirtualInstanceResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccWorkloadsSAPDiscoveryVirtualInstance_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_discovery_virtual_instance", "test")
	r := WorkloadsSAPDiscoveryVirtualInstanceResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
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
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r WorkloadsSAPDiscoveryVirtualInstanceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r WorkloadsSAPDiscoveryVirtualInstanceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sapvis-%d"
  location = "%s"
}

resource "azurerm_workloads_sap_discovery_virtual_instance" "test" {
  name                              = "%s"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = azurerm_resource_group.test.location
  environment                       = "NonProd"
  sap_product                       = "S4HANA"
  central_server_virtual_machine_id = "%s"
  managed_storage_account_name      = "acctestmanagedsa%s"

  identity {
    type = "UserAssigned"

    identity_ids = [
      "%s",
    ]
  }

  lifecycle {
    ignore_changes = [managed_resource_group_name]
  }
}
`, data.RandomInteger, data.Locations.Primary, os.Getenv("ARM_TEST_SAP_VIRTUAL_INSTANCE_NAME"), os.Getenv("ARM_TEST_CENTRAL_SERVER_VM_ID"), data.RandomString, os.Getenv("ARM_TEST_IDENTITY_ID"))
}

func (r WorkloadsSAPDiscoveryVirtualInstanceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_workloads_sap_discovery_virtual_instance" "import" {
  name                              = azurerm_workloads_sap_discovery_virtual_instance.test.name
  resource_group_name               = azurerm_workloads_sap_discovery_virtual_instance.test.resource_group_name
  location                          = azurerm_workloads_sap_discovery_virtual_instance.test.location
  environment                       = azurerm_workloads_sap_discovery_virtual_instance.test.environment
  sap_product                       = azurerm_workloads_sap_discovery_virtual_instance.test.sap_product
  central_server_virtual_machine_id = "%s"
  managed_storage_account_name      = "acctestmanagedsa%s"

  identity {
    type = "UserAssigned"

    identity_ids = [
      "%s",
    ]
  }
}
`, r.basic(data), os.Getenv("ARM_TEST_CENTRAL_SERVER_VM_ID"), data.RandomString, os.Getenv("ARM_TEST_IDENTITY_ID"))
}

func (r WorkloadsSAPDiscoveryVirtualInstanceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sapvis-%d"
  location = "%s"
}

resource "azurerm_workloads_sap_discovery_virtual_instance" "test" {
  name                              = "%s"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = azurerm_resource_group.test.location
  environment                       = "NonProd"
  sap_product                       = "S4HANA"
  managed_resource_group_name       = "acctestmanagedRG%d"
  central_server_virtual_machine_id = "%s"
  managed_storage_account_name      = "acctestmanagedsa%s"

  identity {
    type = "UserAssigned"

    identity_ids = [
      "%s",
    ]
  }

  tags = {
    env = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, os.Getenv("ARM_TEST_SAP_VIRTUAL_INSTANCE_NAME"), data.RandomInteger, os.Getenv("ARM_TEST_CENTRAL_SERVER_VM_ID"), data.RandomString, os.Getenv("ARM_TEST_IDENTITY_ID"))
}

func (r WorkloadsSAPDiscoveryVirtualInstanceResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sapvis-%d"
  location = "%s"
}

resource "azurerm_workloads_sap_discovery_virtual_instance" "test" {
  name                              = "%s"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = azurerm_resource_group.test.location
  environment                       = "NonProd"
  sap_product                       = "S4HANA"
  managed_resource_group_name       = "acctestmanagedRG%d"
  central_server_virtual_machine_id = "%s"
  managed_storage_account_name      = "acctestmanagedsa%s"

  tags = {
    env = "Test2"
  }
}
`, data.RandomInteger, data.Locations.Primary, os.Getenv("ARM_TEST_SAP_VIRTUAL_INSTANCE_NAME"), data.RandomInteger, os.Getenv("ARM_TEST_CENTRAL_SERVER_VM_ID"), data.RandomString)
}
