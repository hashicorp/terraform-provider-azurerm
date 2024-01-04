package resourceguards

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

type GetResourcesInSubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ResourceGuardResource
}

type GetResourcesInSubscriptionCompleteResult struct {
	Items []ResourceGuardResource
}

// GetResourcesInSubscription ...
func (c ResourceGuardsClient) GetResourcesInSubscription(ctx context.Context, id commonids.SubscriptionId) (result GetResourcesInSubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.DataProtection/resourceGuards", id.ID()),
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
		Values *[]ResourceGuardResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetResourcesInSubscriptionComplete retrieves all the results into a single object
func (c ResourceGuardsClient) GetResourcesInSubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (GetResourcesInSubscriptionCompleteResult, error) {
	return c.GetResourcesInSubscriptionCompleteMatchingPredicate(ctx, id, ResourceGuardResourceOperationPredicate{})
}

// GetResourcesInSubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ResourceGuardsClient) GetResourcesInSubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate ResourceGuardResourceOperationPredicate) (result GetResourcesInSubscriptionCompleteResult, err error) {
	items := make([]ResourceGuardResource, 0)

	resp, err := c.GetResourcesInSubscription(ctx, id)
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

	result = GetResourcesInSubscriptionCompleteResult{
		Items: items,
	}
	return
}
