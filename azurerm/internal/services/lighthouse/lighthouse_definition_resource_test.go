package lighthouse_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/lighthouse/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LighthouseDefinitionResource struct {
}

func TestAccLighthouseDefinition_basic(t *testing.T) {
	// Multiple tenants are needed to test this resource.
	// Second tenant ID needs to be set as a environment variable ARM_TENANT_ID_ALT.
	// ObjectId for user, usergroup or service principal from second Tenant needs to be set as a environment variable ARM_PRINCIPAL_ID_ALT_TENANT.
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	if secondTenantID == "" || principalID == "" {
		t.Skip("Skipping as ARM_TENANT_ID_ALT and/or ARM_PRINCIPAL_ID_ALT_TENANT are not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_lighthouse_definition", "test")
	r := LighthouseDefinitionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(uuid.New().String(), secondTenantID, principalID, data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				resource.TestMatchResourceAttr(data.ResourceName, "lighthouse_definition_id", validate.UUIDRegExp),
				check.That(data.ResourceName).Key("authorization.0.principal_display_name").HasValue("Tier 1 Support"),
			),
		},
	})
}

func TestAccLighthouseDefinition_requiresImport(t *testing.T) {
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	if secondTenantID == "" || principalID == "" {
		t.Skip("Skipping as ARM_TENANT_ID_ALT and/or ARM_PRINCIPAL_ID_ALT_TENANT are not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_lighthouse_definition", "test")
	r := LighthouseDefinitionResource{}
	id := uuid.New().String()

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(id, secondTenantID, principalID, data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				resource.TestMatchResourceAttr(data.ResourceName, "lighthouse_definition_id", validate.UUIDRegExp),
			),
		},
		{
			Config:      r.requiresImport(id, secondTenantID, principalID, data),
			ExpectError: acceptance.RequiresImportError("azurerm_lighthouse_definition"),
		},
	})
}

func TestAccLighthouseDefinition_complete(t *testing.T) {
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	if secondTenantID == "" || principalID == "" {
		t.Skip("Skipping as ARM_TENANT_ID_ALT and/or ARM_PRINCIPAL_ID_ALT_TENANT are not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_lighthouse_definition", "test")
	r := LighthouseDefinitionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(uuid.New().String(), secondTenantID, principalID, data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				resource.TestMatchResourceAttr(data.ResourceName, "lighthouse_definition_id", validate.UUIDRegExp),
				check.That(data.ResourceName).Key("description").HasValue("Acceptance Test Lighthouse Definition"),
			),
		},
		data.ImportStep("lighthouse_definition_id"),
	})
}

