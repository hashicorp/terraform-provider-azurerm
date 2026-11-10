// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_namespace_channel_topic

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/channels"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9-]{3,50}$`),
				"must be between 3 and 50 characters and contain only letters, numbers and hyphens",
			),
		},
		"partner_namespace_id": commonschema.ResourceIDReferenceRequiredForceNew(&channels.PartnerNamespaceId{}),
		"partner_topic": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[a-zA-Z0-9-]{3,50}$`),
							"must be between 3 and 50 characters and contain only letters, numbers and hyphens",
						),
					},
					"resource_group_name": commonschema.ResourceGroupName(),
					"source": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"subscription_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},
		"event_type": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 128),
					},
					"display_name": {
						Type: pluginsdk.TypeString,
						// Optional in spec, but API will return 400 "DisplayName is null or empty" if not provided
						// https://github.com/Azure/azure-rest-api-specs/issues/47482
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"data_schema_url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
					"description": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringLenBetween(1, 256),
					},
					"documentation_url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		// Property is writable at / after creation. Azure will set a default activation message if not provided. Designing
		// this as a read-only: after activation, this property is no longer returned by the API, creating a potentially
		// confusing plan diff.
		"activation_message": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		// Property is writable at / after creation. Maximum value and default if not specified is 7 days from
		// creation. Designing this as a read-only: after activation, this property is no longer returned by the
		// API, creating a potentially confusing plan diff.
		"expiration_time_if_not_activated": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"readiness_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			return pluginsdk.CustomDiffInSequence(
				validateUniqueEventTypeNames,
				validateAlphabeticalEventTypeNames,
			)(ctx, metadata.ResourceDiff, metadata)
		},
	}
}

func validateUniqueEventTypeNames(_ context.Context, _ *pluginsdk.ResourceDiff, meta any) error {
	var config EventGridPartnerNamespaceChannelTopicResourceModel
	if err := meta.(sdk.ResourceMetaData).DecodeDiff(&config); err != nil {
		return fmt.Errorf("decoding: %+v", err)
	}
	names := make(map[string]struct{})
	for _, event := range config.EventTypes {
		if event.Name == "" {
			continue
		}
		if _, exists := names[event.Name]; exists {
			return fmt.Errorf("`event_type` names must be unique: %q is specified more than once", event.Name)
		}
		names[event.Name] = struct{}{}
	}
	return nil
}

func validateAlphabeticalEventTypeNames(_ context.Context, _ *pluginsdk.ResourceDiff, meta any) error {
	var config EventGridPartnerNamespaceChannelTopicResourceModel
	if err := meta.(sdk.ResourceMetaData).DecodeDiff(&config); err != nil {
		return fmt.Errorf("decoding: %+v", err)
	}
	previous := ""
	for _, event := range config.EventTypes {
		if event.Name == "" {
			continue
		}
		if event.Name < previous {
			return fmt.Errorf("`event_type` blocks must be ordered alphabetically by `name` (case-sensitive): %q must appear before %q", event.Name, previous)
		}
		previous = event.Name
	}
	return nil
}
