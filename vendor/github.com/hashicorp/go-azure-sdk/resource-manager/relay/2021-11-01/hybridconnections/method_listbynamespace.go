package hybridconnections

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByNamespaceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]HybridConnection
}

type ListByNamespaceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []HybridConnection
}

type ListByNamespaceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByNamespaceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByNamespace ...
func (c HybridConnectionsClient) ListByNamespace(ctx context.Context, id NamespaceId) (result ListByNamespaceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByNamespaceCustomPager{},
		Path:       fmt.Sprintf("%s/hybridConnections", id.ID()),
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
		Values *[]HybridConnection `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByNamespaceComplete retrieves all the results into a single object
func (c HybridConnectionsClient) ListByNamespaceComplete(ctx context.Context, id NamespaceId) (ListByNamespaceCompleteResult, error) {
	return c.ListByNamespaceCompleteMatchingPredicate(ctx, id, HybridConnectionOperationPredicate{})
}

// ListByNamespaceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c HybridConnectionsClient) ListByNamespaceCompleteMatchingPredicate(ctx context.Context, id NamespaceId, predicate HybridConnectionOperationPredicate) (result ListByNamespaceCompleteResult, err error) {
	items := make([]HybridConnection, 0)

	resp, err := c.ListByNamespace(ctx, id)
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

	result = ListByNamespaceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
