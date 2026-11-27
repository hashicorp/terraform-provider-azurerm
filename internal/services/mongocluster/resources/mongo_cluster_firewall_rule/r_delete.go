// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mongo_cluster_firewall_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2026-06-01/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	mongoclusterClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/mongocluster/client"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := mongoclusterClient.From(metadata.Client).FirewallRulesClient

			id, err := firewallrules.ParseFirewallRuleID(metadata.ResourceData.Id())
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
