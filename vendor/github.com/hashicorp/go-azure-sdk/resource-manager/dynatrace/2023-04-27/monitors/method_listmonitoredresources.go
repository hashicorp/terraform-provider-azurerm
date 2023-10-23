package monitors

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListMonitoredResourcesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]MonitoredResource
}

type ListMonitoredResourcesCompleteResult struct {
	Items []MonitoredResource
}

// ListMonitoredResources ...
func (c MonitorsClient) ListMonitoredResources(ctx context.Context, id MonitorId) (result ListMonitoredResourcesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/listMonitoredResources", id.ID()),
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
		Values *[]MonitoredResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListMonitoredResourcesComplete retrieves all the results into a single object
func (c MonitorsClient) ListMonitoredResourcesComplete(ctx context.Context, id MonitorId) (ListMonitoredResourcesCompleteResult, error) {
	return c.ListMonitoredResourcesCompleteMatchingPredicate(ctx, id, MonitoredResourceOperationPredicate{})
}

// ListMonitoredResourcesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c MonitorsClient) ListMonitoredResourcesCompleteMatchingPredicate(ctx context.Context, id MonitorId, predicate MonitoredResourceOperationPredicate) (result ListMonitoredResourcesCompleteResult, err error) {
	items := make([]MonitoredResource, 0)

	resp, err := c.ListMonitoredResources(ctx, id)
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

	result = ListMonitoredResourcesCompleteResult{
		Items: items,
	}
	return
}
