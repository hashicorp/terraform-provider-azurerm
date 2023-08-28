// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/apitag"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementApiTagResource struct{}

func TestAccApiManagementApiTag_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_tag", "test")
	r := ApiManagementApiTagResource{}

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

func TestAccApiManagementApiTag_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_tag", "test")
	r := ApiManagementApiTagResource{}

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

func (ApiManagementApiTagResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := apitag.ParseApiTagID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.ApiTagClient.TagGetByApi(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (r ApiManagementApiTagResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_api_management_tag" "test" {
  api_management_id = azurerm_api_management.test.id
  name              = "acctest-Tag-%d"
}

resource "azurerm_api_management_api_tag" "test" {
  api_id = azurerm_api_management_api.test.id
  name   = "acctest-Tag-%d"
}
`, ApiManagementApiResource{}.basic(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementApiTagResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_api_management_api_tag" "import" {
  api_id = azurerm_api_management_api.test.id
  name   = azurerm_api_management_api_tag.test.name
}
`, r.basic(data))
}
