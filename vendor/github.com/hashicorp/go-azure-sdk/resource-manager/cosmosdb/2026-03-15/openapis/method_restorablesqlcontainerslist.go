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

type RestorableSqlContainersListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RestorableSqlContainerGetResult
}

type RestorableSqlContainersListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RestorableSqlContainerGetResult
}

type RestorableSqlContainersListOperationOptions struct {
	EndTime                  *string
	RestorableSqlDatabaseRid *string
	StartTime                *string
}

func DefaultRestorableSqlContainersListOperationOptions() RestorableSqlContainersListOperationOptions {
	return RestorableSqlContainersListOperationOptions{}
}

func (o RestorableSqlContainersListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RestorableSqlContainersListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o RestorableSqlContainersListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.EndTime != nil {
		out.Append("endTime", fmt.Sprintf("%v", *o.EndTime))
	}
	if o.RestorableSqlDatabaseRid != nil {
		out.Append("restorableSqlDatabaseRid", fmt.Sprintf("%v", *o.RestorableSqlDatabaseRid))
	}
	if o.StartTime != nil {
		out.Append("startTime", fmt.Sprintf("%v", *o.StartTime))
	}
	return &out
}

type RestorableSqlContainersListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *RestorableSqlContainersListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// RestorableSqlContainersList ...
func (c OpenapisClient) RestorableSqlContainersList(ctx context.Context, id RestorableDatabaseAccountId, options RestorableSqlContainersListOperationOptions) (result RestorableSqlContainersListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &RestorableSqlContainersListCustomPager{},
		Path:          fmt.Sprintf("%s/restorableSqlContainers", id.ID()),
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
		Values *[]RestorableSqlContainerGetResult `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RestorableSqlContainersListComplete retrieves all the results into a single object
func (c OpenapisClient) RestorableSqlContainersListComplete(ctx context.Context, id RestorableDatabaseAccountId, options RestorableSqlContainersListOperationOptions) (RestorableSqlContainersListCompleteResult, error) {
	return c.RestorableSqlContainersListCompleteMatchingPredicate(ctx, id, options, RestorableSqlContainerGetResultOperationPredicate{})
}

// RestorableSqlContainersListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) RestorableSqlContainersListCompleteMatchingPredicate(ctx context.Context, id RestorableDatabaseAccountId, options RestorableSqlContainersListOperationOptions, predicate RestorableSqlContainerGetResultOperationPredicate) (result RestorableSqlContainersListCompleteResult, err error) {
	items := make([]RestorableSqlContainerGetResult, 0)

	resp, err := c.RestorableSqlContainersList(ctx, id, options)
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

	result = RestorableSqlContainersListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
