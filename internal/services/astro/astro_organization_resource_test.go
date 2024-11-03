package astro_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/astronomer/2023-08-01/organizations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AstroOrganizationResource struct{}

func TestAccAstroOrganization_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_astro_organization", "test")
	r := AstroOrganizationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAstroOrganization_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_astro_organization", "test")
	r := AstroOrganizationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAstroOrganization_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_astro_organization", "test")
	r := AstroOrganizationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAstroOrganization_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_astro_organization", "test")
	r := AstroOrganizationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AstroOrganizationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := organizations.ParseOrganizationID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Astro.OrganizationsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("%s does not exist", id)
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r AstroOrganizationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

data "azurerm_subscription" "current" {}

`, data.RandomInteger, data.Locations.Primary)
}

func (r AstroOrganizationResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_astro_organization" "test" {
  name                = "acctest-ao-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  marketplace {
    subscription_id     = data.azurerm_subscription.current.subscription_id
    subscription_status = "Subscribed"
    offer {
      offer_id     = "example-offer-id"
      plan_id      = "example-plan-id"
      plan_name    = "example-plan-name"
      publisher_id = "example-publisher-id"
      term_id      = "example-term-id"
      term_unit    = "example-term-unit"
    }
  }
  partner_organization {
    single_sign_on {
      aad_domains          = ["mpliftrlogz20210811outlook.onmicrosoft.com"]
	}
  }
  user {
    email_address = "user@example.com"
    first_name    = "John"
    last_name     = "Doe"
    phone_number  = "+1234567890"
    principal_name = "john.doe@example.com"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AstroOrganizationResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_astro_organization" "import" {
  name                = azurerm_astro_organization.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  marketplace {
    subscription_id     = data.azurerm_subscription.current.subscription_id
    subscription_status = "Subscribed"
    offer {
      offer_id     = "example-offer-id"
      plan_id      = "example-plan-id"
      plan_name    = "example-plan-name"
      publisher_id = "example-publisher-id"
      term_id      = "example-term-id"
      term_unit    = "example-term-unit"
    }
  }
  partner_organization {
    single_sign_on {
      aad_domains          = ["mpliftrlogz20210811outlook.onmicrosoft.com"]
	}
  }
  user {
    email_address = "user@example.com"
    first_name    = "John"
    last_name     = "Doe"
    phone_number  = "+1234567890"
    principal_name = "john.doe@example.com"
  }
}
`, config, data.Locations.Primary)
}

func (r AstroOrganizationResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_astro_organization" "test" {
  name                = "acctest-ao-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  marketplace {
    subscription_id     = data.azurerm_subscription.current.subscription_id
    subscription_status = "Subscribed"
    offer {
      offer_id     = "example-offer-id"
      plan_id      = "example-plan-id"
      plan_name    = "example-plan-name"
      publisher_id = "example-publisher-id"
      term_id      = "example-term-id"
      term_unit    = "example-term-unit"
    }
  }
  partner_organization {
    organization_id   = "example-organization-id"
    organization_name = "example-organization-name"
    workspace_id      = "example-workspace-id"
    workspace_name    = "example-workspace-name"
    single_sign_on {
      enterprise_app_id    = "00000000-0000-0000-0000-000000000000"
      single_sign_on_state = "Enable"
      single_sign_on_url   = "https://example.com/sso"
      aad_domains          = ["mpliftrlogz20210811outlook.onmicrosoft.com"]
    }
  }
  user {
    email_address = "user@example.com"
    first_name    = "John"
    last_name     = "Doe"
    phone_number  = "+1234567890"
    principal_name = "john.doe@example.com"
  }
  tags = {
    environment = "production"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AstroOrganizationResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_astro_organization" "test" {
  name                = "acctest-ao-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  marketplace {
    subscription_id     = data.azurerm_subscription.current.subscription_id
    subscription_status = "Subscribed"
    offer {
      offer_id     = "example-offer-id"
      plan_id      = "example-plan-id"
      plan_name    = "example-plan-name"
      publisher_id = "example-publisher-id"
      term_id      = "example-term-id"
      term_unit    = "example-term-unit"
    }
  }
  partner_organization {
    organization_id   = "updated-organization-id"
    organization_name = "updated-organization-name"
    workspace_id      = "updated-workspace-id"
    workspace_name    = "updated-workspace-name"
    single_sign_on {
      enterprise_app_id    = "updated-enterprise-app-id"
      single_sign_on_state = "Enable"
      single_sign_on_url   = "https://updated-example.com/sso"
      aad_domains          = ["mpliftrlogz20210811outlook.onmicrosoft.com"]
    }
  }
  user {
    email_address = "updated_user@example.com"
    first_name    = "Jane"
    last_name     = "Smith"
    phone_number  = "+0987654321"
    principal_name = "jane.smith@example.com"
  }
  tags = {
    environment = "staging"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}
