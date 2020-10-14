package sql

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2018-06-01-preview/sql"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSqlInstancePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSqlInstancePoolCreateUpdate,
		Read:   resourceArmSqlInstancePoolRead,
		Update: resourceArmSqlInstancePoolCreateUpdate,
		Delete: resourceArmSqlInstancePoolDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(24 * time.Hour),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(24 * time.Hour),
			Delete: schema.DefaultTimeout(24 * time.Hour),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateMsSqlServerName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"GP_Gen4",
					"GP_Gen5",
					"BC_Gen4",
					"BC_Gen5",
				}, false),
			},

			"vcores": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validate.IntInSlice([]int{
					4,
					8,
					16,
					24,
					32,
					40,
					64,
					80,
				}),
			},

			"license_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"LicenseIncluded",
					"BasePrice",
				}, true),
			},

			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmSqlInstancePoolCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.InstancePoolsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	subnetId := d.Get("subnet_id").(string)
	metadata := tags.Expand(d.Get("tags").(map[string]interface{}))

	sku, err := expandManagedInstanceSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("error expanding sku_name for SQL Managed Instance Server %q (Resource Group %q): %v", name, resGroup, err)
	}

	parameters := sql.InstancePool{
		Sku:      sku,
		Location: utils.String(location),
		Tags:     metadata,
		InstancePoolProperties: &sql.InstancePoolProperties{
			SubnetID:    utils.String(subnetId),
			VCores:      utils.Int32(8),
			LicenseType: sql.BasePrice,
		},
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, parameters)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasConflict(future.Response()) {
			return fmt.Errorf("SQL Server names need to be globally unique and %q is already in use.", name)
		}

		return err
	}

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)

	return resourceArmSqlInstancePoolRead(d, meta)
}

func resourceArmSqlInstancePoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.InstancePoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["instancePools"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading SQL Server %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading SQL Server %s: %v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if miServerProperties := resp.InstancePoolProperties; miServerProperties != nil {
		d.Set("license_type", miServerProperties.LicenseType)
		d.Set("subnet_id", miServerProperties.SubnetID)
		d.Set("vcores", miServerProperties.VCores)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmSqlInstancePoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.InstancePoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["instancePools"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting SQL Server %s: %+v", name, err)
	}

	return future.WaitForCompletionRef(ctx, client.Client)
}
