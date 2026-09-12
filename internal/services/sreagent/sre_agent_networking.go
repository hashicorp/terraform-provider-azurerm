// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sreagent

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-01-01/agents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SreAgentNetworking struct {
	EgressMode string               `tfschema:"egress_mode"`
	SubnetID   string               `tfschema:"subnet_id"`
	PrivateDNS []SreAgentPrivateDNS `tfschema:"private_dns"`
}

type SreAgentPrivateDNS struct {
	Enabled bool `tfschema:"enabled"`
}

func sreAgentNetworkingSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type: pluginsdk.TypeList, Optional: true, Computed: true, MaxItems: 1,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"egress_mode": {
				Type: pluginsdk.TypeString, Required: true,
				ValidateFunc: validation.StringInSlice(agents.PossibleValuesForSandboxEgressMode(), false),
			},
			"subnet_id": {
				Type: pluginsdk.TypeString, Optional: true,
				ValidateFunc: commonids.ValidateSubnetID,
			},
			"private_dns": {
				Type: pluginsdk.TypeList, Optional: true, Computed: true, MaxItems: 1,
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"enabled": {Type: pluginsdk.TypeBool, Required: true},
				}},
			},
		}},
	}
}

func validateSreAgentNetworking(network SreAgentNetworking) error {
	switch agents.SandboxEgressMode(network.EgressMode) {
	case agents.SandboxEgressModeAzureVNet:
		if network.SubnetID == "" {
			return fmt.Errorf("`networking.subnet_id` is required when `egress_mode` is AzureVNet")
		}
		if _, err := commonids.ParseSubnetID(network.SubnetID); err != nil {
			return fmt.Errorf("parsing `networking.subnet_id`: %+v", err)
		}
	case agents.SandboxEgressModeLimited, agents.SandboxEgressModeUnrestricted:
		if network.SubnetID != "" {
			return fmt.Errorf("`networking.subnet_id` must be omitted when `egress_mode` is %s", network.EgressMode)
		}
	default:
		return fmt.Errorf("`networking.egress_mode` must explicitly select AzureVNet, Limited or Unrestricted")
	}
	return nil
}

func sreAgentNetworkingRemovalError() error {
	return fmt.Errorf("removing the entire `networking` block does not detach the subnet: explicitly select a non-VNet `egress_mode` and omit `subnet_id`")
}

func sreAgentPrivateDNSRemovalError() error {
	return fmt.Errorf("removing `networking.private_dns` does not reset the service setting: explicitly configure `enabled` as true or false; null reset is not supported")
}

func expandSreAgentNetworking(network SreAgentNetworking, props *agents.AgentPatchProperties, attachmentChanged, dnsChanged bool) error {
	if err := validateSreAgentNetworking(network); err != nil {
		return err
	}
	props.SandboxConfiguration = &agents.SandboxConfiguration{
		Egress: &agents.SandboxEgressConfiguration{},
	}
	if attachmentChanged {
		props.SandboxConfiguration.Egress.Mode = pointer.To(agents.SandboxEgressMode(network.EgressMode))
		if network.SubnetID != "" {
			props.VnetConfiguration = &agents.VnetConfiguration{SubnetResourceId: pointer.To(network.SubnetID)}
		}
	}
	if dnsChanged {
		if len(network.PrivateDNS) != 1 {
			return sreAgentPrivateDNSRemovalError()
		}
		props.SandboxConfiguration.Egress.VnetConfiguration = &agents.SandboxVnetConfiguration{
			UsePrivateDnsResolution: pointer.To(network.PrivateDNS[0].Enabled),
		}
	}
	return nil
}

func flattenSreAgentNetworking(props *agents.AgentProperties) []SreAgentNetworking {
	network := SreAgentNetworking{}
	if props.VnetConfiguration != nil {
		network.SubnetID = pointer.From(props.VnetConfiguration.SubnetResourceId)
	}
	if props.SandboxConfiguration != nil && props.SandboxConfiguration.Egress != nil {
		network.EgressMode = string(pointer.From(props.SandboxConfiguration.Egress.Mode))
		if dns := props.SandboxConfiguration.Egress.VnetConfiguration; dns != nil && dns.UsePrivateDnsResolution != nil {
			network.PrivateDNS = []SreAgentPrivateDNS{{Enabled: *dns.UsePrivateDnsResolution}}
		}
	}
	if network.SubnetID == "" && network.EgressMode == "" && len(network.PrivateDNS) == 0 {
		return nil
	}
	return []SreAgentNetworking{network}
}

func sreAgentNeedsNetworkDetach(props *agents.AgentPatchProperties) bool {
	if props == nil || props.VnetConfiguration != nil || props.SandboxConfiguration == nil || props.SandboxConfiguration.Egress == nil {
		return false
	}
	mode := pointer.From(props.SandboxConfiguration.Egress.Mode)
	return mode == agents.SandboxEgressModeLimited || mode == agents.SandboxEgressModeUnrestricted
}

func marshalSreAgentPatchProperties(props *agents.AgentPatchProperties) (json.RawMessage, error) {
	if props == nil {
		return nil, nil
	}
	if !sreAgentNeedsNetworkDetach(props) {
		return json.Marshal(props)
	}
	// Generated optional pointers omit nil; the public detach contract requires an explicit null.
	return json.Marshal(struct {
		agents.AgentPatchProperties
		VnetConfiguration json.RawMessage `json:"vnetConfiguration"`
	}{AgentPatchProperties: *props})
}
