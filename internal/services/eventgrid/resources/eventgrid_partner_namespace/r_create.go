// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_namespace

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/partnernamespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
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

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, param, metadata.SetIDCallback(&id)); err != nil {
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
			Action: pointer.ToEnum[partnernamespaces.IPActionType](v.Action),
			IPMask: pointer.To(v.IpMask),
		})
	}
	return &ipRules
}
