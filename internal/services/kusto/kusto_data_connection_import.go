package kusto

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2022-02-01/kusto"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func importDataConnection(kind kusto.KindBasicDataConnection) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := parse.DataConnectionID(d.Id())
		if err != nil {
			return []*pluginsdk.ResourceData{}, err
		}

		client := meta.(*clients.Client).Kusto.DataConnectionsClient
		dataConnection, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.Name)
		if err != nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving Kusto Data Connection %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		if _, ok := dataConnection.Value.AsEventHubDataConnection(); ok && kind != kusto.KindBasicDataConnectionKindEventHub {
			return nil, fmt.Errorf(`kusto data connection "kind" mismatch, expected "%s", got "%s"`, kind, kusto.KindBasicDataConnectionKindEventHub)
		}

		if _, ok := dataConnection.Value.AsIotHubDataConnection(); ok && kind != kusto.KindBasicDataConnectionKindIotHub {
			return nil, fmt.Errorf(`kusto data connection "kind" mismatch, expected "%s", got "%s"`, kind, kusto.KindBasicDataConnectionKindIotHub)
		}

		if _, ok := dataConnection.Value.AsEventGridDataConnection(); ok && kind != kusto.KindBasicDataConnectionKindEventGrid {
			return nil, fmt.Errorf(`kusto data connection "kind" mismatch, expected "%s", got "%s"`, kind, kusto.KindBasicDataConnectionKindEventGrid)
		}

		return []*pluginsdk.ResourceData{d}, nil
	}
}
