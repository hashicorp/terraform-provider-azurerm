// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_hub

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func webPubsubHubSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Web Pubsub hub service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebPubSubHubName(),
		},

		"web_pubsub_id": commonschema.ResourceIDReferenceRequiredForceNew(&webpubsub.WebPubSubId{}, "Web Pub/Sub Service"),

		"event_handler": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"url_template": {
						Description:  "The Event Handler URL Template. Two predefined parameters `{hub}` and `{event}` are available to use in the template. The value of the EventHandler URL is dynamically calculated when the client request comes in. Example: `http://example.com/api/{hub}/{event}`.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"user_event_pattern": {
						Description:  "The matching event names. There are 3 kinds of patterns supported: * `*` matches any event name * `,` Combine multiple events with `,` for example `event1,event2`, it matches event `event1` and `event2` * The single event name, for example `event1`, it matches `event1`.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"system_events": {
						Description: "The list of system events.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
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
						Optional: true,
						MinItems: 1,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"managed_identity_id": {
									Description:  "The identity ID of the target resource.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
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
						Description: "The list of matching user event names. `[\"*\"]` can be used to match all events.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"system_event_name_filter": {
						Description: "The list of system events.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"connected",
								"disconnected",
							}, false),
						},
					},

					"eventhub_namespace_name": {
						Description:  "The Event Hub namespace name to receive the events.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: eventhubValidate.ValidateEventHubNamespaceName(),
					},

					"eventhub_name": {
						Description:  "The Event Hub name to receive the events.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: eventhubValidate.ValidateEventHubName(),
					},
				},
			},
		},

		"anonymous_connections_enabled": {
			Description: "Whether to allow anonymous connections to this hub.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},
	}
}
