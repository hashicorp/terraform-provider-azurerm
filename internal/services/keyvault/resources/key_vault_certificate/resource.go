// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate

import (
	"context"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DeleteAndPurgeNestedItem interface {
	DeleteNestedItem(ctx context.Context) (autorest.Response, error)
	NestedItemHasBeenDeleted(ctx context.Context) (autorest.Response, error)
	PurgeNestedItem(ctx context.Context) (autorest.Response, error)
	NestedItemHasBeenPurged(ctx context.Context) (autorest.Response, error)
}

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		// TODO: support Updating additional properties once we have more information about what can be updated
		Create: resourceKeyVaultCertificateCreate,
		Read:   resourceKeyVaultCertificateRead,
		Delete: resourceKeyVaultCertificateDelete,
		Update: resourceKeyVaultCertificateUpdate,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := keyvault.ParseNestedItemID(id, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeCertificate)
			return err
		}, helpers.NestedItemResourceImporter),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			// TODO: Change this back to 5min, once https://github.com/hashicorp/terraform-provider-azurerm/issues/11059 is addressed.
			Read:   pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: keyVaultCertificateSchema(),
	}
}
