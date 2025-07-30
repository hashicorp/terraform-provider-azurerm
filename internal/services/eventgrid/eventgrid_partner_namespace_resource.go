// Copyright (c) HashiCorp, Inc.
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
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/partnernamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/partnerregistrations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = EventGridPartnerNamespaceResource{}

type EventGridPartnerNamespaceResource struct{}

type EventGridPartnerNamespaceResourceModel struct {
	PartnerNamespaceName                string                      `tfschema:"name"`
	Location                            string                      `tfschema:"location"`
	ResourceGroup                       string                      `tfschema:"resource_group_name"`
	InboundIPRules                      []PartnerInboundIpRuleModel `tfschema:"inbound_ip_rule"`
	LocalAuthEnabled                    bool                        `tfschema:"local_auth_enabled"`
	PartnerRegistrationFullyQualifiedID string                      `tfschema:"partner_registration_id"`
	PartnerTopicRoutingMode             string                      `tfschema:"partner_topic_routing_mode"`
	PublicNetworkAccessEnabled          bool                        `tfschema:"public_network_access_enabled"`
	Tags                                map[string]string           `tfschema:"tags"`
}

type PartnerInboundIpRuleModel struct {
	IpMask string `tfschema:"ip_mask"`
	Action string `tfschema:"action"`
}

func (EventGridPartnerNamespaceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": &schema.Schema{
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{3,50}$"),
					"EventGrid Partner Namespace name must be 3 - 50 characters long, contain only letters, numbers and hyphens.",
				),
			),
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"location":            commonschema.Location(),
		"partner_registration_id": &schema.Schema{
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: ValidatePartnerRegistrationFullyQualifiedID,
		},
		"inbound_ip_rule": &schema.Schema{
			Type:     pluginsdk.TypeList,
			Optional: true,
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
		"local_auth_enabled": &schema.Schema{
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},
		"partner_topic_routing_mode": &schema.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(partnernamespaces.PossibleValuesForPartnerTopicRoutingMode(), false),
			Default:      string(partnernamespaces.PartnerTopicRoutingModeChannelNameHeader),
		},
		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},
		"tags": tags.Schema(),
	}
}

func (EventGridPartnerNamespaceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
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

			publicNetworkAccess := partnernamespaces.PublicNetworkAccessEnabled
			if !config.PublicNetworkAccessEnabled {
				publicNetworkAccess = partnernamespaces.PublicNetworkAccessDisabled
			}

			param := partnernamespaces.PartnerNamespace{
				Location: config.Location,
				Properties: &partnernamespaces.PartnerNamespaceProperties{
					DisableLocalAuth:                    pointer.To(!config.LocalAuthEnabled),
					InboundIPRules:                      expandPartnerInboundIPRules(config.InboundIPRules),
					PartnerRegistrationFullyQualifiedId: pointer.To(config.PartnerRegistrationFullyQualifiedID),
					PartnerTopicRoutingMode:             pointer.To(partnernamespaces.PartnerTopicRoutingMode(config.PartnerTopicRoutingMode)),
					PublicNetworkAccess:                 pointer.To(publicNetworkAccess),
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

func expandPartnerInboundIPRules(input []PartnerInboundIpRuleModel) *[]partnernamespaces.InboundIPRule {
	if len(input) == 0 {
		return nil
	}

	ipRules := make([]partnernamespaces.InboundIPRule, 0)
	for _, v := range input {
		ipRules = append(ipRules, partnernamespaces.InboundIPRule{
			Action: pointer.To(partnernamespaces.IPActionType(v.Action)),
			IPMask: pointer.To(v.IpMask),
		})
	}
	return &ipRules
}

func flattenPartnerInboundIPRules(ipRules *[]partnernamespaces.InboundIPRule) []PartnerInboundIpRuleModel {
	output := make([]PartnerInboundIpRuleModel, 0)

	if ipRules == nil || len(*ipRules) == 0 {
		return output
	}

	for _, v := range *ipRules {
		output = append(output, PartnerInboundIpRuleModel{
			IpMask: pointer.From(v.IPMask),
			Action: string(pointer.From(v.Action)),
		})
	}
	return output
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

			payload := existing.Model

			param := partnernamespaces.PartnerNamespace{
				Location: payload.Location,
				Properties: &partnernamespaces.PartnerNamespaceProperties{
					DisableLocalAuth:                    payload.Properties.DisableLocalAuth,
					InboundIPRules:                      payload.Properties.InboundIPRules,
					PartnerRegistrationFullyQualifiedId: payload.Properties.PartnerRegistrationFullyQualifiedId,
					PartnerTopicRoutingMode:             payload.Properties.PartnerTopicRoutingMode,
					PublicNetworkAccess:                 payload.Properties.PublicNetworkAccess,
				},
				Tags: payload.Tags,
			}

			if metadata.ResourceData.HasChange("local_auth_enabled") {
				param.Properties.DisableLocalAuth = pointer.To(!config.LocalAuthEnabled)
			}
			if metadata.ResourceData.HasChange("inbound_ip_rule") {
				param.Properties.InboundIPRules = expandPartnerInboundIPRules(config.InboundIPRules)
			}
			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				if config.PublicNetworkAccessEnabled {
					param.Properties.PublicNetworkAccess = pointer.To(partnernamespaces.PublicNetworkAccessEnabled)
				} else {
					param.Properties.PublicNetworkAccess = pointer.To(partnernamespaces.PublicNetworkAccessDisabled)
				}
			}
			if metadata.ResourceData.HasChange("tags") {
				param.Tags = pointer.To(config.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, param); err != nil {
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
				return fmt.Errorf("parsing %q: %+v", metadata.ResourceData.Id(), err)
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
				state.Location = location.NormalizeNilable(pointer.To(model.Location))
				if props := model.Properties; props != nil {
					state.LocalAuthEnabled = !pointer.From(props.DisableLocalAuth)
					state.InboundIPRules = flattenPartnerInboundIPRules(props.InboundIPRules)
					state.PartnerRegistrationFullyQualifiedID = pointer.From(props.PartnerRegistrationFullyQualifiedId)
					state.PartnerTopicRoutingMode = string(pointer.From(props.PartnerTopicRoutingMode))
					state.PublicNetworkAccessEnabled = pointer.From(props.PublicNetworkAccess) == partnernamespaces.PublicNetworkAccessEnabled
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

func ValidatePartnerRegistrationFullyQualifiedID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := partnerregistrations.ParsePartnerRegistrationID(v); err != nil {
		errors = append(errors, fmt.Errorf("expected %q to be a valid Partner Registration Fully Qualified ID, got %v: %v", k, i, err))
	}
	return warnings, errors
}
