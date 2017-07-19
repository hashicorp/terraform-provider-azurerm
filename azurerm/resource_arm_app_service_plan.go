package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"bytes"

	"github.com/Azure/azure-sdk-for-go/arm/web"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
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
			"maximum_number_of_workers": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceArmAppServicePlanCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	AppServicePlanClient := client.appServicePlansClient
	//AppServicePlanClient := meta.(*ArmClient).appServicePlansClient

	log.Printf("[INFO] preparing arguments for AzureRM App Service Plan creation.")

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := d.Get("location").(string)

	sku := expandAzureRmAppServicePlanSku(d)

	properties := web.AppServicePlanProperties{}
	if v, ok := d.GetOk("maximum_number_of_workers"); ok {
		maximumNumberOfWorkers := v.(int32)
		properties.MaximumNumberOfWorkers = &maximumNumberOfWorkers
	}

	appServicePlan := web.AppServicePlan{
		Location:                 &location,
		AppServicePlanProperties: &properties,
		Sku: &sku,
	}

	_, error := AppServicePlanClient.CreateOrUpdate(resGroup, name, appServicePlan, make(chan struct{}))
	err := <-error
	if err != nil {
		return err
	}

	read, err := AppServicePlanClient.Get(resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read AzureRM App Service Plan %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	log.Printf("[DEBUG] Waiting for App Service Plan (%s) to become available", name)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Accepted", "Updating"},
		Target:  []string{"Succeeded"},
		Refresh: appServicePlanStateRefreshFunc(client, resGroup, name),
		Timeout: 10 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for App Service Plan (%s) to become available: %s", name, err)
	}

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
		return fmt.Errorf("Error making Read request on Azure App Service Plan %s: %s", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))

	if props := resp.AppServicePlanProperties; props != nil {
		d.Set("maximum_number_of_workers", props.MaximumNumberOfWorkers)
	}

	sku := flattenAzureRmAppServicePlanSku(*resp.Sku)
	d.Set("sku", &sku)

	return nil
}

func resourceArmAppServicePlanDelete(d *schema.ResourceData, meta interface{}) error {
	AppServicePlanClient := meta.(*ArmClient).appServicePlansClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["serverfarms"]

	log.Printf("[DEBUG] Deleting app service plan %s: %s", resGroup, name)

	_, err = AppServicePlanClient.Delete(resGroup, name)

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

func flattenAzureRmAppServicePlanSku(profile web.SkuDescription) *schema.Set {
	skus := &schema.Set{
		F: resourceAzureRMAppServicePlanSkuHash,
	}

	sku := make(map[string]interface{}, 3)

	sku["tier"] = *profile.Tier
	sku["size"] = *profile.Size

	skus.Add(sku)

	return skus
}
