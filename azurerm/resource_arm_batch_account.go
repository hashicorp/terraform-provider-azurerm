package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2017-09-01/batch"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmBatchAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmBatchAccountCreate,
		Read:   resourceArmBatchAccountRead,
		Update: resourceArmBatchAccountUpdate,
		Delete: resourceArmBatchAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMBatchAccountName,
			},
			"resource_group_name": resourceGroupNameSchema(),
			"location":            locationSchema(),
			"storage_account_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},
			"pool_allocation_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(batch.BatchService),
				ValidateFunc: validation.StringInSlice([]string{
					string(batch.BatchService),
					string(batch.UserSubscription),
				}, false),
			},
			"primary_access_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"secondary_access_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"account_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceArmBatchAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).batchAccountClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure Batch account creation.")

	resourceGroupName := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	storageAccountId := d.Get("storage_account_id").(string)
	poolAllocationMode := d.Get("pool_allocation_mode").(string)
	tags := d.Get("tags").(map[string]interface{})

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroupName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Batch Account %q (Resource Group %q): %s", name, resourceGroupName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_batch_account", *existing.ID)
		}
	}

	parameters := batch.AccountCreateParameters{
		Location: &location,
		AccountCreateProperties: &batch.AccountCreateProperties{
			PoolAllocationMode: batch.PoolAllocationMode(poolAllocationMode),
		},
		Tags: expandTags(tags),
	}

	if storageAccountId != "" {
		parameters.AccountCreateProperties.AutoStorage = &batch.AutoStorageBaseProperties{
			StorageAccountID: &storageAccountId,
		}
	}

	future, err := client.Create(ctx, resourceGroupName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Batch account %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Batch account %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	read, err := client.Get(ctx, resourceGroupName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Batch account %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Batch account %q (resource group %q) ID", name, resourceGroupName)
	}

	d.SetId(*read.ID)

	return resourceArmBatchAccountRead(d, meta)
}

func resourceArmBatchAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).batchAccountClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["batchAccounts"]
	resourceGroupName := id.ResourceGroup

	resp, err := client.Get(ctx, resourceGroupName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			log.Printf("[DEBUG] Batch Account %q was not found in Resource Group %q - removing from state!", name, resourceGroupName)
			return nil
		}
		return fmt.Errorf("Error reading the state of Batch account %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroupName)
	d.Set("account_endpoint", resp.AccountEndpoint)

	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.AccountProperties; props != nil {
		if autoStorage := props.AutoStorage; autoStorage != nil {
			d.Set("storage_account_id", autoStorage.StorageAccountID)
		}
		d.Set("pool_allocation_mode", props.PoolAllocationMode)
	}

	if d.Get("pool_allocation_mode").(string) == string(batch.BatchService) {
		keys, err := client.GetKeys(ctx, resourceGroupName, name)

		if err != nil {
			return err
		}

		d.Set("primary_access_key", keys.Primary)
		d.Set("secondary_access_key", keys.Secondary)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmBatchAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).batchAccountClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure Batch account update.")

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["batchAccounts"]
	resourceGroupName := id.ResourceGroup

	storageAccountId := d.Get("storage_account_id").(string)
	tags := d.Get("tags").(map[string]interface{})

	parameters := batch.AccountUpdateParameters{
		AccountUpdateProperties: &batch.AccountUpdateProperties{
			AutoStorage: &batch.AutoStorageBaseProperties{
				StorageAccountID: &storageAccountId,
			},
		},
		Tags: expandTags(tags),
	}

	if _, err = client.Update(ctx, resourceGroupName, name, parameters); err != nil {
		return fmt.Errorf("Error updating Batch account %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	read, err := client.Get(ctx, resourceGroupName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Batch account %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Batch account %q (resource group %q) ID", name, resourceGroupName)
	}

	d.SetId(*read.ID)

	return resourceArmBatchAccountRead(d, meta)
}

func resourceArmBatchAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).batchAccountClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["batchAccounts"]
	resourceGroupName := id.ResourceGroup

	future, err := client.Delete(ctx, resourceGroupName, name)
	if err != nil {
		return fmt.Errorf("Error deleting Batch account %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deletion of Batch account %q (Resource Group %q): %+v", name, resourceGroupName, err)
		}
	}

	return nil
}

func validateAzureRMBatchAccountName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-z0-9]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("lowercase letters and numbers only are allowed in %q: %q", k, value))
	}

	if 3 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 3 characters: %q", k, value))
	}

	if len(value) > 24 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 24 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}
