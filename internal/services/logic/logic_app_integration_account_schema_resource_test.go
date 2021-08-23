package logic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogicAppIntegrationAccountSchemaResource struct {
}

func TestAccLogicAppIntegrationAccountSchema_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_schema", "test")
	r := LogicAppIntegrationAccountSchemaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("content"), // not returned from the API
	})
}

func TestAccLogicAppIntegrationAccountSchema_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_schema", "test")
	r := LogicAppIntegrationAccountSchemaResource{}

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

func TestAccLogicAppIntegrationAccountSchema_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_schema", "test")
	r := LogicAppIntegrationAccountSchemaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("file_name", "content"), // not returned from the API
	})
}

func TestAccLogicAppIntegrationAccountSchema_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_schema", "test")
	r := LogicAppIntegrationAccountSchemaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("content"), // not returned from the API
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("file_name", "content"), // not returned from the API
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("content"), // not returned from the API
	})
}

func (LogicAppIntegrationAccountSchemaResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.IntegrationAccountSchemaID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Logic.IntegrationAccountSchemaClient.Get(ctx, id.ResourceGroup, id.IntegrationAccountName, id.SchemaName)
	if err != nil {
		return nil, fmt.Errorf("reading %q: %+v", id, err)
	}

	return utils.Bool(resp.IntegrationAccountSchemaProperties != nil), nil
}

func (r LogicAppIntegrationAccountSchemaResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-logic-%d"
  location = "%s"
}

