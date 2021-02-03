package databoxedge

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/databoxedge/mgmt/2019-08-01/databoxedge"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databoxedge/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databoxedge/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataboxEdgeOrder() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataboxEdgeOrderCreateUpdate,
		Read:   resourceDataboxEdgeOrderRead,
		Update: resourceDataboxEdgeOrderCreateUpdate,
		Delete: resourceDataboxEdgeOrderDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DataboxEdgeOrderID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"device_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataboxEdgeName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"contact_information": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"company_name": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.DataboxEdgeCompanyName,
						},

						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.DataboxEdgeContactName,
						},

						"emails": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"phone_number": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.DataboxEdgePhoneNumber,
						},
					},
				},
			},

			"current_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"comments": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"additional_order_details": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"update_date_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"shipping_address": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_line1": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.DataboxEdgeStreetAddress,
						},

						"city": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.DataboxEdgeCity,
						},

						"country": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.DataboxEdgeCountry,
						},

						"state": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.DataboxEdgeState,
						},

						"postal_code": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.DataboxEdgePostalCode,
						},

						"address_line2": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validate.DataboxEdgeStreetAddress,
						},

						"address_line3": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validate.DataboxEdgeStreetAddress,
						},
					},
				},
			},

			"delivery_tracking_info": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"carrier_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"serial_number": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tracking_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tracking_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"order_history": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additional_order_details": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"comments": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"update_date_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"return_tracking_info": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"carrier_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"serial_number": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tracking_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tracking_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"serial_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceDataboxEdgeOrderCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataboxEdge.OrderClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	deviceName := d.Get("device_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, deviceName, resourceGroup)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Databox Edge Order (Resource Group %q / Device Name %q): %+v", resourceGroup, deviceName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_databox_edge_order", *existing.ID)
		}
	}

	order := databoxedge.Order{
		OrderProperties: &databoxedge.OrderProperties{
			ContactInformation: expandOrderContactDetails(d.Get("contact_information").([]interface{})),
			ShippingAddress:    expandOrderAddress(d.Get("shipping_address").([]interface{})),
		},
	}

	future, err := client.CreateOrUpdate(ctx, deviceName, order, resourceGroup)
	if err != nil {
		return fmt.Errorf("creating/updating Databox Edge Order (Resource Group %q / Device Name %q): %+v", resourceGroup, deviceName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating/updating future for Databox Edge Order (Resource Group %q / Device Name %q): %+v", resourceGroup, deviceName, err)
	}

	resp, err := client.Get(ctx, deviceName, resourceGroup)
	if err != nil {
		return fmt.Errorf("retrieving Databox Edge Order (Resource Group %q / Device Name %q): %+v", resourceGroup, deviceName, err)
	}

	id, err := parse.DataboxEdgeOrderID(*resp.ID)
	if err != nil {
		return err
	}

	d.SetId(id.ID(subscriptionId))

	return resourceDataboxEdgeOrderRead(d, meta)
}

func resourceDataboxEdgeOrderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataboxEdge.OrderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataboxEdgeOrderID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.DeviceName, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Databox Edge %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Databox Edge Order (Resource Group %q / Device Name %q): %+v", id.ResourceGroup, id.DeviceName, err)
	}

	d.Set("name", "default")
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("device_name", id.DeviceName)

	if props := resp.OrderProperties; props != nil {
		if err := d.Set("contact_information", flattenOrderContactDetails(props.ContactInformation)); err != nil {
			return fmt.Errorf("setting `contact_information`: %+v", err)
		}
		if err := d.Set("current_status", flattenOrderStatus(props.CurrentStatus)); err != nil {
			return fmt.Errorf("setting `current_status`: %+v", err)
		}
		if err := d.Set("shipping_address", flattenOrderAddress(props.ShippingAddress)); err != nil {
			return fmt.Errorf("setting `shipping_address`: %+v", err)
		}
		if err := d.Set("delivery_tracking_info", flattenOrderTrackingInfoArray(props.DeliveryTrackingInfo)); err != nil {
			return fmt.Errorf("setting `delivery_tracking_info`: %+v", err)
		}
		if err := d.Set("order_history", flattenOrderStatusArray(props.OrderHistory)); err != nil {
			return fmt.Errorf("setting `order_history`: %+v", err)
		}
		if err := d.Set("return_tracking_info", flattenOrderTrackingInfoArray(props.ReturnTrackingInfo)); err != nil {
			return fmt.Errorf("setting `return_tracking_info`: %+v", err)
		}
		d.Set("serial_number", props.SerialNumber)
	}

	return nil
}

func resourceDataboxEdgeOrderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataboxEdge.OrderClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataboxEdgeOrderID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.DeviceName, id.ResourceGroup)
	if err != nil {
		return fmt.Errorf("deleting Databox Edge Order (Resource Group %q / Device Name %q): %+v", id.ResourceGroup, id.DeviceName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on deleting future for Databox Edge Order (Resource Group %q / Device Name %q): %+v", id.ResourceGroup, id.DeviceName, err)
	}
	return nil
}

