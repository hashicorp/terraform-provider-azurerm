// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_dataset_azure_sql_table

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
)

func (Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			d := metadata.ResourceData
			client := metadata.Client.DataFactory.DatasetClient

			id, err := parse.DataSetID(d.Id())
			if err != nil {
				return err
			}

			resp, err := client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name)
			if err != nil {
				if !response.WasNotFound(resp.Response) {
					return fmt.Errorf("deleting %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}
