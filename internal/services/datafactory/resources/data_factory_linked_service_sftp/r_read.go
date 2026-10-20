// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_sftp

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryLinkedServiceSFTPRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LinkedServiceID(d.Id())
	if err != nil {
		return err
	}

	dataFactoryId := factories.NewFactoryID(id.SubscriptionId, id.ResourceGroup, id.FactoryName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Data Factory SFTP %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	if resp.Properties != nil {
		sftp, ok := resp.Properties.AsSftpServerLinkedService()
		if !ok {
			return fmt.Errorf("classifying Data Factory Linked Service SFTP %s: Expected: %q Received: %q", id, datafactory.TypeBasicLinkedServiceTypeSftp, pointer.From(resp.Type))
		}

		d.Set("authentication_type", sftp.AuthenticationType)
		d.Set("username", sftp.UserName)
		d.Set("port", sftp.Port)
		d.Set("host", sftp.Host)

		d.Set("additional_properties", sftp.AdditionalProperties)
		d.Set("description", sftp.Description)

		if sftp.Password != nil {
			if v, ok := sftp.Password.AsAzureKeyVaultSecretReference(); ok {
				d.Set("key_vault_password", helpers.FlattenAzureKeyVaultSecretReference(v))
			}
		}

		if sftp.PrivateKeyContent != nil {
			if v, ok := sftp.PrivateKeyContent.AsAzureKeyVaultSecretReference(); ok {
				d.Set("key_vault_private_key_content_base64", helpers.FlattenAzureKeyVaultSecretReference(v))
			}
		}

		if sftp.PassPhrase != nil {
			if v, ok := sftp.PassPhrase.AsAzureKeyVaultSecretReference(); ok {
				d.Set("key_vault_private_key_passphrase", helpers.FlattenAzureKeyVaultSecretReference(v))
			}
		}

		if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(sftp.Annotations)); err != nil {
			return fmt.Errorf("setting `annotations`: %+v", err)
		}

		if err := d.Set("parameters", helpers.FlattenLinkedServiceParameters(sftp.Parameters)); err != nil {
			return fmt.Errorf("setting `parameters`: %+v", err)
		}

		if connectVia := sftp.ConnectVia; connectVia != nil {
			if connectVia.ReferenceName != nil {
				d.Set("integration_runtime_name", connectVia.ReferenceName)
			}
		}

		if props := sftp.SftpServerLinkedServiceTypeProperties; props != nil {
			if skipHostKeyValidation := props.SkipHostKeyValidation; skipHostKeyValidation != nil {
				d.Set("skip_host_key_validation", skipHostKeyValidation.(bool))
			}

			if hostKeyFingerprint := props.HostKeyFingerprint; hostKeyFingerprint != nil {
				d.Set("host_key_fingerprint", hostKeyFingerprint)
			}
		}
	}

	return nil
}
