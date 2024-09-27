// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedhsm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KeyVaultManagedHSMRoleAssignmentResource struct{}

// NOTE: all fields within the Role Assignment are ForceNew, therefore any Update tests aren't going to do much..

func testAccKeyVaultManagedHardwareSecurityModuleRoleAssignment_builtInRole(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_hardware_security_module_role_assignment", "test")
	r := KeyVaultManagedHSMRoleAssignmentResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.builtInRole(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccKeyVaultManagedHardwareSecurityModuleRoleAssignment_customRole(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_hardware_security_module_role_assignment", "test")
	r := KeyVaultManagedHSMRoleAssignmentResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.customRole(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccKeyVaultManagedHardwareSecurityModuleRoleAssignment_legacyBuiltInRole(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("This test isn't applicable in 4.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_hardware_security_module_role_assignment", "test")
	r := KeyVaultManagedHSMRoleAssignmentResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.legacyBuiltInRole(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccKeyVaultManagedHardwareSecurityModuleRoleAssignment_legacyCustomRole(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("This test isn't applicable in 4.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_hardware_security_module_role_assignment", "test")
	r := KeyVaultManagedHSMRoleAssignmentResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.legacyCustomRole(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// real test nested in TestAccKeyVaultManagedHardwareSecurityModule, only provide Exists logic here
func (r KeyVaultManagedHSMRoleAssignmentResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	domainSuffix, ok := client.Account.Environment.ManagedHSM.DomainSuffix()
	if !ok {
		return nil, fmt.Errorf("this Environment doesn't specify the Domain Suffix for Managed HSM")
	}
	id, err := parse.ManagedHSMDataPlaneRoleAssignmentID(state.ID, domainSuffix)
	if err != nil {
		return nil, err
	}
	resp, err := client.ManagedHSMs.DataPlaneRoleAssignmentsClient.Get(ctx, id.BaseURI(), id.Scope, id.RoleAssignmentName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Properties != nil), nil
}

func (r KeyVaultManagedHSMRoleAssignmentResource) builtInRole(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

locals {
  assignmentOfficerName = "706c03c7-69ad-33e5-2796-b3380d3a6e1a"
}

data "azurerm_key_vault_managed_hardware_security_module_role_definition" "officer" {
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  name           = "515eb02d-2335-4d2d-92f2-b1cbdf9c3778"
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "test" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = local.assignmentOfficerName
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.officer.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id
}
`, KeyVaultManagedHardwareSecurityModuleResource{}.download(data, 3))
}

func (r KeyVaultManagedHSMRoleAssignmentResource) customRole(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

locals {
  assignmentTestName = "1e243909-064c-6ac3-84e9-1c8bf8d6ad52"
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "test" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = local.assignmentTestName
  scope              = "/keys"
  role_definition_id = azurerm_key_vault_managed_hardware_security_module_role_definition.test.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id
}
`, KeyVaultMHSMRoleDefinitionResource{}.basic(data))
}

func (r KeyVaultManagedHSMRoleAssignmentResource) legacyBuiltInRole(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

locals {
  assignmentOfficerName = "706c03c7-69ad-33e5-2796-b3380d3a6e1a"
}

data "azurerm_key_vault_managed_hardware_security_module_role_definition" "officer" {
  vault_base_url = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  name           = "515eb02d-2335-4d2d-92f2-b1cbdf9c3778"
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "test" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = local.assignmentOfficerName
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.officer.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id
}
`, KeyVaultManagedHardwareSecurityModuleResource{}.download(data, 3))
}

func (r KeyVaultManagedHSMRoleAssignmentResource) legacyCustomRole(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

locals {
  assignmentTestName = "1e243909-064c-6ac3-84e9-1c8bf8d6ad52"
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "test" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = local.assignmentTestName
  scope              = "/keys"
  role_definition_id = azurerm_key_vault_managed_hardware_security_module_role_definition.test.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id
}
`, KeyVaultMHSMRoleDefinitionResource{}.basic(data))
}
