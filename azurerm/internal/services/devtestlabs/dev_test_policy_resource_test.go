package devtestlabs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DevTestPolicyResource struct {
}

func TestAccDevTestPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_policy", "test")
	r := DevTestPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDevTestPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_policy", "test")
	r := DevTestPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_dev_test_policy"),
		},
	})
}

func TestAccDevTestPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_policy", "test")
	r := DevTestPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Acceptance").HasValue("Test"),
			),
		},
		data.ImportStep(),
	})
}

func (DevTestPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	labName := id.Path["labs"]
	policySetName := id.Path["policysets"]
	name := id.Path["policies"]

	resp, err := clients.DevTestLabs.PoliciesClient.Get(ctx, id.ResourceGroup, labName, policySetName, name, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving DevTest Policy %q (Policy Set %q / Lab %q / Resource Group: %q) does not exist", name, policySetName, labName, id.ResourceGroup)
	}

	return utils.Bool(resp.PolicyProperties != nil), nil
}

func (DevTestPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dev_test_policy" "test" {
  name                = "LabVmCount"
  policy_set_name     = "default"
  lab_name            = azurerm_dev_test_lab.test.name
  resource_group_name = azurerm_resource_group.test.name
  threshold           = "999"
  evaluator_type      = "MaxValuePolicy"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r DevTestPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_policy" "import" {
  name                = azurerm_dev_test_policy.test.name
  policy_set_name     = "$[azurerm_dev_test_policy.test.policy_set_name}"
  lab_name            = azurerm_dev_test_policy.test.lab_name
  resource_group_name = azurerm_dev_test_policy.test.resource_group_name
  threshold           = "999"
  evaluator_type      = "MaxValuePolicy"
}
`, r.basic(data))
}

func (DevTestPolicyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dev_test_policy" "test" {
  name                = "LabVmCount"
  policy_set_name     = "default"
  lab_name            = azurerm_dev_test_lab.test.name
  resource_group_name = azurerm_resource_group.test.name
  threshold           = "999"
  evaluator_type      = "MaxValuePolicy"
  description         = "Aloha this is the max number of VM's'"

  tags = {
    "Acceptance" = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
