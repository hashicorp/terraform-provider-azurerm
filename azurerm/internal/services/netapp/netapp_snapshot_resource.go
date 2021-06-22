package netapp

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2020-09-01/netapp"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/netapp/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/netapp/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceNetAppSnapshot() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetAppSnapshotCreate,
		Read:   resourceNetAppSnapshotRead,
		Update: resourceNetAppSnapshotUpdate,
		Delete: resourceNetAppSnapshotDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SnapshotID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SnapshotName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName,
			},

			"pool_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PoolName,
			},

			"volume_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.VolumeName,
			},

			// TODO: remove this in a next breaking changes release since tags are
			// not supported anymore on Snapshots (todo 3.0)
			"tags": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
				Deprecated: "This property as been deprecated as the API no longer supports tags and will be removed in version 3.0 of the provider.",
			},
		},
	}
}

func resourceNetAppSnapshotCreate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	if tags.Expand(d.Get("tags").(map[string]interface{})) != nil {
		log.Printf("[WARN] Tags are not supported on snaphots anymore, ignoring values.")
	}

	parameters := netapp.Snapshot{
		Location: utils.String(location),
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

	return resourceNetAppSnapshotRead(d, meta)
}

func resourceNetAppSnapshotRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

	return nil
}

func resourceNetAppSnapshotUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	// Snapshot resource in Azure changed its type to proxied resource, therefore
	// tags are not supported anymore, ignoring any tags.
	if tags.Expand(d.Get("tags").(map[string]interface{})) == nil {
		log.Printf("[WARN] Tags are not supported on snaphots anymore, no update will happen in a snapshot at this time.")
	}

	return resourceNetAppSnapshotRead(d, meta)
}

func resourceNetAppSnapshotDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"200", "202"},
		Target:                    []string{"204", "404"},
		Refresh:                   netappSnapshotDeleteStateRefreshFunc(ctx, client, id.ResourceGroup, id.NetAppAccountName, id.CapacityPoolName, id.VolumeName, id.Name),
		Timeout:                   d.Timeout(pluginsdk.TimeoutDelete),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("Error waiting for NetApp Snapshot %q (Resource Group %q) to be deleted: %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func netappSnapshotDeleteStateRefreshFunc(ctx context.Context, client *netapp.SnapshotsClient, resourceGroupName string, accountName string, poolName string, volumeName string, name string) pluginsdk.StateRefreshFunc {
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
