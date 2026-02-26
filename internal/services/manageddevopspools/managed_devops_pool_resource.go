// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

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
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyvaultvalidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/manageddevopspools/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.ResourceWithUpdate        = ManagedDevOpsPoolResource{}
	_ sdk.ResourceWithCustomizeDiff = ManagedDevOpsPoolResource{}
)

type ManagedDevOpsPoolResource struct{}

type ManagedDevOpsPoolModel struct {
	DevCenterProjectId      string                         `tfschema:"dev_center_project_id"`
	VmssFabric              []VmssFabricModel              `tfschema:"vmss_fabric"`
	Identity                []identity.ModelUserAssigned   `tfschema:"identity"`
	Location                string                         `tfschema:"location"`
	MaximumConcurrency      int64                          `tfschema:"maximum_concurrency"`
	Name                    string                         `tfschema:"name"`
	AzureDevOpsOrganization []AzureDevOpsOrganizationModel `tfschema:"azure_devops_organization"`
	ResourceGroupName       string                         `tfschema:"resource_group_name"`
	Tags                    map[string]string              `tfschema:"tags"`
	StatefulAgent           []StatefulAgentModel           `tfschema:"stateful_agent"`
	StatelessAgent          []StatelessAgentModel          `tfschema:"stateless_agent"`
}

func (ManagedDevOpsPoolResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 44),
				validation.StringMatch(
					regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-.]*[a-zA-Z0-9-]$`),
					"`name` can only include alphanumeric characters, periods (.) and hyphens (-). It must also start with alphanumeric characters and cannot end with periods (.).",
				),
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"azure_devops_organization": {
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
								"url": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.All(
										validation.IsURLWithHTTPS,
										validation.StringMatch(
											regexp.MustCompile(`[a-zA-Z0-9]$`),
											"url must end with a letter or number",
										),
									),
								},

								"parallelism": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntBetween(1, 10000),
								},

								"projects": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validate.ProjectName,
									},
								},
							},
						},
					},

					"permission": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"kind": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(pools.PossibleValuesForAzureDevOpsPermissionType(), false),
								},

								"administrator_account": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"groups": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validate.Email,
												},
												AtLeastOneOf: []string{"azure_devops_organization.0.permission.0.administrator_account.0.groups", "azure_devops_organization.0.permission.0.administrator_account.0.users"},
											},

											"users": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validate.Email,
												},
												AtLeastOneOf: []string{"azure_devops_organization.0.permission.0.administrator_account.0.groups", "azure_devops_organization.0.permission.0.administrator_account.0.users"},
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

		"dev_center_project_id": commonschema.ResourceIDReferenceRequired(&projects.ProjectId{}),

		"maximum_concurrency": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 10000),
		},

		"vmss_fabric": {
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
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},

								"buffer": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "*",
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile(`^(?:\*|[0-9][0-9]?|100)$`),
										`Buffer must be "*" or value between 0 and 100.`,
									),
								},

								"id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: azure.ValidateResourceID,
								},

								"well_known_image_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"sku_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"os_disk_storage_account_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(pools.OsDiskStorageAccountTypeStandard),
						ValidateFunc: validation.StringInSlice(pools.PossibleValuesForOsDiskStorageAccountType(), false),
					},

					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},

					"security": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"interactive_logon_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"key_vault_management": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"key_export_enabled": {
												Type:     pluginsdk.TypeBool,
												Required: true,
											},

											"key_vault_certificate_ids": {
												Type:     pluginsdk.TypeList,
												Required: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: keyvaultvalidate.NestedItemIdWithOptionalVersion,
												},
											},

											"certificate_store_location": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"certificate_store_name": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringInSlice(pools.PossibleValuesForCertificateStoreNameOption(), false),
											},
										},
									},
								},
							},
						},
					},

					"storage": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"caching": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringInSlice(pools.PossibleValuesForCachingType(), false),
								},

								"disk_size_in_gb": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 65536),
								},

								"drive_letter": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile(`^[FGHIJKLMNOPQRSTUVWXYZfghijklmnopqrstuvwxyz]$`),
										"drive_letter must be a single letter and cannot be A, C, D, or E (case insensitive)",
									),
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
				},
			},
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"stateful_agent": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"grace_period_time_span": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "00:00:00",
						ValidateFunc: validate.AgentLifetime,
					},

					"max_agent_lifetime": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "7.00:00:00",
						ValidateFunc: validate.AgentLifetime,
					},

					"manual_resource_prediction": manualResourcePredictionSchema("stateful_agent.0"),

					"automatic_resource_prediction": automaticResourcePredictionSchema("stateful_agent.0"),
				},
			},
			ExactlyOneOf: []string{"stateful_agent", "stateless_agent"},
		},

		"stateless_agent": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"manual_resource_prediction": manualResourcePredictionSchema("stateless_agent.0"),

					"automatic_resource_prediction": automaticResourcePredictionSchema("stateless_agent.0"),
				},
			},
			ExactlyOneOf: []string{"stateful_agent", "stateless_agent"},
		},

		"tags": commonschema.Tags(),
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
			if _, ok := metadata.ResourceData.GetOk("stateful_agent"); ok {
				agentProfile = expandStatefulAgentModel(config.StatefulAgent)
			} else {
				agentProfile = expandStatelessAgentModel(config.StatelessAgent)
			}

			azureDevOpsOrganization := expandAzureDevOpsOrganizationModel(config.AzureDevOpsOrganization)
			fabricProfile := expandVmssFabricModel(config.VmssFabric)
			payload := pools.Pool{
				Name:     pointer.To(config.Name),
				Location: config.Location,
				Identity: expandedIdentity,
				Properties: &pools.PoolProperties{
					DevCenterProjectResourceId: config.DevCenterProjectId,
					MaximumConcurrency:         config.MaximumConcurrency,
					AgentProfile:               agentProfile,
					OrganizationProfile:        azureDevOpsOrganization,
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
				expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("dev_center_project_id") {
				payload.Properties.DevCenterProjectResourceId = config.DevCenterProjectId
			}

			if metadata.ResourceData.HasChange("maximum_concurrency") {
				payload.Properties.MaximumConcurrency = config.MaximumConcurrency
			}

			if metadata.ResourceData.HasChange("stateful_agent") || metadata.ResourceData.HasChange("stateless_agent") {
				var agentProfile pools.AgentProfile

				if _, ok := metadata.ResourceData.GetOk("stateful_agent"); ok {
					agentProfile = expandStatefulAgentModel(config.StatefulAgent)
				} else {
					agentProfile = expandStatelessAgentModel(config.StatelessAgent)
				}

				payload.Properties.AgentProfile = agentProfile
			}

			if metadata.ResourceData.HasChange("azure_devops_organization") {
				organization := expandAzureDevOpsOrganizationModel(config.AzureDevOpsOrganization)
				payload.Properties.OrganizationProfile = organization
			}

			if metadata.ResourceData.HasChange("vmss_fabric") {
				vmssFabric := expandVmssFabricModel(config.VmssFabric)
				payload.Properties.FabricProfile = vmssFabric
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
					state.DevCenterProjectId = props.DevCenterProjectResourceId
					state.MaximumConcurrency = props.MaximumConcurrency

					if agentProfile := props.AgentProfile; agentProfile != nil {
						if stateful, ok := agentProfile.(pools.Stateful); ok {
							state.StatefulAgent = flattenStatefulAgentToModel(stateful)
						} else if stateless, ok := agentProfile.(pools.StatelessAgentProfile); ok {
							state.StatelessAgent = flattenStatelessAgentToModel(stateless)
						}
					}

					if organizationProfile := props.OrganizationProfile; organizationProfile != nil {
						if azureDevOpsOrganization, ok := organizationProfile.(pools.AzureDevOpsOrganizationProfile); ok {
							state.AzureDevOpsOrganization = flattenAzureDevOpsOrganizationToModel(azureDevOpsOrganization)
						}
					}

					if fabricProfile := props.FabricProfile; fabricProfile != nil {
						if vmssFabric, ok := fabricProfile.(pools.VMSSFabricProfile); ok {
							state.VmssFabric = flattenVmssFabricToModel(vmssFabric)
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

			if err := validateVmssFabricImages(metadata, model.VmssFabric); err != nil {
				return err
			}

			for _, org := range model.AzureDevOpsOrganization {
				for _, perm := range org.Permission {
					if perm.Kind != string(pools.AzureDevOpsPermissionTypeSpecificAccounts) {
						if len(perm.AdministratorAccounts) > 0 {
							return fmt.Errorf("`administrator_account` block is not required when `permission` kind is `%s`", perm.Kind)
						}
					}
				}
			}

			maxConcurrency := model.MaximumConcurrency

			for _, stateful := range model.StatefulAgent {
				for _, manualPredictions := range stateful.ManualResourcePrediction {
					if err := validateManualAgentCounts("stateful_agent", manualPredictions, maxConcurrency); err != nil {
						return err
					}
				}
			}

			for _, stateless := range model.StatelessAgent {
				for _, manualPredictions := range stateless.ManualResourcePrediction {
					if err := validateManualAgentCounts("stateless_agent", manualPredictions, maxConcurrency); err != nil {
						return err
					}
				}
			}

			return nil
		},
		Timeout: 5 * time.Minute,
	}
}

func validateVmssFabricImages(metadata sdk.ResourceMetaData, vmssFabrics []VmssFabricModel) error {
	rawConfig := metadata.ResourceDiff.GetRawConfig().AsValueMap()

	vmssFabricValue, exists := rawConfig["vmss_fabric"]
	if !exists || vmssFabricValue.IsNull() {
		return nil
	}

	if !vmssFabricValue.IsWhollyKnown() {
		return nil
	}

	for _, vmssFabric := range vmssFabrics {
		for i, image := range vmssFabric.Images {
			haveResourceId := image.Id != ""
			haveWellKnownImageName := image.WellKnownImageName != ""

			if !haveResourceId && !haveWellKnownImageName {
				return fmt.Errorf("one of `id` or `well_known_image_name` must be specified for image %d in `vmss_fabric`", i)
			}

			if haveResourceId && haveWellKnownImageName {
				return fmt.Errorf("only one of `id` or `well_known_image_name` can be specified for image %d in `vmss_fabric`", i)
			}
		}
	}

	return nil
}

func validateManualAgentCounts(profileType string, manualPredictions ManualResourcePredictionModel, maxConcurrency int64) error {
	hasAllWeek := manualPredictions.AllWeekSchedule > 0

	daySchedules := []struct {
		name     string
		schedule []DayScheduleModel
	}{
		{"sunday_schedule", manualPredictions.SundaySchedule},
		{"monday_schedule", manualPredictions.MondaySchedule},
		{"tuesday_schedule", manualPredictions.TuesdaySchedule},
		{"wednesday_schedule", manualPredictions.WednesdaySchedule},
		{"thursday_schedule", manualPredictions.ThursdaySchedule},
		{"friday_schedule", manualPredictions.FridaySchedule},
		{"saturday_schedule", manualPredictions.SaturdaySchedule},
	}

	hasDaySchedule := false
	for _, sched := range daySchedules {
		if len(sched.schedule) > 0 {
			hasDaySchedule = true
			break
		}
	}

	if !hasAllWeek && !hasDaySchedule {
		return fmt.Errorf("%s `manual_resource_prediction` must specify either `all_week_schedule` or at least one day schedule", profileType)
	}

	if hasAllWeek {
		if manualPredictions.AllWeekSchedule > maxConcurrency {
			return fmt.Errorf("%s `all_week_schedule` agent count (%d) cannot exceed `maximum_concurrency` (%d)",
				profileType, manualPredictions.AllWeekSchedule, maxConcurrency)
		}
	}

	for _, sched := range daySchedules {
		if err := validateScheduleAgentCounts(profileType, sched.name, sched.schedule, maxConcurrency); err != nil {
			return err
		}
	}

	return nil
}

func validateScheduleAgentCounts(profileType, scheduleName string, schedule []DayScheduleModel, maxConcurrency int64) error {
	seen := make(map[string]bool, len(schedule))
	for _, entry := range schedule {
		if seen[entry.Time] {
			return fmt.Errorf("%s %s has duplicate time slot %q - entries with the same time must be merged into a single block",
				profileType, scheduleName, entry.Time)
		}
		seen[entry.Time] = true
		if entry.Count > maxConcurrency {
			return fmt.Errorf("%s %s time slot %s has agent count (%d) that exceeds `maximum_concurrency` (%d)",
				profileType, scheduleName, entry.Time, entry.Count, maxConcurrency)
		}
	}
	return nil
}
