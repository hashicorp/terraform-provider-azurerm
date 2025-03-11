package manageddevopspools

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2024-10-19/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = ManagedDevOpsPoolDataSource{}

type ManagedDevOpsPoolDataSource struct{}

func (ManagedDevOpsPoolDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"agent_profile": AgentProfileSchema(),
		"dev_center_project_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"fabric_profile": FabricProfileSchema(),
		"identity":       commonschema.SystemAssignedUserAssignedIdentityOptional(),
		"location":       commonschema.Location(),
		"maximum_concurrency": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"organization_profile": OrganizationProfileSchema(),
		"resource_group_name":  commonschema.ResourceGroupNameForDataSource(),
		"tags":                 commonschema.Tags(),
		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (ManagedDevOpsPoolDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

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

			id, err := pools.NewPoolID(subscriptionId, state.ResourceGroupName, state.Name)
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			state := ManagedDevOpsPoolModel{
				Name: id.PoolName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Location = location.NormalizeNilable(model.Location)
					state.Tags = pointer.From(model.Tags)
					state.Type = model.Type
					state.ResourceGroupName = pools.ResourceGroupNameFromID(model.ID)
					// if there are properties to set into state do that here
					state.DevCenterProjectResourceId = props.DevCenterProjectResourceId
					state.MaximumConcurrency = props.MaximumConcurrency
					// Add other props
				}
			}

			return metadata.Encode(&state)
		},
	}
}
