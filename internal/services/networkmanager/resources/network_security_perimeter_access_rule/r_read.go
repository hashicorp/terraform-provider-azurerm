// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_security_perimeter_access_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecurityperimeteraccessrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecurityperimeterprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.NetworkSecurityPerimeterAccessRulesClient

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
					state.Direction = pointer.FromEnum(props.Direction)
					state.FullyQualifiedDomainNames = pointer.From(props.FullyQualifiedDomainNames)
					state.ServiceTags = pointer.From(props.ServiceTags)
					state.Subscriptions = flattenAccessRuleSubscriptionIDs(props.Subscriptions)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenAccessRuleSubscriptionIDs(subscriptions *[]networksecurityperimeteraccessrules.SubscriptionId) []string {
	if subscriptions == nil || len(*subscriptions) == 0 {
		return []string{}
	}

	result := make([]string, 0, len(*subscriptions))
	for _, s := range *subscriptions {
		if s.Id != nil {
			result = append(result, *s.Id)
		}
	}

	return result
}
