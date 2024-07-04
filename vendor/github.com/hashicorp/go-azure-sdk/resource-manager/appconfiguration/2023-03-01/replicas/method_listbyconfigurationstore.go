package replicas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByConfigurationStoreOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Replica
}

type ListByConfigurationStoreCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Replica
}

type ListByConfigurationStoreCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByConfigurationStoreCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByConfigurationStore ...
func (c ReplicasClient) ListByConfigurationStore(ctx context.Context, id ConfigurationStoreId) (result ListByConfigurationStoreOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByConfigurationStoreCustomPager{},
		Path:       fmt.Sprintf("%s/replicas", id.ID()),
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
		Values *[]Replica `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByConfigurationStoreComplete retrieves all the results into a single object
func (c ReplicasClient) ListByConfigurationStoreComplete(ctx context.Context, id ConfigurationStoreId) (ListByConfigurationStoreCompleteResult, error) {
	return c.ListByConfigurationStoreCompleteMatchingPredicate(ctx, id, ReplicaOperationPredicate{})
}

// ListByConfigurationStoreCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ReplicasClient) ListByConfigurationStoreCompleteMatchingPredicate(ctx context.Context, id ConfigurationStoreId, predicate ReplicaOperationPredicate) (result ListByConfigurationStoreCompleteResult, err error) {
	items := make([]Replica, 0)

	resp, err := c.ListByConfigurationStore(ctx, id)
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

	result = ListByConfigurationStoreCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
