package datafactory

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataFactoryCustomerManagedKeyResource struct{}

type DataFactoryCustomerManagedKeyModel struct {
	DataFactoryID          string `tfschema:"data_factory_id"`
	CustomerManagedKeyID   string `tfschema:"customer_managed_key_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

var _ sdk.ResourceWithUpdate = DataFactoryCustomerManagedKeyResource{}

func (r DataFactoryCustomerManagedKeyResource) ModelObject() interface{} {
	return &DataFactoryCustomerManagedKeyModel{}
}

func (r DataFactoryCustomerManagedKeyResource) ResourceType() string {
	return "azurerm_data_factory_customer_managed_key"
}

func (r DataFactoryCustomerManagedKeyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return factories.ValidateFactoryID
}

func (r DataFactoryCustomerManagedKeyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"data_factory_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},
		"customer_managed_key_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: keyVaultValidate.NestedItemId,
		},
		"user_assigned_identity_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},
	}
}

func (r DataFactoryCustomerManagedKeyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataFactoryCustomerManagedKeyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.Factories

			var customerManagedKey DataFactoryCustomerManagedKeyModel
			if err := metadata.Decode(&customerManagedKey); err != nil {
				return err
			}

			id, err := factories.ParseFactoryID(customerManagedKey.DataFactoryID)
			if err != nil {
				return err
			}

			dataFactory, err := client.Get(ctx, *id, factories.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(dataFactory.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if dataFactory.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if dataFactory.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model.Properties` was nil", id)
			}

			payload := dataFactory.Model
			if payload.Properties.Encryption != nil {
				return tf.ImportAsExistsError("azurerm_data_factory_customer_managed_key", id.ID())
			}

			keyVaultKey, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(customerManagedKey.CustomerManagedKeyID)
			if err != nil {
				return fmt.Errorf("could not parse Key Vault Key ID: %+v", err)
			}
			encryption := &factories.EncryptionConfiguration{
				VaultBaseURL: keyVaultKey.KeyVaultBaseUrl,
				KeyName:      keyVaultKey.Name,
				KeyVersion:   pointer.To(keyVaultKey.Version),
			}

			if identityId := customerManagedKey.UserAssignedIdentityID; identityId != "" {
				encryption.Identity = &factories.CMKIdentityDefinition{
					UserAssignedIdentity: pointer.To(identityId),
				}
			}

			payload.Properties.Encryption = encryption
			if _, err = client.CreateOrUpdate(ctx, *id, *payload, factories.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}

func (r DataFactoryCustomerManagedKeyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.Factories

			id, err := factories.ParseFactoryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			dataFactory, err := client.Get(ctx, *id, factories.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(dataFactory.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if dataFactory.Model == nil || dataFactory.Model.Properties == nil || dataFactory.Model.Properties.Encryption == nil {
				return fmt.Errorf("retrieving encryption for %s", id)
			}

			encryption := dataFactory.Model.Properties.Encryption

			customerManagedKeyId := ""
			customerManagedKeyIdentityId := ""
			if encryption.VaultBaseURL != "" && encryption.KeyName != "" && encryption.KeyVersion != nil {
				version := pointer.From(encryption.KeyVersion)

				keyId, err := keyVaultParse.NewNestedKeyID(encryption.VaultBaseURL, encryption.KeyName, version)
				if err != nil {
					return fmt.Errorf("parsing Nested Item ID: %+v", err)
				}
				customerManagedKeyId = keyId.ID()
			}

			if encIdentity := encryption.Identity; encIdentity != nil && encIdentity.UserAssignedIdentity != nil {
				parsed, err := commonids.ParseUserAssignedIdentityIDInsensitively(pointer.From(encIdentity.UserAssignedIdentity))
				if err != nil {
					return fmt.Errorf("parsing %q: %+v", *encIdentity.UserAssignedIdentity, err)
				}
				customerManagedKeyIdentityId = parsed.ID()
			}

			state := DataFactoryCustomerManagedKeyModel{
				DataFactoryID:          id.ID(),
				CustomerManagedKeyID:   customerManagedKeyId,
				UserAssignedIdentityID: customerManagedKeyIdentityId,
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DataFactoryCustomerManagedKeyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.Factories

			var customerManagedKey DataFactoryCustomerManagedKeyModel
			if err := metadata.Decode(&customerManagedKey); err != nil {
				return err
			}

			id, err := factories.ParseFactoryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			dataFactory, err := client.Get(ctx, *id, factories.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(dataFactory.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if dataFactory.Model == nil || dataFactory.Model.Properties == nil || dataFactory.Model.Properties.Encryption == nil {
				return fmt.Errorf("retrieving encryption for %s", id)
			}

			payload := dataFactory.Model
			encryption := payload.Properties.Encryption

			if metadata.ResourceData.HasChange("customer_managed_key_id") {
				keyVaultKey, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(customerManagedKey.CustomerManagedKeyID)
				if err != nil {
					return fmt.Errorf("could not parse Key Vault Key ID: %+v", err)
				}
				encryption.VaultBaseURL = keyVaultKey.KeyVaultBaseUrl
				encryption.KeyName = keyVaultKey.Name
				encryption.KeyVersion = pointer.To(keyVaultKey.Version)
			}

			if metadata.ResourceData.HasChange("user_assigned_identity_id") {
				encryption.Identity = &factories.CMKIdentityDefinition{
					UserAssignedIdentity: pointer.To(customerManagedKey.UserAssignedIdentityID),
				}
			}

			payload.Properties.Encryption = encryption
			if _, err = client.CreateOrUpdate(ctx, *id, *payload, factories.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r DataFactoryCustomerManagedKeyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			log.Printf(`[INFO] The Customer Managed Key cannot be removed from the Data Factory once added. To remove the Customer Managed Key delete and recreate the parent Data Factory.`)
			return nil
		},
	}
}
