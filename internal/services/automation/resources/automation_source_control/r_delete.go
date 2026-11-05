// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_source_control

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/sourcecontrol"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := sourcecontrol.ParseSourceControlID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			client := meta.Client.Automation.SourceControl
			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", *id, err)
			}
			return nil
		},
	}
}
