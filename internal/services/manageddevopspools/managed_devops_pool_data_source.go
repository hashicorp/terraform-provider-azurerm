package manageddevopspools

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2025-01-21/pools"
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
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (ManagedDevOpsPoolDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"agent_profile": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"kind": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"grace_period_time_span": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"max_agent_lifetime": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"resource_predictions":         ResourcePredictionsSchema(),
					"resource_predictions_profile": ResourcePredictionsProfileSchema(),
				},
			},
		},
		"dev_center_project_resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"fabric_profile": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"image": ImageSchema(),
					"kind": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string("Vmss"),
						}, false),
					},
					"network_profile": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"subnet_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
					"os_profile": OsProfileSchema(),
					"sku": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
					"storage_profile": StorageProfileSchema(),
				},
			},
		},
		"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),
		"location": commonschema.LocationComputed(),
		"maximum_concurrency": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},
		"organization_profile": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"kind": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string("AzureDevOps"),
						}, false),
					},
					"organization": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"parallelism": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
								"projects": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"url": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsURLWithHTTPS,
								},
							},
						},
					},
					"permission_profile": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"groups": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"kind": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"users": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},
		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
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
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				expandedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return err
				}

				state.Identity = expandedIdentity

				if props := model.Properties; props != nil {
					state.DevCenterProjectResourceId = props.DevCenterProjectResourceId
					state.MaximumConcurrency = props.MaximumConcurrency
					state.ProvisioningState = string(pointer.From(props.ProvisioningState))

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
