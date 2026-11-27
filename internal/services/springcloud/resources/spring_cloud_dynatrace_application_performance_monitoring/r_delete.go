// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_dynatrace_application_performance_monitoring

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (s Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.AppPlatformClient

			id, err := appplatform.ParseApmID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.ApmsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
