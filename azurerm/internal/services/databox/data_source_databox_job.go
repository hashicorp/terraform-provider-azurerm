package databox

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmDataBoxJob() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmDataBoxJobRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateDataBoxJobName,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"contact_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"emails": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"mobile": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification_preference": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"at_azure_dc": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"data_copied": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"delivered": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"device_prepared": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"dispatched": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"picked_up": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"phone_extension": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"phone_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"databox_preferred_disk_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"databox_preferred_disk_size_in_tb": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"datacenter_region_preference": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"delivery_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"delivery_scheduled_date_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"destination_account": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"share_password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"staging_storage_account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"device_password": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"expected_data_size_in_tb": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"preferred_shipment_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"shipping_address": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"city": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"company_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"postal_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"postal_code_ext": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_or_province": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"street_address_1": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"street_address_2": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"street_address_3": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"sku_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmDataBoxJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBox.JobClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name, "Details")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: DataBox Job (DataBox Job Name %q / Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading DataBox Job (DataBox Job Name %q / Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("sku_name", resp.Sku.Name)
	if props := resp.JobProperties; props != nil {
		if props.DeliveryInfo != nil && props.DeliveryInfo.ScheduledDateTime != nil {
			d.Set("delivery_scheduled_date_time", (*props.DeliveryInfo.ScheduledDateTime).Format(time.RFC3339))
		}
		d.Set("delivery_type", props.DeliveryType)

		if details := props.Details; details != nil {
			if v, ok := details.AsJobDetailsType(); ok && v != nil {
				d.Set("device_password", v.DevicePassword)

				if err := d.Set("contact_details", flattenArmDataBoxJobContactDetails(v.ContactDetails)); err != nil {
					return fmt.Errorf("Error setting `contact_details`: %+v", err)
				}
				if err := d.Set("destination_account", flattenArmDataBoxJobDestinationAccount(v.DestinationAccountDetails)); err != nil {
					return fmt.Errorf("Error setting `destination_account`: %+v", err)
				}
				if v.Preferences != nil {
					if err := d.Set("datacenter_region_preference", utils.FlattenStringSlice(v.Preferences.PreferredDataCenterRegion)); err != nil {
						return fmt.Errorf("Error setting `datacenter_region_preference`: %+v", err)
					}
					if v.Preferences.TransportPreferences != nil {
						if err := d.Set("preferred_shipment_type", v.Preferences.TransportPreferences.PreferredShipmentType); err != nil {
							return fmt.Errorf("Error setting `preferred_shipment_type`: %+v", err)
						}
					}
				}
				if err := d.Set("shipping_address", flattenArmDataBoxJobShippingAddress(v.ShippingAddress)); err != nil {
					return fmt.Errorf("Error setting `shipping_address`: %+v", err)
				}
			} else if v, ok := details.AsDiskJobDetails(); ok && v != nil {
				for k, v := range v.PreferredDisks {
					if err := d.Set("databox_preferred_disk_count", k); err != nil {
						return fmt.Errorf("Error setting `databox_preferred_disk_count`: %+v", err)
					}
					if err := d.Set("databox_preferred_disk_size_in_tb", v); err != nil {
						return fmt.Errorf("Error setting `databox_preferred_disk_size_in_tb`: %+v", err)
					}
				}
				if err := d.Set("contact_details", flattenArmDataBoxJobContactDetails(v.ContactDetails)); err != nil {
					return fmt.Errorf("Error setting `contact_details`: %+v", err)
				}
				if err := d.Set("destination_account", flattenArmDataBoxJobDestinationAccount(v.DestinationAccountDetails)); err != nil {
					return fmt.Errorf("Error setting `destination_account`: %+v", err)
				}
				if v.Preferences != nil {
					if v.Preferences.TransportPreferences != nil {
						if err := d.Set("preferred_shipment_type", v.Preferences.TransportPreferences.PreferredShipmentType); err != nil {
							return fmt.Errorf("Error setting `preferred_shipment_type`: %+v", err)
						}
					}
					if err := d.Set("datacenter_region_preference", utils.FlattenStringSlice(v.Preferences.PreferredDataCenterRegion)); err != nil {
						return fmt.Errorf("Error setting `datacenter_region_preference`: %+v", err)
					}
				}
				if err := d.Set("shipping_address", flattenArmDataBoxJobShippingAddress(v.ShippingAddress)); err != nil {
					return fmt.Errorf("Error setting `shipping_address`: %+v", err)
				}
			} else if v, ok := details.AsHeavyJobDetails(); ok && v != nil {
				d.Set("device_password", v.DevicePassword)

				if err := d.Set("contact_details", flattenArmDataBoxJobContactDetails(v.ContactDetails)); err != nil {
					return fmt.Errorf("Error setting `contact_details`: %+v", err)
				}
				if err := d.Set("destination_account", flattenArmDataBoxJobDestinationAccount(v.DestinationAccountDetails)); err != nil {
					return fmt.Errorf("Error setting `destination_account`: %+v", err)
				}
				if v.Preferences != nil {
					if v.Preferences.TransportPreferences != nil {
						if err := d.Set("preferred_shipment_type", v.Preferences.TransportPreferences.PreferredShipmentType); err != nil {
							return fmt.Errorf("Error setting `preferred_shipment_type`: %+v", err)
						}
					}
					if err := d.Set("datacenter_region_preference", utils.FlattenStringSlice(v.Preferences.PreferredDataCenterRegion)); err != nil {
						return fmt.Errorf("Error setting `datacenter_region_preference`: %+v", err)
					}
				}
				if err := d.Set("shipping_address", flattenArmDataBoxJobShippingAddress(v.ShippingAddress)); err != nil {
					return fmt.Errorf("Error setting `shipping_address`: %+v", err)
				}
			}
		}
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("API returns a nil/empty id on DataBox Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	d.SetId(*resp.ID)

	return tags.FlattenAndSet(d, resp.Tags)
}
