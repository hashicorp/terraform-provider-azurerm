package resourceproviders

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

type ListPremierAddOnOffersOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PremierAddOnOffer
}

type ListPremierAddOnOffersCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PremierAddOnOffer
}

type ListPremierAddOnOffersCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListPremierAddOnOffersCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListPremierAddOnOffers ...
func (c ResourceProvidersClient) ListPremierAddOnOffers(ctx context.Context, id commonids.SubscriptionId) (result ListPremierAddOnOffersOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListPremierAddOnOffersCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Web/premieraddonoffers", id.ID()),
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
		Values *[]PremierAddOnOffer `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListPremierAddOnOffersComplete retrieves all the results into a single object
func (c ResourceProvidersClient) ListPremierAddOnOffersComplete(ctx context.Context, id commonids.SubscriptionId) (ListPremierAddOnOffersCompleteResult, error) {
	return c.ListPremierAddOnOffersCompleteMatchingPredicate(ctx, id, PremierAddOnOfferOperationPredicate{})
}

// ListPremierAddOnOffersCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ResourceProvidersClient) ListPremierAddOnOffersCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate PremierAddOnOfferOperationPredicate) (result ListPremierAddOnOffersCompleteResult, err error) {
	items := make([]PremierAddOnOffer, 0)

	resp, err := c.ListPremierAddOnOffers(ctx, id)
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

	result = ListPremierAddOnOffersCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
