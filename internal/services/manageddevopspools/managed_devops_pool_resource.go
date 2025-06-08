package manageddevopspools

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projects"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2025-01-21/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = ManagedDevOpsPoolResource{}

type ManagedDevOpsPoolResource struct{}

func (ManagedDevOpsPoolResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-.]*[a-zA-Z0-9-]$`),
				"`name` can only include alphanumeric characters, periods (.) and hyphens (-). It must also start with alphanumeric characters and cannot end with periods (.).",
			),
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"location":            commonschema.Location(),
		"agent_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"kind": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string("Stateless"),
							string("Stateful"),
						}, false),
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
		"dev_center_project_resource_id": commonschema.ResourceIDReferenceRequired(&projects.ProjectId{}),
		"fabric_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"images": ImagesSchema(),
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
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile(`^/subscriptions/[0-9a-fA-F-]{36}/resourceGroups/[-\w._()]+/providers/Microsoft\.Network/virtualNetworks/[-\w._()]+/subnets/[-\w._()]+$`),
										"Subnet ID must match the format '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualNetworks/{vnetName}/subnets/{subnetName}'.",
									),
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
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
					"storage_profile": StorageProfileSchema(),
				},
			},
		},
		"maximum_concurrency": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},
		"organization_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"kind": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string("AzureDevOps"),
						}, false),
					},
					"organizations": {
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
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(pools.PossibleValuesForAzureDevOpsPermissionType(), false),
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
		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),
		"tags":     commonschema.Tags(),
	}
}

func (ManagedDevOpsPoolResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (ManagedDevOpsPoolResource) ModelObject() interface{} {
	return &ManagedDevOpsPoolModel{}
}

func (ManagedDevOpsPoolResource) ResourceType() string {
	return "azurerm_managed_devops_pool"
}

func (r ManagedDevOpsPoolResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedDevOpsPools.PoolsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ManagedDevOpsPoolModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := pools.NewPoolID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(config.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			agentProfile, err := expandAgentProfileModel(config.AgentProfile)
			if err != nil {
				return fmt.Errorf("expanding `agent_profile`: %+v", err)
			}

			organizationProfile, err := expandOrganizationProfileModel(config.OrganizationProfile)
			if err != nil {
				return fmt.Errorf("expanding `organization_profile`: %+v", err)
			}

			fabricProfile, err := expandFabricProfileModel(config.FabricProfile)
			if err != nil {
				return fmt.Errorf("expanding `fabric_profile`: %+v", err)
			}

			payload := pools.Pool{
				Name:     pointer.To(config.Name),
				Location: location.Normalize(config.Location),
				Identity: identity,
				Properties: &pools.PoolProperties{
					DevCenterProjectResourceId: config.DevCenterProjectResourceId,
					MaximumConcurrency:         config.MaximumConcurrency,
					AgentProfile:               agentProfile,
					OrganizationProfile:        organizationProfile,
					FabricProfile:              fabricProfile,
				},
				Tags: pointer.To(config.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagedDevOpsPoolResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedDevOpsPools.PoolsClient

			var config ManagedDevOpsPoolModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := pools.ParsePoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			payload := existing.Model

			if metadata.ResourceData.HasChange("identity") {
				identity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(config.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = identity
			}

			if metadata.ResourceData.HasChange("dev_center_project_resource_id") {
				payload.Properties.DevCenterProjectResourceId = config.DevCenterProjectResourceId
			}

			if metadata.ResourceData.HasChange("maximum_concurrency") {
				payload.Properties.MaximumConcurrency = config.MaximumConcurrency
			}

			if metadata.ResourceData.HasChange("agent_profile") {
				agentProfile, err := expandAgentProfileModel(config.AgentProfile)
				if err != nil {
					return fmt.Errorf("expanding `agent_profile`: %+v", err)
				}
				payload.Properties.AgentProfile = agentProfile
			}

			if metadata.ResourceData.HasChange("organization_profile") {
				organizationProfile, err := expandOrganizationProfileModel(config.OrganizationProfile)
				if err != nil {
					return fmt.Errorf("expanding `organization_profile`: %+v", err)
				}
				payload.Properties.OrganizationProfile = organizationProfile
			}

			if metadata.ResourceData.HasChange("fabric_profile") {
				fabricProfile, err := expandFabricProfileModel(config.FabricProfile)
				if err != nil {
					return fmt.Errorf("expanding `fabric_profile`: %+v", err)
				}
				payload.Properties.FabricProfile = fabricProfile
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(config.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (ManagedDevOpsPoolResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedDevOpsPools.PoolsClient

			id, err := pools.ParsePoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ManagedDevOpsPoolModel{
				Name:              id.PoolName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if modelIdentity := model.Identity; modelIdentity != nil {
					identity, err := identity.FlattenSystemAndUserAssignedMapToModel(pointer.To((identity.SystemAndUserAssignedMap)(pointer.From(model.Identity))))
					if err != nil {
						return err
					}

					state.Identity = pointer.From(identity)
				}

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

func (ManagedDevOpsPoolResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedDevOpsPools.PoolsClient

			id, err := pools.ParsePoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (ManagedDevOpsPoolResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return pools.ValidatePoolID
}
