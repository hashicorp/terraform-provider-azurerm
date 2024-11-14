package cognitive_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2024-10-01/raipolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type RaiPolicyTestResource struct{}

func TestCognitiveAccountRaiPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_rai_policy", "test")
	r := RaiPolicyTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.ApplyStep(r.basicConfig, r),
		data.ImportStep(),
	})
}

func TestCognitiveAccountRaiPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_rai_policy", "test")
	r := RaiPolicyTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.ApplyStep(r.basicConfig, r),
		data.RequiresImportErrorStep(r.requiresImportConfig),
	})
}

func TestCognitiveAccountRaiPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_rai_policy", "test")
	r := RaiPolicyTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.ApplyStep(r.completeConfig, r),
		data.ImportStep(),
		data.ApplyStep(r.basicConfig, r),
		data.ImportStep(),
		data.ApplyStep(r.completeConfig, r),
		data.ImportStep(),
	})
}

func (r RaiPolicyTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := raipolicies.ParseRaiPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Cognitive.RaiPoliciesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r RaiPolicyTestResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "OpenAI"
  sku_name            = "S0"
}

resource "azurerm_cognitive_account_rai_policy" "test" {
  name                 = "acctestraip-%s"
  cognitive_account_id = azurerm_cognitive_account.test.id
  base_policy_name     = "Microsoft.Default"
  content_filter {
    name               = "Hate"
    filter_enabled     = true
    block_enabled      = true
    severity_threshold = "High"
    source             = "Prompt"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (r RaiPolicyTestResource) completeConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "OpenAI"
  sku_name            = "S0"
}

resource "azurerm_cognitive_account_rai_policy" "test" {
  name                 = "acctestraip-%s"
  cognitive_account_id = azurerm_cognitive_account.test.id
  base_policy_name     = "Microsoft.Default"
  content_filter {
    name               = "Hate"
    filter_enabled     = true
    block_enabled      = true
    severity_threshold = "High"
    source             = "Prompt"
  }
  mode = "Asynchronous_filter"
  tags = {
    Environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (r RaiPolicyTestResource) requiresImportConfig(data acceptance.TestData) string {
	template := r.basicConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_rai_policy" "import" {
  name                 = azurerm_cognitive_account_rai_policy.test.name
  cognitive_account_id = azurerm_cognitive_account.test.id
  base_policy_name     = azurerm_cognitive_account_rai_policy.test.base_policy_name
  content_filter {
    name               = "Hate"
    filter_enabled     = true
    block_enabled      = true
    severity_threshold = "High"
    source             = "Prompt"
  }
}
`, template)
}
