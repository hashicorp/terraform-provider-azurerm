// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_workbook_template

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-11-20/workbooktemplatesapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppInsights.WorkbookTemplateClient

			id, err := workbooktemplatesapis.ParseWorkbookTemplateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkbookTemplatesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r Resource) Identity() resourceids.ResourceId {
	return &workbooktemplatesapis.WorkbookTemplateId{}
}
