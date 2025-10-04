// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource                   = DesktopVirtualizationHostPoolResource{}
	_ sdk.ResourceWithUpdate         = DesktopVirtualizationHostPoolResource{}
	_ sdk.ResourceWithStateMigration = DesktopVirtualizationHostPoolResource{}
)

type DesktopVirtualizationHostPoolResource struct{}

func (DesktopVirtualizationHostPoolResource) ModelObject() interface{} {
	return &DesktopVirtualizationHostPoolModel{}
}

func (DesktopVirtualizationHostPoolResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return hostpool.ValidateHostPoolID
}

func (DesktopVirtualizationHostPoolResource) ResourceType() string {
	return "azurerm_virtual_desktop_host_pool"
}

type DesktopVirtualizationHostPoolModel struct {
	Name                          string                                               `tfschema:"name"`
	Location                      string                                               `tfschema:"location"`
	ResourceGroupName             string                                               `tfschema:"resource_group_name"`
	Type                          string                                               `tfschema:"type"`
	LoadBalancerType              string                                               `tfschema:"load_balancer_type"`
	FriendlyName                  string                                               `tfschema:"friendly_name"`
	Description                   string                                               `tfschema:"description"`
	ValidateEnvironment           bool                                                 `tfschema:"validate_environment"`
	CustomRdpProperties           string                                               `tfschema:"custom_rdp_properties"`
	PersonalDesktopAssignmentType string                                               `tfschema:"personal_desktop_assignment_type"`
	PublicNetworkAccess           string                                               `tfschema:"public_network_access"`
	MaximumSessionsAllowed        int64                                                `tfschema:"maximum_sessions_allowed"`
	StartVmOnConnect              bool                                                 `tfschema:"start_vm_on_connect"`
	PreferredAppGroupType         string                                               `tfschema:"preferred_app_group_type"`
	ScheduledAgentUpdates         []DesktopVirtualizationHostPoolScheduledAgentUpdates `tfschema:"scheduled_agent_updates"`
	VmTemplate                    string                                               `tfschema:"vm_template"`
	Tags                          map[string]string                                    `tfschema:"tags"`
}
type DesktopVirtualizationHostPoolScheduledAgentUpdates struct {
	Enabled                bool                                                         `tfschema:"enabled"`
	Timezone               string                                                       `tfschema:"timezone"`
	UseSessionHostTimezone bool                                                         `tfschema:"use_session_host_timezone"`
	Schedule               []DesktopVirtualizationHostPoolScheduledAgentUpdatesSchedule `tfschema:"schedule"`
}

type DesktopVirtualizationHostPoolScheduledAgentUpdatesSchedule struct {
	DayOfWeek string `tfschema:"day_of_week"`
	HourOfDay int64  `tfschema:"hour_of_day"`
}

func (DesktopVirtualizationHostPoolResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.HostPoolV0ToV1{},
		},
	}
}

func (r DesktopVirtualizationHostPoolResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(hostpool.PossibleValuesForHostPoolType(), false),
		},

		"load_balancer_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(hostpool.PossibleValuesForLoadBalancerType(), false),
		},

		"friendly_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 64),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 512),
		},

		"validate_environment": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"custom_rdp_properties": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"personal_desktop_assignment_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(hostpool.PossibleValuesForPersonalDesktopAssignmentType(), false),
		},

		"public_network_access": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(hostpool.PossibleValuesForHostpoolPublicNetworkAccess(), false),
			Default:      string(hostpool.HostpoolPublicNetworkAccessEnabled),
		},

		"maximum_sessions_allowed": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      999999,
			ValidateFunc: validation.IntBetween(0, 999999),
		},

		"start_vm_on_connect": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"preferred_app_group_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Description:  "Preferred App Group type to display",
			ValidateFunc: validation.StringInSlice(hostpool.PossibleValuesForPreferredAppGroupType(), false),
			Default:      string(hostpool.PreferredAppGroupTypeDesktop),
		},

		"scheduled_agent_updates": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"timezone": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "UTC",
					},

					"use_session_host_timezone": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"schedule": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 2,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"day_of_week": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(hostpool.PossibleValuesForDayOfWeek(), false),
								},

								"hour_of_day": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(0, 23),
								},
							},
						},
					},
				},
			},
		},

		"vm_template": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"tags": commonschema.Tags(),
	}
}

func (r DesktopVirtualizationHostPoolResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DesktopVirtualizationHostPoolResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.HostPoolsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DesktopVirtualizationHostPoolModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := hostpool.NewHostPoolID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			personalDesktopAssignmentType := hostpool.PersonalDesktopAssignmentType(model.PersonalDesktopAssignmentType)
			vmTemplate := model.VmTemplate
			if vmTemplate != "" {
				// we have no use with the json object as azure accepts string only
				// merely here for validation
				_, err := pluginsdk.ExpandJsonFromString(vmTemplate)
				if err != nil {
					return fmt.Errorf("expanding JSON for `vm_template`: %+v", err)
				}
			}
			payload := hostpool.HostPool{
				Location: location.Normalize(model.Location),
				Tags:     pointer.To(model.Tags),
				Properties: hostpool.HostPoolProperties{
					HostPoolType:                  hostpool.HostPoolType(model.Type),
					FriendlyName:                  pointer.To(model.FriendlyName),
					Description:                   pointer.To(model.Description),
					ValidationEnvironment:         pointer.To(model.ValidateEnvironment),
					CustomRdpProperty:             pointer.To(model.CustomRdpProperties),
					MaxSessionLimit:               pointer.To(model.MaximumSessionsAllowed),
					StartVMOnConnect:              pointer.To(model.StartVmOnConnect),
					LoadBalancerType:              hostpool.LoadBalancerType(model.LoadBalancerType),
					PersonalDesktopAssignmentType: &personalDesktopAssignmentType,
					PreferredAppGroupType:         hostpool.PreferredAppGroupType(model.PreferredAppGroupType),
					PublicNetworkAccess:           pointer.To(hostpool.HostpoolPublicNetworkAccess(model.PublicNetworkAccess)),
					AgentUpdate:                   expandAgentUpdateCreate(model.ScheduledAgentUpdates),
					VMTemplate:                    pointer.To(vmTemplate),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DesktopVirtualizationHostPoolResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.HostPoolsClient

			var model DesktopVirtualizationHostPoolModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := hostpool.ParseHostPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.HostPoolName, r.ResourceType())
			defer locks.UnlockByName(id.HostPoolName, r.ResourceType())

			payload := hostpool.HostPoolPatch{}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if metadata.ResourceData.HasChanges("custom_rdp_properties", "description", "friendly_name", "load_balancer_type", "maximum_sessions_allowed", "preferred_app_group_type", "public_network_access", "start_vm_on_connect", "validate_environment", "scheduled_agent_updates") {
				payload.Properties = &hostpool.HostPoolPatchProperties{}

				if metadata.ResourceData.HasChange("custom_rdp_properties") {
					payload.Properties.CustomRdpProperty = pointer.To(model.CustomRdpProperties)
				}

				if metadata.ResourceData.HasChange("description") {
					payload.Properties.Description = pointer.To(model.Description)
				}

				if metadata.ResourceData.HasChange("friendly_name") {
					payload.Properties.FriendlyName = pointer.To(model.FriendlyName)
				}

				if metadata.ResourceData.HasChange("load_balancer_type") {
					loadBalancerType := hostpool.LoadBalancerType(model.LoadBalancerType)
					payload.Properties.LoadBalancerType = &loadBalancerType
				}

				if metadata.ResourceData.HasChange("maximum_sessions_allowed") {
					payload.Properties.MaxSessionLimit = pointer.To(model.MaximumSessionsAllowed)
				}

				if metadata.ResourceData.HasChange("preferred_app_group_type") {
					preferredAppGroupType := hostpool.PreferredAppGroupType(model.PreferredAppGroupType)
					payload.Properties.PreferredAppGroupType = &preferredAppGroupType
				}

				if metadata.ResourceData.HasChange("public_network_access") {
					payload.Properties.PublicNetworkAccess = pointer.To(hostpool.HostpoolPublicNetworkAccess(model.PublicNetworkAccess))
				}

				if metadata.ResourceData.HasChange("start_vm_on_connect") {
					payload.Properties.StartVMOnConnect = pointer.To(model.StartVmOnConnect)
				}

				if metadata.ResourceData.HasChange("validate_environment") {
					payload.Properties.ValidationEnvironment = pointer.To(model.ValidateEnvironment)
				}

				if metadata.ResourceData.HasChanges("scheduled_agent_updates") {
					payload.Properties.AgentUpdate = expandAgentUpdatePatch(model.ScheduledAgentUpdates)
				}

				if metadata.ResourceData.HasChanges("vm_template") {
					payload.Properties.VMTemplate = pointer.To(model.VmTemplate)
				}
			}

			if _, err := client.Update(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r DesktopVirtualizationHostPoolResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.HostPoolsClient

			state := DesktopVirtualizationHostPoolModel{}

			id, err := hostpool.ParseHostPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					log.Printf("[DEBUG] %s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state.Name = id.HostPoolName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
				props := model.Properties

				state.CustomRdpProperties = pointer.From(props.CustomRdpProperty)
				state.Description = pointer.From(props.Description)
				state.FriendlyName = pointer.From(props.FriendlyName)
				state.MaximumSessionsAllowed = pointer.From(props.MaxSessionLimit)
				state.LoadBalancerType = string(props.LoadBalancerType)
				personalDesktopAssignmentType := ""
				if props.PersonalDesktopAssignmentType != nil {
					personalDesktopAssignmentType = string(*props.PersonalDesktopAssignmentType)
				}
				state.PersonalDesktopAssignmentType = personalDesktopAssignmentType
				state.PreferredAppGroupType = string(props.PreferredAppGroupType)
				state.PublicNetworkAccess = string(pointer.From(props.PublicNetworkAccess))
				state.StartVmOnConnect = pointer.From(props.StartVMOnConnect)
				state.Type = string(props.HostPoolType)
				state.ValidateEnvironment = pointer.From(props.ValidationEnvironment)
				state.ScheduledAgentUpdates = flattenAgentUpdate(props.AgentUpdate)
				state.VmTemplate = pointer.From(props.VMTemplate)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DesktopVirtualizationHostPoolResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.HostPoolsClient
			timeout, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("context is missing a timeout")
			}

			id, err := hostpool.ParseHostPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.HostPoolName, r.ResourceType())
			defer locks.UnlockByName(id.HostPoolName, r.ResourceType())

			options := hostpool.DeleteOperationOptions{
				Force: pointer.To(true),
			}

			err = pluginsdk.Retry(time.Until(timeout), func() *pluginsdk.RetryError {
				_, err := client.Delete(ctx, *id, options)
				if err == nil {
					return nil
				}
				if strings.Contains(err.Error(), fmt.Sprintf("The SessionHostPool %s could not be deleted because it still has ApplicationGroups associated with it. Please remove all ApplicationGroups from the SessionHostPool and retry the operation.", id.HostPoolName)) {
					return pluginsdk.RetryableError(fmt.Errorf(" %s still has ApplicationGroups attached, retrying", id.HostPoolName))
				}
				return pluginsdk.NonRetryableError(fmt.Errorf("deleting %s: %+v", *id, err))
			})
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func expandAgentUpdateSchedule(input []DesktopVirtualizationHostPoolScheduledAgentUpdatesSchedule) *[]hostpool.MaintenanceWindowProperties {
	if len(input) == 0 {
		return nil
	}

	results := make([]hostpool.MaintenanceWindowProperties, 0)
	for _, item := range input {
		dayOfWeek := hostpool.DayOfWeek(item.DayOfWeek)
		hourOfDay := item.HourOfDay

		results = append(results, hostpool.MaintenanceWindowProperties{
			DayOfWeek: &dayOfWeek,
			Hour:      pointer.To(hourOfDay),
		})
	}

	return &results
}

func expandAgentUpdateCreate(input []DesktopVirtualizationHostPoolScheduledAgentUpdates) *hostpool.AgentUpdateProperties {
	if len(input) == 0 {
		return nil
	}

	raw := input[0]

	props := hostpool.AgentUpdateProperties{}
	updatesScheduled := hostpool.SessionHostComponentUpdateTypeScheduled
	updatesDefault := hostpool.SessionHostComponentUpdateTypeDefault

	useSessionHostLocalTime := raw.UseSessionHostTimezone
	updateScheduleTimeZone := pointer.To(raw.Timezone)

	if raw.Enabled {
		props.Type = &updatesScheduled
		if !useSessionHostLocalTime { // based on the priority used in the Azure Portal, if Session Host time is selected, this overrides the explicit TimeZone setting
			props.MaintenanceWindowTimeZone = updateScheduleTimeZone
			props.UseSessionHostLocalTime = &useSessionHostLocalTime
			props.MaintenanceWindows = expandAgentUpdateSchedule(raw.Schedule)
		}
	} else {
		props.Type = &updatesDefault
		props.MaintenanceWindows = &[]hostpool.MaintenanceWindowProperties{}
		props.UseSessionHostLocalTime = &useSessionHostLocalTime // required by REST API even when set to Default/Disabled
		props.MaintenanceWindowTimeZone = updateScheduleTimeZone // required by REST API even when set to Default/Disabled
	}

	return &props
}

func expandAgentUpdatePatch(input []DesktopVirtualizationHostPoolScheduledAgentUpdates) *hostpool.AgentUpdatePatchProperties {
	if len(input) == 0 {
		return nil
	}

	raw := input[0]

	props := hostpool.AgentUpdatePatchProperties{}
	updatesScheduled := hostpool.SessionHostComponentUpdateTypeScheduled
	updatesDefault := hostpool.SessionHostComponentUpdateTypeDefault

	props.MaintenanceWindowTimeZone = pointer.To(raw.Timezone)
	props.UseSessionHostLocalTime = pointer.To(raw.UseSessionHostTimezone)

	if raw.Enabled {
		props.Type = &updatesScheduled
		props.MaintenanceWindows = expandAgentUpdateSchedulePatch(raw.Schedule)
	} else {
		props.Type = &updatesDefault
		props.MaintenanceWindows = &[]hostpool.MaintenanceWindowPatchProperties{}
	}

	return &props
}

func expandAgentUpdateSchedulePatch(input []DesktopVirtualizationHostPoolScheduledAgentUpdatesSchedule) *[]hostpool.MaintenanceWindowPatchProperties {
	if len(input) == 0 {
		return nil
	}

	results := make([]hostpool.MaintenanceWindowPatchProperties, 0)
	for _, item := range input {
		dayOfWeek := hostpool.DayOfWeek(item.DayOfWeek)
		hourOfDay := item.HourOfDay

		results = append(results, hostpool.MaintenanceWindowPatchProperties{
			DayOfWeek: &dayOfWeek,
			Hour:      pointer.To(hourOfDay),
		})
	}

	return &results
}

func flattenAgentUpdateSchedule(input *[]hostpool.MaintenanceWindowProperties) []DesktopVirtualizationHostPoolScheduledAgentUpdatesSchedule {
	results := make([]DesktopVirtualizationHostPoolScheduledAgentUpdatesSchedule, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		dayOfWeek := ""
		if item.DayOfWeek != nil {
			dayOfWeek = *pointer.To(string(*item.DayOfWeek))
		}
		hourOfDay := pointer.To(int64(0))
		if item.Hour != nil {
			hourOfDay = item.Hour
		}
		results = append(results, DesktopVirtualizationHostPoolScheduledAgentUpdatesSchedule{
			DayOfWeek: dayOfWeek,
			HourOfDay: pointer.From(hourOfDay),
		})
	}
	return results
}

func flattenAgentUpdate(input *hostpool.AgentUpdateProperties) []DesktopVirtualizationHostPoolScheduledAgentUpdates {
	if input == nil {
		return []DesktopVirtualizationHostPoolScheduledAgentUpdates{}
	}
	enabled := false
	if input.Type != nil {
		if *input.Type == hostpool.SessionHostComponentUpdateTypeScheduled {
			enabled = true
		}
	}

	return []DesktopVirtualizationHostPoolScheduledAgentUpdates{
		{
			Enabled:                enabled,
			Timezone:               pointer.From(input.MaintenanceWindowTimeZone),
			UseSessionHostTimezone: pointer.From(input.UseSessionHostLocalTime),
			Schedule:               flattenAgentUpdateSchedule(input.MaintenanceWindows),
		},
	}
}
