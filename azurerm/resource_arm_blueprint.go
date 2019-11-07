package azurerm

import (
	"bytes"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/preview/blueprint/mgmt/2018-11-01-preview/blueprint"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"log"
	"net/http"
	"strings"
)

func resourceArmBlueprint() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmBlueprintCreateOrUpdate,
		Read:   resourceArmBlueprintRead,
		Update: resourceArmBlueprintCreateOrUpdate,
		Delete: resourceArmBlueprintDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateBlueprintScope,
			},
			"properties": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"display_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parameters": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(blueprint.Array),
											string(blueprint.Bool),
											string(blueprint.Int),
											string(blueprint.Object),
											string(blueprint.SecureObject),
											string(blueprint.SecureString),
											string(blueprint.String),
										}, false),
									},
									"default_value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"display_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"resource_groups": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"location": {
										Type:             schema.TypeString,
										Optional:         true,
										StateFunc:        azure.NormalizeLocation,
										DiffSuppressFunc: azure.SuppressLocationDiff,
									},
									"display_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"tags": tags.Schema(),
								},
							},
						},
						"target_scope": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(blueprint.Subscription),
								// Target scope of managementGroup reserved for future use, currently rejected by API
								//string(blueprint.ManagementGroup),
							}, false),
						},
						"versions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
						},
						"status": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"last_modified": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"time_created": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceArmBlueprintCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Blueprint.BlueprintsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	bpType := d.Get("type").(string)
	properties := expandBlueprintProperties(d)

	scope := d.Get("scope").(string)
	model := blueprint.Model{
		Properties: properties,
		Name:       utils.String(name),
		Type:       utils.String(bpType),
		ID:         utils.String(""),
	}

	read, err := client.CreateOrUpdate(ctx, scope, name, model)

	if err != nil {
		return fmt.Errorf("Error creating or updating blueprint: %+v", err)
	}

	d.SetId(*read.ID)
	return resourceArmBlueprintRead(d, meta)
}

func resourceArmBlueprintRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Blueprint.BlueprintsClient
	ctx := meta.(*ArmClient).StopContext

	// Can't use ParseAzureResourceID as normal, as management group id value doesn't start "/subscriptions"
	scope := d.Get("scope").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, scope, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] blueprint %q was not found in scope %q", name, scope)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading blueprint %q in scope %q", name, scope)
	}

	d.Set("name", resp.Name)
	d.Set("type", resp.Type)
	d.Set("scope", scope)

	properties := flattenBlueprintProperties(resp.Properties)
	err = d.Set("properties", properties)

	return nil
}

func resourceArmBlueprintDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Blueprint.BlueprintsClient
	ctx := meta.(*ArmClient).StopContext

	scope := d.Get("scope").(string)
	name := d.Get("name").(string)
	resp, err := client.Delete(ctx, scope, name)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting Blueprint %s in scope %s", name, scope)
	}
	return nil
}

func expandBlueprintProperties(d *schema.ResourceData) *blueprint.Properties {

	ret := blueprint.Properties{}

	v := d.Get("properties")
	props := v.(*schema.Set).List()
	for key := range props {
		prop := props[key].(map[string]interface{})
		if displayName, ok := prop["display_name"].(string); ok {
			ret.DisplayName = &displayName
		}

		if description, ok := prop["description"].(string); ok {
			ret.Description = &description
		}

		if layout, ok := prop["layout"].([]interface{}); ok {
			ret.Layout = layout
		} else {
			ret.Layout = map[string]string{}
		}

		if ts, ok := prop["target_scope"].(string); ok {
			switch ts {
			case "subscription":
				ret.TargetScope = blueprint.Subscription
			case "managementGroup":
				ret.TargetScope = blueprint.ManagementGroup
			}
		}
		if _, ok := d.GetOk("properties"); ok {
			ret.Parameters = expandBlueprintPropertiesParameters(d)
			ret.ResourceGroups = expandBlueprintPropertiesResourceGroups(d)
		}

		if _, ok := prop["versions"].(map[string]*blueprint.Status); ok {
			// todo - handle Versions object when I figure out structure
			ret.Versions = map[string]string{}
		} else {
			ret.Versions = map[string]string{}
		}
	}

	return &ret
}

