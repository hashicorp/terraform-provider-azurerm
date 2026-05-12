// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/partnernamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/partnerregistrations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = EventGridPartnerNamespaceResource{}

type EventGridPartnerNamespaceResource struct{}

type EventGridPartnerNamespaceResourceModel struct {
	PartnerNamespaceName                string                      `tfschema:"name"`
	ResourceGroup                       string                      `tfschema:"resource_group_name"`
	Location                            string                      `tfschema:"location"`
	InboundIPRules                      []PartnerInboundIpRuleModel `tfschema:"inbound_ip_rule"`
	LocalAuthEnabled                    bool                        `tfschema:"local_authentication_enabled"`
	PartnerRegistrationFullyQualifiedID string                      `tfschema:"partner_registration_id"`
	PartnerTopicRoutingMode             string                      `tfschema:"partner_topic_routing_mode"`
	PublicNetworkAccess                 string                      `tfschema:"public_network_access"`
	Endpoint                            string                      `tfschema:"endpoint"`
	Tags                                map[string]string           `tfschema:"tags"`
}

type PartnerInboundIpRuleModel struct {
	IpMask string `tfschema:"ip_mask"`
	Action string `tfschema:"action"`
}

func (EventGridPartnerNamespaceResource) Arguments() map[string]*pluginsdk.Schema {
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
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(partnernamespaces.IPActionTypeAllow),
						ValidateFunc: validation.StringInSlice([]string{
							string(partnernamespaces.IPActionTypeAllow),
						}, false),
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
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(partnernamespaces.PublicNetworkAccessEnabled),
				string(partnernamespaces.PublicNetworkAccessDisabled),
			}, false),
			Default: string(partnernamespaces.PublicNetworkAccessEnabled),
		},
		"tags": commonschema.Tags(),
	}
}

func (EventGridPartnerNamespaceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r EventGridPartnerNamespaceResource) ModelObject() interface{} {
	return &EventGridPartnerNamespaceResourceModel{}
}

func (EventGridPartnerNamespaceResource) ResourceType() string {
	return "azurerm_eventgrid_partner_namespace"
}

func (r EventGridPartnerNamespaceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.PartnerNamespaces

			subscriptionId := metadata.Client.Account.SubscriptionId

			var config EventGridPartnerNamespaceResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			id := partnernamespaces.NewPartnerNamespaceID(subscriptionId, config.ResourceGroup, config.PartnerNamespaceName)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := partnernamespaces.PartnerNamespace{
				Location: config.Location,
				Properties: &partnernamespaces.PartnerNamespaceProperties{
					DisableLocalAuth:                    pointer.To(!config.LocalAuthEnabled),
					InboundIPRules:                      expandPartnerInboundIPRules(config.InboundIPRules),
					PartnerRegistrationFullyQualifiedId: pointer.To(config.PartnerRegistrationFullyQualifiedID),
					PartnerTopicRoutingMode:             pointer.ToEnum[partnernamespaces.PartnerTopicRoutingMode](config.PartnerTopicRoutingMode),
					PublicNetworkAccess:                 pointer.ToEnum[partnernamespaces.PublicNetworkAccess](config.PublicNetworkAccess),
				},
				Tags: pointer.To(config.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r EventGridPartnerNamespaceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.PartnerNamespaces

			id, err := partnernamespaces.ParsePartnerNamespaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config EventGridPartnerNamespaceResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model.Properties` was nil", *id)
			}

			model := existing.Model

			if metadata.ResourceData.HasChange("local_authentication_enabled") {
				model.Properties.DisableLocalAuth = pointer.To(!config.LocalAuthEnabled)
			}
			if metadata.ResourceData.HasChange("inbound_ip_rule") {
				model.Properties.InboundIPRules = expandPartnerInboundIPRules(config.InboundIPRules)
			}
			if metadata.ResourceData.HasChange("public_network_access") {
				model.Properties.PublicNetworkAccess = pointer.ToEnum[partnernamespaces.PublicNetworkAccess](config.PublicNetworkAccess)
			}
			if metadata.ResourceData.HasChange("tags") {
				model.Tags = pointer.To(config.Tags)
			}

			// endpoint is read-only and will throw an error if we keep it in the update
			model.Properties.Endpoint = nil

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r EventGridPartnerNamespaceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.PartnerNamespaces

			id, err := partnernamespaces.ParsePartnerNamespaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			state := EventGridPartnerNamespaceResourceModel{
				ResourceGroup:        id.ResourceGroupName,
				PartnerNamespaceName: id.PartnerNamespaceName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				if props := model.Properties; props != nil {
					state.LocalAuthEnabled = !pointer.From(props.DisableLocalAuth)
					state.InboundIPRules = flattenPartnerInboundIPRules(props.InboundIPRules)
					state.PartnerRegistrationFullyQualifiedID = pointer.From(props.PartnerRegistrationFullyQualifiedId)
					state.PartnerTopicRoutingMode = pointer.FromEnum(props.PartnerTopicRoutingMode)
					state.PublicNetworkAccess = pointer.FromEnum(props.PublicNetworkAccess)
					state.Endpoint = pointer.From(props.Endpoint)
				}

				if model.Tags != nil {
					state.Tags = pointer.From(model.Tags)
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (r EventGridPartnerNamespaceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.PartnerNamespaces

			id, err := partnernamespaces.ParsePartnerNamespaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (EventGridPartnerNamespaceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return partnernamespaces.ValidatePartnerNamespaceID
}

func expandPartnerInboundIPRules(input []PartnerInboundIpRuleModel) *[]partnernamespaces.InboundIPRule {
	if len(input) == 0 {
		return nil
	}

	ipRules := make([]partnernamespaces.InboundIPRule, 0)
	for _, v := range input {
		ipRules = append(ipRules, partnernamespaces.InboundIPRule{
			Action: pointer.ToEnum[partnernamespaces.IPActionType](v.Action),
			IPMask: pointer.To(v.IpMask),
		})
	}
	return &ipRules
}

func flattenPartnerInboundIPRules(ipRules *[]partnernamespaces.InboundIPRule) []PartnerInboundIpRuleModel {
	output := make([]PartnerInboundIpRuleModel, 0)

	if ipRules == nil {
		return output
	}

	for _, v := range *ipRules {
		output = append(output, PartnerInboundIpRuleModel{
			IpMask: pointer.From(v.IPMask),
			Action: pointer.FromEnum(v.Action),
		})
	}
	return output
}
