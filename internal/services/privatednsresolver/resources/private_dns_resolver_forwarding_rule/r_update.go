// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_resolver_forwarding_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/forwardingrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.ForwardingRulesClient

			id, err := forwardingrules.ParseForwardingRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PrivateDNSResolverForwardingRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("domain_name") {
				properties.Properties.DomainName = model.DomainName
			}

			if metadata.ResourceData.HasChange("enabled") {
				forwardingRuleState := forwardingrules.ForwardingRuleStateEnabled
				if !model.Enabled {
					forwardingRuleState = forwardingrules.ForwardingRuleStateDisabled
				}
				properties.Properties.ForwardingRuleState = &forwardingRuleState
			}

			if metadata.ResourceData.HasChange("metadata") {
				properties.Properties.Metadata = &model.Metadata
			}

			if metadata.ResourceData.HasChange("target_dns_servers") {
				targetDnsServersValue := expandTargetDnsServerModel(model.TargetDnsServers)

				if targetDnsServersValue != nil {
					properties.Properties.TargetDnsServers = *targetDnsServersValue
				}
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *properties, forwardingrules.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
