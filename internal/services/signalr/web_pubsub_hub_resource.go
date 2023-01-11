package signalr

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2021-10-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceWebPubSubHub() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceWebPubSubHubCreateUpdate,
		Read:   resourceWebPubSubHubRead,
		Update: resourceWebPubSubHubCreateUpdate,
		Delete: resourceWebPubSubHubDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := webpubsub.ParseHubID(id)
			return err
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.WebPubsubHubV0ToV1{},
		}),
		SchemaVersion: 1,

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
				ValidateFunc: validate.WebPubSubHubName(),
			},

			"web_pubsub_id": commonschema.ResourceIDReferenceRequiredForceNew(webpubsub.WebPubSubId{}),

			"event_handler": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"url_template": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"user_event_pattern": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"system_events": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"connect",
									"connected",
									"disconnected",
								}, false),
							},
						},

						"auth": {
							Type:     pluginsdk.TypeList,
							MaxItems: 1,
							MinItems: 1,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"managed_identity_id": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.Any(
											validation.IsUUID,
											commonids.ValidateUserAssignedIdentityID,
										),
									},
								},
							},
						},
					},
				},
			},

			"anonymous_connections_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceWebPubSubHubCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.WebPubSubClient.WebPubSub
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	webPubSubIdRaw := d.Get("web_pubsub_id").(string)
	webPubSubId, err := webpubsub.ParseWebPubSubID(webPubSubIdRaw)
	if err != nil {
		return fmt.Errorf("parsing ID of %q: %+v", webPubSubIdRaw, err)
	}

	id := webpubsub.NewHubID(subscriptionId, webPubSubId.ResourceGroupName, webPubSubId.ResourceName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.HubsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %q: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_web_pubsub_hub", id.ID())
		}
	}

	anonymousPolicyEnabled := "Deny"
	if d.Get("anonymous_connections_enabled").(bool) {
		anonymousPolicyEnabled = "Allow"
	}

	parameters := webpubsub.WebPubSubHub{
		Properties: webpubsub.WebPubSubHubProperties{
			EventHandlers:          expandEventHandler(d.Get("event_handler").([]interface{})),
			AnonymousConnectPolicy: &anonymousPolicyEnabled,
		},
	}

	if err := client.HubsCreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceWebPubSubHubRead(d, meta)
}

func resourceWebPubSubHubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.WebPubSubClient.WebPubSub
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webpubsub.ParseHubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.HubsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return fmt.Errorf("%q was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.HubName)
	d.Set("web_pubsub_id", webpubsub.NewWebPubSubID(id.SubscriptionId, id.ResourceGroupName, id.ResourceName).ID())

	if model := resp.Model; model != nil {
		if err := d.Set("event_handler", flattenEventHandler(model.Properties.EventHandlers)); err != nil {
			return fmt.Errorf("setting `event_handler`: %+v", err)
		}
		d.Set("anonymous_connections_enabled", strings.EqualFold(*model.Properties.AnonymousConnectPolicy, "Allow"))
	}

	return nil
}

func resourceWebPubSubHubDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.WebPubSubClient.WebPubSub
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webpubsub.ParseHubID(d.Id())
	if err != nil {
		return err
	}

	if err := client.HubsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	return nil
}

func expandEventHandler(input []interface{}) *[]webpubsub.EventHandler {
	if len(input) == 0 {
		return nil
	}

	results := make([]webpubsub.EventHandler, 0)

	for _, eventHandlerItem := range input {
		block := eventHandlerItem.(map[string]interface{})
		eventHandlerSettings := webpubsub.EventHandler{
			UrlTemplate: block["url_template"].(string),
		}

		if v, ok := block["user_event_pattern"]; ok {
			eventHandlerSettings.UserEventPattern = utils.String(v.(string))
		}

		if v, ok := block["system_events"]; ok {
			systemEvents := make([]string, 0)
			for _, item := range v.(*pluginsdk.Set).List() {
				systemEvents = append(systemEvents, item.(string))
			}
			eventHandlerSettings.SystemEvents = &systemEvents
		}

		if v, ok := block["auth"].([]interface{}); ok {
			if len(v) > 0 {
				eventHandlerSettings.Auth = expandAuth(v)
			}
		}

		results = append(results, eventHandlerSettings)
	}
	return &results
}

func flattenEventHandler(input *[]webpubsub.EventHandler) []interface{} {
	eventHandlerBlock := make([]interface{}, 0)
	if input == nil {
		return eventHandlerBlock
	}

	for _, item := range *input {
		userEventPatten := ""
		if item.UserEventPattern != nil {
			userEventPatten = *item.UserEventPattern
		}

		sysEvents := make([]interface{}, 0)
		if item.SystemEvents != nil {
			sysEvents = utils.FlattenStringSlice(item.SystemEvents)
		}

		authBlock := make([]interface{}, 0)
		if item.Auth != nil {
			authBlock = flattenAuth(item.Auth)
		}

		eventHandlerBlock = append(eventHandlerBlock, map[string]interface{}{
			"url_template":       item.UrlTemplate,
			"user_event_pattern": userEventPatten,
			"system_events":      sysEvents,
			"auth":               authBlock,
		})
	}
	return eventHandlerBlock
}

func expandAuth(input []interface{}) *webpubsub.UpstreamAuthSettings {
	if len(input) == 0 || input[0] == nil {
		authType := webpubsub.UpstreamAuthTypeNone
		return &webpubsub.UpstreamAuthSettings{
			Type: &authType,
		}
	}

	authRaw := input[0].(map[string]interface{})
	authId := authRaw["managed_identity_id"].(string)
	authType := webpubsub.UpstreamAuthTypeManagedIdentity
	return &webpubsub.UpstreamAuthSettings{
		Type: &authType,
		ManagedIdentity: &webpubsub.ManagedIdentitySettings{
			Resource: &authId,
		},
	}
}

func flattenAuth(input *webpubsub.UpstreamAuthSettings) []interface{} {
	if input == nil || input.Type == nil || *input.Type == webpubsub.UpstreamAuthTypeNone || input.ManagedIdentity == nil || input.ManagedIdentity.Resource == nil {
		return make([]interface{}, 0)
	}

	authId := *input.ManagedIdentity.Resource

	return []interface{}{
		map[string]interface{}{
			"managed_identity_id": authId,
		},
	}
}
