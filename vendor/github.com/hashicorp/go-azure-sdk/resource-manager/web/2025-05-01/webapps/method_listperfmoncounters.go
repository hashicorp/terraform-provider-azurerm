package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListPerfMonCountersOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PerfMonResponse
}

type ListPerfMonCountersCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PerfMonResponse
}

type ListPerfMonCountersOperationOptions struct {
	Filter *string
}

func DefaultListPerfMonCountersOperationOptions() ListPerfMonCountersOperationOptions {
	return ListPerfMonCountersOperationOptions{}
}

func (o ListPerfMonCountersOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListPerfMonCountersOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListPerfMonCountersOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type ListPerfMonCountersCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListPerfMonCountersCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListPerfMonCounters ...
func (c WebAppsClient) ListPerfMonCounters(ctx context.Context, id commonids.AppServiceId, options ListPerfMonCountersOperationOptions) (result ListPerfMonCountersOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListPerfMonCountersCustomPager{},
		Path:          fmt.Sprintf("%s/perfcounters", id.ID()),
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
		Values *[]PerfMonResponse `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListPerfMonCountersComplete retrieves all the results into a single object
func (c WebAppsClient) ListPerfMonCountersComplete(ctx context.Context, id commonids.AppServiceId, options ListPerfMonCountersOperationOptions) (ListPerfMonCountersCompleteResult, error) {
	return c.ListPerfMonCountersCompleteMatchingPredicate(ctx, id, options, PerfMonResponseOperationPredicate{})
}

// ListPerfMonCountersCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListPerfMonCountersCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceId, options ListPerfMonCountersOperationOptions, predicate PerfMonResponseOperationPredicate) (result ListPerfMonCountersCompleteResult, err error) {
	items := make([]PerfMonResponse, 0)

	resp, err := c.ListPerfMonCounters(ctx, id, options)
	if err != nil {
		result.LatestHttpResponse = resp.HttpResponse
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

	result = ListPerfMonCountersCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
