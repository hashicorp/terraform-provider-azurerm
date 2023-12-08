package replicationprotectioncontainermappings

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
	Model        *[]ProtectionContainerMapping
}

type ListByReplicationProtectionContainersCompleteResult struct {
	Items []ProtectionContainerMapping
}

// ListByReplicationProtectionContainers ...
func (c ReplicationProtectionContainerMappingsClient) ListByReplicationProtectionContainers(ctx context.Context, id ReplicationProtectionContainerId) (result ListByReplicationProtectionContainersOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/replicationProtectionContainerMappings", id.ID()),
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
		Values *[]ProtectionContainerMapping `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByReplicationProtectionContainersComplete retrieves all the results into a single object
func (c ReplicationProtectionContainerMappingsClient) ListByReplicationProtectionContainersComplete(ctx context.Context, id ReplicationProtectionContainerId) (ListByReplicationProtectionContainersCompleteResult, error) {
	return c.ListByReplicationProtectionContainersCompleteMatchingPredicate(ctx, id, ProtectionContainerMappingOperationPredicate{})
}

// ListByReplicationProtectionContainersCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ReplicationProtectionContainerMappingsClient) ListByReplicationProtectionContainersCompleteMatchingPredicate(ctx context.Context, id ReplicationProtectionContainerId, predicate ProtectionContainerMappingOperationPredicate) (result ListByReplicationProtectionContainersCompleteResult, err error) {
	items := make([]ProtectionContainerMapping, 0)

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
