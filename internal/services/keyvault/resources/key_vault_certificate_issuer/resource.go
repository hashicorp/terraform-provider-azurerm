// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate_issuer

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKeyVaultCertificateIssuerCreateOrUpdate,
		Update: resourceKeyVaultCertificateIssuerCreateOrUpdate,
		Read:   resourceKeyVaultCertificateIssuerRead,
		Delete: resourceKeyVaultCertificateIssuerDelete,
		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.IssuerID(id)
			return err
		}, helpers.NestedItemResourceImporter),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: keyVaultCertificateIssuerSchema(),
	}
}
