package azurerm

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/sql"
	"github.com/hashicorp/terraform/helper/schema"
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

	version := d.Get("version").(string)
	var serverVersion sql.ServerVersion
	if version == string(sql.OneTwoFullStopZero) {
		serverVersion = sql.OneTwoFullStopZero
	}
	if version == string(sql.TwoFullStopZero) {
		serverVersion = sql.TwoFullStopZero
	}
	if serverVersion == "" {
		return fmt.Errorf("Invalid server version provided. It must be one of 12.0 or 2.0, passed as string: %s", version)
	}

	tags := d.Get("tags").(map[string]interface{})
	metadata := expandTags(tags)

	props := sql.ServerProperties{
		FullyQualifiedDomainName:   &fqdn,
		Version:                    serverVersion,
		AdministratorLogin:         &admin,
		AdministratorLoginPassword: &adminPw,
	}

	parameters := sql.Server{
		Name:             &name,
		ServerProperties: &props,
		Tags:             metadata,
		Location:         &location,
	}

	//last parameter is set to empty to allow updates to records after creation
	// (per SDK, set it to '*' to prevent updates, all other values are ignored)
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

	serverProperties := *result.ServerProperties

	serverVersion := serverProperties.Version

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("version", string(serverVersion))
	d.Set("location", *result.Location)
	d.Set("administrator_login", *serverProperties.AdministratorLogin)
	d.Set("fully_qualified_domain_name", *serverProperties.FullyQualifiedDomainName)

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
		return fmt.Errorf("Error deleting SQL Server %s: %s", name, error)
	}

	return nil
}
