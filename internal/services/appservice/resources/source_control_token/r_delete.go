// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package source_control_token

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.ResourceProvidersClient

			id, err := resourceproviders.ParseSourceControlID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			sourceControlOAuth := resourceproviders.SourceControl{
				Properties: &resourceproviders.SourceControlProperties{
					Token:       pointer.To(""),
					TokenSecret: pointer.To(""),
				},
			}

			if _, err := client.UpdateSourceControl(ctx, *id, sourceControlOAuth); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
