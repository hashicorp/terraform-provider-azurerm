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

type RestorableTableResourcesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RestorableTableResourcesGetResult
}

type RestorableTableResourcesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RestorableTableResourcesGetResult
}

type RestorableTableResourcesListOperationOptions struct {
	RestoreLocation       *string
	RestoreTimestampInUtc *string
}

func DefaultRestorableTableResourcesListOperationOptions() RestorableTableResourcesListOperationOptions {
	return RestorableTableResourcesListOperationOptions{}
}

func (o RestorableTableResourcesListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RestorableTableResourcesListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o RestorableTableResourcesListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.RestoreLocation != nil {
		out.Append("restoreLocation", fmt.Sprintf("%v", *o.RestoreLocation))
	}
	if o.RestoreTimestampInUtc != nil {
		out.Append("restoreTimestampInUtc", fmt.Sprintf("%v", *o.RestoreTimestampInUtc))
	}
	return &out
}

type RestorableTableResourcesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *RestorableTableResourcesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// RestorableTableResourcesList ...
func (c OpenapisClient) RestorableTableResourcesList(ctx context.Context, id RestorableDatabaseAccountId, options RestorableTableResourcesListOperationOptions) (result RestorableTableResourcesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &RestorableTableResourcesListCustomPager{},
		Path:          fmt.Sprintf("%s/restorableTableResources", id.ID()),
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
		Values *[]RestorableTableResourcesGetResult `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RestorableTableResourcesListComplete retrieves all the results into a single object
func (c OpenapisClient) RestorableTableResourcesListComplete(ctx context.Context, id RestorableDatabaseAccountId, options RestorableTableResourcesListOperationOptions) (RestorableTableResourcesListCompleteResult, error) {
	return c.RestorableTableResourcesListCompleteMatchingPredicate(ctx, id, options, RestorableTableResourcesGetResultOperationPredicate{})
}

// RestorableTableResourcesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) RestorableTableResourcesListCompleteMatchingPredicate(ctx context.Context, id RestorableDatabaseAccountId, options RestorableTableResourcesListOperationOptions, predicate RestorableTableResourcesGetResultOperationPredicate) (result RestorableTableResourcesListCompleteResult, err error) {
	items := make([]RestorableTableResourcesGetResult, 0)

	resp, err := c.RestorableTableResourcesList(ctx, id, options)
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

	result = RestorableTableResourcesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
