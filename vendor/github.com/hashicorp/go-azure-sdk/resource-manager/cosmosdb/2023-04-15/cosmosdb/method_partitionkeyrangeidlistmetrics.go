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

type PartitionKeyRangeIdListMetricsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *PartitionMetricListResult
}

type PartitionKeyRangeIdListMetricsOperationOptions struct {
	Filter *string
}

func DefaultPartitionKeyRangeIdListMetricsOperationOptions() PartitionKeyRangeIdListMetricsOperationOptions {
	return PartitionKeyRangeIdListMetricsOperationOptions{}
}

func (o PartitionKeyRangeIdListMetricsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o PartitionKeyRangeIdListMetricsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o PartitionKeyRangeIdListMetricsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

// PartitionKeyRangeIdListMetrics ...
func (c CosmosDBClient) PartitionKeyRangeIdListMetrics(ctx context.Context, id PartitionKeyRangeIdId, options PartitionKeyRangeIdListMetricsOperationOptions) (result PartitionKeyRangeIdListMetricsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/metrics", id.ID()),
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
