package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-10-15/fleetmembers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = KubernetesFleetMemberResource{}
var _ sdk.ResourceWithUpdate = KubernetesFleetMemberResource{}

type KubernetesFleetMemberResource struct{}

func (r KubernetesFleetMemberResource) ModelObject() interface{} {
	return &KubernetesFleetMemberResourceSchema{}
}

type KubernetesFleetMemberResourceSchema struct {
	Name                     string `tfschema:"name"`
	KubernetesFleetManagerId string `tfschema:"kubernetes_fleet_manager_id"`
	KubernetesClusterId      string `tfschema:"kubernetes_cluster_id"`
	Group                    string `tfschema:"group"`
}

func (r KubernetesFleetMemberResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return fleetmembers.ValidateMemberID
}

func (r KubernetesFleetMemberResource) ResourceType() string {
	return "azurerm_kubernetes_fleet_member"
}

func (r KubernetesFleetMemberResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},

		"kubernetes_fleet_manager_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.FleetId{}),

		"kubernetes_cluster_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KubernetesClusterId{}),

		"group": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
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
			client := metadata.Client.Containers.FleetMembersClient

			var config KubernetesFleetMemberResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			fleetId, err := commonids.ParseFleetID(config.KubernetesFleetManagerId)
			if err != nil {
				return err
			}

			id := fleetmembers.NewMemberID(fleetId.SubscriptionId, fleetId.ResourceGroupName, fleetId.FleetName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := fleetmembers.FleetMember{
				Properties: &fleetmembers.FleetMemberProperties{
					ClusterResourceId: config.KubernetesClusterId,
				},
			}

			if config.Group != "" {
				payload.Properties.Group = pointer.To(config.Group)
			}

			if err := client.CreateThenPoll(ctx, id, payload, fleetmembers.DefaultCreateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r KubernetesFleetMemberResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.FleetMembersClient

			id, err := fleetmembers.ParseMemberID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config KubernetesFleetMemberResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload := fleetmembers.FleetMemberUpdate{
				Properties: &fleetmembers.FleetMemberUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("group") {
				payload.Properties.Group = pointer.To(config.Group)
			}

			if err := client.UpdateThenPoll(ctx, *id, payload, fleetmembers.DefaultUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r KubernetesFleetMemberResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.FleetMembersClient
			schema := KubernetesFleetMemberResourceSchema{}

			id, err := fleetmembers.ParseMemberID(metadata.ResourceData.Id())
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
				schema.Name = id.MemberName
				schema.KubernetesFleetManagerId = commonids.NewFleetID(id.SubscriptionId, id.ResourceGroupName, id.FleetName).ID()
				if model.Properties != nil {
					schema.Group = pointer.From(model.Properties.Group)
					clusterId, err := commonids.ParseKubernetesClusterID(model.Properties.ClusterResourceId)
					if err != nil {
						return err
					}
					schema.KubernetesClusterId = clusterId.ID()
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
			client := metadata.Client.Containers.FleetMembersClient

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
