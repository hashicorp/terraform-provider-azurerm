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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.NetworkSecurityPerimeterAccessRulesClient

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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
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
