package schemaz

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/features"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
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
	out := &pluginsdk.Schema{
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
	if !features.ThreePointOhBeta() {
		s := out.Elem.(*pluginsdk.Resource)
		s.Schema["sample"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeString,
			Optional:   true,
			Computed:   true,
			Deprecated: "Deprecated in favour of `example`",
		}
	}
	return out
}

func ExpandApiManagementOperationRepresentation(d *pluginsdk.ResourceData, schemaPath string, input []interface{}) (*[]apimanagement.RepresentationContract, error) {
	if len(input) == 0 {
		return &[]apimanagement.RepresentationContract{}, nil
	}

	outputs := make([]apimanagement.RepresentationContract, 0)

	for i, v := range input {
		vs := v.(map[string]interface{})

		contentType := vs["content_type"].(string)
		formParametersRaw := vs["form_parameter"].([]interface{})
		formParameters := ExpandApiManagementOperationParameterContract(d, fmt.Sprintf("%s.%d.form_parameter", schemaPath, i), formParametersRaw)
		schemaId := vs["schema_id"].(string)
		typeName := vs["type_name"].(string)

		examples := make(map[string]*apimanagement.ParameterExampleContract)
		if vs["example"] != nil {
			examplesRaw := vs["example"].([]interface{})
			examples = ExpandApiManagementOperationParameterExampleContract(examplesRaw)
		} else if !features.ThreePointOhBeta() && vs["sample"] != nil {
			defaultExample := map[string]interface{}{
				"name":  "default",
				"value": vs["sample"],
			}
			examples = ExpandApiManagementOperationParameterExampleContract([]interface{}{defaultExample})
		}

		output := apimanagement.RepresentationContract{
			ContentType: utils.String(contentType),
			Examples:    examples,
		}

		contentTypeIsFormData := strings.EqualFold(contentType, "multipart/form-data") || strings.EqualFold(contentType, "application/x-www-form-urlencoded")

		// Representation formParameters can only be specified for form data content types (multipart/form-data, application/x-www-form-urlencoded)
		if contentTypeIsFormData {
			output.FormParameters = formParameters
		} else if len(*formParameters) > 0 {
			return nil, fmt.Errorf("`form_parameter` cannot be specified for form data content types (multipart/form-data, application/x-www-form-urlencoded)")
		}

		// Representation schemaId can only be specified for non form data content types (multipart/form-data, application/x-www-form-urlencoded).
		// Representation typeName can only be specified for non form data content types (multipart/form-data, application/x-www-form-urlencoded).
		// nolint gocritic
		if !contentTypeIsFormData {
			output.SchemaID = utils.String(schemaId)
			output.TypeName = utils.String(typeName)
		} else if schemaId != "" {
			return nil, fmt.Errorf("`schema_id` cannot be specified for non-form data content types (multipart/form-data, application/x-www-form-urlencoded)")
		} else if typeName != "" {
			return nil, fmt.Errorf("`type_name` cannot be specified for non-form data content types (multipart/form-data, application/x-www-form-urlencoded)")
		}

		outputs = append(outputs, output)
	}

	return &outputs, nil
}

func FlattenApiManagementOperationRepresentation(input *[]apimanagement.RepresentationContract) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	outputs := make([]interface{}, 0)

	for _, v := range *input {
		output := make(map[string]interface{})

		if v.ContentType != nil {
			output["content_type"] = *v.ContentType
		}

		output["form_parameter"] = FlattenApiManagementOperationParameterContract(v.FormParameters)

		if v.Examples != nil {
			output["example"] = FlattenApiManagementOperationParameterExampleContract(v.Examples)

			if features.ThreePointOhBeta() && v.Examples["default"] != nil && v.Examples["default"].Value != nil {
				output["sample"] = v.Examples["default"].Value.(string)
			}
		}

		if v.SchemaID != nil {
			output["schema_id"] = *v.SchemaID
		}

		if v.TypeName != nil {
			output["type_name"] = *v.TypeName
		}

		outputs = append(outputs, output)
	}

	return outputs
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
			},
		},
	}
}