func TestAccLighthouseDefinition_update(t *testing.T) {
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	if secondTenantID == "" || principalID == "" {
		t.Skip("Skipping as ARM_TENANT_ID_ALT and/or ARM_PRINCIPAL_ID_ALT_TENANT are not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_lighthouse_definition", "test")
	r := LighthouseDefinitionResource{}
	id := uuid.New().String()

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(id, secondTenantID, principalID, data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				resource.TestMatchResourceAttr(data.ResourceName, "lighthouse_definition_id", validate.UUIDRegExp),
			),
		},
		{
			Config: r.complete(id, secondTenantID, principalID, data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				resource.TestMatchResourceAttr(data.ResourceName, "lighthouse_definition_id", validate.UUIDRegExp),
				check.That(data.ResourceName).Key("description").HasValue("Acceptance Test Lighthouse Definition"),
			),
		},
		// multiple DelegatedRoleDefinitionIds
		{
			Config: r.updateDelegatedRoleDefinitionIds(id, secondTenantID, principalID, data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.complete(id, secondTenantID, principalID, data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basic(id, secondTenantID, principalID, data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccLighthouseDefinition_emptyID(t *testing.T) {
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	if secondTenantID == "" || principalID == "" {
		t.Skip("Skipping as ARM_TENANT_ID_ALT and/or ARM_PRINCIPAL_ID_ALT_TENANT are not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_lighthouse_definition", "test")
	r := LighthouseDefinitionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.emptyId(secondTenantID, principalID, data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("lighthouse_definition_id").Exists(),
			),
		},
	})
}

func (LighthouseDefinitionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.LighthouseDefinitionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Lighthouse.DefinitionsClient.Get(ctx, id.Scope, id.LighthouseDefinitionID)
	if err != nil {
		return nil, fmt.Errorf("retrieving Lighthouse Definition %q (Scope: %q) does not exist", id.LighthouseDefinitionID, id.Scope)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (LighthouseDefinitionResource) basic(id string, secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_role_definition" "contributor" {
  role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
}

data "azurerm_subscription" "test" {}

resource "azurerm_lighthouse_definition" "test" {
  lighthouse_definition_id = "%s"
  name                     = "acctest-LD-%d"
  managing_tenant_id       = "%s"
  scope                    = data.azurerm_subscription.test.id

  authorization {
    principal_id           = "%s"
    role_definition_id     = data.azurerm_role_definition.contributor.role_definition_id
    principal_display_name = "Tier 1 Support"
  }
}
`, id, data.RandomInteger, secondTenantID, principalID)
}

func (r LighthouseDefinitionResource) requiresImport(id string, secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lighthouse_definition" "import" {
  name                     = azurerm_lighthouse_definition.test.name
  lighthouse_definition_id = azurerm_lighthouse_definition.test.lighthouse_definition_id
  managing_tenant_id       = azurerm_lighthouse_definition.test.managing_tenant_id
  scope                    = azurerm_lighthouse_definition.test.scope
  authorization {
    principal_id       = azurerm_lighthouse_definition.test.managing_tenant_id
    role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
  }
}
`, r.basic(id, secondTenantID, principalID, data))
}

func (LighthouseDefinitionResource) complete(id string, secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_role_definition" "user_access_administrator" {
  role_definition_id = "18d7d88d-d35e-4fb5-a5c3-7773c20a72d9"
}

data "azurerm_role_definition" "contributor" {
  role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
}

data "azurerm_subscription" "test" {}

resource "azurerm_lighthouse_definition" "test" {
  lighthouse_definition_id = "%s"
  name                     = "acctest-LD-%d"
  description              = "Acceptance Test Lighthouse Definition"
  managing_tenant_id       = "%s"
  scope                    = data.azurerm_subscription.test.id

  authorization {
    principal_id                  = "%s"
    role_definition_id            = data.azurerm_role_definition.user_access_administrator.role_definition_id
    principal_display_name        = "Tier 2 Support"
    delegated_role_definition_ids = [data.azurerm_role_definition.contributor.role_definition_id]
  }
}
`, id, data.RandomInteger, secondTenantID, principalID)
}

func (LighthouseDefinitionResource) updateDelegatedRoleDefinitionIds(id string, secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_role_definition" "user_access_administrator" {
  role_definition_id = "18d7d88d-d35e-4fb5-a5c3-7773c20a72d9"
}

data "azurerm_role_definition" "contributor" {
  role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
}

data "azurerm_role_definition" "reader" {
  role_definition_id = "acdd72a7-3385-48ef-bd42-f606fba81ae7"
}

data "azurerm_subscription" "test" {}

resource "azurerm_lighthouse_definition" "test" {
  lighthouse_definition_id = "%s"
  name                     = "acctest-LD-%d"
  description              = "Acceptance Test Lighthouse Definition"
  managing_tenant_id       = "%s"
  scope                    = data.azurerm_subscription.test.id

  authorization {
    principal_id           = "%s"
    role_definition_id     = data.azurerm_role_definition.user_access_administrator.role_definition_id
    principal_display_name = "Tier 2 Support"
    delegated_role_definition_ids = [
      data.azurerm_role_definition.contributor.role_definition_id,
      data.azurerm_role_definition.reader.role_definition_id,
    ]
  }
}
`, id, data.RandomInteger, secondTenantID, principalID)
}

func (LighthouseDefinitionResource) emptyId(secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_role_definition" "contributor" {
  role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
}

data "azurerm_subscription" "test" {}

resource "azurerm_lighthouse_definition" "test" {
  name               = "acctest-LD-%d"
  description        = "Acceptance Test Lighthouse Definition"
  managing_tenant_id = "%s"
  scope              = data.azurerm_subscription.test.id

  authorization {
    principal_id       = "%s"
    role_definition_id = data.azurerm_role_definition.contributor.role_definition_id
  }
}
`, data.RandomInteger, secondTenantID, principalID)
}
