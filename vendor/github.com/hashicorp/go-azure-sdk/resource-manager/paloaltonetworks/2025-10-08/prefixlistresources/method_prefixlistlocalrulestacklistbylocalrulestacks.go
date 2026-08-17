package prefixlistresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrefixListLocalRulestackListByLocalRulestacksOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PrefixListResource
}

type PrefixListLocalRulestackListByLocalRulestacksCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PrefixListResource
}

type PrefixListLocalRulestackListByLocalRulestacksCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *PrefixListLocalRulestackListByLocalRulestacksCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// PrefixListLocalRulestackListByLocalRulestacks ...
func (c PrefixListResourcesClient) PrefixListLocalRulestackListByLocalRulestacks(ctx context.Context, id LocalRulestackId) (result PrefixListLocalRulestackListByLocalRulestacksOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &PrefixListLocalRulestackListByLocalRulestacksCustomPager{},
		Path:       fmt.Sprintf("%s/prefixLists", id.ID()),
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
		Values *[]PrefixListResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// PrefixListLocalRulestackListByLocalRulestacksComplete retrieves all the results into a single object
func (c PrefixListResourcesClient) PrefixListLocalRulestackListByLocalRulestacksComplete(ctx context.Context, id LocalRulestackId) (PrefixListLocalRulestackListByLocalRulestacksCompleteResult, error) {
	return c.PrefixListLocalRulestackListByLocalRulestacksCompleteMatchingPredicate(ctx, id, PrefixListResourceOperationPredicate{})
}

// PrefixListLocalRulestackListByLocalRulestacksCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PrefixListResourcesClient) PrefixListLocalRulestackListByLocalRulestacksCompleteMatchingPredicate(ctx context.Context, id LocalRulestackId, predicate PrefixListResourceOperationPredicate) (result PrefixListLocalRulestackListByLocalRulestacksCompleteResult, err error) {
	items := make([]PrefixListResource, 0)

	resp, err := c.PrefixListLocalRulestackListByLocalRulestacks(ctx, id)
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

	result = PrefixListLocalRulestackListByLocalRulestacksCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
