// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2022-10-01-preview/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ServiceBusNamespaceCustomerManagedKeyResource struct{}

type ServiceBusNamespaceCustomerManagedKeyModel struct {
	NamespaceID                     string `tfschema:"namespace_id"`
	KeyVaultKeyID                   string `tfschema:"key_vault_key_id"`
	InfrastructureEncryptionEnabled bool   `tfschema:"infrastructure_encryption_enabled"`
}

var _ sdk.ResourceWithUpdate = ServiceBusNamespaceCustomerManagedKeyResource{}

func (r ServiceBusNamespaceCustomerManagedKeyResource) ModelObject() interface{} {
	return &ServiceBusNamespaceCustomerManagedKeyModel{}
}

func (r ServiceBusNamespaceCustomerManagedKeyResource) ResourceType() string {
	return "azurerm_servicebus_namespace_customer_managed_key"
}

func (r ServiceBusNamespaceCustomerManagedKeyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return namespaces.ValidateNamespaceID
}

func (r ServiceBusNamespaceCustomerManagedKeyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"namespace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: namespaces.ValidateNamespaceID,
		},

		"key_vault_key_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
		},

		"infrastructure_encryption_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},
	}
}

func (r ServiceBusNamespaceCustomerManagedKeyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ServiceBusNamespaceCustomerManagedKeyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceBus.NamespacesClient
			var cmk ServiceBusNamespaceCustomerManagedKeyModel

			if err := metadata.Decode(&cmk); err != nil {
				return err
			}

			id, err := namespaces.ParseNamespaceID(cmk.NamespaceID)
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", *id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` is nil", *id)
			}

			if resp.Model.Properties != nil && resp.Model.Properties.Encryption != nil && resp.Model.Properties.Encryption.KeyVaultProperties != nil && len(*resp.Model.Properties.Encryption.KeyVaultProperties) > 0 {
				return metadata.ResourceRequiresImport(r.ResourceType(), *id)
			}

			payload := resp.Model

			keyId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(cmk.KeyVaultKeyID)
			if err != nil {
				return err
			}

			payload.Properties.Encryption = &namespaces.Encryption{
				RequireInfrastructureEncryption: pointer.To(cmk.InfrastructureEncryptionEnabled),
				KeySource:                       pointer.To(namespaces.KeySourceMicrosoftPointKeyVault),
				KeyVaultProperties: &[]namespaces.KeyVaultProperties{
					{
						KeyName:     pointer.To(keyId.Name),
						KeyVersion:  pointer.To(keyId.Version),
						KeyVaultUri: pointer.To(keyId.KeyVaultBaseUrl),
					},
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("creating Customer Managed Key for %s: %+v", *id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ServiceBusNamespaceCustomerManagedKeyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceBus.NamespacesClient

			id, err := namespaces.ParseNamespaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var state ServiceBusNamespaceCustomerManagedKeyModel
			state.NamespaceID = id.ID()

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil && props.Encryption != nil {
					encryption := props.Encryption
					if keyVaultProperties := encryption.KeyVaultProperties; keyVaultProperties != nil && len(*keyVaultProperties) > 0 {
						keyVaultKeyId, err := keyVaultParse.NewNestedItemID(pointer.From((*keyVaultProperties)[0].KeyVaultUri), keyVaultParse.NestedItemTypeKey, pointer.From((*keyVaultProperties)[0].KeyName), pointer.From((*keyVaultProperties)[0].KeyVersion))
						if err != nil {
							return fmt.Errorf("parsing `key_vault_key_id`: %+v", err)
						}
						state.KeyVaultKeyID = keyVaultKeyId.ID()
					}
					state.InfrastructureEncryptionEnabled = pointer.From(encryption.RequireInfrastructureEncryption)
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (r ServiceBusNamespaceCustomerManagedKeyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceBus.NamespacesClient
			var cmk ServiceBusNamespaceCustomerManagedKeyModel

			if err := metadata.Decode(&cmk); err != nil {
				return err
			}

			id, err := namespaces.ParseNamespaceID(cmk.NamespaceID)
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", *id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` is nil", *id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` is nil", *id)
			}

			if resp.Model.Properties.Encryption == nil || resp.Model.Properties.Encryption.KeyVaultProperties == nil || len(*resp.Model.Properties.Encryption.KeyVaultProperties) == 0 {
				return fmt.Errorf("retrieving %s: Customer Managed Key was not found", *id)
			}

			payload := resp.Model

			if metadata.ResourceData.HasChange("key_vault_key_id") {
				keyId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(cmk.KeyVaultKeyID)
				if err != nil {
					return err
				}

				(*payload.Properties.Encryption.KeyVaultProperties)[0].KeyName = pointer.To(keyId.Name)
				(*payload.Properties.Encryption.KeyVaultProperties)[0].KeyVersion = pointer.To(keyId.Version)
				(*payload.Properties.Encryption.KeyVaultProperties)[0].KeyVaultUri = pointer.To(keyId.KeyVaultBaseUrl)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating Customer Managed Key for %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ServiceBusNamespaceCustomerManagedKeyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			log.Printf(`[INFO] Customer Managed Keys cannot be removed from Servicebus Namespaces once added. To remove the Customer Managed Key, delete and recreate the parent Servicebus Namespace`)
			return nil
		},
	}
}
