package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceArmNetworkConnectionMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetworkConnectionMonitorCreateUpdate,
		Read:   resourceArmNetworkConnectionMonitorRead,
		Update: resourceArmNetworkConnectionMonitorCreateUpdate,
		Delete: resourceArmNetworkConnectionMonitorDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"network_watcher_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocation(),

			"auto_start": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},

			"interval_in_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60,
				ValidateFunc: validation.IntAtLeast(30),
			},

			"source": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_machine_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validate.PortNumberOrZero,
						},
					},
				},
			},

			"destination": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_machine_id": {
							Type:          schema.TypeString,
							Optional:      true,
							ValidateFunc:  azure.ValidateResourceID,
							ConflictsWith: []string{"destination.0.address"},
						},
						"address": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"destination.0.virtual_machine_id"},
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validate.PortNumber,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmNetworkConnectionMonitorCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.ConnectionMonitorsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	watcherName := d.Get("network_watcher_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	autoStart := d.Get("auto_start").(bool)
	intervalInSeconds := int32(d.Get("interval_in_seconds").(int))

	source, err := expandArmNetworkConnectionMonitorSource(d)
	if err != nil {
		return err
	}

	dest, err := expandArmNetworkConnectionMonitorDestination(d)
	if err != nil {
		return err
	}

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, watcherName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Connection Monitor %q (Watcher %q / Resource Group %q): %s", name, watcherName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_network_connection_monitor", *existing.ID)
		}
	}

	t := d.Get("tags").(map[string]interface{})

	properties := network.ConnectionMonitor{
		Location: utils.String(location),
		Tags:     tags.Expand(t),
		ConnectionMonitorParameters: &network.ConnectionMonitorParameters{
			Source:                      source,
			Destination:                 dest,
			AutoStart:                   utils.Bool(autoStart),
			MonitoringIntervalInSeconds: utils.Int32(intervalInSeconds),
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, watcherName, name, properties)
	if err != nil {
		return fmt.Errorf("Error creating Connection Monitor %q (Watcher %q / Resource Group %q): %+v", name, watcherName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Connection Monitor %q (Watcher %q / Resource Group %q): %+v", name, watcherName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, watcherName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Connection Monitor %q (Watcher %q / Resource Group %q): %+v", name, watcherName, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Connection Monitor %q (Watcher %q / Resource Group %q) ID", name, watcherName, resourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmNetworkConnectionMonitorRead(d, meta)
}

func resourceArmNetworkConnectionMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.ConnectionMonitorsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	watcherName := id.Path["networkWatchers"]
	name := id.Path["NetworkConnectionMonitors"]

	resp, err := client.Get(ctx, resourceGroup, watcherName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Connection Monitor %q (Watcher %q / Resource Group %q) %+v", name, watcherName, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("network_watcher_name", watcherName)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.ConnectionMonitorResultProperties; props != nil {
		d.Set("auto_start", props.AutoStart)
		d.Set("interval_in_seconds", props.MonitoringIntervalInSeconds)

		source := flattenArmNetworkConnectionMonitorSource(props.Source)
		if err := d.Set("source", source); err != nil {
			return fmt.Errorf("Error setting `source`: %+v", err)
		}

		dest := flattenArmNetworkConnectionMonitorDestination(props.Destination)
		if err := d.Set("destination", dest); err != nil {
			return fmt.Errorf("Error setting `destination`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmNetworkConnectionMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.ConnectionMonitorsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	watcherName := id.Path["networkWatchers"]
	name := id.Path["NetworkConnectionMonitors"]

	future, err := client.Delete(ctx, resourceGroup, watcherName, name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting Connection Monitor %q (Watcher %q / Resource Group %q): %+v", name, watcherName, resourceGroup, err)
		}
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deletion of Connection Monitor %q (Watcher %q / Resource Group %q): %+v", name, watcherName, resourceGroup, err)
	}

	return nil
}

func flattenArmNetworkConnectionMonitorSource(input *network.ConnectionMonitorSource) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})

	if resourceID := input.ResourceID; resourceID != nil {
		output["virtual_machine_id"] = *resourceID
	}
	if port := input.Port; port != nil {
		output["port"] = *port
	}

	return []interface{}{output}
}

func expandArmNetworkConnectionMonitorSource(d *schema.ResourceData) (*network.ConnectionMonitorSource, error) {
	sources := d.Get("source").([]interface{})
	source := sources[0].(map[string]interface{})

	monitorSource := network.ConnectionMonitorSource{}
	if v := source["virtual_machine_id"]; v != "" {
		monitorSource.ResourceID = utils.String(v.(string))
	}
	if v := source["port"]; v != "" {
		monitorSource.Port = utils.Int32(int32(v.(int)))
	}

	return &monitorSource, nil
}

func flattenArmNetworkConnectionMonitorDestination(input *network.ConnectionMonitorDestination) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})

	// When monitoring a VM, the address field will contain the current address
	// of the VM. We only want to copy over the address field if the virtual
	// machine field is not set to avoid unwanted diffs.
	if resourceID := input.ResourceID; resourceID != nil {
		output["virtual_machine_id"] = *resourceID
	} else if address := input.Address; address != nil {
		output["address"] = *address
	}

	if port := input.Port; port != nil {
		output["port"] = *port
	}

	return []interface{}{output}
}

func expandArmNetworkConnectionMonitorDestination(d *schema.ResourceData) (*network.ConnectionMonitorDestination, error) {
	dests := d.Get("destination").([]interface{})
	dest := dests[0].(map[string]interface{})

	monitorDest := network.ConnectionMonitorDestination{}

	if v := dest["virtual_machine_id"]; v != "" {
		monitorDest.ResourceID = utils.String(v.(string))
	}
	if v := dest["address"]; v != "" {
		monitorDest.Address = utils.String(v.(string))
	}
	if v := dest["port"]; v != "" {
		monitorDest.Port = utils.Int32(int32(v.(int)))
	}

	if monitorDest.ResourceID == nil && monitorDest.Address == nil {
		return nil, fmt.Errorf("Error: either `destination.virtual_machine_id` or `destination.address` must be specified")
	}

	return &monitorDest, nil
}
