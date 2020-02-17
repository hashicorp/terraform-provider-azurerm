package databox

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/databox/mgmt/2019-09-01/databox"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databox/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataBoxJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataBoxJobCreate,
		Read:   resourceArmDataBoxJobRead,
		Update: resourceArmDataBoxJobUpdate,
		Delete: resourceArmDataBoxJobDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ParseDataBoxJobID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateDataBoxJobName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"contact_details": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"contact_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: ValidateDataBoxJobContactName,
						},
						"emails": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 10,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: ValidateDataBoxJobEmail,
							},
						},
						"phone_number": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: ValidateDataBoxJobPhoneNumber,
						},
						"mobile": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"notification_preference": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							MinItems: 6,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"send_notification": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"stage_name": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(databox.AtAzureDC),
											string(databox.DataCopy),
											string(databox.Delivered),
											string(databox.DevicePrepared),
											string(databox.Dispatched),
											string(databox.PickedUp),
										}, false),
									},
								},
							},
						},
						"phone_extension": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: ValidateDataBoxJobPhoneExtension,
						},
					},
				},
			},

			"destination_account_details": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_destination_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(databox.DataDestinationTypeStorageAccount),
								string(databox.DataDestinationTypeManagedDisk),
							}, false),
						},
						"resource_group_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"share_password": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ForceNew:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"staging_storage_account_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"storage_account_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: azure.ValidateResourceID,
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
						"city": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: ValidateDataBoxJobCity,
						},
						"country": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"postal_code": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: ValidateDataBoxJobPostCode,
						},
						"state_or_province": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"street_address_1": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: ValidateDataBoxJobStreetAddress,
						},
						"address_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(databox.Commercial),
								string(databox.None),
								string(databox.Residential),
							}, false),
						},
						"company_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: ValidateDataBoxJobCompanyName,
						},
						"street_address_2": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: ValidateDataBoxJobStreetAddress,
						},
						"street_address_3": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: ValidateDataBoxJobStreetAddress,
						},
						"zip_extended_code": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: ValidateDataBoxJobPostCode,
						},
					},
				},
			},

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(databox.DataBox),
					string(databox.DataBoxDisk),
					string(databox.DataBoxHeavy),
				}, false),
			},

			"additional_preferred_disks_properties": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"delivery_scheduled_date_time": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validate.RFC3339Time,
			},

			"delivery_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(databox.NonScheduled),
					string(databox.Scheduled),
				}, false),
			},

			"device_password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"disk_pass_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: ValidateDataBoxJobDiskPassKey,
			},

			"expected_data_size_in_tb": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 1000000),
			},

			"preferences": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"preferred_shipment_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(databox.CustomerManaged),
								string(databox.MicrosoftManaged),
							}, false),
						},
						"preferred_data_center_region": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDataBoxJobCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBox.JobClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, "Details")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Data Box Job (Data Box Job Name %q / Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_box_job", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	contactDetails := d.Get("contact_details").([]interface{})
	deliveryType := d.Get("delivery_type").(string)
	destinationAccountDetails := d.Get("destination_account_details").([]interface{})
	devicePassword := d.Get("device_password").(string)
	diskPassKey := d.Get("disk_pass_key").(string)
	expectedDataSizeInTB := d.Get("expected_data_size_in_tb").(int)
	preferences := d.Get("preferences").([]interface{})
	preferredDisks := d.Get("additional_preferred_disks_properties").(map[string]interface{})
	shippingAddress := d.Get("shipping_address").([]interface{})
	skuName := d.Get("sku_name").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := databox.JobResource{
		Location: utils.String(location),
		JobProperties: &databox.JobProperties{
			DeliveryType: databox.JobDeliveryType(deliveryType),
		},
		Sku: &databox.Sku{
			Name: databox.SkuName(skuName),
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("delivery_scheduled_date_time"); ok {
		t, _ := time.Parse(time.RFC3339, v.(string))
		parameters.JobProperties.DeliveryInfo = &databox.JobDeliveryInfo{
			ScheduledDateTime: &date.Time{Time: t},
		}
	}

	switch skuName {
	case string(databox.DataBox):
		parameters.JobProperties.Details = &databox.JobDetailsType{
			ContactDetails:              expandArmDataBoxJobContactDetails(contactDetails),
			DestinationAccountDetails:   expandArmDataBoxJobDestinationAccountDetails(destinationAccountDetails),
			DevicePassword:              utils.String(devicePassword),
			ExpectedDataSizeInTerabytes: utils.Int32(int32(expectedDataSizeInTB)),
			Preferences:                 expandArmDataBoxJobPreferences(preferences),
			ShippingAddress:             expandArmDataBoxJobShippingAddress(shippingAddress),
		}
	case string(databox.DataBoxDisk):
		parameters.JobProperties.Details = &databox.DiskJobDetails{
			ContactDetails:              expandArmDataBoxJobContactDetails(contactDetails),
			DestinationAccountDetails:   expandArmDataBoxJobDestinationAccountDetails(destinationAccountDetails),
			ExpectedDataSizeInTerabytes: utils.Int32(int32(expectedDataSizeInTB)),
			Passkey:                     utils.String(diskPassKey),
			Preferences:                 expandArmDataBoxJobPreferences(preferences),
			PreferredDisks:              expandArmDataBoxJobPreferredDisks(preferredDisks),
			ShippingAddress:             expandArmDataBoxJobShippingAddress(shippingAddress),
		}
	case string(databox.DataBoxHeavy):
		parameters.JobProperties.Details = &databox.HeavyJobDetails{
			ContactDetails:              expandArmDataBoxJobContactDetails(contactDetails),
			DevicePassword:              utils.String(devicePassword),
			DestinationAccountDetails:   expandArmDataBoxJobDestinationAccountDetails(destinationAccountDetails),
			ExpectedDataSizeInTerabytes: utils.Int32(int32(expectedDataSizeInTB)),
			Preferences:                 expandArmDataBoxJobPreferences(preferences),
			ShippingAddress:             expandArmDataBoxJobShippingAddress(shippingAddress),
		}
	}

	future, err := client.Create(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Data Box Job (Data Box Job Name %q / Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Data Box Job (Data Box Job Name %q / Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name, "Details")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Box Job (Data Box Job Name %q / Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read Data Box Job (Data Box Job Name %q / Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmDataBoxJobRead(d, meta)
}

func resourceArmDataBoxJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBox.JobClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseDataBoxJobID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "Details")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Data Box Job %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Data Box Job (Data Box Job Name %q / Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
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

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDataBoxJobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBox.JobClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Data Box Job update.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	contactDetails := d.Get("contact_details").([]interface{})
	shippingAddress := d.Get("shipping_address").([]interface{})
	t := d.Get("tags").(map[string]interface{})

	parameters := databox.JobResourceUpdateParameter{
		UpdateJobProperties: &databox.UpdateJobProperties{
			Details: &databox.UpdateJobDetails{
				ContactDetails:  expandArmDataBoxJobContactDetails(contactDetails),
				ShippingAddress: expandArmDataBoxJobShippingAddress(shippingAddress),
			},
		},
		Tags: tags.Expand(t),
	}

	future, err := client.Update(ctx, resourceGroup, name, parameters, "")
	if err != nil {
		return fmt.Errorf("Error updating Data Box Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of Data Box Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name, "Details")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Box Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read Data Box Job %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmDataBoxJobRead(d, meta)
}

func resourceArmDataBoxJobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBox.JobClient
	lockClient := meta.(*clients.Client).Resource.LocksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseDataBoxJobID(d.Id())
	if err != nil {
		return err
	}

	reason := &databox.CancellationReason{
		Reason: utils.String("Cancel the order for deleting"),
	}

	resp, err := client.Cancel(ctx, id.ResourceGroup, id.Name, *reason)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error cancelling Order (Data Box Job Name %q / Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Data Box Job (Data Box Job Name %q / Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Data Box Job (Data Box Job Name %q / Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	destinationAccountDetails := d.Get("destination_account_details").([]interface{})
	for _, item := range destinationAccountDetails {
		if item != nil {
			v := item.(map[string]interface{})
			dataDestinationType := v["data_destination_type"].(string)
			lockName := "DATABOX_SERVICE"
			scope := ""

			switch dataDestinationType {
			case string(databox.DataDestinationTypeStorageAccount):
				scope = v["storage_account_id"].(string)
			case string(databox.DataDestinationTypeManagedDisk):
				scope = v["staging_storage_account_id"].(string)
			}

			if scope != "" {
				resp, err := lockClient.DeleteByScope(ctx, scope, lockName)
				if err != nil {
					if utils.ResponseWasNotFound(resp) {
						return nil
					}

					return fmt.Errorf("Error deleting Management Lock (Lock Name %q / Scope %q): %+v", lockName, scope, err)
				}
			}
		}
	}

	return nil
}

func expandArmDataBoxJobShippingAddress(input []interface{}) *databox.ShippingAddress {
	v := input[0].(map[string]interface{})

	return &databox.ShippingAddress{
		AddressType:     databox.AddressType(v["address_type"].(string)),
		City:            utils.String(v["city"].(string)),
		CompanyName:     utils.String(v["company_name"].(string)),
		Country:         utils.String(v["country"].(string)),
		PostalCode:      utils.String(v["postal_code"].(string)),
		StateOrProvince: utils.String(v["state_or_province"].(string)),
		StreetAddress1:  utils.String(v["street_address_1"].(string)),
		StreetAddress2:  utils.String(v["street_address_2"].(string)),
		StreetAddress3:  utils.String(v["street_address_3"].(string)),
		ZipExtendedCode: utils.String(v["zip_extended_code"].(string)),
	}
}

func expandArmDataBoxJobPreferences(input []interface{}) *databox.Preferences {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &databox.Preferences{
		TransportPreferences: &databox.TransportPreferences{
			PreferredShipmentType: databox.TransportShipmentTypes(v["preferred_shipment_type"].(string)),
		},
		PreferredDataCenterRegion: utils.ExpandStringSlice(v["preferred_data_center_region"].(*schema.Set).List()),
	}
}

func expandArmDataBoxJobContactDetails(input []interface{}) *databox.ContactDetails {
	v := input[0].(map[string]interface{})

	return &databox.ContactDetails{
		ContactName:            utils.String(v["contact_name"].(string)),
		EmailList:              utils.ExpandStringSlice(v["emails"].(*schema.Set).List()),
		Mobile:                 utils.String(v["mobile"].(string)),
		NotificationPreference: expandArmDataBoxJobNotificationPreference(v["notification_preference"].(*schema.Set).List()),
		Phone:                  utils.String(v["phone_number"].(string)),
		PhoneExtension:         utils.String(v["phone_extension"].(string)),
	}
}

func expandArmDataBoxJobDestinationAccountDetails(input []interface{}) *[]databox.BasicDestinationAccountDetails {
	results := make([]databox.BasicDestinationAccountDetails, 0)
	for _, item := range input {
		if item != nil {
			v := item.(map[string]interface{})
			dataDestinationType := v["data_destination_type"].(string)

			switch dataDestinationType {
			case string(databox.DataDestinationTypeStorageAccount):
				result := &databox.DestinationStorageAccountDetails{
					DataDestinationType: databox.DataDestinationTypeBasicDestinationAccountDetails(dataDestinationType),
					SharePassword:       utils.String(v["share_password"].(string)),
					StorageAccountID:    utils.String(v["storage_account_id"].(string)),
				}
				results = append(results, result)
			case string(databox.DataDestinationTypeManagedDisk):
				result := &databox.DestinationManagedDiskDetails{
					DataDestinationType:     databox.DataDestinationTypeBasicDestinationAccountDetails(dataDestinationType),
					ResourceGroupID:         utils.String(v["resource_group_id"].(string)),
					SharePassword:           utils.String(v["share_password"].(string)),
					StagingStorageAccountID: utils.String(v["staging_storage_account_id"].(string)),
				}
				results = append(results, result)
			}
		}
	}

	return &results
}

func expandArmDataBoxJobNotificationPreference(input []interface{}) *[]databox.NotificationPreference {
	results := make([]databox.NotificationPreference, 0)
	for _, item := range input {
		if item != nil {
			v := item.(map[string]interface{})
			sendNotification := v["send_notification"].(bool)
			stageName := v["stage_name"].(string)

			result := databox.NotificationPreference{
				SendNotification: utils.Bool(sendNotification),
				StageName:        databox.NotificationStageName(stageName),
			}

			results = append(results, result)
		}
	}

	return &results
}

func expandArmDataBoxJobPreferredDisks(input map[string]interface{}) map[string]*int32 {
	results := make(map[string]*int32)

	for k, v := range input {
		results[k] = utils.Int32(int32(v.(int)))
	}

	return results
}

func flattenArmDataBoxJobShippingAddress(input *databox.ShippingAddress) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	city := ""
	if input.City != nil {
		city = *input.City
	}

	companyName := ""
	if input.CompanyName != nil {
		companyName = *input.CompanyName
	}

	country := ""
	if input.Country != nil {
		country = *input.Country
	}

	postalCode := ""
	if input.PostalCode != nil {
		postalCode = *input.PostalCode
	}

	stateOrProvince := ""
	if input.StateOrProvince != nil {
		stateOrProvince = *input.StateOrProvince
	}

	streetAddress1 := ""
	if input.StreetAddress1 != nil {
		streetAddress1 = *input.StreetAddress1
	}

	streetAddress2 := ""
	if input.StreetAddress2 != nil {
		streetAddress2 = *input.StreetAddress2
	}

	streetAddress3 := ""
	if input.StreetAddress3 != nil {
		streetAddress3 = *input.StreetAddress3
	}

	zipExtendedCode := ""
	if input.ZipExtendedCode != nil {
		zipExtendedCode = *input.ZipExtendedCode
	}

	results = append(results, map[string]interface{}{
		"address_type":      input.AddressType,
		"city":              city,
		"company_name":      companyName,
		"country":           country,
		"postal_code":       postalCode,
		"state_or_province": stateOrProvince,
		"street_address_1":  streetAddress1,
		"street_address_2":  streetAddress2,
		"street_address_3":  streetAddress3,
		"zip_extended_code": zipExtendedCode,
	})

	return results
}

func flattenArmDataBoxJobPreferredDisks(input map[string]*int32) map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range input {
		output[k] = v
	}

	return output
}

func flattenArmDataBoxJobDestinationAccountDetails(input *[]databox.BasicDestinationAccountDetails) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if item != nil {
			if v, ok := item.AsDestinationStorageAccountDetails(); ok && v != nil {
				dataDestinationType := v.DataDestinationType

				sharePassword := ""
				if v.SharePassword != nil {
					sharePassword = *v.SharePassword
				}

				storageAccountID := ""
				if v.StorageAccountID != nil {
					storageAccountID = *v.StorageAccountID
				}

				results = append(results, map[string]interface{}{
					"data_destination_type": dataDestinationType,
					"share_password":        sharePassword,
					"storage_account_id":    storageAccountID,
				})
			} else if v, ok := item.AsDestinationManagedDiskDetails(); ok && v != nil {
				dataDestinationType := v.DataDestinationType

				resourceGroupID := ""
				if v.ResourceGroupID != nil {
					resourceGroupID = *v.ResourceGroupID
				}

				sharePassword := ""
				if v.SharePassword != nil {
					sharePassword = *v.SharePassword
				}

				stagingStorageAccountID := ""
				if v.StagingStorageAccountID != nil {
					stagingStorageAccountID = *v.StagingStorageAccountID
				}

				results = append(results, map[string]interface{}{
					"data_destination_type":      dataDestinationType,
					"resource_group_id":          resourceGroupID,
					"share_password":             sharePassword,
					"staging_storage_account_id": stagingStorageAccountID,
				})
			}
		}
	}

	return results
}

func flattenArmDataBoxJobContactDetails(input *databox.ContactDetails) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	contactName := ""
	if input.ContactName != nil {
		contactName = *input.ContactName
	}

	emails := []string{}
	if v := input.EmailList; v != nil {
		emails = *v
	}

	mobile := ""
	if input.Mobile != nil {
		mobile = *input.Mobile
	}

	phoneExtension := ""
	if input.PhoneExtension != nil {
		phoneExtension = *input.PhoneExtension
	}

	phoneNumber := ""
	if input.Phone != nil {
		phoneNumber = *input.Phone
	}

	results = append(results, map[string]interface{}{
		"contact_name":            contactName,
		"emails":                  utils.FlattenStringSlice(&emails),
		"mobile":                  mobile,
		"notification_preference": flattenArmDataBoxJobNotificationPreference(input.NotificationPreference),
		"phone_extension":         phoneExtension,
		"phone_number":            phoneNumber,
	})

	return results
}

func flattenArmDataBoxJobNotificationPreference(input *[]databox.NotificationPreference) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		sendNotification := false
		if item.SendNotification != nil {
			sendNotification = *item.SendNotification
		}

		results = append(results, map[string]interface{}{
			"send_notification": sendNotification,
			"stage_name":        item.StageName,
		})
	}

	return results
}

func flattenArmDataBoxJobPreferences(input *databox.Preferences) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	preferredDataCenterRegion := []string{}
	if v := input.PreferredDataCenterRegion; v != nil {
		preferredDataCenterRegion = *v
	}

	preferredShipmentType := input.TransportPreferences.PreferredShipmentType

	results = append(results, map[string]interface{}{
		"preferred_data_center_region": utils.FlattenStringSlice(&preferredDataCenterRegion),
		"preferred_shipment_type":      preferredShipmentType,
	})

	return results
}