resource "azurerm_logic_app_integration_account" "test" {
  name                = "acctest-ia-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r LogicAppIntegrationAccountSchemaResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_schema" "test" {
  name                     = "acctest-iaschema-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name

  content = <<XML
<xs:schema xmlns:b="http://schemas.microsoft.com/BizTalk/2003"
           xmlns="http://Inbound_EDI.OrderFile"
           targetNamespace="http://Inbound_EDI.OrderFile"
           xmlns:xs="http://www.w3.org/2001/XMLSchema">
<xs:annotation>
<xs:appinfo>
<b:schemaInfo default_pad_char=" "
              count_positions_by_byte="false"
              parser_optimization="speed"
              lookahead_depth="3"
              suppress_empty_nodes="false"
              generate_empty_nodes="true"
              allow_early_termination="false"
              early_terminate_optional_fields="false"
              allow_message_breakup_of_infix_root="false"
              compile_parse_tables="false"
              standard="Flat File"
              root_reference="OrderFile" />
<schemaEditorExtension:schemaInfo namespaceAlias="b"
                                  extensionClass="Microsoft.BizTalk.FlatFileExtension.FlatFileExtension"
                                  standardName="Flat File"
                                  xmlns:schemaEditorExtension="http://schemas.microsoft.com/BizTalk/2003/SchemaEditorExtensions" />
</xs:appinfo>
</xs:annotation>
<xs:element name="OrderFile">
<xs:annotation>
<xs:appinfo>
<b:recordInfo structure="delimited"
              preserve_delimiter_for_empty_data="true"
              suppress_trailing_delimiters="false"
              sequence_number="1" />
</xs:appinfo>
</xs:annotation>
<xs:complexType>
<xs:sequence>
<xs:annotation>
<xs:appinfo>
<b:groupInfo sequence_number="0" />
</xs:appinfo>
</xs:annotation>
<xs:element name="Order">
<xs:annotation>
<xs:appinfo>
<b:recordInfo sequence_number="1"
              structure="delimited"
              preserve_delimiter_for_empty_data="true"
              suppress_trailing_delimiters="false"
              child_delimiter_type="hex"
              child_delimiter="0x0D 0x0A"
              child_order="infix" />
</xs:appinfo>
</xs:annotation>
<xs:complexType>
<xs:sequence>
<xs:annotation>
<xs:appinfo>
<b:groupInfo sequence_number="0" />
</xs:appinfo>
</xs:annotation>
<xs:element name="Header">
<xs:annotation>
<xs:appinfo>
<b:recordInfo sequence_number="1"
              structure="delimited"
              preserve_delimiter_for_empty_data="true"
              suppress_trailing_delimiters="false"
              child_delimiter_type="char"
              child_delimiter="|"
              child_order="infix"
              tag_name="HDR|" />
</xs:appinfo>
</xs:annotation>
<xs:complexType>
<xs:sequence>
<xs:annotation>
<xs:appinfo>
<b:groupInfo sequence_number="0" />
</xs:appinfo>
</xs:annotation>
<xs:element name="PODate"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="1"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="PONumber"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo justification="left"
             sequence_number="2" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="CustomerID"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="3"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="CustomerContactName"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="5"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="CustomerContactPhone"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="5"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
</xs:sequence>
</xs:complexType>
</xs:element>
<xs:element minOccurs="1"
            maxOccurs="unbounded"
            name="LineItems">
<xs:annotation>
<xs:appinfo>
<b:recordInfo sequence_number="2"
              structure="delimited"
              preserve_delimiter_for_empty_data="true"
              suppress_trailing_delimiters="false"
              child_delimiter_type="char"
              child_delimiter="|"
              child_order="infix"
              tag_name="DTL|" />
</xs:appinfo>
</xs:annotation>
<xs:complexType>
<xs:sequence>
<xs:annotation>
<xs:appinfo>
<b:groupInfo sequence_number="0" />
</xs:appinfo>
</xs:annotation>
<xs:element name="PONumber"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="1"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="ItemOrdered"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="2"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="Quantity"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="3"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="UOM"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="4"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="Price"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="5"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="ExtendedPrice"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="6"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="Description"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="7"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
</xs:sequence>
</xs:complexType>
</xs:element>
</xs:sequence>
</xs:complexType>
</xs:element>
</xs:sequence>
</xs:complexType>
</xs:element>
</xs:schema>
XML
}
`, r.template(data), data.RandomInteger)
}

func (r LogicAppIntegrationAccountSchemaResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_schema" "test" {
  name                     = "acctest-iaschema-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name

  content = <<XML
<xs:schema xmlns:b="http://schemas.microsoft.com/BizTalk/2003"
           xmlns="http://Inbound_EDI.OrderFile"
           targetNamespace="http://Inbound_EDI.OrderFile"
           xmlns:xs="http://www.w3.org/2001/XMLSchema">
<xs:annotation>
<xs:appinfo>
<b:schemaInfo default_pad_char=" "
              count_positions_by_byte="false"
              parser_optimization="speed"
              lookahead_depth="3"
              suppress_empty_nodes="false"
              generate_empty_nodes="true"
              allow_early_termination="false"
              early_terminate_optional_fields="false"
              allow_message_breakup_of_infix_root="false"
              compile_parse_tables="false"
              standard="Flat File"
              root_reference="OrderFile" />
<schemaEditorExtension:schemaInfo namespaceAlias="b"
                                  extensionClass="Microsoft.BizTalk.FlatFileExtension.FlatFileExtension"
                                  standardName="Flat File"
                                  xmlns:schemaEditorExtension="http://schemas.microsoft.com/BizTalk/2003/SchemaEditorExtensions" />
</xs:appinfo>
</xs:annotation>
<xs:element name="OrderFile">
<xs:annotation>
<xs:appinfo>
<b:recordInfo structure="delimited"
              preserve_delimiter_for_empty_data="true"
              suppress_trailing_delimiters="false"
              sequence_number="1" />
</xs:appinfo>
</xs:annotation>
<xs:complexType>
<xs:sequence>
<xs:annotation>
<xs:appinfo>
<b:groupInfo sequence_number="0" />
</xs:appinfo>
</xs:annotation>
<xs:element name="Order">
<xs:annotation>
<xs:appinfo>
<b:recordInfo sequence_number="1"
              structure="delimited"
              preserve_delimiter_for_empty_data="true"
              suppress_trailing_delimiters="false"
              child_delimiter_type="hex"
              child_delimiter="0x0D 0x0A"
              child_order="infix" />
</xs:appinfo>
</xs:annotation>
<xs:complexType>
<xs:sequence>
<xs:annotation>
<xs:appinfo>
<b:groupInfo sequence_number="0" />
</xs:appinfo>
</xs:annotation>
<xs:element name="Header">
<xs:annotation>
<xs:appinfo>
<b:recordInfo sequence_number="1"
              structure="delimited"
              preserve_delimiter_for_empty_data="true"
              suppress_trailing_delimiters="false"
              child_delimiter_type="char"
              child_delimiter="|"
              child_order="infix"
              tag_name="HDR|" />
</xs:appinfo>
</xs:annotation>
<xs:complexType>
<xs:sequence>
<xs:annotation>
<xs:appinfo>
<b:groupInfo sequence_number="0" />
</xs:appinfo>
</xs:annotation>
<xs:element name="PODate"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="1"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="PONumber"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo justification="left"
             sequence_number="2" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="CustomerID"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="3"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="CustomerContactName"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="4"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="CustomerContactPhone"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="5"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
</xs:sequence>
</xs:complexType>
</xs:element>
<xs:element minOccurs="1"
            maxOccurs="unbounded"
            name="LineItems">
<xs:annotation>
<xs:appinfo>
<b:recordInfo sequence_number="2"
              structure="delimited"
              preserve_delimiter_for_empty_data="true"
              suppress_trailing_delimiters="false"
              child_delimiter_type="char"
              child_delimiter="|"
              child_order="infix"
              tag_name="DTL|" />
</xs:appinfo>
</xs:annotation>
<xs:complexType>
<xs:sequence>
<xs:annotation>
<xs:appinfo>
<b:groupInfo sequence_number="0" />
</xs:appinfo>
</xs:annotation>
<xs:element name="PONumber"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="1"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="ItemOrdered"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="2"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="Quantity"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="3"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="UOM"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="4"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="Price"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="5"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="ExtendedPrice"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="6"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
<xs:element name="Description"
            type="xs:string">
<xs:annotation>
<xs:appinfo>
<b:fieldInfo sequence_number="7"
             justification="left" />
</xs:appinfo>
</xs:annotation>
</xs:element>
</xs:sequence>
</xs:complexType>
</xs:element>
</xs:sequence>
</xs:complexType>
</xs:element>
</xs:sequence>
</xs:complexType>
</xs:element>
</xs:schema>
XML

  file_name = "TestFile.xsd"

  metadata = <<METADATA
    {
        "foo": "bar"
    }
METADATA
}
`, r.template(data), data.RandomInteger)
}

func (r LogicAppIntegrationAccountSchemaResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_schema" "import" {
  name                     = azurerm_logic_app_integration_account_schema.test.name
  resource_group_name      = azurerm_logic_app_integration_account_schema.test.resource_group_name
  integration_account_name = azurerm_logic_app_integration_account_schema.test.integration_account_name
  content                  = azurerm_logic_app_integration_account_schema.test.content
}
`, r.basic(data))
}
