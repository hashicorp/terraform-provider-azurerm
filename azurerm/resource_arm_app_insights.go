package azurerm

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/appinsights"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceArmApplicationInsights() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApplicationInsightsCreateOrUpdate,
		Read:   resourceArmApplicationInsightsRead,
		Update: resourceArmApplicationInsightsCreateOrUpdate,
		Delete: resourceArmApplicationInsightsDelete,

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
			"application_type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(appinsights.Web),
					string(appinsights.Other),
				}, true),
			},
			"location": locationSchema(),
			"tags":     tagsSchema(),
			"instrumentation_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmApplicationInsightsCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	AppInsightsClient := client.appInsightsClient

	log.Printf("[INFO] preparing arguments for AzureRM Application Insights creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	applicationType := d.Get("application_type").(string)
	location := d.Get("location").(string)
	tags := d.Get("tags").(map[string]interface{})

	applicationInsightsComponentProperties := appinsights.ApplicationInsightsComponentProperties{
		ApplicationID:   &name,
		ApplicationType: appinsights.ApplicationType(applicationType),
	}

	insightProperties := appinsights.ApplicationInsightsComponent{
		Name:     &name,
		Location: &location,
		Tags:     expandTags(tags),
		Kind:     &applicationType,
		ApplicationInsightsComponentProperties: &applicationInsightsComponentProperties,
	}

	_, err := AppInsightsClient.CreateOrUpdate(resGroup, name, insightProperties)
	if err != nil {
		return err
	}

	read, err := AppInsightsClient.Get(resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read application insights %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmApplicationInsightsRead(d, meta)
}

func resourceArmApplicationInsightsRead(d *schema.ResourceData, meta interface{}) error {
	AppInsightsClient := meta.(*ArmClient).appInsightsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading application insights %s", id)

	resGroup := id.ResourceGroup
	name := id.Path["components"]

	resp, err := AppInsightsClient.Get(resGroup, name)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM application insights %s: %s", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	if resp.ApplicationInsightsComponentProperties != nil {
		d.Set("application_id", resp.ApplicationInsightsComponentProperties.ApplicationID)
		d.Set("application_type", string(resp.ApplicationInsightsComponentProperties.ApplicationType))
		d.Set("instrumentation_key", resp.ApplicationInsightsComponentProperties.InstrumentationKey)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmApplicationInsightsDelete(d *schema.ResourceData, meta interface{}) error {
	AppInsightsClient := meta.(*ArmClient).appInsightsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["components"]

	log.Printf("[DEBUG] Deleting application insights %s: %s", resGroup, name)

	resp, err := AppInsightsClient.Delete(resGroup, name)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Error issuing AzureRM delete request for Application Insights '%s': %+v", name, err)
	}

	return err
}
