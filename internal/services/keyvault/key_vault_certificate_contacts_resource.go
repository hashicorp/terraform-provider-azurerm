// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

type KeyVaultCertificateContactsResource struct{}

var (
	_ sdk.ResourceWithUpdate = KeyVaultCertificateContactsResource{}
)

type KeyVaultCertificateContactsResourceModel struct {
	KeyVaultId string    `tfschema:"key_vault_id"`
	Contact    []Contact `tfschema:"contact"`
}

type Contact struct {
	Email string `tfschema:"email"`
	Name  string `tfschema:"name"`
	Phone string `tfschema:"phone"`
}

func (r KeyVaultCertificateContactsResource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"key_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KeyVaultId{}),
	}

	if features.FourPointOhBeta() {
		schema["contact"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"phone": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		}
	} else {
		schema["contact"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"phone": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		}
	}

	return schema
}

func (r KeyVaultCertificateContactsResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r KeyVaultCertificateContactsResource) ResourceType() string {
	return "azurerm_key_vault_certificate_contacts"
}

func (r KeyVaultCertificateContactsResource) ModelObject() interface{} {
	return &KeyVaultCertificateContactsResourceModel{}
}

func (r KeyVaultCertificateContactsResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.CertificateContactsID
}

func (r KeyVaultCertificateContactsResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			vaultClient := metadata.Client.KeyVault
			client := metadata.Client.KeyVault.ManagementClient
			var state KeyVaultCertificateContactsResourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			keyVaultId, err := commonids.ParseKeyVaultID(state.KeyVaultId)
			if err != nil {
				return fmt.Errorf("parsing `key_vault_id`, %+v", err)
			}

			keyVaultBaseUri, err := vaultClient.BaseUriForKeyVault(ctx, *keyVaultId)
			if err != nil {
				return fmt.Errorf("looking up Base URI for Key Vault Certificate Contacts from %s: %+v", *keyVaultId, err)
			}

			id, err := parse.NewCertificateContactsID(*keyVaultBaseUri)
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			existing, err := client.GetCertificateContacts(ctx, *keyVaultBaseUri)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing Certificate Contacts (Key Vault %q): %s", *keyVaultBaseUri, err)
				}
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				if existing.ContactList != nil && len(*existing.ContactList) != 0 {
					return tf.ImportAsExistsError(r.ResourceType(), id.ID())
				}
			}

			contacts := keyvault.Contacts{
				ContactList: expandKeyVaultCertificateContactsContact(state.Contact),
			}

			if features.FourPointOhBeta() {
				if len(*contacts.ContactList) == 0 {
					if _, err := client.DeleteCertificateContacts(ctx, id.KeyVaultBaseUrl); err != nil {
						return fmt.Errorf("removing Key Vault Certificate Contacts %s: %+v", id, err)
					}
				} else {
					if _, err := client.SetCertificateContacts(ctx, *keyVaultBaseUri, contacts); err != nil {
						return fmt.Errorf("creating Key Vault Certificate Contacts %s: %+v", id, err)
					}
				}
			} else {
				if _, err := client.SetCertificateContacts(ctx, *keyVaultBaseUri, contacts); err != nil {
					return fmt.Errorf("creating Key Vault Certificate Contacts %s: %+v", id, err)
				}
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r KeyVaultCertificateContactsResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			vaultClient := metadata.Client.KeyVault
			client := metadata.Client.KeyVault.ManagementClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id, err := parse.CertificateContactsID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
			keyVaultIdRaw, err := vaultClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, id.KeyVaultBaseUrl)
			if err != nil {
				return fmt.Errorf("retrieving resource ID of the Key Vault at URL %s: %+v", id.KeyVaultBaseUrl, err)
			}
			if keyVaultIdRaw == nil {
				metadata.Logger.Infof("Unable to determine the Resource ID for the Key Vault at URL %s - removing from state!", id.KeyVaultBaseUrl)
				return metadata.MarkAsGone(id)
			}
			keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
			if err != nil {
				return fmt.Errorf("parsing Key Vault ID: %+v", err)
			}

			existing, err := client.GetCertificateContacts(ctx, id.KeyVaultBaseUrl)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					metadata.Logger.Infof("No Certificate Contacts could be found at %s - removing from state!", id.KeyVaultBaseUrl)
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("checking for presence of existing Certificate Contacts (Key Vault %q): %s", id.KeyVaultBaseUrl, err)
			}

			state := KeyVaultCertificateContactsResourceModel{
				KeyVaultId: keyVaultId.ID(),
				Contact:    flattenKeyVaultCertificateContactsContact(existing.ContactList),
			}

			return metadata.Encode(&state)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r KeyVaultCertificateContactsResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.KeyVault.ManagementClient

			var state KeyVaultCertificateContactsResourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id, err := parse.CertificateContactsID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			existing, err := client.GetCertificateContacts(ctx, id.KeyVaultBaseUrl)
			if err != nil {
				return fmt.Errorf("checking for presence of existing Certificate Contacts (Key Vault %q): %s", id.KeyVaultBaseUrl, err)
			}

			if metadata.ResourceData.HasChange("contact") {
				existing.ContactList = expandKeyVaultCertificateContactsContact(state.Contact)
			}

			if features.FourPointOhBeta() {
				if len(*existing.ContactList) == 0 {
					if _, err := client.DeleteCertificateContacts(ctx, id.KeyVaultBaseUrl); err != nil {
						return fmt.Errorf("removing Key Vault Certificate Contacts %s: %+v", id, err)
					}
				} else {
					if _, err := client.SetCertificateContacts(ctx, id.KeyVaultBaseUrl, existing); err != nil {
						return fmt.Errorf("updating Key Vault Certificate Contacts %s: %+v", id, err)
					}
				}
			} else {
				if _, err := client.SetCertificateContacts(ctx, id.KeyVaultBaseUrl, existing); err != nil {
					return fmt.Errorf("updating Key Vault Certificate Contacts %s: %+v", id, err)
				}
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r KeyVaultCertificateContactsResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.KeyVault.ManagementClient

			id, err := parse.CertificateContactsID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			if _, err := client.DeleteCertificateContacts(ctx, id.KeyVaultBaseUrl); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func expandKeyVaultCertificateContactsContact(input []Contact) *[]keyvault.Contact {
	results := make([]keyvault.Contact, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		results = append(results, keyvault.Contact{
			EmailAddress: utils.String(item.Email),
			Name:         utils.String(item.Name),
			Phone:        utils.String(item.Phone),
		})
	}

	return &results
}

func flattenKeyVaultCertificateContactsContact(input *[]keyvault.Contact) []Contact {
	result := make([]Contact, 0)
	if input == nil {
		return result
	}

	for _, item := range *input {
		emailAddress := ""
		if item.EmailAddress != nil {
			emailAddress = *item.EmailAddress
		}

		name := ""
		if item.Name != nil {
			name = *item.Name
		}

		phone := ""
		if item.Phone != nil {
			phone = *item.Phone
		}

		result = append(result, Contact{
			Email: emailAddress,
			Name:  name,
			Phone: phone,
		})
	}

	return result
}
