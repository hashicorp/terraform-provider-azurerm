// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

func dataSourceKeyVaultCertificates() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceKeyVaultCertificatesRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),

			"names": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"include_pending": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"certificates": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"tags": tags.SchemaDataSource(),
					},
				},
			},
		},
	}
}

func dataSourceKeyVaultCertificatesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
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

	certificateList, err := client.GetCertificatesComplete(ctx, *keyVaultBaseUri, utils.Int32(25), &includePending)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *keyVaultId, err)
	}

	d.SetId(keyVaultId.ID())

	var names []string
	var certs []map[string]interface{}
	if certificateList.Response().Value != nil {
		for certificateList.NotDone() {
			for _, v := range *certificateList.Response().Value {
				nestedItem, err := parse.ParseOptionallyVersionedNestedItemID(*v.ID)
				if err != nil {
					return err
				}
				names = append(names, nestedItem.Name)
				certs = append(certs, expandCertificate(nestedItem.Name, v))
				err = certificateList.NextWithContext(ctx)
				if err != nil {
					return fmt.Errorf("retrieving next page of Certificates from %s: %+v", *keyVaultId, err)
				}
			}
		}
	}

	d.Set("names", names)
	d.Set("certificates", certs)
	d.Set("key_vault_id", keyVaultId.ID())

	return nil
}

func expandCertificate(name string, item keyvault.CertificateItem) map[string]interface{} {
	var cert = map[string]interface{}{
		"name": name,
		"id":   *item.ID,
	}

	if item.Attributes != nil && item.Attributes.Enabled != nil {
		cert["enabled"] = *item.Attributes.Enabled
	}

	if item.Tags != nil {
		cert["tags"] = tags.Flatten(item.Tags)
	}

	return cert
}
