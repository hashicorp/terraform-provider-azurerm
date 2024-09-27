// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package signalr

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2023-02-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
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

			"web_pubsub_id": commonschema.ResourceIDReferenceRequiredForceNew(&webpubsub.WebPubSubId{}),

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

			"event_listener": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"user_event_name_filter": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"system_event_name_filter": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"connected",
									"disconnected",
								}, false),
							},
						},

						"eventhub_namespace_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: eventhubValidate.ValidateEventHubNamespaceName(),
						},

						"eventhub_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: eventhubValidate.ValidateEventHubName(),
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

	id := webpubsub.NewHubID(subscriptionId, webPubSubId.ResourceGroupName, webPubSubId.WebPubSubName, d.Get("name").(string))
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

	eventListener, err := expandEventListener(d.Get("event_listener").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding event listener for web pubsub %s: %+v", id, err)
	}

	parameters.Properties.EventListeners = eventListener

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
	d.Set("web_pubsub_id", webpubsub.NewWebPubSubID(id.SubscriptionId, id.ResourceGroupName, id.WebPubSubName).ID())

	if model := resp.Model; model != nil {
		if err := d.Set("event_handler", flattenEventHandler(model.Properties.EventHandlers)); err != nil {
			return fmt.Errorf("setting `event_handler`: %+v", err)
		}
		d.Set("anonymous_connections_enabled", strings.EqualFold(*model.Properties.AnonymousConnectPolicy, "Allow"))
		if err := d.Set("event_listener", flattenEventListener(model.Properties.EventListeners)); err != nil {
			return fmt.Errorf("setting `event_listener`: %+v", err)
		}
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

func expandEventListener(input []interface{}) (*[]webpubsub.EventListener, error) {
	result := make([]webpubsub.EventListener, 0)
	if len(input) == 0 {
		return &result, nil
	}

	for _, eventListenerItem := range input {
		block := eventListenerItem.(map[string]interface{})
		systemEvents := make([]string, 0)
		userEventPattern := ""
		if v, ok := block["user_event_name_filter"]; ok && len(v.([]interface{})) > 0 {
			userEventPatternList := utils.ExpandStringSlice(v.([]interface{}))
			userEventPattern = strings.Join(*userEventPatternList, ",")
		}

		if v, ok := block["system_event_name_filter"]; ok {
			for _, item := range v.([]interface{}) {
				systemEvents = append(systemEvents, item.(string))
			}
		}
		filter := webpubsub.EventNameFilter{
			SystemEvents:     &systemEvents,
			UserEventPattern: utils.String(userEventPattern),
		}

		endpointName := block["eventhub_namespace_name"].(string)
		fullQualifiedName := endpointName + ".servicebus.windows.net"
		if _, ok := block["eventhub_name"]; !ok {
			return nil, fmt.Errorf("no event hub is specified")
		}
		ehName := block["eventhub_name"].(string)
		endpoint := webpubsub.EventHubEndpoint{
			FullyQualifiedNamespace: fullQualifiedName,
			EventHubName:            ehName,
		}

		result = append(result, webpubsub.EventListener{
			Filter:   filter,
			Endpoint: endpoint,
		})
	}
	return &result, nil
}

func flattenEventListener(listener *[]webpubsub.EventListener) []interface{} {
	eventListenerBlocks := make([]interface{}, 0)
	if listener == nil {
		return eventListenerBlocks
	}

	for _, item := range *listener {
		listenerBlock := make(map[string]interface{}, 0)
		// todo use the type Assertion or Type field in sdk to get the different sub-class
		if eventFilter := item.Filter; eventFilter != nil {
			eventNameFilter := item.Filter.(webpubsub.EventNameFilter)
			userNameFilterList := make([]interface{}, 0)
			if eventNameFilter.SystemEvents != nil {
				listenerBlock["system_event_name_filter"] = utils.FlattenStringSlice(eventNameFilter.SystemEvents)
			}
			if eventNameFilter.UserEventPattern != nil && *eventNameFilter.UserEventPattern != "" {
				v := strings.Split(*eventNameFilter.UserEventPattern, ",")
				for _, s := range v {
					userNameFilterList = append(userNameFilterList, s)
				}
				listenerBlock["user_event_name_filter"] = userNameFilterList
			}
		}

		if eventEndpoint := item.Endpoint; eventEndpoint != nil {
			eventhubEndpoint := item.Endpoint.(webpubsub.EventHubEndpoint)
			listenerBlock["eventhub_namespace_name"] = strings.TrimSuffix(eventhubEndpoint.FullyQualifiedNamespace, ".servicebus.windows.net")
			listenerBlock["eventhub_name"] = eventhubEndpoint.EventHubName
		}
		eventListenerBlocks = append(eventListenerBlocks, listenerBlock)
	}

	return eventListenerBlocks
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
