package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/blueprint/mgmt/2018-11-01-preview/blueprint"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/blueprint/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const (
	blueprintDefinitionNamePattern = "acctest-definition-%d"
	resourceGroupNamePattern       = "acctest-RG-blueprint-%d"
)

func TestAccAzureRMBlueprintAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_blueprint_assignment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBlueprintAssignmentAndDependenciesDestroy(data),
		Steps: []resource.TestStep{
			// create and publish a blueprint definition using go SDK, since the blueprint definition resource is not implemented in terraform yet
			{
				Config: testAccAzureRMBlueprintAssignment_subscription(),
				Check: resource.ComposeTestCheckFunc(
					createBlueprintDefinition(data, "data.azurerm_subscription.current"),
					publishBlueprintDefinition(data, "data.azurerm_subscription.current"),
				),
			},
			// create and test the blueprint assignment resource
			{
				Config: testAccAzureRMBlueprintAssignment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBlueprintAssignmentExists(data.ResourceName),
					// test the resources created by the blueprint
					testCheckAzureRMAssignedResourceGroupsExists(data),
				),
			},
			// clean up definitions and artifacts
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMBlueprintAssignment_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_blueprint_assignment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBlueprintAssignmentAndDependenciesDestroy(data),
		Steps: []resource.TestStep{
			// create and publish a blueprint definition using go SDK, since the blueprint definition resource is not implemented in terraform yet
			{
				Config: testAccAzureRMBlueprintAssignment_subscription(),
				Check: resource.ComposeTestCheckFunc(
					createBlueprintDefinition(data, "data.azurerm_subscription.current"),
					publishBlueprintDefinition(data, "data.azurerm_subscription.current"),
				),
			},
			// create and test the blueprint assignment resource
			{
				Config: testAccAzureRMBlueprintAssignment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBlueprintAssignmentExists(data.ResourceName),
					// test the resources created by the blueprint
					testCheckAzureRMAssignedResourceGroupsExists(data),
				),
			},
			// clean up definitions and artifacts
			data.RequiresImportErrorStep(testAccAzureRMBlueprintAssignment_requiresImport),
		},
	})
}

func testCheckAzureRMBlueprintAssignmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("bad: Blueprint Assignment not found: %s", resourceName)
		}

		id, err := parse.BlueprintAssignmentID(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Blueprint.AssignmentClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, id.ScopeId, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Blueprint Assignment %q (Scope %q) does not exist", id.Name, id.ScopeId)
			}
			return fmt.Errorf("bad: Get on Blueprint.AssignmentClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMAssignedResourceGroupsExists(data acceptance.TestData) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Resource.GroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resourceGroupName := fmt.Sprintf(resourceGroupNamePattern, data.RandomInteger)
		if resp, err := client.Get(ctx, resourceGroupName); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Resource Group %q does not exist", resourceGroupName)
			}
			return fmt.Errorf("bad: Get on GroupsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMBlueprintAssignmentAndDependenciesDestroy(data acceptance.TestData) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if err := testCheckAzureRMBlueprintAssignmentDestroy(s); err != nil {
			return err
		}

		// destroy the resource group that is created by the blueprint
		resourceGroupName := fmt.Sprintf(resourceGroupNamePattern, data.RandomInteger)
		if err := destroyAzureRMAssignedResourceGroups(resourceGroupName); err != nil {
			return err
		}

		// destroy the blueprint definition
		definitionName := fmt.Sprintf(blueprintDefinitionNamePattern, data.RandomInteger)
		if err := destroyAzureRMBlueprintAssignmentDependencies(definitionName); err != nil {
			return err
		}

		return nil
	}
}

func testCheckAzureRMBlueprintAssignmentDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Blueprint.AssignmentClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_blueprint_assignment" {
			continue
		}

		id, err := parse.BlueprintAssignmentID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ScopeId, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on Blueprint.AssignmentClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func destroyAzureRMBlueprintAssignmentDependencies(definitionName string) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Blueprint.BlueprintClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	subscriptionId := acceptance.AzureProvider.Meta().(*clients.Client).Account.SubscriptionId

	log.Printf("[DEBUG] Deleting the blueprint definition")
	scope := fmt.Sprintf("/subscriptions/%s", subscriptionId)
	_, err := client.Delete(ctx, scope, definitionName)
	if err != nil {
		return fmt.Errorf("bad: error deleting Blueprint Definition %q (Scope %q): %+v", definitionName, scope, err)
	}

	return nil
}

