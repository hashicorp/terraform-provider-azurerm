package automation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/python3package"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type Python3PackageResource struct{}

func (a Python3PackageResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := python3package.ParsePython3PackageID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Automation.Python3Package.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving Python3Package %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func TestAccPython3Package_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.Python3PackageResource{}.ResourceType(), "test")
	r := Python3PackageResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("content_uri", "content_version"),
	})
}

func TestAccPython3Package_update(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.Python3PackageResource{}.ResourceType(), "test")
	r := Python3PackageResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("content_uri", "content_version"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("content_uri", "content_version", "hash_algorithm", "hash_value"),
	})
}

func (a Python3PackageResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_automation_python3_package" "test" {
  name                    = "acctest-%[2]d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  content_uri             = "https://pypi.org/packages/source/r/requests/requests-2.31.0.tar.gz"
  content_version         = "2.31.0"
  tags = {
    key = "foo"
  }
}
`, a.template(data), data.RandomInteger)
}

func (a Python3PackageResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_automation_python3_package" "test" {
  name                    = "acctest-%[2]d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  content_uri             = "https://pypi.org/packages/source/r/requests/requests-2.31.0.tar.gz"
  content_version         = "2.31.0"
  hash_algorithm          = "sha256"
  hash_value              = "942c5a758f98d790eaed1a29cb6eefc7ffb0d1cf7af05c3d2791656dbd6ad1e1"
  tags = {
    key = "bar"
  }
}
`, a.template(data), data.RandomInteger)
}

func (a Python3PackageResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%[1]d"
  location = "%[2]s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}
`, data.RandomInteger, data.Locations.Primary)
}
