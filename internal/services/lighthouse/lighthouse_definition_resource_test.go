// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lighthouse_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2022-10-01/registrationdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LighthouseDefinitionResource struct{}

func TestAccLighthouseDefinition_basic(t *testing.T) {
	// Multiple tenants are needed to test this acceptance.
	// Second tenant ID needs to be set as a environment variable ARM_TENANT_ID_ALT.
	// ObjectId for user, usergroup or service principal from second Tenant needs to be set as a environment variable ARM_PRINCIPAL_ID_ALT_TENANT.
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	if secondTenantID == "" || principalID == "" {
		t.Skip("Skipping as ARM_TENANT_ID_ALT and/or ARM_PRINCIPAL_ID_ALT_TENANT are not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_lighthouse_definition", "test")
	r := LighthouseDefinitionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(uuid.New().String(), secondTenantID, principalID, data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				check.That(data.ResourceName).Key("lighthouse_definition_id").IsUUID(),
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(id, secondTenantID, principalID, data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				check.That(data.ResourceName).Key("lighthouse_definition_id").IsUUID(),
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(uuid.New().String(), secondTenantID, principalID, data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				check.That(data.ResourceName).Key("lighthouse_definition_id").IsUUID(),
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(id, secondTenantID, principalID, data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				check.That(data.ResourceName).Key("lighthouse_definition_id").IsUUID(),
			),
		},
		{
			Config: r.complete(id, secondTenantID, principalID, data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				check.That(data.ResourceName).Key("lighthouse_definition_id").IsUUID(),
				check.That(data.ResourceName).Key("description").HasValue("Acceptance Test Lighthouse Definition"),
			),
		},
		// multiple DelegatedRoleDefinitionIds
		{
			Config: r.updateDelegatedRoleDefinitionIds(id, secondTenantID, principalID, data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.complete(id, secondTenantID, principalID, data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basic(id, secondTenantID, principalID, data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.emptyId(secondTenantID, principalID, data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("lighthouse_definition_id").Exists(),
			),
		},
	})
}

func TestAccLighthouseDefinition_plan(t *testing.T) {
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	planName := os.Getenv("ARM_PLAN_NAME")
	planPublisher := os.Getenv("ARM_PLAN_PUBLISHER")
	planProduct := os.Getenv("ARM_PLAN_PRODUCT")
	planVersion := os.Getenv("ARM_PLAN_VERSION")
	if secondTenantID == "" || principalID == "" || planName == "" || planPublisher == "" || planProduct == "" || planVersion == "" {
		t.Skip("Skipping as ARM_TENANT_ID_ALT, ARM_PRINCIPAL_ID_ALT_TENANT, ARM_PLAN_NAME, ARM_PLAN_PUBLISHER, ARM_PLAN_PRODUCT or ARM_PLAN_VERSION are not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_lighthouse_definition", "test")
	r := LighthouseDefinitionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.plan(data, secondTenantID, principalID, planName, planPublisher, planProduct, planVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("lighthouse_definition_id").Exists(),
			),
		},
	})
}

func TestAccLighthouseDefinition_eligibleAuthorization(t *testing.T) {
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_USER_GROUP_ID_ALT_TENANT")
	if secondTenantID == "" || principalID == "" {
		t.Skip("Skipping as ARM_TENANT_ID_ALT and/or ARM_USER_GROUP_ID_ALT_TENANT are not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_lighthouse_definition", "test")
	r := LighthouseDefinitionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.eligibleAuthorization(uuid.New().String(), secondTenantID, principalID, data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("lighthouse_definition_id"),
	})
}

func (LighthouseDefinitionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := registrationdefinitions.ParseScopedRegistrationDefinitionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Lighthouse.DefinitionsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
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

func (LighthouseDefinitionResource) plan(data acceptance.TestData, secondTenantID, principalID, planName, planPublisher, planProduct, planVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_role_definition" "reader" {
  role_definition_id = "acdd72a7-3385-48ef-bd42-f606fba81ae7"
}

data "azurerm_subscription" "test" {}

resource "azurerm_lighthouse_definition" "test" {
  name               = "acctest-LD-%d"
  description        = "Acceptance Test Lighthouse Definition"
  managing_tenant_id = "%s"
  scope              = data.azurerm_subscription.test.id

  authorization {
    principal_id           = "%s"
    role_definition_id     = data.azurerm_role_definition.reader.role_definition_id
    principal_display_name = "Reader"
  }

  plan {
    name      = "%s"
    publisher = "%s"
    product   = "%s"
    version   = "%s"
  }
}
`, data.RandomInteger, secondTenantID, principalID, planName, planPublisher, planProduct, planVersion)
}

func (LighthouseDefinitionResource) eligibleAuthorization(id string, secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_role_definition" "contributor" {
  role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c" // Contributor role
}

data "azurerm_role_definition" "reader" {
  role_definition_id = "acdd72a7-3385-48ef-bd42-f606fba81ae7"
}

data "azurerm_subscription" "test" {}

resource "azurerm_lighthouse_definition" "test" {
  lighthouse_definition_id = "%s"
  name                     = "acctest-LD-%d"
  managing_tenant_id       = "%s"
  scope                    = data.azurerm_subscription.test.id

  authorization {
    principal_id           = "%s"
    role_definition_id     = data.azurerm_role_definition.reader.role_definition_id
    principal_display_name = "Reader"
  }

  eligible_authorization {
    principal_id           = "%s"
    role_definition_id     = data.azurerm_role_definition.contributor.role_definition_id
    principal_display_name = "Tier 1 Support"

    just_in_time_access_policy {
      multi_factor_auth_provider  = "Azure"
      maximum_activation_duration = "PT7H"

      approver {
        principal_id           = "%s"
        principal_display_name = "Tier 2 Support"
      }
    }
  }
}
`, id, data.RandomInteger, secondTenantID, principalID, principalID, principalID)
}
