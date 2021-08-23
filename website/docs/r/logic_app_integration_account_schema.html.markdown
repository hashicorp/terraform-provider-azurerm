---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_integration_account_schema"
description: |-
  Manages a Logic App Integration Account Schema.
---

# azurerm_logic_app_integration_account_schema

Manages a Logic App Integration Account Schema.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_logic_app_integration_account" "example" {
  name                = "example-ia"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Basic"
}

resource "azurerm_logic_app_integration_account_schema" "example" {
  name                     = "example-ias"
  resource_group_name      = azurerm_resource_group.example.name
  integration_account_name = azurerm_logic_app_integration_account.example.name

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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Logic App Integration Account Schema. Changing this forces a new Logic App Integration Account Schema to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Logic App Integration Account Schema should exist. Changing this forces a new Logic App Integration Account Schema to be created.

* `integration_account_name` - (Required) The name of the Logic App Integration Account. Changing this forces a new Logic App Integration Account Schema to be created.

* `content` - (Required) The content of the Logic App Integration Account Schema.

* `file_name` - (Optional) The file name of the Logic App Integration Account Schema.

* `metadata` - (Optional) The metadata of the Logic App Integration Account Schema.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Logic App Integration Account Schema.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Logic App Integration Account Schema.
* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Integration Account Schema.
* `update` - (Defaults to 30 minutes) Used when updating the Logic App Integration Account Schema.
* `delete` - (Defaults to 30 minutes) Used when deleting the Logic App Integration Account Schema.

## Import

Logic App Integration Account Schemas can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logic_app_integration_account_schema.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Logic/integrationAccounts/account1/schemas/schema1
```
