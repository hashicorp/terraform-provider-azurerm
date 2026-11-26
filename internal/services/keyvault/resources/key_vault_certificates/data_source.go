// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificates

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/certificates"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to retrieve a list of certificate names from an existing Key Vault.",
		Read:        dataSourceKeyVaultCertificatesRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}, "Key Vault"),

			"names": {
				Description: "List containing names of certificates that exist in this Key Vault.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"include_pending": {
				Description: "Whether to include certificates which are not completely provisioned.",
				Type:        pluginsdk.TypeBool,
				Optional:    true,
				Default:     true,
			},

			"certificates": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The ID of this certificate.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"name": {
							Description: "The name of certificate.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"enabled": {
							Description: "Whether this certificate is enabled.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},

						"tags": commonschema.TagsDataSource(),
					},
				},
			},
		},
	}
}

func dataSourceKeyVaultCertificatesRead(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	includePending := d.Get("include_pending").(bool)

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("fetching base vault url from id %q: %+v", *keyVaultId, err)
	}

	client := meta.(*clients.Client).KeyVault.DataPlaneKeyVaultClient.Certificates.Clone(*keyVaultBaseUri)
	opts := certificates.GetCertificatesOperationOptions{
		IncludePending: &includePending,
		Maxresults:     pointer.To(int64(25)),
	}
	certificateList, err := client.GetCertificatesComplete(ctx, opts)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *keyVaultId, err)
	}

	d.SetId(keyVaultId.ID())

	var names []string
	var certs []map[string]any
	for _, v := range certificateList.Items {
		nestedItem, err := keyvault.ParseNestedItemID(pointer.From(v.Id), keyvault.VersionTypeAny, keyvault.NestedItemTypeCertificate)
		if err != nil {
			return err
		}

		names = append(names, nestedItem.Name)
		certs = append(certs, expandCertificate(nestedItem.Name, v))
	}

	d.Set("names", names)
	d.Set("certificates", certs)
	d.Set("key_vault_id", keyVaultId.ID())

	return nil
}

func expandCertificate(name string, item certificates.CertificateItem) map[string]any {
	cert := map[string]any{
		"name": name,
		"id":   pointer.From(item.Id),
	}

	if item.Attributes != nil && item.Attributes.Enabled != nil {
		cert["enabled"] = *item.Attributes.Enabled
	}

	if item.Tags != nil {
		cert["tags"] = tags.Flatten(tags.FromTypedObject(pointer.From(item.Tags)))
	}

	return cert
}
