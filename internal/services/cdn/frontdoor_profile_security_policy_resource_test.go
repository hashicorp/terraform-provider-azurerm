package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/securitypolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FrontdoorProfileSecurityPolicyResource struct{}

func TestAccFrontdoorProfileSecurityPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_profile_security_policy", "test")
	r := FrontdoorProfileSecurityPolicyResource{}
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

func TestAccFrontdoorProfileSecurityPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_profile_security_policy", "test")
	r := FrontdoorProfileSecurityPolicyResource{}
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

func TestAccFrontdoorProfileSecurityPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_profile_security_policy", "test")
	r := FrontdoorProfileSecurityPolicyResource{}
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

func TestAccFrontdoorProfileSecurityPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_profile_security_policy", "test")
	r := FrontdoorProfileSecurityPolicyResource{}
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

func (r FrontdoorProfileSecurityPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := securitypolicies.ParseSecurityPoliciesID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontdoorProfileSecurityPoliciesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r FrontdoorProfileSecurityPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-cdn-%d"
  location = "%s"
}
resource "azurerm_frontdoor_profile_profile" "test" {
  name                = "acctest-c-%d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r FrontdoorProfileSecurityPolicyResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_frontdoor_profile_security_policy" "test" {
  name           = "acctest-c-%d"
  cdn_profile_id = azurerm_frontdoor_profile_profile.test.id
  parameters {
    type = ""
  }
}
`, template, data.RandomInteger)
}

func (r FrontdoorProfileSecurityPolicyResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_profile_security_policy" "import" {
  name           = azurerm_frontdoor_profile_security_policy.test.name
  cdn_profile_id = azurerm_frontdoor_profile_profile.test.id
  parameters {
    type = ""
  }
}
`, config)
}

func (r FrontdoorProfileSecurityPolicyResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_profile_security_policy" "test" {
  name           = "acctest-c-%d"
  cdn_profile_id = azurerm_frontdoor_profile_profile.test.id
  parameters {
    type = ""
  }
}
`, template, data.RandomInteger)
}

func (r FrontdoorProfileSecurityPolicyResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_profile_security_policy" "test" {
  name           = "acctest-c-%d"
  cdn_profile_id = azurerm_frontdoor_profile_profile.test.id
  parameters {
    type = ""
  }
}
`, template, data.RandomInteger)
}
