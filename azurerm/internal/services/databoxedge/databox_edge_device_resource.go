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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databoxedge/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databoxedge/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
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

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
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

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"model_description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"sku_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataboxEdgeDeviceSkuName,
			},

			"configured_role_types": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"culture": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"device_hcs_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"device_local_capacity": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"device_model": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"device_status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"device_software_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"device_type": {
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
		DeviceProperties: &databoxedge.DeviceProperties{
			Description:      utils.String(d.Get("description").(string)),
			FriendlyName:     utils.String(d.Get("friendly_name").(string)),
			ModelDescription: utils.String(d.Get("model_description").(string)),
		},
		Sku:  expandDeviceSku(d.Get("sku_name").(string)),
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
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
		d.Set("device_status", props.DataBoxEdgeDeviceStatus)
		d.Set("description", props.Description)
		d.Set("friendly_name", props.FriendlyName)
		d.Set("model_description", props.ModelDescription)

		configuredRoleTypes := make([]string, 0)
		if props.ConfiguredRoleTypes != nil {
			for _, item := range *props.ConfiguredRoleTypes {
				configuredRoleTypes = append(configuredRoleTypes, (string)(item))
			}
		}

		d.Set("configured_role_types", utils.FlattenStringSlice(&configuredRoleTypes))
		d.Set("culture", props.Culture)
		d.Set("device_hcs_version", props.DeviceHcsVersion)
		d.Set("device_local_capacity", props.DeviceLocalCapacity)
		d.Set("device_model", props.DeviceModel)
		d.Set("device_software_version", props.DeviceSoftwareVersion)
		d.Set("device_type", props.DeviceType)
		d.Set("node_count", props.NodeCount)
		d.Set("serial_number", props.SerialNumber)
		d.Set("time_zone", props.TimeZone)
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

func flattenDeviceSku(input *databoxedge.Sku) *string {
	if input == nil {
		return nil
	}

	var name databoxedge.SkuName
	if input.Name != "" {
		name = input.Name
	}
	var tier databoxedge.SkuTier
	if input.Tier != "" {
		tier = input.Tier
	} else {
		tier = databoxedge.Standard
	}

	skuName := fmt.Sprintf("%s-%s", name, tier)

	return &skuName
}
