// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package databoxedge

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databoxedge/2022-03-01/orders"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databoxedge/migration"
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

		DeprecationMessage: `Creating DataBox Edge Orders are not supported via the Azure API - as such the 'azurerm_databox_edge_order' resource is deprecated and will be removed in v4.0 of the AzureRM Provider`,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := orders.ParseDataBoxEdgeDeviceID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.DataBoxEdgeOrderV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"device_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataboxEdgeName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

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
	client := meta.(*clients.Client).DataboxEdge.OrdersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := orders.NewDataBoxEdgeDeviceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("device_name").(string)) // TODO: state migration
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_databox_edge_order", id.ID())
		}
	}

	order := orders.Order{
		Properties: &orders.OrderProperties{
			ContactInformation: expandOrderContactDetails(d.Get("contact").([]interface{})),
			ShippingAddress:    expandOrderAddress(d.Get("shipment_address").([]interface{})),
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, order); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceOrderRead(d, meta)
}

func resourceOrderRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataboxEdge.OrdersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := orders.ParseDataBoxEdgeDeviceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", "default") // only one possible value
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("device_name", id.DataBoxEdgeDeviceName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if err := d.Set("contact", flattenOrderContactDetails(props.ContactInformation)); err != nil {
				return fmt.Errorf("setting `contact`: %+v", err)
			}
			currentStatus, err := flattenOrderStatus(props.CurrentStatus)
			if err != nil {
				return fmt.Errorf("flattening `status`: %+v", err)
			}
			if err := d.Set("status", currentStatus); err != nil {
				return fmt.Errorf("setting `status`: %+v", err)
			}
			if err := d.Set("shipment_address", flattenOrderAddress(props.ShippingAddress)); err != nil {
				return fmt.Errorf("setting `shipment_address`: %+v", err)
			}
			if err := d.Set("shipment_tracking", flattenOrderTrackingInfo(props.DeliveryTrackingInfo)); err != nil {
				return fmt.Errorf("setting `shipment_tracking`: %+v", err)
			}
			shipmentHistory, err := flattenOrderHistory(props.OrderHistory)
			if err != nil {
				return fmt.Errorf("flattening `shipment_history`: %+v", err)
			}
			if err := d.Set("shipment_history", shipmentHistory); err != nil {
				return fmt.Errorf("setting `shipment_history`: %+v", err)
			}
			if err := d.Set("return_tracking", flattenOrderTrackingInfo(props.ReturnTrackingInfo)); err != nil {
				return fmt.Errorf("setting `return_tracking`: %+v", err)
			}
			d.Set("serial_number", props.SerialNumber)
		}
	}

	return nil
}

func resourceOrderDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataboxEdge.OrdersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := orders.ParseDataBoxEdgeDeviceID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %v", *id, err)
	}

	return nil
}

func expandOrderContactDetails(input []interface{}) orders.ContactDetails {
	v := input[0].(map[string]interface{})
	emailList := make([]string, 0)

	for _, val := range v["emails"].(*pluginsdk.Set).List() {
		emailList = append(emailList, val.(string))
	}

	return orders.ContactDetails{
		ContactPerson: v["name"].(string),
		CompanyName:   v["company_name"].(string),
		Phone:         v["phone_number"].(string),
		EmailList:     emailList,
	}
}

func expandOrderAddress(input []interface{}) *orders.Address {
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

	return &orders.Address{
		AddressLine1: utils.String(address1),
		AddressLine2: utils.String(address2),
		AddressLine3: utils.String(address3),
		PostalCode:   utils.String(v["postal_code"].(string)),
		City:         utils.String(v["city"].(string)),
		State:        utils.String(v["state"].(string)),
		Country:      v["country"].(string),
	}
}

func flattenOrderContactDetails(input orders.ContactDetails) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"company_name": input.CompanyName,
			"name":         input.ContactPerson,
			"emails":       input.EmailList,
			"phone_number": input.Phone,
		},
	}
}

func flattenOrderStatus(input *orders.OrderStatus) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
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
		for k, v := range *input.AdditionalOrderDetails {
			additionalOrderDetails[k] = v
		}
	}
	var updateDateTime string
	d, err := input.GetUpdateDateTimeAsTime()
	if err != nil {
		return nil, fmt.Errorf("parsing UpdateDateTime: %+v", err)
	}
	if d != nil {
		updateDateTime = d.Format(time.RFC3339)
	}
	return &[]interface{}{
		map[string]interface{}{
			"info":               status,
			"comments":           comments,
			"additional_details": additionalOrderDetails,
			"last_update":        updateDateTime,
		},
	}, nil
}

func flattenOrderAddress(input *orders.Address) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
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
			"city":        input.City,
			"country":     input.Country,
			"postal_code": postalCode,
			"state":       state,
		},
	}
}

func flattenOrderTrackingInfo(input *[]orders.TrackingInfo) []interface{} {
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
		if item.TrackingId != nil {
			trackingId = *item.TrackingId
		}
		var trackingUrl string
		if item.TrackingUrl != nil {
			trackingUrl = *item.TrackingUrl
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

func flattenOrderHistory(input *[]orders.OrderStatus) (*[]interface{}, error) {
	results := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			additionalOrderDetails := make(map[string]interface{})
			if item.AdditionalOrderDetails != nil {
				for k, v := range *item.AdditionalOrderDetails {
					additionalOrderDetails[k] = v
				}
			}
			var comments string
			if item.Comments != nil {
				comments = *item.Comments
			}
			var updateDateTime string
			d, err := item.GetUpdateDateTimeAsTime()
			if err != nil {
				return nil, fmt.Errorf("parsing UpdateDateTime: %+v", err)
			}
			if d != nil {
				updateDateTime = d.Format(time.RFC3339)
			}
			results = append(results, map[string]interface{}{
				"additional_details": additionalOrderDetails,
				"comments":           comments,
				"last_update":        updateDateTime,
			})
		}
	}

	return &results, nil
}
