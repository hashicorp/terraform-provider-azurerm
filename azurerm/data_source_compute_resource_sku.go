package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmComputeResourceSku() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmComputeResourceSkuRead,

		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocation(),

			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tier": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"size": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"family": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"capacity": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"minimum": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"maximum": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"default": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scale_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"locations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"location_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"api_versions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"costs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"meter_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"quantity": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"extended_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"capabilities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"restrictions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"restriction_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"restriction_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"locations": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},

									"zones": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"reason_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmComputeResourceSkuRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).compute.ResourceSkusClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	location := d.Get("location").(string)

	resourceSkus, err := client.ListComplete(ctx)

	if err != nil {
		return fmt.Errorf("Error loading Compute Resource SKUs: %+v", err)
	}

	filteredResourceSkus := make([]compute.ResourceSku, 0)
	for resourceSkus.NotDone() {
		sku := resourceSkus.Value()

		if v, ok := d.GetOk("location"); ok {
			if location := v.(string); location != "" {
				if utils.StringSliceContains(*sku.Locations, location) && *sku.Name == name {
					filteredResourceSkus = append(filteredResourceSkus, sku)
				}
			}
		}

		err = resourceSkus.NextWithContext(ctx)
		if err != nil {
			return fmt.Errorf("Error loading Compute Resource SKUs: %+v", err)
		}
	}

	d.SetId(time.Now().UTC().String())
	d.Set("name", name)

	log.Printf("[DEBUG] filteredResourceSkus: %+v\n", filteredResourceSkus)

	if len(filteredResourceSkus) == 0 {
		return fmt.Errorf("Error loading Resource SKU: could not find SKU '%s'. Invalid SKU Name or not valid in this subscription or location.", name)
	} else if len(filteredResourceSkus) > 1 {
		return fmt.Errorf("Error loading Resource SKU: multiple results found for %s/%s.", location, name)
	} else {

		frs := filteredResourceSkus[0]

		d.Set("resource_type", frs.ResourceType)
		d.Set("tier", frs.Tier)
		d.Set("size", frs.Size)
		d.Set("family", frs.Family)
		d.Set("kind", frs.Kind)

		capacity := flattenCapacity(frs.Capacity)
		if err := d.Set("capacity", capacity); err != nil {
			return fmt.Errorf("Error setting `capacity`: %+v", err)
		}

		locationInfo := flattenLocationInfo(frs.LocationInfo)
		if err := d.Set("location_info", locationInfo); err != nil {
			return fmt.Errorf("Error setting `location_info`: %+v", err)
		}

		d.Set("api_versions", frs.APIVersions)

		costs := flattenCosts(frs.Costs)
		if err := d.Set("costs", costs); err != nil {
			return fmt.Errorf("Error setting `costs`: %+v", err)
		}

		capabilities := flattenCapabilities(frs.Capabilities)
		if err := d.Set("capabilities", capabilities); err != nil {
			return fmt.Errorf("Error setting `capabilities`: %+v", err)
		}

		restrictions := flattenRestrictions(frs.Restrictions)
		if err := d.Set("restrictions", restrictions); err != nil {
			return fmt.Errorf("Error setting `restrictions`: %+v", err)
		}
	}

	return nil
}

func flattenCapacity(input *compute.ResourceSkuCapacity) interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	prop := *input
	output := make(map[string]interface{})

	if prop.Minimum != nil {
		output["minimum"] = *prop.Minimum
	}

	if prop.Maximum != nil {
		output["maximum"] = *prop.Maximum
	}

	if prop.Default != nil {
		output["default"] = *prop.Default
	}

	output["scale_type"] = prop.ScaleType

	result = append(result, output)
	return result
}

func flattenCapabilities(input *[]compute.ResourceSkuCapabilities) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	for _, prop := range *input {
		output := make(map[string]interface{})

		if prop.Name != nil {
			output["name"] = *prop.Name
		}

		if prop.Value != nil {
			output["value"] = *prop.Value
		}

		result = append(result, output)
	}
	return result
}

func flattenCosts(input *[]compute.ResourceSkuCosts) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	for _, prop := range *input {
		output := make(map[string]interface{})

		if prop.MeterID != nil {
			output["meter_id"] = *prop.MeterID
		}

		if prop.Quantity != nil {
			output["quantity"] = *prop.Quantity
		}

		if prop.ExtendedUnit != nil {
			output["extended_unit"] = *prop.ExtendedUnit
		}

		result = append(result, output)
	}
	return result
}

func flattenLocationInfo(input *[]compute.ResourceSkuLocationInfo) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	for _, prop := range *input {
		output := make(map[string]interface{})

		if prop.Location != nil {
			output["location"] = azure.NormalizeLocation(*prop.Location)
		}

		if prop.Zones != nil {
			output["zones"] = *prop.Zones
		}

		result = append(result, output)
	}
	return result
}

func flattenRestrictions(input *[]compute.ResourceSkuRestrictions) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	for _, prop := range *input {
		output := make(map[string]interface{})

		output["restriction_type"] = prop.Type

		if prop.Values != nil {
			output["values"] = *prop.Values
		}

		if prop.RestrictionInfo != nil {
			output["restriction_info"] = flattenRestrictionInfo(prop.RestrictionInfo)
		}

		output["reason_code"] = prop.ReasonCode

		result = append(result, output)
	}
	return result
}

func flattenRestrictionInfo(input *compute.ResourceSkuRestrictionInfo) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	prop := *input
	output := make(map[string]interface{})

	if prop.Locations != nil {
		locations := make([]string, 0)
		if loc := prop.Locations; loc != nil {
			locations = *loc
		}
		output["locations"] = locations
	}

	if prop.Zones != nil {
		output["zones"] = *prop.Zones
	}

	result = append(result, output)

	return result
}