func expandBlueprintPropertiesParameters(d *schema.ResourceData) map[string]*blueprint.ParameterDefinition {
	properties := d.Get("properties").(*schema.Set).List()
	parametersRaw := properties[0].(map[string]interface{})
	parameters := parametersRaw["parameters"].(*schema.Set).List()

	blueprintParameters := make(map[string]*blueprint.ParameterDefinition)
	for _, v := range parameters {
		param := v.(map[string]interface{})
		name := param["name"].(string)
		displayName := param["display_name"].(string)
		desc := param["description"].(string)

		paramMeta := &blueprint.ParameterDefinitionMetadata{
			Description: &desc,
			DisplayName: &displayName,
		}

		paramTypeRaw := param["type"].(string)
		//var paramType blueprint.TemplateParameterType
		paramType := stringToTemplateParameterType(paramTypeRaw)

		defaultValue := param["default_value"].(interface{})

		p := &blueprint.ParameterDefinition{
			Type:         paramType,
			DefaultValue: defaultValue,
		}
		p.ParameterDefinitionMetadata = paramMeta

		blueprintParameters[name] = p
	}

	return blueprintParameters
}

func expandBlueprintPropertiesResourceGroups(d *schema.ResourceData) map[string]*blueprint.ResourceGroupDefinition {
	properties := d.Get("properties").(*schema.Set).List()
	resourceGroupsRaw := properties[0].(map[string]interface{})
	resourceGroups := resourceGroupsRaw["resource_groups"].(*schema.Set).List()

	blueprintResourceGroups := make(map[string]*blueprint.ResourceGroupDefinition)
	for _, v := range resourceGroups {
		resourceGroup := v.(map[string]interface{})
		name := resourceGroup["name"].(string)
		location := resourceGroup["location"].(string)
		displayName := resourceGroup["display_name"].(string)
		description := resourceGroup["description"].(string)

		rgMetaData := &blueprint.ParameterDefinitionMetadata{
			Description: &description,
			DisplayName: &displayName,
		}

		//tagsRaw := resourceGroup["tags"].(map[string]interface{})

		rg := &blueprint.ResourceGroupDefinition{
			Name:     &name,
			Location: &location,
			//Tags: tags,
		}
		rg.ParameterDefinitionMetadata = rgMetaData

		blueprintResourceGroups[name] = rg

	}
	return blueprintResourceGroups
}

func flattenBlueprintProperties(blueprintProperties *blueprint.Properties) *schema.Set {

	properties := &schema.Set{
		F: resourceArmBlueprintPropertiesHash,
	}

	props := make(map[string]interface{}, 0)

	props["display_name"] = blueprintProperties.DisplayName
	props["description"] = blueprintProperties.Description
	props["target_scope"] = string(blueprintProperties.TargetScope)

	props["status"] = flattenBlueprintPropertiesStatus(blueprintProperties.Status)

	if blueprintProperties.Parameters != nil {
		for parameterName, blueprintParameter := range blueprintProperties.Parameters {
			params := make(map[string]interface{})
			params["name"] = parameterName
			params["type"] = string(blueprintParameter.Type)
			params["default_value"] = blueprintParameter.DefaultValue
			params["allowed_values"] = &blueprintParameter.AllowedValues
			params["display_name"] = blueprintParameter.DisplayName
			params["description"] = blueprintParameter.Description
		}
	}

	props["parameters"] = flattenBlueprintPropertiesParameters(blueprintProperties.Parameters)

	props["resource_groups"] = flattenBlueprintPropertiesResourceGroups(blueprintProperties.ResourceGroups)

	properties.Add(props)
	return properties
}

func flattenBlueprintPropertiesResourceGroups(input map[string]*blueprint.ResourceGroupDefinition) *schema.Set {
	resourceGroups := &schema.Set{
		F: resourceArmBlueprintPropertiesResourceGroupHash,
	}

	//rgs := make(map[string]interface{}, 0)

	if input != nil {
		for _, resourceGroupDefinition := range input {
			rg := make(map[string]interface{})
			rg["name"] = resourceGroupDefinition.Name
			rg["location"] = resourceGroupDefinition.Location
			rg["tags"] = resourceGroupDefinition.Tags
			rg["display_name"] = resourceGroupDefinition.DisplayName
			rg["description"] = resourceGroupDefinition.Description
			rg["depends_on"] = resourceGroupDefinition.DependsOn

			resourceGroups.Add(rg)
		}
	}
	return resourceGroups
}

