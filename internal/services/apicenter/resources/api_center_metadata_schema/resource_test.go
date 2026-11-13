// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_metadata_schema_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/metadataschemas"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiCenterMetadataSchemaResource struct{}

func TestAccApiCenterMetadataSchema_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_metadata_schema", "test")
	r := ApiCenterMetadataSchemaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_id").MatchesOtherKey(check.That("azurerm_api_center_service.test").Key("id")),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiCenterMetadataSchema_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_metadata_schema", "test")
	r := ApiCenterMetadataSchemaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("metadata_assignment.0.entity").HasValue("api"),
				check.That(data.ResourceName).Key("metadata_assignment.0.deprecated").HasValue("true"),
				check.That(data.ResourceName).Key("metadata_assignment.0.required").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiCenterMetadataSchema_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_metadata_schema", "test")
	r := ApiCenterMetadataSchemaResource{}

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

func TestAccApiCenterMetadataSchema_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_metadata_schema", "test")
	r := ApiCenterMetadataSchemaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("metadata_assignment.0.deprecated").HasValue("true"),
				check.That(data.ResourceName).Key("metadata_assignment.0.required").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func (ApiCenterMetadataSchemaResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := metadataschemas.ParseMetadataSchemaID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiCenter.MetadataSchemasClient.Get(ctx, *id)
	if response.WasNotFound(resp.HttpResponse) {
		return pointer.To(false), nil
	}
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (ApiCenterMetadataSchemaResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apicenter-%d"
  location = "%s"
}

resource "azurerm_api_center_service" "test" {
  name                = "acctestApiCSvc%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ApiCenterMetadataSchemaResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_center_metadata_schema" "test" {
  name       = "acctest%d"
  service_id = azurerm_api_center_service.test.id

  metadata_assignment {
    entity = "api"
  }

  schema = jsonencode({
    type = "object"
    properties = {
      name = {
        type = "string"
      }
    }
  })
}
`, r.template(data), data.RandomInteger)
}

func (r ApiCenterMetadataSchemaResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_center_metadata_schema" "test" {
  name       = "acctest%d"
  service_id = azurerm_api_center_service.test.id

  metadata_assignment {
    entity     = "api"
    deprecated = true
    required   = true
  }

  schema = jsonencode({
    type = "object"
    properties = {
      description = {
        type = "string"
      }
    }
  })
}
`, r.template(data), data.RandomInteger)
}

func (r ApiCenterMetadataSchemaResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_center_metadata_schema" "import" {
  name       = azurerm_api_center_metadata_schema.test.name
  service_id = azurerm_api_center_metadata_schema.test.service_id

  metadata_assignment {
    entity = azurerm_api_center_metadata_schema.test.metadata_assignment[0].entity
  }

  schema = azurerm_api_center_metadata_schema.test.schema
}
`, r.basic(data))
}
