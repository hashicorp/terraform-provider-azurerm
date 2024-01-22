package managedidentities

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

type UserAssignedIdentitiesListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Identity
}

type UserAssignedIdentitiesListBySubscriptionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Identity
}

// UserAssignedIdentitiesListBySubscription ...
func (c ManagedIdentitiesClient) UserAssignedIdentitiesListBySubscription(ctx context.Context, id commonids.SubscriptionId) (result UserAssignedIdentitiesListBySubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities", id.ID()),
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
		Values *[]Identity `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// UserAssignedIdentitiesListBySubscriptionComplete retrieves all the results into a single object
func (c ManagedIdentitiesClient) UserAssignedIdentitiesListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (UserAssignedIdentitiesListBySubscriptionCompleteResult, error) {
	return c.UserAssignedIdentitiesListBySubscriptionCompleteMatchingPredicate(ctx, id, IdentityOperationPredicate{})
}

// UserAssignedIdentitiesListBySubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagedIdentitiesClient) UserAssignedIdentitiesListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate IdentityOperationPredicate) (result UserAssignedIdentitiesListBySubscriptionCompleteResult, err error) {
	items := make([]Identity, 0)

	resp, err := c.UserAssignedIdentitiesListBySubscription(ctx, id)
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

	result = UserAssignedIdentitiesListBySubscriptionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
