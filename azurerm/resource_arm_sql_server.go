package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2015-05-01-preview/sql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSqlServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSqlServerCreateUpdate,
		Read:   resourceArmSqlServerRead,
		Update: resourceArmSqlServerCreateUpdate,
		Delete: resourceArmSqlServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-z0-9]{3,50}$"),
					"SQL server name must be 3 - 50 characters long, contain only letters, numbers and hyphens.",
				),
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string("2.0"),
					string("12.0"),
				}, true),
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

			"fully_qualified_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmSqlServerCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlServersClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	adminUsername := d.Get("administrator_login").(string)
	adminPassword := d.Get("administrator_login_password").(string)
	version := d.Get("version").(string)

	tags := d.Get("tags").(map[string]interface{})
	metadata := expandTags(tags)

	parameters := sql.Server{
		Location: utils.String(location),
		Tags:     metadata,
		ServerProperties: &sql.ServerProperties{
			Version:                    utils.String(version),
			AdministratorLogin:         utils.String(adminUsername),
			AdministratorLoginPassword: utils.String(adminPassword),
		},
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, parameters)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {

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

func resourceArmSqlServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlServersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["servers"]

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

	if serverProperties := resp.ServerProperties; serverProperties != nil {
		d.Set("version", serverProperties.Version)
		d.Set("administrator_login", serverProperties.AdministratorLogin)
		d.Set("fully_qualified_domain_name", serverProperties.FullyQualifiedDomainName)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmSqlServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlServersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["servers"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting SQL Server %s: %+v", name, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}

	return nil
}
