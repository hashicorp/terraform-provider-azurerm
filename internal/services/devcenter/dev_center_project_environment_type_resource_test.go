// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/environmenttypes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DevCenterProjectEnvironmentTypeTestResource struct{}

func TestAccDevCenterProjectEnvironmentType_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_project_environment_type", "test")
	r := DevCenterProjectEnvironmentTypeTestResource{}

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

func TestAccDevCenterProjectEnvironmentType_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_project_environment_type", "test")
	r := DevCenterProjectEnvironmentTypeTestResource{}

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

func TestAccDevCenterProjectEnvironmentType_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_project_environment_type", "test")
	r := DevCenterProjectEnvironmentTypeTestResource{}

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

func TestAccDevCenterProjectEnvironmentType_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_project_environment_type", "test")
	r := DevCenterProjectEnvironmentTypeTestResource{}

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
			Config: r.update(data),
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

func (r DevCenterProjectEnvironmentTypeTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := environmenttypes.ParseEnvironmentTypeID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DevCenter.V20250201.EnvironmentTypes.ProjectEnvironmentTypesGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r DevCenterProjectEnvironmentTypeTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center_project_environment_type" "test" {
  name                  = "acctest-et-%s"
  location              = azurerm_resource_group.test.location
  dev_center_project_id = azurerm_dev_center_project.test.id
  deployment_target_id  = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"

  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomString)
}

func (r DevCenterProjectEnvironmentTypeTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dev_center_project_environment_type" "import" {
  name                  = azurerm_dev_center_project_environment_type.test.name
  location              = azurerm_dev_center_project_environment_type.test.location
  dev_center_project_id = azurerm_dev_center_project_environment_type.test.dev_center_project_id
  deployment_target_id  = azurerm_dev_center_project_environment_type.test.deployment_target_id

  identity {
    type = "SystemAssigned"
  }
}
`, r.basic(data))
}

func (r DevCenterProjectEnvironmentTypeTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

data "azurerm_role_definition" "test" {
  name = "Owner"
}

resource "azurerm_dev_center_project_environment_type" "test" {
  name                  = "acctest-et-%s"
  location              = azurerm_resource_group.test.location
  dev_center_project_id = azurerm_dev_center_project.test.id
  deployment_target_id  = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  creator_role_assignment_roles = [split("/", data.azurerm_role_definition.test.id)[length(split("/", data.azurerm_role_definition.test.id)) - 1]]

  user_role_assignment {
    user_id = azurerm_user_assigned_identity.test.principal_id
    roles   = [split("/", data.azurerm_role_definition.test.id)[length(split("/", data.azurerm_role_definition.test.id)) - 1]]
  }

  tags = {
    Env = "Test"
  }
}
`, r.template(data), data.RandomString, data.RandomString)
}

func (r DevCenterProjectEnvironmentTypeTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctestuai2%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

data "azurerm_role_definition" "test" {
  name = "Owner"
}

data "azurerm_role_definition" "test2" {
  name = "Contributor"
}

resource "azurerm_dev_center_project_environment_type" "test" {
  name                  = "acctest-et-%s"
  location              = azurerm_resource_group.test.location
  dev_center_project_id = azurerm_dev_center_project.test.id
  deployment_target_id  = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id, azurerm_user_assigned_identity.test2.id]
  }

  creator_role_assignment_roles = [split("/", data.azurerm_role_definition.test.id)[length(split("/", data.azurerm_role_definition.test.id)) - 1], split("/", data.azurerm_role_definition.test2.id)[length(split("/", data.azurerm_role_definition.test2.id)) - 1]]

  user_role_assignment {
    user_id = azurerm_user_assigned_identity.test.principal_id
    roles   = [split("/", data.azurerm_role_definition.test.id)[length(split("/", data.azurerm_role_definition.test.id)) - 1], split("/", data.azurerm_role_definition.test2.id)[length(split("/", data.azurerm_role_definition.test2.id)) - 1]]
  }

  user_role_assignment {
    user_id = azurerm_user_assigned_identity.test2.principal_id
    roles   = [split("/", data.azurerm_role_definition.test2.id)[length(split("/", data.azurerm_role_definition.test2.id)) - 1]]
  }

  tags = {
    Env = "Test2"
  }
}
`, r.template(data), data.RandomString, data.RandomString, data.RandomString)
}

func (r DevCenterProjectEnvironmentTypeTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-dcpet-%d"
  location = "%s"
}

resource "azurerm_dev_center" "test" {
  name                = "acctest-dc-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_dev_center_environment_type" "test" {
  name          = "acctest-et-%s"
  dev_center_id = azurerm_dev_center.test.id
}

resource "azurerm_dev_center_project" "test" {
  name                = "acctest-dcp-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  dev_center_id       = azurerm_dev_center.test.id

  depends_on = [azurerm_dev_center_environment_type.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString)
}
