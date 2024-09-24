package replicationprotectioncontainers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByReplicationFabricsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ProtectionContainer
}

type ListByReplicationFabricsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ProtectionContainer
}

type ListByReplicationFabricsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByReplicationFabricsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByReplicationFabrics ...
func (c ReplicationProtectionContainersClient) ListByReplicationFabrics(ctx context.Context, id ReplicationFabricId) (result ListByReplicationFabricsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByReplicationFabricsCustomPager{},
		Path:       fmt.Sprintf("%s/replicationProtectionContainers", id.ID()),
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
		Values *[]ProtectionContainer `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByReplicationFabricsComplete retrieves all the results into a single object
func (c ReplicationProtectionContainersClient) ListByReplicationFabricsComplete(ctx context.Context, id ReplicationFabricId) (ListByReplicationFabricsCompleteResult, error) {
	return c.ListByReplicationFabricsCompleteMatchingPredicate(ctx, id, ProtectionContainerOperationPredicate{})
}

// ListByReplicationFabricsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ReplicationProtectionContainersClient) ListByReplicationFabricsCompleteMatchingPredicate(ctx context.Context, id ReplicationFabricId, predicate ProtectionContainerOperationPredicate) (result ListByReplicationFabricsCompleteResult, err error) {
	items := make([]ProtectionContainer, 0)

	resp, err := c.ListByReplicationFabrics(ctx, id)
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

	result = ListByReplicationFabricsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
