package afdendpoints

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListResourceUsageOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Usage
}

type ListResourceUsageCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Usage
}

type ListResourceUsageCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListResourceUsageCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListResourceUsage ...
func (c AFDEndpointsClient) ListResourceUsage(ctx context.Context, id AfdEndpointId) (result ListResourceUsageOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &ListResourceUsageCustomPager{},
		Path:       fmt.Sprintf("%s/usages", id.ID()),
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
		Values *[]Usage `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListResourceUsageComplete retrieves all the results into a single object
func (c AFDEndpointsClient) ListResourceUsageComplete(ctx context.Context, id AfdEndpointId) (ListResourceUsageCompleteResult, error) {
	return c.ListResourceUsageCompleteMatchingPredicate(ctx, id, UsageOperationPredicate{})
}

// ListResourceUsageCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AFDEndpointsClient) ListResourceUsageCompleteMatchingPredicate(ctx context.Context, id AfdEndpointId, predicate UsageOperationPredicate) (result ListResourceUsageCompleteResult, err error) {
	items := make([]Usage, 0)

	resp, err := c.ListResourceUsage(ctx, id)
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

	result = ListResourceUsageCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
