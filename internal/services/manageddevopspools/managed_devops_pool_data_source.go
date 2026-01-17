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

type ManagedDevOpsPoolDataSourceModel struct {
	DevCenterProjectResourceId     string                                `tfschema:"dev_center_project_resource_id"`
	VmssFabricProfile              []VmssFabricProfileModel              `tfschema:"vmss_fabric_profile"`
	Identity                       []identity.ModelUserAssigned          `tfschema:"identity"`
	Location                       string                                `tfschema:"location"`
	MaximumConcurrency             int64                                 `tfschema:"maximum_concurrency"`
	Name                           string                                `tfschema:"name"`
	AzureDevOpsOrganizationProfile []AzureDevOpsOrganizationProfileModel `tfschema:"azure_devops_organization_profile"`
	ResourceGroupName              string                                `tfschema:"resource_group_name"`
	Tags                           map[string]string                     `tfschema:"tags"`
	StatefulAgentProfile           []StatefulAgentProfileModel           `tfschema:"stateful_agent_profile"`
	StatelessAgentProfile          []StatelessAgentProfileModel          `tfschema:"stateless_agent_profile"`
}

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
		"location": commonschema.LocationComputed(),

		"azure_devops_organization_profile": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"organization": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"parallelism": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},

								"projects": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"url": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"permission_profile": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"kind": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"administrator_account": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"groups": {
												Type:     pluginsdk.TypeList,
												Computed: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},

											"users": {
												Type:     pluginsdk.TypeList,
												Computed: true,
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
				},
			},
		},

		"dev_center_project_resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"identity": commonschema.UserAssignedIdentityComputed(),

		"maximum_concurrency": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"stateful_agent_profile": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"grace_period_time_span": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"max_agent_lifetime": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"manual_resource_predictions_profile": manualResourcePredictionsProfileSchemaComputed(),

					"automatic_resource_predictions_profile": automaticResourcePredictionsProfileSchemaComputed(),
				},
			},
		},

		"stateless_agent_profile": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"manual_resource_predictions_profile": manualResourcePredictionsProfileSchemaComputed(),

					"automatic_resource_predictions_profile": automaticResourcePredictionsProfileSchemaComputed(),
				},
			},
		},

		"tags": commonschema.TagsDataSource(),

		"vmss_fabric_profile": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"image": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"aliases": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"buffer": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"resource_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"well_known_image_name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"os_profile": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"logon_type": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"secrets_management": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"certificate_store_location": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},

											"certificate_store_name": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},

											"key_export_enabled": {
												Type:     pluginsdk.TypeBool,
												Computed: true,
											},

											"observed_certificates": {
												Type:     pluginsdk.TypeList,
												Computed: true,
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

					"sku_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"storage_profile": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"data_disk": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"caching": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},

											"disk_size_gb": {
												Type:     pluginsdk.TypeInt,
												Computed: true,
											},

											"drive_letter": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},

											"storage_account_type": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
										},
									},
								},

								"os_disk_storage_account_type": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (ManagedDevOpsPoolDataSource) ModelObject() interface{} {
	return &ManagedDevOpsPoolDataSourceModel{}
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

			var state ManagedDevOpsPoolDataSourceModel
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

				if model.Identity != nil {
					flattenedIdentity, err := flattenManagedDevopsUserAssignedToLegacyIdentity(model.Identity)
					if err != nil {
						return fmt.Errorf("flattening `identity`: %+v", err)
					}
					state.Identity = flattenedIdentity
				}

				if props := model.Properties; props != nil {
					state.DevCenterProjectResourceId = props.DevCenterProjectResourceId
					state.MaximumConcurrency = props.MaximumConcurrency

					if agentProfile := props.AgentProfile; agentProfile != nil {
						if stateful, ok := agentProfile.(pools.Stateful); ok {
							state.StatefulAgentProfile = flattenStatefulAgentProfileToModel(stateful)
						} else if stateless, ok := agentProfile.(pools.StatelessAgentProfile); ok {
							state.StatelessAgentProfile = flattenStatelessAgentProfileToModel(stateless)
						}
					}

					if organizationProfile := props.OrganizationProfile; organizationProfile != nil {
						if azureDevOpsOrganizationProfile, ok := organizationProfile.(pools.AzureDevOpsOrganizationProfile); ok {
							state.AzureDevOpsOrganizationProfile = flattenAzureDevOpsOrganizationProfileToModel(azureDevOpsOrganizationProfile)
						}
					}

					if fabricProfile := props.FabricProfile; fabricProfile != nil {
						if vmssFabricProfile, ok := fabricProfile.(pools.VMSSFabricProfile); ok {
							state.VmssFabricProfile = flattenVmssFabricProfileToModel(vmssFabricProfile)
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}
