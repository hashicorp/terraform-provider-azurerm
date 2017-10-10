package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/arm/mysql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMySQLConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMySQLConfigurationCreate,
		Read:   resourceArmMySQLConfigurationRead,
		Delete: resourceArmMySQLConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"server_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"value": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceArmMySQLConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mysqlConfigurationsClient

	log.Printf("[INFO] preparing arguments for AzureRM MySQL Configuration creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)

	value := d.Get("value").(string)

	properties := mysql.Configuration{
		ConfigurationProperties: &mysql.ConfigurationProperties{
			Value: utils.String(value),
		},
	}

	_, error := client.CreateOrUpdate(resGroup, serverName, name, properties, make(chan struct{}))
	err := <-error
	if err != nil {
		return err
	}

	read, err := client.Get(resGroup, serverName, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read MySQL Configuration %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmMySQLConfigurationRead(d, meta)
}

func resourceArmMySQLConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mysqlConfigurationsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["configurations"]

	resp, err := client.Get(resGroup, serverName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] MySQL Configuration '%s' was not found (resource group '%s')", name, resGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure MySQL Configuration %s: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("server_name", serverName)
	d.Set("resource_group_name", resGroup)
	d.Set("value", resp.ConfigurationProperties.Value)

	return nil
}

func resourceArmMySQLConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mysqlConfigurationsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["configurations"]

	// "delete" = resetting this to the default value
	resp, err := client.Get(resGroup, serverName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving MySQL Configuration '%s': %+v", name, err)
	}

	properties := mysql.Configuration{
		ConfigurationProperties: &mysql.ConfigurationProperties{
			// we can alternatively set `source: "system-default"`
			Value: resp.DefaultValue,
		},
	}

	_, error := client.CreateOrUpdate(resGroup, serverName, name, properties, make(chan struct{}))
	err = <-error
	return err
}
