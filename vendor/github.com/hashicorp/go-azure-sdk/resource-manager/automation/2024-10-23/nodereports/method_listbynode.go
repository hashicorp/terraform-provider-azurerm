package nodereports

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByNodeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DscNodeReport
}

type ListByNodeCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DscNodeReport
}

type ListByNodeOperationOptions struct {
	Filter *string
}

func DefaultListByNodeOperationOptions() ListByNodeOperationOptions {
	return ListByNodeOperationOptions{}
}

func (o ListByNodeOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByNodeOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByNodeOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type ListByNodeCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByNodeCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByNode ...
func (c NodeReportsClient) ListByNode(ctx context.Context, id NodeId, options ListByNodeOperationOptions) (result ListByNodeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByNodeCustomPager{},
		Path:          fmt.Sprintf("%s/reports", id.ID()),
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
		Values *[]DscNodeReport `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByNodeComplete retrieves all the results into a single object
func (c NodeReportsClient) ListByNodeComplete(ctx context.Context, id NodeId, options ListByNodeOperationOptions) (ListByNodeCompleteResult, error) {
	return c.ListByNodeCompleteMatchingPredicate(ctx, id, options, DscNodeReportOperationPredicate{})
}

// ListByNodeCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NodeReportsClient) ListByNodeCompleteMatchingPredicate(ctx context.Context, id NodeId, options ListByNodeOperationOptions, predicate DscNodeReportOperationPredicate) (result ListByNodeCompleteResult, err error) {
	items := make([]DscNodeReport, 0)

	resp, err := c.ListByNode(ctx, id, options)
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

	result = ListByNodeCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
