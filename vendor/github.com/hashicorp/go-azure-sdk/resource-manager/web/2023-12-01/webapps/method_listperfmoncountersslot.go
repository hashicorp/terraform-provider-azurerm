package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListPerfMonCountersSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PerfMonResponse
}

type ListPerfMonCountersSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PerfMonResponse
}

type ListPerfMonCountersSlotOperationOptions struct {
	Filter *string
}

func DefaultListPerfMonCountersSlotOperationOptions() ListPerfMonCountersSlotOperationOptions {
	return ListPerfMonCountersSlotOperationOptions{}
}

func (o ListPerfMonCountersSlotOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListPerfMonCountersSlotOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListPerfMonCountersSlotOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type ListPerfMonCountersSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListPerfMonCountersSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListPerfMonCountersSlot ...
func (c WebAppsClient) ListPerfMonCountersSlot(ctx context.Context, id SlotId, options ListPerfMonCountersSlotOperationOptions) (result ListPerfMonCountersSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListPerfMonCountersSlotCustomPager{},
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

// ListPerfMonCountersSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListPerfMonCountersSlotComplete(ctx context.Context, id SlotId, options ListPerfMonCountersSlotOperationOptions) (ListPerfMonCountersSlotCompleteResult, error) {
	return c.ListPerfMonCountersSlotCompleteMatchingPredicate(ctx, id, options, PerfMonResponseOperationPredicate{})
}

// ListPerfMonCountersSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListPerfMonCountersSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, options ListPerfMonCountersSlotOperationOptions, predicate PerfMonResponseOperationPredicate) (result ListPerfMonCountersSlotCompleteResult, err error) {
	items := make([]PerfMonResponse, 0)

	resp, err := c.ListPerfMonCountersSlot(ctx, id, options)
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

	result = ListPerfMonCountersSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
