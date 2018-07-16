package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2015-05-01-preview/sql"
	"github.com/hashicorp/terraform/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSqlAdministrator() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSqlActiveDirectoryAdministratorCreateUpdate,
		Read:   resourceArmSqlActiveDirectoryAdministratorRead,
		Update: resourceArmSqlActiveDirectoryAdministratorCreateUpdate,
		Delete: resourceArmSqlActiveDirectoryAdministratorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"server_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"login": {
				Type:     schema.TypeString,
				Required: true,
			},

			"object_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateUUID,
			},

			"tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateUUID,
			},
		},
	}
}

func resourceArmSqlActiveDirectoryAdministratorCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlServerAzureADAdministratorsClient
	ctx := meta.(*ArmClient).StopContext

	serverName := d.Get("server_name").(string)
	resGroup := d.Get("resource_group_name").(string)
	login := d.Get("login").(string)
	objectId := uuid.FromStringOrNil(d.Get("object_id").(string))
	tenantId := uuid.FromStringOrNil(d.Get("tenant_id").(string))
	parameters := sql.ServerAzureADAdministrator{
		ServerAdministratorProperties: &sql.ServerAdministratorProperties{
			AdministratorType: utils.String("ActiveDirectory"),
			Login:             utils.String(login),
			Sid:               &objectId,
			TenantID:          &tenantId,
		},
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, serverName, parameters)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, resGroup, serverName)
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)

	return nil
}

func resourceArmSqlActiveDirectoryAdministratorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlServerAzureADAdministratorsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serverName := id.Path["servers"]

	resp, err := client.Get(ctx, resourceGroup, serverName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading SQL AD administrator %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading SQL AD administrator: %+v", err)
	}

	d.Set("resource_group_name", resourceGroup)
	d.Set("server_name", serverName)
	d.Set("login", resp.Login)
	d.Set("object_id", resp.Sid.String())
	d.Set("tenant_id", resp.TenantID.String())

	return nil
}

func resourceArmSqlActiveDirectoryAdministratorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlServerAzureADAdministratorsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serverName := id.Path["servers"]

	_, err = client.Delete(ctx, resourceGroup, serverName)
	if err != nil {
		return fmt.Errorf("Error deleting SQL AD administrator: %+v", err)
	}

	return nil
}
