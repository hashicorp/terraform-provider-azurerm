// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package service_plan

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.ServicePlanClient
			id, err := commonids.ParseAppServicePlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			servicePlan, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(servicePlan.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			return r.flatten(metadata, id, servicePlan.Model)
		},
	}
}

type OSType string

const (
	OSTypeLinux            OSType = "Linux"
	OSTypeWindows          OSType = "Windows"
	OSTypeWindowsContainer OSType = "WindowsContainer"
)
