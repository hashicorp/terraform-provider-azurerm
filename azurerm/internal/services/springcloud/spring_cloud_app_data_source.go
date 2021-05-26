package springcloud

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceSpringCloudApp() *pluginsdk.Resource {
	return &pluginsdk.Resource{
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

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"service_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.SpringCloudServiceName,
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"https_only": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"identity": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"tenant_id": {
							Type:     pluginsdk.TypeString,
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

	id := parse.NewSpringCloudAppID(subscriptionId, d.Get("resource_group_name").(string), d.Get("service_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.AppName)
	d.Set("service_name", id.SpringName)
	d.Set("resource_group_name", id.ResourceGroup)
	if err := d.Set("identity", flattenSpringCloudAppIdentity(resp.Identity)); err != nil {
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
	}

	return nil
}
