// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package customproviders

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/customproviders/2018-09-01-preview/customresourceprovider"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/customproviders/validate"
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
			_, err := customresourceprovider.ParseResourceProviderID(id)
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

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

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
							Default:  string(customresourceprovider.ResourceTypeRoutingProxy),
							ValidateFunc: validation.StringInSlice([]string{
								string(customresourceprovider.ResourceTypeRoutingProxy),
								string(customresourceprovider.ResourceTypeRoutingProxyCache),
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

			"tags": commonschema.TagsForceNew(),
		},
	}
}

func resourceCustomProviderCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).CustomProviders.CustomProviderClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	location := azure.NormalizeLocation(d.Get("location").(string))
	id := customresourceprovider.NewResourceProviderID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_custom_resource_provider", id.ID())
		}
	}

	provider := customresourceprovider.CustomRPManifest{
		Properties: &customresourceprovider.CustomRPManifestProperties{
			ResourceTypes: expandCustomProviderResourceType(d.Get("resource_type").(*pluginsdk.Set).List()),
			Actions:       expandCustomProviderAction(d.Get("action").(*pluginsdk.Set).List()),
			Validations:   expandCustomProviderValidation(d.Get("validation").(*pluginsdk.Set).List()),
		},
		Location: location,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, provider); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCustomProviderRead(d, meta)
}

func resourceCustomProviderRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).CustomProviders.CustomProviderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := customresourceprovider.ParseResourceProviderID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ResourceProviderName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", azure.NormalizeLocation(model.Location))

		if props := model.Properties; props != nil {
			if err := d.Set("resource_type", flattenCustomProviderResourceType(props.ResourceTypes)); err != nil {
				return fmt.Errorf("setting `resource_type`: %+v", err)
			}

			if err := d.Set("action", flattenCustomProviderAction(props.Actions)); err != nil {
				return fmt.Errorf("setting `action`: %+v", err)
			}

			if err := d.Set("validation", flattenCustomProviderValidation(props.Validations)); err != nil {
				return fmt.Errorf("setting `validation`: %+v", err)
			}
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceCustomProviderDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).CustomProviders.CustomProviderClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := customresourceprovider.ParseResourceProviderID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandCustomProviderResourceType(input []interface{}) *[]customresourceprovider.CustomRPResourceTypeRouteDefinition {
	if len(input) == 0 {
		return nil
	}
	definitions := make([]customresourceprovider.CustomRPResourceTypeRouteDefinition, 0)

	for _, v := range input {
		if v == nil {
			continue
		}

		attrs := v.(map[string]interface{})
		definitions = append(definitions, customresourceprovider.CustomRPResourceTypeRouteDefinition{
			RoutingType: utils.ToPtr(customresourceprovider.ResourceTypeRouting(attrs["routing_type"].(string))),
			Name:        attrs["name"].(string),
			Endpoint:    attrs["endpoint"].(string),
		})
	}

	return &definitions
}

func flattenCustomProviderResourceType(input *[]customresourceprovider.CustomRPResourceTypeRouteDefinition) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	definitions := make([]interface{}, 0)
	for _, v := range *input {
		definition := make(map[string]interface{})

		definition["routing_type"] = v.RoutingType
		definition["name"] = v.Name
		definition["endpoint"] = v.Endpoint

		definitions = append(definitions, definition)
	}
	return definitions
}

func expandCustomProviderAction(input []interface{}) *[]customresourceprovider.CustomRPActionRouteDefinition {
	if len(input) == 0 {
		return nil
	}
	definitions := make([]customresourceprovider.CustomRPActionRouteDefinition, 0)

	for _, v := range input {
		if v == nil {
			continue
		}

		attrs := v.(map[string]interface{})
		definitions = append(definitions, customresourceprovider.CustomRPActionRouteDefinition{
			Name:     attrs["name"].(string),
			Endpoint: attrs["endpoint"].(string),
		})
	}

	return &definitions
}

func flattenCustomProviderAction(input *[]customresourceprovider.CustomRPActionRouteDefinition) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	definitions := make([]interface{}, 0)
	for _, v := range *input {
		definition := make(map[string]interface{})

		definition["name"] = v.Name
		definition["endpoint"] = v.Endpoint

		definitions = append(definitions, definition)
	}
	return definitions
}

func expandCustomProviderValidation(input []interface{}) *[]customresourceprovider.CustomRPValidations {
	if len(input) == 0 {
		return nil
	}

	validations := make([]customresourceprovider.CustomRPValidations, 0)

	for _, v := range input {
		if v == nil {
			continue
		}

		attrs := v.(map[string]interface{})
		validations = append(validations, customresourceprovider.CustomRPValidations{
			Specification: attrs["specification"].(string),
		})
	}

	return &validations
}

func flattenCustomProviderValidation(input *[]customresourceprovider.CustomRPValidations) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	validations := make([]interface{}, 0)
	for _, v := range *input {
		validation := make(map[string]interface{})

		validation["specification"] = v.Specification

		validations = append(validations, validation)
	}
	return validations
}
