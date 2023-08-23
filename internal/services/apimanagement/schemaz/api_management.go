// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schemaz

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/apioperation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// not in service package as migrate package required this

func SchemaApiManagementName() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validate.ApiManagementServiceName,
	}
}

func SchemaApiManagementDataSourceName() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Required:     true,
		ValidateFunc: validate.ApiManagementServiceName,
	}
}

// SchemaApiManagementChildName returns the Schema for the identifier
// used by resources within nested under the API Management Service resource
func SchemaApiManagementChildName() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validate.ApiManagementChildName,
	}
}

// SchemaApiManagementChildName returns the Schema for the identifier
// used by resources within nested under the API Management Service resource
func SchemaApiManagementApiName() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validate.ApiManagementApiName,
	}
}

// SchemaApiManagementChildDataSourceName returns the Schema for the identifier
// used by resources within nested under the API Management Service resource
func SchemaApiManagementChildDataSourceName() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Required:     true,
		ValidateFunc: validate.ApiManagementChildName,
	}
}

func SchemaApiManagementUserName() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validate.ApiManagementUserName,
	}
}

func SchemaApiManagementUserDataSourceName() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Required:     true,
		ValidateFunc: validate.ApiManagementUserName,
	}
}

func SchemaApiManagementOperationRepresentation() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"content_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},

				"form_parameter": SchemaApiManagementOperationParameterContract(),

				"example": SchemaApiManagementOperationParameterExampleContract(),

				"schema_id": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"type_name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func ExpandApiManagementOperationRepresentation(d *pluginsdk.ResourceData, schemaPath string, input []interface{}) (*[]apioperation.RepresentationContract, error) {
	if len(input) == 0 {
		return &[]apioperation.RepresentationContract{}, nil
	}

	outputs := make([]apioperation.RepresentationContract, 0)

	for i, v := range input {
		vs := v.(map[string]interface{})

		contentType := vs["content_type"].(string)
		formParametersRaw := vs["form_parameter"].([]interface{})
		formParameters := ExpandApiManagementOperationParameterContract(d, fmt.Sprintf("%s.%d.form_parameter", schemaPath, i), formParametersRaw)
		schemaId := vs["schema_id"].(string)
		typeName := vs["type_name"].(string)

		examples := make(map[string]apioperation.ParameterExampleContract)
		if vs["example"] != nil {
			examplesRaw := vs["example"].([]interface{})
			examples = ExpandApiManagementOperationParameterExampleContract(examplesRaw)
		}

		output := apioperation.RepresentationContract{
			ContentType: contentType,
			Examples:    pointer.To(examples),
		}

		contentTypeIsFormData := strings.EqualFold(contentType, "multipart/form-data") || strings.EqualFold(contentType, "application/x-www-form-urlencoded")

		// Representation formParameters can only be specified for form data content types (multipart/form-data, application/x-www-form-urlencoded)
		if contentTypeIsFormData {
			output.FormParameters = formParameters
		} else if len(*formParameters) > 0 {
			return nil, fmt.Errorf("`form_parameter` can only be specified for form data content types (multipart/form-data, application/x-www-form-urlencoded)")
		}

		// Representation schemaId can only be specified for non form data content types (multipart/form-data, application/x-www-form-urlencoded).
		// Representation typeName can only be specified for non form data content types (multipart/form-data, application/x-www-form-urlencoded).
		// nolint gocritic
		if !contentTypeIsFormData {
			output.SchemaId = pointer.To(schemaId)
			output.TypeName = pointer.To(typeName)
		} else if schemaId != "" {
			return nil, fmt.Errorf("`schema_id` cannot be specified for non-form data content types (multipart/form-data, application/x-www-form-urlencoded)")
		} else if typeName != "" {
			return nil, fmt.Errorf("`type_name` cannot be specified for non-form data content types (multipart/form-data, application/x-www-form-urlencoded)")
		}

		outputs = append(outputs, output)
	}

	return &outputs, nil
}

func FlattenApiManagementOperationRepresentation(input *[]apioperation.RepresentationContract) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	outputs := make([]interface{}, 0)

	for _, v := range *input {
		output := make(map[string]interface{})

		output["content_type"] = v.ContentType

		formParameter, err := FlattenApiManagementOperationParameterContract(v.FormParameters)
		if err != nil {
			return nil, err
		}
		output["form_parameter"] = formParameter

		if v.Examples != nil {
			example, err := FlattenApiManagementOperationParameterExampleContract(*v.Examples)
			if err != nil {
				return nil, err
			}
			output["example"] = example
		}

		if v.SchemaId != nil {
			output["schema_id"] = *v.SchemaId
		}

		if v.TypeName != nil {
			output["type_name"] = *v.TypeName
		}

		outputs = append(outputs, output)
	}

	return outputs, nil
}

