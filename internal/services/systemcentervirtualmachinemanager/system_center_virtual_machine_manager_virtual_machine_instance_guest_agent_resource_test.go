package systemcentervirtualmachinemanager_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource struct{}

func TestAccSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentSequential(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests because the testing is against the same Hybrid Machine

	if os.Getenv("ARM_TEST_SCOPE_ID") == "" || os.Getenv("ARM_TEST_USERNAME") == "" || os.Getenv("ARM_TEST_PASSWORD") == "" {
		t.Skip("Skipping as one of `ARM_TEST_SCOPE_ID`, `ARM_TEST_USERNAME`, `ARM_TEST_PASSWORD` was not specified")
	}

	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"scvmmInstanceGuestAgent": {
			"basic":          testAccSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgent_basic,
			"requiresImport": testAccSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgent_requiresImport,
			"complete":       testAccSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgent_complete,
			"update":         testAccSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgent_update,
		},
	})
}

func testAccSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgent_basic(t *testing.T) {
	if os.Getenv("ARM_TEST_SCOPE_ID") == "" || os.Getenv("ARM_TEST_USERNAME") == "" || os.Getenv("ARM_TEST_PASSWORD") == "" {
		t.Skip("Skipping as one of `ARM_TEST_SCOPE_ID`, `ARM_TEST_USERNAME`, `ARM_TEST_PASSWORD` was not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent", "test")
	r := SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("credential.0.password"),
	})
}

func testAccSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgent_requiresImport(t *testing.T) {
	if os.Getenv("ARM_TEST_SCOPE_ID") == "" || os.Getenv("ARM_TEST_USERNAME") == "" || os.Getenv("ARM_TEST_PASSWORD") == "" {
		t.Skip("Skipping as one of `ARM_TEST_SCOPE_ID`, `ARM_TEST_USERNAME`, `ARM_TEST_PASSWORD` was not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent", "test")
	r := SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource{}

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

func testAccSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgent_complete(t *testing.T) {
	if os.Getenv("ARM_TEST_SCOPE_ID") == "" || os.Getenv("ARM_TEST_USERNAME") == "" || os.Getenv("ARM_TEST_PASSWORD") == "" {
		t.Skip("Skipping as one of `ARM_TEST_SCOPE_ID`, `ARM_TEST_USERNAME`, `ARM_TEST_PASSWORD` was not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent", "test")
	r := SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("credential.0.password"),
	})
}

func testAccSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgent_update(t *testing.T) {
	if os.Getenv("ARM_TEST_SCOPE_ID") == "" || os.Getenv("ARM_TEST_USERNAME") == "" || os.Getenv("ARM_TEST_PASSWORD") == "" {
		t.Skip("Skipping as one of `ARM_TEST_SCOPE_ID`, `ARM_TEST_USERNAME`, `ARM_TEST_PASSWORD` was not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent", "test")
	r := SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("credential.0.password"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("credential.0.password"),
	})
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SystemCenterVirtualMachineManager.VMInstanceGuestAgents.Get(ctx, commonids.NewScopeID(id.Scope))
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent" "test" {
  scoped_resource_id = "%s"

  credential {
    username = "%s"
    password = "%s"
  }
}
`, r.template(data), os.Getenv("ARM_TEST_SCOPE_ID"), os.Getenv("ARM_TEST_USERNAME"), os.Getenv("ARM_TEST_PASSWORD"))
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent" "import" {
  scoped_resource_id = azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent.test.scope

  credential {
    username = "%s"
    password = "%s"
  }
}
`, r.basic(data), os.Getenv("ARM_TEST_USERNAME"), os.Getenv("ARM_TEST_PASSWORD"))
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent" "test" {
  scoped_resource_id  = "%s"
  https_proxy         = ""
  provisioning_action = "install"

  credential {
    username = "%s"
    password = "%s"
  }
}
`, r.template(data), os.Getenv("ARM_TEST_SCOPE_ID"), os.Getenv("ARM_TEST_USERNAME"), os.Getenv("ARM_TEST_PASSWORD"))
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-scvmmiga-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
