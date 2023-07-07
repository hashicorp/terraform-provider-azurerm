// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/datasources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func importLogAnalyticsDataSource(kind datasources.DataSourceKind) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := datasources.ParseDataSourceID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).LogAnalytics.DataSourcesClient
		resp, err := client.Get(ctx, *id)
		if err != nil {
			return nil, fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.Model != nil && resp.Model.Kind != kind {
			return nil, fmt.Errorf(`log analytics Data Source "kind" mismatch, expected "%s", got "%s"`, kind, resp.Model.Kind)
		}
		return []*pluginsdk.ResourceData{d}, nil
	}
}
