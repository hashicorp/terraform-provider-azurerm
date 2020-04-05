package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/blob/containers"
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
			"has_immutability_policy": {
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

	azureClient := storageClient.BlobContainersClient
	giovanniClient, err := storageClient.ContainersClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building Giovanni Client: %s", err)
	}

	d.SetId(giovanniClient.GetResourceID(accountName, containerName))

	if storageClient.StorageUseAzureAD {
		azureProps, err := azureClient.Get(ctx, account.ResourceGroup, accountName, containerName)
		if err != nil {
			if utils.ResponseWasNotFound(azureProps.Response) {
				return fmt.Errorf("Container %q was not found in Account %q / Resource Group %q with Azure client", containerName, accountName, account.ResourceGroup)
			}

			return fmt.Errorf("Error retrieving Container %q (Account %q / Resource Group %q) with Azure client: %s", containerName, accountName, account.ResourceGroup, err)
		}
		azPropErr := setContainerPropertiesByAzure(azureProps, d)
		if azPropErr != nil {
			return azPropErr
		}

		resourceManagerId := giovanniClient.GetResourceManagerResourceID(storageClient.SubscriptionId, account.ResourceGroup, accountName, containerName)
		d.Set("resource_manager_id", resourceManagerId)
	} else {
		gvnProps, err := giovanniClient.GetProperties(ctx, accountName, containerName)
		if err != nil {
			log.Printf("[WARN] Error reading Container %q (Storage Account %q / Resource Group %q) with Giovanni client: %s", containerName, accountName, account.ResourceGroup, err)
			azureProps, err := azureClient.Get(ctx, account.ResourceGroup, accountName, containerName)
			if err != nil {
				if utils.ResponseWasNotFound(azureProps.Response) {
					return fmt.Errorf("Container %q was not found in Account %q / Resource Group %q with Azure client", containerName, accountName, account.ResourceGroup)
				}

				return fmt.Errorf("Error retrieving Container %q (Account %q / Resource Group %q) with Azure client: %s", containerName, accountName, account.ResourceGroup, err)
			}
			azPropErr := setContainerPropertiesByAzure(azureProps, d)
			if azPropErr != nil {
				return azPropErr
			}

			resourceManagerId := giovanniClient.GetResourceManagerResourceID(storageClient.SubscriptionId, account.ResourceGroup, accountName, containerName)
			d.Set("resource_manager_id", resourceManagerId)
		} else {
			gvnPropsErr := setContainerPropertiesByGiovanni(gvnProps, d)
			if gvnPropsErr != nil {
				return gvnPropsErr
			}
			resourceManagerId := giovanniClient.GetResourceManagerResourceID(storageClient.SubscriptionId, account.ResourceGroup, accountName, containerName)
			d.Set("resource_manager_id", resourceManagerId)
		}
	}

	d.Set("name", containerName)

	d.Set("storage_account_name", accountName)

	return nil
}

func setContainerPropertiesByAzure(props storage.BlobContainer, d *schema.ResourceData) error {
	accessLevel := flattenStorageContainerAccessLevelByAzure(props.PublicAccess)
	d.Set("container_access_type", accessLevel)

	if err := d.Set("metadata", flattenMetaDataByAzure(props.Metadata)); err != nil {
		return fmt.Errorf("Error setting `metadata`: %+v", err)
	}

	d.Set("has_immutability_policy", props.HasImmutabilityPolicy)
	d.Set("has_legal_hold", props.HasLegalHold)
	return nil
}

func setContainerPropertiesByGiovanni(props containers.ContainerProperties, d *schema.ResourceData) error {
	accessLevel := flattenStorageContainerAccessLevelByGiovanni(props.AccessLevel)
	d.Set("container_access_type", accessLevel)

	if err := d.Set("metadata", FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("Error setting `metadata`: %+v", err)
	}

	d.Set("has_immutability_policy", props.HasImmutabilityPolicy)
	d.Set("has_legal_hold", props.HasLegalHold)

	return nil
}