func expandOrderContactDetails(input []interface{}) *databoxedge.ContactDetails {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &databoxedge.ContactDetails{
		ContactPerson: utils.String(v["name"].(string)),
		CompanyName:   utils.String(v["company_name"].(string)),
		Phone:         utils.String(v["phone_number"].(string)),
		EmailList:     utils.ExpandStringSlice(v["emails"].(*schema.Set).List()),
	}
}

func expandOrderAddress(input []interface{}) *databoxedge.Address {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &databoxedge.Address{
		AddressLine1: utils.String(v["address_line1"].(string)),
		AddressLine2: utils.String(v["address_line2"].(string)),
		AddressLine3: utils.String(v["address_line3"].(string)),
		PostalCode:   utils.String(v["postal_code"].(string)),
		City:         utils.String(v["city"].(string)),
		State:        utils.String(v["state"].(string)),
		Country:      utils.String(v["country"].(string)),
	}
}

func flattenOrderContactDetails(input *databoxedge.ContactDetails) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var companyName string
	if input.CompanyName != nil {
		companyName = *input.CompanyName
	}
	var contactPerson string
	if input.ContactPerson != nil {
		contactPerson = *input.ContactPerson
	}
	var phone string
	if input.Phone != nil {
		phone = *input.Phone
	}
	return []interface{}{
		map[string]interface{}{
			"company_name": companyName,
			"name":         contactPerson,
			"emails":       utils.FlattenStringSlice(input.EmailList),
			"phone_number": phone,
		},
	}
}

func flattenOrderStatus(input *databoxedge.OrderStatus) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var status string
	if input.Status != "" {
		status = string(input.Status)
	}

	var comments string
	if input.Comments != nil {
		comments = *input.Comments
	}
	additionalOrderDetails := make(map[string]interface{})
	if input.AdditionalOrderDetails != nil {
		additionalOrderDetails = utils.FlattenMapStringPtrString(input.AdditionalOrderDetails)
	}
	var updateDateTime string
	if input.UpdateDateTime != nil {
		updateDateTime = input.UpdateDateTime.Format(time.RFC3339)
	}
	return []interface{}{
		map[string]interface{}{
			"status":                   status,
			"comments":                 comments,
			"additional_order_details": additionalOrderDetails,
			"update_date_time":         updateDateTime,
		},
	}
}

func flattenOrderAddress(input *databoxedge.Address) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var addressLine1 string
	if input.AddressLine1 != nil {
		addressLine1 = *input.AddressLine1
	}
	var city string
	if input.City != nil {
		city = *input.City
	}
	var country string
	if input.Country != nil {
		country = *input.Country
	}
	var postalCode string
	if input.PostalCode != nil {
		postalCode = *input.PostalCode
	}
	var state string
	if input.State != nil {
		state = *input.State
	}
	var addressLine2 string
	if input.AddressLine2 != nil {
		addressLine2 = *input.AddressLine2
	}
	var addressLine3 string
	if input.AddressLine3 != nil {
		addressLine3 = *input.AddressLine3
	}
	return []interface{}{
		map[string]interface{}{
			"address_line1": addressLine1,
			"address_line2": addressLine2,
			"address_line3": addressLine3,
			"city":          city,
			"country":       country,
			"postal_code":   postalCode,
			"state":         state,
		},
	}
}

func flattenOrderTrackingInfoArray(input *[]databoxedge.TrackingInfo) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var carrierName string
		if item.CarrierName != nil {
			carrierName = *item.CarrierName
		}
		var serialNumber string
		if item.SerialNumber != nil {
			serialNumber = *item.SerialNumber
		}
		var trackingId string
		if item.TrackingID != nil {
			trackingId = *item.TrackingID
		}
		var trackingUrl string
		if item.TrackingURL != nil {
			trackingUrl = *item.TrackingURL
		}
		results = append(results, map[string]interface{}{
			"carrier_name":  carrierName,
			"serial_number": serialNumber,
			"tracking_id":   trackingId,
			"tracking_url":  trackingUrl,
		})
	}
	return results
}

func flattenOrderStatusArray(input *[]databoxedge.OrderStatus) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		additionalOrderDetails := make(map[string]interface{})
		if item.AdditionalOrderDetails != nil {
			additionalOrderDetails = utils.FlattenMapStringPtrString(item.AdditionalOrderDetails)
		}
		var comments string
		if item.Comments != nil {
			comments = *item.Comments
		}
		var updateDateTime string
		if item.UpdateDateTime != nil {
			updateDateTime = item.UpdateDateTime.Format(time.RFC3339)
		}
		results = append(results, map[string]interface{}{
			"additional_order_details": additionalOrderDetails,
			"comments":                 comments,
			"update_date_time":         updateDateTime,
		})
	}
	return results
}
