package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNetworkWatcher() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetworkWatcherCreateUpdate,
		Read:   resourceArmNetworkWatcherRead,
		Update: resourceArmNetworkWatcherCreateUpdate,
		Delete: resourceArmNetworkWatcherDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"tags": tagsSchema(),
		},
	}
}

func resourceArmNetworkWatcherCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).watcherClient

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := d.Get("location").(string)
	tags := d.Get("tags").(map[string]interface{})

	watcher := network.Watcher{
		Location: utils.String(location),
		Tags:     expandTags(tags),
	}
	_, err := client.CreateOrUpdate(resourceGroup, name, watcher)
	if err != nil {
		return err
	}

	read, err := client.Get(resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Network Watcher %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmNetworkWatcherRead(d, meta)
}

func resourceArmNetworkWatcherRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).watcherClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["networkWatchers"]

	resp, err := client.Get(resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Network Watcher %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmNetworkWatcherDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).watcherClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["networkWatchers"]

	deleteResp, deleteErr := client.Delete(resourceGroup, name, make(chan struct{}))
	resp := <-deleteResp
	err = <-deleteErr

	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting Network Watcher %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}
