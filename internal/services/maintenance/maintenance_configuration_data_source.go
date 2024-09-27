// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package maintenance

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceMaintenanceConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmMaintenanceConfigurationRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"scope": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"visibility": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"window": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"start_date_time": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"expiration_date_time": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"duration": {
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

			"properties": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"install_patches": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*schema.Schema{
						"linux": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*schema.Schema{
									"classifications_to_include": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
									"package_names_mask_to_exclude": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
									"package_names_mask_to_include": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},
						"windows": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*schema.Schema{
									"classifications_to_include": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"kb_numbers_to_exclude": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
									"kb_numbers_to_include": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},
						"reboot": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"in_guest_user_patch_mode": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceArmMaintenanceConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := maintenanceconfigurations.NewMaintenanceConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.MaintenanceConfigurationName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("scope", string(pointer.From(props.MaintenanceScope)))
			d.Set("visibility", string(pointer.From(props.Visibility)))

			properties := flattenExtensionProperties(props.ExtensionProperties)
			if properties["InGuestPatchMode"] != nil {
				d.Set("in_guest_user_patch_mode", properties["InGuestPatchMode"])
				delete(properties, "InGuestPatchMode")
			}
			d.Set("properties", properties)

			window := flattenMaintenanceConfigurationWindow(props.MaintenanceWindow)
			if err := d.Set("window", window); err != nil {
				return fmt.Errorf("setting `window`: %+v", err)
			}

			installPatches := flattenMaintenanceConfigurationInstallPatches(props.InstallPatches)
			if err := d.Set("install_patches", installPatches); err != nil {
				return fmt.Errorf("setting `install_patches`: %+v", err)
			}
		}
		d.Set("location", location.NormalizeNilable(model.Location))
		if err = tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}
