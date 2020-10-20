package tests

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/parse"
)

func TestAccAzureRMManagementGroupSubscriptionAssociation(t *testing.T) {
	// managementGroupData := acceptance.BuildTestData(t, "azurerm_management_group", "test")
	groupData := acceptance.BuildTestData(t, "azurerm_management_group", "test")
	associationData := acceptance.BuildTestData(t, "azurerm_management_group_subscription_association", "test")
	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMManagementGroupSubscriptionAssociation_managementGroupOnly(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementGroupSubscriptionAssociationDoesNotExist(groupData.ResourceName, subscriptionId),
				),
			},
			{
				Config: testAzureRMManagementGroupSubscriptionAssociation_associatedSubscription(subscriptionId),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementGroupSubscriptionAssociationExists(associationData.ResourceName),
				),
			},
			{
				Config: testAzureRMManagementGroupSubscriptionAssociation_managementGroupOnly(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementGroupSubscriptionAssociationDoesNotExist(groupData.ResourceName, subscriptionId),
				),
			},
		},
	})
}

func testCheckAzureRMManagementGroupSubscriptionAssociationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		groupsClient := acceptance.AzureProvider.Meta().(*clients.Client).ManagementGroups.GroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.ManagementGroupSubscriptionAssociationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		recurse := false // Only want immediate children
		resp, err := groupsClient.Get(ctx, id.ManagementGroupName, "children", &recurse, "", "no-cache")
		if err != nil {
			return err
		}

		props := resp.Properties
		if props == nil {
			return fmt.Errorf("properties was nil for Management Group %q", id.ManagementGroupName)
		}

		children := props.Children
		if children == nil {
			return fmt.Errorf("Management Group %q has no children", id.ManagementGroupName)
		}

		exists := false
		if subscriptionIds := *children; subscriptionIds != nil {
			for _, subscriptionId := range subscriptionIds {
				if strings.EqualFold(id.SubscriptionID, *subscriptionId.Name) {
					exists = true
				}
			}
		}

		if !exists {
			return fmt.Errorf("Did not find %q associated to management group %q", id.SubscriptionID, id.ManagementGroupName)
		}

		return nil
	}
}

func testCheckAzureRMManagementGroupSubscriptionAssociationDoesNotExist(resourceName string, subscriptionIdToCheck string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[INFO] testCheckAzureRMManagementGroupSubscriptionAssociationDoesNotExist")

		groupsClient := acceptance.AzureProvider.Meta().(*clients.Client).ManagementGroups.GroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.ManagementGroupID(rs.Primary.ID)
		if err != nil {
			return err
		}

		recurse := false // Only want immediate children
		resp, err := groupsClient.Get(ctx, id.Name, "children", &recurse, "", "no-cache")
		if err != nil {
			return err
		}

		props := resp.Properties
		if props == nil {
			return fmt.Errorf("properties was nil for Management Group %q", id.Name)
		}

		children := props.Children
		if children == nil {
			return nil
		}

		log.Printf("[INFO] Management Group Name: %s", id.Name)
		log.Printf("[INFO] Children: %+v", *children)

		exists := false
		if subscriptionIds := *children; subscriptionIds != nil {
			for _, subscriptionId := range subscriptionIds {
				log.Printf("[INFO] Checking if %s is %s", *subscriptionId.Name, subscriptionIdToCheck)
				if strings.EqualFold(subscriptionIdToCheck, *subscriptionId.Name) {
					exists = true
				}
			}
		}

		if exists {
			return fmt.Errorf("Found %q associated to management group %q", subscriptionIdToCheck, id.Name)
		}

		return nil
	}
}

func testAzureRMManagementGroupSubscriptionAssociation_managementGroupOnly() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
	lifecycle {
    ignore_changes = [subscription_ids, ]
  }
}
`
}

func testAzureRMManagementGroupSubscriptionAssociation_associatedSubscription(subscriptionId string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  lifecycle {
    ignore_changes = [subscription_ids, ]
  }
}

resource "azurerm_management_group_subscription_association" "test" {
  management_group_id = azurerm_management_group.test.id
  subscription_id     = "%s"
}
`, subscriptionId)
}
