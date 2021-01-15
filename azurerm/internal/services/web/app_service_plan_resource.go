package web

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceAppServicePlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppServicePlanCreateUpdate,
		Read:   resourceAppServicePlanRead,
		Update: resourceAppServicePlanCreateUpdate,
		Delete: resourceAppServicePlanDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AppServicePlanID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AppServicePlanName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"kind": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Windows",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					// @tombuildsstuff: I believe `app` is the older representation of `Windows`
					// thus we need to support it to be able to import resources without recreating them.
					// @jcorioland: new SKU and kind 'xenon' have been added for Windows Containers support
					// https://azure.microsoft.com/en-us/blog/announcing-the-public-preview-of-windows-container-support-in-azure-app-service/
					"App",
					"elastic",
					"FunctionApp",
					"Linux",
					"Windows",
					"xenon",
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tier": {
							Type:     schema.TypeString,
							Required: true,
						},
						"size": {
							Type:     schema.TypeString,
							Required: true,
						},
						"capacity": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			// / AppServicePlanProperties
			"app_service_environment_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"per_site_scaling": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"reserved": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"maximum_elastic_worker_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"maximum_number_of_workers": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"is_xenon": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceAppServicePlanCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicePlansClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM App Service Plan creation.")

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing App Service Plan %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_app_service_plan", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	kind := d.Get("kind").(string)
	t := d.Get("tags").(map[string]interface{})

	sku := expandAppServicePlanSku(d)
	properties := &web.AppServicePlanProperties{}

	isXenon := d.Get("is_xenon").(bool)
	properties.IsXenon = &isXenon

	if kind == "xenon" && !isXenon {
		return fmt.Errorf("Creating or updating App Service Plan %q (Resource Group %q): when kind is set to xenon, is_xenon property should be set to true", name, resGroup)
	}

	appServicePlan := web.AppServicePlan{
		Location:                 &location,
		Kind:                     &kind,
		Sku:                      &sku,
		Tags:                     tags.Expand(t),
		AppServicePlanProperties: properties,
	}

	if v := d.Get("app_service_environment_id").(string); v != "" {
		appServicePlan.AppServicePlanProperties.HostingEnvironmentProfile = &web.HostingEnvironmentProfile{
			ID: utils.String(v),
		}
	}

	if v := d.Get("per_site_scaling").(bool); v {
		appServicePlan.AppServicePlanProperties.PerSiteScaling = utils.Bool(v)
	}

	reserved := d.Get("reserved").(bool)
	if strings.EqualFold(kind, "Linux") && !reserved {
		return fmt.Errorf("`reserved` has to be set to true when kind is set to `Linux`")
	}

	if strings.EqualFold(kind, "Windows") && reserved {
		return fmt.Errorf("`reserved` has to be set to false when kind is set to `Windows`")
	}

	if v := d.Get("maximum_elastic_worker_count").(int); v > 0 {
		appServicePlan.AppServicePlanProperties.MaximumElasticWorkerCount = utils.Int32(int32(v))
	}

	if reserved {
		appServicePlan.AppServicePlanProperties.Reserved = utils.Bool(reserved)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, appServicePlan)
	if err != nil {
		return fmt.Errorf("Error creating/updating App Service Plan %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the create/update of App Service Plan %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving App Service Plan %q (Resource Group %q): %+v", name, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read AzureRM App Service Plan %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceAppServicePlanRead(d, meta)
}

func resourceAppServicePlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicePlansClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServicePlanID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerfarmName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Plan %q was not found in Resource Group %q - removnig from state!", id.ServerfarmName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on App Service Plan %q (Resource Group %q): %+v", id.ServerfarmName, id.ResourceGroup, err)
	}

	// A 404 doesn't error from the app service plan sdk so we'll add this check here to catch resource not found responses
	// TODO This block can be removed if https://github.com/Azure/azure-sdk-for-go/issues/5407 gets addressed.
	if utils.ResponseWasNotFound(resp.Response) {
		log.Printf("[DEBUG] App Service Plan %q was not found in Resource Group %q - removing from state!", id.ServerfarmName, id.ResourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("name", id.ServerfarmName)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("kind", resp.Kind)

	if props := resp.AppServicePlanProperties; props != nil {
		appServiceEnvironmentId := ""
		if props.HostingEnvironmentProfile != nil && props.HostingEnvironmentProfile.ID != nil {
			appServiceEnvironmentId = *props.HostingEnvironmentProfile.ID
		}
		d.Set("app_service_environment_id", appServiceEnvironmentId)

		maximumNumberOfWorkers := 0
		if props.MaximumNumberOfWorkers != nil {
			maximumNumberOfWorkers = int(*props.MaximumNumberOfWorkers)
		}
		d.Set("maximum_number_of_workers", maximumNumberOfWorkers)

		maximumElasticWorkerCount := 0
		if props.MaximumElasticWorkerCount != nil {
			maximumElasticWorkerCount = int(*props.MaximumElasticWorkerCount)
		}
		d.Set("maximum_elastic_worker_count", maximumElasticWorkerCount)

		d.Set("per_site_scaling", props.PerSiteScaling)
		d.Set("reserved", props.Reserved)
		d.Set("is_xenon", props.IsXenon)
	}

	if err := d.Set("sku", flattenAppServicePlanSku(resp.Sku)); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceAppServicePlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicePlansClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServicePlanID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting App Service Plan %q (Resource Group %q)", id.ServerfarmName, id.ResourceGroup)

	resp, err := client.Delete(ctx, id.ResourceGroup, id.ServerfarmName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting App Service Plan %q (Resource Group %q): %+v", id.ServerfarmName, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandAppServicePlanSku(d *schema.ResourceData) web.SkuDescription {
	configs := d.Get("sku").([]interface{})
	config := configs[0].(map[string]interface{})

	tier := config["tier"].(string)
	size := config["size"].(string)

	sku := web.SkuDescription{
		Name: utils.String(size),
		Tier: utils.String(tier),
		Size: utils.String(size),
	}

	if v, ok := config["capacity"]; ok {
		capacity := v.(int)
		sku.Capacity = utils.Int32(int32(capacity))
	}

	return sku
}

func flattenAppServicePlanSku(input *web.SkuDescription) []interface{} {
	outputs := make([]interface{}, 0)
	if input == nil {
		return outputs
	}

	output := make(map[string]interface{}, 2)

	if input.Tier != nil {
		output["tier"] = *input.Tier
	}

	if input.Size != nil {
		output["size"] = *input.Size
	}

	if input.Capacity != nil {
		output["capacity"] = *input.Capacity
	}

	outputs = append(outputs, output)

	return outputs
}
