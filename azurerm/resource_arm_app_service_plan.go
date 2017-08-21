package azurerm

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/web"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jen20/riviera/azure"
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"sku": {
				Type:     schema.TypeSet,
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
					},
				},
				Set: resourceAzureRMAppServicePlanSkuHash,
			},

			"properties": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"maximum_number_of_workers": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
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

			"tags": tagsSchema(),
		},
	}
}

func resourceArmAppServicePlanCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicePlansClient

	log.Printf("[INFO] preparing arguments for AzureRM App Service Plan creation.")

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := d.Get("location").(string)
	tags := d.Get("tags").(map[string]interface{})

	sku := expandAzureRmAppServicePlanSku(d)
	properties := expandAppServicePlanProperties(d)

	appServicePlan := web.AppServicePlan{
		Location:                 &location,
		AppServicePlanProperties: properties,
		Tags: expandTags(tags),
		Sku:  &sku,
	}

	_, createErr := client.CreateOrUpdate(resGroup, name, appServicePlan, make(chan struct{}))
	err := <-createErr
	if err != nil {
		return err
	}

	read, err := client.Get(resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read AzureRM App Service Plan %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppServicePlanRead(d, meta)
}

func resourceArmAppServicePlanRead(d *schema.ResourceData, meta interface{}) error {
	AppServicePlanClient := meta.(*ArmClient).appServicePlansClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading Azure App Service Plan %s", id)

	resGroup := id.ResourceGroup
	name := id.Path["serverfarms"]

	resp, err := AppServicePlanClient.Get(resGroup, name)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure App Service Plan %s: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))

	if props := resp.AppServicePlanProperties; props != nil {
		d.Set("properties", flattenAppServiceProperties(props))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", flattenAzureRmAppServicePlanSku(sku))
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmAppServicePlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicePlansClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["serverfarms"]

	log.Printf("[DEBUG] Deleting app service plan %s: %s", resGroup, name)

	_, err = client.Delete(resGroup, name)

	return err
}

func resourceAzureRMAppServicePlanSkuHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	tier := m["tier"].(string)
	size := m["size"].(string)

	buf.WriteString(fmt.Sprintf("%s-", tier))
	buf.WriteString(fmt.Sprintf("%s-", size))

	return hashcode.String(buf.String())
}

func expandAzureRmAppServicePlanSku(d *schema.ResourceData) web.SkuDescription {
	configs := d.Get("sku").(*schema.Set).List()
	config := configs[0].(map[string]interface{})

	tier := config["tier"].(string)
	size := config["size"].(string)

	sku := web.SkuDescription{
		Name: &size,
		Tier: &tier,
		Size: &size,
	}

	return sku
}

func flattenAzureRmAppServicePlanSku(profile *web.SkuDescription) *schema.Set {
	skus := &schema.Set{
		F: resourceAzureRMAppServicePlanSkuHash,
	}

	sku := make(map[string]interface{}, 2)

	sku["tier"] = *profile.Tier
	sku["size"] = *profile.Size

	skus.Add(sku)

	return skus
}

func expandAppServicePlanProperties(d *schema.ResourceData) *web.AppServicePlanProperties {
	configs := d.Get("properties").([]interface{})
	properties := web.AppServicePlanProperties{}
	if len(configs) == 0 {
		return &properties
	}
	config := configs[0].(map[string]interface{})

	perSiteScaling := config["per_site_scaling"].(bool)
	properties.PerSiteScaling = azure.Bool(perSiteScaling)

	reserved := config["reserved"].(bool)
	properties.Reserved = azure.Bool(reserved)

	if v, ok := config["maximum_number_of_workers"]; ok {
		maximumNumberOfWorkers := int32(v.(int))
		properties.MaximumNumberOfWorkers = &maximumNumberOfWorkers
	}

	return &properties
}

func flattenAppServiceProperties(props *web.AppServicePlanProperties) map[string]interface{} {
	properties := make(map[string]interface{}, 0)

	if props.MaximumNumberOfWorkers != nil {
		properties["maximum_number_of_workers"] = int(*props.MaximumNumberOfWorkers)
	}

	if props.PerSiteScaling != nil {
		properties["per_site_scaling"] = *props.PerSiteScaling
	}

	if props.Reserved != nil {
		properties["reserved"] = *props.Reserved
	}

	return properties
}
