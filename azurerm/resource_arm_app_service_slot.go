package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppServiceSlot() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceSlotCreate,
		Read:   resourceArmAppServiceSlotRead,
		Update: resourceArmAppServiceSlotCreate,
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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"identity": azure.SchemaAppServiceIdentity(),

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

			"auth_settings": azure.SchemaAppServiceAuthSettings(),

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
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},
				},
			},

			"tags": tags.Schema(),

			"site_credential": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"default_site_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmAppServiceSlotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).web.AppServicesClient
	ctx := meta.(*ArmClient).StopContext

	slot := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	appServiceName := d.Get("app_service_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.GetSlot(ctx, resourceGroup, appServiceName, slot)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_app_service_slot", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	appServicePlanId := d.Get("app_service_plan_id").(string)
	enabled := d.Get("enabled").(bool)
	httpsOnly := d.Get("https_only").(bool)
	t := d.Get("tags").(map[string]interface{})
	affinity := d.Get("client_affinity_enabled").(bool)

	siteConfig, err := azure.ExpandAppServiceSiteConfig(d.Get("site_config"))
	if err != nil {
		return fmt.Errorf("Error expanding `site_config` for App Service Slot %q (Resource Group %q): %s", slot, resourceGroup, err)
	}
	siteEnvelope := web.Site{
		Location: &location,
		Tags:     tags.Expand(t),
		SiteProperties: &web.SiteProperties{
			ServerFarmID:          utils.String(appServicePlanId),
			Enabled:               utils.Bool(enabled),
			HTTPSOnly:             utils.Bool(httpsOnly),
			SiteConfig:            siteConfig,
			ClientAffinityEnabled: &affinity,
		},
	}

	if _, ok := d.GetOk("identity"); ok {
		appServiceIdentity := azure.ExpandAppServiceIdentity(d)
		siteEnvelope.Identity = appServiceIdentity
	}

	createFuture, err := client.CreateOrUpdateSlot(ctx, resourceGroup, appServiceName, siteEnvelope, slot)
	if err != nil {
		return fmt.Errorf("Error creating Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resourceGroup, err)
	}

	err = createFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for creation of Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resourceGroup, err)
	}

	read, err := client.GetSlot(ctx, resourceGroup, appServiceName, slot)
	if err != nil {
		return fmt.Errorf("Error retrieving Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for Slot %q (App Service %q / Resource Group %q) ID", slot, appServiceName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppServiceSlotUpdate(d, meta)
}

func resourceArmAppServiceSlotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).web.AppServicesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	appServiceName := id.Path["sites"]
	slot := id.Path["slots"]

	location := azure.NormalizeLocation(d.Get("location").(string))
	appServicePlanId := d.Get("app_service_plan_id").(string)
	siteConfig, err := azure.ExpandAppServiceSiteConfig(d.Get("site_config"))
	if err != nil {
		return fmt.Errorf("Error expanding `site_config` for App Service Slot %q (Resource Group %q): %s", slot, resourceGroup, err)
	}
	enabled := d.Get("enabled").(bool)
	httpsOnly := d.Get("https_only").(bool)
	t := d.Get("tags").(map[string]interface{})

	siteEnvelope := web.Site{
		Location: &location,
		Tags:     tags.Expand(t),
		SiteProperties: &web.SiteProperties{
			ServerFarmID: utils.String(appServicePlanId),
			Enabled:      utils.Bool(enabled),
			HTTPSOnly:    utils.Bool(httpsOnly),
			SiteConfig:   siteConfig,
		},
	}
	if v, ok := d.GetOk("client_affinity_enabled"); ok {
		enabled := v.(bool)
		siteEnvelope.SiteProperties.ClientAffinityEnabled = utils.Bool(enabled)
	}
	createFuture, err := client.CreateOrUpdateSlot(ctx, resourceGroup, appServiceName, siteEnvelope, slot)
	if err != nil {
		return fmt.Errorf("Error updating Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resourceGroup, err)
	}

	err = createFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for update of Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resourceGroup, err)
	}

	if d.HasChange("site_config") {
		// update the main configuration
		siteConfig, err := azure.ExpandAppServiceSiteConfig(d.Get("site_config"))
		if err != nil {
			return fmt.Errorf("Error expanding `site_config` for App Service Slot %q (Resource Group %q): %s", slot, resourceGroup, err)
		}
		siteConfigResource := web.SiteConfigResource{
			SiteConfig: siteConfig,
		}
		if _, err := client.CreateOrUpdateConfigurationSlot(ctx, resourceGroup, appServiceName, siteConfigResource, slot); err != nil {
			return fmt.Errorf("Error updating Configuration for App Service Slot %q/%q: %+v", appServiceName, slot, err)
		}
	}

	if d.HasChange("auth_settings") {
		authSettingsRaw := d.Get("auth_settings").([]interface{})
		authSettingsProperties := azure.ExpandAppServiceAuthSettings(authSettingsRaw)
		id := d.Id()
		authSettings := web.SiteAuthSettings{
			ID:                         &id,
			SiteAuthSettingsProperties: &authSettingsProperties,
		}

		if _, err := client.UpdateAuthSettingsSlot(ctx, resourceGroup, appServiceName, authSettings, slot); err != nil {
			return fmt.Errorf("Error updating Authentication Settings for App Service %q: %+v", appServiceName, err)
		}
	}

	if d.HasChange("app_settings") {
		// update the AppSettings
		appSettings := expandAppServiceAppSettings(d)
		settings := web.StringDictionary{
			Properties: appSettings,
		}

		if _, err := client.UpdateApplicationSettingsSlot(ctx, resourceGroup, appServiceName, settings, slot); err != nil {
			return fmt.Errorf("Error updating Application Settings for App Service Slot %q/%q: %+v", appServiceName, slot, err)
		}
	}

	if d.HasChange("connection_string") {
		// update the ConnectionStrings
		connectionStrings := expandAppServiceConnectionStrings(d)
		properties := web.ConnectionStringDictionary{
			Properties: connectionStrings,
		}

		if _, err := client.UpdateConnectionStringsSlot(ctx, resourceGroup, appServiceName, properties, slot); err != nil {
			return fmt.Errorf("Error updating Connection Strings for App Service Slot %q/%q: %+v", appServiceName, slot, err)
		}
	}

	if d.HasChange("identity") {
		identity := azure.ExpandAppServiceIdentity(d)
		sitePatchResource := web.SitePatchResource{
			ID:       utils.String(d.Id()),
			Identity: identity,
		}
		_, err := client.UpdateSlot(ctx, resourceGroup, appServiceName, sitePatchResource, slot)
		if err != nil {
			return fmt.Errorf("Error updating Managed Service Identity for App Service Slot %q/%q: %+v", appServiceName, slot, err)
		}
	}

	return resourceArmAppServiceSlotRead(d, meta)
}

func resourceArmAppServiceSlotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).web.AppServicesClient

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	appServiceName := id.Path["sites"]
	slot := id.Path["slots"]

	ctx := meta.(*ArmClient).StopContext
	resp, err := client.GetSlot(ctx, resourceGroup, appServiceName, slot)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Slot %q (App Service %q / Resource Group %q) were not found - removing from state!", slot, appServiceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resourceGroup, err)
	}

	configResp, err := client.GetConfigurationSlot(ctx, resourceGroup, appServiceName, slot)
	if err != nil {
		if utils.ResponseWasNotFound(configResp.Response) {
			log.Printf("[DEBUG] Configuration for Slot %q (App Service %q / Resource Group %q) were not found - removing from state!", slot, appServiceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Configuration for Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resourceGroup, err)
	}

	authResp, err := client.GetAuthSettingsSlot(ctx, resourceGroup, appServiceName, slot)
	if err != nil {
		return fmt.Errorf("Error reading Auth Settings for Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resourceGroup, err)
	}

	appSettingsResp, err := client.ListApplicationSettingsSlot(ctx, resourceGroup, appServiceName, slot)
	if err != nil {
		if utils.ResponseWasNotFound(appSettingsResp.Response) {
			log.Printf("[DEBUG] App Settings for Slot %q (App Service %q / Resource Group %q) were not found - removing from state!", slot, appServiceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading App Settings for Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resourceGroup, err)
	}

	connectionStringsResp, err := client.ListConnectionStringsSlot(ctx, resourceGroup, appServiceName, slot)
	if err != nil {
		return fmt.Errorf("Error listing Connection Strings for Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resourceGroup, err)
	}

	siteCredFuture, err := client.ListPublishingCredentialsSlot(ctx, resourceGroup, appServiceName, slot)
	if err != nil {
		return fmt.Errorf("Error retrieving publishing credentials for Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resourceGroup, err)
	}
	err = siteCredFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for publishing credentials for Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resourceGroup, err)
	}
	siteCredResp, err := siteCredFuture.Result(*client)
	if err != nil {
		return fmt.Errorf("Error reading publishing credentials for Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resourceGroup, err)
	}

	d.Set("name", slot)
	d.Set("app_service_name", appServiceName)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.SiteProperties; props != nil {
		d.Set("app_service_plan_id", props.ServerFarmID)
		d.Set("client_affinity_enabled", props.ClientAffinityEnabled)
		d.Set("default_site_hostname", props.DefaultHostName)
		d.Set("enabled", props.Enabled)
		d.Set("https_only", props.HTTPSOnly)
	}

	if err := d.Set("app_settings", flattenAppServiceAppSettings(appSettingsResp.Properties)); err != nil {
		return fmt.Errorf("Error setting `app_settings`: %s", err)
	}
	if err := d.Set("connection_string", flattenAppServiceConnectionStrings(connectionStringsResp.Properties)); err != nil {
		return fmt.Errorf("Error setting `connection_string`: %s", err)
	}

	authSettings := azure.FlattenAppServiceAuthSettings(authResp.SiteAuthSettingsProperties)
	if err := d.Set("auth_settings", authSettings); err != nil {
		return fmt.Errorf("Error setting `auth_settings`: %s", err)
	}

	identity := azure.FlattenAppServiceIdentity(resp.Identity)
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("Error setting `identity`: %s", err)
	}

	siteCred := flattenAppServiceSiteCredential(siteCredResp.UserProperties)
	if err := d.Set("site_credential", siteCred); err != nil {
		return fmt.Errorf("Error setting `site_credential`: %s", err)
	}

	siteConfig := azure.FlattenAppServiceSiteConfig(configResp.SiteConfig)
	if err := d.Set("site_config", siteConfig); err != nil {
		return fmt.Errorf("Error setting `site_config`: %s", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmAppServiceSlotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).web.AppServicesClient

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	appServiceName := id.Path["sites"]
	slot := id.Path["slots"]

	log.Printf("[DEBUG] Deleting Slot %q (App Service %q / Resource Group %q)", slot, appServiceName, resourceGroup)

	deleteMetrics := true
	deleteEmptyServerFarm := false
	ctx := meta.(*ArmClient).StopContext
	resp, err := client.DeleteSlot(ctx, resourceGroup, appServiceName, slot, &deleteMetrics, &deleteEmptyServerFarm)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting Slot %q (App Service %q / Resource Group %q): %s", slot, appServiceName, resourceGroup, err)
		}
	}

	return nil
}
