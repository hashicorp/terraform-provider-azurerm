// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_sftp

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryLinkedServiceSFTPCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	subscriptionId := meta.(*clients.Client).DataFactory.LinkedServiceClient.SubscriptionID
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewLinkedServiceID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of existing Data Factory SFTP %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.Response.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_linked_service_sftp", id.ID())
		}
	}

	sftpProperties := &datafactory.SftpServerLinkedServiceTypeProperties{
		Host:               d.Get("host").(string),
		Port:               d.Get("port").(int),
		AuthenticationType: datafactory.SftpAuthenticationType(d.Get("authentication_type").(string)),
		UserName:           d.Get("username").(string),
	}

	if v, ok := d.GetOk("password"); ok {
		sftpProperties.Password = &datafactory.SecureString{
			Value: pointer.To(v.(string)),
			Type:  datafactory.TypeSecureString,
		}
	}

	if v, ok := d.GetOk("key_vault_password"); ok {
		sftpProperties.Password = helpers.ExpandAzureKeyVaultSecretReference(v.([]any))
	}

	if v, ok := d.GetOk("private_key_content_base64"); ok {
		sftpProperties.PrivateKeyContent = &datafactory.SecureString{
			Value: pointer.To(v.(string)),
			Type:  datafactory.TypeSecureString,
		}
	}

	if v, ok := d.GetOk("key_vault_private_key_content_base64"); ok {
		sftpProperties.PrivateKeyContent = helpers.ExpandAzureKeyVaultSecretReference(v.([]any))
	}

	if v, ok := d.GetOk("private_key_passphrase"); ok {
		sftpProperties.PassPhrase = &datafactory.SecureString{
			Value: pointer.To(v.(string)),
			Type:  datafactory.TypeSecureString,
		}
	}

	if v, ok := d.GetOk("key_vault_private_key_passphrase"); ok {
		sftpProperties.PassPhrase = helpers.ExpandAzureKeyVaultSecretReference(v.([]any))
	}

	if v, ok := d.GetOk("private_key_path"); ok {
		sftpProperties.PrivateKeyPath = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("skip_host_key_validation"); ok {
		sftpProperties.SkipHostKeyValidation = pointer.To(v.(bool))
	}

	if v, ok := d.GetOk("host_key_fingerprint"); ok {
		sftpProperties.HostKeyFingerprint = pointer.To(v.(string))
	}

	sftpLinkedService := &datafactory.SftpServerLinkedService{
		SftpServerLinkedServiceTypeProperties: sftpProperties,
		Type:                                  datafactory.TypeBasicLinkedServiceTypeSftp,
	}

	if v, ok := d.GetOk("description"); ok {
		sftpLinkedService.Description = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("parameters"); ok {
		sftpLinkedService.Parameters = helpers.ExpandLinkedServiceParameters(v.(map[string]any))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		sftpLinkedService.ConnectVia = helpers.ExpandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		sftpLinkedService.AdditionalProperties = v.(map[string]any)
	}

	if v, ok := d.GetOk("annotations"); ok {
		sftpLinkedService.Annotations = pointer.To(v.([]any))
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: sftpLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, linkedService, ""); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryLinkedServiceSFTPRead(d, meta)
}
