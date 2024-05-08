package hdinsights

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterJobsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ClusterJob
}

type ClusterJobsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ClusterJob
}

type ClusterJobsListOperationOptions struct {
	Filter *string
}

func DefaultClusterJobsListOperationOptions() ClusterJobsListOperationOptions {
	return ClusterJobsListOperationOptions{}
}

func (o ClusterJobsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ClusterJobsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ClusterJobsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

// ClusterJobsList ...
func (c HdinsightsClient) ClusterJobsList(ctx context.Context, id ClusterId, options ClusterJobsListOperationOptions) (result ClusterJobsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/jobs", id.ID()),
		OptionsObject: options,
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
		Values *[]ClusterJob `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ClusterJobsListComplete retrieves all the results into a single object
func (c HdinsightsClient) ClusterJobsListComplete(ctx context.Context, id ClusterId, options ClusterJobsListOperationOptions) (ClusterJobsListCompleteResult, error) {
	return c.ClusterJobsListCompleteMatchingPredicate(ctx, id, options, ClusterJobOperationPredicate{})
}

// ClusterJobsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c HdinsightsClient) ClusterJobsListCompleteMatchingPredicate(ctx context.Context, id ClusterId, options ClusterJobsListOperationOptions, predicate ClusterJobOperationPredicate) (result ClusterJobsListCompleteResult, err error) {
	items := make([]ClusterJob, 0)

	resp, err := c.ClusterJobsList(ctx, id, options)
	if err != nil {
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

	result = ClusterJobsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
