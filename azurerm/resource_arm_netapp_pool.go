package azurerm

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2019-06-01/netapp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	aznetapp "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/netapp"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNetAppPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetAppPoolCreateUpdate,
		Read:   resourceArmNetAppPoolRead,
		Update: resourceArmNetAppPoolCreateUpdate,
		Delete: resourceArmNetAppPoolDelete,

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
				ValidateFunc: aznetapp.ValidateNetAppPoolName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: aznetapp.ValidateNetAppAccountName,
			},

			"service_level": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(netapp.Premium),
					string(netapp.Standard),
					string(netapp.Ultra),
				}, true),
			},

			"size_in_tb": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(4, 500),
			},
		},
	}
}

func resourceArmNetAppPoolCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.PoolClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, accountName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing NetApp Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_netapp_pool", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	serviceLevel := d.Get("service_level").(string)
	sizeInTB := int64(d.Get("size_in_tb").(int))
	sizeInMB := sizeInTB * 1024 * 1024
	sizeInBytes := sizeInMB * 1024 * 1024

	capacityPoolParameters := netapp.CapacityPool{
		Location: utils.String(location),
		PoolProperties: &netapp.PoolProperties{
			ServiceLevel: netapp.ServiceLevel(serviceLevel),
			Size:         utils.Int64(sizeInBytes),
		},
	}

	future, err := client.CreateOrUpdate(ctx, capacityPoolParameters, resourceGroup, accountName, name)
	if err != nil {
		return fmt.Errorf("Error creating NetApp Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of NetApp Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving NetApp Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read NetApp Pool %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmNetAppPoolRead(d, meta)
}

func resourceArmNetAppPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.PoolClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	accountName := id.Path["netAppAccounts"]
	name := id.Path["capacityPools"]

	resp, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] NetApp Pools %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading NetApp Pools %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("account_name", accountName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if poolProperties := resp.PoolProperties; poolProperties != nil {
		d.Set("service_level", poolProperties.ServiceLevel)

		sizeInTB := int64(0)
		if poolProperties.Size != nil {
			sizeInBytes := *poolProperties.Size
			sizeInMB := sizeInBytes / 1024 / 1024
			sizeInTB = sizeInMB / 1024 / 1024
		}
		d.Set("size_in_tb", int(sizeInTB))
	}

	return nil
}

func resourceArmNetAppPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.PoolClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	accountName := id.Path["netAppAccounts"]
	name := id.Path["capacityPools"]

	_, err = client.Delete(ctx, resourceGroup, accountName, name)
	if err != nil {
		return fmt.Errorf("Error deleting NetApp Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return waitForNetAppPoolToBeDeleted(ctx, client, resourceGroup, accountName, name)
}

func waitForNetAppPoolToBeDeleted(ctx context.Context, client *netapp.PoolsClient, resourceGroup, accountName, name string) error {
	log.Printf("[DEBUG] Waiting for NetApp Pool Provisioning Service %q (Resource Group %q) to be deleted", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"200", "202"},
		Target:  []string{"404"},
		Refresh: netappPoolDeleteStateRefreshFunc(ctx, client, resourceGroup, accountName, name),
		Timeout: 20 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for NetApp Pool Provisioning Service %q (Resource Group %q) to be deleted: %+v", name, resourceGroup, err)
	}

	return nil
}

func netappPoolDeleteStateRefreshFunc(ctx context.Context, client *netapp.PoolsClient, resourceGroupName string, accountName string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, accountName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(res.Response) {
				return nil, "", fmt.Errorf("Error retrieving NetApp Pool %q (Resource Group %q): %s", name, resourceGroupName, err)
			}
		}

		// The resource NetApp Pool depends on the resource NetApp Account.
		// Although the delete API returns 404 which means the NetApp Pool resource has been deleted.
		// Then it tries to immediately delete NetApp Account but it still throws error `Can not delete resource before nested resources are deleted.`
		// In this case we're going to try triggering the Deletion again, in-case it didn't work prior to this attempt.
		// For more details, see related Bug: https://github.com/Azure/azure-sdk-for-go/issues/6374
		if _, err := client.Delete(ctx, resourceGroupName, accountName, name); err != nil {
			log.Printf("Error reissuing NetApp Pool %q delete request (Resource Group %q): %+v", name, resourceGroupName, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}
