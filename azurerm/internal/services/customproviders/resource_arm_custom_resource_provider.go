package customproviders

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/customproviders/mgmt/2018-09-01-preview/customproviders"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/customproviders/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCustomResourceProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCustomResourceProviderCreateUpdate,
		Read:   resourceArmCustomResourceProviderRead,
		Update: resourceArmCustomResourceProviderCreateUpdate,
		Delete: resourceArmCustomResourceProviderDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.CustomResourceProviderID(id)
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
				ValidateFunc: validateCustomResourceProviderName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"resource_type": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"resource_type", "action"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"endpoint": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.URLIsHTTPS,
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
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"resource_type", "action"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"endpoint": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.URLIsHTTPS,
						},
						"routing_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(customproviders.ResourceTypeRoutingProxy),
							ValidateFunc: validation.StringInSlice([]string{
								string(customproviders.ResourceTypeRoutingProxy),
							}, false),
						},
					},
				},
			},

			"validation": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"specification": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(customproviders.Swagger),
							ValidateFunc: validation.StringInSlice([]string{
								string(customproviders.Swagger),
							}, false),
						},
					},
				},
			},

			"tags": tags.ForceNewSchema(),
		},
	}
}

func resourceArmCustomResourceProviderCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).CustomProviders.CustomResourceProviderClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Custom Resource Provider %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_custom_resource_provider", *existing.ID)
		}
	}

	provider := customproviders.CustomRPManifest{
		CustomRPManifestProperties: &customproviders.CustomRPManifestProperties{
			ResourceTypes: expandCustomResourceProviderResourceType(d.Get("resource_type").([]interface{})),
			Actions:       expandCustomResourceProviderAction(d.Get("action").([]interface{})),
			Validations:   expandCustomResourceProviderValidation(d.Get("validation").([]interface{})),
		},
		Location: &location,
		Tags:     tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, provider)
	if err != nil {
		return fmt.Errorf("creating/updating Custom Resource Provider %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of Custom Resource Provider %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Custom Resource Provider %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("cannot read Custom Resource Provider %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)
	return resourceArmCustomResourceProviderRead(d, meta)
}

func resourceArmCustomResourceProviderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).CustomProviders.CustomResourceProviderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CustomResourceProviderID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Custom Resource Provider %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
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
func resourceArmCustomResourceProviderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).CustomProviders.CustomResourceProviderClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CustomResourceProviderID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Custom Resource Provider %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Custom Resource Provider %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandCustomResourceProviderResourceType(input []interface{}) *[]customproviders.CustomRPResourceTypeRouteDefinition {
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
		definition := make(map[string]interface{}, 0)

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

func expandCustomResourceProviderAction(input []interface{}) *[]customproviders.CustomRPActionRouteDefinition {
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
			RoutingType: customproviders.ActionRouting(attrs["routing_type"].(string)),
			Name:        utils.String(attrs["name"].(string)),
			Endpoint:    utils.String(attrs["endpoint"].(string)),
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
		definition := make(map[string]interface{}, 0)

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

func expandCustomResourceProviderValidation(input []interface{}) *[]customproviders.CustomRPValidations {
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
			ValidationType: customproviders.ValidationType(attrs["type"].(string)),
			Specification:  utils.String(attrs["specification"].(string)),
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
		validation := make(map[string]interface{}, 0)

		validation["type"] = string(v.ValidationType)

		if v.Specification != nil {
			validation["specification"] = *v.Specification
		}

		validations = append(validations, validation)
	}
	return validations
}

func validateCustomResourceProviderName(v interface{}, k string) (warnings []string, errors []error) {
	name := v.(string)

	if regexp.MustCompile(`^[\s]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q must not consist of whitespace", k))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q may only contain letters and digits: %q", k, name))
	}

	if len(name) < 3 || len(name) > 63 {
		errors = append(errors, fmt.Errorf("%q must be (inclusive) between 3 and 63 characters long but is %d", k, len(name)))
	}

	return warnings, errors
}
