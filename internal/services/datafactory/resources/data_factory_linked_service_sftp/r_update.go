// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_sftp

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryLinkedServiceSFTPUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LinkedServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if resp.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` is nil", id)
	}

	sftp, ok := resp.Properties.AsSftpServerLinkedService()
	if !ok {
		return fmt.Errorf("classifying Data Factory Linked Service SFTP %s: Expected: %q Received: %q", id, datafactory.TypeBasicLinkedServiceTypeSftp, pointer.From(resp.Type))
	}

	if d.HasChange("authentication_type") {
		sftp.AuthenticationType = datafactory.SftpAuthenticationType(d.Get("authentication_type").(string))
	}

	if d.HasChange("host") {
		sftp.Host = pointer.To(d.Get("host").(string))
	}

	if d.HasChange("port") {
		sftp.Port = pointer.To(d.Get("port").(string))
	}

	if d.HasChange("username") {
		sftp.UserName = pointer.To(d.Get("username").(string))
	}

	if d.HasChange("password") {
		sftp.Password = datafactory.SecureString{
			Value: pointer.To(d.Get("password").(string)),
			Type:  datafactory.TypeSecureString,
		}
	}

	if d.HasChange("key_vault_password") {
		sftp.Password = helpers.ExpandAzureKeyVaultSecretReference(d.Get("key_vault_password").([]any))
	}

	if d.HasChange("key_vault_private_key_content_base64") {
		sftp.PrivateKeyContent = helpers.ExpandAzureKeyVaultSecretReference(d.Get("key_vault_private_key_content_base64").([]any))
	}

	if d.HasChange("key_vault_private_key_passphrase") {
		sftp.PassPhrase = helpers.ExpandAzureKeyVaultSecretReference(d.Get("key_vault_private_key_passphrase").([]any))
	}

	if d.HasChange("skip_host_key_validation") {
		sftp.SkipHostKeyValidation = pointer.To(d.Get("skip_host_key_validation").(bool))
	}

	if d.HasChange("host_key_fingerprint") {
		sftp.HostKeyFingerprint = pointer.To(d.Get("host_key_fingerprint").(string))
	}

	if d.HasChange("description") {
		sftp.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("parameters") {
		sftp.Parameters = helpers.ExpandLinkedServiceParameters(d.Get("parameters").(map[string]any))
	}

	if d.HasChange("integration_runtime_name") {
		sftp.ConnectVia = helpers.ExpandDataFactoryLinkedServiceIntegrationRuntime(d.Get("integration_runtime_name").(string))
	}

	if d.HasChange("additional_properties") {
		sftp.AdditionalProperties = d.Get("additional_properties").(map[string]any)
	}

	if d.HasChange("annotations") {
		sftp.Annotations = pointer.To(d.Get("annotations").([]any))
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: sftp,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, linkedService, ""); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceDataFactoryLinkedServiceSFTPRead(d, meta)
}
