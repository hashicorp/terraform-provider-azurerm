// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceSpringCloudApp() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		DeprecationMessage: "Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_app` data source is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.",

		Read: dataSourceSpringCloudAppRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.SpringCloudAppName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"service_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.SpringCloudServiceName,
			},

			"addon_json": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"custom_persistent_disk": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"storage_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"mount_path": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"share_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"mount_options": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"read_only_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"https_only": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

			"ingress_settings": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"backend_protocol": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"read_timeout_in_seconds": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"send_timeout_in_seconds": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"session_affinity": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"session_cookie_max_age": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"is_public": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"persistent_disk": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"mount_path": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"size_in_gb": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"public_endpoint_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tls_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSpringCloudAppRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.AppsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := appplatform.NewAppID(subscriptionId, d.Get("resource_group_name").(string), d.Get("service_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.AppName, "")
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.AppName)
	d.Set("service_name", id.SpringName)
	d.Set("resource_group_name", id.ResourceGroupName)
	identity, err := flattenSpringCloudAppIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %s", err)
	}

	if prop := resp.Properties; prop != nil {
		d.Set("fqdn", prop.Fqdn)
		d.Set("https_only", prop.HTTPSOnly)
		d.Set("is_public", prop.Public)
		d.Set("url", prop.URL)
		d.Set("tls_enabled", prop.EnableEndToEndTLS)

		if err := d.Set("persistent_disk", flattenSpringCloudAppPersistentDisk(prop.PersistentDisk)); err != nil {
			return fmt.Errorf("setting `persistent_disk`: %s", err)
		}

		if err := d.Set("addon_json", flattenSpringCloudAppAddon(prop.AddonConfigs)); err != nil {
			return fmt.Errorf("setting `addon_json`: %s", err)
		}

		if err := d.Set("custom_persistent_disk", flattenAppCustomPersistentDiskResourceArray(prop.CustomPersistentDisks)); err != nil {
			return fmt.Errorf("setting `custom_persistent_disk`: %+v", err)
		}

		if err := d.Set("ingress_settings", flattenSpringCloudAppIngressSettings(prop.IngressSettings)); err != nil {
			return fmt.Errorf("setting `ingress_settings`: %+v", err)
		}

		if prop.VnetAddons != nil {
			d.Set("public_endpoint_enabled", pointer.From(prop.VnetAddons.PublicEndpoint))
		}
	}

	return nil
}
