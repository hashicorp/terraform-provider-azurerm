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

type RestorableGremlinResourcesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RestorableGremlinResourcesGetResult
}

type RestorableGremlinResourcesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RestorableGremlinResourcesGetResult
}

type RestorableGremlinResourcesListOperationOptions struct {
	RestoreLocation       *string
	RestoreTimestampInUtc *string
}

func DefaultRestorableGremlinResourcesListOperationOptions() RestorableGremlinResourcesListOperationOptions {
	return RestorableGremlinResourcesListOperationOptions{}
}

func (o RestorableGremlinResourcesListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RestorableGremlinResourcesListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o RestorableGremlinResourcesListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.RestoreLocation != nil {
		out.Append("restoreLocation", fmt.Sprintf("%v", *o.RestoreLocation))
	}
	if o.RestoreTimestampInUtc != nil {
		out.Append("restoreTimestampInUtc", fmt.Sprintf("%v", *o.RestoreTimestampInUtc))
	}
	return &out
}

type RestorableGremlinResourcesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *RestorableGremlinResourcesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// RestorableGremlinResourcesList ...
func (c OpenapisClient) RestorableGremlinResourcesList(ctx context.Context, id RestorableDatabaseAccountId, options RestorableGremlinResourcesListOperationOptions) (result RestorableGremlinResourcesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &RestorableGremlinResourcesListCustomPager{},
		Path:          fmt.Sprintf("%s/restorableGremlinResources", id.ID()),
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
		Values *[]RestorableGremlinResourcesGetResult `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RestorableGremlinResourcesListComplete retrieves all the results into a single object
func (c OpenapisClient) RestorableGremlinResourcesListComplete(ctx context.Context, id RestorableDatabaseAccountId, options RestorableGremlinResourcesListOperationOptions) (RestorableGremlinResourcesListCompleteResult, error) {
	return c.RestorableGremlinResourcesListCompleteMatchingPredicate(ctx, id, options, RestorableGremlinResourcesGetResultOperationPredicate{})
}

// RestorableGremlinResourcesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) RestorableGremlinResourcesListCompleteMatchingPredicate(ctx context.Context, id RestorableDatabaseAccountId, options RestorableGremlinResourcesListOperationOptions, predicate RestorableGremlinResourcesGetResultOperationPredicate) (result RestorableGremlinResourcesListCompleteResult, err error) {
	items := make([]RestorableGremlinResourcesGetResult, 0)

	resp, err := c.RestorableGremlinResourcesList(ctx, id, options)
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

	result = RestorableGremlinResourcesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
