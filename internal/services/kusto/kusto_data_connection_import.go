package kusto

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2022-07-07/dataconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func importDataConnection() pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := dataconnections.ParseDataConnectionID(d.Id())
		if err != nil {
			return []*pluginsdk.ResourceData{}, err
		}

		client := meta.(*clients.Client).Kusto.DataConnectionsClient
		_, err = client.Get(ctx, *id)
		if err != nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving Kusto Data Connection %q (Resource Group %q): %+v", id.DataConnectionName, id.ResourceGroupName, err)
		}

		return []*pluginsdk.ResourceData{d}, nil
	}
}
