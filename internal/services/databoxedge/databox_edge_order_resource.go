package databoxedge

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/databoxedge/mgmt/2020-12-01/databoxedge"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databoxedge/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databoxedge/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceOrder() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceOrderCreateUpdate,
		Read:   resourceOrderRead,
		Update: resourceOrderCreateUpdate,
		Delete: resourceOrderDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.OrderID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"device_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataboxEdgeName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"contact": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"company_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.DataboxEdgeCompanyName,
						},

						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.DataboxEdgeContactName,
						},

						"emails": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							ForceNew: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"phone_number": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.DataboxEdgePhoneNumber,
						},
					},
				},
			},

			"status": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"info": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"comments": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"additional_details": {
							Type:     pluginsdk.TypeMap,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"last_update": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"shipment_address": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"address": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MaxItems: 3,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"city": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.DataboxEdgeCity,
						},

						"country": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.DataboxEdgeCountry,
						},

						"state": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.DataboxEdgeState,
						},

						"postal_code": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.DataboxEdgePostalCode,
						},
					},
				},
			},

			"shipment_tracking": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"carrier_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"serial_number": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"tracking_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"tracking_url": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"shipment_history": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"additional_details": {
							Type:     pluginsdk.TypeMap,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"comments": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"last_update": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"return_tracking": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"carrier_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"serial_number": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"tracking_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"tracking_url": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"serial_number": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(databoxEdgeCustomizeDiff),
	}
}

func resourceOrderCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataboxEdge.OrderClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewOrderID(subscriptionId, d.Get("resource_group_name").(string), d.Get("device_name").(string), "default")
	if d.IsNewResource() {
		// SDK method: Get(ctx context.Context, deviceName string, resourceGroupName string)
		existing, err := client.Get(ctx, id.DataBoxEdgeDeviceName, id.ResourceGroup)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_databox_edge_order", id.ID())
		}
	}

	order := databoxedge.Order{
		OrderProperties: &databoxedge.OrderProperties{
			ContactInformation: expandOrderContactDetails(d.Get("contact").([]interface{})),
			ShippingAddress:    expandOrderAddress(d.Get("shipment_address").([]interface{})),
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.DataBoxEdgeDeviceName, order, id.ResourceGroup)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceOrderRead(d, meta)
}

func resourceOrderRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataboxEdge.OrderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.OrderID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.DataBoxEdgeDeviceName, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s does not exist - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("device_name", id.DataBoxEdgeDeviceName)

	if props := resp.OrderProperties; props != nil {
		if err := d.Set("contact", flattenOrderContactDetails(props.ContactInformation)); err != nil {
			return fmt.Errorf("setting `contact`: %+v", err)
		}
		if err := d.Set("status", flattenOrderStatus(props.CurrentStatus)); err != nil {
			return fmt.Errorf("setting `status`: %+v", err)
		}
		if err := d.Set("shipment_address", flattenOrderAddress(props.ShippingAddress)); err != nil {
			return fmt.Errorf("setting `shipment_address`: %+v", err)
		}
		if err := d.Set("shipment_tracking", flattenOrderTrackingInfo(props.DeliveryTrackingInfo)); err != nil {
			return fmt.Errorf("setting `shipment_tracking`: %+v", err)
		}
		if err := d.Set("shipment_history", flattenOrderHistory(props.OrderHistory)); err != nil {
			return fmt.Errorf("setting `shipment_history`: %+v", err)
		}
		if err := d.Set("return_tracking", flattenOrderTrackingInfo(props.ReturnTrackingInfo)); err != nil {
			return fmt.Errorf("setting `return_tracking`: %+v", err)
		}
		d.Set("serial_number", props.SerialNumber)
	}

	return nil
}

func resourceOrderDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataboxEdge.OrderClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.OrderID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.DataBoxEdgeDeviceName, id.ResourceGroup)
	if err != nil {
		return fmt.Errorf("deleting %s: %v", *id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %v", *id, err)
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
		EmailList:     utils.ExpandStringSlice(v["emails"].(*pluginsdk.Set).List()),
	}
}

func expandOrderAddress(input []interface{}) *databoxedge.Address {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	var address1 string
	var address2 string
	var address3 string

	addressLines := v["address"].([]interface{})

	for i, addressLine := range addressLines {
		if addressLine != "" {
			switch i {
			case 0:
				address1 = addressLine.(string)
			case 1:
				address2 = addressLine.(string)
			case 3:
				address3 = addressLine.(string)
			}
		}
	}

	return &databoxedge.Address{
		AddressLine1: utils.String(address1),
		AddressLine2: utils.String(address2),
		AddressLine3: utils.String(address3),
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
			"info":               status,
			"comments":           comments,
			"additional_details": additionalOrderDetails,
			"last_update":        updateDateTime,
		},
	}
}

func flattenOrderAddress(input *databoxedge.Address) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
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

	address := make([]interface{}, 0)

	if input.AddressLine1 != nil {
		address = append(address, *input.AddressLine1)
	}

	if input.AddressLine2 != nil && *input.AddressLine2 != "" {
		address = append(address, *input.AddressLine2)
	}

	if input.AddressLine3 != nil && *input.AddressLine3 != "" {
		address = append(address, *input.AddressLine3)
	}

	return []interface{}{
		map[string]interface{}{
			"address":     address,
			"city":        city,
			"country":     country,
			"postal_code": postalCode,
			"state":       state,
		},
	}
}

func flattenOrderTrackingInfo(input *[]databoxedge.TrackingInfo) []interface{} {
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

func flattenOrderHistory(input *[]databoxedge.OrderStatus) []interface{} {
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
			"additional_details": additionalOrderDetails,
			"comments":           comments,
			"last_update":        updateDateTime,
		})
	}
	return results
}
