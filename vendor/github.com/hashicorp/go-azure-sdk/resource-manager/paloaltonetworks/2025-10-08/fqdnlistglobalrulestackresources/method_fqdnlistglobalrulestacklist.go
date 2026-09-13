package fqdnlistglobalrulestackresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FqdnListGlobalRulestackListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]FqdnListGlobalRulestackResource
}

type FqdnListGlobalRulestackListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []FqdnListGlobalRulestackResource
}

type FqdnListGlobalRulestackListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *FqdnListGlobalRulestackListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// FqdnListGlobalRulestackList ...
func (c FqdnListGlobalRulestackResourcesClient) FqdnListGlobalRulestackList(ctx context.Context, id GlobalRulestackId) (result FqdnListGlobalRulestackListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &FqdnListGlobalRulestackListCustomPager{},
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
		Values *[]FqdnListGlobalRulestackResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// FqdnListGlobalRulestackListComplete retrieves all the results into a single object
func (c FqdnListGlobalRulestackResourcesClient) FqdnListGlobalRulestackListComplete(ctx context.Context, id GlobalRulestackId) (FqdnListGlobalRulestackListCompleteResult, error) {
	return c.FqdnListGlobalRulestackListCompleteMatchingPredicate(ctx, id, FqdnListGlobalRulestackResourceOperationPredicate{})
}

// FqdnListGlobalRulestackListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c FqdnListGlobalRulestackResourcesClient) FqdnListGlobalRulestackListCompleteMatchingPredicate(ctx context.Context, id GlobalRulestackId, predicate FqdnListGlobalRulestackResourceOperationPredicate) (result FqdnListGlobalRulestackListCompleteResult, err error) {
	items := make([]FqdnListGlobalRulestackResource, 0)

	resp, err := c.FqdnListGlobalRulestackList(ctx, id)
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

	result = FqdnListGlobalRulestackListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
