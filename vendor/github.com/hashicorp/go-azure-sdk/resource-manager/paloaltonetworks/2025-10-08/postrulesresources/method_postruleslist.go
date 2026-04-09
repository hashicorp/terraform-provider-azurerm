package postrulesresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PostRulesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PostRulesResource
}

type PostRulesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PostRulesResource
}

type PostRulesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *PostRulesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// PostRulesList ...
func (c PostRulesResourcesClient) PostRulesList(ctx context.Context, id GlobalRulestackId) (result PostRulesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &PostRulesListCustomPager{},
		Path:       fmt.Sprintf("%s/postRules", id.ID()),
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
		Values *[]PostRulesResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// PostRulesListComplete retrieves all the results into a single object
func (c PostRulesResourcesClient) PostRulesListComplete(ctx context.Context, id GlobalRulestackId) (PostRulesListCompleteResult, error) {
	return c.PostRulesListCompleteMatchingPredicate(ctx, id, PostRulesResourceOperationPredicate{})
}

// PostRulesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PostRulesResourcesClient) PostRulesListCompleteMatchingPredicate(ctx context.Context, id GlobalRulestackId, predicate PostRulesResourceOperationPredicate) (result PostRulesListCompleteResult, err error) {
	items := make([]PostRulesResource, 0)

	resp, err := c.PostRulesList(ctx, id)
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

	result = PostRulesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
