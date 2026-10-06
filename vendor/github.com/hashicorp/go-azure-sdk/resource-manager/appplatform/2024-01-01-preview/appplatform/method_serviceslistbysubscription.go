package appplatform

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

type ServicesListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ServiceResource
}

type ServicesListBySubscriptionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ServiceResource
}

type ServicesListBySubscriptionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ServicesListBySubscriptionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ServicesListBySubscription ...
func (c AppPlatformClient) ServicesListBySubscription(ctx context.Context, id commonids.SubscriptionId) (result ServicesListBySubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ServicesListBySubscriptionCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.AppPlatform/spring", id.ID()),
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
		Values *[]ServiceResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ServicesListBySubscriptionComplete retrieves all the results into a single object
func (c AppPlatformClient) ServicesListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (ServicesListBySubscriptionCompleteResult, error) {
	return c.ServicesListBySubscriptionCompleteMatchingPredicate(ctx, id, ServiceResourceOperationPredicate{})
}

// ServicesListBySubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) ServicesListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate ServiceResourceOperationPredicate) (result ServicesListBySubscriptionCompleteResult, err error) {
	items := make([]ServiceResource, 0)

	resp, err := c.ServicesListBySubscription(ctx, id)
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

	result = ServicesListBySubscriptionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
