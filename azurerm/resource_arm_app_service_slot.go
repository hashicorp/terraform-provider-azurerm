package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
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

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								"SystemAssigned",
							}, true),
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

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

			"site_config": azure.SchemaAppServiceSiteConfig(),

			"client_affinity_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"https_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"app_settings": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},

			"connection_string": {
				Type:     schema.TypeSet,
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

			"tags": tagsSchema(),

			"default_site_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmAppServiceSlotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM App Service Slot creation.")

	slot := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	appServiceName := d.Get("app_service_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.GetSlot(ctx, resGroup, appServiceName, slot)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_app_service_slot", *existing.ID)
		}
	}

	location := azureRMNormalizeLocation(d.Get("location").(string))
	appServicePlanId := d.Get("app_service_plan_id").(string)
	enabled := d.Get("enabled").(bool)
	httpsOnly := d.Get("https_only").(bool)
	tags := d.Get("tags").(map[string]interface{})

	siteConfig := azure.ExpandAppServiceSiteConfig(d.Get("site_config"))
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

	if _, ok := d.GetOk("identity"); ok {
		appServiceIdentity := expandAzureRmAppServiceIdentity(d)
		siteEnvelope.Identity = appServiceIdentity
	}

	if v, ok := d.GetOk("client_affinity_enabled"); ok {
		enabled := v.(bool)
		siteEnvelope.SiteProperties.ClientAffinityEnabled = utils.Bool(enabled)
	}

	resp, err := client.Get(ctx, resGroup, appServiceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("[DEBUG] App Service %q (resource group %q) was not found.", appServiceName, resGroup)
		}
		return fmt.Errorf("Error making Read request on AzureRM App Service %q: %+v", appServiceName, err)
	}

	createFuture, err := client.CreateOrUpdateSlot(ctx, resGroup, appServiceName, siteEnvelope, slot)
	if err != nil {
		return err
	}

	err = createFuture.WaitForCompletionRef(ctx, client.Client)
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
	location := azureRMNormalizeLocation(d.Get("location").(string))
	appServicePlanId := d.Get("app_service_plan_id").(string)
	slot := id.Path["slots"]
	siteConfig := azure.ExpandAppServiceSiteConfig(d.Get("site_config"))
	enabled := d.Get("enabled").(bool)
	httpsOnly := d.Get("https_only").(bool)
	tags := d.Get("tags").(map[string]interface{})
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
	createFuture, err := client.CreateOrUpdateSlot(ctx, resGroup, appServiceName, siteEnvelope, slot)
	if err != nil {
		return err
	}

	err = createFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return err
	}
	if d.HasChange("site_config") {
		// update the main configuration
		siteConfig := azure.ExpandAppServiceSiteConfig(d.Get("site_config"))
		siteConfigResource := web.SiteConfigResource{
			SiteConfig: &siteConfig,
		}
		if _, err := client.CreateOrUpdateConfigurationSlot(ctx, resGroup, appServiceName, siteConfigResource, slot); err != nil {
			return fmt.Errorf("Error updating Configuration for App Service Slot %q/%q: %+v", appServiceName, slot, err)
		}
	}

	if d.HasChange("client_affinity_enabled") {
		affinity := d.Get("client_affinity_enabled").(bool)
		sitePatchResource := web.SitePatchResource{
			ID: utils.String(d.Id()),
			SitePatchResourceProperties: &web.SitePatchResourceProperties{
				ClientAffinityEnabled: &affinity,
			},
		}
		_, err := client.UpdateSlot(ctx, resGroup, appServiceName, sitePatchResource, slot)
		if err != nil {
			return fmt.Errorf("Error updating App Service ARR Affinity setting %q: %+v", slot, err)
		}
	}

	if d.HasChange("app_settings") {
		// update the AppSettings
		appSettings := expandAppServiceAppSettings(d)
		settings := web.StringDictionary{
			Properties: appSettings,
		}

		if _, err := client.UpdateApplicationSettingsSlot(ctx, resGroup, appServiceName, settings, slot); err != nil {
			return fmt.Errorf("Error updating Application Settings for App Service Slot %q/%q: %+v", appServiceName, slot, err)
		}
	}

	if d.HasChange("connection_string") {
		// update the ConnectionStrings
		connectionStrings := expandAppServiceConnectionStrings(d)
		properties := web.ConnectionStringDictionary{
			Properties: connectionStrings,
		}

		if _, err := client.UpdateConnectionStringsSlot(ctx, resGroup, appServiceName, properties, slot); err != nil {
			return fmt.Errorf("Error updating Connection Strings for App Service Slot %q/%q: %+v", appServiceName, slot, err)
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
		if utils.ResponseWasNotFound(configResp.Response) {
			log.Printf("[DEBUG] Configuration of App Service Slot %q/%q (resource group %q) was not found", appServiceName, slot, resGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM App Service Slot Configuration %q/%q: %+v", appServiceName, slot, err)
	}

	appSettingsResp, err := client.ListApplicationSettingsSlot(ctx, resGroup, appServiceName, slot)
	if err != nil {
		if utils.ResponseWasNotFound(appSettingsResp.Response) {
			log.Printf("[DEBUG] Application Settings of App Service Slot %q/%q (resource group %q) were not found", appServiceName, slot, resGroup)
			d.SetId("")
			return nil
		}
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

	siteConfig := azure.FlattenAppServiceSiteConfig(configResp.SiteConfig)
	if err := d.Set("site_config", siteConfig); err != nil {
		return err
	}

	identity := flattenAzureRmAppServiceMachineIdentity(resp.Identity)
	if err := d.Set("identity", identity); err != nil {
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
	ctx := meta.(*ArmClient).StopContext
	resp, err := client.DeleteSlot(ctx, resGroup, appServiceName, slot, &deleteMetrics, &deleteEmptyServerFarm)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return err
		}
	}

	return nil
}
