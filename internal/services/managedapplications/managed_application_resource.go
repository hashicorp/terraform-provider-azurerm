// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedapplications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applicationdefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applications"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedapplications/validate"
	resourcesParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagedApplication() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceManagedApplicationCreateUpdate,
		Read:   resourceManagedApplicationRead,
		Update: resourceManagedApplicationCreateUpdate,
		Delete: resourceManagedApplicationDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := applications.ParseApplicationID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceManagedApplicationSchema(),
	}
}

func resourceManagedApplicationSchema() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ApplicationName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"kind": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"MarketPlace",
				"ServiceCatalog",
			}, false),
		},

		"managed_resource_group_name": commonschema.ResourceGroupName(),

		"application_definition_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: applicationdefinitions.ValidateApplicationDefinitionID,
		},

		"parameter_values": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			ConflictsWith: func() []string {
				if !features.FourPointOhBeta() {
					return []string{"parameters"}
				}
				return []string{}
			}(),
		},

		"plan": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"product": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"publisher": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"version": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"promotion_code": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"outputs": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}

	if !features.FourPointOhBeta() {
		schema["parameters"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeMap,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"parameter_values"},
			Deprecated:    "This property has been deprecated in favour of `parameter_values`",
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		}
	}

	return schema
}

func resourceManagedApplicationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := applications.NewApplicationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("failed to check for presence of %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_managed_application", id.ID())
		}
	}

	parameters := applications.Application{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Kind:     d.Get("kind").(string),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("managed_resource_group_name"); ok {
		targetResourceGroupId := commonids.NewResourceGroupID(meta.(*clients.Client).Account.SubscriptionId, v.(string))
		parameters.Properties = applications.ApplicationProperties{
			ManagedResourceGroupId: pointer.To(targetResourceGroupId.ID()),
		}
	}

	if v, ok := d.GetOk("application_definition_id"); ok {
		parameters.Properties.ApplicationDefinitionId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("plan"); ok {
		parameters.Plan = expandManagedApplicationPlan(v.([]interface{}))
	}

	params, err := expandManagedApplicationParameters(d)
	if err != nil {
		if !features.FourPointOhBeta() {
			return fmt.Errorf("expanding `parameters` or `parameter_values`: %+v", err)
		}
		return fmt.Errorf("expanding `parameter_values`: %+v", err)
	}
	parameters.Properties.Parameters = pointer.To(interface{}(params))

	err = client.CreateOrUpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("failed to create %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceManagedApplicationRead(d, meta)
}

func resourceManagedApplicationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := applications.ParseApplicationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Managed Application %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read Managed Application %s: %+v", id, err)
	}

	d.Set("name", id.ApplicationName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		p := model.Properties

		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("kind", model.Kind)
		if err := d.Set("plan", flattenManagedApplicationPlan(model.Plan)); err != nil {
			return fmt.Errorf("setting `plan`: %+v", err)
		}

		id, err := resourcesParse.ResourceGroupIDInsensitively(pointer.From(p.ManagedResourceGroupId))
		if err != nil {
			return err
		}

		d.Set("managed_resource_group_name", id.ResourceGroup)
		d.Set("application_definition_id", p.ApplicationDefinitionId)

		expendedParams, err := expandManagedApplicationParameters(d)
		if err != nil {
			if !features.FourPointOhBeta() {
				return fmt.Errorf("expanding `parameters` or `parameter_values`: %+v", err)
			}
			return fmt.Errorf("expanding `parameter_values`: %+v", err)
		}

		parameterValues, err := flattenManagedApplicationParameterValuesValueToString(p.Parameters, *expendedParams)
		if err != nil {
			return fmt.Errorf("serializing JSON from `parameter_values`: %+v", err)
		}
		d.Set("parameter_values", parameterValues)

		if !features.FourPointOhBeta() {
			parameters, err := flattenManagedApplicationParameters(p.Parameters, *expendedParams)
			if err != nil {
				return err
			}
			if err = d.Set("parameters", parameters); err != nil {
				return err
			}
		}

		outputs, err := flattenManagedApplicationOutputs(p.Outputs)
		if err != nil {
			return err
		}
		if err = d.Set("outputs", outputs); err != nil {
			return err
		}

		if err = tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceManagedApplicationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := applications.ParseApplicationID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("failed to delete %s: %+v", *id, err)
	}

	return nil
}

func expandManagedApplicationPlan(input []interface{}) *applications.Plan {
	if len(input) == 0 {
		return nil
	}
	plan := input[0].(map[string]interface{})

	return &applications.Plan{
		Name:          plan["name"].(string),
		Product:       plan["product"].(string),
		Publisher:     plan["publisher"].(string),
		Version:       plan["version"].(string),
		PromotionCode: pointer.To(plan["promotion_code"].(string)),
	}
}

