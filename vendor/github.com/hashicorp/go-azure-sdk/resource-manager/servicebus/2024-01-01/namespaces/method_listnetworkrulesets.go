package namespaces

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListNetworkRuleSetsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NetworkRuleSet
}

type ListNetworkRuleSetsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NetworkRuleSet
}

type ListNetworkRuleSetsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListNetworkRuleSetsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListNetworkRuleSets ...
func (c NamespacesClient) ListNetworkRuleSets(ctx context.Context, id NamespaceId) (result ListNetworkRuleSetsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListNetworkRuleSetsCustomPager{},
		Path:       fmt.Sprintf("%s/networkRuleSets", id.ID()),
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
		Values *[]NetworkRuleSet `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListNetworkRuleSetsComplete retrieves all the results into a single object
func (c NamespacesClient) ListNetworkRuleSetsComplete(ctx context.Context, id NamespaceId) (ListNetworkRuleSetsCompleteResult, error) {
	return c.ListNetworkRuleSetsCompleteMatchingPredicate(ctx, id, NetworkRuleSetOperationPredicate{})
}

// ListNetworkRuleSetsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NamespacesClient) ListNetworkRuleSetsCompleteMatchingPredicate(ctx context.Context, id NamespaceId, predicate NetworkRuleSetOperationPredicate) (result ListNetworkRuleSetsCompleteResult, err error) {
	items := make([]NetworkRuleSet, 0)

	resp, err := c.ListNetworkRuleSets(ctx, id)
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

	result = ListNetworkRuleSetsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
