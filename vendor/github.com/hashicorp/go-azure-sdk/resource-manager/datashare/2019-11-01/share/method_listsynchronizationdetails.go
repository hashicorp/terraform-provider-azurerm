package share

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListSynchronizationDetailsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SynchronizationDetails
}

type ListSynchronizationDetailsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SynchronizationDetails
}

type ListSynchronizationDetailsOperationOptions struct {
	Filter  *string
	Orderby *string
}

func DefaultListSynchronizationDetailsOperationOptions() ListSynchronizationDetailsOperationOptions {
	return ListSynchronizationDetailsOperationOptions{}
}

func (o ListSynchronizationDetailsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListSynchronizationDetailsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListSynchronizationDetailsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Orderby != nil {
		out.Append("$orderby", fmt.Sprintf("%v", *o.Orderby))
	}
	return &out
}

type ListSynchronizationDetailsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListSynchronizationDetailsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListSynchronizationDetails ...
func (c ShareClient) ListSynchronizationDetails(ctx context.Context, id ShareId, input ShareSynchronization, options ListSynchronizationDetailsOperationOptions) (result ListSynchronizationDetailsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Pager:         &ListSynchronizationDetailsCustomPager{},
		Path:          fmt.Sprintf("%s/listSynchronizationDetails", id.ID()),
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
		Values *[]SynchronizationDetails `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListSynchronizationDetailsComplete retrieves all the results into a single object
func (c ShareClient) ListSynchronizationDetailsComplete(ctx context.Context, id ShareId, input ShareSynchronization, options ListSynchronizationDetailsOperationOptions) (ListSynchronizationDetailsCompleteResult, error) {
	return c.ListSynchronizationDetailsCompleteMatchingPredicate(ctx, id, input, options, SynchronizationDetailsOperationPredicate{})
}

// ListSynchronizationDetailsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ShareClient) ListSynchronizationDetailsCompleteMatchingPredicate(ctx context.Context, id ShareId, input ShareSynchronization, options ListSynchronizationDetailsOperationOptions, predicate SynchronizationDetailsOperationPredicate) (result ListSynchronizationDetailsCompleteResult, err error) {
	items := make([]SynchronizationDetails, 0)

	resp, err := c.ListSynchronizationDetails(ctx, id, input, options)
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

	result = ListSynchronizationDetailsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
