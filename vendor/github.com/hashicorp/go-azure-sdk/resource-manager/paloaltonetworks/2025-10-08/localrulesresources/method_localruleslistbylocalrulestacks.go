package localrulesresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocalRulesListByLocalRulestacksOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LocalRulesResource
}

type LocalRulesListByLocalRulestacksCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []LocalRulesResource
}

type LocalRulesListByLocalRulestacksCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *LocalRulesListByLocalRulestacksCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// LocalRulesListByLocalRulestacks ...
func (c LocalRulesResourcesClient) LocalRulesListByLocalRulestacks(ctx context.Context, id LocalRulestackId) (result LocalRulesListByLocalRulestacksOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &LocalRulesListByLocalRulestacksCustomPager{},
		Path:       fmt.Sprintf("%s/localRules", id.ID()),
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
		Values *[]LocalRulesResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LocalRulesListByLocalRulestacksComplete retrieves all the results into a single object
func (c LocalRulesResourcesClient) LocalRulesListByLocalRulestacksComplete(ctx context.Context, id LocalRulestackId) (LocalRulesListByLocalRulestacksCompleteResult, error) {
	return c.LocalRulesListByLocalRulestacksCompleteMatchingPredicate(ctx, id, LocalRulesResourceOperationPredicate{})
}

// LocalRulesListByLocalRulestacksCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LocalRulesResourcesClient) LocalRulesListByLocalRulestacksCompleteMatchingPredicate(ctx context.Context, id LocalRulestackId, predicate LocalRulesResourceOperationPredicate) (result LocalRulesListByLocalRulestacksCompleteResult, err error) {
	items := make([]LocalRulesResource, 0)

	resp, err := c.LocalRulesListByLocalRulestacks(ctx, id)
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

	result = LocalRulesListByLocalRulestacksCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
