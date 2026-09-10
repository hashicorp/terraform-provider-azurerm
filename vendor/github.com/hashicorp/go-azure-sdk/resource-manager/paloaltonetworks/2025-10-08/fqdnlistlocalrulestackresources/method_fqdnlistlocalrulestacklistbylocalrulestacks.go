package fqdnlistlocalrulestackresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FqdnListLocalRulestackListByLocalRulestacksOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]FqdnListLocalRulestackResource
}

type FqdnListLocalRulestackListByLocalRulestacksCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []FqdnListLocalRulestackResource
}

type FqdnListLocalRulestackListByLocalRulestacksCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *FqdnListLocalRulestackListByLocalRulestacksCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// FqdnListLocalRulestackListByLocalRulestacks ...
func (c FqdnListLocalRulestackResourcesClient) FqdnListLocalRulestackListByLocalRulestacks(ctx context.Context, id LocalRulestackId) (result FqdnListLocalRulestackListByLocalRulestacksOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &FqdnListLocalRulestackListByLocalRulestacksCustomPager{},
		Path:       fmt.Sprintf("%s/fqdnLists", id.ID()),
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
		Values *[]FqdnListLocalRulestackResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// FqdnListLocalRulestackListByLocalRulestacksComplete retrieves all the results into a single object
func (c FqdnListLocalRulestackResourcesClient) FqdnListLocalRulestackListByLocalRulestacksComplete(ctx context.Context, id LocalRulestackId) (FqdnListLocalRulestackListByLocalRulestacksCompleteResult, error) {
	return c.FqdnListLocalRulestackListByLocalRulestacksCompleteMatchingPredicate(ctx, id, FqdnListLocalRulestackResourceOperationPredicate{})
}

// FqdnListLocalRulestackListByLocalRulestacksCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c FqdnListLocalRulestackResourcesClient) FqdnListLocalRulestackListByLocalRulestacksCompleteMatchingPredicate(ctx context.Context, id LocalRulestackId, predicate FqdnListLocalRulestackResourceOperationPredicate) (result FqdnListLocalRulestackListByLocalRulestacksCompleteResult, err error) {
	items := make([]FqdnListLocalRulestackResource, 0)

	resp, err := c.FqdnListLocalRulestackListByLocalRulestacks(ctx, id)
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

	result = FqdnListLocalRulestackListByLocalRulestacksCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
