// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementWorkspaceGroupTestResource struct{}

func TestAccApiManagementWorkspaceGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_group", "test")
	r := ApiManagementWorkspaceGroupTestResource{}

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

func TestAccApiManagementWorkspaceGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_group", "test")
	r := ApiManagementWorkspaceGroupTestResource{}
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

func TestAccApiManagementWorkspaceGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_group", "test")
	r := ApiManagementWorkspaceGroupTestResource{}

	r.skipIfEnvNotSet(t, "TF_ACC_ARM_CLIENT_ID", "TF_ACC_ARM_CLIENT_SECRET")

	clientId := os.Getenv("TF_ACC_ARM_CLIENT_ID")
	clientSecret := os.Getenv("TF_ACC_ARM_CLIENT_SECRET")

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, clientId, clientSecret),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementWorkspaceGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_group", "test")
	r := ApiManagementWorkspaceGroupTestResource{}

	r.skipIfEnvNotSet(t, "TF_ACC_ARM_CLIENT_ID", "TF_ACC_ARM_CLIENT_SECRET")

	clientId := os.Getenv("TF_ACC_ARM_CLIENT_ID")
	clientSecret := os.Getenv("TF_ACC_ARM_CLIENT_SECRET")

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, clientId, clientSecret),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, clientId, clientSecret),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, clientId, clientSecret),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementWorkspaceGroup_requireTypeSetToExternalWhenExternalIdIsSpecifiedError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_group", "test")
	r := ApiManagementWorkspaceGroupTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.requireTypeSetToExternalWhenExternalIdIsSpecifiedError(data),
			ExpectError: regexp.MustCompile("`type` must be set to `external` when `external_id` is specified`"),
		},
	})
}

func TestAccApiManagementWorkspaceGroup_requireExternalIdWhenTypeIsSpecifiedAsExternalError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_group", "test")
	r := ApiManagementWorkspaceGroupTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.requireExternalIdWhenTypeIsSpecifiedAsExternalError(data),
			ExpectError: regexp.MustCompile("`external_id` must be specified when `type` is set to `external`"),
		},
	})
}

func (ApiManagementWorkspaceGroupTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := group.ParseWorkspaceGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.GroupClient_v2024_05_01.WorkspaceGroupGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ApiManagementWorkspaceGroupTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_group" "test" {
  name                        = "acctest-group-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "Test Group %d"
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementWorkspaceGroupTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_workspace_group" "import" {
  name                        = azurerm_api_management_workspace_group.test.name
  api_management_workspace_id = azurerm_api_management_workspace_group.test.api_management_workspace_id
  display_name                = azurerm_api_management_workspace_group.test.display_name
}
`, r.basic(data))
}

func (r ApiManagementWorkspaceGroupTestResource) update(data acceptance.TestData, clientId string, clientSecret string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

%[1]s

data "azurerm_client_config" "current" {}

resource "azuread_group" "test" {
  display_name     = "acctest-aad-group-%[2]d"
  security_enabled = true
}

resource "azurerm_api_management_identity_provider_aad" "test" {
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  client_id           = "%[3]s"
  client_secret       = "%[4]s"
  allowed_tenants     = [data.azurerm_client_config.current.tenant_id]
}

resource "azurerm_api_management_workspace_group" "test" {
  name                        = "acctest-group-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "Updated Test Group %[2]d"
  description                 = "Updated description for test group"
  external_id                 = "aad://${data.azurerm_client_config.current.tenant_id}/groups/${azuread_group.test.object_id}"
  type                        = "external"

  depends_on = [azurerm_api_management_identity_provider_aad.test]
}
`, r.template(data), data.RandomInteger, clientId, clientSecret)
}

func (r ApiManagementWorkspaceGroupTestResource) complete(data acceptance.TestData, clientId string, clientSecret string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

%[1]s

data "azurerm_client_config" "current" {}

resource "azuread_group" "test" {
  display_name     = "acctest-aad-group-%[2]d"
  security_enabled = true
}

resource "azurerm_api_management_identity_provider_aad" "test" {
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  client_id           = "%[3]s"
  client_secret       = "%[4]s"
  allowed_tenants     = [data.azurerm_client_config.current.tenant_id]
}

resource "azurerm_api_management_workspace_group" "test" {
  name                        = "acctest-group-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "Complete Test Group %[2]d"
  description                 = "A complete test group with all properties"
  external_id                 = "aad://${data.azurerm_client_config.current.tenant_id}/groups/${azuread_group.test.object_id}"
  type                        = "external"

  depends_on = [azurerm_api_management_identity_provider_aad.test]
}
`, r.template(data), data.RandomInteger, clientId, clientSecret)
}

func (r ApiManagementWorkspaceGroupTestResource) requireTypeSetToExternalWhenExternalIdIsSpecifiedError(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

%[1]s

resource "azuread_group" "test" {
  display_name     = "acctest-aad-group-%[2]d"
  security_enabled = true
}

resource "azurerm_api_management_workspace_group" "test" {
  name                        = "acctest-group-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "Test Group %[2]d"
  external_id                 = "external group"
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceGroupTestResource) requireExternalIdWhenTypeIsSpecifiedAsExternalError(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

%[1]s

resource "azuread_group" "test" {
  display_name     = "acctest-aad-group-%[2]d"
  security_enabled = true
}

resource "azurerm_api_management_workspace_group" "test" {
  name                        = "acctest-group-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "Test Group %[2]d"
  type                        = "external"
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceGroupTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apim-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestapim-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Premium_1"
}

resource "azurerm_api_management_workspace" "test" {
  name              = "acctestapimws-%d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Test Workspace %d"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementWorkspaceGroupTestResource) skipIfEnvNotSet(t *testing.T, envVars ...string) {
	for _, env := range envVars {
		if os.Getenv(env) == "" {
			t.Skipf("skipping as %s environment variable is not set", env)
		}
	}
}
