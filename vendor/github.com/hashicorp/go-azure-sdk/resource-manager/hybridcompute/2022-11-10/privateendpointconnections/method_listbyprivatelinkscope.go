package privateendpointconnections

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByPrivateLinkScopeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PrivateEndpointConnection
}

type ListByPrivateLinkScopeCompleteResult struct {
	Items []PrivateEndpointConnection
}

// ListByPrivateLinkScope ...
func (c PrivateEndpointConnectionsClient) ListByPrivateLinkScope(ctx context.Context, id ProviderPrivateLinkScopeId) (result ListByPrivateLinkScopeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/privateEndpointConnections", id.ID()),
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
		Values *[]PrivateEndpointConnection `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByPrivateLinkScopeComplete retrieves all the results into a single object
func (c PrivateEndpointConnectionsClient) ListByPrivateLinkScopeComplete(ctx context.Context, id ProviderPrivateLinkScopeId) (ListByPrivateLinkScopeCompleteResult, error) {
	return c.ListByPrivateLinkScopeCompleteMatchingPredicate(ctx, id, PrivateEndpointConnectionOperationPredicate{})
}

// ListByPrivateLinkScopeCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PrivateEndpointConnectionsClient) ListByPrivateLinkScopeCompleteMatchingPredicate(ctx context.Context, id ProviderPrivateLinkScopeId, predicate PrivateEndpointConnectionOperationPredicate) (result ListByPrivateLinkScopeCompleteResult, err error) {
	items := make([]PrivateEndpointConnection, 0)

	resp, err := c.ListByPrivateLinkScope(ctx, id)
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

	result = ListByPrivateLinkScopeCompleteResult{
		Items: items,
	}
	return
}
