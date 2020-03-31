package managedapplications

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-07-01/managedapplications"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managedapplications/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managedapplications/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceManagedApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceManagedApplicationCreateUpdate,
		Read:   resourceManagedApplicationRead,
		Update: resourceManagedApplicationCreateUpdate,
		Delete: resourceManagedApplicationDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ManagedApplicationID(id)
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
				ValidateFunc: validate.ManagedApplicationName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"kind": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MarketPlace",
					"ServiceCatalog",
				}, false),
			},

			"managed_resource_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: resource.ValidateResourceGroupID,
			},

			"application_definition_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ManagedApplicationDefinitionID,
			},

			"parameters": {
				Type:             schema.TypeString,
				Optional:         true,
				StateFunc:        azure.NormalizeJson,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			"plan": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"product": {
							Type:     schema.TypeString,
							Required: true,
						},
						"publisher": {
							Type:     schema.TypeString,
							Required: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
						},
						"promotion_code": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceManagedApplicationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroupName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("failed to check for present of existing Managed Application Name %q (Resource Group %q): %+v", name, resourceGroupName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_managed_application", *existing.ID)
		}
	}

	parameters := managedapplications.Application{
		Location: utils.String(azure.NormalizeLocation(d.Get("location"))),
		Kind:     utils.String(d.Get("kind").(string)),
		ApplicationProperties: &managedapplications.ApplicationProperties{
			ManagedResourceGroupID: utils.String(d.Get("managed_resource_group_id").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("application_definition_id"); ok {
		parameters.ApplicationDefinitionID = utils.String(v.(string))
	}
	if v, ok := d.GetOk("plan"); ok {
		parameters.Plan = expandManagedApplicationPlan(v.([]interface{}))
	}
	if v, ok := d.GetOk("parameters"); ok {
		expandedParams, err := structure.ExpandJsonFromString(v.(string))
		if err != nil {
			return fmt.Errorf("unable to parse parameters: %s", err)
		}

		parameters.Parameters = &expandedParams
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroupName, name, parameters)
	if err != nil {
		return fmt.Errorf("failed to create Managed Application %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("failed to wait for creation of Managed Application %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	resp, err := client.Get(ctx, resourceGroupName, name)
	if err != nil {
		return fmt.Errorf("failed to retrieve Managed Application %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("cannot read Managed Application %q (Resource Group %q) ID", name, resourceGroupName)
	}
	d.SetId(*resp.ID)

	return resourceManagedApplicationRead(d, meta)
}

func resourceManagedApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedApplicationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Managed Application %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read Managed Application %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("kind", resp.Kind)
	if err := d.Set("plan", flattenManagedApplicationPlan(resp.Plan)); err != nil {
		return fmt.Errorf("setting `plan`: %+v", err)
	}
	if props := resp.ApplicationProperties; props != nil {
		d.Set("managed_resource_group_id", props.ManagedResourceGroupID)
		d.Set("application_definition_id", props.ApplicationDefinitionID)

		if params := props.Parameters; params != nil {
			paramsVal := params.(map[string]interface{})
			for _, v := range paramsVal {
				if v != nil {
					delete(v.(map[string]interface{}), "type")
				}
			}
			json, err := structure.FlattenJsonToString(paramsVal)
			if err != nil {
				return fmt.Errorf("failed to serialize JSON from Parameters: %+v", err)
			}

			d.Set("parameters", json)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceManagedApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedApplicationID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("failed to delete Managed Application %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("failed to wait for deleting Managed Application (Managed Application Name %q / Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandManagedApplicationPlan(input []interface{}) *managedapplications.Plan {
	if len(input) == 0 {
		return nil
	}
	plan := input[0].(map[string]interface{})

	return &managedapplications.Plan{
		Name:          utils.String(plan["name"].(string)),
		Product:       utils.String(plan["product"].(string)),
		Publisher:     utils.String(plan["publisher"].(string)),
		Version:       utils.String(plan["version"].(string)),
		PromotionCode: utils.String(plan["promotion_code"].(string)),
	}
}

func flattenManagedApplicationPlan(input *managedapplications.Plan) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	name := ""
	if input.Name != nil {
		name = *input.Name
	}
	product := ""
	if input.Product != nil {
		product = *input.Product
	}
	publisher := ""
	if input.Publisher != nil {
		publisher = *input.Publisher
	}
	version := ""
	if input.Version != nil {
		version = *input.Version
	}
	promotionCode := ""
	if input.PromotionCode != nil {
		promotionCode = *input.PromotionCode
	}

	results = append(results, map[string]interface{}{
		"name":           name,
		"product":        product,
		"publisher":      publisher,
		"version":        version,
		"promotion_code": promotionCode,
	})

	return results
}
