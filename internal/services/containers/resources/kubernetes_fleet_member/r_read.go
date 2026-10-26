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
				r.mapFleetMemberToKubernetesFleetMemberResourceSchema(*model, &schema)
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

func (r Resource) mapKubernetesFleetMemberResourceSchemaToFleetMember(input KubernetesFleetMemberResourceSchema, output *fleetmembers.FleetMember) {
	if output.Properties == nil {
		output.Properties = &fleetmembers.FleetMemberProperties{}
	}

	output.Properties.Group = new(input.Group)
	output.Properties.ClusterResourceId = input.KubernetesClusterId
}

func (r Resource) mapFleetMemberToKubernetesFleetMemberResourceSchema(input fleetmembers.FleetMember, output *KubernetesFleetMemberResourceSchema) {
	if input.Properties == nil {
		input.Properties = &fleetmembers.FleetMemberProperties{}
	}

	output.Group = pointer.From(input.Properties.Group)
	output.KubernetesClusterId = input.Properties.ClusterResourceId
}

func (r Resource) mapKubernetesFleetMemberResourceSchemaToFleetMemberUpdate(input KubernetesFleetMemberResourceSchema, output *fleetmembers.FleetMemberUpdate) {
	if output.Properties == nil {
		output.Properties = &fleetmembers.FleetMemberUpdateProperties{}
	}

	output.Properties.Group = new(input.Group)
}
