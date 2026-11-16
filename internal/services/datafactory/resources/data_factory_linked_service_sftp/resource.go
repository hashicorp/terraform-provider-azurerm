// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_sftp

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryLinkedServiceSFTPCreate,
		Read:   resourceDataFactoryLinkedServiceSFTPRead,
		Update: resourceDataFactoryLinkedServiceSFTPUpdate,
		Delete: resourceDataFactoryLinkedServiceSFTPDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.LinkedServiceID(id)
			return err
		}, helpers.ImportDataFactoryLinkedService(datafactory.TypeBasicLinkedServiceTypeSftp)),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: dataFactoryLinkedServiceSftpSchema(),
		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
				authCombinations := map[string]string{
					"private_key_content_base64":           string(datafactory.SftpAuthenticationTypeSSHPublicKey),
					"key_vault_private_key_content_base64": string(datafactory.SftpAuthenticationTypeSSHPublicKey),
					"key_vault_private_key_passphrase":     string(datafactory.SftpAuthenticationTypeSSHPublicKey),
					"private_key_path":                     string(datafactory.SftpAuthenticationTypeSSHPublicKey),
					"password":                             string(datafactory.SftpAuthenticationTypeBasic),
					"key_vault_password":                   string(datafactory.SftpAuthenticationTypeBasic),
				}

				for keyType, authType := range authCombinations {
					if _, ok := d.GetOk(keyType); ok {
						if !slices.Contains([]string{authType, string(datafactory.SftpAuthenticationTypeMultiFactor)}, d.Get("authentication_type").(string)) {
							return fmt.Errorf("when `%s` is defined, `authentication_type` must be `%s` or `%s`", keyType, authType, datafactory.SftpAuthenticationTypeMultiFactor)
						}
					}
				}

				return nil
			}),
		),
	}
}
