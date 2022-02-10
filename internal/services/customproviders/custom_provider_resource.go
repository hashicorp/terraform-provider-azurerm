package customproviders

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/customproviders/mgmt/2018-09-01-preview/customproviders"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/customproviders/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/customproviders/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCustomProvider() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCustomProviderCreateUpdate,
		Read:   resourceCustomProviderRead,
		Update: resourceCustomProviderCreateUpdate,
		Delete: resourceCustomProviderDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ResourceProviderID(id)
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
				ValidateFunc: validate.CustomProviderName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"resource_type": {
				Type:         pluginsdk.TypeSet,
				Optional:     true,
				AtLeastOneOf: []string{"resource_type", "action"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"endpoint": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsURLWithHTTPS,
						},
						"routing_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(customproviders.ResourceTypeRoutingProxy),
							ValidateFunc: validation.StringInSlice([]string{
								string(customproviders.ResourceTypeRoutingProxy),
								string(customproviders.ResourceTypeRoutingProxyCache),
							}, false),
						},
					},
				},
			},

			"action": {
				Type:         pluginsdk.TypeSet,
				Optional:     true,
				AtLeastOneOf: []string{"resource_type", "action"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"endpoint": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsURLWithHTTPS,
						},
					},
				},
			},

			"validation": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"specification": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsURLWithHTTPS,
						},
					},
				},
			},

			"tags": tags.ForceNewSchema(),
		},
	}
}

func resourceCustomProviderCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).CustomProviders.CustomProviderClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	location := azure.NormalizeLocation(d.Get("location").(string))
	id := parse.NewResourceProviderID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_custom_resource_provider", id.ID())
		}
	}

	provider := customproviders.CustomRPManifest{
		CustomRPManifestProperties: &customproviders.CustomRPManifestProperties{
			ResourceTypes: expandCustomProviderResourceType(d.Get("resource_type").(*pluginsdk.Set).List()),
			Actions:       expandCustomProviderAction(d.Get("action").(*pluginsdk.Set).List()),
			Validations:   expandCustomProviderValidation(d.Get("validation").(*pluginsdk.Set).List()),
		},
		Location: &location,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, provider)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCustomProviderRead(d, meta)
}

func resourceCustomProviderRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).CustomProviders.CustomProviderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceProviderID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Custom Provider %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if err := d.Set("resource_type", flattenCustomProviderResourceType(resp.ResourceTypes)); err != nil {
		return fmt.Errorf("setting `resource_type`: %+v", err)
	}

	if err := d.Set("action", flattenCustomProviderAction(resp.Actions)); err != nil {
		return fmt.Errorf("setting `action`: %+v", err)
	}

	if err := d.Set("validation", flattenCustomProviderValidation(resp.Validations)); err != nil {
		return fmt.Errorf("setting `validation`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceCustomProviderDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).CustomProviders.CustomProviderClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceProviderID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Custom Provider %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Custom Provider %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandCustomProviderResourceType(input []interface{}) *[]customproviders.CustomRPResourceTypeRouteDefinition {
	if len(input) == 0 {
		return nil
	}
	definitions := make([]customproviders.CustomRPResourceTypeRouteDefinition, 0)

	for _, v := range input {
		if v == nil {
			continue
		}

		attrs := v.(map[string]interface{})
		definitions = append(definitions, customproviders.CustomRPResourceTypeRouteDefinition{
			RoutingType: customproviders.ResourceTypeRouting(attrs["routing_type"].(string)),
			Name:        utils.String(attrs["name"].(string)),
			Endpoint:    utils.String(attrs["endpoint"].(string)),
		})
	}

	return &definitions
}

func flattenCustomProviderResourceType(input *[]customproviders.CustomRPResourceTypeRouteDefinition) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	definitions := make([]interface{}, 0)
	for _, v := range *input {
		definition := make(map[string]interface{})

		definition["routing_type"] = string(v.RoutingType)

		if v.Name != nil {
			definition["name"] = *v.Name
		}

		if v.Endpoint != nil {
			definition["endpoint"] = *v.Endpoint
		}

		definitions = append(definitions, definition)
	}
	return definitions
}

func expandCustomProviderAction(input []interface{}) *[]customproviders.CustomRPActionRouteDefinition {
	if len(input) == 0 {
		return nil
	}
	definitions := make([]customproviders.CustomRPActionRouteDefinition, 0)

	for _, v := range input {
		if v == nil {
			continue
		}

		attrs := v.(map[string]interface{})
		definitions = append(definitions, customproviders.CustomRPActionRouteDefinition{
			Name:     utils.String(attrs["name"].(string)),
			Endpoint: utils.String(attrs["endpoint"].(string)),
		})
	}

	return &definitions
}

func flattenCustomProviderAction(input *[]customproviders.CustomRPActionRouteDefinition) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	definitions := make([]interface{}, 0)
	for _, v := range *input {
		definition := make(map[string]interface{})

		if v.Name != nil {
			definition["name"] = *v.Name
		}

		if v.Endpoint != nil {
			definition["endpoint"] = *v.Endpoint
		}

		definitions = append(definitions, definition)
	}
	return definitions
}

func expandCustomProviderValidation(input []interface{}) *[]customproviders.CustomRPValidations {
	if len(input) == 0 {
		return nil
	}

	validations := make([]customproviders.CustomRPValidations, 0)

	for _, v := range input {
		if v == nil {
			continue
		}

		attrs := v.(map[string]interface{})
		validations = append(validations, customproviders.CustomRPValidations{
			Specification: utils.String(attrs["specification"].(string)),
		})
	}

	return &validations
}

func flattenCustomProviderValidation(input *[]customproviders.CustomRPValidations) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	validations := make([]interface{}, 0)
	for _, v := range *input {
		validation := make(map[string]interface{})

		if v.Specification != nil {
			validation["specification"] = *v.Specification
		}

		validations = append(validations, validation)
	}
	return validations
}
