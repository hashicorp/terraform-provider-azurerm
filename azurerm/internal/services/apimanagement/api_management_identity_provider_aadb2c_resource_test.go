package apimanagement_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

/*
Note that this resource requires a preexisting B2C tenant with matching policies already created.
Accordingly, these tests rely on additional environment variables to be set (and be valid):
* ARM_TEST_B2C_TENANT_ID      - the UUID of the B2C tenant
* ARM_TEST_B2C_TENANT_SLUG    - the first part of the *.onmicrosoft.com domain of the B2C tenant
* ARM_TEST_B2C_CLIENT_ID      - client ID of an application in the B2C domain with privileges to read directory data
* ARM_TEST_B2C_CLIENT_SECRET  - client secret for that application
*/

type ApiManagementIdentityProviderAADB2CResource struct {
}

func TestAccAzureRMApiManagementIdentityProviderAADB2C_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_aadb2c", "test")
	r := ApiManagementIdentityProviderAADB2CResource{}
	b2cConfig := testAccAzureRMApiManagementIdentityProviderAADB2C_getB2CConfig(t)
	env, err := acceptance.Environment()
	if err != nil {
		t.Fatalf("could not load Azure Environment: %+v", err)
	}
	apiDomain := env.APIManagementHostNameSuffix
	if apiDomain == "" {
		t.Fatalf("APIManagementHostNameSuffix was empty")
	}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, b2cConfig, apiDomain),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret"),
	})
}

func TestAccAzureRMApiManagementIdentityProviderAADB2C_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_aadb2c", "test")
	r := ApiManagementIdentityProviderAADB2CResource{}
	b2cConfig := testAccAzureRMApiManagementIdentityProviderAADB2C_getB2CConfig(t)
	apiDomain := acceptance.AzureProvider.Meta().(*clients.Client).Account.Environment.APIManagementHostNameSuffix

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, b2cConfig, apiDomain),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, b2cConfig, apiDomain),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
	})
}

func testAccAzureRMApiManagementIdentityProviderAADB2C_getB2CConfig(t *testing.T) map[string]string {
	config := map[string]string{
		"tenant_id":     "",
		"tenant_slug":   "",
		"client_id":     "",
		"client_secret": "",
	}

	for k := range config {
		e := fmt.Sprintf("ARM_TEST_B2C_%s", strings.ToUpper(k))
		if v := os.Getenv(e); v != "" {
			config[k] = v
			continue
		}
		vars := make([]string, 0, len(config))
		for k := range config {
			v := fmt.Sprintf("ARM_TEST_B2C_%s", strings.ToUpper(k))
			vars = append(vars, v)
		}
		t.Skip(fmt.Sprintf("Acceptance tests for resource `azurerm_api_management_identity_provider_aadb2c` skipped unless environment variables set: %s", strings.Join(vars, ", ")))
	}

	return config
}

func (ApiManagementIdentityProviderAADB2CResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.IdentityProviderID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.IdentityProviderClient.Get(ctx, id.ResourceGroup, id.ServiceName, apimanagement.IdentityProviderType(id.Name))
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagement Identity Provider AADB2C (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ApiManagementIdentityProviderAADB2CResource) basic(data acceptance.TestData, b2cConfig map[string]string, apiDomain string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {
  tenant_id     = "%[1]s"
  client_id     = "%[2]s"
  client_secret = "%[3]s"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-api-%[5]d"
  location = "%[6]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[5]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}

resource "azuread_application" "test" {
  name                       = "acctestAM-%[5]d"
  oauth2_allow_implicit_flow = true
  reply_urls                 = ["https://${azurerm_api_management.test.name}.developer.%[8]s/signin"]
}

resource "azuread_application_password" "test" {
  application_object_id = azuread_application.test.object_id
  end_date_relative     = "36h"
  value                 = "P@55w0rD!%[7]s"
}

resource "azurerm_api_management_identity_provider_aadb2c" "test" {
  resource_group_name    = azurerm_resource_group.test.name
  api_management_name    = azurerm_api_management.test.name
  client_id              = azuread_application.test.application_id
  client_secret          = "P@55w0rD!%[7]s"
  allowed_tenant         = "%[4]s.onmicrosoft.com"
  signin_tenant          = "%[4]s.onmicrosoft.com"
  authority              = "%[4]s.b2clogin.com"
  signin_policy          = "B2C_1_Login"
  signup_policy          = "B2C_1_Signup"
  profile_editing_policy = "B2C_1_EditProfile"
  password_reset_policy  = "B2C_1_ResetPassword"

  depends_on = [azuread_application_password.test]
}
`, b2cConfig["tenant_id"], b2cConfig["client_id"], b2cConfig["client_secret"], b2cConfig["tenant_slug"], data.RandomInteger, data.Locations.Primary, data.RandomString, apiDomain)
}

func (r ApiManagementIdentityProviderAADB2CResource) requiresImport(data acceptance.TestData, b2cConfig map[string]string, apiDomain string) string {
	template := r.basic(data, b2cConfig, apiDomain)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_identity_provider_aadb2c" "import" {
  resource_group_name    = azurerm_api_management_identity_provider_aadb2c.test.resource_group_name
  api_management_name    = azurerm_api_management_identity_provider_aadb2c.test.api_management_name
  client_id              = azurerm_api_management_identity_provider_aadb2c.test.client_id
  client_secret          = azurerm_api_management_identity_provider_aadb2c.test.client_secret
  allowed_tenant         = azurerm_api_management_identity_provider_aadb2c.test.allowed_tenant
  signin_tenant          = azurerm_api_management_identity_provider_aadb2c.test.signin_tenant
  authority              = azurerm_api_management_identity_provider_aadb2c.test.authority
  signup_policy          = azurerm_api_management_identity_provider_aadb2c.test.signup_policy
  signin_policy          = azurerm_api_management_identity_provider_aadb2c.test.signin_policy
  profile_editing_policy = azurerm_api_management_identity_provider_aadb2c.test.profile_editing_policy
  password_reset_policy  = azurerm_api_management_identity_provider_aadb2c.test.password_reset_policy
}
`, template)
}
