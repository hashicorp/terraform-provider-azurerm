package firewallresources

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

type FirewallsListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]FirewallResource
}

type FirewallsListBySubscriptionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []FirewallResource
}

type FirewallsListBySubscriptionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *FirewallsListBySubscriptionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// FirewallsListBySubscription ...
func (c FirewallResourcesClient) FirewallsListBySubscription(ctx context.Context, id commonids.SubscriptionId) (result FirewallsListBySubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &FirewallsListBySubscriptionCustomPager{},
		Path:       fmt.Sprintf("%s/providers/PaloAltoNetworks.Cloudngfw/firewalls", id.ID()),
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
		Values *[]FirewallResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// FirewallsListBySubscriptionComplete retrieves all the results into a single object
func (c FirewallResourcesClient) FirewallsListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (FirewallsListBySubscriptionCompleteResult, error) {
	return c.FirewallsListBySubscriptionCompleteMatchingPredicate(ctx, id, FirewallResourceOperationPredicate{})
}

// FirewallsListBySubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c FirewallResourcesClient) FirewallsListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate FirewallResourceOperationPredicate) (result FirewallsListBySubscriptionCompleteResult, err error) {
	items := make([]FirewallResource, 0)

	resp, err := c.FirewallsListBySubscription(ctx, id)
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

	result = FirewallsListBySubscriptionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
