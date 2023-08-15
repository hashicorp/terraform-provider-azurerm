package diskaccesses

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListPrivateEndpointConnectionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PrivateEndpointConnection
}

type ListPrivateEndpointConnectionsCompleteResult struct {
	Items []PrivateEndpointConnection
}

// ListPrivateEndpointConnections ...
func (c DiskAccessesClient) ListPrivateEndpointConnections(ctx context.Context, id DiskAccessId) (result ListPrivateEndpointConnectionsOperationResponse, err error) {
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

// ListPrivateEndpointConnectionsComplete retrieves all the results into a single object
func (c DiskAccessesClient) ListPrivateEndpointConnectionsComplete(ctx context.Context, id DiskAccessId) (ListPrivateEndpointConnectionsCompleteResult, error) {
	return c.ListPrivateEndpointConnectionsCompleteMatchingPredicate(ctx, id, PrivateEndpointConnectionOperationPredicate{})
}

// ListPrivateEndpointConnectionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DiskAccessesClient) ListPrivateEndpointConnectionsCompleteMatchingPredicate(ctx context.Context, id DiskAccessId, predicate PrivateEndpointConnectionOperationPredicate) (result ListPrivateEndpointConnectionsCompleteResult, err error) {
	items := make([]PrivateEndpointConnection, 0)

	resp, err := c.ListPrivateEndpointConnections(ctx, id)
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

	result = ListPrivateEndpointConnectionsCompleteResult{
		Items: items,
	}
	return
}
