package cognitive_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2024-10-01/raipolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CognitiveRaiPolicyTestResource struct{}

func TestAccCognitiveRaiPolicySequential(t *testing.T) {
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"rai_policy": {
			"basic":          TestAccCognitiveRaiPolicy_basic,
			"requiresImport": TestAccCognitiveRaiPolicy_requiresImport,
			"complete":       TestAccCognitiveRaiPolicy_complete,
			"update":         TestAccCognitiveRaiPolicy_update,
		},
	})
}

func TestAccCognitiveRaiPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_rai_policy", "test")
	testRaiPolicy := CognitiveRaiPolicyTestResource{}

	data.ResourceSequentialTest(t, testRaiPolicy, []acceptance.TestStep{
		{
			Config: testRaiPolicy.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(testRaiPolicy),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveRaiPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_rai_policy", "test")
	testRaiPolicy := CognitiveRaiPolicyTestResource{}

	data.ResourceSequentialTest(t, testRaiPolicy, []acceptance.TestStep{
		{
			Config: testRaiPolicy.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(testRaiPolicy),
			),
		},
		data.RequiresImportErrorStep(testRaiPolicy.requiresImport),
	})
}

func TestAccCognitiveRaiPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_rai_policy", "test")
	testRaiPolicy := CognitiveRaiPolicyTestResource{}

	data.ResourceSequentialTest(t, testRaiPolicy, []acceptance.TestStep{
		{
			Config: testRaiPolicy.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(testRaiPolicy),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveRaiPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_rai_policy", "test")
	testRaiPolicy := CognitiveRaiPolicyTestResource{}

	data.ResourceSequentialTest(t, testRaiPolicy, []acceptance.TestStep{
		{
			Config: testRaiPolicy.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(testRaiPolicy),
			),
		},
		data.ImportStep(),
		{
			Config: testRaiPolicy.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(testRaiPolicy),
			),
		},
		data.ImportStep(),
	})
}

func (r CognitiveRaiPolicyTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := raipolicies.ParseRaiPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cognitive.RaiPoliciesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r CognitiveRaiPolicyTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
		provider "azurerm" {
			features{}
		}

		resource "azurerm_resource_group" "test" {
			name     = "acctest-rg-%d"
			location = "%s"
		}
		
		resource "azurerm_cognitive_account" "test" {
			name                = "acctest-ca-%d"
			location            = azurerm_resource_group.test.location
			resource_group_name = azurerm_resource_group.test.name
			kind                = "OpenAI"
			sku_name            = "S0"
		}

		resource "azurerm_cognitive_deployment" "test" {
			name                 = "acctest-cd-%d"
			cognitive_account_id = azurerm_cognitive_account.test.id

			model {
				format  = "OpenAI"
				name    = "gpt-4o-mini"
				version = "2024-07-18"
			}
			sku {
				name = "Standard"
			}
		}
	`, data.RandomInteger, data.Locations.Secondary, data.RandomInteger, data.RandomInteger)
}

func (r CognitiveRaiPolicyTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
		%s
		resource "azurerm_cognitive_rai_policy" "test"{
			name                 = "acctest-crp-%d"
			cognitive_account_id = azurerm_cognitive_account.test.id
			mode                 = "Default"
		}
	`, template, data.RandomInteger)

}

func (r CognitiveRaiPolicyTestResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
		%s
		resource "azurerm_cognitive_rai_policy" "import"{
			name                 = azurerm_cognitive_rai_policy.test.name
			cognitive_account_id = azurerm_cognitive_account.test.id
			mode                 = "Default"
		}
	`, config)
}

func (r CognitiveRaiPolicyTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
		%s
		resource "azurerm_cognitive_rai_policy" "test"{
			name                 = "acctest-crp-%d"
			cognitive_account_id = azurerm_cognitive_account.test.id
			base_policy_name     = "Microsoft.Default"
			mode                 = "Default"
			type 				 = "SystemManaged"

			content_filters = [
				{
					name               = "Hate"
					blocking           = true
					enabled            = true
					severity_threshold = "High"
					source             = "Prompt"
				},
				{
					name               = "Hate"
					blocking           = true
					enabled            = true
					severity_threshold = "High"
					source             = "Completion"
				},
				{
					name               = "Sexual"
					blocking           = true
					enabled            = true
					severity_threshold = "Medium"
					source             = "Prompt"
				},
				{
					name               = "Sexual"
					blocking           = true
					enabled            = true
					severity_threshold = "High"
					source             = "Completion"
				},
				{
					name               = "Selfharm"
					blocking           = true
					enabled            = true
					severity_threshold = "High"
					source             = "Prompt"
				},
				{
					name               = "Selfharm"
					blocking           = true
					enabled            = true
					severity_threshold = "High"
					source             = "Completion"
				},
				{
					name               = "Violence"
					blocking           = true
					enabled            = true
					severity_threshold = "High"
					source             = "Prompt"
				},
				{
					name               = "Violence"
					blocking           = true
					enabled            = true
					severity_threshold = "High"
					source             = "Completion"
				},
				{
					name               = "Jailbreak"
					blocking           = true
					enabled            = true
					source             = "Prompt"
					severity_threshold = "Medium"
				},
				{
					name               = "Indirect Attack"
					blocking           = true
					enabled            = true
					source             = "Prompt"
					severity_threshold = "Medium"
				},
				{
					name               = "Protected Material Text"
					blocking           = true
					enabled            = true
					source             = "Completion"
					severity_threshold = "High"
				},
				{
					name               = "Protected Material Code"
					blocking           = true
					enabled            = true
					source             = "Completion"
					severity_threshold = "High"
				},
				{
					name               = "Profanity"
					blocking           = true
					enabled            = true
					source             = "Prompt"
					severity_threshold = "High"
				},
				{
					name               = "Profanity"
					blocking           = true
					enabled            = true
					source             = "Completion"
					severity_threshold = "High"
				},
			]
		}
	`, template, data.RandomInteger)
}
