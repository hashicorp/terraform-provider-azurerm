package cosmosdb

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CollectionPartitionRegionListMetricsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *PartitionMetricListResult
}

type CollectionPartitionRegionListMetricsOperationOptions struct {
	Filter *string
}

func DefaultCollectionPartitionRegionListMetricsOperationOptions() CollectionPartitionRegionListMetricsOperationOptions {
	return CollectionPartitionRegionListMetricsOperationOptions{}
}

func (o CollectionPartitionRegionListMetricsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o CollectionPartitionRegionListMetricsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o CollectionPartitionRegionListMetricsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

// CollectionPartitionRegionListMetrics ...
func (c CosmosDBClient) CollectionPartitionRegionListMetrics(ctx context.Context, id DatabaseCollectionId, options CollectionPartitionRegionListMetricsOperationOptions) (result CollectionPartitionRegionListMetricsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/partitions/metrics", id.ID()),
		OptionsObject: options,
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var model PartitionMetricListResult
	result.Model = &model

	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
