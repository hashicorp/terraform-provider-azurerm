// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_member

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleetmembers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20240401.FleetMembers
			schema := KubernetesFleetMemberResourceSchema{}

			id, err := fleetmembers.ParseMemberID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			kubernetesFleetId := commonids.NewKubernetesFleetID(id.SubscriptionId, id.ResourceGroupName, id.FleetName)

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.KubernetesFleetId = kubernetesFleetId.ID()
				schema.Name = id.MemberName
				if err := r.mapFleetMemberToKubernetesFleetMemberResourceSchema(*model, &schema); err != nil {
					return fmt.Errorf("flattening model: %+v", err)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

type KubernetesFleetMemberResourceSchema struct {
	Group               string `tfschema:"group"`
	KubernetesClusterId string `tfschema:"kubernetes_cluster_id"`
	KubernetesFleetId   string `tfschema:"kubernetes_fleet_id"`
	Name                string `tfschema:"name"`
}

func (r Resource) mapKubernetesFleetMemberResourceSchemaToFleetMember(input KubernetesFleetMemberResourceSchema, output *fleetmembers.FleetMember) error {
	if output.Properties == nil {
		output.Properties = &fleetmembers.FleetMemberProperties{}
	}
	if err := r.mapKubernetesFleetMemberResourceSchemaToFleetMemberProperties(input, output.Properties); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "FleetMemberProperties", "Properties", err)
	}

	return nil
}

func (r Resource) mapFleetMemberToKubernetesFleetMemberResourceSchema(input fleetmembers.FleetMember, output *KubernetesFleetMemberResourceSchema) error {
	if input.Properties == nil {
		input.Properties = &fleetmembers.FleetMemberProperties{}
	}
	if err := r.mapFleetMemberPropertiesToKubernetesFleetMemberResourceSchema(*input.Properties, output); err != nil {
		return fmt.Errorf("mapping SDK Field %q / Model %q to Schema: %+v", "FleetMemberProperties", "Properties", err)
	}

	return nil
}

func (r Resource) mapKubernetesFleetMemberResourceSchemaToFleetMemberProperties(input KubernetesFleetMemberResourceSchema, output *fleetmembers.FleetMemberProperties) error {
	output.Group = &input.Group
	output.ClusterResourceId = input.KubernetesClusterId
	return nil
}

func (r Resource) mapFleetMemberPropertiesToKubernetesFleetMemberResourceSchema(input fleetmembers.FleetMemberProperties, output *KubernetesFleetMemberResourceSchema) error {
	output.Group = pointer.From(input.Group)
	output.KubernetesClusterId = input.ClusterResourceId
	return nil
}

func (r Resource) mapKubernetesFleetMemberResourceSchemaToFleetMemberUpdate(input KubernetesFleetMemberResourceSchema, output *fleetmembers.FleetMemberUpdate) error {
	if output.Properties == nil {
		output.Properties = &fleetmembers.FleetMemberUpdateProperties{}
	}
	if err := r.mapKubernetesFleetMemberResourceSchemaToFleetMemberUpdateProperties(input, output.Properties); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "FleetMemberUpdateProperties", "Properties", err)
	}

	return nil
}

func (r Resource) mapFleetMemberUpdateToKubernetesFleetMemberResourceSchema(input fleetmembers.FleetMemberUpdate, output *KubernetesFleetMemberResourceSchema) error {
	if input.Properties == nil {
		input.Properties = &fleetmembers.FleetMemberUpdateProperties{}
	}
	if err := r.mapFleetMemberUpdatePropertiesToKubernetesFleetMemberResourceSchema(*input.Properties, output); err != nil {
		return fmt.Errorf("mapping SDK Field %q / Model %q to Schema: %+v", "FleetMemberUpdateProperties", "Properties", err)
	}

	return nil
}

func (r Resource) mapKubernetesFleetMemberResourceSchemaToFleetMemberUpdateProperties(input KubernetesFleetMemberResourceSchema, output *fleetmembers.FleetMemberUpdateProperties) error {
	output.Group = &input.Group
	return nil
}

func (r Resource) mapFleetMemberUpdatePropertiesToKubernetesFleetMemberResourceSchema(input fleetmembers.FleetMemberUpdateProperties, output *KubernetesFleetMemberResourceSchema) error {
	output.Group = pointer.From(input.Group)
	return nil
}