func flattenBlueprintPropertiesStatus(input *blueprint.Status) *schema.Set {
	status := &schema.Set{
		F: resourceArmBlueprintPropertiesStatusHash,
	}
	stat := make(map[string]interface{})

	stat["last_modified"] = input.LastModified.String()
	stat["time_created"] = input.TimeCreated.String()
	status.Add(stat)

	return status
}

func flattenBlueprintPropertiesParameters(input map[string]*blueprint.ParameterDefinition) *schema.Set {
	parameters := &schema.Set{
		F: resourceBlueprintPropertiesParametersHash,
	}

	if input != nil {
		for name, parameter := range input {
			param := make(map[string]interface{})
			param["name"] = name
			param["display_name"] = parameter.DisplayName
			param["type"] = parameter.Type
			param["default_value"] = parameter.DefaultValue.(string)
			param["allowed_values"] = parameter.AllowedValues
			param["description"] = parameter.Description

			parameters.Add(param)
		}
	}

	return parameters
}

func templateParameterTypeToString(t blueprint.TemplateParameterType) (param string) {

	return string(t)
}

func stringToTemplateParameterType(t string) (paramType blueprint.TemplateParameterType) {
	switch t {
	case "string":
		paramType = blueprint.String
	case "array":
		paramType = blueprint.Array
	case "bool":
		paramType = blueprint.Bool
	case "int":
		paramType = blueprint.Int
	case "object":
		paramType = blueprint.Object
	case "secureObject":
		paramType = blueprint.SecureObject
	case "secureString":
		paramType = blueprint.SecureString
	}
	return paramType
}

func validateBlueprintScope(i interface{}, k string) (warnings []string, errors []error) {
	input := i.(string)

	if strings.HasPrefix(input, "/subscription") {
		_, err := azure.ValidateResourceID(i, input)
		if err != nil {
			errors = append(errors, fmt.Errorf("Subscription specified is not a valid Resource ID: %q", k))
		}
	} else if strings.HasPrefix(input, "/providers/Microsoft.Management/managementGroups/") {
		input = strings.TrimPrefix(input, "/")
		input = strings.TrimSuffix(input, "/")
		components := strings.Split(input, "/")

		if len(components) != 4 {
			errors = append(errors, fmt.Errorf("Invalid management group path, should contain 4 elements not %q", len(components)))
		}
		_, err := validate.UUID(components[3], input)
		if err != nil {
			errors = append(errors, fmt.Errorf("Management group ID not a valid uuid: %q", components[3]))
		}
	} else {
		errors = append(errors, fmt.Errorf("Invalid scope, should be a subscription resource ID or Management Group path: %q", k))
	}

	return warnings, errors
}

func resourceArmBlueprintPropertiesHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		if h, ok := m["display_name"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", h))
		}
	}
	return hashcode.String(buf.String())
}

func resourceArmBlueprintPropertiesResourceGroupHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"]))
		buf.WriteString(fmt.Sprintf("%s-", m["location"]))
		buf.WriteString(fmt.Sprintf("%s-", m["display_name"]))
		buf.WriteString(fmt.Sprintf("%s-", m["description"]))
	}
	return hashcode.String(buf.String())
}

func resourceArmBlueprintPropertiesStatusHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["time_created"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["last_modified"].(string)))
	}
	return hashcode.String(buf.String())
}

func resourceBlueprintPropertiesParametersHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", templateParameterTypeToString(m["type"].(blueprint.TemplateParameterType))))
		buf.WriteString(fmt.Sprintf("%s-", m["display_name"]))
		buf.WriteString(fmt.Sprintf("%s-", m["description"]))
		// todo deal with allowed and default values objects
		//buf.WriteString(fmt.Sprintf("%s-", m["allowed_values"].(string)))
		//buf.WriteString(fmt.Sprintf("%s-", m["default_value"].(string)))
	}
	return hashcode.String(buf.String())
}
