// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_manager

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"hub_profile": {
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"agent_profile": {
						Optional: true,
						// NOTE: O+C - populated by the API when the service creates default hub agent settings.
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"subnet_id": {
									Optional: true,
									// NOTE: O+C - populated by the API on read/import when the service returns the configured subnet.
									Computed:     true,
									ForceNew:     true,
									Type:         pluginsdk.TypeString,
									ValidateFunc: commonids.ValidateSubnetID,
									AtLeastOneOf: []string{"hub_profile.0.agent_profile.0.subnet_id", "hub_profile.0.agent_profile.0.virtual_machine_size"},
								},
								"virtual_machine_size": {
									Optional: true,
									// NOTE: O+C - populated by the API on read/import when the service returns the configured VM size.
									Computed:     true,
									ForceNew:     true,
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.StringIsNotEmpty,
									AtLeastOneOf: []string{"hub_profile.0.agent_profile.0.subnet_id", "hub_profile.0.agent_profile.0.virtual_machine_size"},
								},
							},
						},
						ForceNew: true,
						MaxItems: 1,
						Type:     pluginsdk.TypeList,
					},
					"api_server_access_profile": {
						Optional: true,
						// NOTE: O+C - populated by the API when the service returns the effective hub API server access settings.
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"private_cluster_enabled": {
									Required: true,
									ForceNew: true,
									Type:     pluginsdk.TypeBool,
								},
							},
						},
						ForceNew: true,
						MaxItems: 1,
						Type:     pluginsdk.TypeList,
					},
					"dns_prefix": {
						Optional: true,
						// NOTE: O+C - preserves the DNS prefix returned by Azure during import.
						Computed: true,
						ForceNew: true,
						Type:     pluginsdk.TypeString,
						ValidateFunc: validation.All(
							validation.StringLenBetween(1, 54),
							validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]$|^[a-zA-Z0-9][a-zA-Z0-9-]{0,52}[a-zA-Z0-9]$`), "must match the pattern ^[a-zA-Z0-9]$|^[a-zA-Z0-9][a-zA-Z0-9-]{0,52}[a-zA-Z0-9]$"),
						),
					},
					"fqdn": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
					"kubernetes_version": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
					"portal_fqdn": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
				},
			},
			ForceNew: true,
			MaxItems: 1,
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rawConfig := metadata.ResourceDiff.GetRawConfig()
			if !rawConfig.IsKnown() || rawConfig.IsNull() {
				return nil
			}
			hub := rawConfig.AsValueMap()["hub_profile"]
			if !hub.IsKnown() || hub.IsNull() || len(hub.AsValueSlice()) == 0 {
				return nil
			}
			hubProfile := hub.AsValueSlice()[0]
			if !hubProfile.IsKnown() || hubProfile.IsNull() {
				return nil
			}
			fields := hubProfile.AsValueMap()
			access := fields["api_server_access_profile"]
			if !access.IsKnown() || access.IsNull() || len(access.AsValueSlice()) == 0 {
				return nil
			}
			accessProfile := access.AsValueSlice()[0]
			if !accessProfile.IsKnown() || accessProfile.IsNull() {
				return nil
			}
			private := accessProfile.AsValueMap()["private_cluster_enabled"]
			if !private.IsKnown() || private.IsNull() || !private.True() {
				return nil
			}
			if dns := fields["dns_prefix"]; dns.IsKnown() && !dns.IsNull() && dns.AsString() != "" {
				return errors.New("`hub_profile.dns_prefix` cannot be set when `api_server_access_profile.private_cluster_enabled` is true")
			}
			if metadata.ResourceDiff.Id() != "" && !metadata.ResourceDiff.HasChange("hub_profile") {
				return nil
			}
			agent := fields["agent_profile"]
			if !agent.IsKnown() {
				return nil
			}
			if agent.IsNull() || len(agent.AsValueSlice()) == 0 {
				return errors.New("`hub_profile.agent_profile.subnet_id` is required when creating a private Fleet hub")
			}
			agentProfile := agent.AsValueSlice()[0]
			if !agentProfile.IsKnown() || agentProfile.IsNull() {
				return nil
			}
			subnet := agentProfile.AsValueMap()["subnet_id"]
			if subnet.IsKnown() && (subnet.IsNull() || subnet.AsString() == "") {
				return errors.New("`hub_profile.agent_profile.subnet_id` is required when creating a private Fleet hub")
			}
			return nil
		},
	}
}
