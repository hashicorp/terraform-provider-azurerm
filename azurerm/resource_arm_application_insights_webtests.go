package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"

	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApplicationInsightsWebTests() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApplicationInsightsWebTestsCreateUpdate,
		Read:   resourceArmApplicationInsightsWebTestsRead,
		Update: resourceArmApplicationInsightsWebTestsCreateUpdate,
		Delete: resourceArmApplicationInsightsWebTestsDelete,
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

			"application_insights_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"location": locationSchema(),

			"kind": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					"ping",
					"multistep",
				}, true),
			},

			"frequency": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
				ValidateFunc: validation.IntInSlice([]int{
					300,
					600,
					900,
				}),
			},

			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"retry_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"geo_locations": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateFunc:     validate.NoEmptyStrings,
					StateFunc:        azureRMNormalizeLocation,
					DiffSuppressFunc: azureRMSuppressLocationDiff,
				},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"test_configuration": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.XmlDiff,
			},

			"tags": tagsSchema(),

			"synthetic_monitor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"provisioning_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmApplicationInsightsWebTestsCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appInsightsWebTestsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Application Insights WebTest creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	appInsightsID := d.Get("application_insights_id").(string)

	id, err := parseAzureResourceID(appInsightsID)
	if err != nil {
		return err
	}

	appInsightsName := id.Path["components"]

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Application Insights WebTests %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_application_insights_webTests", *existing.ID)
		}
	}

	location := azureRMNormalizeLocation(d.Get("location").(string))
	kind := d.Get("kind").(string)
	description := d.Get("description").(string)
	frequency := int32(d.Get("frequency").(int))
	timeout := int32(d.Get("timeout").(int))
	isEnabled := d.Get("enabled").(bool)
	retryEnabled := d.Get("retry_enabled").(bool)
	geoLocations := extractGeoLocations(d)
	testConf := d.Get("test_configuration").(string)

	tags := d.Get("tags").(map[string]interface{})
	tagKey := fmt.Sprintf("hidden-link:/subscriptions/%s/resourceGroups/%s/providers/microsoft.insights/components/%s", client.SubscriptionID, resGroup, appInsightsName)
	tags[tagKey] = "Resource"

	testConfiguration := insights.WebTestPropertiesConfiguration{
		WebTest: &testConf,
	}

	webTestProperties := insights.WebTestProperties{
		SyntheticMonitorID: &name,
		WebTestName:        &name,
		Description:        &description,
		Enabled:            &isEnabled,
		Frequency:          &frequency,
		Timeout:            &timeout,
		WebTestKind:        insights.WebTestKind(kind),
		RetryEnabled:       &retryEnabled,
		Locations:          &geoLocations,
		Configuration:      &testConfiguration,
	}

	webTest := insights.WebTest{
		Name:              &name,
		Location:          &location,
		Kind:              insights.WebTestKind(kind),
		WebTestProperties: &webTestProperties,
		Tags:              expandTags(tags),
	}

	resp, err := client.CreateOrUpdate(ctx, resGroup, name, webTest)

	if err != nil {
		return fmt.Errorf("Error creating Application Insights WebTest %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmApplicationInsightsWebTestsRead(d, meta)
}

func resourceArmApplicationInsightsWebTestsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appInsightsWebTestsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading AzureRM Application Insights WebTests '%s'", id)

	resGroup := id.ResourceGroup
	name := id.Path["webtests"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Application Insights WebTests '%s': %+v", name, err)
	}

	log.Printf("[DEBUG] AzureRM Application Insights WebTests name '%s'", name)

	appInsightsId := ""
	tags := resp.Tags
	for i := range tags {
		if strings.HasPrefix(i, "hidden-link") {
			appInsightsId = strings.Split(i, ":")[1]
		}
	}
	d.Set("application_insights_id", appInsightsId)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if kind := resp.Kind; kind != "" {
		d.Set("kind", resp.Kind)
	}

	if props := resp.WebTestProperties; props != nil {
		d.Set("synthetic_monitor_id", props.SyntheticMonitorID)
		d.Set("description", props.Description)
		d.Set("enabled", props.Enabled)
		d.Set("frequency", props.Frequency)
		d.Set("timeout", props.Timeout)
		d.Set("retry_enabled", props.RetryEnabled)
		d.Set("configuration", props.Configuration)
		d.Set("provisioning_state", props.ProvisioningState)
		d.Set("test_configuration", props.Configuration.WebTest)

		if err := d.Set("geo_locations", flattenGeoLocations(props.Locations)); err != nil {
			return fmt.Errorf("Error setting `geo_locations`: %+v", err)
		}
	}

	log.Printf("[DEBUG] AzureRM Application Insights WebTests synetheticmonitorid '%s'", d.Get("synthetic_monitor_id"))

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmApplicationInsightsWebTestsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appInsightsWebTestsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["webtests"]

	log.Printf("[DEBUG] Deleting AzureRM Application Insights WebTest '%s' (resource group '%s')", name, resGroup)

	resp, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Error issuing AzureRM delete request for Application Insights WebTest '%s': %+v", name, err)
	}

	return err
}

func extractGeoLocations(d *schema.ResourceData) []insights.WebTestGeolocation {
	geoLocations := d.Get("geo_locations").([]interface{})
	locations := make([]insights.WebTestGeolocation, 0)

	for _, location := range geoLocations {
		lc := location.(string)
		loc := insights.WebTestGeolocation{
			Location: &lc,
		}
		locations = append(locations, loc)
	}

	return locations
}

func flattenGeoLocations(input *[]insights.WebTestGeolocation) []string {
	results := make([]string, 0)
	if input == nil {
		return results
	}

	for _, prop := range *input {
		if prop.Location != nil {
			results = append(results, azureRMNormalizeLocation(*prop.Location))
		}

	}

	return results
}