func expandManagedApplicationParameters(d *pluginsdk.ResourceData) (*map[string]interface{}, error) {
	newParams := make(map[string]interface{})

	if v, ok := d.GetOk("parameter_values"); ok {
		if err := json.Unmarshal([]byte(v.(string)), &newParams); err != nil {
			return nil, fmt.Errorf("unmarshalling `parameter_values`: %+v", err)
		}
	}

	if !features.FourPointOhBeta() {
		// `parameters` will be available in state as well after first apply when `parameter_values` is used, so getting its value only during creation or when it is changed
		if d.IsNewResource() || d.HasChange("parameters") {
			if v, ok := d.GetOk("parameters"); ok {
				params := v.(map[string]interface{})

				for key, val := range params {
					newParamValue := make(map[string]interface{}, 1)
					newParamValue["value"] = val
					newParams[key] = newParamValue
				}
			}
		}
	}

	return &newParams, nil
}

func flattenManagedApplicationPlan(input *applications.Plan) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	results = append(results, map[string]interface{}{
		"name":           input.Name,
		"product":        input.Product,
		"publisher":      input.Publisher,
		"version":        input.Version,
		"promotion_code": pointer.From(input.PromotionCode),
	})

	return results
}

func flattenManagedApplicationParameters(input *interface{}, localParameters map[string]interface{}) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	if input == nil {
		return results, nil
	}

	attrs := *input
	if _, ok := attrs.(map[string]interface{}); ok {
		for k, val := range attrs.(map[string]interface{}) {
			mapVal, ok := val.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("unexpected managed application parameter type: %+v", mapVal)
			}
			if mapVal != nil {
				v, ok := mapVal["value"]
				if !ok {
					// Secure values are not returned, thus settings it with local value
					v = ""
					if oldValueStruct, oldValueStructOK := localParameters[k]; oldValueStructOK {
						if _, oldValueStructTypeOK := oldValueStruct.(map[string]interface{}); oldValueStructTypeOK {
							if oldValue, oldValueOK := oldValueStruct.(map[string]interface{})["value"]; oldValueOK {
								v = oldValue
							}
						}
					}
				}

				value, err := extractParameterOrOutputValue(v)
				if err != nil {
					return nil, fmt.Errorf("extracting parameters: %+v", err)
				}
				results[k] = value
			}
		}
	}

	return results, nil
}

func flattenManagedApplicationOutputs(input *interface{}) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	if input == nil {
		return results, nil
	}

	attrs := *input
	if _, ok := attrs.(map[string]interface{}); ok {
		for k, val := range attrs.(map[string]interface{}) {
			mapVal, ok := val.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("unexpected managed application output type: %+v", mapVal)
			}
			if mapVal != nil {
				v, ok := mapVal["value"]
				if !ok {
					return nil, fmt.Errorf("missing key 'value' in output map %+v", mapVal)
				}

				value, err := extractParameterOrOutputValue(v)
				if err != nil {
					return nil, fmt.Errorf("extracting outputs: %+v", err)
				}
				results[k] = value
			}
		}
	}

	return results, nil
}

func flattenManagedApplicationParameterValuesValueToString(input *interface{}, localParameters map[string]interface{}) (string, error) {
	if input == nil {
		return "", nil
	}

	attrs := *input
	if _, ok := attrs.(map[string]interface{}); ok {
		for k, v := range attrs.(map[string]interface{}) {
			if v != nil {
				delete(attrs.(map[string]interface{})[k].(map[string]interface{}), "type")

				// Secure values are not returned, thus settings it with local value
				value := attrs.(map[string]interface{})[k].(map[string]interface{})
				if _, ok := value["value"]; !ok {
					value["value"] = ""
					if localParam, localParamOK := localParameters[k]; localParamOK {
						if _, oldValueStructTypeOK := localParam.(map[string]interface{}); oldValueStructTypeOK {
							if localParamValue, localParamValueOK := localParam.(map[string]interface{})["value"]; localParamValueOK {
								value["value"] = localParamValue
							}
						}
					}
				}
			}
		}

		return compactParameterOrOutputValue(input)
	}

	return "", nil
}

func extractParameterOrOutputValue(v interface{}) (string, error) {
	switch t := v.(type) {
	case bool:
		return fmt.Sprintf("%t", v.(bool)), nil
	case float64:
		// use precision 0 since this comes from an int
		return fmt.Sprintf("%.f", v.(float64)), nil
	case string:
		return v.(string), nil
	case map[string]interface{}:
		return compactParameterOrOutputValue(v)
	default:
		return "", fmt.Errorf("unexpected type %T", t)
	}
}

func compactParameterOrOutputValue(v interface{}) (string, error) {
	result, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	compactJson := bytes.Buffer{}
	if err = json.Compact(&compactJson, result); err != nil {
		return "", err
	}
	return compactJson.String(), nil
}
