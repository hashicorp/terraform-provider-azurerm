package labservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/labservices/sdk/2021-10-01-preview/labplan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LabServicesPlanResource struct{}

func TestLabPlan_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LabServicesPlanResource{}

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

// Exists func

func (r LabServicesPlanResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := labplan.ParseLabPlanID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.LabServices.LabPlanClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Lab Plan %s: %+v", id, err)
	}
	if response.WasNotFound(resp.HttpResponse) {
		return utils.Bool(false), nil
	}
	return utils.Bool(true), nil
}

// Configs

func (r LabServicesPlanResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_lab_services_plan" "test" {
  name                = "acctestLSP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (LabServicesPlanResource) baseTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

`, data.RandomInteger, data.Locations.Primary)
}
