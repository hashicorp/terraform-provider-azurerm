package replicationnetworkmappings

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByReplicationNetworksOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NetworkMapping
}

type ListByReplicationNetworksCompleteResult struct {
	Items []NetworkMapping
}

// ListByReplicationNetworks ...
func (c ReplicationNetworkMappingsClient) ListByReplicationNetworks(ctx context.Context, id ReplicationNetworkId) (result ListByReplicationNetworksOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/replicationNetworkMappings", id.ID()),
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
		Values *[]NetworkMapping `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByReplicationNetworksComplete retrieves all the results into a single object
func (c ReplicationNetworkMappingsClient) ListByReplicationNetworksComplete(ctx context.Context, id ReplicationNetworkId) (ListByReplicationNetworksCompleteResult, error) {
	return c.ListByReplicationNetworksCompleteMatchingPredicate(ctx, id, NetworkMappingOperationPredicate{})
}

// ListByReplicationNetworksCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ReplicationNetworkMappingsClient) ListByReplicationNetworksCompleteMatchingPredicate(ctx context.Context, id ReplicationNetworkId, predicate NetworkMappingOperationPredicate) (result ListByReplicationNetworksCompleteResult, err error) {
	items := make([]NetworkMapping, 0)

	resp, err := c.ListByReplicationNetworks(ctx, id)
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

	result = ListByReplicationNetworksCompleteResult{
		Items: items,
	}
	return
}
