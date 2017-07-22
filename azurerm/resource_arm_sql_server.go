package azurerm

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/sql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceArmSqlServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSqlServerCreateOrUpdate,
		Read:   resourceArmSqlServerRead,
		Update: resourceArmSqlServerCreateOrUpdate,
		Delete: resourceArmSqlServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"version": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.OneTwoFullStopZero),
					string(sql.TwoFullStopZero),
				}, true),
			},

			"administrator_login": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"administrator_login_password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"fully_qualified_domain_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmSqlServerCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	sqlServersClient := meta.(*ArmClient).sqlServersClient

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := d.Get("location").(string)
	admin := d.Get("administrator_login").(string)
	adminPw := d.Get("administrator_login_password").(string)
	fqdn := d.Get("fully_qualified_domain_name").(string)

	tags := d.Get("tags").(map[string]interface{})
	metadata := expandTags(tags)

	props := sql.ServerProperties{
		FullyQualifiedDomainName:   &fqdn,
		Version:                    sql.ServerVersion(d.Get("version").(string)),
		AdministratorLogin:         &admin,
		AdministratorLoginPassword: &adminPw,
	}

	parameters := sql.Server{
		Name:             &name,
		ServerProperties: &props,
		Tags:             metadata,
		Location:         &location,
	}

	result, err := sqlServersClient.CreateOrUpdate(resGroup, name, parameters)
	if err != nil {
		return err
	}

	if result.ID == nil {
		return fmt.Errorf("Cannot create Sql Server %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*result.ID)

	return resourceArmSqlServerRead(d, meta)
}

func resourceArmSqlServerRead(d *schema.ResourceData, meta interface{}) error {
	sqlServersClient := meta.(*ArmClient).sqlServersClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["servers"]

	result, err := sqlServersClient.Get(resGroup, name)
	if err != nil {
		return fmt.Errorf("Error reading SQL Server %s: %v", name, err)
	}
	if result.Response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	if serverProperties := result.ServerProperties; serverProperties != nil {
		d.Set("version", string(serverProperties.Version))
		d.Set("administrator_login", serverProperties.AdministratorLogin)
		d.Set("fully_qualified_domain_name", serverProperties.FullyQualifiedDomainName)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", azureRMNormalizeLocation(*result.Location))

	flattenAndSetTags(d, result.Tags)

	return nil
}

func resourceArmSqlServerDelete(d *schema.ResourceData, meta interface{}) error {
	sqlServersClient := meta.(*ArmClient).sqlServersClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["servers"]

	result, error := sqlServersClient.Delete(resGroup, name)
	if result.Response.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting SQL Server %s: %+v", name, error)
	}

	return nil
}
