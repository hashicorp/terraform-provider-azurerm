// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_workbook_template

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-11-20/workbooktemplatesapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppInsights.WorkbookTemplateClient

			id, err := workbooktemplatesapis.ParseWorkbookTemplateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.WorkbookTemplatesDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
