package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListProcessesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ProcessInfo
}

type ListProcessesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ProcessInfo
}

// ListProcesses ...
func (c WebAppsClient) ListProcesses(ctx context.Context, id commonids.AppServiceId) (result ListProcessesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
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

// ListProcessesComplete retrieves all the results into a single object
func (c WebAppsClient) ListProcessesComplete(ctx context.Context, id commonids.AppServiceId) (ListProcessesCompleteResult, error) {
	return c.ListProcessesCompleteMatchingPredicate(ctx, id, ProcessInfoOperationPredicate{})
}

// ListProcessesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListProcessesCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceId, predicate ProcessInfoOperationPredicate) (result ListProcessesCompleteResult, err error) {
	items := make([]ProcessInfo, 0)

	resp, err := c.ListProcesses(ctx, id)
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

	result = ListProcessesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
