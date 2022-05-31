package signalr

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/webpubsub/mgmt/2021-10-01/webpubsub"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	identityValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/msi/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceWebPubsubHub() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceWebPubsubHubCreateUpdate,
		Read:   resourceWebPubSubHubRead,
		Update: resourceWebPubsubHubCreateUpdate,
		Delete: resourceWebPubsubHubDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WebPubsubHubID(id)
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
				ValidateFunc: validate.ValidateWebPubsbHubName(),
			},

			"web_pubsub_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WebPubsubID,
			},

			"event_handler": {
				Type:     pluginsdk.TypeSet,
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
											identityValidate.UserAssignedIdentityID,
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

func resourceWebPubsubHubCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.WebPubsubHubsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	webPubsubID, err := parse.WebPubsubID(d.Get("web_pubsub_id").(string))
	if err != nil {
		return fmt.Errorf("parsing ID of %q: %+v", webPubsubID, err)
	}
	id := parse.NewWebPubsubHubID(subscriptionId, webPubsubID.ResourceGroup, webPubsubID.WebPubSubName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.HubName, id.ResourceGroup, id.WebPubSubName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %q: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_web_pubsub_hub", id.ID())
		}
	}

	anonymousPolicyEnabled := "Deny"
	if d.Get("anonymous_connections_enabled").(bool) {
		anonymousPolicyEnabled = "Allow"
	}

	parameters := webpubsub.Hub{
		Properties: &webpubsub.HubProperties{
			EventHandlers:          expandEventHandler(d.Get("event_handler").(*pluginsdk.Set).List()),
			AnonymousConnectPolicy: &anonymousPolicyEnabled,
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.HubName, parameters, id.ResourceGroup, id.WebPubSubName)
	if err != nil {
		return err
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q: %+v", id, err)
	}
	d.SetId(id.ID())

	return resourceWebPubSubHubRead(d, meta)
}

func resourceWebPubSubHubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.WebPubsubHubsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebPubsubHubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.HubName, id.ResourceGroup, id.WebPubSubName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return fmt.Errorf("%q was not found", id)
		}
		return fmt.Errorf("making request on %q: %+v", id, err)
	}

	d.Set("name", id.HubName)
	d.Set("web_pubsub_id", parse.NewWebPubsubID(id.SubscriptionId, id.ResourceGroup, id.WebPubSubName).ID())

	if props := resp.Properties; props != nil {
		if err := d.Set("event_handler", flattenEventHandler(props.EventHandlers)); err != nil {
			return fmt.Errorf("setting `event_handler`: %+v", err)
		}
		d.Set("anonymous_connections_enabled", strings.EqualFold(*props.AnonymousConnectPolicy, "Allow"))
	}

	return nil
}

func resourceWebPubsubHubDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.WebPubsubHubsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebPubsubHubID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.HubName, id.ResourceGroup, id.WebPubSubName)
	if err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deleting %q: %+v", id, err)
		}
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
			URLTemplate: utils.String(block["url_template"].(string)),
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
		urlTemplate := ""
		if item.URLTemplate != nil {
			urlTemplate = *item.URLTemplate
		}

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
			"url_template":       urlTemplate,
			"user_event_pattern": userEventPatten,
			"system_events":      sysEvents,
			"auth":               authBlock,
		})
	}
	return eventHandlerBlock
}

func expandAuth(input []interface{}) *webpubsub.UpstreamAuthSettings {
	if len(input) == 0 || input[0] == nil {
		return &webpubsub.UpstreamAuthSettings{
			Type: webpubsub.UpstreamAuthTypeNone,
		}
	}

	authRaw := input[0].(map[string]interface{})
	authId := authRaw["managed_identity_id"].(string)

	return &webpubsub.UpstreamAuthSettings{
		Type: webpubsub.UpstreamAuthTypeManagedIdentity,
		ManagedIdentity: &webpubsub.ManagedIdentitySettings{
			Resource: &authId,
		},
	}
}

func flattenAuth(input *webpubsub.UpstreamAuthSettings) []interface{} {
	if input == nil || input.Type == webpubsub.UpstreamAuthTypeNone {
		return make([]interface{}, 0)
	}

	authId := input.ManagedIdentity.Resource

	return []interface{}{
		map[string]interface{}{
			"managed_identity_id": authId,
		},
	}
}
