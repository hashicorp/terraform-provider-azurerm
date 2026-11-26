// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_resolver_forwarding_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsforwardingrulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/forwardingrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PrivateDNSResolverForwardingRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.PrivateDnsResolver.ForwardingRulesClient
			dnsForwardingRulesetId, err := dnsforwardingrulesets.ParseDnsForwardingRulesetID(model.DnsForwardingRulesetId)
			if err != nil {
				return err
			}

			id := forwardingrules.NewForwardingRuleID(dnsForwardingRulesetId.SubscriptionId, dnsForwardingRulesetId.ResourceGroupName, dnsForwardingRulesetId.DnsForwardingRulesetName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			forwardingRuleState := forwardingrules.ForwardingRuleStateEnabled
			if !model.Enabled {
				forwardingRuleState = forwardingrules.ForwardingRuleStateDisabled
			}
			properties := &forwardingrules.ForwardingRule{
				Properties: forwardingrules.ForwardingRuleProperties{
					DomainName:          model.DomainName,
					ForwardingRuleState: &forwardingRuleState,
					Metadata:            &model.Metadata,
				},
			}

			targetDnsServersValue := expandTargetDnsServerModel(model.TargetDnsServers)
			if targetDnsServersValue != nil {
				properties.Properties.TargetDnsServers = *targetDnsServersValue
			}

			if _, err := client.CreateOrUpdate(ctx, id, *properties, forwardingrules.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func expandTargetDnsServerModel(inputList []TargetDnsServerModel) *[]forwardingrules.TargetDnsServer {
	outputList := make([]forwardingrules.TargetDnsServer, 0, len(inputList))
	for _, v := range inputList {
		input := v
		output := forwardingrules.TargetDnsServer{
			IPAddress: input.IPAddress,
			Port:      &input.Port,
		}

		outputList = append(outputList, output)
	}

	return &outputList
}
