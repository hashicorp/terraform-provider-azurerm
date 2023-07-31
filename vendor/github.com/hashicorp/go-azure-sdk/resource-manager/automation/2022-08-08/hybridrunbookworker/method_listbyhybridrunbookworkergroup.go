package hybridrunbookworker

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByHybridRunbookWorkerGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]HybridRunbookWorker
}

type ListByHybridRunbookWorkerGroupCompleteResult struct {
	Items []HybridRunbookWorker
}

type ListByHybridRunbookWorkerGroupOperationOptions struct {
	Filter *string
}

func DefaultListByHybridRunbookWorkerGroupOperationOptions() ListByHybridRunbookWorkerGroupOperationOptions {
	return ListByHybridRunbookWorkerGroupOperationOptions{}
}

func (o ListByHybridRunbookWorkerGroupOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByHybridRunbookWorkerGroupOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByHybridRunbookWorkerGroupOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

// ListByHybridRunbookWorkerGroup ...
func (c HybridRunbookWorkerClient) ListByHybridRunbookWorkerGroup(ctx context.Context, id HybridRunbookWorkerGroupId, options ListByHybridRunbookWorkerGroupOperationOptions) (result ListByHybridRunbookWorkerGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/hybridRunbookWorkers", id.ID()),
		OptionsObject: options,
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.ExecutePaged(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var values struct {
		Values *[]HybridRunbookWorker `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByHybridRunbookWorkerGroupComplete retrieves all the results into a single object
func (c HybridRunbookWorkerClient) ListByHybridRunbookWorkerGroupComplete(ctx context.Context, id HybridRunbookWorkerGroupId, options ListByHybridRunbookWorkerGroupOperationOptions) (ListByHybridRunbookWorkerGroupCompleteResult, error) {
	return c.ListByHybridRunbookWorkerGroupCompleteMatchingPredicate(ctx, id, options, HybridRunbookWorkerOperationPredicate{})
}

// ListByHybridRunbookWorkerGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c HybridRunbookWorkerClient) ListByHybridRunbookWorkerGroupCompleteMatchingPredicate(ctx context.Context, id HybridRunbookWorkerGroupId, options ListByHybridRunbookWorkerGroupOperationOptions, predicate HybridRunbookWorkerOperationPredicate) (result ListByHybridRunbookWorkerGroupCompleteResult, err error) {
	items := make([]HybridRunbookWorker, 0)

	resp, err := c.ListByHybridRunbookWorkerGroup(ctx, id, options)
	if err != nil {
		err = fmt.Errorf("loading results: %+v", err)
		return
	}
	if resp.Model != nil {
		for _, v := range *resp.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	result = ListByHybridRunbookWorkerGroupCompleteResult{
		Items: items,
	}
	return
}
