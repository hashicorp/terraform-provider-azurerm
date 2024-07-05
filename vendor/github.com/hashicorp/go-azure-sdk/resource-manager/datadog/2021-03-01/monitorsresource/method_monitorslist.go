package monitorsresource

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

type MonitorsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DatadogMonitorResource
}

type MonitorsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DatadogMonitorResource
}

type MonitorsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *MonitorsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// MonitorsList ...
func (c MonitorsResourceClient) MonitorsList(ctx context.Context, id commonids.SubscriptionId) (result MonitorsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &MonitorsListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Datadog/monitors", id.ID()),
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
		Values *[]DatadogMonitorResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// MonitorsListComplete retrieves all the results into a single object
func (c MonitorsResourceClient) MonitorsListComplete(ctx context.Context, id commonids.SubscriptionId) (MonitorsListCompleteResult, error) {
	return c.MonitorsListCompleteMatchingPredicate(ctx, id, DatadogMonitorResourceOperationPredicate{})
}

// MonitorsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c MonitorsResourceClient) MonitorsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate DatadogMonitorResourceOperationPredicate) (result MonitorsListCompleteResult, err error) {
	items := make([]DatadogMonitorResource, 0)

	resp, err := c.MonitorsList(ctx, id)
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

	result = MonitorsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
