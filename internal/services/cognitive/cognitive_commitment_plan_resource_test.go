package cognitive_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/commitmentplans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CognitiveCommitmentPlanTestResource struct{}

func TestAccCognitiveCommitmentPlan_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_commitment_plan", "test")
	r := CognitiveCommitmentPlanTestResource{}
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

func TestAccCognitiveCommitmentPlan_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_commitment_plan", "test")
	r := CognitiveCommitmentPlanTestResource{}
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

func TestAccCognitiveCommitmentPlan_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_commitment_plan", "test")
	r := CognitiveCommitmentPlanTestResource{}
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

func TestAccCognitiveCommitmentPlan_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_commitment_plan", "test")
	r := CognitiveCommitmentPlanTestResource{}
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
	})
}

func (r CognitiveCommitmentPlanTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commitmentplans.ParseAccountCommitmentPlanID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cognitive.CommitmentPlansClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r CognitiveCommitmentPlanTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "SpeechServices"
  sku_name            = "S0"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CognitiveCommitmentPlanTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cognitive_commitment_plan" "test" {
  name                 = "acctest-ccp-%d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  hosting_model        = "Web"
  plan_type            = "STT"
  current_tier         = "T1"
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveCommitmentPlanTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_commitment_plan" "import" {
  name                 = azurerm_cognitive_commitment_plan.test.name
  cognitive_account_id = azurerm_cognitive_account.test.id
  hosting_model        = azurerm_cognitive_commitment_plan.test.hosting_model
  plan_type            = azurerm_cognitive_commitment_plan.test.plan_type
  current_tier         = azurerm_cognitive_commitment_plan.test.current_tier
}
`, r.basic(data))
}

func (r CognitiveCommitmentPlanTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cognitive_commitment_plan" "test" {
  name                 = "acctest-ccp-%d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auto_renew_enabled   = true
  hosting_model        = "Web"
  plan_type            = "STT"
  current_tier         = "T1"
  renewal_tier         = "T2"

  tags = {
    key = "value"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveCommitmentPlanTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cognitive_commitment_plan" "test" {
  name                 = "acctest-ccp-%d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auto_renew_enabled   = false
  hosting_model        = "Web"
  plan_type            = "STT"
  current_tier         = "T1"

  tags = {
    key2 = "value2"
  }
}
`, r.template(data), data.RandomInteger)
}
