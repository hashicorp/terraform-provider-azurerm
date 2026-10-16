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

type RestorableSqlResourcesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RestorableSqlResourcesGetResult
}

type RestorableSqlResourcesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RestorableSqlResourcesGetResult
}

type RestorableSqlResourcesListOperationOptions struct {
	RestoreLocation       *string
	RestoreTimestampInUtc *string
}

func DefaultRestorableSqlResourcesListOperationOptions() RestorableSqlResourcesListOperationOptions {
	return RestorableSqlResourcesListOperationOptions{}
}

func (o RestorableSqlResourcesListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RestorableSqlResourcesListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o RestorableSqlResourcesListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.RestoreLocation != nil {
		out.Append("restoreLocation", fmt.Sprintf("%v", *o.RestoreLocation))
	}
	if o.RestoreTimestampInUtc != nil {
		out.Append("restoreTimestampInUtc", fmt.Sprintf("%v", *o.RestoreTimestampInUtc))
	}
	return &out
}

type RestorableSqlResourcesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *RestorableSqlResourcesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// RestorableSqlResourcesList ...
func (c OpenapisClient) RestorableSqlResourcesList(ctx context.Context, id RestorableDatabaseAccountId, options RestorableSqlResourcesListOperationOptions) (result RestorableSqlResourcesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &RestorableSqlResourcesListCustomPager{},
		Path:          fmt.Sprintf("%s/restorableSqlResources", id.ID()),
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
		Values *[]RestorableSqlResourcesGetResult `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RestorableSqlResourcesListComplete retrieves all the results into a single object
func (c OpenapisClient) RestorableSqlResourcesListComplete(ctx context.Context, id RestorableDatabaseAccountId, options RestorableSqlResourcesListOperationOptions) (RestorableSqlResourcesListCompleteResult, error) {
	return c.RestorableSqlResourcesListCompleteMatchingPredicate(ctx, id, options, RestorableSqlResourcesGetResultOperationPredicate{})
}

// RestorableSqlResourcesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) RestorableSqlResourcesListCompleteMatchingPredicate(ctx context.Context, id RestorableDatabaseAccountId, options RestorableSqlResourcesListOperationOptions, predicate RestorableSqlResourcesGetResultOperationPredicate) (result RestorableSqlResourcesListCompleteResult, err error) {
	items := make([]RestorableSqlResourcesGetResult, 0)

	resp, err := c.RestorableSqlResourcesList(ctx, id, options)
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

	result = RestorableSqlResourcesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
