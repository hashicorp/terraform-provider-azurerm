package resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type TemplateSpecResource struct{}

func TestAccTemplateSpec_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_template_spec", "test")
	r := TemplateSpecResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccTemplateSpec_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_template_spec", "test")
	r := TemplateSpecResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccTemplateSpec_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_template_spec", "test")
	r := TemplateSpecResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccTemplateSpec_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_template_spec", "test")
	r := TemplateSpecResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r TemplateSpecResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	clusterClient := client.Resource.TemplateSpecClient
	id, err := parse.TemplateSpecID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clusterClient.Get(ctx, id.ResourceGroup, id.Name, "versions")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving Template Spec %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (r TemplateSpecResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_template_spec" "test" {
  name                = "acctest-TemplateSpec-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, r.template(data), data.RandomInteger)
}

func (r TemplateSpecResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_template_spec" "import" {
  name                = azurerm_template_spec.test.name
  resource_group_name = azurerm_template_spec.test.resource_group_name
  location            = azurerm_template_spec.test.location
}
`, config)
}

func (r TemplateSpecResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_template_spec" "test" {
  name                = "acctest-TemplateSpec-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  description         = "Test Description"
  display_name        = "Test Display Name"

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r TemplateSpecResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_template_spec" "test" {
  name                = "acctest-TemplateSpec-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  description         = "Test Description"
  display_name        = "Test Display Name"

  tags = {
    ENV = "Test2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r TemplateSpecResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-templatespec-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
