package dataset

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByShareOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DataSet
}

type ListByShareCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DataSet
}

type ListByShareOperationOptions struct {
	Filter  *string
	Orderby *string
}

func DefaultListByShareOperationOptions() ListByShareOperationOptions {
	return ListByShareOperationOptions{}
}

func (o ListByShareOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByShareOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByShareOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Orderby != nil {
		out.Append("$orderby", fmt.Sprintf("%v", *o.Orderby))
	}
	return &out
}

// ListByShare ...
func (c DataSetClient) ListByShare(ctx context.Context, id ShareId, options ListByShareOperationOptions) (result ListByShareOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/dataSets", id.ID()),
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
		Values *[]DataSet `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByShareComplete retrieves all the results into a single object
func (c DataSetClient) ListByShareComplete(ctx context.Context, id ShareId, options ListByShareOperationOptions) (ListByShareCompleteResult, error) {
	return c.ListByShareCompleteMatchingPredicate(ctx, id, options, DataSetOperationPredicate{})
}

// ListByShareCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DataSetClient) ListByShareCompleteMatchingPredicate(ctx context.Context, id ShareId, options ListByShareOperationOptions, predicate DataSetOperationPredicate) (result ListByShareCompleteResult, err error) {
	items := make([]DataSet, 0)

	resp, err := c.ListByShare(ctx, id, options)
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

	result = ListByShareCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
