// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_topic

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/topics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func eventgridTopicSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{3,50}$"),
					"EventGrid topic name must be 3 - 50 characters long, contain only letters, numbers and hyphens.",
				),
			),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"input_schema": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(topics.InputSchemaEventGridSchema),
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(topics.PossibleValuesForInputSchema(), false),
		},

		// lintignore:XS003
		"input_mapping_fields": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						ForceNew: true,
						Optional: true,
					},
					"topic": {
						Type:     pluginsdk.TypeString,
						ForceNew: true,
						Optional: true,
					},
					"event_time": {
						Type:     pluginsdk.TypeString,
						ForceNew: true,
						Optional: true,
					},
					"event_type": {
						Type:     pluginsdk.TypeString,
						ForceNew: true,
						Optional: true,
					},
					"subject": {
						Type:     pluginsdk.TypeString,
						ForceNew: true,
						Optional: true,
					},
					"data_version": {
						Type:     pluginsdk.TypeString,
						ForceNew: true,
						Optional: true,
					},
				},
			},
		},

		// lintignore:XS003
		"input_mapping_default_values": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"event_type": {
						Type:     pluginsdk.TypeString,
						ForceNew: true,
						Optional: true,
					},
					"subject": {
						Type:     pluginsdk.TypeString,
						ForceNew: true,
						Optional: true,
					},
					"data_version": {
						Type:     pluginsdk.TypeString,
						ForceNew: true,
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

		"inbound_ip_rule": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			MaxItems:   128,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_mask": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"action": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(topics.IPActionTypeAllow),
						ValidateFunc: validation.StringInSlice(topics.PossibleValuesForIPActionType(), false),
					},
				},
			},
		},

		"data_residency_boundary": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C because Azure defaults value when not specified, to avoid breaking existing deployments, we will not set a default value here.
			Computed: true,
			ValidateFunc: validation.StringInSlice(
				topics.PossibleValuesForDataResidencyBoundary(),
				false,
			),
		},

		"endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"tags": commonschema.Tags(),

		"minimum_tls_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C because Azure defaults value when not specified, to avoid breaking existing deployments, we will not set a default value here.
			Computed: true,
			ValidateFunc: validation.StringInSlice(
				topics.PossibleValuesForTlsVersion(),
				false,
			),
		},
	}
}