func destroyAzureRMAssignedResourceGroups(name string) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Resource.GroupsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	log.Printf("[DEBUG] Deleting the resource group created by the assignment")
	future, err := client.Delete(ctx, name)
	if err != nil {
		return fmt.Errorf("bad: error deleting Resource Group %q: %+v", name, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("bad: error when waiting for Resource Group %q to be deleted: %+v", name, err)
	}

	return nil
}

func createBlueprintDefinition(data acceptance.TestData, scopeSource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		definitionClient := acceptance.AzureProvider.Meta().(*clients.Client).Blueprint.BlueprintClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[scopeSource]
		if !ok {
			return fmt.Errorf("bad: Not found: %s", scopeSource)
		}

		log.Printf("[DEBUG] Creating a blueprint definition")
		scope := rs.Primary.ID
		definitionName := fmt.Sprintf(blueprintDefinitionNamePattern, data.RandomInteger)
		model := blueprint.Model{
			Properties: &blueprint.Properties{
				TargetScope: blueprint.Subscription,
				Parameters: map[string]*blueprint.ParameterDefinition{
					"tagName": {
						Type:         blueprint.String,
						DefaultValue: utils.String("ENV"),
						ParameterDefinitionMetadata: &blueprint.ParameterDefinitionMetadata{
							DisplayName: utils.String("Tag Name"),
							Description: utils.String("Tag name for each resource that gets created"),
						},
					},
					"tagValue": {
						Type:         blueprint.String,
						DefaultValue: utils.String("Acc-test"),
						ParameterDefinitionMetadata: &blueprint.ParameterDefinitionMetadata{
							DisplayName: utils.String("Tag Value"),
							Description: utils.String("Tag value for each resource that gets created"),
						},
					},
				},
				ResourceGroups: map[string]*blueprint.ResourceGroupDefinition{
					"ProdRG": {
						ParameterDefinitionMetadata: &blueprint.ParameterDefinitionMetadata{
							DisplayName: utils.String("Production resource group"),
						},
					},
				},
				DisplayName: utils.String("Common Policies"),
				Description: utils.String("A set of popular policies to apply to a subscription"),
			},
		}
		_, err := definitionClient.CreateOrUpdate(ctx, scope, definitionName, model)
		if err != nil {
			return fmt.Errorf("bad: error creating Blueprint Definition %q (Scope %q): %+v", definitionName, scope, err)
		}

		return nil
	}
}

func publishBlueprintDefinition(data acceptance.TestData, scopeSource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Blueprint.PublishClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[scopeSource]
		if !ok {
			return fmt.Errorf("bad: Not found: %s", scopeSource)
		}

		log.Printf("[DEBUG] Publish a blueprint definition")
		scope := rs.Primary.ID
		definitionName := fmt.Sprintf(blueprintDefinitionNamePattern, data.RandomInteger)
		_, err := client.Create(ctx, scope, definitionName, "v1", nil)
		if err != nil {
			return fmt.Errorf("bad: error publishing Blueprint Definition %q (Scope %q): %+v", definitionName, scope, err)
		}

		return nil
	}
}

func testAccAzureRMBlueprintAssignment_subscription() string {
	return `
provider "azurerm" {
	features {}
}

data "azurerm_subscription" "current" {}
`
}

func testAccAzureRMBlueprintAssignment_basic(data acceptance.TestData) string {
	resourceGroupName := fmt.Sprintf(resourceGroupNamePattern, data.RandomInteger)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_blueprint_assignment" "test" {
  name     = "acctest-blueprint-%[1]d"
  location = "%[2]s"
  scope    = data.azurerm_subscription.current.id

  blueprint_definition_id = "${data.azurerm_subscription.current.id}/providers/Microsoft.Blueprint/blueprints/acctest-definition-%[1]d"

  identity {
    type = "SystemAssigned"
  }

  resource_groups = <<GROUPS
        {
          "prodRG": {
            "name": "%[3]s",
            "location": "%[2]s"
          }
        }
  GROUPS

  parameter_values = <<VALUES
        {
          "tagName": {
            "value": "ENV"
          },
          "tagValue": {
            "value": "Acc-test"
          }
        }
  VALUES
}
`, data.RandomInteger, data.Locations.Primary, resourceGroupName)
}

func testAccAzureRMBlueprintAssignment_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMBlueprintAssignment_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_blueprint_assignment" "import" {
  name     = azurerm_blueprint_assignment.test.name
  location = azurerm_blueprint_assignment.test.location
  scope    = azurerm_blueprint_assignment.test.scope

  blueprint_definition_id = azurerm_blueprint_assignment.test.blueprint_definition_id

  identity {
    type = "SystemAssigned"
  }

  resource_groups = azurerm_blueprint_assignment.test.resource_groups

  parameter_values = azurerm_blueprint_assignment.test.parameter_values
}
`, template)
}
