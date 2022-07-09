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
					"All", // All is still accepted by the API
					string(maintenance.ScopeExtension),
					string(maintenance.ScopeHost),
					string(maintenance.ScopeInGuestPatch),
					string(maintenance.ScopeOSImage),
					string(maintenance.ScopeSQLDB),
					string(maintenance.ScopeSQLManagedInstance),
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

	recurEveryFilter := d.Get("recur_every").(string)
	if recurEveryFilter == recurFridayToSunday {
		recurEveryFilter = "week Friday, Saturday, Sunday"
	} else if recurEveryFilter == recurMondayToThursday {
		recurEveryFilter = "week Monday, Tuesday, Wednesday, Thursday"
	}

	locationFilter := azure.NormalizeLocation(d.Get("location").(string))
	scopeFilter := d.Get("scope").(string)

	if resp.Value != nil {
		for _, maintenanceConfig := range *resp.Value {

			var configLocation, configRecurEvery, configScope string
			if maintenanceConfig.Location != nil {
				configLocation = azure.NormalizeLocation(*maintenanceConfig.Location)
			}
			if props := maintenanceConfig.ConfigurationProperties; props != nil {
				if props.Window != nil && props.Window.RecurEvery != nil {
					configRecurEvery = *props.Window.RecurEvery
				}
				if string(props.MaintenanceScope) != "" {
					configScope = string(props.MaintenanceScope)
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

	if len(filteredPublicConfigs) == 0 {
		return fmt.Errorf("no Public Maintenance Configurations were found")
	}

	if err := d.Set("configs", filteredPublicConfigs); err != nil {
		return fmt.Errorf("setting `configs`: %+v", err)
	}

	d.SetId(time.Now().UTC().String())
	return nil
}

func flattenPublicMaintenanceConfiguration(config maintenance.Configuration) map[string]interface{} {
	output := make(map[string]interface{})

	output["name"] = ""
	if config.Name != nil {
		output["name"] = *config.Name
	}

	output["id"] = ""
	if config.ID != nil {
		output["id"] = *config.ID
	}

	output["location"] = ""
	if config.Location != nil {
		output["location"] = azure.NormalizeLocation(*config.Location)
	}

	output["maintenance_scope"] = string(config.MaintenanceScope)

	var description, recurEvery, timeZone, duration string
	if props := config.ConfigurationProperties; props != nil {
		if props.ExtensionProperties != nil {
			if configDescription, ok := props.ExtensionProperties["description"]; ok {
				description = *configDescription
			}
		}
		if props.Window != nil {
			if props.Window.RecurEvery != nil {
				recurEvery = *props.Window.RecurEvery
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
	output["recur_every"] = recurEvery
	output["time_zone"] = timeZone
	output["duration"] = duration

	return output
}
