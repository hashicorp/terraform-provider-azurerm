package databoxedge

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/databoxedge/mgmt/2020-12-01/databoxedge"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databoxedge/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databoxedge/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDevice() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDeviceCreate,
		Read:   resourceDeviceRead,
		Update: resourceDeviceUpdate,
		Delete: resourceDeviceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DeviceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataboxEdgeName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"sku_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataboxEdgeDeviceSkuName,
			},

			"device_properties": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"configured_role_types": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"culture": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"hcs_version": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"capacity": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"model": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"status": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"software_version": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"node_count": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"serial_number": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"time_zone": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceDeviceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataboxEdge.DeviceClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewDeviceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	// sdk method is Get(ctx context.Context, deviceName string, resourceGroupName string)
	existing, err := client.Get(ctx, id.DataBoxEdgeDeviceName, id.ResourceGroup)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_databox_edge_device", id.ID())
	}

	dataBoxEdgeDevice := databoxedge.Device{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Sku:      expandDeviceSku(d.Get("sku_name").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	future, err := client.CreateOrUpdate(ctx, id.DataBoxEdgeDeviceName, dataBoxEdgeDevice, id.ResourceGroup)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDeviceRead(d, meta)
}

func resourceDeviceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataboxEdge.DeviceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DeviceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.DataBoxEdgeDeviceName, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.DataBoxEdgeDeviceName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.DeviceProperties; props != nil {
		if err := d.Set("device_properties", flattenDeviceProperties(props)); err != nil {
			return fmt.Errorf("flattening 'device_properties': %+v", err)
		}
	}

	if err := d.Set("sku_name", flattenDeviceSku(resp.Sku)); err != nil {
		return fmt.Errorf("setting `sku_name`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceDeviceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataboxEdge.DeviceClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DeviceID(d.Id())
	if err != nil {
		return err
	}

	parameters := databoxedge.DevicePatch{}
	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.DataBoxEdgeDeviceName, parameters, id.ResourceGroup); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceDeviceRead(d, meta)
}

func resourceDeviceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataboxEdge.DeviceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DeviceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.DataBoxEdgeDeviceName, id.ResourceGroup)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
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
