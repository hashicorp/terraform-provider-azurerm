package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2016-09-01/web"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceCreate,
		Read:   resourceArmAppServiceRead,
		Update: resourceArmAppServiceUpdate,
		Delete: resourceArmAppServiceDelete,
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

			// TODO: reusable schema
			"site_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"always_on": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"default_documents": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"dotnet_framework_version": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "v4.0",
							ValidateFunc: validation.StringInSlice([]string{
								"v2.0",
								"v4.0",
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"java_version": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"1.7",
								"1.8",
							}, false),
						},

						"java_container": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"JETTY",
								"TOMCAT",
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"java_container_version": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"local_mysql_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},

						"managed_pipeline_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(web.Classic),
								string(web.Integrated),
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"php_version": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"5.5",
								"5.6",
								"7.0",
								"7.1",
							}, false),
						},

						"python_version": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"2.7",
								"3.4",
							}, false),
						},

						"remote_debugging_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"remote_debugging_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"VS2012",
								"VS2013",
								"VS2015",
								"VS2017",
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"use_32_bit_worker_process": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},

						"websockets_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"scm_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(web.ScmTypeNone),
							ValidateFunc: validation.StringInSlice([]string{
								string(web.ScmTypeNone),
								string(web.ScmTypeLocalGit),
							}, false),
						},
					},
				},
			},

			"client_affinity_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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
							Type:     schema.TypeString,
							Required: true,
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

			"outbound_ip_addresses": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_control": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repo_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"branch": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceArmAppServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM App Service creation.")

	name := d.Get("name").(string)
	availabilityRequest := web.ResourceNameAvailabilityRequest{
		Name: utils.String(name),
		Type: web.CheckNameResourceTypesMicrosoftWebsites,
	}
	available, err := client.CheckNameAvailability(ctx, availabilityRequest)
	if err != nil {
		return fmt.Errorf("Error checking if the name %q was available: %+v", name, err)
	}

	if !*available.NameAvailable {
		return fmt.Errorf("The name %q used for the App Service needs to be globally unique and isn't available: %s", name, *available.Message)
	}

	resGroup := d.Get("resource_group_name").(string)
	location := d.Get("location").(string)
	appServicePlanId := d.Get("app_service_plan_id").(string)
	enabled := d.Get("enabled").(bool)
	tags := d.Get("tags").(map[string]interface{})

	siteConfig := expandAppServiceSiteConfig(d)

	siteEnvelope := web.Site{
		Location: &location,
		Tags:     expandTags(tags),
		SiteProperties: &web.SiteProperties{
			ServerFarmID: utils.String(appServicePlanId),
			Enabled:      utils.Bool(enabled),
			SiteConfig:   &siteConfig,
		},
	}

	if v, ok := d.GetOkExists("client_affinity_enabled"); ok {
		enabled := v.(bool)
		siteEnvelope.SiteProperties.ClientAffinityEnabled = utils.Bool(enabled)
	}

	// NOTE: these seem like sensible defaults, in lieu of any better documentation.
	skipDNSRegistration := false
	forceDNSRegistration := false
	skipCustomDomainVerification := true
	ttlInSeconds := "60"
	createFuture, err := client.CreateOrUpdate(ctx, resGroup, name, siteEnvelope, &skipDNSRegistration, &skipCustomDomainVerification, &forceDNSRegistration, ttlInSeconds)
	if err != nil {
		return err
	}

	err = createFuture.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read App Service %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppServiceUpdate(d, meta)
}

func resourceArmAppServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["sites"]

	if d.HasChange("site_config") {
		// update the main configuration
		siteConfig := expandAppServiceSiteConfig(d)
		siteConfigResource := web.SiteConfigResource{
			SiteConfig: &siteConfig,
		}
		_, err := client.CreateOrUpdateConfiguration(ctx, resGroup, name, siteConfigResource)
		if err != nil {
			return fmt.Errorf("Error updating Configuration for App Service %q: %+v", name, err)
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

		_, err := client.Update(
			ctx,
			resGroup,
			name,
			sitePatchResource,
			nil,
			nil,
			nil,
			"")

		if err != nil {
			return fmt.Errorf("Error updating App Service ARR Affinity setting %q: %+v", name, err)
		}
	}

	if d.HasChange("app_settings") {
		// update the AppSettings
		appSettings := expandAppServiceAppSettings(d)
		settings := web.StringDictionary{
			Properties: appSettings,
		}

		_, err := client.UpdateApplicationSettings(ctx, resGroup, name, settings)
		if err != nil {
			return fmt.Errorf("Error updating Application Settings for App Service %q: %+v", name, err)
		}
	}

	if d.HasChange("connection_string") {
		// update the ConnectionStrings
		connectionStrings := expandAppServiceConnectionStrings(d)
		properties := web.ConnectionStringDictionary{
			Properties: connectionStrings,
		}

		_, err := client.UpdateConnectionStrings(ctx, resGroup, name, properties)
		if err != nil {
			return fmt.Errorf("Error updating Connection Strings for App Service %q: %+v", name, err)
		}
	}

	return resourceArmAppServiceRead(d, meta)
}

func resourceArmAppServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["sites"]

	ctx := meta.(*ArmClient).StopContext
	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service %q (resource group %q) was not found - removing from state", name, resGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM App Service %q: %+v", name, err)
	}

	configResp, err := client.GetConfiguration(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service Configuration %q: %+v", name, err)
	}

	appSettingsResp, err := client.ListApplicationSettings(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service AppSettings %q: %+v", name, err)
	}

	connectionStringsResp, err := client.ListConnectionStrings(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service ConnectionStrings %q: %+v", name, err)
	}

	scmResp, err := client.GetSourceControl(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service Source Control %q: %+v", name, err)
	}

	siteCredFuture, err := client.ListPublishingCredentials(ctx, resGroup, name)
	if err != nil {
		return err
	}
	err = siteCredFuture.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}
	siteCredResp, err := siteCredFuture.Result(client)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM App Service Site Credential %q: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))

	if props := resp.SiteProperties; props != nil {
		d.Set("app_service_plan_id", props.ServerFarmID)
		d.Set("client_affinity_enabled", props.ClientAffinityEnabled)
		d.Set("enabled", props.Enabled)
		d.Set("default_site_hostname", props.DefaultHostName)
		d.Set("outbound_ip_addresses", props.OutboundIPAddresses)
	}

	if err := d.Set("app_settings", flattenAppServiceAppSettings(appSettingsResp.Properties)); err != nil {
		return err
	}
	if err := d.Set("connection_string", flattenAppServiceConnectionStrings(connectionStringsResp.Properties)); err != nil {
		return err
	}

	siteConfig := flattenAppServiceSiteConfig(configResp.SiteConfig)
	if err := d.Set("site_config", siteConfig); err != nil {
		return err
	}

	scm := flattenAppServiceSourceControl(scmResp.SiteSourceControlProperties)
	if err := d.Set("source_control", scm); err != nil {
		return err
	}

	siteCred := flattenAppServiceSiteCredential(siteCredResp.UserProperties)
	if err := d.Set("site_credential", siteCred); err != nil {
		return err
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmAppServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["sites"]

	log.Printf("[DEBUG] Deleting App Service %q (resource group %q)", name, resGroup)

	deleteMetrics := true
	deleteEmptyServerFarm := false
	skipDNSRegistration := true
	ctx := meta.(*ArmClient).StopContext
	resp, err := client.Delete(ctx, resGroup, name, &deleteMetrics, &deleteEmptyServerFarm, &skipDNSRegistration)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return err
		}
	}

	return nil
}

func expandAppServiceSiteConfig(d *schema.ResourceData) web.SiteConfig {
	configs := d.Get("site_config").([]interface{})
	siteConfig := web.SiteConfig{}

	if len(configs) == 0 {
		return siteConfig
	}

	config := configs[0].(map[string]interface{})

	if v, ok := config["always_on"]; ok {
		siteConfig.AlwaysOn = utils.Bool(v.(bool))
	}

	if v, ok := config["default_documents"]; ok {
		input := v.([]interface{})

		documents := make([]string, 0)
		for _, document := range input {
			documents = append(documents, document.(string))
		}

		siteConfig.DefaultDocuments = &documents
	}

	if v, ok := config["dotnet_framework_version"]; ok {
		siteConfig.NetFrameworkVersion = utils.String(v.(string))
	}

	if v, ok := config["java_version"]; ok {
		siteConfig.JavaVersion = utils.String(v.(string))
	}

	if v, ok := config["java_container"]; ok {
		siteConfig.JavaContainer = utils.String(v.(string))
	}

	if v, ok := config["java_container_version"]; ok {
		siteConfig.JavaContainerVersion = utils.String(v.(string))
	}

	if v, ok := config["local_mysql_enabled"]; ok {
		siteConfig.LocalMySQLEnabled = utils.Bool(v.(bool))
	}

	if v, ok := config["managed_pipeline_mode"]; ok {
		siteConfig.ManagedPipelineMode = web.ManagedPipelineMode(v.(string))
	}

	if v, ok := config["php_version"]; ok {
		siteConfig.PhpVersion = utils.String(v.(string))
	}

	if v, ok := config["python_version"]; ok {
		siteConfig.PythonVersion = utils.String(v.(string))
	}

	if v, ok := config["remote_debugging_enabled"]; ok {
		siteConfig.RemoteDebuggingEnabled = utils.Bool(v.(bool))
	}

	if v, ok := config["remote_debugging_version"]; ok {
		siteConfig.RemoteDebuggingVersion = utils.String(v.(string))
	}

	if v, ok := config["use_32_bit_worker_process"]; ok {
		siteConfig.Use32BitWorkerProcess = utils.Bool(v.(bool))
	}

	if v, ok := config["websockets_enabled"]; ok {
		siteConfig.WebSocketsEnabled = utils.Bool(v.(bool))
	}

	if v, ok := config["scm_type"]; ok {
		siteConfig.ScmType = web.ScmType(v.(string))
	}

	return siteConfig
}

