// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/networkanchors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = NetworkAnchorResource{}

type NetworkAnchorResource struct{}

type NetworkAnchorResourceModel struct {
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`
	Zones             zones.Schema      `tfschema:"zones"`

	ResourceAnchorId string `tfschema:"resource_anchor_id"`
	SubnetId         string `tfschema:"subnet_id"`

	OciVcnDnsLabel                     string                   `tfschema:"oci_vcn_dns_label"`
	OciBackupCidrBlock                 string                   `tfschema:"oci_backup_cidr_block"`
	OracleDnsForwardingEndpointEnabled bool                     `tfschema:"oracle_dns_forwarding_endpoint_enabled"`
	OracleDnsListeningEndpointEnabled  bool                     `tfschema:"oracle_dns_listening_endpoint_enabled"`
	OracleToAzureDnsZoneSyncEnabled    bool                     `tfschema:"oracle_to_azure_dns_zone_sync_enabled"`
	DnsListeningEndpointAllowedCidr    string                   `tfschema:"dns_listening_endpoint_allowed_cidr"`
	DnsForwardingRule                  []DnsForwardingRuleModel `tfschema:"dns_forwarding_rule"`
}

func (NetworkAnchorResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),

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

		"resource_anchor_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"oci_vcn_dns_label": {
			Type: pluginsdk.TypeString,
			// O+C if the value not specified, this gets set to the name of the Network Anchor
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"oci_backup_cidr_block": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsCIDR,
		},

		"oracle_dns_forwarding_endpoint_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"oracle_dns_listening_endpoint_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"oracle_to_azure_dns_zone_sync_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"dns_listening_endpoint_allowed_cidr": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.IsIpOrCIDRRangeList,
		},

		"dns_forwarding_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"domain_names": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"forwarding_ip_address": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsIPAddress,
					},
				},
			},
		},

		"tags":  commonschema.Tags(),
		"zones": commonschema.ZonesMultipleRequiredForceNew(),
	}
}

func (NetworkAnchorResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (NetworkAnchorResource) ModelObject() interface{} {
	return &NetworkAnchorResource{}
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
				return err
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
				IsOracleDnsForwardingEndpointEnabled: pointer.To(model.OracleDnsForwardingEndpointEnabled),
				IsOracleDnsListeningEndpointEnabled:  pointer.To(model.OracleDnsListeningEndpointEnabled),
				IsOracleToAzureDnsZoneSyncEnabled:    pointer.To(model.OracleToAzureDnsZoneSyncEnabled),
			}

			if model.OciVcnDnsLabel != "" {
				properties.OciVcnDnsLabel = pointer.To(model.OciVcnDnsLabel)
			}

			if model.OciBackupCidrBlock != "" {
				properties.OciBackupCidrBlock = pointer.To(model.OciBackupCidrBlock)
			}

			if model.DnsListeningEndpointAllowedCidr != "" {
				properties.DnsListeningEndpointAllowedCidrs = pointer.To(model.DnsListeningEndpointAllowedCidr)
			}

			if len(model.DnsForwardingRule) > 0 {
				properties.DnsForwardingRules = expandDnsForwardingRules(model.DnsForwardingRule)
			}

			param := networkanchors.NetworkAnchor{
				Name:       pointer.To(model.Name),
				Location:   location.Normalize(model.Location),
				Tags:       pointer.To(model.Tags),
				Zones:      pointer.To(model.Zones),
				Properties: properties,
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
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
				return err
			}

			_, err = client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if r.hasUpdates(metadata) {
				update := &networkanchors.NetworkAnchorUpdate{
					Properties: &networkanchors.NetworkAnchorUpdateProperties{},
				}

				if metadata.ResourceData.HasChange("tags") {
					update.Tags = pointer.To(model.Tags)
				}

				if metadata.ResourceData.HasChange("oci_backup_cidr_block") {
					update.Properties.OciBackupCidrBlock = pointer.To(model.OciBackupCidrBlock)
				}

				if metadata.ResourceData.HasChange("oracle_dns_forwarding_endpoint_enabled") {
					update.Properties.IsOracleDnsForwardingEndpointEnabled = pointer.To(model.OracleDnsForwardingEndpointEnabled)
				}

				if metadata.ResourceData.HasChange("oracle_dns_listening_endpoint_enabled") {
					update.Properties.IsOracleDnsListeningEndpointEnabled = pointer.To(model.OracleDnsListeningEndpointEnabled)
				}

				if metadata.ResourceData.HasChange("oracle_to_azure_dns_zone_sync_enabled") {
					update.Properties.IsOracleToAzureDnsZoneSyncEnabled = pointer.To(model.OracleToAzureDnsZoneSyncEnabled)
				}
				if err := client.UpdateThenPoll(ctx, *id, *update); err != nil {
					return fmt.Errorf("updating %s: %v", id, err)
				}
			}
			return nil
		},
	}
}

func (NetworkAnchorResource) Read() sdk.ResourceFunc {
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

			state := NetworkAnchorResourceModel{
				Name:              id.NetworkAnchorName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
				state.Zones = pointer.From(model.Zones)
				if props := model.Properties; props != nil {
					state.ResourceAnchorId = props.ResourceAnchorId
					state.SubnetId = props.SubnetId
					state.OciVcnDnsLabel = pointer.From(props.OciVcnDnsLabel)
					state.OciBackupCidrBlock = pointer.From(props.OciBackupCidrBlock)
					state.OracleDnsForwardingEndpointEnabled = pointer.From(props.IsOracleDnsForwardingEndpointEnabled)
					state.OracleDnsListeningEndpointEnabled = pointer.From(props.IsOracleDnsListeningEndpointEnabled)
					state.OracleToAzureDnsZoneSyncEnabled = pointer.From(props.IsOracleToAzureDnsZoneSyncEnabled)
					state.DnsForwardingRule = FlattenDnsForwardingRules(props.DnsForwardingRules)
				}
			}

			return metadata.Encode(&state)
		},
	}
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
	rules := make([]networkanchors.DnsForwardingRule, 0)
	for _, item := range dnsForwardingRules {
		rules = append(rules, networkanchors.DnsForwardingRule{
			DomainNames:         item.DomainNames,
			ForwardingIPAddress: item.ForwardingIPAddress,
		})
	}
	return &rules
}

func (r NetworkAnchorResource) hasUpdates(metadata sdk.ResourceMetaData) bool {
	return metadata.ResourceData.HasChange("tags") ||
		metadata.ResourceData.HasChange("oci_backup_cidr_block") ||
		metadata.ResourceData.HasChange("oracle_dns_forwarding_endpoint_enabled") ||
		metadata.ResourceData.HasChange("oracle_dns_listening_endpoint_enabled") ||
		metadata.ResourceData.HasChange("oracle_to_azure_dns_zone_sync_enabled")
}
