// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func importStreamAnalyticsReferenceInput(expectType string) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := inputs.ParseInputID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).StreamAnalytics.InputsClient
		resp, err := client.Get(ctx, *id)
		if err != nil {
			return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				input, ok := props.(inputs.InputProperties) // nolint: gosimple
				if !ok {
					return nil, fmt.Errorf("failed to convert to Input")
				}

				reference, ok := input.(inputs.ReferenceInputProperties)
				if !ok {
					return nil, fmt.Errorf("failed to convert to Reference Input")
				}

				var actualType string

				if _, ok := reference.Datasource.(inputs.BlobReferenceInputDataSource); ok {
					actualType = "Microsoft.Storage/Blob"
				}
				if _, ok := reference.Datasource.(inputs.AzureSqlReferenceInputDataSource); ok {
					actualType = "Microsoft.Sql/Server/Database"
				}

				if actualType != expectType {
					return nil, fmt.Errorf("stream analytics reference input has mismatched type, expected: %q, got %q", expectType, actualType)
				}
			}
		}
		return []*pluginsdk.ResourceData{d}, nil
	}
}
