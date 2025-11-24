// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudApp() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		DeprecationMessage: features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_app` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information."),

		Create: resourceSpringCloudAppCreate,
		Read:   resourceSpringCloudAppRead,
		Update: resourceSpringCloudAppUpdate,
		Delete: resourceSpringCloudAppDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SpringCloudAppV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SpringCloudAppID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudAppName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"service_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudServiceName,
			},

			"addon_json": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"custom_persistent_disk": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"mount_path": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"share_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"mount_options": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"read_only_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"is_public": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"https_only": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ingress_settings": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"backend_protocol": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(appplatform.BackendProtocolDefault),
							ValidateFunc: validation.StringInSlice([]string{
								string(appplatform.BackendProtocolDefault),
								string(appplatform.BackendProtocolGRPC),
							}, false),
						},

						"read_timeout_in_seconds": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      300,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"send_timeout_in_seconds": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      60,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"session_affinity": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(appplatform.SessionAffinityNone),
							ValidateFunc: validation.StringInSlice([]string{
								string(appplatform.SessionAffinityCookie),
								string(appplatform.SessionAffinityNone),
							}, false),
						},

						"session_cookie_max_age": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
					},
				},
			},

			"persistent_disk": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"size_in_gb": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 50),
						},

						"mount_path": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "/persistent",
							ValidateFunc: validate.MountPath,
						},
					},
				},
			},

			"public_endpoint_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"tls_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSpringCloudAppCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.AppsClient
	servicesClient := meta.(*clients.Client).AppPlatform.ServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewSpringCloudAppID(subscriptionId, d.Get("resource_group_name").(string), d.Get("service_name").(string), d.Get("name").(string))
	serviceResp, err := servicesClient.Get(ctx, id.ResourceGroup, id.SpringName)
	if err != nil {
		return fmt.Errorf("unable to retrieve %q: %+v", id, err)
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %q: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_spring_cloud_app", id.ID())
		}
	}

	identity, err := expandSpringCloudAppIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return err
	}

	addonConfig, err := expandSpringCloudAppAddon(d.Get("addon_json").(string))
	if err != nil {
		return err
	}

	app := appplatform.AppResource{
		Location: serviceResp.Location,
		Identity: identity,
		Properties: &appplatform.AppResourceProperties{
			AddonConfigs:          addonConfig,
			EnableEndToEndTLS:     utils.Bool(d.Get("tls_enabled").(bool)),
			Public:                utils.Bool(d.Get("is_public").(bool)),
			CustomPersistentDisks: expandAppCustomPersistentDiskResourceArray(d.Get("custom_persistent_disk").([]interface{}), id),
		},
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.AppName, app)
	if err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %q: %+v", id, err)
	}

	// HTTPSOnly and PersistentDisk could only be set by update
	app.Properties.HTTPSOnly = utils.Bool(d.Get("https_only").(bool))
	app.Properties.PersistentDisk = expandSpringCloudAppPersistentDisk(d.Get("persistent_disk").([]interface{}))
	// VNetAddons.PublicEndpoint could only be set by update
	if enabled := d.Get("public_endpoint_enabled").(bool); enabled {
		app.Properties.VnetAddons = &appplatform.AppVNetAddons{
			PublicEndpoint: utils.Bool(enabled),
		}
	}
	// IngressSettings could only be set by update
	// Issue: https://github.com/Azure/azure-rest-api-specs/issues/21536
	app.Properties.IngressSettings = expandSpringCloudAppIngressSetting(d.Get("ingress_settings").([]interface{}))
	future, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.AppName, app)
	if err != nil {
		return fmt.Errorf("update %q: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSpringCloudAppRead(d, meta)
}

func resourceSpringCloudAppUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudAppID(d.Id())
	if err != nil {
		return err
	}

	identity, err := expandSpringCloudAppIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return err
	}

	addonConfig, err := expandSpringCloudAppAddon(d.Get("addon_json").(string))
	if err != nil {
		return err
	}

	app := appplatform.AppResource{
		Identity: identity,
		Properties: &appplatform.AppResourceProperties{
			AddonConfigs:          addonConfig,
			EnableEndToEndTLS:     utils.Bool(d.Get("tls_enabled").(bool)),
			Public:                utils.Bool(d.Get("is_public").(bool)),
			HTTPSOnly:             utils.Bool(d.Get("https_only").(bool)),
			IngressSettings:       expandSpringCloudAppIngressSetting(d.Get("ingress_settings").([]interface{})),
			PersistentDisk:        expandSpringCloudAppPersistentDisk(d.Get("persistent_disk").([]interface{})),
			CustomPersistentDisks: expandAppCustomPersistentDiskResourceArray(d.Get("custom_persistent_disk").([]interface{}), *id),
		},
	}
	if enabled := d.Get("public_endpoint_enabled").(bool); enabled {
		app.Properties.VnetAddons = &appplatform.AppVNetAddons{
			PublicEndpoint: utils.Bool(enabled),
		}
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.AppName, app)
	if err != nil {
		return fmt.Errorf("update %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
	}

	return resourceSpringCloudAppRead(d, meta)
}

func resourceSpringCloudAppRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudAppID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud App %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", id.AppName, id.SpringName, id.ResourceGroup, err)
	}

	d.Set("name", id.AppName)
	d.Set("service_name", id.SpringName)
	d.Set("resource_group_name", id.ResourceGroup)

	identity, err := flattenSpringCloudAppIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %s", err)
	}

	if prop := resp.Properties; prop != nil {
		d.Set("is_public", prop.Public)
		d.Set("https_only", prop.HTTPSOnly)
		d.Set("fqdn", prop.Fqdn)
		d.Set("url", prop.URL)
		d.Set("tls_enabled", prop.EnableEndToEndTLS)
		if err := d.Set("addon_json", flattenSpringCloudAppAddon(prop.AddonConfigs)); err != nil {
			return fmt.Errorf("setting `addon_json`: %s", err)
		}
		if err := d.Set("persistent_disk", flattenSpringCloudAppPersistentDisk(prop.PersistentDisk)); err != nil {
			return fmt.Errorf("setting `persistent_disk`: %s", err)
		}
		if err := d.Set("custom_persistent_disk", flattenAppCustomPersistentDiskResourceArray(prop.CustomPersistentDisks)); err != nil {
			return fmt.Errorf("setting `custom_persistent_disk`: %+v", err)
		}
		if err := d.Set("ingress_settings", flattenSpringCloudAppIngressSettings(prop.IngressSettings)); err != nil {
			return fmt.Errorf("setting `ingress_settings`: %+v", err)
		}
		if prop.VnetAddons != nil {
			d.Set("public_endpoint_enabled", prop.VnetAddons.PublicEndpoint)
		}
	}

	return nil
}

func resourceSpringCloudAppDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudAppID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.AppName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandSpringCloudAppIdentity(input []interface{}) (*appplatform.ManagedIdentityProperties, error) {
	config, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := appplatform.ManagedIdentityProperties{
		Type: appplatform.ManagedIdentityType(string(config.Type)),
	}
	if config.Type == identity.TypeUserAssigned || config.Type == identity.TypeSystemAssignedUserAssigned {
		out.UserAssignedIdentities = make(map[string]*appplatform.UserAssignedManagedIdentity)
		for k := range config.IdentityIds {
			out.UserAssignedIdentities[k] = &appplatform.UserAssignedManagedIdentity{
				// intentionally empty
			}
		}
	}

	return &out, nil
}

func expandSpringCloudAppPersistentDisk(input []interface{}) *appplatform.PersistentDisk {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	raw := input[0].(map[string]interface{})
	return &appplatform.PersistentDisk{
		SizeInGB:  utils.Int32(int32(raw["size_in_gb"].(int))),
		MountPath: utils.String(raw["mount_path"].(string)),
	}
}

func expandAppCustomPersistentDiskResourceArray(input []interface{}, id parse.SpringCloudAppId) *[]appplatform.CustomPersistentDiskResource {
	results := make([]appplatform.CustomPersistentDiskResource, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, appplatform.CustomPersistentDiskResource{
			StorageID: utils.String(parse.NewSpringCloudStorageID(id.SubscriptionId, id.ResourceGroup, id.SpringName, v["storage_name"].(string)).ID()),
			CustomPersistentDiskProperties: &appplatform.AzureFileVolume{
				ShareName:    utils.String(v["share_name"].(string)),
				MountPath:    utils.String(v["mount_path"].(string)),
				ReadOnly:     utils.Bool(v["read_only_enabled"].(bool)),
				MountOptions: utils.ExpandStringSlice(v["mount_options"].(*pluginsdk.Set).List()),
			},
		})
	}
	return &results
}

func expandSpringCloudAppAddon(input string) (map[string]interface{}, error) {
	var addonConfig map[string]interface{}
	if len(input) != 0 {
		err := json.Unmarshal([]byte(input), &addonConfig)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal `addon_json`: %+v", err)
		}
	}
	return addonConfig, nil
}

func expandSpringCloudAppIngressSetting(input []interface{}) *appplatform.IngressSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	raw := input[0].(map[string]interface{})

	return &appplatform.IngressSettings{
		ReadTimeoutInSeconds: utils.Int32(int32(raw["read_timeout_in_seconds"].(int))),
		SendTimeoutInSeconds: utils.Int32(int32(raw["send_timeout_in_seconds"].(int))),
		SessionAffinity:      appplatform.SessionAffinity(raw["session_affinity"].(string)),
		SessionCookieMaxAge:  utils.Int32(int32(raw["session_cookie_max_age"].(int))),
		BackendProtocol:      appplatform.BackendProtocol(raw["backend_protocol"].(string)),
	}
}

