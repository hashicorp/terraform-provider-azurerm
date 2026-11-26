// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_domain

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func eventgridDomainSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the EventGrid Domain resource.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{3,50}$"),
					"EventGrid domain name must be 3 - 50 characters long, contain only letters, numbers and hyphens.",
				),
			),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"input_schema": {
			Description:  "The schema in which incoming events will be published to this domain.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(domains.InputSchemaEventGridSchema),
			ValidateFunc: validation.StringInSlice(domains.PossibleValuesForInputSchema(), false),
		},

		// lintignore:XS003
		"input_mapping_fields": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Description: "The ID of the EventGrid Event to associate with the domain.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
					},
					"topic": {
						Description: "The topic of the EventGrid Event to associate with the domain.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
					},
					"event_time": {
						Description: "The event time of the EventGrid Event to associate with the domain.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
					},
					"event_type": {
						Description: "The event type of the EventGrid Event to associate with the domain.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
					},
					"subject": {
						Description: "The subject of the EventGrid Event to associate with the domain.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
					},
					"data_version": {
						Description: "The data version of the EventGrid Event to associate with the domain.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
					},
				},
			},
		},

		// lintignore:XS003
		"input_mapping_default_values": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"event_type": {
						Description: "The default event type of the EventGrid Event to associate with the domain.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
					},
					"subject": {
						Description: "The default subject of the EventGrid Event to associate with the domain.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
					},
					"data_version": {
						Description: "The default data version of the EventGrid Event to associate with the domain.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
					},
				},
			},
		},

		"public_network_access_enabled": {
			Description: "Whether public network access is allowed for this server.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"local_auth_enabled": {
			Description: "Whether local authentication methods are enabled for the EventGrid Domain.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"data_residency_boundary": {
			Description: "The data residency boundary.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// NOTE: O+C because Azure defaults value when not specified, to avoid breaking existing deployments, we will not set a default value here.
			Computed: true,
			ValidateFunc: validation.StringInSlice(
				domains.PossibleValuesForDataResidencyBoundary(),
				false,
			),
		},

		"auto_create_topic_with_first_subscription": {
			Description: "Whether to create the domain topic when the first event subscription at the scope of the domain topic is created.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"auto_delete_topic_with_last_subscription": {
			Description: "Whether to delete the domain topic when the last event subscription at the scope of the domain topic is deleted.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"inbound_ip_rule": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			MaxItems:   128,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_mask": {
						Description: "The IP mask (CIDR) to match on.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"action": {
						Description:  "The action to take when the rule is matched.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(domains.IPActionTypeAllow),
						ValidateFunc: validation.StringInSlice(domains.PossibleValuesForIPActionType(), false),
					},
				},
			},
		},

		"minimum_tls_version": {
			Description: "The minimum TLS version that clients can use.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// NOTE: O+C because Azure defaults value when not specified, to avoid breaking existing deployments, we will not set a default value here.
			Computed:     true,
			ValidateFunc: validation.StringInSlice(domains.PossibleValuesForTlsVersion(), false),
		},

		"endpoint": {
			Description: "The Endpoint associated with the EventGrid Domain.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"primary_access_key": {
			Description: "The Primary Shared Access Key associated with the EventGrid Domain.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_access_key": {
			Description: "The Secondary Shared Access Key associated with the EventGrid Domain.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"tags": commonschema.Tags(),
	}
}
