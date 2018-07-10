package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2016-09-01/web"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppServicePlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServicePlanCreateUpdate,
		Read:   resourceArmAppServicePlanRead,
		Update: resourceArmAppServicePlanCreateUpdate,
		Delete: resourceArmAppServicePlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAppServicePlanName,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"kind": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Windows",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					// @tombuildsstuff: I believe `app` is the older representation of `Windows`
					// thus we need to support it to be able to import resources without recreating them.
					"App",
					"FunctionApp",
					"Linux",
					"Windows",
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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

			"properties": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_service_environment_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"reserved": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"per_site_scaling": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"maximum_number_of_workers": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmAppServicePlanCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicePlansClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM App Service Plan creation.")

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	kind := d.Get("kind").(string)
	tags := d.Get("tags").(map[string]interface{})

	sku := expandAzureRmAppServicePlanSku(d)
	properties := expandAppServicePlanProperties(d, name)

	appServicePlan := web.AppServicePlan{
		Location:                 &location,
		AppServicePlanProperties: properties,
		Kind: &kind,
		Tags: expandTags(tags),
		Sku:  &sku,
	}

	createFuture, err := client.CreateOrUpdate(ctx, resGroup, name, appServicePlan)
	if err != nil {
		return err
	}

	err = createFuture.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read AzureRM App Service Plan %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppServicePlanRead(d, meta)
}

func resourceArmAppServicePlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicePlansClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading Azure App Service Plan %s", id)

	resGroup := id.ResourceGroup
	name := id.Path["serverfarms"]

	ctx := meta.(*ArmClient).StopContext
	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure App Service Plan %s: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}
	d.Set("kind", resp.Kind)

	if props := resp.AppServicePlanProperties; props != nil {
		d.Set("properties", flattenAppServiceProperties(props))

		if props.MaximumNumberOfWorkers != nil {
			d.Set("maximum_number_of_workers", int(*props.MaximumNumberOfWorkers))
		}
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", flattenAppServicePlanSku(sku))
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmAppServicePlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicePlansClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["serverfarms"]

	log.Printf("[DEBUG] Deleting app service plan %s: %s", resGroup, name)

	resp, err := client.Delete(ctx, resGroup, name)

	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return err
	}

	return nil
}

func expandAzureRmAppServicePlanSku(d *schema.ResourceData) web.SkuDescription {
	configs := d.Get("sku").([]interface{})
	config := configs[0].(map[string]interface{})

	tier := config["tier"].(string)
	size := config["size"].(string)

	sku := web.SkuDescription{
		Name: &size,
		Tier: &tier,
		Size: &size,
	}

	if v, ok := config["capacity"]; ok {
		capacity := v.(int)
		sku.Capacity = utils.Int32(int32(capacity))
	}

	return sku
}

func flattenAppServicePlanSku(profile *web.SkuDescription) []interface{} {
	skus := make([]interface{}, 0, 1)
	sku := make(map[string]interface{}, 2)

	sku["tier"] = *profile.Tier
	sku["size"] = *profile.Size

	if profile.Capacity != nil {
		sku["capacity"] = *profile.Capacity
	}

	skus = append(skus, sku)

	return skus
}

func expandAppServicePlanProperties(d *schema.ResourceData, name string) *web.AppServicePlanProperties {
	configs := d.Get("properties").([]interface{})
	properties := web.AppServicePlanProperties{
		Name: &name,
	}
	if len(configs) == 0 {
		return &properties
	}
	config := configs[0].(map[string]interface{})

	appServiceEnvironmentId := config["app_service_environment_id"].(string)
	if appServiceEnvironmentId != "" {
		properties.HostingEnvironmentProfile = &web.HostingEnvironmentProfile{
			ID: utils.String(appServiceEnvironmentId),
		}
	}

	perSiteScaling := config["per_site_scaling"].(bool)
	properties.PerSiteScaling = utils.Bool(perSiteScaling)

	reserved := config["reserved"].(bool)
	properties.Reserved = utils.Bool(reserved)

	return &properties
}

func flattenAppServiceProperties(props *web.AppServicePlanProperties) []interface{} {
	result := make([]interface{}, 0, 1)
	properties := make(map[string]interface{}, 0)

	if props.HostingEnvironmentProfile != nil {
		properties["app_service_environment_id"] = *props.HostingEnvironmentProfile.ID
	}

	if props.PerSiteScaling != nil {
		properties["per_site_scaling"] = *props.PerSiteScaling
	}

	if props.Reserved != nil {
		properties["reserved"] = *props.Reserved
	}

	result = append(result, properties)
	return result
}

func validateAppServicePlanName(v interface{}, k string) (ws []string, es []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[0-9a-zA-Z-_]{1,60}$`).Match([]byte(value)); !matched {
		es = append(es, fmt.Errorf("%q may only contain alphanumeric characters, dashes and underscores up to 60 characters in length", k))
	}

	return
}
