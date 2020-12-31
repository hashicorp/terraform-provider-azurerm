package iothub

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/iothub/mgmt/2020-03-01/devices"
	"github.com/Azure/azure-sdk-for-go/services/provisioningservices/mgmt/2018-01-22/iothub"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceIotHubDPS() *schema.Resource {
	return &schema.Resource{
		Create: resourceIotHubDPSCreateUpdate,
		Read:   resourceIotHubDPSRead,
		Update: resourceIotHubDPSCreateUpdate,
		Delete: resourceIotHubDPSDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
				ValidateFunc: validate.IoTHubName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(), // azure.SchemaResourceGroupNameDiffSuppress(),

			"location": azure.SchemaLocation(),

			"sku": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(devices.S1),
							}, false),
						},

						"capacity": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 200),
						},
					},
				},
			},

			"linked_hub": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection_string": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							ForceNew:     true,
							// Azure returns the key as ****. We'll suppress that here.
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								secretKeyRegex := regexp.MustCompile("(SharedAccessKey)=[^;]+")
								maskedNew := secretKeyRegex.ReplaceAllString(new, "$1=****")
								return (new == d.Get(k).(string)) && (maskedNew == old)
							},
							Sensitive: true,
						},
						"location": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							StateFunc:    azure.NormalizeLocation,
							ForceNew:     true,
						},
						"apply_allocation_policy": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"allocation_weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
						"hostname": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"allocation_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"device_provisioning_host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"id_scope": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_operations_host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceIotHubDPSCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, name, resourceGroup)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing IoT Device Provisioning Service %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_iothub_dps", *existing.ID)
		}
	}

	iotdps := iothub.ProvisioningServiceDescription{
		Location: utils.String(d.Get("location").(string)),
		Name:     utils.String(name),
		Sku:      expandIoTHubDPSSku(d),
		Properties: &iothub.IotDpsPropertiesDescription{
			IotHubs: expandIoTHubDPSIoTHubs(d.Get("linked_hub").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, iotdps)
	if err != nil {
		return fmt.Errorf("Error creating/updating IoT Device Provisioning Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the completion of the creating/updating of IoT Device Provisioning Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, name, resourceGroup)
	if err != nil {
		return fmt.Errorf("Error retrieving IoT Device Provisioning Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read IoT Device Provisioning Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceIotHubDPSRead(d, meta)
}

func resourceIotHubDPSRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["provisioningServices"]
	// the name path can use the ProvisioningServices in older iterations
	if name == "" {
		name = id.Path["ProvisioningServices"]
	}

	resp, err := client.Get(ctx, name, resourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving IoT Device Provisioning Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	sku := flattenIoTHubDPSSku(resp.Sku)
	if err := d.Set("sku", sku); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}

	if props := resp.Properties; props != nil {
		if err := d.Set("linked_hub", flattenIoTHubDPSLinkedHub(props.IotHubs)); err != nil {
			return fmt.Errorf("Error setting `linked_hub`: %+v", err)
		}

		d.Set("service_operations_host_name", props.ServiceOperationsHostName)
		d.Set("device_provisioning_host_name", props.DeviceProvisioningHostName)
		d.Set("id_scope", props.IDScope)
		d.Set("allocation_policy", props.AllocationPolicy)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceIotHubDPSDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["provisioningServices"]

	future, err := client.Delete(ctx, name, resourceGroup)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting IoT Device Provisioning Service %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return waitForIotHubDPSToBeDeleted(ctx, client, resourceGroup, name, d)
}

func waitForIotHubDPSToBeDeleted(ctx context.Context, client *iothub.IotDpsResourceClient, resourceGroup, name string, d *schema.ResourceData) error {
	// we can't use the Waiter here since the API returns a 404 once it's deleted which is considered a polling status code..
	log.Printf("[DEBUG] Waiting for IoT Device Provisioning Service %q (Resource Group %q) to be deleted", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: iothubdpsStateStatusCodeRefreshFunc(ctx, client, resourceGroup, name),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for IoT Device Provisioning Service %q (Resource Group %q) to be deleted: %+v", name, resourceGroup, err)
	}

	return nil
}

func iothubdpsStateStatusCodeRefreshFunc(ctx context.Context, client *iothub.IotDpsResourceClient, resourceGroup, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, name, resourceGroup)

		log.Printf("Retrieving IoT Device Provisioning Service %q (Resource Group %q) returned Status %d", resourceGroup, name, res.StatusCode)

		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, strconv.Itoa(res.StatusCode), nil
			}
			return nil, "", fmt.Errorf("Error polling for the status of the IoT Device Provisioning Service %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}

func expandIoTHubDPSSku(d *schema.ResourceData) *iothub.IotDpsSkuInfo {
	skuList := d.Get("sku").([]interface{})
	skuMap := skuList[0].(map[string]interface{})

	return &iothub.IotDpsSkuInfo{
		Name:     iothub.IotDpsSku(skuMap["name"].(string)),
		Capacity: utils.Int64(int64(skuMap["capacity"].(int))),
	}
}

func expandIoTHubDPSIoTHubs(input []interface{}) *[]iothub.DefinitionDescription {
	linkedHubs := make([]iothub.DefinitionDescription, 0)

	for _, attr := range input {
		linkedHubConfig := attr.(map[string]interface{})
		linkedHub := iothub.DefinitionDescription{
			ConnectionString:      utils.String(linkedHubConfig["connection_string"].(string)),
			AllocationWeight:      utils.Int32(int32(linkedHubConfig["allocation_weight"].(int))),
			ApplyAllocationPolicy: utils.Bool(linkedHubConfig["apply_allocation_policy"].(bool)),
			Location:              utils.String(linkedHubConfig["location"].(string)),
		}

		linkedHubs = append(linkedHubs, linkedHub)
	}

	return &linkedHubs
}

func flattenIoTHubDPSSku(input *iothub.IotDpsSkuInfo) []interface{} {
	output := make(map[string]interface{})

	output["name"] = string(input.Name)
	if capacity := input.Capacity; capacity != nil {
		output["capacity"] = int(*capacity)
	}

	return []interface{}{output}
}

func flattenIoTHubDPSLinkedHub(input *[]iothub.DefinitionDescription) []interface{} {
	linkedHubs := make([]interface{}, 0)
	if input == nil {
		return linkedHubs
	}

	for _, attr := range *input {
		linkedHub := make(map[string]interface{})

		if attr.Name != nil {
			linkedHub["hostname"] = *attr.Name
		}
		if attr.ApplyAllocationPolicy != nil {
			linkedHub["apply_allocation_policy"] = *attr.ApplyAllocationPolicy
		}
		if attr.AllocationWeight != nil {
			linkedHub["allocation_weight"] = *attr.AllocationWeight
		}
		if attr.ConnectionString != nil {
			linkedHub["connection_string"] = *attr.ConnectionString
		}
		if attr.Location != nil {
			linkedHub["location"] = *attr.Location
		}

		linkedHubs = append(linkedHubs, linkedHub)
	}

	return linkedHubs
}
