// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataFactoryLinkedServiceAzureSQLDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryLinkedServiceAzureSQLDatabaseCreateUpdate,
		Read:   resourceDataFactoryLinkedServiceAzureSQLDatabaseRead,
		Update: resourceDataFactoryLinkedServiceAzureSQLDatabaseCreateUpdate,
		Delete: resourceDataFactoryLinkedServiceAzureSQLDatabaseDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.LinkedServiceID(id)
			return err
		}, importDataFactoryLinkedService(datafactory.TypeBasicLinkedServiceTypeAzureSQLDatabase)),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LinkedServiceDatasetName,
			},

			"data_factory_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: factories.ValidateFactoryID,
			},

			"connection_string": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ExactlyOneOf:     []string{"connection_string", "key_vault_connection_string"},
				DiffSuppressFunc: azureRmDataFactoryLinkedServiceConnectionStringDiff,
				ValidateFunc:     validation.StringIsNotEmpty,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"key_vault_connection_string": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				ExactlyOneOf: []string{"connection_string", "key_vault_connection_string"},
				MaxItems:     1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"linked_service_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"secret_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"key_vault_password": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"linked_service_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"secret_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"use_managed_identity": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
				ConflictsWith: []string{
					"service_principal_id",
				},
			},

			"service_principal_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				RequiredWith: []string{"service_principal_key"},
				ConflictsWith: []string{
					"use_managed_identity",
				},
			},

			"service_principal_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				RequiredWith: []string{"service_principal_id"},
			},

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"integration_runtime_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"parameters": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"annotations": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"additional_properties": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func resourceDataFactoryLinkedServiceAzureSQLDatabaseCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Data Factory Azure SQL Database %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_linked_service_azure_sql_database", id.ID())
		}
	}

	sqlDatabaseProperties := &datafactory.AzureSQLDatabaseLinkedServiceTypeProperties{}

	if v, ok := d.GetOk("connection_string"); ok {
		if d.Get("use_managed_identity").(bool) {
			sqlDatabaseProperties.ConnectionString = utils.String(v.(string))
		} else {
			sqlDatabaseProperties.ConnectionString = &datafactory.SecureString{
				Value: utils.String(v.(string)),
				Type:  datafactory.TypeSecureString,
			}
		}
	}

	if v, ok := d.GetOk("key_vault_connection_string"); ok {
		sqlDatabaseProperties.ConnectionString = expandAzureKeyVaultSecretReference(v.([]interface{}))
	}

	if d.Get("use_managed_identity").(bool) {
		sqlDatabaseProperties.Tenant = utils.String(d.Get("tenant_id").(string))
	} else {
		secureString := datafactory.SecureString{
			Value: utils.String(d.Get("service_principal_key").(string)),
			Type:  datafactory.TypeSecureString,
		}

		sqlDatabaseProperties.ServicePrincipalID = utils.String(d.Get("service_principal_id").(string))
		sqlDatabaseProperties.Tenant = utils.String(d.Get("tenant_id").(string))
		sqlDatabaseProperties.ServicePrincipalKey = &secureString
	}

	if v, ok := d.GetOk("key_vault_password"); ok {
		password := v.([]interface{})
		sqlDatabaseProperties.Password = expandAzureKeyVaultSecretReference(password)
	}

	azureSQLDatabaseLinkedService := &datafactory.AzureSQLDatabaseLinkedService{
		Description: utils.String(d.Get("description").(string)),
		AzureSQLDatabaseLinkedServiceTypeProperties: sqlDatabaseProperties,
		Type: datafactory.TypeBasicLinkedServiceTypeAzureSQLDatabase,
	}

	if v, ok := d.GetOk("parameters"); ok {
		azureSQLDatabaseLinkedService.Parameters = expandLinkedServiceParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		azureSQLDatabaseLinkedService.ConnectVia = expandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		azureSQLDatabaseLinkedService.AdditionalProperties = v.(map[string]interface{})
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		azureSQLDatabaseLinkedService.Annotations = &annotations
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

func resourceDataFactoryLinkedServiceAzureSQLDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LinkedServiceID(d.Id())
	if err != nil {
		return err
	}

	dataFactoryId := factories.NewFactoryID(id.SubscriptionId, id.ResourceGroup, id.FactoryName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Data Factory Azure SQL Database %s: %+v", *id, err)
	}

	d.Set("name", resp.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	sql, ok := resp.Properties.AsAzureSQLDatabaseLinkedService()
	if !ok {
		return fmt.Errorf("classifying Data Factory Azure SQL Database %s: Expected: %q Received: %q", *id, datafactory.TypeBasicLinkedServiceTypeAzureSQLDatabase, *resp.Type)
	}

	if sql != nil {
		if sql.Tenant != nil {
			d.Set("tenant_id", sql.Tenant)
		}

		if sql.ServicePrincipalID != nil {
			d.Set("service_principal_id", sql.ServicePrincipalID)
			d.Set("use_managed_identity", false)
		} else {
			d.Set("use_managed_identity", true)
		}
	}

	if sql.ConnectionString != nil {
		if val, ok := sql.ConnectionString.(map[string]interface{}); ok {
			if val["type"] != "SecureString" {
				if err := d.Set("key_vault_connection_string", flattenAzureKeyVaultConnectionString(val)); err != nil {
					return fmt.Errorf("setting `key_vault_connection_string`: %+v", err)
				}
			}
		}
	}

	d.Set("additional_properties", sql.AdditionalProperties)
	d.Set("description", sql.Description)

	if password := sql.Password; password != nil {
		if keyVaultPassword, ok := password.AsAzureKeyVaultSecretReference(); ok {
			if err := d.Set("key_vault_password", flattenAzureKeyVaultSecretReference(keyVaultPassword)); err != nil {
				return fmt.Errorf("setting `key_vault_password`: %+v", err)
			}
		}
	}

	annotations := flattenDataFactoryAnnotations(sql.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	parameters := flattenLinkedServiceParameters(sql.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	if connectVia := sql.ConnectVia; connectVia != nil {
		if connectVia.ReferenceName != nil {
			d.Set("integration_runtime_name", connectVia.ReferenceName)
		}
	}

	return nil
}

func resourceDataFactoryLinkedServiceAzureSQLDatabaseDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LinkedServiceID(d.Id())
	if err != nil {
		return err
	}

	response, err := client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("deleting Data Factory Azure SQL Database %s: %+v", *id, err)
		}
	}

	return nil
}
