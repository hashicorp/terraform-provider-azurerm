package querypackqueries

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueriesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LogAnalyticsQueryPackQuery
}

type QueriesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []LogAnalyticsQueryPackQuery
}

type QueriesListOperationOptions struct {
	IncludeBody *bool
	Top         *int64
}

func DefaultQueriesListOperationOptions() QueriesListOperationOptions {
	return QueriesListOperationOptions{}
}

func (o QueriesListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o QueriesListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o QueriesListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.IncludeBody != nil {
		out.Append("includeBody", fmt.Sprintf("%v", *o.IncludeBody))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type QueriesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *QueriesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// QueriesList ...
func (c QueryPackQueriesClient) QueriesList(ctx context.Context, id QueryPackId, options QueriesListOperationOptions) (result QueriesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &QueriesListCustomPager{},
		Path:          fmt.Sprintf("%s/queries", id.ID()),
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
		Values *[]LogAnalyticsQueryPackQuery `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// QueriesListComplete retrieves all the results into a single object
func (c QueryPackQueriesClient) QueriesListComplete(ctx context.Context, id QueryPackId, options QueriesListOperationOptions) (QueriesListCompleteResult, error) {
	return c.QueriesListCompleteMatchingPredicate(ctx, id, options, LogAnalyticsQueryPackQueryOperationPredicate{})
}

// QueriesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c QueryPackQueriesClient) QueriesListCompleteMatchingPredicate(ctx context.Context, id QueryPackId, options QueriesListOperationOptions, predicate LogAnalyticsQueryPackQueryOperationPredicate) (result QueriesListCompleteResult, err error) {
	items := make([]LogAnalyticsQueryPackQuery, 0)

	resp, err := c.QueriesList(ctx, id, options)
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

	result = QueriesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
