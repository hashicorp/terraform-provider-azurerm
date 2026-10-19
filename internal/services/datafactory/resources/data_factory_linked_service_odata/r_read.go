// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_odata

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

func resourceArmDataFactoryLinkedServiceODataRead(d *pluginsdk.ResourceData, meta any) error {
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

		return fmt.Errorf("retrieving Data Factory OData %s: %+v", *id, err)
	}

	d.Set("name", resp.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	odata, ok := resp.Properties.AsODataLinkedService()
	if !ok {
		return fmt.Errorf("classifying Data Factory OData %s: Expected: %q Received: %q", *id, datafactory.TypeBasicLinkedServiceTypeOData, *resp.Type)
	}

	props := odata.ODataLinkedServiceTypeProperties
	d.Set("url", props.URL)
	if props.AuthenticationType == datafactory.ODataAuthenticationTypeBasic {
		if err := d.Set("basic_authentication", []any{map[string]any{
			"username": props.UserName,
			// `password` isn't returned from the api so we'll set it to `*****` here to be able to check for diffs during plan
			"password": "*****",
		}}); err != nil {
			return fmt.Errorf("setting `basic_authentication`: %+v", err)
		}
	}

	d.Set("additional_properties", odata.AdditionalProperties)
	d.Set("description", odata.Description)

	if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(odata.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if err := d.Set("parameters", helpers.FlattenLinkedServiceParameters(odata.Parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	if connectVia := odata.ConnectVia; connectVia != nil {
		if connectVia.ReferenceName != nil {
			d.Set("integration_runtime_name", connectVia.ReferenceName)
		}
	}

	return nil
}
