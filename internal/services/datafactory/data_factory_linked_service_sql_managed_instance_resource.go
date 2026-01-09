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
	Parameters               map[string]interface{}           `tfschema:"parameters"`
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
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
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
				return err
			}

			id := linkedservices.NewLinkedServiceID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, config.Name)

			existing, err := client.Get(ctx, id, linkedservices.DefaultGetOperationOptions())
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(r.ResourceType(), id.ID())
			}

			sqlMILinkedService := &linkedservices.AzureSqlMILinkedService{
				Description: pointer.To(config.Description),
				TypeProperties: linkedservices.AzureSqlMILinkedServiceTypeProperties{
					Password: expandLinkedServiceSqlManagedInstanceKeyVaultPassword(config.KeyVaultPassword),
				},
			}

			if config.ConnectionString != "" {
				sqlMILinkedService.TypeProperties.ConnectionString = pointer.To(interface{}(config.ConnectionString))
			}

			if len(config.KeyVaultConnectionString) > 0 {
				sqlMILinkedService.TypeProperties.ConnectionString = pointer.To(expandLinkedServiceSqlManagedInstanceKeyVaultConnectionString(config.KeyVaultConnectionString))
			}

			if config.ServicePrincipalID != "" {
				sqlMILinkedService.TypeProperties.ServicePrincipalId = pointer.To(interface{}(config.ServicePrincipalID))
			}

			if config.ServicePrincipalKey != "" {
				secureString := linkedservices.SecureString{
					Value: config.ServicePrincipalKey,
				}
				sqlMILinkedService.TypeProperties.ServicePrincipalKey = secureString
				sqlMILinkedService.TypeProperties.ServicePrincipalCredential = secureString
			}

			if config.Tenant != "" {
				sqlMILinkedService.TypeProperties.Tenant = pointer.To(interface{}(config.Tenant))
			}

			if len(config.Parameters) > 0 {
				sqlMILinkedService.Parameters = expandLinkedServiceSqlManagedInstanceParameters(config.Parameters)
			}

			if config.IntegrationRuntimeName != "" {
				sqlMILinkedService.ConnectVia = &linkedservices.IntegrationRuntimeReference{
					Type:          linkedservices.IntegrationRuntimeReferenceTypeIntegrationRuntimeReference,
					ReferenceName: config.IntegrationRuntimeName,
				}
			}

			if len(config.Annotations) > 0 {
				sqlMILinkedService.Annotations = expandLinkedServiceSqlManagedInstanceAnnotations(config.Annotations)
			}

			linkedService := linkedservices.LinkedServiceResource{
				Properties: sqlMILinkedService,
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

			sqlMILinkedService, ok := resp.Model.Properties.(linkedservices.AzureSqlMILinkedService)
			if !ok {
				return fmt.Errorf("classifying %s: Expected: %q Received: %T", id, "AzureSqlMI", resp.Model.Properties)
			}

			if sqlMILinkedService.Type != "AzureSqlMI" {
				return fmt.Errorf("classifying %s: Expected: %q Received: %q", id, "AzureSqlMI", sqlMILinkedService.Type)
			}

			state := LinkedServiceSqlManagedInstanceModel{
				Name:          id.LinkedServiceName,
				DataFactoryID: factories.NewFactoryID(id.SubscriptionId, id.ResourceGroupName, id.FactoryName).ID(),
			}

			state.Description = pointer.From(sqlMILinkedService.Description)

			props := sqlMILinkedService.TypeProperties

			state.ConnectionString = ""
			if props.ConnectionString != nil {
				if val, ok := pointer.From(props.ConnectionString).(map[string]interface{}); ok {
					state.KeyVaultConnectionString = flattenLinkedServiceSqlManagedInstanceKeyVaultConnectionString(val)
				} else {
					state.ConnectionString = pointer.From(props.ConnectionString).(string)
				}
			}

			state.ServicePrincipalID = ""
			if v := pointer.From(props.ServicePrincipalId); v != nil {
				state.ServicePrincipalID = v.(string)
			}

			state.Tenant = ""
			if v := pointer.From(props.Tenant); v != nil {
				state.Tenant = v.(string)
			}

			// not returned from API
			state.ServicePrincipalKey = ""
			if v, exists := metadata.ResourceData.GetOk("service_principal_key"); exists {
				state.ServicePrincipalKey = v.(string)
			}

			state.KeyVaultPassword = flattenLinkedServiceSqlManagedInstanceKeyVaultPassword(props.Password)
			state.Annotations = flattenLinkedServiceSqlManagedInstanceAnnotations(sqlMILinkedService.Annotations)
			state.Parameters = flattenLinkedServiceSqlManagedInstanceParameters(sqlMILinkedService.Parameters)

			state.IntegrationRuntimeName = ""
			if connectVia := sqlMILinkedService.ConnectVia; connectVia != nil {
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

			props := existing.Model.Properties.(linkedservices.AzureSqlMILinkedService)
			typeProps := props.TypeProperties

			if metadata.ResourceData.HasChanges("connection_string", "key_vault_connection_string") {
				typeProps.ConnectionString = pointer.To(interface{}(config.ConnectionString))

				if len(config.KeyVaultConnectionString) > 0 {
					typeProps.ConnectionString = pointer.To(expandLinkedServiceSqlManagedInstanceKeyVaultConnectionString(config.KeyVaultConnectionString))
				}
			}

			if metadata.ResourceData.HasChange("key_vault_password") {
				typeProps.Password = expandLinkedServiceSqlManagedInstanceKeyVaultPassword(config.KeyVaultPassword)
			}

			if metadata.ResourceData.HasChange("service_principal_id") {
				typeProps.ServicePrincipalId = pointer.To(interface{}(config.ServicePrincipalID))
			}

			if metadata.ResourceData.HasChange("service_principal_key") && config.ServicePrincipalKey != "" {
				secureString := linkedservices.SecureString{
					Value: config.ServicePrincipalKey,
				}
				typeProps.ServicePrincipalKey = secureString
				typeProps.ServicePrincipalCredential = secureString
			} else {
				typeProps.ServicePrincipalKey = nil
				typeProps.ServicePrincipalCredential = nil
			}

			if metadata.ResourceData.HasChange("tenant") {
				typeProps.Tenant = pointer.To(interface{}(config.Tenant))
			}

			if metadata.ResourceData.HasChange("parameters") {
				props.Parameters = expandLinkedServiceSqlManagedInstanceParameters(config.Parameters)
			}

			if metadata.ResourceData.HasChange("integration_runtime_name") {
				if config.IntegrationRuntimeName != "" {
					props.ConnectVia = &linkedservices.IntegrationRuntimeReference{
						Type:          linkedservices.IntegrationRuntimeReferenceTypeIntegrationRuntimeReference,
						ReferenceName: config.IntegrationRuntimeName,
					}
				} else {
					props.ConnectVia = nil
				}
			}

			if metadata.ResourceData.HasChange("description") {
				props.Description = pointer.To(config.Description)
			}

			if metadata.ResourceData.HasChange("annotations") {
				props.Annotations = expandLinkedServiceSqlManagedInstanceAnnotations(config.Annotations)
			}

			props.TypeProperties = typeProps
			payload := linkedservices.LinkedServiceResource{
				Properties: &props,
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

			_, err = client.Delete(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r LinkedServiceSqlManagedInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return linkedservices.ValidateLinkedServiceID
}

func expandLinkedServiceSqlManagedInstanceKeyVaultConnectionString(input []KeyVaultConnectionStringConfig) interface{} {
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

func expandLinkedServiceSqlManagedInstanceKeyVaultPassword(input []KeyVaultPasswordConfig) *linkedservices.AzureKeyVaultSecretReference {
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

func flattenLinkedServiceSqlManagedInstanceKeyVaultConnectionString(input interface{}) []KeyVaultConnectionStringConfig {
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

func flattenLinkedServiceSqlManagedInstanceKeyVaultPassword(input *linkedservices.AzureKeyVaultSecretReference) []KeyVaultPasswordConfig {
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

func expandLinkedServiceSqlManagedInstanceParameters(input map[string]interface{}) *map[string]linkedservices.ParameterSpecification {
	if len(input) == 0 {
		return nil
	}

	parameterSpec := make(map[string]linkedservices.ParameterSpecification)
	for key, value := range input {
		parameterSpec[key] = linkedservices.ParameterSpecification{
			Type:         linkedservices.ParameterTypeString,
			DefaultValue: pointer.To(value),
		}
	}

	return &parameterSpec
}

func flattenLinkedServiceSqlManagedInstanceParameters(input *map[string]linkedservices.ParameterSpecification) map[string]interface{} {
	output := make(map[string]interface{})
	if input == nil {
		return output
	}

	for key, param := range *input {
		if param.DefaultValue != nil {
			if str, ok := pointer.From(param.DefaultValue).(string); ok {
				output[key] = str
			}
		}
	}

	return output
}

func expandLinkedServiceSqlManagedInstanceAnnotations(input []string) *[]interface{} {
	if len(input) == 0 {
		return nil
	}

	annotations := make([]interface{}, len(input))
	for i, v := range input {
		annotations[i] = v
	}

	return &annotations
}

func flattenLinkedServiceSqlManagedInstanceAnnotations(input *[]interface{}) []string {
	annotations := make([]string, 0)
	if input == nil {
		return annotations
	}

	for _, annotation := range *input {
		if str, ok := annotation.(string); ok {
			annotations = append(annotations, str)
		}
	}

	return annotations
}