func flattenAppServiceSiteConfig(input *web.SiteConfig) []interface{} {
	results := make([]interface{}, 0)
	result := make(map[string]interface{}, 0)

	if input == nil {
		log.Printf("[DEBUG] SiteConfig is nil")
		return results
	}

	if input.AlwaysOn != nil {
		result["always_on"] = *input.AlwaysOn
	}

	if input.DefaultDocuments != nil {
		documents := make([]string, 0)
		for _, document := range *input.DefaultDocuments {
			documents = append(documents, document)
		}

		result["default_documents"] = documents
	}

	if input.NetFrameworkVersion != nil {
		result["dotnet_framework_version"] = *input.NetFrameworkVersion
	}

	if input.JavaVersion != nil {
		result["java_version"] = *input.JavaVersion
	}

	if input.JavaContainer != nil {
		result["java_container"] = *input.JavaContainer
	}

	if input.JavaContainerVersion != nil {
		result["java_container_version"] = *input.JavaContainerVersion
	}

	if input.LocalMySQLEnabled != nil {
		result["local_mysql_enabled"] = *input.LocalMySQLEnabled
	}

	result["managed_pipeline_mode"] = string(input.ManagedPipelineMode)

	if input.PhpVersion != nil {
		result["php_version"] = *input.PhpVersion
	}

	if input.PythonVersion != nil {
		result["python_version"] = *input.PythonVersion
	}

	if input.RemoteDebuggingEnabled != nil {
		result["remote_debugging_enabled"] = *input.RemoteDebuggingEnabled
	}

	if input.RemoteDebuggingVersion != nil {
		result["remote_debugging_version"] = *input.RemoteDebuggingVersion
	}

	if input.Use32BitWorkerProcess != nil {
		result["use_32_bit_worker_process"] = *input.Use32BitWorkerProcess
	}

	if input.WebSocketsEnabled != nil {
		result["websockets_enabled"] = *input.WebSocketsEnabled
	}

	result["scm_type"] = string(input.ScmType)

	return append(results, result)
}

func flattenAppServiceSourceControl(input *web.SiteSourceControlProperties) []interface{} {
	results := make([]interface{}, 0)
	result := make(map[string]interface{}, 0)

	if input == nil {
		log.Printf("[DEBUG] SiteSourceControlProperties is nil")
		return results
	}

	if input.RepoURL != nil {
		result["repo_url"] = *input.RepoURL
	}
	if input.Branch != nil && *input.Branch != "" {
		result["branch"] = *input.Branch
	} else {
		result["branch"] = "master"
	}

	return append(results, result)
}

func expandAppServiceAppSettings(d *schema.ResourceData) *map[string]*string {
	input := d.Get("app_settings").(map[string]interface{})
	output := make(map[string]*string, len(input))

	for k, v := range input {
		output[k] = utils.String(v.(string))
	}

	return &output
}

func expandAppServiceConnectionStrings(d *schema.ResourceData) *map[string]*web.ConnStringValueTypePair {
	input := d.Get("connection_string").([]interface{})
	output := make(map[string]*web.ConnStringValueTypePair, len(input))

	for _, v := range input {
		vals := v.(map[string]interface{})

		csName := vals["name"].(string)
		csType := vals["type"].(string)
		csValue := vals["value"].(string)

		output[csName] = &web.ConnStringValueTypePair{
			Value: utils.String(csValue),
			Type:  web.ConnectionStringType(csType),
		}
	}

	return &output
}

func flattenAppServiceConnectionStrings(input *map[string]*web.ConnStringValueTypePair) interface{} {
	results := make([]interface{}, 0)

	for k, v := range *input {
		result := make(map[string]interface{}, 0)
		result["name"] = k
		result["type"] = string(v.Type)
		result["value"] = *v.Value
		results = append(results, result)
	}

	return results
}

func flattenAppServiceAppSettings(input *map[string]*string) map[string]string {
	output := make(map[string]string, 0)
	for k, v := range *input {
		output[k] = *v
	}

	return output
}

func validateAppServiceName(v interface{}, k string) (ws []string, es []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[0-9a-zA-Z-]+$`).Match([]byte(value)); !matched {
		es = append(es, fmt.Errorf("%q may only contain alphanumeric characters and dashes", k))
	}

	return
}

func flattenAppServiceSiteCredential(input *web.UserProperties) []interface{} {
	results := make([]interface{}, 0)
	result := make(map[string]interface{}, 0)

	if input == nil {
		log.Printf("[DEBUG] UserProperties is nil")
		return results
	}

	if input.PublishingUserName != nil {
		result["username"] = *input.PublishingUserName
	}

	if input.PublishingPassword != nil {
		result["password"] = *input.PublishingPassword
	}

	return append(results, result)
}
