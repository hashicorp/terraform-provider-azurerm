package newrelicmonitorresources

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

type MonitorsListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NewRelicMonitorResource
}

type MonitorsListBySubscriptionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NewRelicMonitorResource
}

type MonitorsListBySubscriptionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *MonitorsListBySubscriptionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// MonitorsListBySubscription ...
func (c NewRelicMonitorResourcesClient) MonitorsListBySubscription(ctx context.Context, id commonids.SubscriptionId) (result MonitorsListBySubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &MonitorsListBySubscriptionCustomPager{},
		Path:       fmt.Sprintf("%s/providers/NewRelic.Observability/monitors", id.ID()),
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
		Values *[]NewRelicMonitorResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// MonitorsListBySubscriptionComplete retrieves all the results into a single object
func (c NewRelicMonitorResourcesClient) MonitorsListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (MonitorsListBySubscriptionCompleteResult, error) {
	return c.MonitorsListBySubscriptionCompleteMatchingPredicate(ctx, id, NewRelicMonitorResourceOperationPredicate{})
}

// MonitorsListBySubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NewRelicMonitorResourcesClient) MonitorsListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate NewRelicMonitorResourceOperationPredicate) (result MonitorsListBySubscriptionCompleteResult, err error) {
	items := make([]NewRelicMonitorResource, 0)

	resp, err := c.MonitorsListBySubscription(ctx, id)
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

	result = MonitorsListBySubscriptionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
