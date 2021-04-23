package kusto

import (
	"context"
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2020-09-18/kusto"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
)

func importDataConnection(kind kusto.KindBasicDataConnection) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := parse.DataConnectionID(d.Id())
		if err != nil {
			return []*schema.ResourceData{}, err
		}

		client := meta.(*clients.Client).Kusto.DataConnectionsClient
		dataConnection, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.Name)
		if err != nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving Kusto Data Connection %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		if _, ok := dataConnection.Value.AsEventHubDataConnection(); ok && kind != kusto.KindEventHub {
			return nil, fmt.Errorf(`kusto data connection "kind" mismatch, expected "%s", got "%s"`, kind, kusto.KindEventHub)
		}

		if _, ok := dataConnection.Value.AsIotHubDataConnection(); ok && kind != kusto.KindIotHub {
			return nil, fmt.Errorf(`kusto data connection "kind" mismatch, expected "%s", got "%s"`, kind, kusto.KindIotHub)
		}

		if _, ok := dataConnection.Value.AsEventGridDataConnection(); ok && kind != kusto.KindEventGrid {
			return nil, fmt.Errorf(`kusto data connection "kind" mismatch, expected "%s", got "%s"`, kind, kusto.KindEventGrid)
		}

		return []*pluginsdk.ResourceData{d}, nil
	}
}
