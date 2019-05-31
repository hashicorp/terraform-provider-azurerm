package azurerm

import (
	// "bytes"
	// "crypto/rand"
	// "encoding/base64"
	"fmt"
	// "io"
	// "log"
	// "net/url"
	// "os"
	// "runtime"
	// "strings"
	// "sync"

	// "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"

	"github.com/hashicorp/terraform/helper/validation"

	// "github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"

	// "github.com/Azure/azure-sdk-for-go/storage"
	// "github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmStorageManagementPolicy() *schema.Resource {
	return &schema.Resource{
		Create:        resourceArmStorageManagementPolicyCreate,
		Read:          resourceArmStorageManagementPolicyRead,
		Update:        resourceArmStorageManagementPolicyUpdate,
		Delete:        resourceArmStorageManagementPolicyDelete,
		SchemaVersion: 1,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"storage_account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Lifecycle"}, false),
						},
						"filters": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"prefix_match": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
									},
									"blob_types": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
									},
								},
							},
						},
						"actions": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"base_blob": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tier_to_cool_after_days_since_modification_greater_than": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"tier_to_archive_after_days_since_modification_greater_than": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"delete_after_days_since_modification_greater_than": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
									"snapshot": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"delete_after_days_since_creation_greater_than": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceArmStorageManagementPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).storage.ManagementPoliciesClient
	ctx := meta.(*ArmClient).StopContext

	storageAccountId := d.Get("storage_account_id").(string)

	rid, err := parseAzureResourceID(storageAccountId)
	resourceGroupName := rid.ResourceGroup
	storageAccountName := rid.Path["storageAccounts"]

	name := d.Get("name").(string)

	parameters := storage.ManagementPolicy{
		Name: &name,
	}

	rules := d.Get("rule").([]interface{})
	armRules, err := expandStorageManagementPolicyRules(rules)

	parameters.ManagementPolicyProperties = &storage.ManagementPolicyProperties{
		Policy: &storage.ManagementPolicySchema{
			Rules: armRules,
		},
	}

	result, err := client.CreateOrUpdate(ctx, resourceGroupName, storageAccountName, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Azure Storage Management Policy %q: %+v", storageAccountId, err)
	}

	d.SetId(*result.ID)

	return resourceArmStorageManagementPolicyRead(d, meta)
}

func expandStorageManagementPolicyRules(list []interface{}) (*[]storage.ManagementPolicyRule, error) {
	result := []storage.ManagementPolicyRule{}

	for _, tempItem := range list {
		item := tempItem.(map[string]interface{})
		policyRule, err := expandStorageManagementPolicyRule(item)
		if err != nil {
			return nil, err
		}
		result = append(result, *policyRule)
	}
	return &result, nil
}

