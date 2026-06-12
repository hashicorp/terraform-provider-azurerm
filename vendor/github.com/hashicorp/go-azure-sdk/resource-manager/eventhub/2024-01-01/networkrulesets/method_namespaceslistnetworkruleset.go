package networkrulesets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NamespacesListNetworkRuleSetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NetworkRuleSet
}

type NamespacesListNetworkRuleSetCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NetworkRuleSet
}

type NamespacesListNetworkRuleSetCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *NamespacesListNetworkRuleSetCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// NamespacesListNetworkRuleSet ...
func (c NetworkRuleSetsClient) NamespacesListNetworkRuleSet(ctx context.Context, id NamespaceId) (result NamespacesListNetworkRuleSetOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &NamespacesListNetworkRuleSetCustomPager{},
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

// NamespacesListNetworkRuleSetComplete retrieves all the results into a single object
func (c NetworkRuleSetsClient) NamespacesListNetworkRuleSetComplete(ctx context.Context, id NamespaceId) (NamespacesListNetworkRuleSetCompleteResult, error) {
	return c.NamespacesListNetworkRuleSetCompleteMatchingPredicate(ctx, id, NetworkRuleSetOperationPredicate{})
}

// NamespacesListNetworkRuleSetCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NetworkRuleSetsClient) NamespacesListNetworkRuleSetCompleteMatchingPredicate(ctx context.Context, id NamespaceId, predicate NetworkRuleSetOperationPredicate) (result NamespacesListNetworkRuleSetCompleteResult, err error) {
	items := make([]NetworkRuleSet, 0)

	resp, err := c.NamespacesListNetworkRuleSet(ctx, id)
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

	result = NamespacesListNetworkRuleSetCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
