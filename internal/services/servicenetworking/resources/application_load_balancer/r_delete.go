// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/trafficcontrollerinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	servicenetworkingClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicenetworking/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (t Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := servicenetworkingClient.From(metadata.Client).TrafficControllerInterface

			id, err := trafficcontrollerinterface.ParseTrafficControllerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// a workaround for that some child resources may still exist for seconds before it fully deleted.
			// tracked o https://github.com/Azure/azure-rest-api-specs/issues/26000
			// it will cause the error "Can not delete resource before nested resources are deleted."
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("could not retrieve context deadline for %s", id.ID())
			}
			stateConf := &pluginsdk.StateChangeConf{
				Delay:   5 * time.Minute,
				Pending: []string{"409"},
				Target:  []string{"200", "202"},
				Refresh: func() (result any, state string, err error) {
					resp, err := client.Delete(ctx, *id)
					if err != nil {
						if resp.HttpResponse.StatusCode == http.StatusConflict {
							return nil, "409", nil
						}
						return nil, "", err
					}
					return resp, "200", nil
				},
				MinTimeout: 15 * time.Second,
				Timeout:    time.Until(deadline),
			}

			if future, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for deleting of %s: %+v", id, err)
			} else {
				poller := future.(trafficcontrollerinterface.DeleteOperationResponse).Poller
				if err := poller.PollUntilDone(ctx); err != nil {
					return fmt.Errorf("deleting %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}
