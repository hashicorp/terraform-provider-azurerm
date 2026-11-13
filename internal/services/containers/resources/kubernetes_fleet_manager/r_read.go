// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_manager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20240401.Fleets
			schema := KubernetesFleetManagerResourceSchema{}

			id, err := fleets.ParseFleetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.Name = id.FleetName
				schema.ResourceGroupName = id.ResourceGroupName
				if err := r.mapFleetToKubernetesFleetManagerResourceSchema(*model, &schema); err != nil {
					return err
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

type KubernetesFleetManagerResourceSchema struct {
	Location          string                   `tfschema:"location"`
	Name              string                   `tfschema:"name"`
	ResourceGroupName string                   `tfschema:"resource_group_name"`
	HubProfile        []FleetManagerHubProfile `tfschema:"hub_profile"`
	Tags              map[string]any           `tfschema:"tags"`
}

type FleetManagerHubProfile struct {
	AgentProfile           []FleetManagerHubAgentProfile           `tfschema:"agent_profile"`
	ApiServerAccessProfile []FleetManagerHubAPIServerAccessProfile `tfschema:"api_server_access_profile"`
	DnsPrefix              string                                  `tfschema:"dns_prefix"`
	Fqdn                   string                                  `tfschema:"fqdn"`
	KubernetesVersion      string                                  `tfschema:"kubernetes_version"`
	PortalFqdn             string                                  `tfschema:"portal_fqdn"`
}

type FleetManagerHubAgentProfile struct {
	SubnetId           string `tfschema:"subnet_id"`
	VirtualMachineSize string `tfschema:"virtual_machine_size"`
}

type FleetManagerHubAPIServerAccessProfile struct {
	PrivateClusterEnabled bool `tfschema:"private_cluster_enabled"`
}

func (r Resource) mapKubernetesFleetManagerResourceSchemaToFleet(input KubernetesFleetManagerResourceSchema, output *fleets.Fleet) {
	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Expand(input.Tags)

	if output.Properties == nil {
		output.Properties = &fleets.FleetProperties{}
	}

	if len(input.HubProfile) > 0 {
		output.Properties.HubProfile = expandFleetManagerHubProfile(input.HubProfile)
	}
}

func (r Resource) mapFleetToKubernetesFleetManagerResourceSchema(input fleets.Fleet, output *KubernetesFleetManagerResourceSchema) error {
	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Flatten(input.Tags)

	if input.Properties == nil {
		input.Properties = &fleets.FleetProperties{}
	}

	hubProfile, err := flattenFleetManagerHubProfile(input.Properties.HubProfile)
	if err != nil {
		return err
	}
	output.HubProfile = hubProfile
	return nil
}

func expandFleetManagerHubProfile(input []FleetManagerHubProfile) *fleets.FleetHubProfile {
	if len(input) == 0 {
		return nil
	}

	output := &fleets.FleetHubProfile{
		AgentProfile:           expandFleetManagerHubAgentProfile(input[0].AgentProfile),
		ApiServerAccessProfile: expandFleetManagerHubAPIServerAccessProfile(input[0].ApiServerAccessProfile),
	}

	if input[0].DnsPrefix != "" {
		output.DnsPrefix = pointer.To(input[0].DnsPrefix)
	}

	return output
}

func flattenFleetManagerHubProfile(input *fleets.FleetHubProfile) ([]FleetManagerHubProfile, error) {
	if input == nil {
		return []FleetManagerHubProfile{}, nil
	}
	agentProfile, err := flattenFleetManagerHubAgentProfile(input.AgentProfile)
	if err != nil {
		return []FleetManagerHubProfile{}, err
	}

	return []FleetManagerHubProfile{
		{
			AgentProfile:           agentProfile,
			ApiServerAccessProfile: flattenFleetManagerHubAPIServerAccessProfile(input.ApiServerAccessProfile),
			DnsPrefix:              pointer.From(input.DnsPrefix),
			Fqdn:                   pointer.From(input.Fqdn),
			KubernetesVersion:      pointer.From(input.KubernetesVersion),
			PortalFqdn:             pointer.From(input.PortalFqdn),
		},
	}, nil
}

func expandFleetManagerHubAgentProfile(input []FleetManagerHubAgentProfile) *fleets.AgentProfile {
	if len(input) == 0 {
		return nil
	}

	output := &fleets.AgentProfile{}
	if input[0].SubnetId != "" {
		output.SubnetId = pointer.To(input[0].SubnetId)
	}
	if input[0].VirtualMachineSize != "" {
		output.VMSize = pointer.To(input[0].VirtualMachineSize)
	}

	return output
}

func flattenFleetManagerHubAgentProfile(input *fleets.AgentProfile) ([]FleetManagerHubAgentProfile, error) {
	if input == nil || (input.SubnetId == nil && input.VMSize == nil) {
		return []FleetManagerHubAgentProfile{}, nil
	}

	subnetId := ""
	if input.SubnetId != nil && *input.SubnetId != "" {
		id, err := commonids.ParseSubnetIDInsensitively(*input.SubnetId)
		if err != nil {
			return []FleetManagerHubAgentProfile{}, err
		}
		subnetId = id.ID()
	}

	return []FleetManagerHubAgentProfile{
		{
			SubnetId:           subnetId,
			VirtualMachineSize: pointer.From(input.VMSize),
		},
	}, nil
}

func expandFleetManagerHubAPIServerAccessProfile(input []FleetManagerHubAPIServerAccessProfile) *fleets.APIServerAccessProfile {
	if len(input) == 0 {
		return nil
	}

	return &fleets.APIServerAccessProfile{
		EnablePrivateCluster: pointer.To(input[0].PrivateClusterEnabled),
	}
}

func flattenFleetManagerHubAPIServerAccessProfile(input *fleets.APIServerAccessProfile) []FleetManagerHubAPIServerAccessProfile {
	if input == nil || input.EnablePrivateCluster == nil {
		return []FleetManagerHubAPIServerAccessProfile{}
	}

	return []FleetManagerHubAPIServerAccessProfile{
		{
			PrivateClusterEnabled: pointer.From(input.EnablePrivateCluster),
		},
	}
}
