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

type ListProcessThreadsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ProcessThreadInfo
}

type ListProcessThreadsCompleteResult struct {
	Items []ProcessThreadInfo
}

// ListProcessThreads ...
func (c WebAppsClient) ListProcessThreads(ctx context.Context, id ProcessId) (result ListProcessThreadsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
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

// ListProcessThreadsComplete retrieves all the results into a single object
func (c WebAppsClient) ListProcessThreadsComplete(ctx context.Context, id ProcessId) (ListProcessThreadsCompleteResult, error) {
	return c.ListProcessThreadsCompleteMatchingPredicate(ctx, id, ProcessThreadInfoOperationPredicate{})
}

// ListProcessThreadsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListProcessThreadsCompleteMatchingPredicate(ctx context.Context, id ProcessId, predicate ProcessThreadInfoOperationPredicate) (result ListProcessThreadsCompleteResult, err error) {
	items := make([]ProcessThreadInfo, 0)

	resp, err := c.ListProcessThreads(ctx, id)
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

	result = ListProcessThreadsCompleteResult{
		Items: items,
	}
	return
}
