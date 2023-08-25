package replicationprotecteditems

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByReplicationProtectionContainersOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ReplicationProtectedItem
}

type ListByReplicationProtectionContainersCompleteResult struct {
	Items []ReplicationProtectedItem
}

// ListByReplicationProtectionContainers ...
func (c ReplicationProtectedItemsClient) ListByReplicationProtectionContainers(ctx context.Context, id ReplicationProtectionContainerId) (result ListByReplicationProtectionContainersOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/replicationProtectedItems", id.ID()),
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
		Values *[]ReplicationProtectedItem `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByReplicationProtectionContainersComplete retrieves all the results into a single object
func (c ReplicationProtectedItemsClient) ListByReplicationProtectionContainersComplete(ctx context.Context, id ReplicationProtectionContainerId) (ListByReplicationProtectionContainersCompleteResult, error) {
	return c.ListByReplicationProtectionContainersCompleteMatchingPredicate(ctx, id, ReplicationProtectedItemOperationPredicate{})
}

// ListByReplicationProtectionContainersCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ReplicationProtectedItemsClient) ListByReplicationProtectionContainersCompleteMatchingPredicate(ctx context.Context, id ReplicationProtectionContainerId, predicate ReplicationProtectedItemOperationPredicate) (result ListByReplicationProtectionContainersCompleteResult, err error) {
	items := make([]ReplicationProtectedItem, 0)

	resp, err := c.ListByReplicationProtectionContainers(ctx, id)
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

	result = ListByReplicationProtectionContainersCompleteResult{
		Items: items,
	}
	return
}
