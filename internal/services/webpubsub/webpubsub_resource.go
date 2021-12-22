package webpubsub

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/webpubsub/mgmt/2021-10-01/webpubsub"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/webpubsub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/webpubsub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceWebPubSub() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceWebPubSubCreateUpdate,
		Read:   resourceWebPubSubRead,
		Update: resourceWebPubSubCreateUpdate,
		Delete: resourceWebPubSubDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WebPubsubID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateWebpubsubName(),
			},
			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Standard_S1",
					"Free_F1",
				}, false),
			},

			"capacity": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 5, 10, 20, 50, 100}),
			},

			"live_trace_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},
						"categories": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"ConnectivityLogs",
											"MessagingLogs",
											"HttpRequestLogs",
										}, false),
									},
									"enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  true,
									},
								},
							},
						},
					},
				},
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"local_auth_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"aad_auth_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"tls_client_cert_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"public_port": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"server_port": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"external_ip": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"hostname": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceWebPubSubCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Webpubsub.WebPubsubClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewWebPubsubID(subscriptionId, resourceGroup, name)
	liveTraceConfig := d.Get("live_trace_configuration").([]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.WebPubSubName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %q: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_web_pubsub", id.ID())
		}
	}

	parameters := webpubsub.ResourceType{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &webpubsub.Properties{
			LiveTraceConfiguration: expandLiveTraceConfig(liveTraceConfig),
			PublicNetworkAccess:    utils.String("Enabled"),
			DisableAadAuth:         utils.Bool(d.Get("aad_auth_enabled").(bool)),
			DisableLocalAuth:       utils.Bool(d.Get("local_auth_enabled").(bool)),
			TLS: &webpubsub.TLSSettings{
				ClientCertEnabled: utils.Bool(d.Get("tls_client_cert_enabled").(bool)),
			},
		},
		Sku: &webpubsub.ResourceSku{
			Name:     utils.String(d.Get("sku").(string)),
			Capacity: utils.Int32(int32(d.Get("capacity").(int))),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
		parameters.Properties.PublicNetworkAccess = utils.String("Disabled")
	}

	future, err := client.CreateOrUpdate(ctx, parameters, id.ResourceGroup, id.WebPubSubName)
	if err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceWebPubSubRead(d, meta)
}

func resourceWebPubSubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Webpubsub.WebPubsubClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebPubsubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WebPubSubName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %q does not exists - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	keys, err := client.ListKeys(ctx, id.ResourceGroup, id.WebPubSubName)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", *id, err)
	}

	d.Set("primary_access_key", keys.PrimaryKey)
	d.Set("primary_connection_string", keys.PrimaryConnectionString)
	d.Set("secondary_access_key", keys.SecondaryKey)
	d.Set("secondary_connection_string", keys.SecondaryConnectionString)

	d.Set("name", id.WebPubSubName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("sku", resp.Sku.Name)
	d.Set("capacity", resp.Sku.Capacity)

	if props := resp.Properties; props != nil {
		d.Set("external_ip", props.ExternalIP)
		d.Set("hostname", props.HostName)
		d.Set("public_port", props.PublicPort)
		d.Set("server_port", props.ServerPort)
		d.Set("version", props.Version)
		d.Set("aad_auth_enabled", props.DisableAadAuth)
		d.Set("local_auth_enabled", props.DisableLocalAuth)
		d.Set("public_network_access_enabled", strings.EqualFold(*props.PublicNetworkAccess, "Enabled"))
		d.Set("tls_client_cert_enabled", props.TLS.ClientCertEnabled)

		if err := d.Set("live_trace_configuration", flattenLiveTraceConfig(props.LiveTraceConfiguration)); err != nil {
			return fmt.Errorf("setting `live_trace_configuration`:%+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceWebPubSubDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Webpubsub.WebPubsubClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebPubsubID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.WebPubSubName)
	if err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of %q: %+v", id, err)
		}
	}

	return nil
}

func expandLiveTraceConfig(input []interface{}) *webpubsub.LiveTraceConfiguration {
	resourceCategories := make([]webpubsub.LiveTraceCategory, 0)

	if len(input) != 0 && input[0] != nil {
		v := input[0].(map[string]interface{})
		enabled := "true"
		if !v["enabled"].(bool) {
			enabled = "false"
		}

		for _, item := range v["categories"].(*pluginsdk.Set).List() {
			setting := item.(map[string]interface{})
			liveTraceCategory := webpubsub.LiveTraceCategory{
				Name:    utils.String(setting["name"].(string)),
				Enabled: utils.String("true"),
			}
			if !setting["enabled"].(bool) {
				liveTraceCategory.Enabled = utils.String("false")
			}
			resourceCategories = append(resourceCategories, liveTraceCategory)
		}
		return &webpubsub.LiveTraceConfiguration{
			Enabled:    &enabled,
			Categories: &resourceCategories,
		}
	}

	return nil
}

func flattenLiveTraceConfig(input *webpubsub.LiveTraceConfiguration) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	enabled := strings.EqualFold(*input.Enabled, "true")

	resourceCategories := make([]interface{}, 0)
	if input.Categories != nil {
		for _, item := range *input.Categories {
			block := make(map[string]interface{})

			name := ""
			if item.Name != nil {
				name = *item.Name
			}

			block["name"] = name
			block["enabled"] = strings.EqualFold(*item.Enabled, "true")

			resourceCategories = append(resourceCategories, block)
		}
	}
	return []interface{}{map[string]interface{}{
		"enabled":    enabled,
		"categories": resourceCategories,
	}}
}
