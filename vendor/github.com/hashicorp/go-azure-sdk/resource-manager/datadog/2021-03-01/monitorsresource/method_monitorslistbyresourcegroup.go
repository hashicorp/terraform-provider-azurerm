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

type MonitorsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DatadogMonitorResource
}

type MonitorsListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DatadogMonitorResource
}

type MonitorsListByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *MonitorsListByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// MonitorsListByResourceGroup ...
func (c MonitorsResourceClient) MonitorsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result MonitorsListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &MonitorsListByResourceGroupCustomPager{},
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

// MonitorsListByResourceGroupComplete retrieves all the results into a single object
func (c MonitorsResourceClient) MonitorsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (MonitorsListByResourceGroupCompleteResult, error) {
	return c.MonitorsListByResourceGroupCompleteMatchingPredicate(ctx, id, DatadogMonitorResourceOperationPredicate{})
}

// MonitorsListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c MonitorsResourceClient) MonitorsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate DatadogMonitorResourceOperationPredicate) (result MonitorsListByResourceGroupCompleteResult, err error) {
	items := make([]DatadogMonitorResource, 0)

	resp, err := c.MonitorsListByResourceGroup(ctx, id)
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

	result = MonitorsListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
