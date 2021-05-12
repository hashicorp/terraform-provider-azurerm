package databoxedge

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/databoxedge/mgmt/2020-12-01/databoxedge"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databoxedge/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databoxedge/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataboxEdgeDevice() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataboxEdgeDeviceCreate,
		Read:   resourceDataboxEdgeDeviceRead,
		Update: resourceDataboxEdgeDeviceUpdate,
		Delete: resourceDataboxEdgeDeviceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DataboxEdgeDeviceID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataboxEdgeName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"sku_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataboxEdgeDeviceSkuName,
			},

			"device_properties": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configured_role_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"culture": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"hcs_version": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"model": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"software_version": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"serial_number": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"time_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceDataboxEdgeDeviceCreate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataboxEdge.DeviceClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, name, resourceGroup)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Databox Edge Device %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_databox_edge_device", *existing.ID)
	}

	dataBoxEdgeDevice := databoxedge.Device{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Sku:      expandDeviceSku(d.Get("sku_name").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	future, err := client.CreateOrUpdate(ctx, name, dataBoxEdgeDevice, resourceGroup)
	if err != nil {
		return fmt.Errorf("creating Databox Edge Device %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating future for Databox Edge Device %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, name, resourceGroup)
	if err != nil {
		return fmt.Errorf("retrieving Databox Edge Device %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Databox Edge Device %q (Resource Group %q) ID", name, resourceGroup)
	}

	id, err := parse.DataboxEdgeDeviceID(*resp.ID)
	if err != nil {
		return err
	}

	d.SetId(id.ID(subscriptionId))

	return resourceDataboxEdgeDeviceRead(d, meta)
}

func resourceDataboxEdgeDeviceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataboxEdge.DeviceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataboxEdgeDeviceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.Name, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] databoxedge %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Databox Edge Device %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.DeviceProperties; props != nil {
		if err := d.Set("device_properties", flattenDeviceProperties(props)); err != nil {
			return fmt.Errorf("flattening 'device_properties' Databox Edge Device %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	if err := d.Set("sku_name", flattenDeviceSku(resp.Sku)); err != nil {
		return fmt.Errorf("setting `sku_name`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceDataboxEdgeDeviceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataboxEdge.DeviceClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataboxEdgeDeviceID(d.Id())
	if err != nil {
		return err
	}

	parameters := databoxedge.DevicePatch{}
	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.Name, parameters, id.ResourceGroup); err != nil {
		return fmt.Errorf("updating Databox Edge Device %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return resourceDataboxEdgeDeviceRead(d, meta)
}

func resourceDataboxEdgeDeviceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataboxEdge.DeviceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataboxEdgeDeviceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.Name, id.ResourceGroup)
	if err != nil {
		return fmt.Errorf("deleting Databox Edge Device %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on deleting future for Databox Edge Device %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return nil
}

func expandDeviceSku(input string) *databoxedge.Sku {
	if len(input) == 0 {
		return nil
	}

	v, err := parse.DataboxEdgeDeviceSkuName(input)
	if err != nil {
		return nil
	}

	return &databoxedge.Sku{
		Name: databoxedge.SkuName(v.Name),
		Tier: databoxedge.SkuTier(v.Tier),
	}
}

func flattenDeviceProperties(input *databoxedge.DeviceProperties) *[]interface{} {
	output := make([]interface{}, 0)
	configuredRoleTypes := make([]string, 0)

	var status string
	var culture string
	var hcsVersion string
	var capacity int64
	var model string
	var softwareVersion string
	var deviceType string
	var nodeCount int32
	var serialNumber string
	var timeZone string

	if input != nil {
		if input.ConfiguredRoleTypes != nil {
			for _, item := range *input.ConfiguredRoleTypes {
				configuredRoleTypes = append(configuredRoleTypes, (string)(item))
			}
		}

		if input.DataBoxEdgeDeviceStatus != "" {
			status = string(input.DataBoxEdgeDeviceStatus)
		}

		if input.Culture != nil {
			culture = *input.Culture
		}

		if input.DeviceHcsVersion != nil {
			hcsVersion = *input.DeviceHcsVersion
		}

		if input.DeviceLocalCapacity != nil {
			capacity = *input.DeviceLocalCapacity
		}

		if input.DeviceModel != nil {
			model = *input.DeviceModel
		}

		if input.DeviceSoftwareVersion != nil {
			softwareVersion = *input.DeviceSoftwareVersion
		}

		if input.DeviceType != "" {
			deviceType = string(input.DeviceType)
		}

		if input.NodeCount != nil {
			nodeCount = *input.NodeCount
		}

		if input.SerialNumber != nil {
			serialNumber = *input.SerialNumber
		}

		if input.TimeZone != nil {
			timeZone = *input.TimeZone
		}
	}

	output = append(output, map[string]interface{}{
		"configured_role_types": utils.FlattenStringSlice(&configuredRoleTypes),
		"culture":               culture,
		"hcs_version":           hcsVersion,
		"capacity":              capacity,
		"model":                 model,
		"status":                status,
		"software_version":      softwareVersion,
		"type":                  deviceType,
		"node_count":            nodeCount,
		"serial_number":         serialNumber,
		"time_zone":             timeZone,
	})

	return &output
}

func flattenDeviceSku(input *databoxedge.Sku) *string {
	if input == nil {
		return nil
	}

	var name databoxedge.SkuName
	var tier databoxedge.SkuTier

	if input.Name != "" {
		name = input.Name
	}

	if input.Tier != "" {
		tier = input.Tier
	} else {
		tier = databoxedge.Standard
	}

	skuName := fmt.Sprintf("%s-%s", name, tier)

	return &skuName
}