func ExpandApiManagementOperationParameterContract(d *pluginsdk.ResourceData, schemaPath string, input []interface{}) *[]apimanagement.ParameterContract {
	if len(input) == 0 {
		return &[]apimanagement.ParameterContract{}
	}

	outputs := make([]apimanagement.ParameterContract, 0)

	for i, v := range input {
		vs := v.(map[string]interface{})

		name := vs["name"].(string)
		description := vs["description"].(string)
		paramType := vs["type"].(string)
		required := vs["required"].(bool)
		valuesRaw := vs["values"].(*pluginsdk.Set).List()

		output := apimanagement.ParameterContract{
			Name:         utils.String(name),
			Description:  utils.String(description),
			Type:         utils.String(paramType),
			Required:     utils.Bool(required),
			DefaultValue: nil,
			Values:       utils.ExpandStringSlice(valuesRaw),
		}

		// DefaultValue must be included in Values, else it returns error
		// when DefaultValue is unset, we need to set it nil
		// "" is a valid DefaultValue
		if v, ok := d.GetOk(fmt.Sprintf("%s.%d.default_value", schemaPath, i)); ok {
			output.DefaultValue = utils.String(v.(string))
		}
		outputs = append(outputs, output)
	}

	return &outputs
}

func FlattenApiManagementOperationParameterContract(input *[]apimanagement.ParameterContract) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	outputs := make([]interface{}, 0)
	for _, v := range *input {
		output := map[string]interface{}{}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if v.Description != nil {
			output["description"] = *v.Description
		}

		if v.Type != nil {
			output["type"] = *v.Type
		}

		if v.Required != nil {
			output["required"] = *v.Required
		}

		if v.DefaultValue != nil {
			output["default_value"] = *v.DefaultValue
		}

		output["values"] = pluginsdk.NewSet(pluginsdk.HashString, utils.FlattenStringSlice(v.Values))

		outputs = append(outputs, output)
	}

	return outputs
}

func SchemaApiManagementOperationParameterExampleContract() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: !features.ThreePointOhBeta(),
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
					Type:     pluginsdk.TypeString,
					Optional: true,
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

func ExpandApiManagementOperationParameterExampleContract(input []interface{}) map[string]*apimanagement.ParameterExampleContract {
	if len(input) == 0 {
		return map[string]*apimanagement.ParameterExampleContract{}
	}

	outputs := make(map[string]*apimanagement.ParameterExampleContract)

	for _, v := range input {
		vs := v.(map[string]interface{})

		name := vs["name"].(string)

		outputs[name] = &apimanagement.ParameterExampleContract{}

		if vs["summary"] != nil {
			outputs[name].Summary = utils.String(vs["summary"].(string))
		}

		if vs["description"] != nil {
			outputs[name].Description = utils.String(vs["description"].(string))
		}

		if vs["value"] != nil {
			outputs[name].Value = utils.String(vs["value"].(string))
		}

		if vs["external_value"] != nil {
			outputs[name].ExternalValue = utils.String(vs["external_value"].(string))
		}
	}

	return outputs
}

func FlattenApiManagementOperationParameterExampleContract(input map[string]*apimanagement.ParameterExampleContract) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	outputs := make([]interface{}, 0)
	for k, v := range input {
		output := map[string]interface{}{}

		output["name"] = k

		if v.Summary != nil {
			output["summary"] = *v.Summary
		}

		if v.Description != nil {
			output["description"] = *v.Description
		}

		if v.Value != nil {
			output["value"] = v.Value.(string)
		}

		if v.ExternalValue != nil {
			output["external_value"] = *v.ExternalValue
		}

		outputs = append(outputs, output)
	}

	return outputs
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
