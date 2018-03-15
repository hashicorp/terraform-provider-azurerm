package azurerm

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMServicePrincipal_importSimple(t *testing.T) {
	resourceName := "azurerm_service_principal.test"

	id := uuid.New().String()
	config := testAccAzureRMServicePrincipal_simple(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServicePrincipalDestroy,
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
