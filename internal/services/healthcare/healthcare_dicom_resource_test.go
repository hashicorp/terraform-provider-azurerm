// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcare_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/dicomservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type HealthCareDicomResource struct{}

func TestAccHealthCareDicomResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_dicom_service", "test")
	r := HealthCareDicomResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccHealthCareDicomResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_dicom_service", "test")
	r := HealthCareDicomResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccHealthCareDicomResource_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_dicom_service", "test")
	r := HealthCareDicomResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccHealthCareDicomResource_updateUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_dicom_service", "test")
	r := HealthCareDicomResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
		{
			Config: r.userAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccHealthCareDicomResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_dicom_service", "test")
	r := HealthCareDicomResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (HealthCareDicomResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dicomservices.ParseDicomServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.HealthCare.HealthcareWorkspaceDicomServiceClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s, %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r HealthCareDicomResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthcare_dicom_service" "test" {
  name         = "dicom%d"
  workspace_id = azurerm_healthcare_workspace.test.id
  location     = "%s"
  depends_on   = [azurerm_healthcare_workspace.test]
}
`, r.template(data), data.RandomIntOfLength(10), data.Locations.Primary)
}

func (r HealthCareDicomResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthcare_dicom_service" "test" {
  name         = "dicom%d"
  workspace_id = azurerm_healthcare_workspace.test.id
  location     = "%s"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "None"
  }
  depends_on = [azurerm_healthcare_workspace.test]
}
`, r.template(data), data.RandomIntOfLength(10), data.Locations.Primary)
}

func (r HealthCareDicomResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthcare_dicom_service" "test" {
  name         = "dicom%d"
  workspace_id = azurerm_healthcare_workspace.test.id
  location     = "%s"

  tags = {
    environment = "Prod"
  }
  depends_on = [azurerm_healthcare_workspace.test]
}
`, r.template(data), data.RandomIntOfLength(10), data.Locations.Primary)
}

func (r HealthCareDicomResource) userAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_healthcare_dicom_service" "test" {
  name         = "dicom%d"
  workspace_id = azurerm_healthcare_workspace.test.id
  location     = "%s"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    environment = "None"
  }
  depends_on = [azurerm_healthcare_workspace.test, azurerm_user_assigned_identity.test]
}
`, r.template(data), data.RandomInteger, data.RandomIntOfLength(10), data.Locations.Primary)
}

func (r HealthCareDicomResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_healthcare_dicom_service" "import" {
  name         = azurerm_healthcare_dicom_service.test.name
  workspace_id = azurerm_healthcare_dicom_service.test.workspace_id
  location     = azurerm_healthcare_dicom_service.test.location
}
`, r.basic(data))
}

func (HealthCareDicomResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-dicom-%d"
  location = "%s"
}

resource "azurerm_healthcare_workspace" "test" {
  name                = "wk%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(10))
}
