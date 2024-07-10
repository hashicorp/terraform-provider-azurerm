package containers

// NOTE: this file is generated - manual changes will be overwritten.
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.
import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleetmembers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = KubernetesFleetMemberResource{}
var _ sdk.ResourceWithUpdate = KubernetesFleetMemberResource{}

type KubernetesFleetMemberResource struct{}

func (r KubernetesFleetMemberResource) ModelObject() interface{} {
	return &KubernetesFleetMemberResourceSchema{}
}

type KubernetesFleetMemberResourceSchema struct {
	Group               string `tfschema:"group"`
	KubernetesClusterId string `tfschema:"kubernetes_cluster_id"`
	KubernetesFleetId   string `tfschema:"kubernetes_fleet_id"`
	Name                string `tfschema:"name"`
}

func (r KubernetesFleetMemberResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return fleetmembers.ValidateMemberID
}
func (r KubernetesFleetMemberResource) ResourceType() string {
	return "azurerm_kubernetes_fleet_member"
}
func (r KubernetesFleetMemberResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"kubernetes_cluster_id": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"kubernetes_fleet_id": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"group": {
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
	}
}
func (r KubernetesFleetMemberResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
func (r KubernetesFleetMemberResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20231015.FleetMembers

			var config KubernetesFleetMemberResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			kubernetesFleetId, err := commonids.ParseKubernetesFleetID(config.KubernetesFleetId)
			if err != nil {
				return err
			}

			id := fleetmembers.NewMemberID(subscriptionId, kubernetesFleetId.ResourceGroupName, kubernetesFleetId.FleetName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			var payload fleetmembers.FleetMember
			if err := r.mapKubernetesFleetMemberResourceSchemaToFleetMember(config, &payload); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

			if err := client.CreateThenPoll(ctx, id, payload, fleetmembers.DefaultCreateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
func (r KubernetesFleetMemberResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20231015.FleetMembers
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
func (r KubernetesFleetMemberResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20231015.FleetMembers

			id, err := fleetmembers.ParseMemberID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, fleetmembers.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
func (r KubernetesFleetMemberResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20231015.FleetMembers

			id, err := fleetmembers.ParseMemberID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config KubernetesFleetMemberResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			var payload fleetmembers.FleetMemberUpdate
			if err := r.mapKubernetesFleetMemberResourceSchemaToFleetMemberUpdate(config, &payload); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

			if err := client.UpdateThenPoll(ctx, *id, payload, fleetmembers.DefaultUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r KubernetesFleetMemberResource) mapKubernetesFleetMemberResourceSchemaToFleetMember(input KubernetesFleetMemberResourceSchema, output *fleetmembers.FleetMember) error {

	if output.Properties == nil {
		output.Properties = &fleetmembers.FleetMemberProperties{}
	}
	if err := r.mapKubernetesFleetMemberResourceSchemaToFleetMemberProperties(input, output.Properties); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "FleetMemberProperties", "Properties", err)
	}

	return nil
}

func (r KubernetesFleetMemberResource) mapFleetMemberToKubernetesFleetMemberResourceSchema(input fleetmembers.FleetMember, output *KubernetesFleetMemberResourceSchema) error {

	if input.Properties == nil {
		input.Properties = &fleetmembers.FleetMemberProperties{}
	}
	if err := r.mapFleetMemberPropertiesToKubernetesFleetMemberResourceSchema(*input.Properties, output); err != nil {
		return fmt.Errorf("mapping SDK Field %q / Model %q to Schema: %+v", "FleetMemberProperties", "Properties", err)
	}

	return nil
}

func (r KubernetesFleetMemberResource) mapKubernetesFleetMemberResourceSchemaToFleetMemberProperties(input KubernetesFleetMemberResourceSchema, output *fleetmembers.FleetMemberProperties) error {
	output.Group = &input.Group
	output.ClusterResourceId = input.KubernetesClusterId
	return nil
}

func (r KubernetesFleetMemberResource) mapFleetMemberPropertiesToKubernetesFleetMemberResourceSchema(input fleetmembers.FleetMemberProperties, output *KubernetesFleetMemberResourceSchema) error {
	output.Group = pointer.From(input.Group)
	output.KubernetesClusterId = input.ClusterResourceId
	return nil
}

func (r KubernetesFleetMemberResource) mapKubernetesFleetMemberResourceSchemaToFleetMemberUpdate(input KubernetesFleetMemberResourceSchema, output *fleetmembers.FleetMemberUpdate) error {

	if output.Properties == nil {
		output.Properties = &fleetmembers.FleetMemberUpdateProperties{}
	}
	if err := r.mapKubernetesFleetMemberResourceSchemaToFleetMemberUpdateProperties(input, output.Properties); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "FleetMemberUpdateProperties", "Properties", err)
	}

	return nil
}

func (r KubernetesFleetMemberResource) mapFleetMemberUpdateToKubernetesFleetMemberResourceSchema(input fleetmembers.FleetMemberUpdate, output *KubernetesFleetMemberResourceSchema) error {

	if input.Properties == nil {
		input.Properties = &fleetmembers.FleetMemberUpdateProperties{}
	}
	if err := r.mapFleetMemberUpdatePropertiesToKubernetesFleetMemberResourceSchema(*input.Properties, output); err != nil {
		return fmt.Errorf("mapping SDK Field %q / Model %q to Schema: %+v", "FleetMemberUpdateProperties", "Properties", err)
	}

	return nil
}

func (r KubernetesFleetMemberResource) mapKubernetesFleetMemberResourceSchemaToFleetMemberUpdateProperties(input KubernetesFleetMemberResourceSchema, output *fleetmembers.FleetMemberUpdateProperties) error {
	output.Group = &input.Group
	return nil
}

func (r KubernetesFleetMemberResource) mapFleetMemberUpdatePropertiesToKubernetesFleetMemberResourceSchema(input fleetmembers.FleetMemberUpdateProperties, output *KubernetesFleetMemberResourceSchema) error {
	output.Group = pointer.From(input.Group)
	return nil
}
