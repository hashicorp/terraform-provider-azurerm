// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsforwardingrulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/forwardingrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PrivateDNSResolverForwardingRuleModel struct {
	Name                   string                 `tfschema:"name"`
	DnsForwardingRulesetId string                 `tfschema:"dns_forwarding_ruleset_id"`
	DomainName             string                 `tfschema:"domain_name"`
	Enabled                bool                   `tfschema:"enabled"`
	Metadata               map[string]string      `tfschema:"metadata"`
	TargetDnsServers       []TargetDnsServerModel `tfschema:"target_dns_servers"`
}

type TargetDnsServerModel struct {
	IPAddress string `tfschema:"ip_address"`
	Port      int64  `tfschema:"port"`
}

type PrivateDNSResolverForwardingRuleResource struct{}

var _ sdk.ResourceWithUpdate = PrivateDNSResolverForwardingRuleResource{}

func (r PrivateDNSResolverForwardingRuleResource) ResourceType() string {
	return "azurerm_private_dns_resolver_forwarding_rule"
}

func (r PrivateDNSResolverForwardingRuleResource) ModelObject() interface{} {
	return &PrivateDNSResolverForwardingRuleModel{}
}

func (r PrivateDNSResolverForwardingRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return forwardingrules.ValidateForwardingRuleID
}

func (r PrivateDNSResolverForwardingRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dns_forwarding_ruleset_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: dnsforwardingrulesets.ValidateDnsForwardingRulesetID,
		},

		"domain_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_dns_servers": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_address": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"port": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"metadata": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r PrivateDNSResolverForwardingRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PrivateDNSResolverForwardingRuleResource) Create() sdk.ResourceFunc {
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
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
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
			return nil
		},
	}
}

func (r PrivateDNSResolverForwardingRuleResource) Update() sdk.ResourceFunc {
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

func (r PrivateDNSResolverForwardingRuleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.ForwardingRulesClient

			id, err := forwardingrules.ParseForwardingRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := PrivateDNSResolverForwardingRuleModel{
				Name:                   id.ForwardingRuleName,
				DnsForwardingRulesetId: dnsforwardingrulesets.NewDnsForwardingRulesetID(id.SubscriptionId, id.ResourceGroupName, id.DnsForwardingRulesetName).ID(),
			}

			properties := &model.Properties
			state.DomainName = properties.DomainName

			state.Enabled = false
			if properties.ForwardingRuleState != nil && *properties.ForwardingRuleState == forwardingrules.ForwardingRuleStateEnabled {
				state.Enabled = true
			}

			if properties.Metadata != nil {
				state.Metadata = *properties.Metadata
			}

			state.TargetDnsServers = flattenTargetDnsServerModel(&properties.TargetDnsServers)

			return metadata.Encode(&state)
		},
	}
}

func (r PrivateDNSResolverForwardingRuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.ForwardingRulesClient

			id, err := forwardingrules.ParseForwardingRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id, forwardingrules.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandTargetDnsServerModel(inputList []TargetDnsServerModel) *[]forwardingrules.TargetDnsServer {
	var outputList []forwardingrules.TargetDnsServer
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

func flattenTargetDnsServerModel(inputList *[]forwardingrules.TargetDnsServer) []TargetDnsServerModel {
	var outputList []TargetDnsServerModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := TargetDnsServerModel{
			IPAddress: input.IPAddress,
		}

		if input.Port != nil {
			output.Port = *input.Port
		}

		outputList = append(outputList, output)
	}

	return outputList
}