func expandStorageManagementPolicyRule(ref map[string]interface{}) (*storage.ManagementPolicyRule, error) {
	if len(ref) == 0 {
		return nil, fmt.Errorf("Error: storage management policy rule should be defined")
	}

	name := ref["name"].(string)
	enabled := ref["enabled"].(bool)
	typeVal := ref["type"].(string)
	// TODO - add remaining fields!

	definition := storage.ManagementPolicyDefinition{
		Filters: &storage.ManagementPolicyFilter{},
		Actions: &storage.ManagementPolicyAction{},
	}
	filtersRef := ref["filters"].([]interface{})
	if len(filtersRef) == 1 {
		filterRef := filtersRef[0].(map[string]interface{})

		prefixMatches := []string{}
		prefixMatchesRef := filterRef["prefix_match"].(*schema.Set)
		if prefixMatchesRef != nil {
			for _, prefixMatchRef := range prefixMatchesRef.List() {
				prefixMatches = append(prefixMatches, prefixMatchRef.(string))
			}
		}
		definition.Filters.PrefixMatch = &prefixMatches

		blobTypes := []string{}
		blobTypesRef := filterRef["blob_types"].(*schema.Set)
		if blobTypesRef != nil {
			for _, blobTypeRef := range blobTypesRef.List() {
				blobTypes = append(blobTypes, blobTypeRef.(string))
			}
		}
		definition.Filters.BlobTypes = &blobTypes
	}
	actionsRef := ref["actions"].([]interface{})
	if len(actionsRef) == 1 {
		actionRef := actionsRef[0].(map[string]interface{})

		baseBlobsRef := actionRef["base_blob"].([]interface{})
		if len(baseBlobsRef) == 1 {
			baseBlob := &storage.ManagementPolicyBaseBlob{}
			baseBlobRef := baseBlobsRef[0].(map[string]interface{})
			if v, ok := baseBlobRef["tier_to_cool_after_days_since_modification_greater_than"]; ok {
				vInt := int32(v.(int))
				baseBlob.TierToCool = &storage.DateAfterModification{DaysAfterModificationGreaterThan: &vInt}
			}
			if v, ok := baseBlobRef["tier_to_archive_after_days_since_modification_greater_than"]; ok {
				vInt := int32(v.(int))
				baseBlob.TierToArchive = &storage.DateAfterModification{DaysAfterModificationGreaterThan: &vInt}
			}
			if v, ok := baseBlobRef["delete_after_days_since_modification_greater_than"]; ok {
				vInt := int32(v.(int))
				baseBlob.Delete = &storage.DateAfterModification{DaysAfterModificationGreaterThan: &vInt}
			}
			definition.Actions.BaseBlob = baseBlob
		}

	}

	rule := &storage.ManagementPolicyRule{
		Name:       &name,
		Enabled:    &enabled,
		Type:       &typeVal,
		Definition: &definition,
	}
	return rule, nil
}

func resourceArmStorageManagementPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	// armClient := meta.(*ArmClient)
	// ctx := armClient.StopContext

	// id, err := parseStorageBlobID(d.Id(), armClient.environment)
	// if err != nil {
	// 	return err
	// }

	// resourceGroup, err := determineResourceGroupForStorageAccount(id.storageAccountName, armClient)
	// if err != nil {
	// 	return err
	// }

	// if resourceGroup == nil {
	// 	return fmt.Errorf("Unable to determine Resource Group for Storage Account %q", id.storageAccountName)
	// }

	// blobClient, accountExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, *resourceGroup, id.storageAccountName)
	// if err != nil {
	// 	return fmt.Errorf("Error getting storage account %s: %+v", id.storageAccountName, err)
	// }
	// if !accountExists {
	// 	return fmt.Errorf("Storage account %s not found in resource group %s", id.storageAccountName, *resourceGroup)
	// }

	// container := blobClient.GetContainerReference(id.containerName)
	// blob := container.GetBlobReference(id.blobName)

	// if d.HasChange("content_type") {
	// 	blob.Properties.ContentType = d.Get("content_type").(string)
	// }

	// options := &storage.SetBlobPropertiesOptions{}
	// err = blob.SetProperties(options)
	// if err != nil {
	// 	return fmt.Errorf("Error setting properties of blob %s (container %s, storage account %s): %+v", id.blobName, id.containerName, id.storageAccountName, err)
	// }

	// if d.HasChange("metadata") {
	// 	blob.Metadata = expandStorageAccountBlobMetadata(d)

	// 	opts := &storage.SetBlobMetadataOptions{}
	// 	if err := blob.SetMetadata(opts); err != nil {
	// 		return fmt.Errorf("Error setting metadata for storage blob on Azure: %s", err)
	// 	}
	// }

	// return nil
	return fmt.Errorf("Not implemented!")
}

