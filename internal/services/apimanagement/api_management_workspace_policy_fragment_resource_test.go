// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/policyfragment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementWorkspacePolicyFragmentTestResource struct{}

func TestAccApiManagementWorkspacePolicyFragment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_policy_fragment", "test")
	r := ApiManagementWorkspacePolicyFragmentTestResource{}

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

func TestAccApiManagementWorkspacePolicyFragment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_policy_fragment", "test")
	r := ApiManagementWorkspacePolicyFragmentTestResource{}

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

func TestAccApiManagementWorkspacePolicyFragment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_policy_fragment", "test")
	r := ApiManagementWorkspacePolicyFragmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// Because of API behavior, Workspace Policy Fragments are always imported as `xml`.
		// As a result, the `xml_format` and `xml_content`properties should be ignored when set to `rawxml`.
		data.ImportStep("xml_format", "xml_content"),
	})
}

func TestAccApiManagementWorkspacePolicyFragment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_policy_fragment", "test")
	r := ApiManagementWorkspacePolicyFragmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// Because of API behavior, Workspace Policy Fragments are always imported as `xml`.
		// As a result, the `xml_format` and `xml_content`properties should be ignored when set to `rawxml`.
		data.ImportStep("xml_format", "xml_content"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// Because of API behavior, Workspace Policy Fragments are always imported as `xml`.
		// As a result, the `xml_format` and `xml_content`properties should be ignored when set to `rawxml`.
		data.ImportStep("xml_format", "xml_content"),
	})
}

func (ApiManagementWorkspacePolicyFragmentTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := policyfragment.ParseWorkspacePolicyFragmentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.PolicyFragmentClient_v2024_05_01.WorkspacePolicyFragmentGet(ctx, *id, policyfragment.WorkspacePolicyFragmentGetOperationOptions{})
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ApiManagementWorkspacePolicyFragmentTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_policy_fragment" "test" {
  name                        = "acctestpolicyfragment%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  xml_content                 = file("testdata/api_management_policy_fragment_test_xml.xml")
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspacePolicyFragmentTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_workspace_policy_fragment" "import" {
  name                        = azurerm_api_management_workspace_policy_fragment.test.name
  api_management_workspace_id = azurerm_api_management_workspace_policy_fragment.test.api_management_workspace_id
  xml_content                 = azurerm_api_management_workspace_policy_fragment.test.xml_content
}
`, r.basic(data))
}

func (r ApiManagementWorkspacePolicyFragmentTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_policy_fragment" "test" {
  name                        = "acctestpolicyfragment%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  description                 = "A test policy fragment"
  xml_format                  = "rawxml"
  xml_content                 = file("testdata/api_management_policy_fragment_test_rawxml.xml")
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspacePolicyFragmentTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_policy_fragment" "test" {
  name                        = "acctestpolicyfragment%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  description                 = "Updated policy fragment description"
  xml_format                  = "xml"
  xml_content                 = file("testdata/api_management_policy_fragment_test_xml.xml")
}
`, r.template(data), data.RandomInteger)
}

func (ApiManagementWorkspacePolicyFragmentTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Premium_1"
}

resource "azurerm_api_management_workspace" "test" {
  name              = "acctestAMWorkspace-%d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Test Workspace"
  description       = "Test workspace for policy fragments"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
