package webpubsub

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/webpubsub/mgmt/2021-10-01/webpubsub"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
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
			Read:   pluginsdk.DefaultTimeout(30 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				//TODO name restriction - regex
				ValidateFunc: validate.ValidateWebPubsbHubName(),
			},

			"web_pubsub_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateWebpubsubName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"event_handler": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"url_template": {
							Type:     pluginsdk.TypeString,
							Required: true,
							//todo
							//ValidateFunc:http://example.com/api/{hub}/{event}
						},

						//There are 3 kind of patterns supported:
						//1. \"*\", it to matches any event name
						//2. Combine multiple events with \",\", for example \"event1,event2\", it matches event \"event1\" and \"event2\"\r\n
						//3. The single event name, for example, \"event1\", it matches \"event1\"",
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
									"type": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(webpubsub.UpstreamAuthTypeNone),
											string(webpubsub.UpstreamAuthTypeManagedIdentity),
										}, false),
									},

									"managed_identity_resource": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: identityValidate.UserAssignedIdentityID,
									},
								},
							},
						},
					},
				},
			},

			"anonymous_connect_policy": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "Deny",
				ValidateFunc: validation.StringInSlice([]string{
					"Allow",
					"Deny",
				}, false),
			},
		},
	}
}

func resourceWebPubsubHubCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Webpubsub.WebPubsubHubsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM WebPubsub Hub creation.")

	id := parse.NewWebPubsubHubID(subscriptionId, d.Get("resource_group_name").(string), d.Get("web_pubsub_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.HubName, id.ResourceGroup, id.WebPubsubName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing web pubsub hub (%q): %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_web_pubsub_hub", id.ID())
		}
	}

	anonymousPolicy := d.Get("anonymous_connect_policy").(string)
	eventHandlerData := d.Get("event_handler").(*pluginsdk.Set).List()

	eventHandler, err := expandEventHandler(eventHandlerData)
	if err != nil {
		return fmt.Errorf("setting event handler for hub %s: %+v", id, err)
	}

	parameters := webpubsub.Hub{
		Properties: &webpubsub.HubProperties{
			EventHandlers:          eventHandler,
			AnonymousConnectPolicy: &anonymousPolicy,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.HubName, parameters, id.ResourceGroup, id.WebPubsubName); err != nil {
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

	resp, err := client.Get(ctx, id.HubName, id.ResourceGroup, id.WebPubsubName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("making request on %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.HubName)
	d.Set("web_pubsub_name", id.WebPubsubName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.Properties; props != nil && props.EventHandlers != nil {

		if err := d.Set("event_handler", flattenEventHandler(props.EventHandlers)); err != nil {
			return fmt.Errorf("setting `event_handler`: %+v", err)
		}
		d.Set("anonymous_connect_policy", props.AnonymousConnectPolicy)
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

	locks.ByName(id.WebPubsubName, "azurerm_web_pubsub")
	defer locks.UnlockByName(id.WebPubsubName, "azurerm_web_pubsub")

	resp, err := client.Delete(ctx, id.HubName, id.ResourceGroup, id.WebPubsubName)
	if err != nil {
		if !response.WasNotFound(resp.Response()) {
			return fmt.Errorf("deleting web pubsub hub %q (web pubsub %q / Resource Group %q): %+v", id.HubName, id.WebPubsubName, id.ResourceGroup, err)
		}
	}
	return nil
}

func expandEventHandler(input []interface{}) (*[]webpubsub.EventHandler, error) {
	results := make([]webpubsub.EventHandler, 0)
	systemEvents := make([]string, 0)

	if input == nil {
		return &results, nil
	}

	for _, item := range input {
		v := item.(map[string]interface{})

		urlTemplate := v["url_template"].(string)
		userEventPattern := v["user_event_pattern"].(string)
		systemEventsRaw := v["system_events"].(*pluginsdk.Set).List()

		for _, item := range systemEventsRaw {
			systemEvents = append(systemEvents, item.(string))
		}

		if userEventPattern == "" && len(systemEvents) == 0 {
			return nil, fmt.Errorf("`user_event_pattern` and `system_events` cannot be null at the same time")
		}

		authRaws := v["auth"].([]interface{})

		authSetting, err := expandAuth(authRaws)

		js, _ := json.Marshal(authSetting)
		log.Printf("DDDSetting%s", js)

		if err != nil {
			return nil, err
		}

		results = append(results, webpubsub.EventHandler{
			URLTemplate:      &urlTemplate,
			SystemEvents:     &systemEvents,
			UserEventPattern: &userEventPattern,
			Auth:             authSetting,
		})
	}

	return &results, nil
}

func flattenEventHandler(input *[]webpubsub.EventHandler) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	systemEvent := make([]string, 0)
	userEventPatten := ""

	for _, item := range *input {
		urlTemplate := *item.URLTemplate

		if item.UserEventPattern != nil {
			userEventPatten = *item.UserEventPattern
		}

		if item.SystemEvents != nil {
			for _, item := range *item.SystemEvents {
				systemEvent = append(systemEvent, item)
			}
		}
		systemEventPatten := utils.FlattenStringSlice(&systemEvent)

		results = append(results, map[string]interface{}{
			"url_template":       urlTemplate,
			"user_event_pattern": userEventPatten,
			"system_events":      systemEventPatten,
			"auth":               flattenAuth(item.Auth),
		})
	}

	return results
}

func flattenAuth(input *webpubsub.UpstreamAuthSettings) []interface{} {

	if input == nil {
		return make([]interface{}, 0)
	}

	authType := input.Type
	authId := input.ManagedIdentity.Resource

	return []interface{}{
		map[string]interface{}{
			"type":                      authType,
			"managed_identity_resource": authId,
		},
	}
}

func expandAuth(input []interface{}) (*webpubsub.UpstreamAuthSettings, error) {
	if len(input) == 0 {
		return nil, nil
	}

	authRaw := input[0].(map[string]interface{})

	if authType, ok := authRaw["type"].(string); ok {
		authId := authRaw["managed_identity_resource"].(string)

		if authType == string(webpubsub.UpstreamAuthTypeManagedIdentity) && authId == "" {
			return nil, fmt.Errorf("managed_identity_resource is required when the auth_type is set to `managedIdentity")
		} else if authType == string(webpubsub.UpstreamAuthTypeNone) && authId != "" {
			return nil, fmt.Errorf("managed_identity_type is set to None, no auth_id is needed")
		} else {
			return &webpubsub.UpstreamAuthSettings{
				Type: webpubsub.UpstreamAuthType(authType),
				ManagedIdentity: &webpubsub.ManagedIdentitySettings{
					Resource: &authId,
				},
			}, nil
		}
	}

	return nil, nil
}