func resourceArmStorageManagementPolicyRead(d *schema.ResourceData, meta interface{}) error {
	// armClient := meta.(*ArmClient)
	// ctx := armClient.StopContext

	// id, err := parseStorageBlobID(d.Id(), armClient.environment)
	// if err != nil {
	// 	return err
	// }

	// resourceGroup, err := determineResourceGroupForStorageAccount(id.storageAccountName, armClient)
	// if err != nil {
	// 	return err
	// }

	// if resourceGroup == nil {
	// 	return fmt.Errorf("Unable to determine Resource Group for Storage Account %q", id.storageAccountName)
	// }

	// blobClient, accountExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, *resourceGroup, id.storageAccountName)
	// if err != nil {
	// 	return err
	// }
	// if !accountExists {
	// 	log.Printf("[DEBUG] Storage account %q not found, removing blob %q from state", id.storageAccountName, d.Id())
	// 	d.SetId("")
	// 	return nil
	// }

	// log.Printf("[INFO] Checking for existence of storage blob %q in container %q.", id.blobName, id.containerName)
	// container := blobClient.GetContainerReference(id.containerName)
	// blob := container.GetBlobReference(id.blobName)
	// exists, err := blob.Exists()
	// if err != nil {
	// 	return fmt.Errorf("error checking for existence of storage blob %q: %s", id.blobName, err)
	// }

	// if !exists {
	// 	log.Printf("[INFO] Storage blob %q no longer exists, removing from state...", id.blobName)
	// 	d.SetId("")
	// 	return nil
	// }

	// options := &storage.GetBlobPropertiesOptions{}
	// err = blob.GetProperties(options)
	// if err != nil {
	// 	return fmt.Errorf("Error getting properties of blob %s (container %s, storage account %s): %+v", id.blobName, id.containerName, id.storageAccountName, err)
	// }

	// metadataOptions := &storage.GetBlobMetadataOptions{}
	// err = blob.GetMetadata(metadataOptions)
	// if err != nil {
	// 	return fmt.Errorf("Error getting metadata of blob %s (container %s, storage account %s): %+v", id.blobName, id.containerName, id.storageAccountName, err)
	// }

	// d.Set("name", id.blobName)
	// d.Set("storage_container_name", id.containerName)
	// d.Set("storage_account_name", id.storageAccountName)
	// d.Set("resource_group_name", resourceGroup)

	// d.Set("content_type", blob.Properties.ContentType)

	// d.Set("source_uri", blob.Properties.CopySource)

	// blobType := strings.ToLower(strings.Replace(string(blob.Properties.BlobType), "Blob", "", 1))
	// d.Set("type", blobType)

	// u := blob.GetURL()
	// if u == "" {
	// 	log.Printf("[INFO] URL for %q is empty", id.blobName)
	// }
	// d.Set("url", u)
	// d.Set("metadata", flattenStorageAccountBlobMetadata(blob.Metadata))

	// return nil

	return fmt.Errorf("Not implemented!")

}

func resourceArmStorageManagementPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	// armClient := meta.(*ArmClient)
	// ctx := armClient.StopContext

	// id, err := parseStorageBlobID(d.Id(), armClient.environment)
	// if err != nil {
	// 	return err
	// }

	// resourceGroup, err := determineResourceGroupForStorageAccount(id.storageAccountName, armClient)
	// if err != nil {
	// 	return fmt.Errorf("Unable to determine Resource Group for Storage Account %q: %+v", id.storageAccountName, err)
	// }
	// if resourceGroup == nil {
	// 	log.Printf("[INFO] Resource Group doesn't exist so the blob won't exist")
	// 	return nil
	// }

	// blobClient, accountExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, *resourceGroup, id.storageAccountName)
	// if err != nil {
	// 	return err
	// }
	// if !accountExists {
	// 	log.Printf("[INFO] Storage Account %q doesn't exist so the blob won't exist", id.storageAccountName)
	// 	return nil
	// }

	// log.Printf("[INFO] Deleting storage blob %q", id.blobName)
	// options := &storage.DeleteBlobOptions{}
	// container := blobClient.GetContainerReference(id.containerName)
	// blob := container.GetBlobReference(id.blobName)
	// _, err = blob.DeleteIfExists(options)
	// if err != nil {
	// 	return fmt.Errorf("Error deleting storage blob %q: %s", id.blobName, err)
	// }

	// return nil
	return fmt.Errorf("Not implemented!")
}
