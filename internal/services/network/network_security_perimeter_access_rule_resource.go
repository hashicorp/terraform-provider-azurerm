// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeteraccessrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeterprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = NetworkSecurityPerimeterAccessRuleResource{}

type NetworkSecurityPerimeterAccessRuleResource struct{}

type NetworkSecurityPerimeterAccessRuleResourceModel struct {
	Name                      string   `tfschema:"name"`
	ProfileId                 string   `tfschema:"network_security_perimeter_profile_id"`
	Direction                 string   `tfschema:"direction"`
	AddressPrefixes           []string `tfschema:"address_prefixes"`
	FullyQualifiedDomainNames []string `tfschema:"fqdns"`
	Subscriptions             []string `tfschema:"subscription_ids"`
	ServiceTags               []string `tfschema:"service_tags"`
}

func (NetworkSecurityPerimeterAccessRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`(^[a-zA-Z0-9]+[a-zA-Z0-9_.-]{0,78}[a-zA-Z0-9_]+$)|(^[a-zA-Z0-9]$)`),
				"`name` must be between 1 and 80 characters long, start with a letter or number, end with a letter, number, or underscore, and may contain only letters, numbers, underscores (_), periods (.), or hyphens (-).",
			),
			ForceNew: true,
		},

		"network_security_perimeter_profile_id": {
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
			ExactlyOneOf: []string{"address_prefixes", "fqdns", "service_tags", "subscription_ids"},
		},

		"fqdns": {
			Type: pluginsdk.TypeList,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Optional:     true,
			MinItems:     1,
			ExactlyOneOf: []string{"address_prefixes", "fqdns", "service_tags", "subscription_ids"},
		},

		"service_tags": {
			Type: pluginsdk.TypeList,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Optional:     true,
			MinItems:     1,
			ExactlyOneOf: []string{"address_prefixes", "fqdns", "service_tags", "subscription_ids"},
		},

		"subscription_ids": {
			Type: pluginsdk.TypeList,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: commonids.ValidateSubscriptionID,
			},
			Optional:     true,
			MinItems:     1,
			ExactlyOneOf: []string{"address_prefixes", "fqdns", "service_tags", "subscription_ids"},
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

			if direction == string(networksecurityperimeteraccessrules.AccessRuleDirectionOutbound) {
				if v, ok := rd.GetOk("address_prefixes"); ok && len(v.([]interface{})) > 0 {
					return fmt.Errorf("`address_prefixes` cannot be specified when `direction` is Outbound")
				}

				if v, ok := rd.GetOk("subscription_ids"); ok && len(v.([]interface{})) > 0 {
					return fmt.Errorf("`subscription_ids` cannot be specified when `direction` is Outbound")
				}

				if v, ok := rd.GetOk("service_tags"); ok && len(v.([]interface{})) > 0 {
					return fmt.Errorf("`service_tags` cannot be specified when `direction` is Outbound")
				}
			}

			if direction == string(networksecurityperimeteraccessrules.AccessRuleDirectionInbound) {
				if v, ok := rd.GetOk("fqdns"); ok && len(v.([]interface{})) > 0 {
					return fmt.Errorf("`fqdns` cannot be specified when `direction` is Inbound")
				}
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
				return err
			}

			id := networksecurityperimeteraccessrules.NewAccessRuleID(subscriptionId, profileId.ResourceGroupName, profileId.NetworkSecurityPerimeterName, profileId.ProfileName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := networksecurityperimeteraccessrules.NspAccessRule{
				Properties: &networksecurityperimeteraccessrules.NspAccessRuleProperties{
					Direction:                 pointer.ToEnum[networksecurityperimeteraccessrules.AccessRuleDirection](config.Direction),
					AddressPrefixes:           pointer.To(config.AddressPrefixes),
					FullyQualifiedDomainNames: pointer.To(config.FullyQualifiedDomainNames),
					ServiceTags:               pointer.To(config.ServiceTags),
					Subscriptions:             expandAccessRuleSubscriptionIDs(config.Subscriptions),
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
			if metadata.ResourceData.HasChange("service_tags") {
				existing.Model.Properties.ServiceTags = pointer.To(config.ServiceTags)
			}
			if metadata.ResourceData.HasChange("subscription_ids") {
				existing.Model.Properties.Subscriptions = expandAccessRuleSubscriptionIDs(config.Subscriptions)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
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

			profileId := networksecurityperimeterprofiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.NetworkSecurityPerimeterName, id.ProfileName)

			state := NetworkSecurityPerimeterAccessRuleResourceModel{
				Name:      id.AccessRuleName,
				ProfileId: profileId.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.AddressPrefixes = pointer.From(props.AddressPrefixes)
					state.Direction = string(pointer.From(props.Direction))
					state.FullyQualifiedDomainNames = pointer.From(props.FullyQualifiedDomainNames)
					state.ServiceTags = pointer.From(props.ServiceTags)
					state.Subscriptions = flattenAccessRuleSubscriptionIDs(props.Subscriptions)
				}
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

func expandAccessRuleSubscriptionIDs(subscriptionIDs []string) *[]networksecurityperimeteraccessrules.SubscriptionId {
	if len(subscriptionIDs) == 0 {
		return nil
	}

	result := make([]networksecurityperimeteraccessrules.SubscriptionId, 0, len(subscriptionIDs))
	for _, id := range subscriptionIDs {
		result = append(result, networksecurityperimeteraccessrules.SubscriptionId{
			Id: pointer.To(id),
		})
	}

	return &result
}

func flattenAccessRuleSubscriptionIDs(subscriptions *[]networksecurityperimeteraccessrules.SubscriptionId) []string {
	if subscriptions == nil || len(*subscriptions) == 0 {
		return nil
	}

	result := make([]string, 0, len(*subscriptions))
	for _, s := range *subscriptions {
		if s.Id != nil {
			result = append(result, *s.Id)
		}
	}

	return result
}
