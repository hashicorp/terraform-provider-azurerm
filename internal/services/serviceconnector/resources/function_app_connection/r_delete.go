// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package function_app_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2022-05-01/links"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	serviceconnectorClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/serviceconnector/client"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := serviceconnectorClient.From(metadata.Client).LinksClient
			id, err := links.ParseScopedLinkerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.LinkerDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}
