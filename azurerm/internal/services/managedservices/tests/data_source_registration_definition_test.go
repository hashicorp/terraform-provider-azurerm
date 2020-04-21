package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMRegistrationDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_registration_definition", "test")
	// Multiple tenants are needed to test this resource. 
	// Second tenant ID needs to be set as a enviornment variable ARM_SECOND_TENANT_ID.
	secondTenantID := os.Getenv("ARM_SECOND_TENANT_ID")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRegistrationDefinition_basic(id, secondTenantID, data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
					resource.TestMatchResourceAttr(data.ResourceName, "registration_definition_id", validate.UUIDRegExp),
					resource.TestCheckResourceAttr(data.ResourceName, "registration_definition_name", fmt.Sprintf("acctestrd-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Acceptance Test Registration Definition"),
					resource.TestMatchResourceAttr(data.ResourceName, "managed_by_tenant_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "authorization.0.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "authorization.0.role_definition_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func testAccDataSourceRegistrationDefinition_basic(id string, secondTenantID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

data "azurerm_client_config" "test" {
}

resource "azurerm_registration_definition" "test" {
  registration_definition_id 	= "%s"
  registration_definition_name  = "acctestrd-%d"
  description				 	= "Acceptance Test Registration Definition"
  scope                      	= data.azurerm_subscription.primary.id
  managed_by_tenant_id       	= "%s"

  authorization {
	principal_id        = data.azurerm_client_config.test.object_id
	role_definition_id  = "b24988ac-6180-42a0-ab88-20f7382dd24c"
  }
}

data "azurerm_registration_definition" "test" {
  registration_definition_id = azurerm_registration_definition.test.registration_definition_id
  scope              		 = data.azurerm_subscription.primary.id
}
`, id, data.RandomInteger, secondTenantID)
}
