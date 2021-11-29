package sql

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSqlAdministrator() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSqlActiveDirectoryAdministratorCreateUpdate,
		Read:   resourceSqlActiveDirectoryAdministratorRead,
		Update: resourceSqlActiveDirectoryAdministratorCreateUpdate,
		Delete: resourceSqlActiveDirectoryAdministratorDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"server_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"login": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"object_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"azuread_authentication_only": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceSqlActiveDirectoryAdministratorCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ServerAzureADAdministratorsClient
	aadOnlyAuthentictionsClient := meta.(*clients.Client).Sql.ServerAzureADOnlyAuthenticationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverName := d.Get("server_name").(string)
	resGroup := d.Get("resource_group_name").(string)
	login := d.Get("login").(string)
	objectId := uuid.FromStringOrNil(d.Get("object_id").(string))
	tenantId := uuid.FromStringOrNil(d.Get("tenant_id").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, serverName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing SQL Administrator (Server %q / Resource Group %q): %+v", serverName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_sql_active_directory_administrator", *existing.ID)
		}
	}

	aadOnlyDeleteFuture, err := aadOnlyAuthentictionsClient.Delete(ctx, resGroup, serverName)
	if err != nil {
		if aadOnlyDeleteFuture.Response() == nil || aadOnlyDeleteFuture.Response().StatusCode != http.StatusBadRequest {
			return fmt.Errorf("deleting AD Only Authentications AAD Administrator (Server %q / Resource Group %q): %+v", serverName, resGroup, err)
		}
		log.Printf("[INFO] AD Only Authentication is not removed as AD Admin is not set for AAD Administrator (Server %q / Resource Group %q): %+v", serverName, resGroup, err)
	} else if err = aadOnlyDeleteFuture.WaitForCompletionRef(ctx, aadOnlyAuthentictionsClient.Client); err != nil {
		return fmt.Errorf("waiting for deletion of AD Only Authentications for AAD Administrator (Server %q / Resource Group %q): %+v", serverName, resGroup, err)
	}

	parameters := sql.ServerAzureADAdministrator{
		AdministratorProperties: &sql.AdministratorProperties{
			AdministratorType: utils.String("ActiveDirectory"),
			Login:             utils.String(login),
			Sid:               &objectId,
			TenantID:          &tenantId,
		},
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, serverName, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating AAD Administrator (Server %q / Resource Group %q): %+v", serverName, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of AAD Administrator (Server %q / Resource Group %q): %+v", serverName, resGroup, err)
	}

	if aadOnlyAuthentictionsEnabled, ok := d.GetOk("azuread_authentication_only"); ok && aadOnlyAuthentictionsEnabled != nil && aadOnlyAuthentictionsEnabled.(bool) {
		aadOnlyAuthentictionsParams := sql.ServerAzureADOnlyAuthentication{
			AzureADOnlyAuthProperties: &sql.AzureADOnlyAuthProperties{
				AzureADOnlyAuthentication: utils.Bool(aadOnlyAuthentictionsEnabled.(bool)),
			},
		}
		aadOnlyEnabledFuture, err := aadOnlyAuthentictionsClient.CreateOrUpdate(ctx, resGroup, serverName, aadOnlyAuthentictionsParams)
		if err != nil {
			return fmt.Errorf("setting AAD only authentication for SQL Administrator (Server %q / Resource Group %q): %+v", serverName, resGroup, err)
		}

		if err = aadOnlyEnabledFuture.WaitForCompletionRef(ctx, aadOnlyAuthentictionsClient.Client); err != nil {
			return fmt.Errorf("waiting for setting of AAD only authentication for SQL Administrator (Server %q / Resource Group %q): %+v", serverName, resGroup, err)
		}
	}

	resp, err := client.Get(ctx, resGroup, serverName)
	if err != nil {
		return fmt.Errorf("retrieving SQL Administrator (Resource Group %q, Server %q): %+v", resGroup, serverName, err)
	}
	d.SetId(*resp.ID)

	return nil
}

func resourceSqlActiveDirectoryAdministratorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ServerAzureADAdministratorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AzureActiveDirectoryAdministratorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %q not found - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	d.Set("server_name", id.ServerName)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("login", resp.Login)
	d.Set("object_id", resp.Sid.String())
	d.Set("tenant_id", resp.TenantID.String())
	d.Set("azuread_authentication_only", resp.AzureADOnlyAuthentication)

	return nil
}

func resourceSqlActiveDirectoryAdministratorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ServerAzureADAdministratorsClient
	aadOnlyAuthentictionsClient := meta.(*clients.Client).Sql.ServerAzureADOnlyAuthenticationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AzureActiveDirectoryAdministratorID(d.Id())
	if err != nil {
		return err
	}

	aadOnlyDeleteFuture, err := aadOnlyAuthentictionsClient.Delete(ctx, id.ResourceGroup, id.ServerName)
	if err != nil {
		if aadOnlyDeleteFuture.Response() == nil || aadOnlyDeleteFuture.Response().StatusCode != http.StatusBadRequest {
			return fmt.Errorf("deleting AD Only Authentications for %q: %+v", id, err)
		}
		log.Printf("[INFO] AD Only Authentication is not removed as AD Admin is not set for %q: %+v", id, err)
	} else if err = aadOnlyDeleteFuture.WaitForCompletionRef(ctx, aadOnlyAuthentictionsClient.Client); err != nil {
		return fmt.Errorf("waiting for deletion of AD Only Authentications for %q: %+v", id, err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ServerName)
	if err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %q: %+v", id, err)
	}

	return nil
}
