package azurerm

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMActiveDirectoryApplication_importBasic(t *testing.T) {
	resourceName := "azurerm_azuread_application.test"

	id := uuid.New().String()
	config := testAccAzureRMActiveDirectoryApplication_basic(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMActiveDirectoryApplication_importComplete(t *testing.T) {
	resourceName := "azurerm_azuread_application.test"

	id := uuid.New().String()
	config := testAccAzureRMActiveDirectoryApplication_complete(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
