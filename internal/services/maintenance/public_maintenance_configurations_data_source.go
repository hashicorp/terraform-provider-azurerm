// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package maintenance

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/publicmaintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

const recurMondayToThursday = "Monday-Thursday"
const recurFridayToSunday = "Friday-Sunday"

func dataSourcePublicMaintenanceConfigurations() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourcePublicMaintenanceConfigurationsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{

			"location": {
				Type:      pluginsdk.TypeString,
				Optional:  true,
				StateFunc: azure.NormalizeLocation,
			},

			"scope": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(publicmaintenanceconfigurations.MaintenanceScopeExtension),
					string(publicmaintenanceconfigurations.MaintenanceScopeHost),
					string(publicmaintenanceconfigurations.MaintenanceScopeInGuestPatch),
					string(publicmaintenanceconfigurations.MaintenanceScopeOSImage),
					string(publicmaintenanceconfigurations.MaintenanceScopeSQLDB),
					string(publicmaintenanceconfigurations.MaintenanceScopeSQLManagedInstance),
				}, false),
			},

			"recur_every": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					recurMondayToThursday,
					recurFridayToSunday,
				}, false),
			},

			"configs": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"location": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"description": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"duration": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"maintenance_scope": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"time_zone": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"recur_every": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePublicMaintenanceConfigurationsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.PublicConfigurationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subId := commonids.NewSubscriptionID(subscriptionId)
	resp, err := client.List(ctx, subId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("no Public Maintenance Configurations were found")
		}
		return fmt.Errorf("retrieving Public Maintenance Configurations: %+v", err)
	}

	filteredPublicConfigs := make([]interface{}, 0)

	recurEveryFilterRaw := d.Get("recur_every").(string)
	recurEveryFilter := recurEveryFilterRaw
	if recurEveryFilterRaw == recurFridayToSunday {
		recurEveryFilter = "week Friday, Saturday, Sunday"
	} else if recurEveryFilterRaw == recurMondayToThursday {
		recurEveryFilter = "week Monday, Tuesday, Wednesday, Thursday"
	}

	locationFilter := azure.NormalizeLocation(d.Get("location").(string))
	scopeFilter := d.Get("scope").(string)

	if resp.Model != nil {
		if resp.Model.Value != nil {
			for _, maintenanceConfig := range *resp.Model.Value {
				var configLocation, configRecurEvery, configScope string
				if maintenanceConfig.Location != nil {
					configLocation = azure.NormalizeLocation(*maintenanceConfig.Location)
				}
				if props := maintenanceConfig.Properties; props != nil {
					if props.MaintenanceWindow != nil && props.MaintenanceWindow.RecurEvery != nil {
						configRecurEvery = *props.MaintenanceWindow.RecurEvery
					}

					if props.MaintenanceScope != nil {
						configScope = string(*props.MaintenanceScope)
					}
				}

				if locationFilter != "" && locationFilter != configLocation {
					continue
				}
				if recurEveryFilter != "" && recurEveryFilter != configRecurEvery {
					continue
				}
				if scopeFilter != "" && scopeFilter != configScope {
					continue
				}

				filteredPublicConfigs = append(filteredPublicConfigs, flattenPublicMaintenanceConfiguration(maintenanceConfig))
			}
		}
	}
	if len(filteredPublicConfigs) == 0 {
		return fmt.Errorf("no Public Maintenance Configurations were found")
	}

	if err := d.Set("configs", filteredPublicConfigs); err != nil {
		return fmt.Errorf("setting `configs`: %+v", err)
	}

	id := fmt.Sprintf("publicMaintenanceConfigurations/location=%s;scope=%s;recurEvery=%s", locationFilter, scopeFilter, recurEveryFilterRaw)
	d.SetId(base64.StdEncoding.EncodeToString([]byte(id)))
	return nil
}

func flattenPublicMaintenanceConfiguration(config publicmaintenanceconfigurations.MaintenanceConfiguration) map[string]interface{} {
	output := make(map[string]interface{})

	output["name"] = ""
	if config.Name != nil {
		output["name"] = *config.Name
	}

	output["id"] = ""
	if config.Id != nil {
		output["id"] = *config.Id
	}

	output["location"] = ""
	if config.Location != nil {
		output["location"] = azure.NormalizeLocation(*config.Location)
	}

	var description, recurEvery, timeZone, duration, scope string
	if props := config.Properties; props != nil {
		if props.ExtensionProperties != nil {
			extensionProperties := *props.ExtensionProperties
			if configDescription, ok := extensionProperties["description"]; ok {
				description = configDescription
			}
		}

		if config.Properties.MaintenanceScope != nil {
			scope = string(*config.Properties.MaintenanceScope)
		}

		if props.MaintenanceWindow != nil {
			if props.MaintenanceWindow.RecurEvery != nil {
				recurEvery = *props.MaintenanceWindow.RecurEvery
			}
			if props.MaintenanceWindow.TimeZone != nil {
				timeZone = *props.MaintenanceWindow.TimeZone
			}
			if props.MaintenanceWindow.Duration != nil {
				duration = *props.MaintenanceWindow.Duration
			}
		}
	}

	output["description"] = description
	output["recur_every"] = recurEvery
	output["time_zone"] = timeZone
	output["duration"] = duration
	output["maintenance_scope"] = scope

	return output
}
