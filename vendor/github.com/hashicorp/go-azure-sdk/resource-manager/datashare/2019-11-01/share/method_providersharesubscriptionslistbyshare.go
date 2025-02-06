package share

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProviderShareSubscriptionsListByShareOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ProviderShareSubscription
}

type ProviderShareSubscriptionsListByShareCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ProviderShareSubscription
}

type ProviderShareSubscriptionsListByShareCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ProviderShareSubscriptionsListByShareCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ProviderShareSubscriptionsListByShare ...
func (c ShareClient) ProviderShareSubscriptionsListByShare(ctx context.Context, id ShareId) (result ProviderShareSubscriptionsListByShareOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ProviderShareSubscriptionsListByShareCustomPager{},
		Path:       fmt.Sprintf("%s/providerShareSubscriptions", id.ID()),
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
		Values *[]ProviderShareSubscription `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ProviderShareSubscriptionsListByShareComplete retrieves all the results into a single object
func (c ShareClient) ProviderShareSubscriptionsListByShareComplete(ctx context.Context, id ShareId) (ProviderShareSubscriptionsListByShareCompleteResult, error) {
	return c.ProviderShareSubscriptionsListByShareCompleteMatchingPredicate(ctx, id, ProviderShareSubscriptionOperationPredicate{})
}

// ProviderShareSubscriptionsListByShareCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ShareClient) ProviderShareSubscriptionsListByShareCompleteMatchingPredicate(ctx context.Context, id ShareId, predicate ProviderShareSubscriptionOperationPredicate) (result ProviderShareSubscriptionsListByShareCompleteResult, err error) {
	items := make([]ProviderShareSubscription, 0)

	resp, err := c.ProviderShareSubscriptionsListByShare(ctx, id)
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

	result = ProviderShareSubscriptionsListByShareCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
