package azurerm

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApplicationInsights() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApplicationInsightsCreateOrUpdate,
		Read:   resourceArmApplicationInsightsRead,
		Update: resourceArmApplicationInsightsCreateOrUpdate,
		Delete: resourceArmApplicationInsightsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"application_type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(insights.Web),
					string(insights.Other),
				}, true),
			},

			"tags": tagsSchema(),

			"app_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instrumentation_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceArmApplicationInsightsCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appInsightsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Application Insights creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	applicationType := d.Get("application_type").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})

	applicationInsightsComponentProperties := insights.ApplicationInsightsComponentProperties{
		ApplicationID:   &name,
		ApplicationType: insights.ApplicationType(applicationType),
	}

	insightProperties := insights.ApplicationInsightsComponent{
		Name:     &name,
		Location: &location,
		Kind:     &applicationType,
		ApplicationInsightsComponentProperties: &applicationInsightsComponentProperties,
		Tags: expandTags(tags),
	}

	_, err := client.CreateOrUpdate(ctx, resGroup, name, insightProperties)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read AzureRM Application Insights '%s' (Resource Group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmApplicationInsightsRead(d, meta)
}

func resourceArmApplicationInsightsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appInsightsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading AzureRM Application Insights '%s'", id)

	resGroup := id.ResourceGroup
	name := id.Path["components"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Application Insights '%s': %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.ApplicationInsightsComponentProperties; props != nil {
		d.Set("application_type", string(props.ApplicationType))
		d.Set("app_id", props.AppID)
		d.Set("instrumentation_key", props.InstrumentationKey)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmApplicationInsightsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appInsightsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["components"]

	log.Printf("[DEBUG] Deleting AzureRM Application Insights '%s' (resource group '%s')", name, resGroup)

	resp, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Error issuing AzureRM delete request for Application Insights '%s': %+v", name, err)
	}

	return err
}
