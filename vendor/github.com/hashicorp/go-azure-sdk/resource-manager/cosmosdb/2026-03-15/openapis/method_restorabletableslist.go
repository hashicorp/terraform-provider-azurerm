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

type RestorableTablesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RestorableTableGetResult
}

type RestorableTablesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RestorableTableGetResult
}

type RestorableTablesListOperationOptions struct {
	EndTime   *string
	StartTime *string
}

func DefaultRestorableTablesListOperationOptions() RestorableTablesListOperationOptions {
	return RestorableTablesListOperationOptions{}
}

func (o RestorableTablesListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RestorableTablesListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o RestorableTablesListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.EndTime != nil {
		out.Append("endTime", fmt.Sprintf("%v", *o.EndTime))
	}
	if o.StartTime != nil {
		out.Append("startTime", fmt.Sprintf("%v", *o.StartTime))
	}
	return &out
}

type RestorableTablesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *RestorableTablesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// RestorableTablesList ...
func (c OpenapisClient) RestorableTablesList(ctx context.Context, id RestorableDatabaseAccountId, options RestorableTablesListOperationOptions) (result RestorableTablesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &RestorableTablesListCustomPager{},
		Path:          fmt.Sprintf("%s/restorableTables", id.ID()),
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
		Values *[]RestorableTableGetResult `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RestorableTablesListComplete retrieves all the results into a single object
func (c OpenapisClient) RestorableTablesListComplete(ctx context.Context, id RestorableDatabaseAccountId, options RestorableTablesListOperationOptions) (RestorableTablesListCompleteResult, error) {
	return c.RestorableTablesListCompleteMatchingPredicate(ctx, id, options, RestorableTableGetResultOperationPredicate{})
}

// RestorableTablesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) RestorableTablesListCompleteMatchingPredicate(ctx context.Context, id RestorableDatabaseAccountId, options RestorableTablesListOperationOptions, predicate RestorableTableGetResultOperationPredicate) (result RestorableTablesListCompleteResult, err error) {
	items := make([]RestorableTableGetResult, 0)

	resp, err := c.RestorableTablesList(ctx, id, options)
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

	result = RestorableTablesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
