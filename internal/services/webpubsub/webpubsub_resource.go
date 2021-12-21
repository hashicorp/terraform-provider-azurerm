package webpubsub

import (
	"fmt"
	"log"
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
			_, err := parse.WebPubSubID(id)
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
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Standard_S1",
								"Free_F1",
							}, false),
						},

						"tier": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Standard",
								"Free",
							}, false),
						},

						"capacity": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntInSlice([]int{1, 2, 5, 10, 20, 50, 100}),
						},

						"size": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"live_trace_configuration": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "true",
							ValidateFunc: validation.StringInSlice([]string{
								"true",
								"false",
							}, false),
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
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  "true",
										ValidateFunc: validation.StringInSlice([]string{
											"true",
											"false",
										}, false),
									},
								},
							},
						},
					},
				},
			},
			// Enable or disable public network access. Default to "Enabled". When it's Enabled, network ACLs still apply.
			// When it's Disabled, public network access is always disabled no matter what you set in network ACLs.
			"public_network_access": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "Enabled",
				ValidateFunc: validation.StringInSlice([]string{
					"Enabled",
					"Disabled",
				}, false),
			},

			"disable_local_auth": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"disable_aad_auth": {
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

			"host_name_prefix": {
				Type:     pluginsdk.TypeString,
				Optional: true,
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

	id := parse.NewWebPubSubID(subscriptionId, resourceGroup, name)
	sku := d.Get("sku").([]interface{})
	liveTraceConfig := d.Get("live_trace_configuration").(*pluginsdk.Set).List()

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroupId, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing Web Pubsub (%q):%+v", id, err)
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
		},
		Sku:  expandWebPubsubSku(sku),
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	publicNetworkAccess := d.Get("public_network_access").(string)
	if publicNetworkAccess != "" {
		parameters.Properties.PublicNetworkAccess = utils.String(publicNetworkAccess)
	}

	disableAADAuth := d.Get("disable_aad_auth").(bool)
	if disableAADAuth {
		parameters.Properties.DisableAadAuth = utils.Bool(disableAADAuth)
	}

	disableLocalAuth := d.Get("disable_local_auth").(bool)
	if disableLocalAuth {
		parameters.Properties.DisableLocalAuth = utils.Bool(disableLocalAuth)
	}

	tlsCertEnabled := d.Get("tls_client_cert_enabled").(bool)

	tlsSetting := webpubsub.TLSSettings{
		ClientCertEnabled: utils.Bool(tlsCertEnabled),
	}
	parameters.Properties.TLS = &tlsSetting

	future, err := client.CreateOrUpdate(ctx, parameters, id.ResourceGroupId, id.Name)
	if err != nil {
		return fmt.Errorf("creating Web Pubsub (%q): %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of the Web Pubsub (%q):%+v", id, err)
	}

	d.SetId(id.ID())
	return resourceWebPubSubRead(d, meta)
}

func resourceWebPubSubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Webpubsub.WebPubsubClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebPubSubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupId, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Web Pubsub %q does not exists - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retriving Web Pubsub (%q): %+v", id, err)
	}

	keys, err := client.ListKeys(ctx, id.ResourceGroupId, id.Name)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", *id, err)
	}

	d.Set("primary_access_key", keys.PrimaryKey)
	d.Set("primary_connection_string", keys.PrimaryConnectionString)
	d.Set("secondary_access_key", keys.SecondaryKey)
	d.Set("secondary_connection_string", keys.SecondaryConnectionString)

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroupId)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if err = d.Set("sku", flattenWebPubsubSku(resp.Sku)); err != nil {
		return fmt.Errorf("setting `sku`: %+v", err)
	}

	if props := resp.Properties; props != nil {
		d.Set("external_ip", props.ExternalIP)
		d.Set("hostname", props.HostName)
		d.Set("public_port", props.PublicPort)
		d.Set("server_port", props.ServerPort)
		d.Set("version", props.Version)
		d.Set("disable_aad_auth", props.DisableAadAuth)
		d.Set("disable_local_auth", props.DisableLocalAuth)
		d.Set("public_network_access", props.PublicNetworkAccess)
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

	id, err := parse.WebPubSubID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroupId, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
		}
	}

	return nil
}

func expandLiveTraceConfig(input []interface{}) *webpubsub.LiveTraceConfiguration {
	resourceCategories := make([]webpubsub.LiveTraceCategory, 0)

	if len(input) != 0 && input[0] != nil {
		v := input[0].(map[string]interface{})
		enabled := v["enabled"].(string)
		for _, item := range v["categories"].(*pluginsdk.Set).List() {
			setting := item.(map[string]interface{})
			resourceCategory := webpubsub.LiveTraceCategory{
				Name:    utils.String(setting["name"].(string)),
				Enabled: utils.String(setting["enabled"].(string)),
			}
			resourceCategories = append(resourceCategories, resourceCategory)
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

	enabled := "false"
	if input.Enabled != nil {
		enabled = *input.Enabled
	}

	resourceCategories := make([]interface{}, 0)
	if input.Categories != nil {
		for _, item := range *input.Categories {
			block := make(map[string]interface{})

			name := ""
			if item.Name != nil {
				name = *item.Name
			}
			block["name"] = name

			if v := item.Enabled; v != nil {
				block["enabled"] = *v
			}
			resourceCategories = append(resourceCategories, block)
		}
	}
	return []interface{}{map[string]interface{}{
		"enabled":    enabled,
		"categories": resourceCategories,
	}}
}

func expandWebPubsubSku(input []interface{}) *webpubsub.ResourceSku {
	v := input[0].(map[string]interface{})
	return &webpubsub.ResourceSku{
		Name:     utils.String(v["name"].(string)),
		Capacity: utils.Int32(int32(v["capacity"].(int))),
	}
}

func flattenWebPubsubSku(input *webpubsub.ResourceSku) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	capacity := 1
	if input.Capacity != nil {
		capacity = int(*input.Capacity)
	}

	return []interface{}{
		map[string]interface{}{
			"capacity": capacity,
			"name":     input.Name,
		},
	}
}
