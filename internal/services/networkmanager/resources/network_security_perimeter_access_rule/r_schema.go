// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_security_perimeter_access_rule

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecurityperimeteraccessrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
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

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff
			direction := rd.Get("direction").(string)

			if direction == string(networksecurityperimeteraccessrules.AccessRuleDirectionOutbound) {
				if v, ok := rd.GetOk("address_prefixes"); ok && len(v.([]any)) > 0 {
					return fmt.Errorf("`address_prefixes` cannot be specified when `direction` is Outbound")
				}

				if v, ok := rd.GetOk("subscription_ids"); ok && len(v.([]any)) > 0 {
					return fmt.Errorf("`subscription_ids` cannot be specified when `direction` is Outbound")
				}

				if v, ok := rd.GetOk("service_tags"); ok && len(v.([]any)) > 0 {
					return fmt.Errorf("`service_tags` cannot be specified when `direction` is Outbound")
				}
			}

			if direction == string(networksecurityperimeteraccessrules.AccessRuleDirectionInbound) {
				if v, ok := rd.GetOk("fqdns"); ok && len(v.([]any)) > 0 {
					return fmt.Errorf("`fqdns` cannot be specified when `direction` is Inbound")
				}
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
