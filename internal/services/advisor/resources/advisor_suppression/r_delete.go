// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package advisor_suppression

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/advisor/2023-01-01/suppressions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	advisorClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/advisor/client"
)

func (Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := advisorClient.From(metadata.Client).SuppressionsClient

			id, err := suppressions.ParseScopedSuppressionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
