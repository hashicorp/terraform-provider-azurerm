// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web_app_active_slot

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// Nothing to do here - there's no actual resource to delete
			// Note: deleting does not change the active slot, nor revert to any previous state.
			return nil
		},
	}
}
