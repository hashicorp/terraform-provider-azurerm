// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/apioperationpolicy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementApiOperationPolicyResource struct{}

func TestAccApiManagementAPIOperationPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation_policy", "test")
	r := ApiManagementApiOperationPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			ResourceName:            data.ResourceName,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"xml_link"},
		},
	})
}

func TestAccApiManagementAPIOperationPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation_policy", "test")
	r := ApiManagementApiOperationPolicyResource{}

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

func TestAccApiManagementAPIOperationPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation_policy", "test")
	r := ApiManagementApiOperationPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			ResourceName:            data.ResourceName,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"xml_link"},
		},
	})
}

func TestAccApiManagementAPIOperationPolicy_rawXml(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation_policy", "test")
	r := ApiManagementApiOperationPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.rawXml(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ApiManagementApiOperationPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := apioperationpolicy.ParseOperationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.ApiOperationPoliciesClient.Get(ctx, *id, apioperationpolicy.GetOperationOptions{Format: pointer.To(apioperationpolicy.PolicyExportFormatXml)})
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (r ApiManagementApiOperationPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation_policy" "test" {
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  operation_id        = azurerm_api_management_api_operation.test.operation_id
  xml_link            = "https://gist.githubusercontent.com/riordanp/ca22f8113afae0eb38cc12d718fd048d/raw/d6ac89a2f35a6881a7729f8cb4883179dc88eea1/example.xml"
}
`, ApiManagementApiOperationResource{}.basic(data))
}

func (r ApiManagementApiOperationPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation_policy" "import" {
  api_name            = azurerm_api_management_api_operation_policy.test.api_name
  api_management_name = azurerm_api_management_api_operation_policy.test.api_management_name
  resource_group_name = azurerm_api_management_api_operation_policy.test.resource_group_name
  operation_id        = azurerm_api_management_api_operation_policy.test.operation_id
  xml_link            = azurerm_api_management_api_operation_policy.test.xml_link
}
`, r.basic(data))
}

func (r ApiManagementApiOperationPolicyResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation_policy" "test" {
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  operation_id        = azurerm_api_management_api_operation.test.operation_id

  xml_content = <<XML
<policies>
  <inbound>
    <set-variable name="abc" value="@(context.Request.Headers.GetValueOrDefault("X-Header-Name", ""))" />
    <find-and-replace from="xyz" to="abc" />
  </inbound>
</policies>
XML

}
`, ApiManagementApiOperationResource{}.basic(data))
}

func (r ApiManagementApiOperationPolicyResource) rawXml(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation_policy" "test" {
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  operation_id        = azurerm_api_management_api_operation.test.operation_id

  xml_content = file("testdata/api_management_api_operation_policy.xml")
}
`, ApiManagementApiOperationResource{}.basic(data))
}
