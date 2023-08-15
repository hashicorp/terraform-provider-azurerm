// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/apioperationtag"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementApiOperationTagResource struct{}

func TestAccApiManagementApiOperationTag_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation_tag", "test")
	r := ApiManagementApiOperationTagResource{}

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

func TestAccApiManagementApiOperationTag_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation_tag", "test")
	r := ApiManagementApiOperationTagResource{}

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

func TestAccApiManagementApiOperationTag_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation_tag", "test")
	r := ApiManagementApiOperationTagResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ApiManagementApiOperationTagResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := apioperationtag.ParseOperationTagID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.ApiOperationTagClient.TagGetByOperation(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (r ApiManagementApiOperationTagResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation_tag" "test" {
  api_operation_id = azurerm_api_management_api_operation.test.id
  name             = "acctest-Op-Tag-%d"
  display_name     = "Display-Op-Tag"
}
`, ApiManagementApiOperationResource{}.basic(data), data.RandomInteger)
}

func (r ApiManagementApiOperationTagResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation_tag" "import" {
  api_operation_id = azurerm_api_management_api_operation_tag.test.api_operation_id
  name             = azurerm_api_management_api_operation_tag.test.name
  display_name     = azurerm_api_management_api_operation_tag.test.display_name
}
`, r.basic(data))
}

func (r ApiManagementApiOperationTagResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation_tag" "test" {
  api_operation_id = azurerm_api_management_api_operation.test.id
  name             = "acctest-Op-Tag-%d"

  display_name = "Display-Op-Tag Updated"
}
`, ApiManagementApiOperationResource{}.basic(data), data.RandomInteger)
}
