// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/identityprovider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

/*
Note that this resource requires a preexisting B2C tenant with matching policies already created.
Accordingly, these tests rely on additional environment variables to be set (and be valid):
* ARM_TEST_B2C_TENANT_ID      - the UUID of the B2C tenant
* ARM_TEST_B2C_TENANT_SLUG    - the first part of the *.onmicrosoft.com domain of the B2C tenant
* ARM_TEST_B2C_CLIENT_ID      - client ID of an application in the B2C domain with privileges to read directory data
* ARM_TEST_B2C_CLIENT_SECRET  - client secret for that application
*/

type ApiManagementIdentityProviderAADB2CResource struct{}

func TestAccAzureRMApiManagementIdentityProviderAADB2C_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_aadb2c", "test")
	r := ApiManagementIdentityProviderAADB2CResource{}
	b2cConfig := testAccAzureRMApiManagementIdentityProviderAADB2C_getB2CConfig(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, b2cConfig),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret"),
	})
}

func TestAccAzureRMApiManagementIdentityProviderAADB2C_clientLibrary(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_aadb2c", "test")
	r := ApiManagementIdentityProviderAADB2CResource{}
	b2cConfig := testAccAzureRMApiManagementIdentityProviderAADB2C_getB2CConfig(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.clientLibrary(data, b2cConfig),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret"),
		{
			Config: r.clientLibraryUpdate(data, b2cConfig),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret"),
		{
			Config: r.clientLibrary(data, b2cConfig),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, b2cConfig),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, b2cConfig),
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
		t.Skipf("Acceptance tests for resource `azurerm_api_management_identity_provider_aadb2c` skipped unless environment variables set: %s", strings.Join(vars, ", "))
	}

	return config
}

func (ApiManagementIdentityProviderAADB2CResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := identityprovider.ParseIdentityProviderID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.IdentityProviderClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (ApiManagementIdentityProviderAADB2CResource) basic(data acceptance.TestData, b2cConfig map[string]string) string {
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
  display_name = "acctestAM-%[5]d"
  web {
    redirect_uris = [azurerm_api_management.test.developer_portal_url]

    implicit_grant {
      access_token_issuance_enabled = true
    }
  }
}

resource "azuread_application_password" "test" {
  application_object_id = azuread_application.test.object_id
}

resource "azurerm_api_management_identity_provider_aadb2c" "test" {
  resource_group_name    = azurerm_resource_group.test.name
  api_management_name    = azurerm_api_management.test.name
  client_id              = azuread_application.test.application_id
  client_secret          = azuread_application_password.test.value
  allowed_tenant         = "%[4]s.onmicrosoft.com"
  signin_tenant          = "%[4]s.onmicrosoft.com"
  authority              = "%[4]s.b2clogin.com"
  signin_policy          = "B2C_1_Login"
  signup_policy          = "B2C_1_Signup"
  profile_editing_policy = "B2C_1_EditProfile"
  password_reset_policy  = "B2C_1_ResetPassword"

  depends_on = [azuread_application_password.test]
}
`, b2cConfig["tenant_id"], b2cConfig["client_id"], b2cConfig["client_secret"], b2cConfig["tenant_slug"], data.RandomInteger, data.Locations.Primary)
}

func (ApiManagementIdentityProviderAADB2CResource) clientLibrary(data acceptance.TestData, b2cConfig map[string]string) string {
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
  display_name = "acctestAM-%[5]d"
  web {
    redirect_uris = [azurerm_api_management.test.developer_portal_url]

    implicit_grant {
      access_token_issuance_enabled = true
    }
  }
}

resource "azuread_application_password" "test" {
  application_object_id = azuread_application.test.object_id
}

resource "azurerm_api_management_identity_provider_aadb2c" "test" {
  resource_group_name    = azurerm_resource_group.test.name
  api_management_name    = azurerm_api_management.test.name
  client_id              = azuread_application.test.application_id
  client_library         = "MSAL"
  client_secret          = azuread_application_password.test.value
  allowed_tenant         = "%[4]s.onmicrosoft.com"
  signin_tenant          = "%[4]s.onmicrosoft.com"
  authority              = "%[4]s.b2clogin.com"
  signin_policy          = "B2C_1_Login"
  signup_policy          = "B2C_1_Signup"
  profile_editing_policy = "B2C_1_EditProfile"
  password_reset_policy  = "B2C_1_ResetPassword"

  depends_on = [azuread_application_password.test]
}
`, b2cConfig["tenant_id"], b2cConfig["client_id"], b2cConfig["client_secret"], b2cConfig["tenant_slug"], data.RandomInteger, data.Locations.Primary)
}

func (ApiManagementIdentityProviderAADB2CResource) clientLibraryUpdate(data acceptance.TestData, b2cConfig map[string]string) string {
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
  display_name = "acctestAM-%[5]d"
  web {
    redirect_uris = [azurerm_api_management.test.developer_portal_url]

    implicit_grant {
      access_token_issuance_enabled = true
    }
  }
}

resource "azuread_application_password" "test" {
  application_object_id = azuread_application.test.object_id
}

resource "azurerm_api_management_identity_provider_aadb2c" "test" {
  resource_group_name    = azurerm_resource_group.test.name
  api_management_name    = azurerm_api_management.test.name
  client_id              = azuread_application.test.application_id
  client_library         = "MSAL-2"
  client_secret          = azuread_application_password.test.value
  allowed_tenant         = "%[4]s.onmicrosoft.com"
  signin_tenant          = "%[4]s.onmicrosoft.com"
  authority              = "%[4]s.b2clogin.com"
  signin_policy          = "B2C_1_Login"
  signup_policy          = "B2C_1_Signup"
  profile_editing_policy = "B2C_1_EditProfile"
  password_reset_policy  = "B2C_1_ResetPassword"

  depends_on = [azuread_application_password.test]
}
`, b2cConfig["tenant_id"], b2cConfig["client_id"], b2cConfig["client_secret"], b2cConfig["tenant_slug"], data.RandomInteger, data.Locations.Primary)
}

func (r ApiManagementIdentityProviderAADB2CResource) requiresImport(data acceptance.TestData, b2cConfig map[string]string) string {
	template := r.basic(data, b2cConfig)
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
