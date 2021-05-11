package sql

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2018-06-01-preview/sql"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSqlMiServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSqlMiServerCreateUpdate,
		Read:   resourceArmSqlMiServerRead,
		Update: resourceArmSqlMiServerCreateUpdate,
		Delete: resourceArmSqlMiServerDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ManagedInstanceID(id)
			return err
		}),

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
				ValidateFunc: validate.ValidateMsSqlServerName,
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

			"administrator_login": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"administrator_login_password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"vcores": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntInSlice([]int{
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

			"storage_size_in_gb": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(32, 8192),
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

func resourceArmSqlMiServerCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ManagedInstancesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	adminUsername := d.Get("administrator_login").(string)
	licenseType := d.Get("license_type").(string)
	subnetId := d.Get("subnet_id").(string)
	metadata := tags.Expand(d.Get("tags").(map[string]interface{}))

	id := parse.NewManagedInstanceID(subscriptionId, resGroup, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Managed Instance %q: %s", id.ID(), err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_sql_managed_instance", *existing.ID)
		}
	}

	sku, err := expandManagedInstanceSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("error expanding sku_name for SQL Managed Instance Server %q (Resource Group %q): %v", name, resGroup, err)
	}

	parameters := sql.ManagedInstance{
		Sku:      sku,
		Location: utils.String(location),
		Tags:     metadata,
		ManagedInstanceProperties: &sql.ManagedInstanceProperties{
			LicenseType:        sql.ManagedInstanceLicenseType(licenseType),
			AdministratorLogin: utils.String(adminUsername),
			SubnetID:           utils.String(subnetId),
			StorageSizeInGB:    utils.Int32(int32(d.Get("storage_size_in_gb").(int))),
			VCores:             utils.Int32(int32(d.Get("vcores").(int))),
		},
	}

	if d.HasChange("administrator_login_password") {
		adminPassword := d.Get("administrator_login_password").(string)
		parameters.ManagedInstanceProperties.AdministratorLoginPassword = utils.String(adminPassword)
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

	d.SetId(id.ID())

	return resourceArmSqlMiServerRead(d, meta)
}

func resourceArmSqlMiServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ManagedInstancesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedInstanceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading SQL Managed Instance %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading SQL Managed Instance %q: %v", id.ID(), err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	if miServerProperties := resp.ManagedInstanceProperties; miServerProperties != nil {
		d.Set("license_type", miServerProperties.LicenseType)
		d.Set("administrator_login", miServerProperties.AdministratorLogin)
		d.Set("subnet_id", miServerProperties.SubnetID)
		d.Set("storage_size_in_gb", miServerProperties.StorageSizeInGB)
		d.Set("vcores", miServerProperties.VCores)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmSqlMiServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ManagedInstancesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedInstanceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting SQL Managed Instance %q: %+v", id.ID(), err)
	}

	return future.WaitForCompletionRef(ctx, client.Client)
}

func expandManagedInstanceSkuName(skuName string) (*sql.Sku, error) {
	parts := strings.Split(skuName, "_")
	if len(parts) != 2 {
		return nil, fmt.Errorf("sku_name (%s) has the wrong number of parts (%d) after splitting on _", skuName, len(parts))
	}

	var tier string
	switch parts[0] {
	case "GP":
		tier = "GeneralPurpose"
	case "BC":
		tier = "BusinessCritical"
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", skuName, parts[0])
	}

	return &sql.Sku{
		Name:   utils.String(skuName),
		Tier:   utils.String(tier),
		Family: utils.String(parts[1]),
	}, nil
}
