// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementGlobalSchemaResource struct{}

func TestAccApiManagementGlobalSchema_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_global_schema", "test")
	r := ApiManagementGlobalSchemaResource{}

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

func TestAccApiManagementGlobalSchema_jsonSchema(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_global_schema", "test")
	r := ApiManagementGlobalSchemaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.jsonSchema(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementGlobalSchema_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_global_schema", "test")
	r := ApiManagementGlobalSchemaResource{}

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
	})
}

func TestAccApiManagementGlobalSchema_requireImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_global_schema", "test")
	r := ApiManagementGlobalSchemaResource{}

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

func (ApiManagementGlobalSchemaResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := schema.ParseSchemaID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.GlobalSchemaClient.GlobalSchemaGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (r ApiManagementGlobalSchemaResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}

resource "azurerm_api_management_global_schema" "test" {
  schema_id           = "accetestSchema-%[1]d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  type                = "xml"
  value               = <<XML
    <xsd:schema xmlns:xsd="http://www.w3.org/2001/XMLSchema"
    xmlns:tns="http://tempuri.org/PurchaseOrderSchema.xsd" targetNamespace="http://tempuri.org/PurchaseOrderSchema.xsd" elementFormDefault="qualified">
    <xsd:element name="PurchaseOrder" type="tns:PurchaseOrderType"/>
    <xsd:complexType name="PurchaseOrderType">
        <xsd:sequence>
            <xsd:element name="ShipTo" type="tns:USAddress" maxOccurs="2"/>
            <xsd:element name="BillTo" type="tns:USAddress"/>
        </xsd:sequence>
        <xsd:attribute name="OrderDate" type="xsd:date"/>
    </xsd:complexType>
    <xsd:complexType name="USAddress">
        <xsd:sequence>
            <xsd:element name="name" type="xsd:string"/>
            <xsd:element name="street" type="xsd:string"/>
            <xsd:element name="city" type="xsd:string"/>
            <xsd:element name="state" type="xsd:string"/>
            <xsd:element name="zip" type="xsd:integer"/>
        </xsd:sequence>
        <xsd:attribute name="country" type="xsd:NMTOKEN" fixed="US"/>
    </xsd:complexType>
</xsd:schema>
XML
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ApiManagementGlobalSchemaResource) jsonSchema(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}

resource "azurerm_api_management_global_schema" "test" {
  schema_id           = "accetestSchema-%[1]d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  type                = "json"
  value               = <<JSON
{
    "schema-bug-example": {
        "properties": {
            "Field2": {
                "description": "Field2",
                "type": "string"
            },
            "field1": {
                "description": "Field1",
                "type": "string"
            }
        },
        "required": [
            "field1",
            "Field2"
        ],
        "type": "object"
    }
}
JSON
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ApiManagementGlobalSchemaResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}

resource "azurerm_api_management_global_schema" "test" {
  schema_id           = "accetestSchema-%[1]d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  type                = "xml"
  value               = <<XML
    <xsd:schema xmlns:xsd="http://www.w3.org/2001/XMLSchema"
    xmlns:tns="http://tempuri.org/PurchaseOrderSchema.xsd" targetNamespace="http://tempuri.org/PurchaseOrderSchema.xsd" elementFormDefault="qualified">
    <xsd:element name="PurchaseOrder" type="tns:PurchaseOrderType"/>
    <xsd:complexType name="PurchaseOrderType">
        <xsd:sequence>
            <xsd:element name="ShipTo" type="tns:USAddress" maxOccurs="2"/>
            <xsd:element name="BillTo" type="tns:USAddress"/>
        </xsd:sequence>
        <xsd:attribute name="OrderDate" type="xsd:date"/>
    </xsd:complexType>
    <xsd:complexType name="USAddress">
        <xsd:sequence>
            <xsd:element name="name" type="xsd:string"/>
            <xsd:element name="street" type="xsd:string"/>
            <xsd:element name="city" type="xsd:string"/>
        </xsd:sequence>
        <xsd:attribute name="country" type="xsd:NMTOKEN" fixed="US"/>
    </xsd:complexType>
</xsd:schema>
XML
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ApiManagementGlobalSchemaResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_global_schema" "import" {
  schema_id           = azurerm_api_management_global_schema.test.schema_id
  api_management_name = azurerm_api_management_global_schema.test.api_management_name
  resource_group_name = azurerm_api_management_global_schema.test.resource_group_name
  type                = azurerm_api_management_global_schema.test.type
  value               = azurerm_api_management_global_schema.test.value
}
`, r.basic(data))
}
