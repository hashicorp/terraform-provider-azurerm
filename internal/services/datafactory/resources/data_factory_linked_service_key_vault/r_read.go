// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_key_vault

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryLinkedServiceKeyVaultRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
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

		return fmt.Errorf("retrieving Data Factory Key Vault %s: %+v", *id, err)
	}

	d.Set("name", resp.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	keyVault, ok := resp.Properties.AsAzureKeyVaultLinkedService()
	if !ok {
		return fmt.Errorf("classifying Data Factory Key Vault %s: Expected: %q Received: %q", *id, datafactory.TypeBasicLinkedServiceTypeAzureKeyVault, *resp.Type)
	}

	d.Set("additional_properties", keyVault.AdditionalProperties)
	d.Set("description", keyVault.Description)

	if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(keyVault.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if err := d.Set("parameters", helpers.FlattenLinkedServiceParameters(keyVault.Parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	if connectVia := keyVault.ConnectVia; connectVia != nil {
		if connectVia.ReferenceName != nil {
			d.Set("integration_runtime_name", connectVia.ReferenceName)
		}
	}

	baseUrl := ""
	if properties := keyVault.AzureKeyVaultLinkedServiceTypeProperties; properties != nil {
		if properties.BaseURL != nil {
			val, ok := properties.BaseURL.(string)
			if ok {
				baseUrl = val
			} else {
				log.Printf("[DEBUG] Skipping base url string %q since it's not a string", val)
			}
		}
	}

	var keyVaultId *string
	if baseUrl != "" {
		subscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)
		keyVaultId, err = keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionId, baseUrl)
		if err != nil {
			return err
		}
	}
	d.Set("key_vault_id", keyVaultId)

	return nil
}
