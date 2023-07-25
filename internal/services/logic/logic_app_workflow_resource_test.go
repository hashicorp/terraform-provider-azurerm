// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogicAppWorkflowResource struct{}

func TestAccLogicAppWorkflow_empty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_workflow", "test")
	r := LogicAppWorkflowResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("parameters.%").HasValue("0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("connector_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("connector_outbound_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_outbound_ip_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppWorkflow_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_workflow", "test")
	r := LogicAppWorkflowResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_logic_app_workflow"),
		},
	})
}

func TestAccLogicAppWorkflow_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_workflow", "test")
	r := LogicAppWorkflowResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("parameters.%").HasValue("0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.tags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("parameters.%").HasValue("0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Source").HasValue("AcceptanceTests"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppWorkflow_integrationAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_workflow", "test")
	r := LogicAppWorkflowResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.integrationAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("logic_app_integration_account_id"),
		{
			Config: r.integrationAccountUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("logic_app_integration_account_id"),
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("logic_app_integration_account_id"),
	})
}

func TestAccLogicAppWorkflow_integrationServiceEnvironment(t *testing.T) {
	t.Skip("skip as Integration Service Environment is being deprecated")

	data := acceptance.BuildTestData(t, "azurerm_logic_app_workflow", "test")
	r := LogicAppWorkflowResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.integrationServiceEnvironment(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.integrationServiceEnvironmentUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppWorkflow_parameters(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_workflow", "test")
	r := LogicAppWorkflowResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.parameters(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parameters.secstr", "parameters.secobj"),
	})
}

func TestAccLogicAppWorkflow_accessControl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_workflow", "test")
	r := LogicAppWorkflowResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.accessControl(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateAccessControl(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateAccessControlWithStep(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccLogicAppWorkflow_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_workflow", "test")
	r := LogicAppWorkflowResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.userAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (LogicAppWorkflowResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := workflows.ParseWorkflowID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Logic.WorkflowClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving Logic App Workflow %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (LogicAppWorkflowResource) empty(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-logic-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r LogicAppWorkflowResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_workflow" "import" {
  name                = azurerm_logic_app_workflow.test.name
  location            = azurerm_logic_app_workflow.test.location
  resource_group_name = azurerm_logic_app_workflow.test.resource_group_name
}
`, r.empty(data))
}

func (LogicAppWorkflowResource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "Source" = "AcceptanceTests"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogicAppWorkflowResource) integrationAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-logic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_logic_app_integration_account" "test" {
  name                = "acctest-IA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_logic_app_integration_account" "test2" {
  name                = "acctest-IA2-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_logic_app_workflow" "test" {
  name                             = "acctestlaw-%[1]d"
  location                         = azurerm_resource_group.test.location
  resource_group_name              = azurerm_resource_group.test.name
  logic_app_integration_account_id = azurerm_logic_app_integration_account.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (LogicAppWorkflowResource) integrationAccountUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-logic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_logic_app_integration_account" "test" {
  name                = "acctest-IA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_logic_app_integration_account" "test2" {
  name                = "acctest-IA2-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_logic_app_workflow" "test" {
  name                             = "acctestlaw-%[1]d"
  location                         = azurerm_resource_group.test.location
  resource_group_name              = azurerm_resource_group.test.name
  logic_app_integration_account_id = azurerm_logic_app_integration_account.test2.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r LogicAppWorkflowResource) integrationServiceEnvironment(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_workflow" "test" {
  name                               = "acctestlaw-%d"
  location                           = azurerm_resource_group.test.location
  resource_group_name                = azurerm_resource_group.test.name
  integration_service_environment_id = azurerm_integration_service_environment.test.id
}
`, IntegrationServiceEnvironmentResource{}.basic(data), data.RandomInteger)
}

func (r LogicAppWorkflowResource) integrationServiceEnvironmentUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_workflow" "test" {
  name                               = "acctestlaw-%d"
  location                           = azurerm_resource_group.test.location
  resource_group_name                = azurerm_resource_group.test.name
  integration_service_environment_id = azurerm_integration_service_environment.test.id

  tags = {
    "Source" = "AcceptanceTests"
  }
}
`, IntegrationServiceEnvironmentResource{}.basic(data), data.RandomInteger)
}

func (LogicAppWorkflowResource) parameters(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-logic-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workflow_parameters = {
    b = jsonencode({
      type = "Bool"
    })
    str = jsonencode({
      type = "String"
    })
    int = jsonencode({
      type = "Int"
    })
    float = jsonencode({
      type = "Float"
    })
    obj = jsonencode({
      type = "Object"
    })
    array = jsonencode({
      type = "Array"
    })
    secstr = jsonencode({
      type = "SecureString"
    })
    secobj = jsonencode({
      type = "SecureObject"
    })
  }

  parameters = {
    b     = "true"
    str   = "value"
    int   = "123"
    float = "1.23"
    obj = jsonencode({
      s     = "foo"
      array = [1, 2, 3]
      obj = {
        i = 123
      }
    })
    array = jsonencode([
      1, "string", {}, []
    ])
    secstr = "value"
    secobj = jsonencode({
      foo = "foo"
    })
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogicAppWorkflowResource) accessControl(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-logic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  access_control {
    content {
      allowed_caller_ip_address_range = ["10.0.5.0-10.0.5.10"]
    }

    action {
      allowed_caller_ip_address_range = ["10.0.6.0-10.0.6.10"]
    }

    trigger {
      allowed_caller_ip_address_range = ["10.0.7.0-10.0.7.10"]

      open_authentication_policy {
        name = "testpolicy1"

        claim {
          name  = "iss"
          value = "https://sts.windows.net/"
        }

        claim {
          name  = "aud"
          value = "https://management.core.windows.net/"
        }
      }
    }

    workflow_management {
      allowed_caller_ip_address_range = ["10.0.8.0-10.0.8.10"]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (LogicAppWorkflowResource) updateAccessControl(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-logic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  enabled             = false

  access_control {
    content {
      allowed_caller_ip_address_range = ["10.10.3.0/24"]
    }

    action {
      allowed_caller_ip_address_range = ["10.10.4.0/24"]
    }

    trigger {
      allowed_caller_ip_address_range = ["10.10.5.0/24"]

      open_authentication_policy {
        name = "testpolicy4"

        claim {
          name  = "iss"
          value = "https://sts.windows.net/"
        }

        claim {
          name  = "testclaimname"
          value = "testclaimvalue"
        }
      }
    }

    workflow_management {
      allowed_caller_ip_address_range = ["10.10.6.0/24"]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r LogicAppWorkflowResource) updateAccessControlWithStep(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_custom" "test" {
  name         = "test"
  logic_app_id = azurerm_logic_app_workflow.test.id

  body = <<BODY
{
    "description": "A variable to configure the auto expiration age in days. Configured in negative number. Default is -30 (30 days old).",
    "inputs": {
        "variables": [
            {
                "name": "ExpirationAgeInDays",
                "type": "Integer",
                "value": -30
            }
        ]
    },
    "runAfter": {},
    "type": "InitializeVariable"
}
BODY
}

`, r.updateAccessControl(data))
}

func (LogicAppWorkflowResource) systemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-logic-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogicAppWorkflowResource) userAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-logic-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-user-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = "Logic App Operator"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  depends_on = [azurerm_role_assignment.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}
