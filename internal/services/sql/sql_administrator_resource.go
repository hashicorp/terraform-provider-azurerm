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

		DeprecationMessage: "The `azurerm_sql_active_directory_administrator` resource is deprecated and will be removed in version 4.0 of the AzureRM provider. Please use the `azuread_administrator` block of the `azurerm_mssql_server` resource instead.",

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AzureActiveDirectoryAdministratorID(id)
			return err
		}),

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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	aadOnlyAuthClient := meta.(*clients.Client).Sql.ServerAzureADOnlyAuthenticationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	login := d.Get("login").(string)
	objectId := uuid.FromStringOrNil(d.Get("object_id").(string))
	tenantId := uuid.FromStringOrNil(d.Get("tenant_id").(string))

	id := parse.NewAzureActiveDirectoryAdministratorID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), "ActiveDirectory")
	serverId := parse.NewServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, serverId.ResourceGroup, serverId.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_sql_active_directory_administrator", id.ID())
		}
	}

	aadOnlyDeleteFuture, err := aadOnlyAuthClient.Delete(ctx, serverId.ResourceGroup, serverId.Name)
	if err != nil {
		if aadOnlyDeleteFuture.Response() == nil || aadOnlyDeleteFuture.Response().StatusCode != http.StatusBadRequest {
			return fmt.Errorf("deleting AD Only Authentications AAD Administrator for %s: %+v", serverId, err)
		}
		log.Printf("[INFO] AD Only Authentication is not removed as AD Admin is not set for AAD Administrator %s: %+v", serverId, err)
	} else if err = aadOnlyDeleteFuture.WaitForCompletionRef(ctx, aadOnlyAuthClient.Client); err != nil {
		return fmt.Errorf("waiting for deletion of AD Only Authentications for %s: %+v", serverId, err)
	}

	parameters := sql.ServerAzureADAdministrator{
		AdministratorProperties: &sql.AdministratorProperties{
			AdministratorType: utils.String("ActiveDirectory"),
			Login:             utils.String(login),
			Sid:               &objectId,
			TenantID:          &tenantId,
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	if aadOnlyAuthentictionsEnabled, ok := d.GetOk("azuread_authentication_only"); ok && aadOnlyAuthentictionsEnabled != nil && aadOnlyAuthentictionsEnabled.(bool) {
		aadOnlyAuthentictionsParams := sql.ServerAzureADOnlyAuthentication{
			AzureADOnlyAuthProperties: &sql.AzureADOnlyAuthProperties{
				AzureADOnlyAuthentication: utils.Bool(aadOnlyAuthentictionsEnabled.(bool)),
			},
		}
		aadOnlyEnabledFuture, err := aadOnlyAuthClient.CreateOrUpdate(ctx, serverId.ResourceGroup, serverId.Name, aadOnlyAuthentictionsParams)
		if err != nil {
			return fmt.Errorf("setting AAD only authentication for %s: %+v", serverId, err)
		}

		if err = aadOnlyEnabledFuture.WaitForCompletionRef(ctx, aadOnlyAuthClient.Client); err != nil {
			return fmt.Errorf("waiting for AAD only authentication to be set for %s: %+v", serverId, err)
		}
	}

	d.SetId(id.ID())
	return resourceSqlActiveDirectoryAdministratorRead(d, meta)
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

	if props := resp.AdministratorProperties; props != nil {
		d.Set("azuread_authentication_only", props.AzureADOnlyAuthentication)
		d.Set("login", props.Login)
		d.Set("object_id", props.Sid.String())
		d.Set("tenant_id", props.TenantID.String())
	}

	return nil
}

func resourceSqlActiveDirectoryAdministratorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ServerAzureADAdministratorsClient
	aadOnlyAuthClient := meta.(*clients.Client).Sql.ServerAzureADOnlyAuthenticationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AzureActiveDirectoryAdministratorID(d.Id())
	if err != nil {
		return err
	}

	serverId := parse.NewServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)
	aadOnlyDeleteFuture, err := aadOnlyAuthClient.Delete(ctx, serverId.ResourceGroup, serverId.Name)
	if err != nil {
		if aadOnlyDeleteFuture.Response() == nil || aadOnlyDeleteFuture.Response().StatusCode != http.StatusBadRequest {
			return fmt.Errorf("deleting AD Only Authentications for %s: %+v", serverId, err)
		}
		log.Printf("[INFO] AD Only Authentication is not removed as AD Admin is not set for %s: %+v", serverId, err)
	} else if err = aadOnlyDeleteFuture.WaitForCompletionRef(ctx, aadOnlyAuthClient.Client); err != nil {
		return fmt.Errorf("waiting for deletion of AD Only Authentications for %s: %+v", serverId, err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ServerName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}

	return nil
}
