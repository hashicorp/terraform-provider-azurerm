// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_namespace

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/partnernamespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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
