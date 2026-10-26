// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_file_storage

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

func resourceDataFactoryLinkedServiceAzureFileStorageCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	subscriptionId := meta.(*clients.Client).DataFactory.LinkedServiceClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewLinkedServiceID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing Data Factory Azure File Storage Anonymous %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_data_factory_linked_service_azure_file_storage", id.ID())
			}
		}
	}

	fileStorageProperties := &datafactory.AzureFileStorageLinkedServiceTypeProperties{
		ConnectionString: &datafactory.SecureString{
			Value: pointer.To(d.Get("connection_string").(string)),
			Type:  datafactory.TypeSecureString,
		},
		FileShare: d.Get("file_share").(string),
	}

	if host := d.Get("host").(string); host != "" {
		fileStorageProperties.Host = host
	}

	if userId := d.Get("user_id").(string); userId != "" {
		fileStorageProperties.UserID = userId
	}

	password := d.Get("password").(string)
	if password != "" {
		fileStorageProperties.Password = &datafactory.SecureString{
			Value: pointer.To(d.Get("password").(string)),
			Type:  datafactory.TypeSecureString,
		}
	}

	fileStorageLinkedService := &datafactory.AzureFileStorageLinkedService{
		Description: pointer.To(d.Get("description").(string)),
		AzureFileStorageLinkedServiceTypeProperties: fileStorageProperties,
		Type: datafactory.TypeBasicLinkedServiceTypeAzureFileStorage,
	}

	if v, ok := d.GetOk("parameters"); ok {
		fileStorageLinkedService.Parameters = helpers.ExpandLinkedServiceParameters(v.(map[string]any))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		fileStorageLinkedService.ConnectVia = helpers.ExpandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("key_vault_password"); ok {
		password := v.([]any)
		fileStorageProperties.Password = helpers.ExpandAzureKeyVaultSecretReference(password)
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		fileStorageLinkedService.AdditionalProperties = v.(map[string]any)
	}

	if v, ok := d.GetOk("annotations"); ok {
		fileStorageLinkedService.Annotations = pointer.To(v.([]any))
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: fileStorageLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, linkedService, ""); err != nil {
		return fmt.Errorf("creating/updating Data Factory Azure File Storage %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryLinkedServiceAzureFileStorageRead(d, meta)
}
