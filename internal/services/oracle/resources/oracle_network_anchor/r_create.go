// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_network_anchor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/networkanchors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.NetworkAnchors
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model NetworkAnchorResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			id := networkanchors.NewNetworkAnchorID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &networkanchors.NetworkAnchorProperties{
				ResourceAnchorId:                     model.ResourceAnchorId,
				SubnetId:                             model.SubnetId,
				IsOracleDnsForwardingEndpointEnabled: pointer.To(len(model.DnsForwardingRule) > 0),
				IsOracleDnsListeningEndpointEnabled:  pointer.To(len(model.DnsListeningEndpointAllowedCidrs) > 0),
				IsOracleToAzureDnsZoneSyncEnabled:    pointer.To(model.DnsPrivateZoneReplicationEnabled),
			}

			if model.BackupSubnetCidr != "" {
				properties.OciBackupCidrBlock = pointer.To(model.BackupSubnetCidr)
			}

			if len(model.DnsForwardingRule) > 0 {
				properties.DnsForwardingRules = expandDnsForwardingRules(model.DnsForwardingRule)
			}

			if len(model.DnsListeningEndpointAllowedCidrs) > 0 {
				properties.DnsListeningEndpointAllowedCidrs = pointer.To(strings.Join(model.DnsListeningEndpointAllowedCidrs, ","))
			}

			param := networkanchors.NetworkAnchor{
				Name:       pointer.To(model.Name),
				Location:   location.Normalize(model.Location),
				Tags:       pointer.To(model.Tags),
				Zones:      pointer.To(model.Zones),
				Properties: properties,
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, param, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}
