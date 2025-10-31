// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type LinkedServiceSqlManagedInstanceResource struct{}

type LinkedServiceSqlManagedInstanceModel struct {
	Name                     string                           `tfschema:"name"`
	DataFactoryID            string                           `tfschema:"data_factory_id"`
	Annotations              []string                         `tfschema:"annotations"`
	ConnectionString         string                           `tfschema:"connection_string"`
	Description              string                           `tfschema:"description"`
	IntegrationRuntimeName   string                           `tfschema:"integration_runtime_name"`
	KeyVaultConnectionString []KeyVaultConnectionStringConfig `tfschema:"key_vault_connection_string"`
	KeyVaultPassword         []KeyVaultPasswordConfig         `tfschema:"key_vault_password"`
	Parameters               map[string]string                `tfschema:"parameters"`
	ServicePrincipalID       string                           `tfschema:"service_principal_id"`
	ServicePrincipalKey      string                           `tfschema:"service_principal_key"`
	Tenant                   string                           `tfschema:"tenant"`
}

type KeyVaultConnectionStringConfig struct {
	LinkedServiceName string `tfschema:"linked_service_name"`
	SecretName        string `tfschema:"secret_name"`
}

type KeyVaultPasswordConfig struct {
	LinkedServiceName string `tfschema:"linked_service_name"`
	SecretName        string `tfschema:"secret_name"`
}

var _ sdk.ResourceWithUpdate = LinkedServiceSqlManagedInstanceResource{}

func (r LinkedServiceSqlManagedInstanceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
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

		"annotations": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
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

		"integration_runtime_name": {
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

		"parameters": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"service_principal_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"service_principal_key", "tenant"},
		},

		"service_principal_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"service_principal_id", "tenant"},
		},

		"tenant": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"service_principal_id", "service_principal_key"},
		},
	}
}

func (r LinkedServiceSqlManagedInstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LinkedServiceSqlManagedInstanceResource) ModelObject() interface{} {
	return &LinkedServiceSqlManagedInstanceModel{}
}

func (r LinkedServiceSqlManagedInstanceResource) ResourceType() string {
	return "azurerm_data_factory_linked_service_sql_managed_instance"
}

func (r LinkedServiceSqlManagedInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.LinkedServicesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config LinkedServiceSqlManagedInstanceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			dataFactoryId, err := factories.ParseFactoryID(config.DataFactoryID)
			if err != nil {
				return fmt.Errorf("parsing Data Factory ID %q: %+v", config.DataFactoryID, err)
			}

			id := linkedservices.NewLinkedServiceID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, config.Name)

			existing, err := client.Get(ctx, id, linkedservices.DefaultGetOperationOptions())
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_data_factory_linked_service_sql_managed_instance", id.ID())
			}

			parameters := &linkedservices.AzureSqlMILinkedService{
				Description: pointer.To(config.Description),
				TypeProperties: linkedservices.AzureSqlMILinkedServiceTypeProperties{
					Password: expandKeyVaultPasswordFromConfig(config.KeyVaultPassword),
				},
			}

			if config.ConnectionString != "" {
				connStr := interface{}(config.ConnectionString)
				parameters.TypeProperties.ConnectionString = &connStr
			}

			if len(config.KeyVaultConnectionString) > 0 {
				keyVaultConnStr := expandKeyVaultConnectionStringFromConfig(config.KeyVaultConnectionString)
				parameters.TypeProperties.ConnectionString = &keyVaultConnStr
			}

			if config.ServicePrincipalID != "" {
				spID := interface{}(config.ServicePrincipalID)
				parameters.TypeProperties.ServicePrincipalId = &spID
			}

			if config.ServicePrincipalKey != "" {
				secureString := linkedservices.SecureString{
					Value: config.ServicePrincipalKey,
				}
				parameters.TypeProperties.ServicePrincipalKey = secureString
				parameters.TypeProperties.ServicePrincipalCredential = secureString
			}

			if config.Tenant != "" {
				parameters.TypeProperties.Tenant = pointer.To(interface{}(config.Tenant))
			}

			if len(config.Parameters) > 0 {
				parameterSpec := make(map[string]linkedservices.ParameterSpecification)
				for key, value := range config.Parameters {
					val := interface{}(value)
					parameterSpec[key] = linkedservices.ParameterSpecification{
						Type:         linkedservices.ParameterTypeString,
						DefaultValue: &val,
					}
				}
				parameters.Parameters = &parameterSpec
			}

			if config.IntegrationRuntimeName != "" {
				parameters.ConnectVia = &linkedservices.IntegrationRuntimeReference{
					Type:          linkedservices.IntegrationRuntimeReferenceTypeIntegrationRuntimeReference,
					ReferenceName: config.IntegrationRuntimeName,
				}
			}

			if config.IntegrationRuntimeName != "" {
				parameters.ConnectVia = &linkedservices.IntegrationRuntimeReference{
					Type:          linkedservices.IntegrationRuntimeReferenceTypeIntegrationRuntimeReference,
					ReferenceName: config.IntegrationRuntimeName,
				}
			}

			if len(config.Annotations) > 0 {
				annotations := make([]interface{}, len(config.Annotations))
				for i, v := range config.Annotations {
					annotations[i] = v
				}
				parameters.Annotations = &annotations
			}

			linkedService := linkedservices.LinkedServiceResource{
				Properties: parameters,
			}

			if _, err := client.CreateOrUpdate(ctx, id, linkedService, linkedservices.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r LinkedServiceSqlManagedInstanceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.LinkedServicesClient

			id, err := linkedservices.ParseLinkedServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id, linkedservices.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			existing, ok := resp.Model.Properties.(linkedservices.AzureSqlMILinkedService)
			if !ok {
				return fmt.Errorf("classifying %s: Expected: %q Received: %T", id, "AzureSqlMI", resp.Model.Properties)
			}

			if existing.Type != "AzureSqlMI" {
				return fmt.Errorf("classifying %s: Expected: %q Received: %q", id, "AzureSqlMI", existing.Type)
			}

			state := LinkedServiceSqlManagedInstanceModel{
				Name:          id.LinkedServiceName,
				DataFactoryID: factories.NewFactoryID(id.SubscriptionId, id.ResourceGroupName, id.FactoryName).ID(),
			}

			state.Description = pointer.From(existing.Description)

			props := existing.TypeProperties

			if props.ConnectionString != nil {
				if val, ok := (*props.ConnectionString).(map[string]interface{}); ok {
					state.KeyVaultConnectionString = flattenKeyVaultConnectionStringToConfig(val)
				}
				if val, ok := (*props.ConnectionString).(string); ok {
					state.ConnectionString = val
				}
			}

			if props.ServicePrincipalId != nil {
				if id, ok := (*props.ServicePrincipalId).(string); ok {
					state.ServicePrincipalID = id
				}
			}

			if props.Tenant != nil {
				if tenant, ok := (*props.Tenant).(string); ok {
					state.Tenant = tenant
				}
			}

			if v, exists := metadata.ResourceData.GetOk("service_principal_key"); exists && v.(string) != "" {
				state.ServicePrincipalKey = v.(string)
			}

			state.KeyVaultPassword = flattenKeyVaultPasswordToConfig(props.Password)

			if existing.Annotations != nil {
				annotations := make([]string, 0)
				for _, annotation := range *existing.Annotations {
					if str, ok := annotation.(string); ok {
						annotations = append(annotations, str)
					}
				}
				state.Annotations = annotations
			}

			if existing.Parameters != nil {
				parameters := make(map[string]string)
				for key, param := range *existing.Parameters {
					if param.DefaultValue != nil {
						if str, ok := (*param.DefaultValue).(string); ok {
							parameters[key] = str
						}
					}
				}
				state.Parameters = parameters
			}

			if connectVia := existing.ConnectVia; connectVia != nil {
				state.IntegrationRuntimeName = connectVia.ReferenceName
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LinkedServiceSqlManagedInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.LinkedServicesClient

			id, err := linkedservices.ParseLinkedServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config LinkedServiceSqlManagedInstanceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id, linkedservices.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model.Properties` was nil", *id)
			}

			payload := linkedservices.LinkedServiceResource{
				Properties: &linkedservices.AzureSqlMILinkedService{
					Description: pointer.To(config.Description),
					TypeProperties: linkedservices.AzureSqlMILinkedServiceTypeProperties{
						Password: expandKeyVaultPasswordFromConfig(config.KeyVaultPassword),
					},
				},
			}
			props := payload.Properties.(*linkedservices.AzureSqlMILinkedService)

			if metadata.ResourceData.HasChanges("connection_string", "key_vault_connection_string") {
				if config.ConnectionString != "" {
					connStr := interface{}(config.ConnectionString)
					props.TypeProperties.ConnectionString = &connStr
				}

				if len(config.KeyVaultConnectionString) > 0 {
					keyVaultConnStr := expandKeyVaultConnectionStringFromConfig(config.KeyVaultConnectionString)
					props.TypeProperties.ConnectionString = &keyVaultConnStr
				}
			}

			if metadata.ResourceData.HasChanges("service_principal_id", "service_principal_key", "tenant") {
				if config.ServicePrincipalID != "" {
					spID := interface{}(config.ServicePrincipalID)
					props.TypeProperties.ServicePrincipalId = &spID
				}

				if config.ServicePrincipalKey != "" {
					secureString := linkedservices.SecureString{
						Value: config.ServicePrincipalKey,
					}
					props.TypeProperties.ServicePrincipalKey = secureString
					props.TypeProperties.ServicePrincipalCredential = secureString
				}

				if config.Tenant != "" {
					tenant := interface{}(config.Tenant)
					props.TypeProperties.Tenant = &tenant
				}
			}

			if metadata.ResourceData.HasChange("parameters") {
				if len(config.Parameters) > 0 {
					parameterSpec := make(map[string]linkedservices.ParameterSpecification)
					for key, value := range config.Parameters {
						val := interface{}(value)
						parameterSpec[key] = linkedservices.ParameterSpecification{
							Type:         linkedservices.ParameterTypeString,
							DefaultValue: &val,
						}
					}
					props.Parameters = &parameterSpec
				}
			}

			if metadata.ResourceData.HasChange("integration_runtime_name") {
				if config.IntegrationRuntimeName != "" {
					props.ConnectVia = &linkedservices.IntegrationRuntimeReference{
						Type:          linkedservices.IntegrationRuntimeReferenceTypeIntegrationRuntimeReference,
						ReferenceName: config.IntegrationRuntimeName,
					}
				}
			}

			if metadata.ResourceData.HasChange("annotations") {
				if len(config.Annotations) > 0 {
					annotations := make([]interface{}, len(config.Annotations))
					for i, v := range config.Annotations {
						annotations[i] = v
					}
					props.Annotations = &annotations
				}
			}

			if _, err := client.CreateOrUpdate(ctx, *id, payload, linkedservices.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r LinkedServiceSqlManagedInstanceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.LinkedServicesClient

			id, err := linkedservices.ParseLinkedServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Delete(ctx, *id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("deleting %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (r LinkedServiceSqlManagedInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return linkedservices.ValidateLinkedServiceID
}

func expandKeyVaultConnectionStringFromConfig(input []KeyVaultConnectionStringConfig) interface{} {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	return &linkedservices.AzureKeyVaultSecretReference{
		SecretName: config.SecretName,
		Store: linkedservices.LinkedServiceReference{
			Type:          linkedservices.TypeLinkedServiceReference,
			ReferenceName: config.LinkedServiceName,
		},
	}
}

func expandKeyVaultPasswordFromConfig(input []KeyVaultPasswordConfig) *linkedservices.AzureKeyVaultSecretReference {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	return &linkedservices.AzureKeyVaultSecretReference{
		SecretName: config.SecretName,
		Store: linkedservices.LinkedServiceReference{
			Type:          linkedservices.TypeLinkedServiceReference,
			ReferenceName: config.LinkedServiceName,
		},
	}
}

func flattenKeyVaultConnectionStringToConfig(input interface{}) []KeyVaultConnectionStringConfig {
	if input == nil {
		return []KeyVaultConnectionStringConfig{}
	}

	flattened := flattenAzureKeyVaultConnectionString(input.(map[string]interface{}))
	if len(flattened) == 0 {
		return []KeyVaultConnectionStringConfig{}
	}

	configMap := flattened[0].(map[string]interface{})
	return []KeyVaultConnectionStringConfig{{
		LinkedServiceName: configMap["linked_service_name"].(string),
		SecretName:        configMap["secret_name"].(string),
	}}
}

func flattenKeyVaultPasswordToConfig(input *linkedservices.AzureKeyVaultSecretReference) []KeyVaultPasswordConfig {
	if input == nil {
		return []KeyVaultPasswordConfig{}
	}

	config := KeyVaultPasswordConfig{}

	if secretName, ok := input.SecretName.(string); ok {
		config.SecretName = secretName
	}

	config.LinkedServiceName = input.Store.ReferenceName

	return []KeyVaultPasswordConfig{config}
}
