package azurerm

import (
	"fmt"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/arm/sql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"log"
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDBAccountName,
			},

			"location": locationSchema(),

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"version": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.TwoFullStopZero),
					string(sql.OneTwoFullStopZero),
				}, true),
				// TODO: is this ForceNew?
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

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := d.Get("location").(string)
	adminUsername := d.Get("administrator_login").(string)
	adminPassword := d.Get("administrator_login_password").(string)
	version := d.Get("version").(string)

	tags := d.Get("tags").(map[string]interface{})
	metadata := expandTags(tags)

	parameters := sql.Server{
		Location: &location,
		Tags:     metadata,
		ServerProperties: &sql.ServerProperties{
			Version:                    sql.ServerVersion(version),
			AdministratorLogin:         &adminUsername,
			AdministratorLoginPassword: &adminPassword,
		},
	}

	response, err := client.CreateOrUpdate(resGroup, name, parameters)
	if err != nil {
		// if the name is in-use, Azure returns a 409 "Unknown Service Error" which is a bad UX
		if utils.ResponseWasConflict(response.Response) {
			return fmt.Errorf("SQL Server names need to be globally unique and '%s' is already in use.", name)
		}

		return err
	}

	if response.ID == nil {
		return fmt.Errorf("Cannot create SQL Server %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*response.ID)

	return resourceArmSqlServerRead(d, meta)
}

func resourceArmSqlServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlServersClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["servers"]

	resp, err := client.Get(resGroup, name)
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
	d.Set("location", azureRMNormalizeLocation(*resp.Location))

	if serverProperties := resp.ServerProperties; serverProperties != nil {
		d.Set("version", string(serverProperties.Version))
		d.Set("administrator_login", serverProperties.AdministratorLogin)
		d.Set("fully_qualified_domain_name", serverProperties.FullyQualifiedDomainName)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmSqlServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlServersClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["servers"]

	response, err := client.Delete(resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(response) {
			return nil
		}

		return fmt.Errorf("Error deleting SQL Server %s: %+v", name, err)
	}

	return nil
}
