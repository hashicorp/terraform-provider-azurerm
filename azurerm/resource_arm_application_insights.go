package azurerm

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 30),
			Update: schema.DefaultTimeout(time.Minute * 30),
			Delete: schema.DefaultTimeout(time.Minute * 30),
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
					"web",
					"other",
					"java",
					"MobileCenter",
					"phone",
					"store",
					"ios",
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

	if d.IsNewResource() {
		// first check if there's one in this subscription requiring import
		resp, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for the existence of Application Insights %q (Resource Group %q): %+v", name, resGroup, err)
			}
		}

		if resp.ID != nil {
			return tf.ImportAsExistsError("azurerm_application_insights", *resp.ID)
		}
	}

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

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(tf.TimeoutForCreateUpdate(d)))
	defer cancel()
	resp, err := client.CreateOrUpdate(waitCtx, resGroup, name, insightProperties)
	if err != nil {
		// @tombuildsstuff - from 2018-08-14 the Create call started returning a 201 instead of 200
		// which doesn't match the Swagger - this works around it until that's fixed
		// BUG: https://github.com/Azure/azure-sdk-for-go/issues/2465
		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("Error creating Application Insights %q (Resource Group %q): %+v", name, resGroup, err)
		}
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Application Insights %q (Resource Group %q): %+v", name, resGroup, err)
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

	log.Printf("[DEBUG] Reading AzureRM Application Insights %q", id)

	resGroup := id.ResourceGroup
	name := id.Path["components"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Application Insights %q: %+v", name, err)
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

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutDelete))
	defer cancel()
	resp, err := client.Delete(waitCtx, resGroup, name)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Error issuing AzureRM delete request for Application Insights '%s': %+v", name, err)
	}

	return err
}
