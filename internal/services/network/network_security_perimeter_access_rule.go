// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeterprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeteraccessrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = NetworkSecurityPerimeterAccessRuleResource{}

type NetworkSecurityPerimeterAccessRuleResource struct{}

type NetworkSecurityPerimeterAccessRuleResourceModel struct {
	Name                      string   `tfschema:"name"`
	ProfileId                 string   `tfschema:"profile_id"`
	Direction                 string   `tfschema:"direction"`
	AddressPrefixes           []string `tfschema:"address_prefixes"`
	FullyQualifiedDomainNames []string `tfschema:"fqdns"`
	Subscriptions             []string `tfschema:"subscription_ids"`
}

func (NetworkSecurityPerimeterAccessRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ForceNew:     true,
		},

		"profile_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: networksecurityperimeteraccessrules.ValidateProfileID,
			ForceNew:     true,
		},

		"direction": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice(
				networksecurityperimeteraccessrules.PossibleValuesForAccessRuleDirection(),
				false,
			),
			ForceNew: true,
		},

		"address_prefixes": {
			Type: pluginsdk.TypeList,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsCIDR,
			},
			Optional:     true,
			MinItems:     1,
			ExactlyOneOf: []string{"address_prefixes", "fqdns", "subscription_ids"},
		},

		"fqdns": {
			Type: pluginsdk.TypeList,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Optional:     true,
			MinItems:     1,
			ExactlyOneOf: []string{"address_prefixes", "fqdns", "subscription_ids"},
		},

		"subscription_ids": {
			Type: pluginsdk.TypeList,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: commonids.ValidateSubscriptionID,
			},
			Optional:     true,
			MinItems:     1,
			ExactlyOneOf: []string{"address_prefixes", "fqdns", "subscription_ids"},
		},
	}
}

func (NetworkSecurityPerimeterAccessRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (NetworkSecurityPerimeterAccessRuleResource) ModelObject() interface{} {
	return &NetworkSecurityPerimeterAccessRuleResourceModel{}
}

func (NetworkSecurityPerimeterAccessRuleResource) ResourceType() string {
	return "azurerm_network_security_perimeter_access_rule"
}

func (r NetworkSecurityPerimeterAccessRuleResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff
			direction := rd.Get("direction").(string)

			if direction == string(networksecurityperimeteraccessrules.AccessRuleDirectionOutbound) && rd.HasChange("address_prefixes") {
				return fmt.Errorf("`address_prefixes` can only be set when `direction` is Inbound")
			}

			if direction == string(networksecurityperimeteraccessrules.AccessRuleDirectionOutbound) && rd.HasChange("subscription_id") {
				return fmt.Errorf("`subscription_ids` can only be set when `direction` is Inbound")
			}

			if direction == string(networksecurityperimeteraccessrules.AccessRuleDirectionInbound) && rd.HasChange("fqdns") {
				return fmt.Errorf("`fqdns` cannot be specified when `direction` is Outbound")
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r NetworkSecurityPerimeterAccessRuleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{

		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkSecurityPerimeterAccessRulesClient

			subscriptionId := metadata.Client.Account.SubscriptionId

			var config NetworkSecurityPerimeterAccessRuleResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			profileId, err := networksecurityperimeterprofiles.ParseProfileID(config.ProfileId)
			if err != nil {
				return fmt.Errorf("parsing profile ID: %+v", err)
			}
			nspId := networksecurityperimeters.NewNetworkSecurityPerimeterID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.NetworkSecurityPerimeterName)

			id := networksecurityperimeteraccessrules.NewAccessRuleID(subscriptionId, nspId.ResourceGroupName, nspId.NetworkSecurityPerimeterName, profileId.ProfileName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			direction := networksecurityperimeteraccessrules.AccessRuleDirection(config.Direction)

			subscriptions := make([]networksecurityperimeteraccessrules.SubscriptionId, len(config.Subscriptions))
			for i, s := range config.Subscriptions {
				subscriptions[i] = networksecurityperimeteraccessrules.SubscriptionId{
					Id: pointer.To(s),
				}
			}
			param := networksecurityperimeteraccessrules.NspAccessRule{
				Properties: &networksecurityperimeteraccessrules.NspAccessRuleProperties{
					Direction:                 &direction,
					AddressPrefixes:           pointer.To(config.AddressPrefixes),
					FullyQualifiedDomainNames: pointer.To(config.FullyQualifiedDomainNames),
					Subscriptions:             &subscriptions,
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NetworkSecurityPerimeterAccessRuleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkSecurityPerimeterAccessRulesClient

			id, err := networksecurityperimeteraccessrules.ParseAccessRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config NetworkSecurityPerimeterAccessRuleResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			if metadata.ResourceData.HasChange("address_prefixes") {
				existing.Model.Properties.AddressPrefixes = pointer.To(config.AddressPrefixes)
			}
			if metadata.ResourceData.HasChange("fqdns") {
				existing.Model.Properties.FullyQualifiedDomainNames = pointer.To(config.FullyQualifiedDomainNames)
			}
			if metadata.ResourceData.HasChange("direction") {
				direction := networksecurityperimeteraccessrules.AccessRuleDirection(config.Direction)
				existing.Model.Properties.Direction = &direction
			}
			if metadata.ResourceData.HasChange("subscription_ids") {
				subs := make([]networksecurityperimeteraccessrules.SubscriptionId, len(config.Subscriptions))
				for i, s := range config.Subscriptions {
					subs[i] = networksecurityperimeteraccessrules.SubscriptionId{
						Id: pointer.To(s),
					}
				}
				existing.Model.Properties.Subscriptions = &subs
			}

			param := networksecurityperimeteraccessrules.NspAccessRule{
				Properties: existing.Model.Properties,
			}
			if _, err := client.CreateOrUpdate(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (NetworkSecurityPerimeterAccessRuleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkSecurityPerimeterAccessRulesClient

			id, err := networksecurityperimeteraccessrules.ParseAccessRuleID(metadata.ResourceData.Id())
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

			var subscriptions []string
			if resp.Model.Properties.Subscriptions != nil {
				subscriptions = make([]string, len(*resp.Model.Properties.Subscriptions))
				for i, s := range *resp.Model.Properties.Subscriptions {
					if s.Id != nil {
						subscriptions[i] = *s.Id
					}
				}
			}

			profileId := networksecurityperimeterprofiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.NetworkSecurityPerimeterName, id.ProfileName)

			state := NetworkSecurityPerimeterAccessRuleResourceModel{
				Name:                      id.AccessRuleName,
				ProfileId:                 profileId.ID(),
				AddressPrefixes:           pointer.From(resp.Model.Properties.AddressPrefixes),
				Direction:                 string(pointer.From(resp.Model.Properties.Direction)),
				FullyQualifiedDomainNames: pointer.From(resp.Model.Properties.FullyQualifiedDomainNames),
				Subscriptions:             subscriptions,
			}

			return metadata.Encode(&state)
		},
	}
}

func (NetworkSecurityPerimeterAccessRuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkSecurityPerimeterAccessRulesClient

			id, err := networksecurityperimeteraccessrules.ParseAccessRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (NetworkSecurityPerimeterAccessRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return networksecurityperimeteraccessrules.ValidateAccessRuleID
}
