package manageddevopspools

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projects"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2025-01-21/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.ResourceWithUpdate        = ManagedDevOpsPoolResource{}
	_ sdk.ResourceWithCustomizeDiff = ManagedDevOpsPoolResource{}
)

type ManagedDevOpsPoolResource struct{}

type ManagedDevOpsPoolModel struct {
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
		"stateful_agent_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"grace_period_time_span": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "00:00:00",
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"max_agent_lifetime": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "7.00:00:00",
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"manual_resource_predictions_profile":    manualResourcePredictionsProfileSchema("stateful_agent_profile.0"),
					"automatic_resource_predictions_profile": automaticResourcePredictionsProfileSchema("stateful_agent_profile.0"),
				},
			},
			ExactlyOneOf: []string{"stateful_agent_profile", "stateless_agent_profile"},
		},
		"stateless_agent_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"manual_resource_predictions_profile":    manualResourcePredictionsProfileSchema("stateless_agent_profile.0"),
					"automatic_resource_predictions_profile": automaticResourcePredictionsProfileSchema("stateless_agent_profile.0"),
				},
			},
			ExactlyOneOf: []string{"stateful_agent_profile", "stateless_agent_profile"},
		},
		"dev_center_project_resource_id": commonschema.ResourceIDReferenceRequired(&projects.ProjectId{}),
		"vmss_fabric_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"image": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"aliases": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"buffer": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile(`^(?:\*|[0-9][0-9]?|100)$`),
										`Buffer must be "*" or value between 0 and 100.`,
									),
								},
								"resource_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"well_known_image_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
					"network_profile": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"subnet_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: commonids.ValidateSubnetID,
								},
							},
						},
					},
					"os_profile": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"logon_type": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      string(pools.LogonTypeService),
									ValidateFunc: validation.StringInSlice(pools.PossibleValuesForLogonType(), false),
								},
								"secrets_management": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"certificate_store_location": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},
											"certificate_store_name": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringInSlice(pools.PossibleValuesForCertificateStoreNameOption(), false),
											},
											"key_export_enabled": {
												Type:     pluginsdk.TypeBool,
												Required: true,
											},
											"observed_certificates": {
												Type:     pluginsdk.TypeSet,
												Required: true,
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"storage_profile": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"data_disk": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"caching": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringInSlice(pools.PossibleValuesForCachingType(), false),
											},
											"disk_size_gb": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(1, 32767),
											},
											"drive_letter": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},
											"storage_account_type": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												Default:      pools.StorageAccountTypeStandardLRS,
												ValidateFunc: validation.StringInSlice(pools.PossibleValuesForStorageAccountType(), false),
											},
										},
									},
								},
								"os_disk_storage_account_type": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringInSlice(pools.PossibleValuesForOsDiskStorageAccountType(), false),
								},
							},
						},
					},
				},
			},
		},
		"maximum_concurrency": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},
		"azure_devops_organization_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
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
						MaxItems: 1,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"kind": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(pools.PossibleValuesForAzureDevOpsPermissionType(), false),
								},
								"administrator_accounts": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"groups": {
												Type:     pluginsdk.TypeSet,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
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
				},
			},
		},
		"identity": commonschema.UserAssignedIdentityOptional(),
		"tags":     commonschema.Tags(),
	}
}

func (ManagedDevOpsPoolResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
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

			expandedIdentity, err := expandManagedDevopsToUserAssignedIdentity(config.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			var agentProfile pools.AgentProfile
			if config.StatefulAgentProfile != nil {
				agentProfile = expandStatefulAgentProfileModel(config.StatefulAgentProfile)
			} else {
				agentProfile = expandStatelessAgentProfileModel(config.StatelessAgentProfile)
			}

			azureDevOpsOrganizationProfile := expandAzureDevOpsOrganizationProfileModel(config.AzureDevOpsOrganizationProfile)

			fabricProfile := expandVmssFabricProfileModel(config.VmssFabricProfile)

			payload := pools.Pool{
				Name:     pointer.To(config.Name),
				Location: location.Normalize(config.Location),
				Identity: expandedIdentity,
				Properties: &pools.PoolProperties{
					DevCenterProjectResourceId: config.DevCenterProjectResourceId,
					MaximumConcurrency:         config.MaximumConcurrency,
					AgentProfile:               agentProfile,
					OrganizationProfile:        azureDevOpsOrganizationProfile,
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
				expandedIdentity, err := expandManagedDevopsToUserAssignedIdentity(config.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("dev_center_project_resource_id") {
				payload.Properties.DevCenterProjectResourceId = config.DevCenterProjectResourceId
			}

			if metadata.ResourceData.HasChange("maximum_concurrency") {
				payload.Properties.MaximumConcurrency = config.MaximumConcurrency
			}

			if metadata.ResourceData.HasChange("stateful_agent_profile") || metadata.ResourceData.HasChange("stateless_agent_profile") {
				var agentProfile pools.AgentProfile

				if len(config.StatefulAgentProfile) > 0 {
					agentProfile = expandStatefulAgentProfileModel(config.StatefulAgentProfile)
				} else if len(config.StatelessAgentProfile) > 0 {
					agentProfile = expandStatelessAgentProfileModel(config.StatelessAgentProfile)
				}

				payload.Properties.AgentProfile = agentProfile
			}

			if metadata.ResourceData.HasChange("azure_devops_organization_profile") {
				organizationProfile := expandAzureDevOpsOrganizationProfileModel(config.AzureDevOpsOrganizationProfile)
				payload.Properties.OrganizationProfile = organizationProfile
			}

			if metadata.ResourceData.HasChange("vmss_fabric_profile") {
				vmssFabricProfile := expandVmssFabricProfileModel(config.VmssFabricProfile)
				payload.Properties.FabricProfile = vmssFabricProfile
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

func (ManagedDevOpsPoolResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagedDevOpsPoolModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("DecodeDiff: %+v", err)
			}

			// Validate vmss_fabric_profile images
			for _, vmssFabricProfile := range model.VmssFabricProfile {
				for i, image := range vmssFabricProfile.Images {
					resourceIdSet := image.ResourceId != nil && *image.ResourceId != ""
					wellKnownImageNameSet := image.WellKnownImageName != nil && *image.WellKnownImageName != ""

					if resourceIdSet && wellKnownImageNameSet {
						return fmt.Errorf("only one of `resource_id` or `well_known_image_name` can be specified for image %d in vmss_fabric_profile", i)
					}
				}
			}

			for _, orgProfile := range model.AzureDevOpsOrganizationProfile {
				for _, permProfile := range orgProfile.PermissionProfile {
					if permProfile.Kind != "SpecificAccounts" {
						if len(permProfile.AdministratorAccounts) > 0 {
							return fmt.Errorf("administrator_accounts block is not required when permission_profile kind is '%s'", permProfile.Kind)
						}
					}
				}
			}

			// Validate agent counts don't exceed maximum_concurrency
			maxConcurrency := model.MaximumConcurrency

			// Validate stateful agent profile schedules
			for _, statefulProfile := range model.StatefulAgentProfile {
				for _, manualProfile := range statefulProfile.ManualResourcePredictionsProfile {
					if err := validateManualProfileAgentCounts("stateful_agent_profile", manualProfile, maxConcurrency); err != nil {
						return err
					}
				}
			}

			// Validate stateless agent profile schedules
			for _, statelessProfile := range model.StatelessAgentProfile {
				for _, manualProfile := range statelessProfile.ManualResourcePredictionsProfile {
					if err := validateManualProfileAgentCounts("stateless_agent_profile", manualProfile, maxConcurrency); err != nil {
						return err
					}
				}
			}

			return nil
		},
		Timeout: 5 * time.Minute,
	}
}

func validateManualProfileAgentCounts(profileType string, manualProfile ManualResourcePredictionsProfileModel, maxConcurrency int64) error {
	// Check daily schedules
	schedules := []struct {
		name     string
		schedule map[string]interface{}
	}{
		{"sunday_schedule", manualProfile.SundaySchedule},
		{"monday_schedule", manualProfile.MondaySchedule},
		{"tuesday_schedule", manualProfile.TuesdaySchedule},
		{"wednesday_schedule", manualProfile.WednesdaySchedule},
		{"thursday_schedule", manualProfile.ThursdaySchedule},
		{"friday_schedule", manualProfile.FridaySchedule},
		{"saturday_schedule", manualProfile.SaturdaySchedule},
	}

	// Validate that either all_week_schedule or at least one day schedule is specified
	hasAllWeekSchedule := manualProfile.AllWeekSchedule > 0
	hasDaySchedule := false

	for _, sched := range schedules {
		if len(sched.schedule) > 0 {
			hasDaySchedule = true
			break
		}
	}

	if !hasAllWeekSchedule && !hasDaySchedule {
		return fmt.Errorf("%s manual_resource_predictions_profile must specify either all_week_schedule or at least one day schedule (sunday_schedule, monday_schedule, tuesday_schedule, wednesday_schedule, thursday_schedule, friday_schedule, saturday_schedule)", profileType)
	}

	// Check all_week_schedule
	if manualProfile.AllWeekSchedule > 0 && manualProfile.AllWeekSchedule > maxConcurrency {
		return fmt.Errorf("%s all_week_schedule agent count (%d) cannot exceed maximum_concurrency (%d)",
			profileType, manualProfile.AllWeekSchedule, maxConcurrency)
	}

	for _, sched := range schedules {
		if err := validateScheduleAgentCounts(profileType, sched.name, sched.schedule, maxConcurrency); err != nil {
			return err
		}
	}

	return nil
}

func validateScheduleAgentCounts(profileType, scheduleName string, schedule map[string]interface{}, maxConcurrency int64) error {
	if len(schedule) == 0 {
		return nil
	}

	for timeSlot, countInterface := range schedule {
		var agentCount int64

		switch count := countInterface.(type) {
		case int:
			agentCount = int64(count)
		case int64:
			agentCount = count
		case float64:
			agentCount = int64(count)
		default:
			continue
		}

		if agentCount > maxConcurrency {
			return fmt.Errorf("%s %s time slot %s has agent count (%d) that exceeds maximum_concurrency (%d)",
				profileType, scheduleName, timeSlot, agentCount, maxConcurrency)
		}
	}

	return nil
}
