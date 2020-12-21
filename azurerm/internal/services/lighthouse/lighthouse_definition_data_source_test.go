package lighthouse_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type LighthouseDefinitionDataSourceResource struct {
}

func TestAccLighthouseDefinitionDataSource_basic(t *testing.T) {
	// Multiple tenants are needed to test this resource.
	// Second tenant ID needs to be set as a environment variable ARM_TENANT_ID_ALT.
	// ObjectId for user, usergroup or service principal from second Tenant needs to be set as a environment variable ARM_PRINCIPAL_ID_ALT_TENANT.
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	data := acceptance.BuildTestData(t, "data.azurerm_lighthouse_definition", "test")
	r := LighthouseDefinitionDataSourceResource{}
	id := uuid.New().String()

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(id, secondTenantID, principalID, data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("scope").Exists(),
				resource.TestMatchResourceAttr(data.ResourceName, "lighthouse_definition_id", validate.UUIDRegExp),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctest-LD-%d")),
				check.That(data.ResourceName).Key("description").HasValue("Acceptance Test Lighthouse Definition"),
				resource.TestMatchResourceAttr(data.ResourceName, "managing_tenant_id", validate.UUIDRegExp),
				resource.TestMatchResourceAttr(data.ResourceName, "authorization.0.principal_id", validate.UUIDRegExp),
				resource.TestMatchResourceAttr(data.ResourceName, "authorization.0.role_definition_id", validate.UUIDRegExp),
			),
		},
	})
}

func (LighthouseDefinitionDataSourceResource) basic(id string, secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_role_definition" "contributor" {
  role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
}

resource "azurerm_lighthouse_definition" "test" {
  lighthouse_definition_id = "%s"
  name                     = "acctest-LD-%d"
  description              = "Acceptance Test Lighthouse Definition"
  managing_tenant_id       = "%s"

  authorization {
    principal_id       = "%s"
    role_definition_id = data.azurerm_role_definition.contributor.role_definition_id
  }
}

data "azurerm_lighthouse_definition" "test" {
  lighthouse_definition_id = azurerm_lighthouse_definition.test.lighthouse_definition_id
}
`, id, data.RandomInteger, secondTenantID, principalID)
}
