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

type ListByApiOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ReportRecordContract
}

type ListByApiCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ReportRecordContract
}

type ListByApiOperationOptions struct {
	Filter  *string
	Orderby *string
	Skip    *int64
	Top     *int64
}

func DefaultListByApiOperationOptions() ListByApiOperationOptions {
	return ListByApiOperationOptions{}
}

func (o ListByApiOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByApiOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByApiOperationOptions) ToQuery() *client.QueryParams {
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

type ListByApiCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByApiCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByApi ...
func (c ReportsClient) ListByApi(ctx context.Context, id ServiceId, options ListByApiOperationOptions) (result ListByApiOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByApiCustomPager{},
		Path:          fmt.Sprintf("%s/reports/byApi", id.ID()),
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

// ListByApiComplete retrieves all the results into a single object
func (c ReportsClient) ListByApiComplete(ctx context.Context, id ServiceId, options ListByApiOperationOptions) (ListByApiCompleteResult, error) {
	return c.ListByApiCompleteMatchingPredicate(ctx, id, options, ReportRecordContractOperationPredicate{})
}

// ListByApiCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ReportsClient) ListByApiCompleteMatchingPredicate(ctx context.Context, id ServiceId, options ListByApiOperationOptions, predicate ReportRecordContractOperationPredicate) (result ListByApiCompleteResult, err error) {
	items := make([]ReportRecordContract, 0)

	resp, err := c.ListByApi(ctx, id, options)
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

	result = ListByApiCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
