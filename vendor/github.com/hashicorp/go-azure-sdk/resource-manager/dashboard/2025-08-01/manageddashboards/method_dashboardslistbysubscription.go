package manageddashboards

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

type DashboardsListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ManagedDashboard
}

type DashboardsListBySubscriptionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ManagedDashboard
}

type DashboardsListBySubscriptionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DashboardsListBySubscriptionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DashboardsListBySubscription ...
func (c ManagedDashboardsClient) DashboardsListBySubscription(ctx context.Context, id commonids.SubscriptionId) (result DashboardsListBySubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &DashboardsListBySubscriptionCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Dashboard/dashboards", id.ID()),
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
		Values *[]ManagedDashboard `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DashboardsListBySubscriptionComplete retrieves all the results into a single object
func (c ManagedDashboardsClient) DashboardsListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (DashboardsListBySubscriptionCompleteResult, error) {
	return c.DashboardsListBySubscriptionCompleteMatchingPredicate(ctx, id, ManagedDashboardOperationPredicate{})
}

// DashboardsListBySubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagedDashboardsClient) DashboardsListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate ManagedDashboardOperationPredicate) (result DashboardsListBySubscriptionCompleteResult, err error) {
	items := make([]ManagedDashboard, 0)

	resp, err := c.DashboardsListBySubscription(ctx, id)
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

	result = DashboardsListBySubscriptionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
