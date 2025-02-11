package monitoredresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorsListMonitoredResourcesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]MonitoredResource
}

type MonitorsListMonitoredResourcesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []MonitoredResource
}

type MonitorsListMonitoredResourcesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *MonitorsListMonitoredResourcesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// MonitorsListMonitoredResources ...
func (c MonitoredResourcesClient) MonitorsListMonitoredResources(ctx context.Context, id MonitorId) (result MonitorsListMonitoredResourcesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &MonitorsListMonitoredResourcesCustomPager{},
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

// MonitorsListMonitoredResourcesComplete retrieves all the results into a single object
func (c MonitoredResourcesClient) MonitorsListMonitoredResourcesComplete(ctx context.Context, id MonitorId) (MonitorsListMonitoredResourcesCompleteResult, error) {
	return c.MonitorsListMonitoredResourcesCompleteMatchingPredicate(ctx, id, MonitoredResourceOperationPredicate{})
}

// MonitorsListMonitoredResourcesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c MonitoredResourcesClient) MonitorsListMonitoredResourcesCompleteMatchingPredicate(ctx context.Context, id MonitorId, predicate MonitoredResourceOperationPredicate) (result MonitorsListMonitoredResourcesCompleteResult, err error) {
	items := make([]MonitoredResource, 0)

	resp, err := c.MonitorsListMonitoredResources(ctx, id)
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

	result = MonitorsListMonitoredResourcesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
