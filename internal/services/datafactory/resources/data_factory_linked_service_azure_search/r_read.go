// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_search

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

func resourceDataFactoryLinkedServiceAzureSearchRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LinkedServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", factories.NewFactoryID(subscriptionId, id.ResourceGroup, id.FactoryName).ID())

	linkedService, ok := resp.Properties.AsAzureSearchLinkedService()
	if !ok {
		return fmt.Errorf("classifying %s: Expected: %q", id, datafactory.TypeBasicLinkedServiceTypeAzureSearch)
	}

	if prop := linkedService.AzureSearchLinkedServiceTypeProperties; prop != nil {
		url := ""
		if v, ok := prop.URL.(string); ok {
			url = v
		}
		d.Set("url", url)

		d.Set("encrypted_credential", pointer.From(prop.EncryptedCredential))
	}

	d.Set("additional_properties", linkedService.AdditionalProperties)
	d.Set("description", linkedService.Description)

	if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(linkedService.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if err := d.Set("parameters", helpers.FlattenLinkedServiceParameters(linkedService.Parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	integrationRuntimeName := ""
	if linkedService.ConnectVia != nil && linkedService.ConnectVia.ReferenceName != nil {
		integrationRuntimeName = *linkedService.ConnectVia.ReferenceName
	}
	d.Set("integration_runtime_name", integrationRuntimeName)

	return nil
}
