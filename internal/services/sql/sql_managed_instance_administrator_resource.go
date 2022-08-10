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

func resourceSqlManagedInstanceAdministrator() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSqlManagedInstanceActiveDirectoryAdministratorCreateUpdate,
		Read:   resourceSqlManagedInstanceActiveDirectoryAdministratorRead,
		Update: resourceSqlManagedInstanceActiveDirectoryAdministratorCreateUpdate,
		Delete: resourceSqlManagedInstanceActiveDirectoryAdministratorDelete,

		DeprecationMessage: "The `azurerm_sql_managed_instance_active_directory_administrator` resource is deprecated and will be removed in version 4.0 of the AzureRM provider. Please use the `azurerm_mssql_managed_instance_active_directory_administrator` resource instead.",

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ManagedInstanceAzureActiveDirectoryAdministratorID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"managed_instance_name": {
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
				Default:  false,
			},
		},
	}
}

func resourceSqlManagedInstanceActiveDirectoryAdministratorCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ManagedInstanceAdministratorsClient
	aadOnlyAuthentictionsClient := meta.(*clients.Client).Sql.ManagedInstanceAzureADOnlyAuthenticationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewManagedInstanceAzureActiveDirectoryAdministratorID(subscriptionId, d.Get("resource_group_name").(string), d.Get("managed_instance_name").(string), "activeDirectory")
	login := d.Get("login").(string)
	objectId := uuid.FromStringOrNil(d.Get("object_id").(string))
	tenantId := uuid.FromStringOrNil(d.Get("tenant_id").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ManagedInstanceName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %q: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_sql_managed_instance_active_directory_administrator", id.ID())
		}
	}

	if !d.IsNewResource() {
		aadOnlyDeleteFuture, err := aadOnlyAuthentictionsClient.Delete(ctx, id.ResourceGroup, id.ManagedInstanceName)
		if err != nil {
			if aadOnlyDeleteFuture.Response() == nil || aadOnlyDeleteFuture.Response().StatusCode != http.StatusBadRequest {
				return fmt.Errorf("deleting AD Only Authentications %s: %+v", id.String(), err)
			}
			log.Printf("[INFO] AD Only Authentication is not removed as AD Admin is not set for %s: %+v", id.String(), err)
		} else if err = aadOnlyDeleteFuture.WaitForCompletionRef(ctx, aadOnlyAuthentictionsClient.Client); err != nil {
			return fmt.Errorf("waiting for deletion of AD Only Authentications %s: %+v", id.String(), err)
		}
	}

	parameters := sql.ManagedInstanceAdministrator{
		ManagedInstanceAdministratorProperties: &sql.ManagedInstanceAdministratorProperties{
			AdministratorType: utils.String("ActiveDirectory"),
			Login:             utils.String(login),
			Sid:               &objectId,
			TenantID:          &tenantId,
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedInstanceName, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q: %+v", id, err)
	}

	aadOnlyAuthentictionsParams := sql.ManagedInstanceAzureADOnlyAuthentication{
		ManagedInstanceAzureADOnlyAuthProperties: &sql.ManagedInstanceAzureADOnlyAuthProperties{
			AzureADOnlyAuthentication: utils.Bool(d.Get("azuread_authentication_only").(bool)),
		},
	}
	aadOnlyEnabledFuture, err := aadOnlyAuthentictionsClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedInstanceName, aadOnlyAuthentictionsParams)
	if err != nil {
		return fmt.Errorf("setting AAD only authentication for %s: %+v", id.String(), err)
	}

	if err = aadOnlyEnabledFuture.WaitForCompletionRef(ctx, aadOnlyAuthentictionsClient.Client); err != nil {
		return fmt.Errorf("waiting for setting of AAD only authentication for %s: %+v", id.String(), err)
	}

	d.SetId(id.ID())

	return nil
}

func resourceSqlManagedInstanceActiveDirectoryAdministratorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ManagedInstanceAdministratorsClient
	aadOnlyAuthentictionsClient := meta.(*clients.Client).Sql.ManagedInstanceAzureADOnlyAuthenticationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedInstanceAzureActiveDirectoryAdministratorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ManagedInstanceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %q was not found - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	d.Set("managed_instance_name", id.ManagedInstanceName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("login", resp.Login)
	d.Set("object_id", resp.Sid.String())
	d.Set("tenant_id", resp.TenantID.String())

	respAadOnly, err := aadOnlyAuthentictionsClient.Get(ctx, id.ResourceGroup, id.ManagedInstanceName)
	if err != nil {
		return fmt.Errorf("reading AAD only authentication for %s: %+v", id.String(), err)
	}
	aadOnly := false
	if authProps := respAadOnly.ManagedInstanceAzureADOnlyAuthProperties; authProps != nil && authProps.AzureADOnlyAuthentication != nil {
		aadOnly = *authProps.AzureADOnlyAuthentication
	}
	d.Set("azuread_authentication_only", aadOnly)

	return nil
}

func resourceSqlManagedInstanceActiveDirectoryAdministratorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ManagedInstanceAdministratorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedInstanceAzureActiveDirectoryAdministratorID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ManagedInstanceName)
	if err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %q: %+v", id, err)
	}

	return nil
}
