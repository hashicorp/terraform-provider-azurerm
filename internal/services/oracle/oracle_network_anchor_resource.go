// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/networkanchors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/resourceanchors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity

var (
	_ sdk.ResourceWithUpdate   = NetworkAnchorResource{}
	_ sdk.ResourceWithIdentity = NetworkAnchorResource{}
)

type NetworkAnchorResource struct{}

type NetworkAnchorResourceModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	Location          string `tfschema:"location"`

	ResourceAnchorId string       `tfschema:"resource_anchor_id"`
	SubnetId         string       `tfschema:"subnet_id"`
	Zones            zones.Schema `tfschema:"zones"`

	BackupSubnetCidr                 string                   `tfschema:"backup_subnet_cidr"`
	DnsForwardingRule                []DnsForwardingRuleModel `tfschema:"dns_forwarding_rule"`
	DnsListeningEndpointAllowedCidrs []string                 `tfschema:"dns_listening_endpoint_allowed_cidrs"`
	DnsPrivateZoneReplicationEnabled bool                     `tfschema:"dns_private_zone_replication_enabled"`

	DnsForwardingEndpointIpAddress  string            `tfschema:"dns_forwarding_endpoint_ip_address"`
	DnsForwardingEndpointNsgRuleUrl string            `tfschema:"dns_forwarding_endpoint_nsg_rule_url"`
	DnsForwardingRuleUrl            string            `tfschema:"dns_forwarding_rule_url"`
	DnsListeningEndpointIpAddress   string            `tfschema:"dns_listening_endpoint_ip_address"`
	DnsListeningEndpointNsgRuleUrl  string            `tfschema:"dns_listening_endpoint_nsg_rule_url"`
	OciVcnDnsLabel                  string            `tfschema:"oci_vcn_dns_label"`
	Tags                            map[string]string `tfschema:"tags"`
}

type DnsForwardingRuleModel struct {
	DomainNames         []string `tfschema:"domain_names"`
	ForwardingIPAddress string   `tfschema:"forwarding_ip_address"`
}

func (NetworkAnchorResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 24),
				validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-]+$`), "Name may include letters, numbers, or hyphens only"),
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"resource_anchor_id": commonschema.ResourceIDReferenceRequiredForceNew(&resourceanchors.ResourceAnchorId{}),

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"zones": commonschema.ZonesMultipleRequiredForceNew(),

		"backup_subnet_cidr": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsCIDR,
		},

		"dns_forwarding_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"domain_names": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						ForceNew: true,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.DomainName,
						},
					},
					"forwarding_ip_address": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsIPAddress,
					},
				},
			},
		},

		"dns_listening_endpoint_allowed_cidrs": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsCIDR,
			},
		},

		"dns_private_zone_replication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},

		"tags": commonschema.Tags(),
	}
}

func (NetworkAnchorResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dns_forwarding_endpoint_ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"dns_forwarding_endpoint_nsg_rule_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"dns_forwarding_rule_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"dns_listening_endpoint_ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"dns_listening_endpoint_nsg_rule_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"oci_vcn_dns_label": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (NetworkAnchorResource) ModelObject() any {
	return &NetworkAnchorResourceModel{}
}

func (NetworkAnchorResource) Identity() resourceids.ResourceId {
	return &networkanchors.NetworkAnchorId{}
}

func (NetworkAnchorResource) ResourceType() string {
	return "azurerm_oracle_network_anchor"
}

func (r NetworkAnchorResource) Create() sdk.ResourceFunc {
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

func (r NetworkAnchorResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.NetworkAnchors
			id, err := networkanchors.ParseNetworkAnchorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model NetworkAnchorResourceModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			update := *existing.Model

			if metadata.ResourceData.HasChange("tags") {
				update.Tags = pointer.To(model.Tags)
			}

			if metadata.ResourceData.HasChange("backup_subnet_cidr") {
				update.Properties.OciBackupCidrBlock = nil
				if model.BackupSubnetCidr != "" {
					update.Properties.OciBackupCidrBlock = pointer.To(model.BackupSubnetCidr)
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, update); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r NetworkAnchorResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := networkanchors.ParseNetworkAnchorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			client := metadata.Client.Oracle.OracleClient.NetworkAnchors
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (NetworkAnchorResource) flatten(metadata sdk.ResourceMetaData, id *networkanchors.NetworkAnchorId, model *networkanchors.NetworkAnchor) error {
	state := NetworkAnchorResourceModel{
		Name:              id.NetworkAnchorName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.Tags = pointer.From(model.Tags)
		state.Zones = pointer.From(model.Zones)
		if props := model.Properties; props != nil {
			resourceAnchorID, err := resourceanchors.ParseResourceAnchorIDInsensitively(props.ResourceAnchorId)
			if err != nil {
				return err
			}
			state.ResourceAnchorId = resourceAnchorID.ID()

			subnetID, err := commonids.ParseSubnetIDInsensitively(props.SubnetId)
			if err != nil {
				return err
			}
			state.SubnetId = subnetID.ID()

			state.OciVcnDnsLabel = pointer.From(props.OciVcnDnsLabel)
			state.DnsPrivateZoneReplicationEnabled = pointer.From(props.IsOracleToAzureDnsZoneSyncEnabled)
			state.DnsForwardingEndpointIpAddress = pointer.From(props.DnsForwardingEndpointIPAddress)
			state.DnsForwardingRuleUrl = pointer.From(props.DnsForwardingRulesURL)
			state.DnsForwardingEndpointNsgRuleUrl = pointer.From(props.DnsForwardingEndpointNsgRulesURL)
			state.DnsListeningEndpointIpAddress = pointer.From(props.DnsListeningEndpointIPAddress)
			state.DnsListeningEndpointNsgRuleUrl = pointer.From(props.DnsListeningEndpointNsgRulesURL)

			state.BackupSubnetCidr = pointer.From(props.OciBackupCidrBlock)
			state.DnsForwardingRule = flattenDnsForwardingRules(props.DnsForwardingRules)
			state.DnsListeningEndpointAllowedCidrs = flattenNetworkAnchorCommaSeparatedValues(pointer.From(props.DnsListeningEndpointAllowedCidrs))
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (NetworkAnchorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.NetworkAnchors
			id, err := networkanchors.ParseNetworkAnchorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (NetworkAnchorResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return networkanchors.ValidateNetworkAnchorID
}

func expandDnsForwardingRules(dnsForwardingRules []DnsForwardingRuleModel) *[]networkanchors.DnsForwardingRule {
	results := make([]networkanchors.DnsForwardingRule, 0)
	for _, item := range dnsForwardingRules {
		results = append(results, networkanchors.DnsForwardingRule{
			DomainNames:         strings.Join(item.DomainNames, ","),
			ForwardingIPAddress: item.ForwardingIPAddress,
		})
	}
	return &results
}

func flattenDnsForwardingRules(dnsForwardingRules *[]networkanchors.DnsForwardingRule) []DnsForwardingRuleModel {
	results := make([]DnsForwardingRuleModel, 0)
	if dnsForwardingRules != nil {
		for _, item := range *dnsForwardingRules {
			results = append(results, DnsForwardingRuleModel{
				DomainNames:         flattenNetworkAnchorCommaSeparatedValues(item.DomainNames),
				ForwardingIPAddress: item.ForwardingIPAddress,
			})
		}
	}

	return results
}

func flattenNetworkAnchorCommaSeparatedValues(input string) []string {
	result := make([]string, 0)
	for value := range strings.SplitSeq(input, ",") {
		if value = strings.TrimSpace(value); value != "" {
			result = append(result, value)
		}
	}
	return result
}
