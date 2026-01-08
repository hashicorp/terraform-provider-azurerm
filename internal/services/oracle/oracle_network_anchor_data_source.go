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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/networkanchors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type NetworkAnchorDataSource struct{}

type NetworkAnchorDataModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	Location          string `tfschema:"location"`

	CidrBlock                          string            `tfschema:"cidr_block"`
	DnsForwardingEndpointIpAddress     string            `tfschema:"dns_forwarding_endpoint_ip_address"`
	DnsForwardingEndpointNsgRuleUrl    string            `tfschema:"dns_forwarding_endpoint_nsg_rule_url"`
	DnsForwardingRuleUrl               string            `tfschema:"dns_forwarding_rule_url"`
	DnsListeningEndpointAllowedCidrs   string            `tfschema:"dns_listening_endpoint_allowed_cidrs"`
	DnsListeningEndpointIpAddress      string            `tfschema:"dns_listening_endpoint_ip_address"`
	DnsListeningEndpointNsgRuleUrl     string            `tfschema:"dns_listening_endpoint_nsg_rule_url"`
	OciBackupCidrBlock                 string            `tfschema:"oci_backup_cidr_block"`
	OciSubnetId                        string            `tfschema:"oci_subnet_id"`
	OciVcnDnsLabel                     string            `tfschema:"oci_vcn_dns_label"`
	OciVcnId                           string            `tfschema:"oci_vcn_id"`
	OracleDnsForwardingEndpointEnabled bool              `tfschema:"oracle_dns_forwarding_endpoint_enabled"`
	OracleDnsListeningEndpointEnabled  bool              `tfschema:"oracle_dns_listening_endpoint_enabled"`
	OracleToAzureDnsZoneSyncEnabled    bool              `tfschema:"oracle_to_azure_dns_zone_sync_enabled"`
	ProvisioningState                  string            `tfschema:"provisioning_state"`
	ResourceAnchorId                   string            `tfschema:"resource_anchor_id"`
	SubnetId                           string            `tfschema:"subnet_id"`
	VnetId                             string            `tfschema:"vnet_id"`
	Tags                               map[string]string `tfschema:"tags"`
	Zones                              zones.Schema      `tfschema:"zones"`
}

func (d NetworkAnchorDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 24),
				validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-]+$`), "Name may include letters, numbers, or hyphens only"),
			),
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d NetworkAnchorDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"resource_anchor_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"vnet_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"subnet_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"cidr_block": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"oci_vcn_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"oci_vcn_dns_label": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"oci_subnet_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"oci_backup_cidr_block": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"oracle_dns_forwarding_endpoint_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"oracle_dns_listening_endpoint_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"oracle_to_azure_dns_zone_sync_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"dns_forwarding_rule_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"dns_forwarding_endpoint_ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"dns_listening_endpoint_allowed_cidrs": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"dns_listening_endpoint_nsg_rule_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"dns_listening_endpoint_ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"dns_forwarding_endpoint_nsg_rule_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"zones": commonschema.ZonesMultipleComputed(),
		"tags":  commonschema.TagsDataSource(),
	}
}

func (d NetworkAnchorDataSource) ModelObject() interface{} {
	return &NetworkAnchorDataSource{}
}

func (d NetworkAnchorDataSource) ResourceType() string {
	return "azurerm_oracle_network_anchor"
}

func (d NetworkAnchorDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return networkanchors.ValidateNetworkAnchorID
}

func (d NetworkAnchorDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.NetworkAnchors
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state NetworkAnchorDataModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := networkanchors.NewNetworkAnchorID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return err
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
				state.Zones = pointer.From(model.Zones)
				if props := model.Properties; props != nil {
					state.ResourceAnchorId = props.ResourceAnchorId
					state.ProvisioningState = pointer.FromEnum(props.ProvisioningState)
					state.VnetId = pointer.From(props.VnetId)
					state.SubnetId = props.SubnetId
					state.CidrBlock = pointer.From(props.CidrBlock)
					state.OciVcnId = pointer.From(props.OciVcnId)
					state.OciVcnDnsLabel = pointer.From(props.OciVcnDnsLabel)
					state.OciSubnetId = pointer.From(props.OciSubnetId)
					state.OciBackupCidrBlock = pointer.From(props.OciBackupCidrBlock)
					state.OracleDnsForwardingEndpointEnabled = pointer.From(props.IsOracleDnsForwardingEndpointEnabled)
					state.OracleDnsListeningEndpointEnabled = pointer.From(props.IsOracleDnsListeningEndpointEnabled)
					state.OracleToAzureDnsZoneSyncEnabled = pointer.From(props.IsOracleToAzureDnsZoneSyncEnabled)
					state.DnsForwardingEndpointIpAddress = pointer.From(props.DnsForwardingEndpointIPAddress)
					state.DnsForwardingRuleUrl = pointer.From(props.DnsForwardingRulesURL)
					state.DnsForwardingEndpointNsgRuleUrl = pointer.From(props.DnsForwardingEndpointNsgRulesURL)
					state.DnsListeningEndpointAllowedCidrs = pointer.From(props.DnsListeningEndpointAllowedCidrs)
					state.DnsListeningEndpointIpAddress = pointer.From(props.DnsListeningEndpointIPAddress)
					state.DnsListeningEndpointNsgRuleUrl = pointer.From(props.DnsListeningEndpointNsgRulesURL)
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
