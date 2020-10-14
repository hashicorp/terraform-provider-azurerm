package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMManagementGroupSubscriptionAssociation(t *testing.T) {
	managementGroupData := acceptance.BuildTestData(t, "azurerm_management_group", "test")
	managementGroupSubscriptionAssociationData := acceptance.BuildTestData(t, "azurerm_management_group_subscription_association", "test")
	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMManagementGroup_managementGroupOnly(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementGroupExists(managementGroupData.ResourceName),
					testCheckAzureRMManagementGroupSubscriptionAssociationDoesNotExist(managementGroupSubscriptionAssociationData.ResourceName),
					resource.TestCheckResourceAttr(managementGroupData.ResourceName, "subscription_ids.#", "0"),
				),
			},
			{
				Config: testAzureRMManagementGroup_associatedSubscription(subscriptionId),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementGroupExists(managementGroupData.ResourceName),
					testCheckAzureRMManagementGroupSubscriptionAssociationExists(managementGroupSubscriptionAssociationData.ResourceName),
					resource.TestCheckResourceAttr(managementGroupData.ResourceName, "subscription_ids.#", "1"),
				),
			},
			{
				Config: testAzureRMManagementGroup_managementGroupOnly(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementGroupExists(managementGroupData.ResourceName),
					testCheckAzureRMManagementGroupSubscriptionAssociationDoesNotExist(managementGroupSubscriptionAssociationData.ResourceName),
					resource.TestCheckResourceAttr(managementGroupData.ResourceName, "subscription_ids.#", "0"),
				),
			},
		},
	})
}

func testCheckAzureRMManagementGroupSubscriptionAssociationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// No GET for (*clients.Client).ManagementGroups.GroupsClient so just check Terraform state plus subscription_ids.# increment / decrement

		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		return nil
	}
}

func testCheckAzureRMManagementGroupSubscriptionAssociationDoesNotExist(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// No GET for (*clients.Client).ManagementGroups.GroupsClient so just check Terraform state plus subscription_ids.# increment / decrement

		_, ok := s.RootModule().Resources[resourceName]
		if ok {
			return fmt.Errorf("unexpectedly found: %s", resourceName)
		}

		return nil
	}
}

func testAzureRMManagementGroup_managementGroupOnly() string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {}

`)
}

func testAzureRMManagementGroup_associatedSubscription(subscriptionId string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {}

resource "azurerm_management_group_subscription_association" "test" {
  management_group_id = azurerm_management_group.test.id
  subscription_id     = "%s"
}

`, subscriptionId)
}
