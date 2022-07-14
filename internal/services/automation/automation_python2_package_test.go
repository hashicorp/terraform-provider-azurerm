package automation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type Python2PackageResource struct{}

func (a Python2PackageResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.Python2PackageID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Automation.Python2PackageClient.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Type %s: %+v", id, err)
	}
	return utils.Bool(resp.ModuleProperties != nil), nil
}

func (a Python2PackageResource) template(data acceptance.TestData) string {
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

func (a Python2PackageResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_automation_python2_package" "test" {
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  name                    = "acctest-py-%[2]d"
  tags = {
    foo = "bar"
  }
  content {
    uri            = "https://files.pythonhosted.org/packages/3e/93/02056aca45162f9fc275d1eaad12a2a07ef92375afb48eabddc4134b8315/azure_graphrbac-0.61.1-py2.py3-none-any.whl"
    hash_algorithm = "sha256"
    hash_value     = "7b4e0f05676acc912f2b33c71c328d9fb2e4dc8e70ebadc9d3de8ab08bf0b175"
    version        = "1.0.0.0"
  }
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a Python2PackageResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_automation_python2_package" "test" {
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  name                    = "acctest-py-%[2]d"
  tags = {
    foo = "bar2"
  }
  content {
    uri            = "https://files.pythonhosted.org/packages/3e/93/02056aca45162f9fc275d1eaad12a2a07ef92375afb48eabddc4134b8315/azure_graphrbac-0.61.1-py2.py3-none-any.whl"
    hash_algorithm = "sha256"
    hash_value     = "7b4e0f05676acc912f2b33c71c328d9fb2e4dc8e70ebadc9d3de8ab08bf0b175"
    version        = "1.0.0.0"
  }
}
`, a.template(data), data.RandomInteger)
}

func TestAccPython2Package_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.Python2PackageResource{}.ResourceType(), "test")
	r := Python2PackageResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.foo").HasValue("bar"),
			),
		},
		data.ImportStep("content"),
	})
}

func TestAccPython2Package_update(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.Python2PackageResource{}.ResourceType(), "test")
	r := Python2PackageResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.foo").HasValue("bar"),
			),
		},
		data.ImportStep("content"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.foo").HasValue("bar2"),
			),
		},
		data.ImportStep("content", "tags"),
	})
}
