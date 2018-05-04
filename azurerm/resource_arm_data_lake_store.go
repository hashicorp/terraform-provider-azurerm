package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/datalake/store/mgmt/2016-11-01/account"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDateLakeStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDateLakeStoreCreate,
		Read:   resourceArmDateLakeStoreRead,
		Update: resourceArmDateLakeStoreUpdate,
		Delete: resourceArmDateLakeStoreDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"tier": {
				Type:     schema.TypeString,
				Required: false,
				ForceNew: false,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmDateLakeStoreCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataLakeStoreAccountClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM Date Lake Store creation.")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	tier := d.Get("tier").(string)
	tags := d.Get("tags").(map[string]interface{})

	dateLakeStore := accounts.CreateDataLakeStoreAccountParameters{
		Location: &location,
		Tags:     expandTags(tags),
		CreateDataLakeStoreAccountProperties: &accounts.CreateDataLakeStoreAccountProperties{
			NewTier: tier,
		},
	}

	future, err := client.Create(ctx, resGroup, name, dateLakeStore)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Data Lake Store %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmDateLakeStoreRead(d, meta)
}

func resourceArmDateLakeStoreUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataLakeStoreAccountClient
	ctx := meta.(*ArmClient).StopContext

	if !d.HasChange("tags") {
		return nil
	}

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	newTags := d.Get("tags").(map[string]interface{})
	newTier := d.Get("tier").(string)

	props := accounts.UpdateDataLakeStoreAccountParameters{
		Tags: expandTags(newTags),
		UpdateDataLakeStoreAccountProperties: &accounts.UpdateDataLakeStoreAccountProperties{
			NewTier: newTier,
		}
	}

	future, err := client.Update(ctx, resourceGroup, name, props)
	if err != nil {
		return fmt.Errorf("Error issuing update request for Data Lake Store %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the update of Data Lake Store %q (Resource Group %q) to commplete: %+v", name, resourceGroup, err)
	}

	return resourceArmDateLakeStoreRead(d, meta)
}

func resourceArmDateLakeStoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataLakeStoreAccountClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["accounts"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Data Lake Store %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	d.Set("tier", resp.Type)

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmDateLakeStoreDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataLakeStoreAccountClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["accounts"]
	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing delete request for Azure Data Lake Store %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error waiting for Azure Data Lake Store %q (Resource Group %q) to be deleted: %+v", name, resourceGroup, err)
	}

	return err
}
