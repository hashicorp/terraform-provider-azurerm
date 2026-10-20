// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry_task_schedule_run_now

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			return nil
		},
	}
}
