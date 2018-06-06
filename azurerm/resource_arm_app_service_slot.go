package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2016-09-01/web"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppServiceSlot() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceSlotCreate,
		Read:   resourceArmAppServiceSlotRead,
		Update: resourceArmAppServiceSlotUpdate,
		Delete: resourceArmAppServiceSlotDelete,
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

			"app_service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"app_service_plan_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"site_config": azSchema.AppServiceSiteConfigSchema(),

			"client_affinity_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,

				// TODO: (tombuildsstuff) support Update once the API is fixed:
				// https://github.com/Azure/azure-rest-api-specs/issues/1697
				ForceNew: true,
			},

			"https_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,

				// TODO: (tombuildsstuff) support Update once the API is fixed:
				// https://github.com/Azure/azure-rest-api-specs/issues/1697
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

			"app_settings": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},

			"connection_string": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(web.APIHub),
								string(web.Custom),
								string(web.DocDb),
								string(web.EventHub),
								string(web.MySQL),
								string(web.NotificationHub),
								string(web.PostgreSQL),
								string(web.RedisCache),
								string(web.ServiceBus),
								string(web.SQLAzure),
								string(web.SQLServer),
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},
					},
				},
			},

			// TODO: (tombuildsstuff) support Update once the API is fixed:
			// https://github.com/Azure/azure-rest-api-specs/issues/1697
			"tags": tagsForceNewSchema(),

			"default_site_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmAppServiceSlotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicesClient

	log.Printf("[INFO] preparing arguments for AzureRM App Service Slot creation.")

	slot := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	appServiceName := d.Get("app_service_name").(string)
	appServicePlanId := d.Get("app_service_plan_id").(string)
	enabled := d.Get("enabled").(bool)
	httpsOnly := d.Get("https_only").(bool)
	tags := d.Get("tags").(map[string]interface{})

	siteConfig := azSchema.ExpandAppServiceSiteConfig(d.Get("site_config"))
	siteEnvelope := web.Site{
		Location: &location,
		Tags:     expandTags(tags),
		SiteProperties: &web.SiteProperties{
			ServerFarmID: utils.String(appServicePlanId),
			Enabled:      utils.Bool(enabled),
			HTTPSOnly:    utils.Bool(httpsOnly),
			SiteConfig:   &siteConfig,
		},
	}

	if v, ok := d.GetOk("client_affinity_enabled"); ok {
		enabled := v.(bool)
		siteEnvelope.SiteProperties.ClientAffinityEnabled = utils.Bool(enabled)
	}

	// NOTE: these seem like sensible defaults, in lieu of any better documentation.
	skipDNSRegistration := false
	forceDNSRegistration := false
	skipCustomDomainVerification := true
	ttlInSeconds := "60"
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Get(ctx, resGroup, appServiceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("[DEBUG] App Service %q (resource group %q) was not found.", appServiceName, resGroup)
		}
		return fmt.Errorf("Error making Read request on AzureRM App Service %q: %+v", appServiceName, err)
	}

	createFuture, err := client.CreateOrUpdateSlot(ctx, resGroup, appServiceName, siteEnvelope, slot, &skipDNSRegistration, &skipCustomDomainVerification, &forceDNSRegistration, ttlInSeconds)
	if err != nil {
		return err
	}

	err = createFuture.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}

	read, err := client.GetSlot(ctx, resGroup, appServiceName, slot)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read App Service Slot %q/%q (resource group %q) ID", appServiceName, slot, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppServiceSlotUpdate(d, meta)
}

func resourceArmAppServiceSlotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	appServiceName := id.Path["sites"]
	slot := id.Path["slots"]

	if d.HasChange("site_config") {
		// update the main configuration
		siteConfig := azSchema.ExpandAppServiceSiteConfig(d.Get("site_config"))
		siteConfigResource := web.SiteConfigResource{
			SiteConfig: &siteConfig,
		}
		_, err := client.CreateOrUpdateConfigurationSlot(ctx, resGroup, appServiceName, siteConfigResource, slot)
		if err != nil {
			return fmt.Errorf("Error updating Configuration for App Service Slot %q/%q: %+v", appServiceName, slot, err)
		}
	}

	if d.HasChange("app_settings") {
		// update the AppSettings
		appSettings := expandAppServiceAppSettings(d)
		settings := web.StringDictionary{
			Properties: appSettings,
		}

		_, err := client.UpdateApplicationSettingsSlot(ctx, resGroup, appServiceName, settings, slot)
		if err != nil {
			return fmt.Errorf("Error updating Application Settings for App Service Slot %q/%q: %+v", appServiceName, slot, err)
		}
	}

	if d.HasChange("connection_string") {
		// update the ConnectionStrings
		connectionStrings := expandAppServiceConnectionStrings(d)
		properties := web.ConnectionStringDictionary{
			Properties: connectionStrings,
		}

		_, err := client.UpdateConnectionStringsSlot(ctx, resGroup, appServiceName, properties, slot)
		if err != nil {
			return fmt.Errorf("Error updating Connection Strings for App Service %q/%q: %+v", appServiceName, slot, err)
		}
	}

	return resourceArmAppServiceSlotRead(d, meta)
}

func resourceArmAppServiceSlotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	appServiceName := id.Path["sites"]
	slot := id.Path["slots"]

	ctx := meta.(*ArmClient).StopContext
	resp, err := client.GetSlot(ctx, resGroup, appServiceName, slot)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Slot %q/%q (resource group %q) was not found - removing from state", appServiceName, slot, resGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM App Service Slot %q/%q: %+v", appServiceName, slot, err)
	}

	configResp, err := client.GetConfigurationSlot(ctx, resGroup, appServiceName, slot)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service Slot Configuration %q/%q: %+v", appServiceName, slot, err)
	}

	appSettingsResp, err := client.ListApplicationSettingsSlot(ctx, resGroup, appServiceName, slot)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service Slot AppSettings %q/%q: %+v", appServiceName, slot, err)
	}

	connectionStringsResp, err := client.ListConnectionStringsSlot(ctx, resGroup, appServiceName, slot)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service Slot ConnectionStrings %q/%q: %+v", appServiceName, slot, err)
	}

	d.Set("name", slot)
	d.Set("app_service_name", appServiceName)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.SiteProperties; props != nil {
		d.Set("app_service_plan_id", props.ServerFarmID)
		d.Set("client_affinity_enabled", props.ClientAffinityEnabled)
		d.Set("enabled", props.Enabled)
		d.Set("https_only", props.HTTPSOnly)
		d.Set("default_site_hostname", props.DefaultHostName)
	}

	if err := d.Set("app_settings", flattenAppServiceAppSettings(appSettingsResp.Properties)); err != nil {
		return err
	}
	if err := d.Set("connection_string", flattenAppServiceConnectionStrings(connectionStringsResp.Properties)); err != nil {
		return err
	}

	siteConfig := azSchema.FlattenAppServiceSiteConfig(configResp.SiteConfig)
	if err := d.Set("site_config", siteConfig); err != nil {
		return err
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmAppServiceSlotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	appServiceName := id.Path["sites"]
	slot := id.Path["slots"]

	log.Printf("[DEBUG] Deleting App Service Slot %q/%q (resource group %q)", appServiceName, slot, resGroup)

	deleteMetrics := true
	deleteEmptyServerFarm := false
	skipDNSRegistration := true
	ctx := meta.(*ArmClient).StopContext
	resp, err := client.DeleteSlot(ctx, resGroup, appServiceName, slot, &deleteMetrics, &deleteEmptyServerFarm, &skipDNSRegistration)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return err
		}
	}

	return nil
}
