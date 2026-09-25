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

type RestorableMongodbResourcesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RestorableMongodbResourcesGetResult
}

type RestorableMongodbResourcesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RestorableMongodbResourcesGetResult
}

type RestorableMongodbResourcesListOperationOptions struct {
	RestoreLocation       *string
	RestoreTimestampInUtc *string
}

func DefaultRestorableMongodbResourcesListOperationOptions() RestorableMongodbResourcesListOperationOptions {
	return RestorableMongodbResourcesListOperationOptions{}
}

func (o RestorableMongodbResourcesListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RestorableMongodbResourcesListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o RestorableMongodbResourcesListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.RestoreLocation != nil {
		out.Append("restoreLocation", fmt.Sprintf("%v", *o.RestoreLocation))
	}
	if o.RestoreTimestampInUtc != nil {
		out.Append("restoreTimestampInUtc", fmt.Sprintf("%v", *o.RestoreTimestampInUtc))
	}
	return &out
}

type RestorableMongodbResourcesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *RestorableMongodbResourcesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// RestorableMongodbResourcesList ...
func (c OpenapisClient) RestorableMongodbResourcesList(ctx context.Context, id RestorableDatabaseAccountId, options RestorableMongodbResourcesListOperationOptions) (result RestorableMongodbResourcesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &RestorableMongodbResourcesListCustomPager{},
		Path:          fmt.Sprintf("%s/restorableMongodbResources", id.ID()),
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
		Values *[]RestorableMongodbResourcesGetResult `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RestorableMongodbResourcesListComplete retrieves all the results into a single object
func (c OpenapisClient) RestorableMongodbResourcesListComplete(ctx context.Context, id RestorableDatabaseAccountId, options RestorableMongodbResourcesListOperationOptions) (RestorableMongodbResourcesListCompleteResult, error) {
	return c.RestorableMongodbResourcesListCompleteMatchingPredicate(ctx, id, options, RestorableMongodbResourcesGetResultOperationPredicate{})
}

// RestorableMongodbResourcesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) RestorableMongodbResourcesListCompleteMatchingPredicate(ctx context.Context, id RestorableDatabaseAccountId, options RestorableMongodbResourcesListOperationOptions, predicate RestorableMongodbResourcesGetResultOperationPredicate) (result RestorableMongodbResourcesListCompleteResult, err error) {
	items := make([]RestorableMongodbResourcesGetResult, 0)

	resp, err := c.RestorableMongodbResourcesList(ctx, id, options)
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

	result = RestorableMongodbResourcesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
