package applicationgateways

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

type ListAvailableSslPredefinedPoliciesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ApplicationGatewaySslPredefinedPolicy
}

type ListAvailableSslPredefinedPoliciesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ApplicationGatewaySslPredefinedPolicy
}

type ListAvailableSslPredefinedPoliciesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAvailableSslPredefinedPoliciesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAvailableSslPredefinedPolicies ...
func (c ApplicationGatewaysClient) ListAvailableSslPredefinedPolicies(ctx context.Context, id commonids.SubscriptionId) (result ListAvailableSslPredefinedPoliciesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListAvailableSslPredefinedPoliciesCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Network/applicationGatewayAvailableSslOptions/default/predefinedPolicies", id.ID()),
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
		Values *[]ApplicationGatewaySslPredefinedPolicy `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListAvailableSslPredefinedPoliciesComplete retrieves all the results into a single object
func (c ApplicationGatewaysClient) ListAvailableSslPredefinedPoliciesComplete(ctx context.Context, id commonids.SubscriptionId) (ListAvailableSslPredefinedPoliciesCompleteResult, error) {
	return c.ListAvailableSslPredefinedPoliciesCompleteMatchingPredicate(ctx, id, ApplicationGatewaySslPredefinedPolicyOperationPredicate{})
}

// ListAvailableSslPredefinedPoliciesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApplicationGatewaysClient) ListAvailableSslPredefinedPoliciesCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate ApplicationGatewaySslPredefinedPolicyOperationPredicate) (result ListAvailableSslPredefinedPoliciesCompleteResult, err error) {
	items := make([]ApplicationGatewaySslPredefinedPolicy, 0)

	resp, err := c.ListAvailableSslPredefinedPolicies(ctx, id)
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

	result = ListAvailableSslPredefinedPoliciesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
