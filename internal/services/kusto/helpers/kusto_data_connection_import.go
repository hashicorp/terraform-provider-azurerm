// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/dataconnections"
	kustoClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func ImportDataConnection(kind string) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta any) (data []*pluginsdk.ResourceData, err error) {
		id, err := dataconnections.ParseDataConnectionID(d.Id())
		if err != nil {
			return []*pluginsdk.ResourceData{}, err
		}

		client := kustoClient.From(meta).DataConnectionsClient
		dataConnection, err := client.Get(ctx, *id)
		if err != nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving Kusto Data Connection %q (Resource Group %q): %+v", id.DataConnectionName, id.ResourceGroupName, err)
		}

		if dataConnection.Model != nil {
			if dataCon, ok := dataConnection.Model.(dataconnections.EventHubDataConnection); ok && kind != "EventHub" {
				return nil, fmt.Errorf(`kusto data connection "kind" mismatch, expected "%T", got "%T"`, kind, dataCon)
			}

			if dataCon, ok := dataConnection.Model.(dataconnections.IotHubDataConnection); ok && kind != "IotHub" {
				return nil, fmt.Errorf(`kusto data connection "kind" mismatch, expected "%T", got "%T"`, kind, dataCon)
			}

			if dataCon, ok := dataConnection.Model.(dataconnections.EventGridDataConnection); ok && kind != "EventGrid" {
				return nil, fmt.Errorf(`kusto data connection "kind" mismatch, expected "%T", got "%T"`, kind, dataCon)
			}
		}

		return []*pluginsdk.ResourceData{d}, nil
	}
}
