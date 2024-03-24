// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package maintenance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MaintenanceDynamicScopeResource struct{}

type Tag struct {
	Tag    string   `tfschema:"tag"`
	Values []string `tfschema:"values"`
}

type Filter struct {
	Locations      []string `tfschema:"locations"`
	OsTypes        []string `tfschema:"os_types"`
	ResourceGroups []string `tfschema:"resource_groups"`
	ResourceTypes  []string `tfschema:"resource_types"`
	Tags           []Tag    `tfschema:"tags"`
	TagFilter      string   `tfschema:"tag_filter"`
}

type MaintenanceDynamicScopeModel struct {
	MaintenanceConfigurationId string   `tfschema:"maintenance_configuration_id"`
	Scope                      string   `tfschema:"subscription_id"`
	Location                   string   `tfschema:"location"`
	Filter                     []Filter `tfschema:"filter"`
}

var _ sdk.Resource = MaintenanceDynamicScopeResource{}

func (MaintenanceDynamicScopeResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"maintenance_configuration_id": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateFunc:     maintenanceconfigurations.ValidateMaintenanceConfigurationID,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"subscription_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubscriptionID,
		},

		"location": commonschema.Location(),

		"filter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"locations": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"os_types": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"resource_groups": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"resource_types": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"tags": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"tag": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"values": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},

					"tag_filter": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(configurationassignments.TagOperatorsAny),
							string(configurationassignments.TagOperatorsAll),
						}, false),
					},
				},
			},
		},
	}
}

func (MaintenanceDynamicScopeResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (MaintenanceDynamicScopeResource) ModelObject() interface{} {
	return &MaintenanceDynamicScopeModel{}
}

func (MaintenanceDynamicScopeResource) ResourceType() string {
	return "azurerm_maintenance_assignment_dynamic_scope"
}

func (r MaintenanceDynamicScopeResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Maintenance.ConfigurationAssignmentsClient

			var model MaintenanceDynamicScopeModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			subscriptionId, err := commonids.ParseSubscriptionID(model.Scope)
			if err != nil {
				return err
			}

			maintenanceConfigurationId, err := maintenanceconfigurations.ParseMaintenanceConfigurationID(model.MaintenanceConfigurationId)
			if err != nil {
				return err
			}

			id := configurationassignments.NewScopedConfigurationAssignmentID(subscriptionId.ID(), maintenanceConfigurationId.MaintenanceConfigurationName)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			configurationAssignment := configurationassignments.ConfigurationAssignment{
				Name:     pointer.To(id.ConfigurationAssignmentName),
				Location: pointer.To(location.Normalize(model.Location)),
				Properties: &configurationassignments.ConfigurationAssignmentProperties{
					MaintenanceConfigurationId: pointer.To(maintenanceConfigurationId.ID()),
					// test to see if required but I think it's for specific resources i.e vmss or vms
					//ResourceId: pointer.To(subscriptionId.ID()),
				},
			}

			if len(model.Filter) > 1 {
				filter := model.Filter[0]
				filterProperties := configurationassignments.ConfigurationAssignmentFilterProperties{}

				if len(filter.Locations) > 0 {
					filterProperties.Locations = pointer.To(filter.Locations)
				}

				if len(filter.OsTypes) > 0 {
					filterProperties.OsTypes = pointer.To(filter.OsTypes)
				}

				if len(filter.ResourceGroups) > 0 {
					filterProperties.ResourceGroups = pointer.To(filter.ResourceGroups)
				}

				if len(filter.ResourceTypes) > 0 {
					filterProperties.ResourceTypes = pointer.To(filter.ResourceTypes)
				}

				if len(filter.Tags) > 0 || filter.TagFilter != "" {
					tags := make(map[string][]string)
					for _, tag := range filter.Tags {
						tags[tag.Tag] = tag.Values
					}

					tagProperties := &configurationassignments.TagSettingsProperties{
						FilterOperator: pointer.To(configurationassignments.TagOperators(filter.TagFilter)),
						Tags:           pointer.To(tags),
					}
					filterProperties.TagSettings = tagProperties
				}
				configurationAssignment.Properties.Filter = pointer.To(filterProperties)
			}

			if _, err = client.CreateOrUpdate(ctx, id, configurationAssignment); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (MaintenanceDynamicScopeResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Maintenance.ConfigurationAssignmentsClient

			var state MaintenanceDynamicScopeModel
			id, err := configurationassignments.ParseScopedConfigurationAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return err
			}

			subscriptionId, err := commonids.ParseSubscriptionID(id.Scope)
			if err != nil {
				return err
			}

			state.Scope = subscriptionId.ID()

			if model := resp.Model; model != nil {
				state.Location = location.NormalizeNilable(model.Location)

				if properties := model.Properties; properties != nil {

					if properties.MaintenanceConfigurationId != nil {
						maintenanceConfigurationId, err := maintenanceconfigurations.ParseMaintenanceConfigurationIDInsensitively(pointer.From(properties.MaintenanceConfigurationId))
						if err != nil {
							return fmt.Errorf("parsing %q: %+v", pointer.From(properties.MaintenanceConfigurationId), err)
						}
						state.MaintenanceConfigurationId = maintenanceConfigurationId.ID()
					}

					if filter := properties.Filter; filter != nil {
						filterProp := make([]Filter, 0)
						tagsListProp := make([]Tag, 0)
						tagFilterProp := ""
						if tags := filter.TagSettings; tags != nil {
							tagFilterProp = string(pointer.From(tags.FilterOperator))
							for k, v := range pointer.From(tags.Tags) {
								tagsListProp = append(tagsListProp, Tag{
									Tag:    k,
									Values: v,
								})
							}
						}
						filterProp = append(filterProp, Filter{
							Locations:      pointer.From(filter.Locations),
							OsTypes:        pointer.From(filter.OsTypes),
							ResourceGroups: pointer.From(filter.ResourceGroups),
							ResourceTypes:  pointer.From(filter.ResourceTypes),
							Tags:           tagsListProp,
							TagFilter:      tagFilterProp,
						})
						state.Filter = filterProp
					}
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (MaintenanceDynamicScopeResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Maintenance.ConfigurationAssignmentsClient

			id, err := configurationassignments.ParseScopedConfigurationAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model MaintenanceDynamicScopeModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			existing := resp.Model

			if metadata.ResourceData.HasChange("maintenance_configuration_id") {
				existing.Properties.MaintenanceConfigurationId = pointer.To(model.MaintenanceConfigurationId)
			}

			if metadata.ResourceData.HasChange("location") {
				existing.Location = pointer.To(model.Location)
			}

			if metadata.ResourceData.HasChange("filter") {
				if len(model.Filter) > 1 {
					filter := model.Filter[0]
					filterProperties := configurationassignments.ConfigurationAssignmentFilterProperties{}

					if len(filter.Locations) > 0 {
						filterProperties.Locations = pointer.To(filter.Locations)
					}

					if len(filter.OsTypes) > 0 {
						filterProperties.OsTypes = pointer.To(filter.OsTypes)
					}

					if len(filter.ResourceGroups) > 0 {
						filterProperties.ResourceGroups = pointer.To(filter.ResourceGroups)
					}

					if len(filter.ResourceTypes) > 0 {
						filterProperties.ResourceTypes = pointer.To(filter.ResourceTypes)
					}

					if len(filter.Tags) > 0 || filter.TagFilter != "" {
						tags := make(map[string][]string)
						for _, tag := range filter.Tags {
							tags[tag.Tag] = tag.Values
						}

						tagProperties := &configurationassignments.TagSettingsProperties{
							FilterOperator: pointer.To(configurationassignments.TagOperators(filter.TagFilter)),
							Tags:           pointer.To(tags),
						}
						filterProperties.TagSettings = tagProperties
					}
					existing.Properties.Filter = pointer.To(filterProperties)
				} else {
					existing.Properties.Filter = &configurationassignments.ConfigurationAssignmentFilterProperties{}
				}
			}

			if _, err = client.CreateOrUpdate(ctx, pointer.From(id), *existing); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}
			return nil
		},
	}
}

func (MaintenanceDynamicScopeResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Maintenance.ConfigurationAssignmentsClient

			id, err := configurationassignments.ParseScopedConfigurationAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (MaintenanceDynamicScopeResource) IDValidationFunc() func(interface{}, string) ([]string, []error) {
	return configurationassignments.ValidateConfigurationAssignmentID
}
