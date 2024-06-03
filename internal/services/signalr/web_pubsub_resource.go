// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package signalr

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2023-02-01/webpubsub"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/validate"
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

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.WebPubsubV0ToV1{},
		}),
		SchemaVersion: 1,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := webpubsub.ParseWebPubSubID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WebPubSubName(),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Premium_P1",
					"Premium_P2",
					"Standard_S1",
					"Free_F1",
				}, false),
			},

			"capacity": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Default:  1,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 200,
					300, 400, 500, 600, 700, 800, 900, 1000})},

			"live_trace": {
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

						"connectivity_logs_enabled": {
							Type:     pluginsdk.TypeBool,
							Default:  true,
							Optional: true,
						},

						"messaging_logs_enabled": {
							Type:     pluginsdk.TypeBool,
							Default:  true,
							Optional: true,
						},

						"http_request_logs_enabled": {
							Type:     pluginsdk.TypeBool,
							Default:  true,
							Optional: true,
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
				Default:  true,
			},

			"aad_auth_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"tls_client_cert_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

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

			"tags": commonschema.Tags(),
		},
	}
}

func resourceWebPubSubCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.WebPubSubClient.WebPubSub
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := webpubsub.NewWebPubSubID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	liveTraceConfig := d.Get("live_trace").([]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %q: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_web_pubsub", id.ID())
		}
	}

	publicNetworkAcc := "Enabled"
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAcc = "Disabled"
	}

	identity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	parameters := webpubsub.WebPubSubResource{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Identity: identity,
		Properties: &webpubsub.WebPubSubProperties{
			LiveTraceConfiguration: expandLiveTraceConfig(liveTraceConfig),
			PublicNetworkAccess:    utils.String(publicNetworkAcc),
			DisableAadAuth:         utils.Bool(!d.Get("aad_auth_enabled").(bool)),
			DisableLocalAuth:       utils.Bool(!d.Get("local_auth_enabled").(bool)),
			Tls: &webpubsub.WebPubSubTlsSettings{
				ClientCertEnabled: utils.Bool(d.Get("tls_client_cert_enabled").(bool)),
			},
		},
		Sku: &webpubsub.ResourceSku{
			Name:     d.Get("sku").(string),
			Capacity: pointer.To(int64(d.Get("capacity").(int))),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{
			string(webpubsub.ProvisioningStateUpdating),
			string(webpubsub.ProvisioningStateCreating),
			string(webpubsub.ProvisioningStateMoving),
			string(webpubsub.ProvisioningStateRunning),
		},
		Target:                    []string{string(webpubsub.ProvisioningStateSucceeded)},
		Refresh:                   webPubsubProvisioningStateRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
		PollInterval:              10 * time.Second,
		ContinuousTargetOccurence: 5,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return err
	}

	d.SetId(id.ID())
	return resourceWebPubSubRead(d, meta)
}

func resourceWebPubSubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.WebPubSubClient.WebPubSub
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webpubsub.ParseWebPubSubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Web Pubsub %s does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	keys, err := client.ListKeys(ctx, *id)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", *id, err)
	}

	d.Set("name", id.WebPubSubName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		skuName := ""
		skuCapacity := int64(0)
		if model.Sku != nil {
			skuName = model.Sku.Name
			skuCapacity = *model.Sku.Capacity
		}
		d.Set("sku", skuName)
		d.Set("capacity", skuCapacity)

		if props := model.Properties; props != nil {
			d.Set("external_ip", props.ExternalIP)
			d.Set("hostname", props.HostName)
			d.Set("public_port", props.PublicPort)
			d.Set("server_port", props.ServerPort)
			d.Set("version", props.Version)

			aadAuthEnabled := true
			if props.DisableAadAuth != nil {
				aadAuthEnabled = !(*props.DisableAadAuth)
			}
			d.Set("aad_auth_enabled", aadAuthEnabled)

			disableLocalAuth := false
			if props.DisableLocalAuth != nil {
				disableLocalAuth = !(*props.DisableLocalAuth)
			}
			d.Set("local_auth_enabled", disableLocalAuth)

			publicNetworkAccessEnabled := true
			if props.PublicNetworkAccess != nil {
				publicNetworkAccessEnabled = strings.EqualFold(*props.PublicNetworkAccess, "Enabled")
			}
			d.Set("public_network_access_enabled", publicNetworkAccessEnabled)

			tlsClientCertEnabled := false
			if props.Tls != nil && props.Tls.ClientCertEnabled != nil {
				tlsClientCertEnabled = *props.Tls.ClientCertEnabled
			}
			d.Set("tls_client_cert_enabled", tlsClientCertEnabled)

			if err := d.Set("live_trace", flattenLiveTraceConfig(props.LiveTraceConfiguration)); err != nil {
				return fmt.Errorf("setting `live_trace`:%+v", err)
			}

			identity, err := identity.FlattenSystemOrUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}
			if err := d.Set("identity", identity); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			if err := tags.FlattenAndSet(d, model.Tags); err != nil {
				return fmt.Errorf("setting `tags`: %+v", err)
			}
		}
	}

	if model := keys.Model; model != nil {
		d.Set("primary_access_key", model.PrimaryKey)
		d.Set("primary_connection_string", model.PrimaryConnectionString)
		d.Set("secondary_access_key", model.SecondaryKey)
		d.Set("secondary_connection_string", model.SecondaryConnectionString)
	}

	return nil
}

func resourceWebPubSubDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.WebPubSubClient.WebPubSub
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webpubsub.ParseWebPubSubID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	return nil
}

func expandLiveTraceConfig(input []interface{}) *webpubsub.LiveTraceConfiguration {
	resourceCategories := make([]webpubsub.LiveTraceCategory, 0)
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	enabled := "false"
	if v["enabled"].(bool) {
		enabled = "true"
	}

	messageLogEnabled := "false"
	if v["messaging_logs_enabled"].(bool) {
		messageLogEnabled = "true"
	}
	resourceCategories = append(resourceCategories, webpubsub.LiveTraceCategory{
		Name:    utils.String("MessagingLogs"),
		Enabled: utils.String(messageLogEnabled),
	})

	connectivityLogEnabled := "false"
	if v["connectivity_logs_enabled"].(bool) {
		connectivityLogEnabled = "true"
	}
	resourceCategories = append(resourceCategories, webpubsub.LiveTraceCategory{
		Name:    utils.String("ConnectivityLogs"),
		Enabled: utils.String(connectivityLogEnabled),
	})

	httpLogEnabled := "false"
	if v["http_request_logs_enabled"].(bool) {
		httpLogEnabled = "true"
	}
	resourceCategories = append(resourceCategories, webpubsub.LiveTraceCategory{
		Name:    utils.String("HttpRequestLogs"),
		Enabled: utils.String(httpLogEnabled),
	})

	return &webpubsub.LiveTraceConfiguration{
		Enabled:    &enabled,
		Categories: &resourceCategories,
	}
}

func flattenLiveTraceConfig(input *webpubsub.LiveTraceConfiguration) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	var enabled bool
	if input.Enabled != nil {
		enabled = strings.EqualFold(*input.Enabled, "true")
	}

	var (
		messagingLogEnabled    bool
		connectivityLogEnabled bool
		httpLogsEnabled        bool
	)

	if input.Categories != nil {
		for _, item := range *input.Categories {
			name := ""
			if item.Name != nil {
				name = *item.Name
			}

			var cateEnabled string
			if item.Enabled != nil {
				cateEnabled = *item.Enabled
			}

			switch name {
			case "MessagingLogs":
				messagingLogEnabled = strings.EqualFold(cateEnabled, "true")
			case "ConnectivityLogs":
				connectivityLogEnabled = strings.EqualFold(cateEnabled, "true")
			case "HttpRequestLogs":
				httpLogsEnabled = strings.EqualFold(cateEnabled, "true")
			default:
				continue
			}
		}
	}
	return []interface{}{map[string]interface{}{
		"enabled":                   enabled,
		"messaging_logs_enabled":    messagingLogEnabled,
		"connectivity_logs_enabled": connectivityLogEnabled,
		"http_request_logs_enabled": httpLogsEnabled,
	}}
}

func webPubsubProvisioningStateRefreshFunc(ctx context.Context, client *webpubsub.WebPubSubClient, id webpubsub.WebPubSubId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)

		provisioningState := "Pending"
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return res, provisioningState, nil
			}
			return nil, "Error", fmt.Errorf("polling for the provisioning state of %s: %+v", id, err)
		}

		if res.Model != nil && res.Model.Properties.ProvisioningState != nil {
			provisioningState = string(*res.Model.Properties.ProvisioningState)
		}

		return res, provisioningState, nil
	}
}
