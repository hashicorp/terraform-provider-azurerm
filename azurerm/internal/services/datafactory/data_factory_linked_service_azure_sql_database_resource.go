package datafactory

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataFactoryLinkedServiceAzureSQLDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryLinkedServiceAzureSQLDatabaseCreateUpdate,
		Read:   resourceDataFactoryLinkedServiceAzureSQLDatabaseRead,
		Update: resourceDataFactoryLinkedServiceAzureSQLDatabaseCreateUpdate,
		Delete: resourceDataFactoryLinkedServiceAzureSQLDatabaseDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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

			"data_factory_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataFactoryName(),
			},

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/5788
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

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
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	dataFactoryName := d.Get("data_factory_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Data Factory Linked Service AzureSQLDatabase %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_linked_service_azure_sql_database", *existing.ID)
		}
	}

	sqlDatabaseProperties := &datafactory.AzureSQLDatabaseLinkedServiceTypeProperties{}

	if v, ok := d.GetOk("connection_string"); ok {
		sqlDatabaseProperties.ConnectionString = &datafactory.SecureString{
			Value: utils.String(v.(string)),
			Type:  datafactory.TypeSecureString,
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
		azureSQLDatabaseLinkedService.Parameters = expandDataFactoryParameters(v.(map[string]interface{}))
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

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, dataFactoryName, name, linkedService, ""); err != nil {
		return fmt.Errorf("Error creating/updating Data Factory Linked Service AzureSQLDatabase %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Factory Linked Service AzureSQLDatabase %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Data Factory Linked Service AzureSQLDatabase %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

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

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Data Factory Linked Service AzureSQLDatabase %q (Data Factory %q / Resource Group %q): %+v", id.Name, id.FactoryName, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("data_factory_name", id.FactoryName)

	sql, ok := resp.Properties.AsAzureSQLDatabaseLinkedService()
	if !ok {
		return fmt.Errorf("Error classifiying Data Factory Linked Service AzureSQLDatabase %q (Data Factory %q / Resource Group %q): Expected: %q Received: %q", id.Name, id.FactoryName, id.ResourceGroup, datafactory.TypeBasicLinkedServiceTypeAzureSQLDatabase, *resp.Type)
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
		return fmt.Errorf("Error setting `annotations`: %+v", err)
	}

	parameters := flattenDataFactoryParameters(sql.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("Error setting `parameters`: %+v", err)
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
			return fmt.Errorf("deleting Data Factory Linked Service AzureSQLDatabase %q (Data Factory %q / Resource Group %q): %+v", id.Name, id.FactoryName, id.ResourceGroup, err)
		}
	}

	return nil
}
