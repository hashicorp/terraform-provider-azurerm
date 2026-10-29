// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_kusto

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryLinkedServiceKustoRead(d *pluginsdk.ResourceData, meta any) error {
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

	linkedService, ok := resp.Properties.AsAzureDataExplorerLinkedService()
	if !ok {
		return fmt.Errorf("classifying %s: Expected: %q", *id, datafactory.TypeBasicLinkedServiceTypeAzureDataExplorer)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", factories.NewFactoryID(subscriptionId, id.ResourceGroup, id.FactoryName).ID())
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

	if prop := linkedService.AzureDataExplorerLinkedServiceTypeProperties; prop != nil {
		d.Set("kusto_endpoint", prop.Endpoint)
		d.Set("kusto_database_name", prop.Database)
		d.Set("tenant", prop.Tenant)
		d.Set("service_principal_id", prop.ServicePrincipalID)

		useManagedIdentity := true
		if prop.ServicePrincipalID != nil {
			useManagedIdentity = false
		}
		d.Set("use_managed_identity", useManagedIdentity)
	}

	return nil
}
