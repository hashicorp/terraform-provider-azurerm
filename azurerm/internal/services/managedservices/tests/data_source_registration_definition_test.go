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
	// Multiple tenants are needed to test this resource.
	// Second tenant ID needs to be set as a environment variable ARM_TENANT_ID_ALT.
	// ObjectId for user, usergroup or service principal in second tenant needs to be set as a environment variable ARM_PRINCIPAL_ID_ALT_TENANT.
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	data := acceptance.BuildTestData(t, "data.azurerm_registration_definition", "test")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRegistrationDefinition_basic(id, secondTenantID, principalID, data),
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

func testAccDataSourceRegistrationDefinition_basic(id string, secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_registration_definition" "test" {
  registration_definition_id   = "%s"
  registration_definition_name = "acctestrd-%d"
  description                  = "Acceptance Test Registration Definition"
  managed_by_tenant_id         = "%s"

  authorization {
    principal_id       = "%s"
    role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
  }
}

data "azurerm_registration_definition" "test" {
  registration_definition_id = azurerm_registration_definition.test.registration_definition_id
}
`, id, data.RandomInteger, secondTenantID, principalID)
}
