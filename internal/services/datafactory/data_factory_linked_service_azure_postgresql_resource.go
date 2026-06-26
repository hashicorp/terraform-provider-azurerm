// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource                  = LinkedServiceAzurePostgreSQLResource{}
	_ sdk.ResourceWithIdentity      = LinkedServiceAzurePostgreSQLResource{}
	_ sdk.ResourceWithCustomizeDiff = LinkedServiceAzurePostgreSQLResource{}
	_ sdk.ResourceWithUpdate        = LinkedServiceAzurePostgreSQLResource{}
)

type LinkedServiceAzurePostgreSQLResource struct{}

type LinkedServiceAzurePostgreSQLResourceModel struct {
	Name                   string                   `tfschema:"name"`
	DataFactoryID          string                   `tfschema:"data_factory_id"`
	AuthenticationType     string                   `tfschema:"authentication_type"`
	Database               string                   `tfschema:"database_name"`
	Port                   int                      `tfschema:"port"`
	Server                 string                   `tfschema:"server"`
	Annotations            []string                 `tfschema:"annotations"`
	CredentialName         string                   `tfschema:"credential_name"`
	Description            string                   `tfschema:"description"`
	IntegrationRuntimeName string                   `tfschema:"integration_runtime_name"`
	KeyVaultPassword       []KeyVaultPasswordConfig `tfschema:"key_vault_password"`
	Parameters             map[string]interface{}   `tfschema:"parameters"`
	SslMode                string                   `tfschema:"ssl_mode"`
	Username               string                   `tfschema:"username"`
}

func (LinkedServiceAzurePostgreSQLResource) Arguments() map[string]*pluginsdk.Schema {
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

		// authentication_type should be supported by sdk but hasn't due to bug
		// TODO: support
		"authentication_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Basic",
				"SystemAssignedManagedIdentity",
				"UserAssignedManagedIdentity",
			}, false),
		},

		"database_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"port": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: azValidate.PortNumberOrZero,
		},

		"server": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"annotations": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"credential_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ConflictsWith: []string{
				"key_vault_password",
			},
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
			ConflictsWith: []string{
				"credential_name",
			},
		},

		"parameters": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"ssl_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"0",
				"1",
				"2",
				"3",
				"4",
				"5",
			}, false),
		},

		"username": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r LinkedServiceAzurePostgreSQLResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diff := metadata.ResourceDiff

			authenticationType := diff.Get("authentication_type").(string)

			if authenticationType == "Basic" {
				if _, ok := diff.GetOk("username"); !ok {
					return errors.New("`username` must be specified when `authentication_type` is `Basic`")
				}
				if _, ok := diff.GetOk("key_vault_password"); !ok {
					return errors.New("`key_vault_password` must be specified when `authentication_type` is `Basic`")
				}
			}

			if authenticationType == "UserAssignedManagedIdentity" {
				if _, ok := diff.GetOk("credential_name"); !ok {
					return errors.New("`credential_name` must be specified when `authentication_type` is `UserAssignedManagedIdentity`")
				}
			}

			return nil
		},
	}
}

func (LinkedServiceAzurePostgreSQLResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (LinkedServiceAzurePostgreSQLResource) ModelObject() interface{} {
	return &LinkedServiceAzurePostgreSQLResourceModel{}
}

func (LinkedServiceAzurePostgreSQLResource) ResourceType() string {
	return "azurerm_data_factory_linked_service_azure_postgresql"
}

func (r LinkedServiceAzurePostgreSQLResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.LinkedServicesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config LinkedServiceAzurePostgreSQLResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			dataFactoryId, err := linkedservices.ParseFactoryID(config.DataFactoryID)
			if err != nil {
				return err
			}

			id := linkedservices.NewLinkedServiceID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id, linkedservices.DefaultGetOperationOptions())
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}

			}

			azurePostgreSQLLinkedService := linkedservices.AzurePostgreSqlLinkedService{
				Annotations:    expandTypedLinkedServiceAnnotations(config.Annotations),
				ConnectVia:     expandTypedLinkedServiceIntegrationRuntimeName(config.IntegrationRuntimeName),
				Description:    pointer.To(config.Description),
				Parameters:     expandTypedLinkedServiceParameters(config.Parameters),
				TypeProperties: expandLinkedServiceAzurePostgreSQLTypeProperties(config),
				Version:        pointer.To("2.0"),
			}

			linkedService := linkedservices.LinkedServiceResource{
				Properties: azurePostgreSQLLinkedService,
			}

			if _, err := client.CreateOrUpdate(ctx, id, linkedService, linkedservices.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}
			return nil
		},
	}
}

func (LinkedServiceAzurePostgreSQLResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.LinkedServicesClient

			id, err := linkedservices.ParseLinkedServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config LinkedServiceAzurePostgreSQLResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			azurePostgreSQLLinkedService := linkedservices.AzurePostgreSqlLinkedService{
				Annotations:    expandTypedLinkedServiceAnnotations(config.Annotations),
				ConnectVia:     expandTypedLinkedServiceIntegrationRuntimeName(config.IntegrationRuntimeName),
				Description:    pointer.To(config.Description),
				Parameters:     expandTypedLinkedServiceParameters(config.Parameters),
				TypeProperties: expandLinkedServiceAzurePostgreSQLTypeProperties(config),
				Version:        pointer.To("2.0"),
			}

			linkedService := linkedservices.LinkedServiceResource{
				Properties: azurePostgreSQLLinkedService,
			}

			if _, err := client.CreateOrUpdate(ctx, *id, linkedService, linkedservices.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (LinkedServiceAzurePostgreSQLResource) Read() sdk.ResourceFunc {
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

			return LinkedServiceAzurePostgreSQLResource{}.flatten(metadata, id, resp.Model)
		},
	}
}

func (LinkedServiceAzurePostgreSQLResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.LinkedServicesClient

			id, err := linkedservices.ParseLinkedServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (LinkedServiceAzurePostgreSQLResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return linkedservices.ValidateLinkedServiceID
}

func (r LinkedServiceAzurePostgreSQLResource) Identity() resourceids.ResourceId {
	return &linkedservices.LinkedServiceId{}
}

func expandLinkedServiceAzurePostgreSQLTypeProperties(config LinkedServiceAzurePostgreSQLResourceModel) linkedservices.AzurePostgreSqlLinkedServiceTypeProperties {
	typeProperties := &linkedservices.AzurePostgreSqlLinkedServiceTypeProperties{
		Database: pointer.To(interface{}(config.Database)),
		Server:   pointer.To(interface{}(config.Server)),
		Port:     pointer.To(interface{}(config.Port)),
	}

	if config.SslMode != "" {
		typeProperties.SslMode = pointer.To(interface{}(config.SslMode))
	}
	if config.AuthenticationType == "Basic" {
		typeProperties.Password = expandTypedLinkedServiceKeyVaultPassword(config.KeyVaultPassword)
		typeProperties.Username = pointer.To(interface{}(config.Username))
	}

	if config.AuthenticationType == "UserAssignedManagedIdentity" {
		typeProperties.Credential = &linkedservices.CredentialReference{
			ReferenceName: config.CredentialName,
		}
	}

	return *typeProperties
}

func (LinkedServiceAzurePostgreSQLResource) flatten(metadata sdk.ResourceMetaData, id *linkedservices.LinkedServiceId, model *linkedservices.LinkedServiceResource) error {
	dataFactoryId := linkedservices.NewFactoryID(id.SubscriptionId, id.ResourceGroupName, id.FactoryName)
	state := LinkedServiceAzurePostgreSQLResourceModel{
		Name:          id.LinkedServiceName,
		DataFactoryID: dataFactoryId.ID(),
	}

	if model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	linkedService, ok := model.Properties.(linkedservices.AzurePostgreSqlLinkedService)
	if !ok {
		return fmt.Errorf("classifying %s: expected `LinkedServiceAzurePostgreSql` received `%T`", id, model.Properties)
	}

	props := linkedService.TypeProperties

	if database, ok := pointer.From(props.Database).(string); ok {
		state.Database = database
	}

	switch port := pointer.From(props.Port).(type) {
	case int:
		state.Port = port
	case int64:
		state.Port = int(port)
	case float64:
		state.Port = int(port)
	}

	if server, ok := pointer.From(props.Server).(string); ok {
		state.Server = server
	}

	state.Annotations = flattenTypedLinkedServiceAnnotations(linkedService.Annotations)
	if props.Credential != nil {
		state.CredentialName = props.Credential.ReferenceName
	}
	state.Description = pointer.From(linkedService.Description)
	state.IntegrationRuntimeName = flattenTypedLinkedServiceIntegrationRuntimeName(linkedService.ConnectVia)
	state.KeyVaultPassword = flattenTypedLinkedServiceKeyVaultPassword(props.Password)
	state.Parameters = flattenTypedLinkedServiceParameters(linkedService.Parameters)
	if props.SslMode != nil {
		if sslMode, ok := pointer.From(props.SslMode).(string); ok {
			state.SslMode = sslMode
		}
	}
	if props.Username != nil {
		if username, ok := pointer.From(props.Username).(string); ok {
			state.Username = username
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}
