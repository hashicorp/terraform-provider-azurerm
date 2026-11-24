// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_namespace

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/partnernamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/partnerregistrations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{3,50}$"),
					"`name` must be 3 - 50 characters long, contain only letters, numbers and hyphens.",
				),
			),
		},
		"resource_group_name":     commonschema.ResourceGroupName(),
		"location":                commonschema.Location(),
		"partner_registration_id": commonschema.ResourceIDReferenceRequiredForceNew(&partnerregistrations.PartnerRegistrationId{}),
		"inbound_ip_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 16,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_mask": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsCIDR,
					},
					"action": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(partnernamespaces.IPActionTypeAllow),
						ValidateFunc: validation.StringInSlice(partnernamespaces.PossibleValuesForIPActionType(), false),
					},
				},
			},
		},
		"local_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},
		"partner_topic_routing_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(partnernamespaces.PossibleValuesForPartnerTopicRoutingMode(), false),
			Default:      string(partnernamespaces.PartnerTopicRoutingModeChannelNameHeader),
		},
		"public_network_access": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(partnernamespaces.PossibleValuesForPublicNetworkAccess(), false),
			Default:      string(partnernamespaces.PublicNetworkAccessEnabled),
		},
		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
