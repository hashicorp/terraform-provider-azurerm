package managedapplications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applicationdefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applications"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
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

		Schema: map[string]*pluginsdk.Schema{
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

			"parameters": {
				Type:          pluginsdk.TypeMap,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"parameter_values"},
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"parameter_values": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
				ConflictsWith:    []string{"parameters"},
			},

			"plan": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"product": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"publisher": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"version": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"promotion_code": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
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
		},
	}
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
		Location: pointer.To(azure.NormalizeLocation(d.Get("location"))),
		Kind:     d.Get("kind").(string),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("managed_resource_group_name"); ok {
		targetResourceGroupId := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", meta.(*clients.Client).Account.SubscriptionId, v)
		parameters.Properties = applications.ApplicationProperties{
			ManagedResourceGroupId: pointer.To(targetResourceGroupId),
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
		return fmt.Errorf("expanding `parameters` or `parameter_values`: %+v", err)
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

	if m := resp.Model; m != nil {
		p := m.Properties

		d.Set("name", m.Name)
		d.Set("resource_group_name", id.ResourceGroupName)
		d.Set("location", location.NormalizeNilable(m.Location))
		d.Set("kind", m.Kind)
		if err := d.Set("plan", flattenManagedApplicationPlan(m.Plan)); err != nil {
			return fmt.Errorf("setting `plan`: %+v", err)
		}

		id, err := resourcesParse.ResourceGroupIDInsensitively(*p.ManagedResourceGroupId)
		if err != nil {
			return err
		}

		d.Set("managed_resource_group_name", id.ResourceGroup)
		d.Set("application_definition_id", p.ApplicationDefinitionId)

		parameterValues, err := flattenManagedApplicationParameterValuesValueToString(p.Parameters)
		if err != nil {
			return fmt.Errorf("serializing JSON from `parameter_values`: %+v", err)
		}
		d.Set("parameter_values", parameterValues)

		parameters, err := flattenManagedApplicationParametersOrOutputs(p.Parameters)
		if err != nil {
			return err
		}
		if err = d.Set("parameters", parameters); err != nil {
			return err
		}

		outputs, err := flattenManagedApplicationParametersOrOutputs(p.Outputs)
		if err != nil {
			return err
		}
		if err = d.Set("outputs", outputs); err != nil {
			return err
		}

		return tags.FlattenAndSet(d, m.Tags)
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

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("failed to delete Managed Application %s: %+v", id, err)
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

	if v, ok := d.GetOk("parameters"); ok {
		params := v.(map[string]interface{})

		for key, val := range params {
			newParams[key] = struct {
				Value interface{} `json:"value"`
			}{
				Value: val,
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

func flattenManagedApplicationParametersOrOutputs(input interface{}) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	if input == nil {
		return results, nil
	}

	for k, val := range input.(map[string]interface{}) {
		mapVal, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected managed application parameter or output type: %+v", mapVal)
		}
		if mapVal != nil {
			v, ok := mapVal["value"]
			if !ok {
				return nil, fmt.Errorf("missing key 'value' in parameters or output map %+v", mapVal)
			}
			switch t := v.(type) {
			case float64:
				results[k] = v.(float64)
			case string:
				results[k] = v.(string)
			case map[string]interface{}:
				// Azure NVA managed applications read call returns empty map[string]interface{} parameter 'tags'
				// Do not return an error if the parameter is unsupported type, but is empty
				if len(v.(map[string]interface{})) == 0 {
					log.Printf("parameter '%s' is unexpected type %T, but we're ignoring it because of the empty value", k, t)
				} else {
					return nil, fmt.Errorf("unexpected parameter type %T", t)
				}
			default:
				return nil, fmt.Errorf("unexpected parameter type %T", t)
			}
		}
	}

	return results, nil
}

func flattenManagedApplicationParameterValuesValueToString(input interface{}) (string, error) {
	if input == nil {
		return "", nil
	}

	for k, v := range input.(map[string]interface{}) {
		if v != nil {
			delete(input.(map[string]interface{})[k].(map[string]interface{}), "type")
		}
	}

	result, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	compactJson := bytes.Buffer{}
	if err := json.Compact(&compactJson, result); err != nil {
		return "", err
	}

	return compactJson.String(), nil
}
