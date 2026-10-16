package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListInstanceProcessThreadsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ProcessThreadInfo
}

type ListInstanceProcessThreadsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ProcessThreadInfo
}

type ListInstanceProcessThreadsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListInstanceProcessThreadsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListInstanceProcessThreads ...
func (c WebAppsClient) ListInstanceProcessThreads(ctx context.Context, id InstanceProcessId) (result ListInstanceProcessThreadsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListInstanceProcessThreadsCustomPager{},
		Path:       fmt.Sprintf("%s/threads", id.ID()),
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
		Values *[]ProcessThreadInfo `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListInstanceProcessThreadsComplete retrieves all the results into a single object
func (c WebAppsClient) ListInstanceProcessThreadsComplete(ctx context.Context, id InstanceProcessId) (ListInstanceProcessThreadsCompleteResult, error) {
	return c.ListInstanceProcessThreadsCompleteMatchingPredicate(ctx, id, ProcessThreadInfoOperationPredicate{})
}

// ListInstanceProcessThreadsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListInstanceProcessThreadsCompleteMatchingPredicate(ctx context.Context, id InstanceProcessId, predicate ProcessThreadInfoOperationPredicate) (result ListInstanceProcessThreadsCompleteResult, err error) {
	items := make([]ProcessThreadInfo, 0)

	resp, err := c.ListInstanceProcessThreads(ctx, id)
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

	result = ListInstanceProcessThreadsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
