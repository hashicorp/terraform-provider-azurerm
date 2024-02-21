// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/policyfragment"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementPolicyFragmentResource struct{}

func TestAccApiManagementPolicyFragment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_policy_fragment", "test")
	r := ApiManagementPolicyFragmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue("Some descriptive text"),
				check.That(data.ResourceName).Key("format").HasValue("xml"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementPolicyFragment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_policy_fragment", "test")
	r := ApiManagementPolicyFragmentResource{}

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

func TestAccApiManagementPolicyFragment_updateDescription(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_policy_fragment", "test")
	r := ApiManagementPolicyFragmentResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue("Some descriptive text"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue("Some descriptive text which is updated"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementPolicyFragment_rawxml(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_policy_fragment", "test")
	r := ApiManagementPolicyFragmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.rawxml(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ApiManagementPolicyFragmentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := policyfragment.ParsePolicyFragmentIDInsensitively(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.PolicyFragmentClient.Get(ctx, *id, policyfragment.GetOperationOptions{})
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (r ApiManagementPolicyFragmentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_policy_fragment" "test" {
  name                = azurerm_api_management.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  description         = "Some descriptive text"
  value               = file("testdata/api_management_policy_fragment_test_xml.xml")
}
`, r.template(data))
}

func (r ApiManagementPolicyFragmentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_policy_fragment" "import" {
  name                = azurerm_api_management_policy_fragment.test.name
  api_management_name = azurerm_api_management_policy_fragment.test.api_management_name
  resource_group_name = azurerm_api_management_policy_fragment.test.resource_group_name
  description         = azurerm_api_management_policy_fragment.test.description
  value               = azurerm_api_management_policy_fragment.test.value
}
`, r.basic(data))
}

func (r ApiManagementPolicyFragmentResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_policy_fragment" "test" {
  name                = azurerm_api_management.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  description         = "Some descriptive text which is updated"
  value               = file("testdata/api_management_policy_fragment_test_xml.xml")
}
`, r.template(data))
}

func (r ApiManagementPolicyFragmentResource) rawxml(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_policy_fragment" "test" {
  name                = azurerm_api_management.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  description         = "Some descriptive text"
  format              = "rawxml"
  value               = file("testdata/api_management_policy_fragment_test_rawxml.xml")
}
`, r.template(data))
}

func (ApiManagementPolicyFragmentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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
  sku_name            = "Developer_1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
