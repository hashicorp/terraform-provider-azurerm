package customproviders

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/customproviders/mgmt/2018-09-01-preview/customproviders"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/customproviders/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/customproviders/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCustomProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCustomProviderCreateUpdate,
		Read:   resourceArmCustomProviderRead,
		Update: resourceArmCustomProviderCreateUpdate,
		Delete: resourceArmCustomProviderDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.CustomProviderID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CustomProviderName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"resource_type": {
				Type:         schema.TypeSet,
				Optional:     true,
				AtLeastOneOf: []string{"resource_type", "action"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"endpoint": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsURLWithHTTPS,
						},
						"routing_type": {
							Type:     schema.TypeString,
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
				Type:         schema.TypeSet,
				Optional:     true,
				AtLeastOneOf: []string{"resource_type", "action"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"endpoint": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsURLWithHTTPS,
						},
					},
				},
			},

			"validation": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"specification": {
							Type:         schema.TypeString,
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

func resourceArmCustomProviderCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).CustomProviders.CustomProviderClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Custom Provider %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_custom_resource_provider", *existing.ID)
		}
	}

	provider := customproviders.CustomRPManifest{
		CustomRPManifestProperties: &customproviders.CustomRPManifestProperties{
			ResourceTypes: expandCustomProviderResourceType(d.Get("resource_type").(*schema.Set).List()),
			Actions:       expandCustomProviderAction(d.Get("action").(*schema.Set).List()),
			Validations:   expandCustomProviderValidation(d.Get("validation").(*schema.Set).List()),
		},
		Location: &location,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, provider)
	if err != nil {
		return fmt.Errorf("creating/updating Custom Provider %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of Custom Provider %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Custom Provider %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("cannot read Custom Provider %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)
	return resourceArmCustomProviderRead(d, meta)
}

func resourceArmCustomProviderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).CustomProviders.CustomProviderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CustomProviderID(d.Id())
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
func resourceArmCustomProviderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).CustomProviders.CustomProviderClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CustomProviderID(d.Id())
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
