// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_iothub_data_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/dataconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoIotHubDataConnectionUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataconnections.ParseDataConnectionID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return err
	}

	model := existing.Model
	if model == nil {
		return fmt.Errorf("retrieving existing %s: model was empty", id)
	}

	iothubDataConnection, ok := model.(dataconnections.IotHubDataConnection)
	if !ok {
		return fmt.Errorf("the Data Connection %s is not aa IoTHub Data Connection", id)
	}

	props := iothubDataConnection.Properties
	if props == nil {
		return fmt.Errorf("retrieving existing %s: properties were empty", id)
	}

	props.ProvisioningState = nil // clear read-only field

	if d.HasChange("retrieval_start_date") {
		props.RetrievalStartDate = pointer.To(d.Get("retrieval_start_date").(string))
	}

	dataConnection := dataconnections.IotHubDataConnection{
		Location:   pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: props,
	}

	if err = client.CreateOrUpdateThenPoll(ctx, *id, dataConnection); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceKustoIotHubDataConnectionRead(d, meta)
}
