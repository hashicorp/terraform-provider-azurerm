package maintenance

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/maintenance/mgmt/2021-05-01/maintenance"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const recurMondayToThursday = "weekMondayToThursday"
const recurFridayToSunday = "weekFridayToSunday"

func dataSourcePublicMaintenanceConfigurations() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourcePublicMaintenanceConfigurationsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{

			"location_filter": {
				Type:      pluginsdk.TypeString,
				Optional:  true,
				StateFunc: azure.NormalizeLocation,
			},

			"scope_filter": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"All", // All is still accepted by the API
					string(maintenance.ScopeExtension),
					string(maintenance.ScopeHost),
					string(maintenance.ScopeInGuestPatch),
					string(maintenance.ScopeOSImage),
					string(maintenance.ScopeSQLDB),
					string(maintenance.ScopeSQLManagedInstance),
				}, false),
			},

			"recur_window_filter": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					recurMondayToThursday,
					recurFridayToSunday,
				}, false),
			},

			"public_maintenance_configurations": {
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
						"recur_window": {
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
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.List(ctx)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("no Public Maintenance Configurations were found")
		}
		return fmt.Errorf("retrieving Public Maintenance Configurations: %+v", err)
	}

	filteredPublicConfigs := make([]interface{}, 0)

	recurWindowFilter := d.Get("recur_window_filter").(string)
	if recurWindowFilter == recurFridayToSunday {
		recurWindowFilter = "week Friday, Saturday, Sunday"
	}
	if recurWindowFilter == recurMondayToThursday {
		recurWindowFilter = "week Monday, Tuesday, Wednesday, Thursday"
	}

	locationFilter := d.Get("location_filter").(string)
	scopeFilter := d.Get("scope_filter").(string)
	for _, maintenanceConfig := range *resp.Value {

		var configLocation, configRecurWindow, configScope string
		if maintenanceConfig.Location != nil {
			configLocation = *maintenanceConfig.Location
		}
		if maintenanceConfig.ConfigurationProperties != nil && maintenanceConfig.ConfigurationProperties.Window != nil && maintenanceConfig.ConfigurationProperties.Window.RecurEvery != nil {
			configRecurWindow = *maintenanceConfig.ConfigurationProperties.Window.RecurEvery
		}
		if maintenanceConfig.ConfigurationProperties != nil && string(maintenanceConfig.ConfigurationProperties.MaintenanceScope) != "" {
			configScope = string(maintenanceConfig.ConfigurationProperties.MaintenanceScope)
		}

		if locationFilter == "" || locationFilter == configLocation {
			if recurWindowFilter == "" || recurWindowFilter == configRecurWindow {
				if scopeFilter == "" || scopeFilter == configScope {
					filteredPublicConfigs = append(filteredPublicConfigs, flattenPublicMaintenanceConfiguration(maintenanceConfig))
				}
			}
		}
	}

	d.Set("public_maintenance_configurations", filteredPublicConfigs)

	d.SetId(time.Now().UTC().String())
	return nil
}

func flattenPublicMaintenanceConfiguration(config maintenance.Configuration) map[string]interface{} {
	output := make(map[string]interface{})

	output["name"] = *config.Name
	output["id"] = *config.ID
	output["location"] = *config.Location
	output["maintenance_scope"] = string(config.MaintenanceScope)

	var description, recurWindow, timeZone, duration string
	if props := config.ConfigurationProperties; props != nil {
		if props.ExtensionProperties != nil {
			if configDescription, ok := props.ExtensionProperties["description"]; ok {
				description = *configDescription
			}
		}
		if props.Window != nil {
			if props.Window.RecurEvery != nil {
				recurWindow = *props.Window.RecurEvery
			}
			if props.Window.TimeZone != nil {
				timeZone = *props.Window.TimeZone
			}
			if props.Window.Duration != nil {
				duration = *props.Window.Duration
			}

		}
	}

	output["description"] = description
	output["recur_window"] = recurWindow
	output["time_zone"] = timeZone
	output["duration"] = duration

	return output
}
