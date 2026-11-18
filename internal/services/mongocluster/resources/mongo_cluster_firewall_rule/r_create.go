// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mongo_cluster_firewall_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2026-06-01/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MongoCluster.FirewallRulesClient

			var state MongoClusterFirewallRuleResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			mongoClusterId, err := firewallrules.ParseMongoClusterID(state.MongoClusterId)
			if err != nil {
				return err
			}

			id := firewallrules.NewFirewallRuleID(mongoClusterId.SubscriptionId, mongoClusterId.ResourceGroupName, mongoClusterId.MongoClusterName, state.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := firewallrules.FirewallRule{
				Properties: &firewallrules.FirewallRuleProperties{
					EndIPAddress:   state.EndIpAddress,
					StartIPAddress: state.StartIpAddress,
				},
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, properties, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
