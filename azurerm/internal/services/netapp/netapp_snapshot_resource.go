package netapp

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2019-10-01/netapp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/netapp/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNetAppSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetAppSnapshotCreate,
		Read:   resourceArmNetAppSnapshotRead,
		Update: resourceArmNetAppSnapshotUpdate,
		Delete: resourceArmNetAppSnapshotDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SnapshotID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateNetAppSnapshotName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateNetAppAccountName,
			},

			"pool_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateNetAppPoolName,
			},

			"volume_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateNetAppVolumeName,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmNetAppSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.SnapshotClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)
	poolName := d.Get("pool_name").(string)
	volumeName := d.Get("volume_name").(string)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, accountName, poolName, volumeName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing NetApp Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_netapp_snapshot", *resp.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	parameters := netapp.Snapshot{
		Location: utils.String(location),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Create(ctx, parameters, resourceGroup, accountName, poolName, volumeName, name)
	if err != nil {
		return fmt.Errorf("Error creating NetApp Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of NetApp Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, accountName, poolName, volumeName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving NetApp Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read NetApp Snapshot %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmNetAppSnapshotRead(d, meta)
}

func resourceArmNetAppSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.SnapshotClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SnapshotID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.NetAppAccountName, id.CapacityPoolName, id.VolumeName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] NetApp Snapshots %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading NetApp Snapshots %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.NetAppAccountName)
	d.Set("pool_name", id.CapacityPoolName)
	d.Set("volume_name", id.VolumeName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmNetAppSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.SnapshotClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SnapshotID(d.Id())
	if err != nil {
		return err
	}

	parameters := netapp.SnapshotPatch{
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err = client.Update(ctx, parameters, id.ResourceGroup, id.NetAppAccountName, id.CapacityPoolName, id.VolumeName, id.Name); err != nil {
		return fmt.Errorf("Error updating NetApp Snapshot %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.NetAppAccountName, id.CapacityPoolName, id.VolumeName, id.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving NetApp Snapshot %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read NetApp Snapshot %q (Resource Group %q) ID", id.Name, id.ResourceGroup)
	}

	return resourceArmNetAppSnapshotRead(d, meta)
}

func resourceArmNetAppSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.SnapshotClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SnapshotID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, id.ResourceGroup, id.NetAppAccountName, id.CapacityPoolName, id.VolumeName, id.Name); err != nil {
		return fmt.Errorf("Error deleting NetApp Snapshot %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	// The resource NetApp Snapshot depends on the resource NetApp Volume.
	// Although the delete API returns 404 which means the NetApp Snapshot resource has been deleted.
	// Then it tries to immediately delete NetApp Volume but it still throws error `Can not delete resource before nested resources are deleted.`
	// In this case we're going to re-check status code again.
	// For more details, see related Bug: https://github.com/Azure/azure-sdk-for-go/issues/11475
	log.Printf("[DEBUG] Waiting for NetApp Snapshot %q (Resource Group %q) to be deleted", id.Name, id.ResourceGroup)
	stateConf := &resource.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"200", "202"},
		Target:                    []string{"204", "404"},
		Refresh:                   netappSnapshotDeleteStateRefreshFunc(ctx, client, id.ResourceGroup, id.NetAppAccountName, id.CapacityPoolName, id.VolumeName, id.Name),
		Timeout:                   d.Timeout(schema.TimeoutDelete),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for NetApp Snapshot %q (Resource Group %q) to be deleted: %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func netappSnapshotDeleteStateRefreshFunc(ctx context.Context, client *netapp.SnapshotsClient, resourceGroupName string, accountName string, poolName string, volumeName string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, accountName, poolName, volumeName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(res.Response) {
				return nil, "", fmt.Errorf("Error retrieving NetApp Snapshot %q (Resource Group %q): %s", name, resourceGroupName, err)
			}
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}
