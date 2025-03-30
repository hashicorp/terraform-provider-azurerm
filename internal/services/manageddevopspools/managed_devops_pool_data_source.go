package manageddevopspools

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2024-10-19/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = ManagedDevOpsPoolDataSource{}

type ManagedDevOpsPoolDataSource struct{}

func (ManagedDevOpsPoolDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

func (ManagedDevOpsPoolDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"location": commonschema.LocationComputed(),
		"dev_center_project_resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"maximum_concurrency": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},
		"agent_profile": AgentProfileComputedSchema(),
		"fabric_profile": FabricProfileComputedSchema(),
		"identity":       commonschema.SystemAssignedUserAssignedIdentityComputed(),
		"organization_profile": OrganizationProfileComputedSchema(),
		"tags": commonschema.TagsDataSource(),
	}
}

func (ManagedDevOpsPoolDataSource) ModelObject() interface{} {
	return &ManagedDevOpsPoolModel{}
}

func (ManagedDevOpsPoolDataSource) ResourceType() string {
	return "azurerm_managed_devops_pool"
}

func (ManagedDevOpsPoolDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedDevOpsPools.PoolsClient

			subscriptionId := metadata.Client.Account.SubscriptionId

			var state ManagedDevOpsPoolModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := pools.NewPoolID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				state.Type = pointer.From(model.Type)
				state.Location = model.Location
				state.Tags = pointer.From(model.Tags)

				if modelIdentity := model.Identity; modelIdentity != nil {
					identity, err := identity.FlattenSystemAndUserAssignedMapToModel(pointer.To((identity.SystemAndUserAssignedMap)(*model.Identity)))
					if err != nil {
						return err
					}

					state.Identity = pointer.From(identity)
				}

				if props := model.Properties; props != nil {
					state.DevCenterProjectResourceId = props.DevCenterProjectResourceId
					state.MaximumConcurrency = props.MaximumConcurrency

					if agentProfile := props.AgentProfile; agentProfile != nil {
						state.AgentProfile = flattenAgentProfileToModel(agentProfile)
					}

					if organizationProfile := props.OrganizationProfile; organizationProfile != nil {
						state.OrganizationProfile = flattenOrganizationProfileToModel(organizationProfile)
					}

					if fabricProfile := props.FabricProfile; fabricProfile != nil {
						state.FabricProfile = flattenFabricProfileToModel(fabricProfile)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}
