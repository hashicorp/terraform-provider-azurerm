package storage

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmStorageContainer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmStorageContainerRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"storage_account_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"container_access_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"metadata": MetaDataComputedSchema(),

			// TODO: support for ACL's, Legal Holds and Immutability Policies
			"has_extended_immutability_policy": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"has_legal_hold": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"resource_manager_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmStorageContainerRead(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	mgmtContainerClient := storageClient.BlobContainersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	containerName := d.Get("name").(string)
	accountName := d.Get("storage_account_name").(string)

	account, err := storageClient.FindAccount(ctx, accountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Container %q: %s", accountName, containerName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Account %q for Storage Container %q", accountName, containerName)
	}

	client, err := storageClient.ContainersClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building Containers Client for Storage Account %q (Resource Group %q): %s", accountName, account.ResourceGroup, err)
	}

	d.SetId(client.GetResourceID(accountName, containerName))

	resp, err := mgmtContainerClient.Get(ctx, account.ResourceGroup, accountName, containerName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("container %q was not found in Account %q / Resource Group %q", containerName, accountName, account.ResourceGroup)
		}

		return fmt.Errorf("retrieving Container %q (Account %q / Resource Group %q): %s", containerName, accountName, account.ResourceGroup, err)
	}

	d.Set("name", containerName)

	d.Set("storage_account_name", accountName)

	if props := resp.ContainerProperties; props != nil {
		d.Set("container_access_type", flattenStorageContainerAccessLevel(props.PublicAccess))

		if err := d.Set("metadata", FlattenMetaDataPtr(props.Metadata)); err != nil {
			return fmt.Errorf("setting `metadata`: %+v", err)
		}

		d.Set("has_extended_immutability_policy", props.HasImmutabilityPolicy)
		d.Set("has_legal_hold", props.HasLegalHold)
	}

	resourceManagerId := client.GetResourceManagerResourceID(storageClient.SubscriptionId, account.ResourceGroup, accountName, containerName)
	d.Set("resource_manager_id", resourceManagerId)

	return nil
}
