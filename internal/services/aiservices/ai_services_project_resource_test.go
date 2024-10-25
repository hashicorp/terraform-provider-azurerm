package aiservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AIServicesProject struct{}

func TestAccAIServicesProject_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_services_project", "test")
	r := AIServicesProject{}

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

func TestAccAIServicesProject_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_services_project", "test")
	r := AIServicesProject{}

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

func TestAccAIServicesProject_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_services_project", "test")
	r := AIServicesProject{}

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

func TestAccAIServicesProject_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_services_project", "test")
	r := AIServicesProject{}

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

func (AIServicesProject) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := workspaces.ParseWorkspaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MachineLearning.Workspaces.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r AIServicesProject) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_ai_services_project" "test" {
  name               = "acctestaip-%[2]d"
  location           = azurerm_ai_services_hub.test.location
  ai_services_hub_id = azurerm_ai_services_hub.test.id

  identity {
    type = "SystemAssigned"
  }
}
`, AIServicesHub{}.basic(data), data.RandomInteger)
}

func (r AIServicesProject) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_ai_services_project" "test" {
  name               = "acctestaip-%[2]d"
  location           = azurerm_ai_services_hub.test.location
  ai_services_hub_id = azurerm_ai_services_hub.test.id

  image_build_compute_name     = "projectbuild"
  description                  = "AI Project created by Terraform"
  friendly_name                = "AI Project"
  high_business_impact_enabled = false

  identity {
    type = "SystemAssigned"
  }

  tags = {
    model = "regression"
  }
}
`, AIServicesHub{}.complete(data), data.RandomInteger)
}

func (r AIServicesProject) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test_project" {
  name                = "acctestuaip-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_ai_services_project" "test" {
  name               = "acctestaip-%[2]d"
  location           = azurerm_ai_services_hub.test.location
  ai_services_hub_id = azurerm_ai_services_hub.test.id

  image_build_compute_name     = "projectbuildupdate"
  description                  = "AI Project updated by Terraform"
  friendly_name                = "AI Project for OS models"
  high_business_impact_enabled = false

  identity {
    type = "SystemAssigned"
  }

  tags = {
    model = "regression"
    env   = "test"
  }
}
`, AIServicesHub{}.complete(data), data.RandomInteger)
}

func (AIServicesProject) requiresImport(data acceptance.TestData) string {
	template := AIServicesProject{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_ai_services_project" "import" {
  name               = azurerm_ai_services_project.test.name
  location           = azurerm_ai_services_project.test.location
  ai_services_hub_id = azurerm_ai_services_project.test.ai_services_hub_id

  identity {
    type = "SystemAssigned"
  }
}
`, template)
}
