// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedapplications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applicationdefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applications"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedapplications/validate"
	resourcesParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagedApplication() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceManagedApplicationCreate,
		Read:   resourceManagedApplicationRead,
		Update: resourceManagedApplicationUpdate,
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

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"tags": commonschema.Tags(),

		"outputs": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}

	return schema
}

func resourceManagedApplicationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	if _, ok := d.GetOk("identity"); ok {
		managedApplicationIdentity, err := expandManagedApplicationIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		parameters.Identity = managedApplicationIdentity
	}

	params, err := expandManagedApplicationParameters(d)
	if err != nil {
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

func resourceManagedApplicationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := applications.ParseApplicationID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	payload := existing.Model

	if d.HasChange("application_definition_id") {
		payload.Properties.ApplicationDefinitionId = pointer.To(d.Get("application_definition_id").(string))
	}

	if d.HasChange("identity") {
		managedApplicationIdentity, err := expandManagedApplicationIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		payload.Identity = managedApplicationIdentity
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	params, err := expandManagedApplicationParameters(d)
	if err != nil {
		return fmt.Errorf("expanding `parameter_values`: %+v", err)
	}
	payload.Properties.Parameters = pointer.To(interface{}(params))

	err = client.CreateOrUpdateThenPoll(ctx, *id, *payload)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

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
			return fmt.Errorf("expanding `parameter_values`: %+v", err)
		}

		parameterValues, err := flattenManagedApplicationParameterValuesValueToString(p.Parameters, *expendedParams)
		if err != nil {
			return fmt.Errorf("serializing JSON from `parameter_values`: %+v", err)
		}
		d.Set("parameter_values", parameterValues)

		managedApplicationIdentity, err := flattenManagedApplicationIdentity(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		d.Set("identity", managedApplicationIdentity)

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
		return strconv.FormatBool(v.(bool)), nil
	case float64:
		// use precision 0 since this comes from an int
		return fmt.Sprintf("%.f", v.(float64)), nil
	case string:
		return v.(string), nil
	case map[string]interface{}:
		return compactParameterOrOutputValue(v)
	case []interface{}:
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

func expandManagedApplicationIdentity(input []interface{}) (*applications.Identity, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	resourceType := applications.ResourceIdentityType(expanded.Type)
	out := &applications.Identity{
		Type: &resourceType,
	}

	if expanded.Type == identity.TypeUserAssigned || expanded.Type == identity.TypeSystemAssignedUserAssigned {
		userAssignedIdentities := make(map[string]applications.UserAssignedResourceIdentity)
		for k := range expanded.IdentityIds {
			userAssignedIdentities[k] = applications.UserAssignedResourceIdentity{}
		}

		out.UserAssignedIdentities = &userAssignedIdentities
	}

	return out, nil
}

func flattenManagedApplicationIdentity(input *applications.Identity) ([]interface{}, error) {
	var config *identity.SystemAndUserAssignedMap

	if input != nil {
		config = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(*input.Type),
			IdentityIds: nil,
		}

		if input.PrincipalId != nil {
			config.PrincipalId = *input.PrincipalId
		}

		if input.TenantId != nil {
			config.TenantId = *input.TenantId
		}

		identityIds := make(map[string]identity.UserAssignedIdentityDetails)
		if input.UserAssignedIdentities != nil {
			for k, v := range *input.UserAssignedIdentities {
				identityIds[k] = identity.UserAssignedIdentityDetails{
					PrincipalId: v.PrincipalId,
				}
			}
		}

		config.IdentityIds = identityIds
	}

	result, err := identity.FlattenSystemAndUserAssignedMap(config)
	if err != nil {
		return nil, err
	}

	return *result, nil
}
