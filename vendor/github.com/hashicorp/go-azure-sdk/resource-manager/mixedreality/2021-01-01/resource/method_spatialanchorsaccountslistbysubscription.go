package resource

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

type SpatialAnchorsAccountsListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SpatialAnchorsAccount
}

type SpatialAnchorsAccountsListBySubscriptionCompleteResult struct {
	Items []SpatialAnchorsAccount
}

// SpatialAnchorsAccountsListBySubscription ...
func (c ResourceClient) SpatialAnchorsAccountsListBySubscription(ctx context.Context, id commonids.SubscriptionId) (result SpatialAnchorsAccountsListBySubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.MixedReality/spatialAnchorsAccounts", id.ID()),
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
		Values *[]SpatialAnchorsAccount `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SpatialAnchorsAccountsListBySubscriptionComplete retrieves all the results into a single object
func (c ResourceClient) SpatialAnchorsAccountsListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (SpatialAnchorsAccountsListBySubscriptionCompleteResult, error) {
	return c.SpatialAnchorsAccountsListBySubscriptionCompleteMatchingPredicate(ctx, id, SpatialAnchorsAccountOperationPredicate{})
}

// SpatialAnchorsAccountsListBySubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ResourceClient) SpatialAnchorsAccountsListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate SpatialAnchorsAccountOperationPredicate) (result SpatialAnchorsAccountsListBySubscriptionCompleteResult, err error) {
	items := make([]SpatialAnchorsAccount, 0)

	resp, err := c.SpatialAnchorsAccountsListBySubscription(ctx, id)
	if err != nil {
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

	result = SpatialAnchorsAccountsListBySubscriptionCompleteResult{
		Items: items,
	}
	return
}
