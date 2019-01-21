package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2015-05-01-preview/sql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSqlMiServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSqlMiServerCreateUpdate,
		Read:   resourceArmSqlMiServerRead,
		Update: resourceArmSqlMiServerCreateUpdate,
		Delete: resourceArmSqlMiServerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateMsSqlServerName,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"sku": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"GP_Gen4",
								"GP_Gen5",
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"tier": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"family": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"administrator_login": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"administrator_login_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"vcores": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  8,
				ValidateFunc: validate.IntInSlice([]int{
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
				Optional:     true,
				Default:      32,
				ValidateFunc: validate.IntBetweenAndDivisibleBy(32, 8000, 32),
			},

			"license_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "LicenseIncluded",
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

			"tags": tagsSchema(),
		},
	}
}

func resourceArmSqlMiServerCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlMiServersClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	adminUsername := d.Get("administrator_login").(string)
	licenseType := d.Get("license_type").(string)
	subnetId := d.Get("subnet_id").(string)

	tags := d.Get("tags").(map[string]interface{})
	metadata := expandTags(tags)

	parameters := sql.ManagedInstance{
		Location: utils.String(location),
		Tags:     metadata,
		ManagedInstanceProperties: &sql.ManagedInstanceProperties{
			LicenseType:        utils.String(licenseType),
			AdministratorLogin: utils.String(adminUsername),
			SubnetID:           utils.String(subnetId),
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

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)

	return resourceArmSqlServerRead(d, meta)
}

func resourceArmSqlMiServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlMiServersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["managedInstances"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading SQL Server %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading SQL Server %s: %v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if miServerProperties := resp.ManagedInstanceProperties; miServerProperties != nil {
		d.Set("license_type", miServerProperties.LicenseType)
		d.Set("administrator_login", miServerProperties.AdministratorLogin)
		d.Set("fully_qualified_domain_name", miServerProperties.FullyQualifiedDomainName)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmSqlMiServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlMiServersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["managedInstances"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting SQL Server %s: %+v", name, err)
	}

	return future.WaitForCompletionRef(ctx, client.Client)
}
