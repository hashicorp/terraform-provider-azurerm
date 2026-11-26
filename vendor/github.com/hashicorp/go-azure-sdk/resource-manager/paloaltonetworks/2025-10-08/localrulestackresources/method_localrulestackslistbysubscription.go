package localrulestackresources

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

type LocalRulestacksListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LocalRulestackResource
}

type LocalRulestacksListBySubscriptionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []LocalRulestackResource
}

type LocalRulestacksListBySubscriptionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *LocalRulestacksListBySubscriptionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// LocalRulestacksListBySubscription ...
func (c LocalRulestackResourcesClient) LocalRulestacksListBySubscription(ctx context.Context, id commonids.SubscriptionId) (result LocalRulestacksListBySubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &LocalRulestacksListBySubscriptionCustomPager{},
		Path:       fmt.Sprintf("%s/providers/PaloAltoNetworks.Cloudngfw/localRulestacks", id.ID()),
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
		Values *[]LocalRulestackResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LocalRulestacksListBySubscriptionComplete retrieves all the results into a single object
func (c LocalRulestackResourcesClient) LocalRulestacksListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (LocalRulestacksListBySubscriptionCompleteResult, error) {
	return c.LocalRulestacksListBySubscriptionCompleteMatchingPredicate(ctx, id, LocalRulestackResourceOperationPredicate{})
}

// LocalRulestacksListBySubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LocalRulestackResourcesClient) LocalRulestacksListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate LocalRulestackResourceOperationPredicate) (result LocalRulestacksListBySubscriptionCompleteResult, err error) {
	items := make([]LocalRulestackResource, 0)

	resp, err := c.LocalRulestacksListBySubscription(ctx, id)
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

	result = LocalRulestacksListBySubscriptionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
