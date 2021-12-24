package webpubsub

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/webpubsub/mgmt/2021-10-01/webpubsub"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	identityValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/msi/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/webpubsub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/webpubsub/validate"
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
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"setting": {
							Type:     pluginsdk.TypeSet,
							Required: true,
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
													Type:         pluginsdk.TypeString,
													Required:     true,
													ValidateFunc: identityValidate.UserAssignedIdentityID,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"anonymous_connect_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceWebPubsubHubCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Webpubsub.WebPubsubHubsClient
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
	if d.Get("anonymous_connect_enabled").(bool) {
		anonymousPolicyEnabled = "Allow"
	}

	parameters := webpubsub.Hub{
		Properties: &webpubsub.HubProperties{
			EventHandlers:          expandEventHandler(d.Get("event_handler").([]interface{})),
			AnonymousConnectPolicy: &anonymousPolicyEnabled,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.HubName, parameters, id.ResourceGroup, id.WebPubSubName); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceWebPubSubHubRead(d, meta)
}

func resourceWebPubSubHubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Webpubsub.WebPubsubHubsClient
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
		d.Set("anonymous_connect_enabled", strings.EqualFold(*props.AnonymousConnectPolicy, "Allow"))
	}

	return nil
}

func resourceWebPubsubHubDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Webpubsub.WebPubsubHubsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebPubsubHubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.HubName, id.ResourceGroup, id.WebPubSubName)
	if err != nil {
		if !response.WasNotFound(resp.Response()) {
			return fmt.Errorf("deleting %q: %+v", id, err)
		}
	}
	return nil
}

func expandEventHandler(input []interface{}) *[]webpubsub.EventHandler {
	if len(input) == 0 {
		return nil
	}

	results := make([]webpubsub.EventHandler, 0)
	block := input[0].(map[string]interface{})

	setting := block["setting"].(*pluginsdk.Set).List()
	for _, r := range setting {
		rblock := r.(map[string]interface{})
		systemEvents := make([]string, 0)

		urlTemplate := rblock["url_template"].(string)
		userEventPattern := rblock["user_event_pattern"].(string)
		systemEventsRaw := rblock["system_events"].(*pluginsdk.Set).List()
		for _, item := range systemEventsRaw {
			systemEvents = append(systemEvents, item.(string))
		}
		authRaws := rblock["auth"].([]interface{})

		results = append(results, webpubsub.EventHandler{
			URLTemplate:      &urlTemplate,
			SystemEvents:     &systemEvents,
			UserEventPattern: &userEventPattern,
			Auth:             expandAuth(authRaws),
		})
	}
	return &results
}

func flattenEventHandler(input *[]webpubsub.EventHandler) []interface{} {
	eventHandlerBlock := make([]interface{}, 0)
	if input == nil {
		return eventHandlerBlock
	}

	setting := make([]interface{}, 0)
	for _, item := range *input {
		settingBlock := make(map[string]interface{})

		urlTemplate := ""
		if item.URLTemplate != nil {
			urlTemplate = *item.URLTemplate
		}
		settingBlock["url_template"] = urlTemplate

		if userEventPatten := item.UserEventPattern; userEventPatten != nil {
			settingBlock["user_event_pattern"] = userEventPatten
		}

		settingBlock["system_events"] = utils.FlattenStringSlice(item.SystemEvents)

		if auth := item.Auth; auth != nil {
			settingBlock["auth"] = flattenAuth(auth)
		}

		setting = append(setting, settingBlock)
	}
	eventHandlerBlock = append(eventHandlerBlock, map[string]interface{}{
		"setting": setting,
	})
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
