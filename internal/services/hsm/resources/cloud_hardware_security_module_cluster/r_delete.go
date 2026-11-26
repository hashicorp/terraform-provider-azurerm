// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cloud_hardware_security_module_cluster

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2025-03-31/cloudhsmclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HSM.CloudHsmClustersClient

			id, err := cloudhsmclusters.ParseCloudHsmClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// the deletion of Cloud HSM may return 409 if the linked Private Endpoint is deleted recently
			deadline, ok := ctx.Deadline()
			if !ok {
				return errors.New("internal-error: context had no deadline")
			}

			if err = pluginsdk.Retry(time.Until(deadline), func() *pluginsdk.RetryError {
				resp, err := client.Delete(ctx, *id)
				if err != nil {
					if response.WasConflict(resp.HttpResponse) {
						return pluginsdk.RetryableError(err)
					}
					return pluginsdk.NonRetryableError(err)
				}
				if err := resp.Poller.PollUntilDone(ctx); err != nil {
					return pluginsdk.NonRetryableError(err)
				}
				return nil
			}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
