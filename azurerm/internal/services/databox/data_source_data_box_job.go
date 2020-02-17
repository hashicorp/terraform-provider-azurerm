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
						"contact_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"emails": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"phone_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mobile": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification_preference": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"send_notification": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"stage_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"phone_extension": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"destination_account_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_destination_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"share_password": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"staging_storage_account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"shipping_address": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"city": {
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
						"state_or_province": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"street_address_1": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"company_name": {
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
						"zip_extended_code": {
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

			"additional_preferred_disks_properties": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"delivery_scheduled_date_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"delivery_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"device_password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
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
			return fmt.Errorf("Error: Data Box Job (Data Box Job Name %q / Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading Data Box Job (Data Box Job Name %q / Resource Group %q): %+v", name, resourceGroup, err)
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
				if err := d.Set("destination_account_details", flattenArmDataBoxJobDestinationAccountDetails(v.DestinationAccountDetails)); err != nil {
					return fmt.Errorf("Error setting `destination_account_details`: %+v", err)
				}
				if err := d.Set("preferences", flattenArmDataBoxJobPreferences(v.Preferences)); err != nil {
					return fmt.Errorf("Error setting `preferences`: %+v", err)
				}
				if err := d.Set("shipping_address", flattenArmDataBoxJobShippingAddress(v.ShippingAddress)); err != nil {
					return fmt.Errorf("Error setting `shipping_address`: %+v", err)
				}
			} else if v, ok := details.AsDiskJobDetails(); ok && v != nil {
				if err := d.Set("additional_preferred_disks_properties", flattenArmDataBoxJobPreferredDisks(v.PreferredDisks)); err != nil {
					return fmt.Errorf("Error setting `additional_preferred_disks_properties`: %+v", err)
				}
				if err := d.Set("contact_details", flattenArmDataBoxJobContactDetails(v.ContactDetails)); err != nil {
					return fmt.Errorf("Error setting `contact_details`: %+v", err)
				}
				if err := d.Set("destination_account_details", flattenArmDataBoxJobDestinationAccountDetails(v.DestinationAccountDetails)); err != nil {
					return fmt.Errorf("Error setting `destination_account_details`: %+v", err)
				}
				if err := d.Set("preferences", flattenArmDataBoxJobPreferences(v.Preferences)); err != nil {
					return fmt.Errorf("Error setting `preferences`: %+v", err)
				}
				if err := d.Set("shipping_address", flattenArmDataBoxJobShippingAddress(v.ShippingAddress)); err != nil {
					return fmt.Errorf("Error setting `shipping_address`: %+v", err)
				}
			} else if v, ok := details.AsHeavyJobDetails(); ok && v != nil {
				d.Set("device_password", v.DevicePassword)
				if err := d.Set("contact_details", flattenArmDataBoxJobContactDetails(v.ContactDetails)); err != nil {
					return fmt.Errorf("Error setting `contact_details`: %+v", err)
				}
				if err := d.Set("destination_account_details", flattenArmDataBoxJobDestinationAccountDetails(v.DestinationAccountDetails)); err != nil {
					return fmt.Errorf("Error setting `destination_account_details`: %+v", err)
				}
				if err := d.Set("preferences", flattenArmDataBoxJobPreferences(v.Preferences)); err != nil {
					return fmt.Errorf("Error setting `preferences`: %+v", err)
				}
				if err := d.Set("shipping_address", flattenArmDataBoxJobShippingAddress(v.ShippingAddress)); err != nil {
					return fmt.Errorf("Error setting `shipping_address`: %+v", err)
				}
			}
		}
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("API returns a nil/empty id on Data Box Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	d.SetId(*resp.ID)

	return tags.FlattenAndSet(d, resp.Tags)
}
