// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package gallery_application

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplications"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.GalleryApplicationsClient
			id, err := galleryapplications.ParseApplicationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
