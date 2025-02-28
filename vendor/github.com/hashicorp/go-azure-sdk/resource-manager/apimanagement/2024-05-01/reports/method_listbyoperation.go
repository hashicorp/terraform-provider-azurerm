package reports

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByOperationOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ReportRecordContract
}

type ListByOperationCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ReportRecordContract
}

type ListByOperationOperationOptions struct {
	Filter  *string
	Orderby *string
	Skip    *int64
	Top     *int64
}

func DefaultListByOperationOperationOptions() ListByOperationOperationOptions {
	return ListByOperationOperationOptions{}
}

func (o ListByOperationOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByOperationOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByOperationOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Orderby != nil {
		out.Append("$orderby", fmt.Sprintf("%v", *o.Orderby))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListByOperationCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByOperationCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByOperation ...
func (c ReportsClient) ListByOperation(ctx context.Context, id ServiceId, options ListByOperationOperationOptions) (result ListByOperationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByOperationCustomPager{},
		Path:          fmt.Sprintf("%s/reports/byOperation", id.ID()),
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
		Values *[]ReportRecordContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByOperationComplete retrieves all the results into a single object
func (c ReportsClient) ListByOperationComplete(ctx context.Context, id ServiceId, options ListByOperationOperationOptions) (ListByOperationCompleteResult, error) {
	return c.ListByOperationCompleteMatchingPredicate(ctx, id, options, ReportRecordContractOperationPredicate{})
}

// ListByOperationCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ReportsClient) ListByOperationCompleteMatchingPredicate(ctx context.Context, id ServiceId, options ListByOperationOperationOptions, predicate ReportRecordContractOperationPredicate) (result ListByOperationCompleteResult, err error) {
	items := make([]ReportRecordContract, 0)

	resp, err := c.ListByOperation(ctx, id, options)
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

	result = ListByOperationCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
