package openapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestorableGremlinGraphsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RestorableGremlinGraphGetResult
}

type RestorableGremlinGraphsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RestorableGremlinGraphGetResult
}

type RestorableGremlinGraphsListOperationOptions struct {
	EndTime                      *string
	RestorableGremlinDatabaseRid *string
	StartTime                    *string
}

func DefaultRestorableGremlinGraphsListOperationOptions() RestorableGremlinGraphsListOperationOptions {
	return RestorableGremlinGraphsListOperationOptions{}
}

func (o RestorableGremlinGraphsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RestorableGremlinGraphsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o RestorableGremlinGraphsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.EndTime != nil {
		out.Append("endTime", fmt.Sprintf("%v", *o.EndTime))
	}
	if o.RestorableGremlinDatabaseRid != nil {
		out.Append("restorableGremlinDatabaseRid", fmt.Sprintf("%v", *o.RestorableGremlinDatabaseRid))
	}
	if o.StartTime != nil {
		out.Append("startTime", fmt.Sprintf("%v", *o.StartTime))
	}
	return &out
}

type RestorableGremlinGraphsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *RestorableGremlinGraphsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// RestorableGremlinGraphsList ...
func (c OpenapisClient) RestorableGremlinGraphsList(ctx context.Context, id RestorableDatabaseAccountId, options RestorableGremlinGraphsListOperationOptions) (result RestorableGremlinGraphsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &RestorableGremlinGraphsListCustomPager{},
		Path:          fmt.Sprintf("%s/restorableGraphs", id.ID()),
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
		Values *[]RestorableGremlinGraphGetResult `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RestorableGremlinGraphsListComplete retrieves all the results into a single object
func (c OpenapisClient) RestorableGremlinGraphsListComplete(ctx context.Context, id RestorableDatabaseAccountId, options RestorableGremlinGraphsListOperationOptions) (RestorableGremlinGraphsListCompleteResult, error) {
	return c.RestorableGremlinGraphsListCompleteMatchingPredicate(ctx, id, options, RestorableGremlinGraphGetResultOperationPredicate{})
}

// RestorableGremlinGraphsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) RestorableGremlinGraphsListCompleteMatchingPredicate(ctx context.Context, id RestorableDatabaseAccountId, options RestorableGremlinGraphsListOperationOptions, predicate RestorableGremlinGraphGetResultOperationPredicate) (result RestorableGremlinGraphsListCompleteResult, err error) {
	items := make([]RestorableGremlinGraphGetResult, 0)

	resp, err := c.RestorableGremlinGraphsList(ctx, id, options)
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

	result = RestorableGremlinGraphsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
