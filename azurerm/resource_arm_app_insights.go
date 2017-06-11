package azurerm

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/appinsights"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmAppInsights() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppInsightsCreateOrUpdate,
		Read:   resourceArmAppInsightsRead,
		Update: resourceArmAppInsightsCreateOrUpdate,
		Delete: resourceArmAppInsightsDelete,

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
			"kind": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

func resourceArmAppInsightsCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	AppInsightsClient := client.appInsightsClient

	log.Printf("[INFO] preparing arguments for Azure ARM App Insights creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	insightsType := "microsoft.insights/components"
	kind := d.Get("kind").(string)
	location := d.Get("location").(string)
	tags := d.Get("tags").(map[string]interface{})
	requestSource := "IbizaAIExtension"

	applicationInsightsComponentProperties := appinsights.ApplicationInsightsComponentProperties{
		ApplicationID:   &name,
		ApplicationType: appinsights.ApplicationType(insightsType),
		RequestSource:   appinsights.RequestSource(requestSource),
	}

	insightProperties := appinsights.ApplicationInsightsComponent{
		Name:     &name,
		Type:     &insightsType,
		Location: &location,
		Tags:     expandTags(tags),
		Kind:     &kind,
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
		return fmt.Errorf("Cannot read app insights %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppInsightsRead(d, meta)
}

func resourceArmAppInsightsRead(d *schema.ResourceData, meta interface{}) error {
	AppInsightsClient := meta.(*ArmClient).appInsightsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading app insights %s", id)

	resGroup := id.ResourceGroup
	name := id.Path["components"]

	resp, err := AppInsightsClient.Get(resGroup, name)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure app insights %s: %s", name, err)
	}

	instrumentationKey := resp.ApplicationInsightsComponentProperties.InstrumentationKey

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("instrumentation_key", instrumentationKey)

	return nil
}

func resourceArmAppInsightsDelete(d *schema.ResourceData, meta interface{}) error {
	AppInsightsClient := meta.(*ArmClient).appInsightsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["components"]

	log.Printf("[DEBUG] Deleting app insights %s: %s", resGroup, name)

	_, err = AppInsightsClient.Delete(resGroup, name)

	return err
}
