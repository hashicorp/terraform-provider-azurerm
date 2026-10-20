// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_sql_database

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryLinkedServiceAzureSQLDatabaseCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	subscriptionId := meta.(*clients.Client).DataFactory.LinkedServiceClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewLinkedServiceID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing Data Factory Azure SQL Database %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_data_factory_linked_service_azure_sql_database", id.ID())
			}
		}
	}

	sqlDatabaseProperties := &datafactory.AzureSQLDatabaseLinkedServiceTypeProperties{}

	if v, ok := d.GetOk("connection_string"); ok {
		if d.Get("use_managed_identity").(bool) {
			sqlDatabaseProperties.ConnectionString = pointer.To(v.(string))
		} else {
			sqlDatabaseProperties.ConnectionString = &datafactory.SecureString{
				Value: pointer.To(v.(string)),
				Type:  datafactory.TypeSecureString,
			}
		}
	}

	if v, ok := d.GetOk("key_vault_connection_string"); ok {
		sqlDatabaseProperties.ConnectionString = helpers.ExpandAzureKeyVaultSecretReference(v.([]any))
	}

	if d.Get("use_managed_identity").(bool) {
		sqlDatabaseProperties.Tenant = pointer.To(d.Get("tenant_id").(string))
	} else {
		secureString := datafactory.SecureString{
			Value: pointer.To(d.Get("service_principal_key").(string)),
			Type:  datafactory.TypeSecureString,
		}

		sqlDatabaseProperties.ServicePrincipalID = pointer.To(d.Get("service_principal_id").(string))
		sqlDatabaseProperties.ServicePrincipalKey = &secureString
		if v := d.Get("tenant_id").(string); v != "" {
			sqlDatabaseProperties.Tenant = pointer.To(v)
		}
	}

	if v, ok := d.GetOk("key_vault_password"); ok {
		password := v.([]any)
		sqlDatabaseProperties.Password = helpers.ExpandAzureKeyVaultSecretReference(password)
	}

	azureSQLDatabaseLinkedService := &datafactory.AzureSQLDatabaseLinkedService{
		Description: pointer.To(d.Get("description").(string)),
		AzureSQLDatabaseLinkedServiceTypeProperties: sqlDatabaseProperties,
		Type: datafactory.TypeBasicLinkedServiceTypeAzureSQLDatabase,
	}

	if v, ok := d.GetOk("parameters"); ok {
		azureSQLDatabaseLinkedService.Parameters = helpers.ExpandLinkedServiceParameters(v.(map[string]any))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		azureSQLDatabaseLinkedService.ConnectVia = helpers.ExpandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		azureSQLDatabaseLinkedService.AdditionalProperties = v.(map[string]any)
	}

	if v, ok := d.GetOk("annotations"); ok {
		azureSQLDatabaseLinkedService.Annotations = pointer.To(v.([]any))
	}

	if credentialName := d.Get("credential_name").(string); credentialName != "" {
		azureSQLDatabaseLinkedService.Credential = &datafactory.CredentialReference{
			ReferenceName: pointer.To(credentialName),
		}
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: azureSQLDatabaseLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, linkedService, ""); err != nil {
		return fmt.Errorf("creating/updating Data Factory Azure SQL Database %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryLinkedServiceAzureSQLDatabaseRead(d, meta)
}
