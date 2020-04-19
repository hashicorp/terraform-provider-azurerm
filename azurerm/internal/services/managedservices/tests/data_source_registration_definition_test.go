package tests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMRegistrationDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_registration_definition", "test")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRegistrationDefinition_basic(id, data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "description"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "type"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "registrationDefinitionName"),
					resource.TestMatchResourceAttr(data.ResourceName, "managed_by_tenant_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "authorization.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "authorization.role_definition_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func testAccDataSourceRegistrationDefinition_basic(id string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

resource "azuread_application" "test" {
	name = "acctestspa-%d"
}
		
resource "azuread_service_principal" "test" {
	application_id = azuread_application.test.application_id
}
	  
data "azurerm_role_definition" "builtin" {
	name = "Contributor"
}

resource "azurerm_registration_definition" "test" {
  registration_definition_id = "%s"
  name                       = "acctestrd-%d"
  scope                      = data.azurerm_subscription.primary.id
  managed_by_tenant_id       = "70b0b1ee-ab4d-4c83-816a-2957afa9ea0b"

  authorization {
	principal_id        = azuread_service_principal.test.id
	role_definition_id  = data.azurerm_role_definition.builtin.name
  }
}

data "azurerm_registration_definition" "test" {
  registration_definition_id = azurerm_registration_definition.test.registration_definition_id
  scope              		 = data.azurerm_subscription.primary.id
}
`, data.RandomInteger, id, data.RandomInteger)
}
