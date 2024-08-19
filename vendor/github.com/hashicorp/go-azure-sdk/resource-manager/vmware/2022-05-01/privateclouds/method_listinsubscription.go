package privateclouds

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

type ListInSubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PrivateCloud
}

type ListInSubscriptionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PrivateCloud
}

type ListInSubscriptionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListInSubscriptionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListInSubscription ...
func (c PrivateCloudsClient) ListInSubscription(ctx context.Context, id commonids.SubscriptionId) (result ListInSubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListInSubscriptionCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.AVS/privateClouds", id.ID()),
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
		Values *[]PrivateCloud `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListInSubscriptionComplete retrieves all the results into a single object
func (c PrivateCloudsClient) ListInSubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (ListInSubscriptionCompleteResult, error) {
	return c.ListInSubscriptionCompleteMatchingPredicate(ctx, id, PrivateCloudOperationPredicate{})
}

// ListInSubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PrivateCloudsClient) ListInSubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate PrivateCloudOperationPredicate) (result ListInSubscriptionCompleteResult, err error) {
	items := make([]PrivateCloud, 0)

	resp, err := c.ListInSubscription(ctx, id)
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

	result = ListInSubscriptionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