func flattenSpringCloudAppIdentity(input *appplatform.ManagedIdentityProperties) (*[]interface{}, error) {
	var transform *identity.SystemAndUserAssignedMap
	if input != nil {
		transform = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
		for k, v := range input.UserAssignedIdentities {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientID,
				PrincipalId: v.PrincipalID,
			}
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(transform)
}

func flattenSpringCloudAppPersistentDisk(input *appplatform.PersistentDisk) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	sizeInGB := 0
	if input.SizeInGB != nil {
		sizeInGB = int(*input.SizeInGB)
	}

	mountPath := ""
	if input.MountPath != nil {
		mountPath = *input.MountPath
	}

	return []interface{}{
		map[string]interface{}{
			"size_in_gb": sizeInGB,
			"mount_path": mountPath,
		},
	}
}

func flattenAppCustomPersistentDiskResourceArray(input *[]appplatform.CustomPersistentDiskResource) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var storageName string
		if item.StorageID != nil {
			// The returned value has inconsistent casing
			// TODO: Remove the normalization codes once the following issue is fixed.
			// Issue: https://github.com/Azure/azure-rest-api-specs/issues/22205
			if id, err := parse.SpringCloudStorageIDInsensitively(*item.StorageID); err == nil {
				storageName = id.StorageName
			}
		}
		var mountPath string
		var shareName string
		var readOnly bool
		var mountOptions *[]string
		if item.CustomPersistentDiskProperties != nil {
			if prop, ok := item.CustomPersistentDiskProperties.AsAzureFileVolume(); ok && prop != nil {
				if prop.MountPath != nil {
					mountPath = *prop.MountPath
				}
				if prop.ShareName != nil {
					shareName = *prop.ShareName
				}
				if prop.ReadOnly != nil {
					readOnly = *prop.ReadOnly
				}
				mountOptions = prop.MountOptions
			}
		}

		results = append(results, map[string]interface{}{
			"storage_name":      storageName,
			"mount_path":        mountPath,
			"share_name":        shareName,
			"mount_options":     set.FromStringSliceNilable(mountOptions),
			"read_only_enabled": readOnly,
		})
	}
	return results
}

func flattenSpringCloudAppAddon(configs map[string]interface{}) *string {
	if len(configs) == 0 {
		return nil
	}
	// The returned value has inconsistent casing
	// TODO: Remove the normalization codes once the following issue is fixed.
	// Issue: https://github.com/Azure/azure-rest-api-specs/issues/22481
	if raw, ok := configs["applicationConfigurationService"]; ok && raw != nil {
		if applicationConfigurationService, ok := raw.(map[string]interface{}); ok && len(applicationConfigurationService) != 0 {
			if resourceId, ok := applicationConfigurationService["resourceId"]; ok && resourceId != nil {
				applicationConfigurationServiceId, err := parse.SpringCloudConfigurationServiceIDInsensitively(resourceId.(string))
				if err == nil {
					applicationConfigurationService["resourceId"] = applicationConfigurationServiceId.ID()
					configs["applicationConfigurationService"] = applicationConfigurationService
				}
			}
		}
	}
	if raw, ok := configs["serviceRegistry"]; ok && raw != nil {
		if serviceRegistry, ok := raw.(map[string]interface{}); ok && len(serviceRegistry) != 0 {
			if resourceId, ok := serviceRegistry["resourceId"]; ok && resourceId != nil {
				serviceRegistryId, err := parse.SpringCloudServiceRegistryIDInsensitively(resourceId.(string))
				if err == nil {
					serviceRegistry["resourceId"] = serviceRegistryId.ID()
					configs["serviceRegistry"] = serviceRegistry
				}
			}
		}
	}
	addonConfig, _ := json.Marshal(configs)
	return utils.String(string(addonConfig))
}

func flattenSpringCloudAppIngressSettings(input *appplatform.IngressSettings) interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}
	var readTimeout, sendTimeout, maxAge int32
	backendProtocol := string(input.BackendProtocol)
	sessionAffinity := string(input.SessionAffinity)
	if input.ReadTimeoutInSeconds != nil {
		readTimeout = *input.ReadTimeoutInSeconds
	}
	if input.SendTimeoutInSeconds != nil {
		sendTimeout = *input.SendTimeoutInSeconds
	}
	if input.SessionCookieMaxAge != nil {
		maxAge = *input.SessionCookieMaxAge
	}
	return []interface{}{map[string]interface{}{
		"backend_protocol":        backendProtocol,
		"read_timeout_in_seconds": readTimeout,
		"send_timeout_in_seconds": sendTimeout,
		"session_affinity":        sessionAffinity,
		"session_cookie_max_age":  maxAge,
	}}
}
