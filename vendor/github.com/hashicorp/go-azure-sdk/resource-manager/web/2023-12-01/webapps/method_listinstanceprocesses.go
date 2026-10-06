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

type ListInstanceProcessesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ProcessInfo
}

type ListInstanceProcessesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ProcessInfo
}

type ListInstanceProcessesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListInstanceProcessesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListInstanceProcesses ...
func (c WebAppsClient) ListInstanceProcesses(ctx context.Context, id InstanceId) (result ListInstanceProcessesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListInstanceProcessesCustomPager{},
		Path:       fmt.Sprintf("%s/processes", id.ID()),
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
		Values *[]ProcessInfo `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListInstanceProcessesComplete retrieves all the results into a single object
func (c WebAppsClient) ListInstanceProcessesComplete(ctx context.Context, id InstanceId) (ListInstanceProcessesCompleteResult, error) {
	return c.ListInstanceProcessesCompleteMatchingPredicate(ctx, id, ProcessInfoOperationPredicate{})
}

// ListInstanceProcessesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListInstanceProcessesCompleteMatchingPredicate(ctx context.Context, id InstanceId, predicate ProcessInfoOperationPredicate) (result ListInstanceProcessesCompleteResult, err error) {
	items := make([]ProcessInfo, 0)

	resp, err := c.ListInstanceProcesses(ctx, id)
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

	result = ListInstanceProcessesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