func SchemaApiManagementOperationParameterContract() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"required": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},

				"description": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"default_value": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"values": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Set: pluginsdk.HashString,
				},

				"example": SchemaApiManagementOperationParameterExampleContract(),

				"schema_id": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"type_name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func ExpandApiManagementOperationParameterContract(d *pluginsdk.ResourceData, schemaPath string, input []interface{}) *[]apioperation.ParameterContract {
	if len(input) == 0 {
		return &[]apioperation.ParameterContract{}
	}

	outputs := make([]apioperation.ParameterContract, 0)

	for i, v := range input {
		vs := v.(map[string]interface{})

		name := vs["name"].(string)
		description := vs["description"].(string)
		paramType := vs["type"].(string)
		required := vs["required"].(bool)
		valuesRaw := vs["values"].(*pluginsdk.Set).List()

		schemaId := vs["schema_id"].(string)
		typeName := vs["type_name"].(string)
		examples := make(map[string]apioperation.ParameterExampleContract)
		if vs["example"] != nil {
			examplesRaw := vs["example"].([]interface{})
			examples = ExpandApiManagementOperationParameterExampleContract(examplesRaw)
		}

		output := apioperation.ParameterContract{
			Name:         name,
			Description:  pointer.To(description),
			Type:         paramType,
			Required:     pointer.To(required),
			DefaultValue: nil,
			Values:       utils.ExpandStringSlice(valuesRaw),
			SchemaId:     pointer.To(schemaId),
			TypeName:     pointer.To(typeName),
			Examples:     pointer.To(examples),
		}

		// DefaultValue must be included in Values, else it returns error
		// when DefaultValue is unset, we need to set it nil
		// "" is a valid DefaultValue
		if v, ok := d.GetOk(fmt.Sprintf("%s.%d.default_value", schemaPath, i)); ok {
			output.DefaultValue = pointer.To(v.(string))
		}
		outputs = append(outputs, output)
	}

	return &outputs
}

func FlattenApiManagementOperationParameterContract(input *[]apioperation.ParameterContract) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	outputs := make([]interface{}, 0)
	for _, v := range *input {
		output := map[string]interface{}{}

		output["name"] = v.Name
		output["description"] = pointer.From(v.Description)
		output["type"] = v.Type
		output["required"] = pointer.From(v.Required)
		output["default_value"] = pointer.From(v.DefaultValue)
		output["values"] = pluginsdk.NewSet(pluginsdk.HashString, utils.FlattenStringSlice(v.Values))

		if v.Examples != nil {
			example, err := FlattenApiManagementOperationParameterExampleContract(*v.Examples)
			if err != nil {
				return nil, err
			}
			output["example"] = example
		}

		output["schema_id"] = pointer.From(v.SchemaId)
		output["type_name"] = pointer.From(v.TypeName)

		outputs = append(outputs, output)
	}

	return outputs, nil
}

func SchemaApiManagementOperationParameterExampleContract() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},

				"summary": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"description": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"value": {
					Type:             pluginsdk.TypeString,
					Optional:         true,
					DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
				},

				"external_value": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
				},
			},
		},
	}
}

func ExpandApiManagementOperationParameterExampleContract(input []interface{}) map[string]apioperation.ParameterExampleContract {
	if len(input) == 0 {
		return map[string]apioperation.ParameterExampleContract{}
	}

	outputs := make(map[string]apioperation.ParameterExampleContract)

	for _, v := range input {
		vs := v.(map[string]interface{})

		example := apioperation.ParameterExampleContract{}

		if vs["summary"] != nil {
			example.Summary = pointer.To(vs["summary"].(string))
		}

		if vs["description"] != nil {
			example.Description = pointer.To(vs["description"].(string))
		}

		if vs["value"] != nil {
			var js interface{}
			if json.Unmarshal([]byte(vs["value"].(string)), &js) == nil {
				example.Value = pointer.To(js)
			} else {
				example.Value = pointer.To(vs["value"])
			}
		}

		if vs["external_value"] != nil {
			example.ExternalValue = pointer.To(vs["external_value"].(string))
		}

		outputs[vs["name"].(string)] = example
	}

	return outputs
}

func FlattenApiManagementOperationParameterExampleContract(input map[string]apioperation.ParameterExampleContract) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	outputs := make([]interface{}, 0)
	for k, v := range input {
		output := map[string]interface{}{}

		output["name"] = k
		output["summary"] = pointer.From(v.Summary)
		output["description"] = pointer.From(v.Description)

		// value can be any type, may be a primitive value or an object
		// https://github.com/Azure/azure-sdk-for-go/blob/main/services/apimanagement/mgmt/2021-08-01/apimanagement/models.go#L10236
		if v.Value != nil {
			value, err := convert2Json(*v.Value)
			if err != nil {
				return nil, err
			}

			output["value"] = value
		}

		output["external_value"] = pointer.From(v.ExternalValue)
		outputs = append(outputs, output)
	}

	return outputs, nil
}

// CopyCertificateAndPassword copies any certificate and password attributes
// from the old config to the current to avoid state diffs.
// Iterate through old state to find sensitive props not returned by API.
// This must be done in order to avoid state diffs.
// NOTE: this information won't be available during times like Import, so this is a best-effort.
func CopyCertificateAndPassword(vals []interface{}, hostName string, output map[string]interface{}) {
	for _, val := range vals {
		oldConfig := val.(map[string]interface{})

		if oldConfig["host_name"] == hostName {
			output["certificate_password"] = oldConfig["certificate_password"]
			output["certificate"] = oldConfig["certificate"]
			break
		}
	}
}

func convert2Json(rawVal interface{}) (string, error) {
	value := ""
	if val, ok := rawVal.(string); ok {
		value = val
	} else {
		val, err := json.Marshal(rawVal)
		if err != nil {
			return "", fmt.Errorf("failed to marshal `representations.examples.value` to json: %+v", err)
		}
		value = string(val)
	}
	return value, nil
}
