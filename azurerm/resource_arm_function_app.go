package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/arm/web"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// Azure Function App shares the same infrastructure with Azure App Service.
// So this resource will reuse most of the App Service code, but remove the configurations which are not applicable for Function App.
func resourceArmFunctionApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmFunctionAppCreate,
		Read:   resourceArmFunctionAppRead,
		Update: resourceArmFunctionAppUpdate,
		Delete: resourceArmFunctionAppDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAppServiceName,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"app_service_plan_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,

				// TODO: (tombuildsstuff) support Update once the API is fixed:
				// https://github.com/Azure/azure-rest-api-specs/issues/1697
				ForceNew: true,
			},

			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "~1",
				ValidateFunc: validation.StringInSlice([]string{
					"~1",
					"beta",
				}, false),
			},

			"storage_connection_string": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			"app_settings": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			// TODO: (tombuildsstuff) support Update once the API is fixed:
			// https://github.com/Azure/azure-rest-api-specs/issues/1697
			"tags": tagsForceNewSchema(),

			"default_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmFunctionAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicesClient

	log.Printf("[INFO] preparing arguments for AzureRM Function App creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := d.Get("location").(string)
	kind := "functionapp"
	appServicePlanID := d.Get("app_service_plan_id").(string)
	enabled := d.Get("enabled").(bool)
	tags := d.Get("tags").(map[string]interface{})
	basicAppSettings := getBasicFunctionAppAppSettings(d)

	siteEnvelope := web.Site{
		Kind:     &kind,
		Location: &location,
		Tags:     expandTags(tags),
		SiteProperties: &web.SiteProperties{
			ServerFarmID: utils.String(appServicePlanID),
			Enabled:      utils.Bool(enabled),
			SiteConfig: &web.SiteConfig{
				AppSettings: &basicAppSettings,
			},
		},
	}

	skipDNSRegistration := false
	forceDNSRegistration := false
	skipCustomDomainVerification := true
	ttlInSeconds := "60"
	_, createErr := client.CreateOrUpdate(resGroup, name, siteEnvelope, &skipDNSRegistration, &skipCustomDomainVerification, &forceDNSRegistration, ttlInSeconds, make(chan struct{}))
	err := <-createErr
	if err != nil {
		return err
	}

	read, err := client.Get(resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Function App %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmFunctionAppUpdate(d, meta)
}

func resourceArmFunctionAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["sites"]

	if d.HasChange("app_settings") || d.HasChange("version") {
		appSettings := expandFunctionAppAppSettings(d)
		settings := web.StringDictionary{
			Properties: appSettings,
		}

		_, err := client.UpdateApplicationSettings(resGroup, name, settings)
		if err != nil {
			return fmt.Errorf("Error updating Application Settings for Function App %q: %+v", name, err)
		}
	}

	return resourceArmFunctionAppRead(d, meta)
}

func resourceArmFunctionAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["sites"]

	resp, err := client.Get(resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Function App %q (resource group %q) was not found - removing from state", name, resGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Function App %q: %+v", name, err)
	}

	appSettingsResp, err := client.ListApplicationSettings(resGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM Function App AppSettings %q: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))

	if props := resp.SiteProperties; props != nil {
		d.Set("app_service_plan_id", props.ServerFarmID)
		d.Set("enabled", props.Enabled)
		d.Set("default_hostname", props.DefaultHostName)
	}

	appSettings := flattenAppServiceAppSettings(appSettingsResp.Properties)

	d.Set("storage_connection_string", appSettings["AzureWebJobsStorage"])
	d.Set("version", appSettings["FUNCTIONS_EXTENSION_VERSION"])

	delete(appSettings, "AzureWebJobsDashboard")
	delete(appSettings, "AzureWebJobsStorage")
	delete(appSettings, "FUNCTIONS_EXTENSION_VERSION")
	delete(appSettings, "WEBSITE_CONTENTSHARE")
	delete(appSettings, "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING")

	if err := d.Set("app_settings", appSettings); err != nil {
		return err
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmFunctionAppDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["sites"]

	log.Printf("[DEBUG] Deleting Function App %q (resource group %q)", name, resGroup)

	deleteMetrics := true
	deleteEmptyServerFarm := false
	skipDNSRegistration := true
	resp, err := client.Delete(resGroup, name, &deleteMetrics, &deleteEmptyServerFarm, &skipDNSRegistration)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return err
		}
	}

	return nil
}

func getBasicFunctionAppAppSettings(d *schema.ResourceData) []web.NameValuePair {
	dashboardPropName := "AzureWebJobsDashboard"
	storagePropName := "AzureWebJobsStorage"
	functionVersionPropName := "FUNCTIONS_EXTENSION_VERSION"
	contentSharePropName := "WEBSITE_CONTENTSHARE"
	contentFileConnStringPropName := "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"

	storageConnection := d.Get("storage_connection_string").(string)
	functionVersion := d.Get("version").(string)
	contentShare := d.Get("name").(string) + "-content"

	return []web.NameValuePair{
		{Name: &dashboardPropName, Value: &storageConnection},
		{Name: &storagePropName, Value: &storageConnection},
		{Name: &functionVersionPropName, Value: &functionVersion},
		{Name: &contentSharePropName, Value: &contentShare},
		{Name: &contentFileConnStringPropName, Value: &storageConnection},
	}
}

func expandFunctionAppAppSettings(d *schema.ResourceData) *map[string]*string {
	output := expandAppServiceAppSettings(d)

	basicAppSettings := getBasicFunctionAppAppSettings(d)
	for _, p := range basicAppSettings {
		(*output)[*p.Name] = p.Value
	}

	return output
}
