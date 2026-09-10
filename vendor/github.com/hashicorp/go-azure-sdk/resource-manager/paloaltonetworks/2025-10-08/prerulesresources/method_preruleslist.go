package prerulesresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PreRulesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PreRulesResource
}

type PreRulesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PreRulesResource
}

type PreRulesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *PreRulesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// PreRulesList ...
func (c PreRulesResourcesClient) PreRulesList(ctx context.Context, id GlobalRulestackId) (result PreRulesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &PreRulesListCustomPager{},
		Path:       fmt.Sprintf("%s/preRules", id.ID()),
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
		Values *[]PreRulesResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// PreRulesListComplete retrieves all the results into a single object
func (c PreRulesResourcesClient) PreRulesListComplete(ctx context.Context, id GlobalRulestackId) (PreRulesListCompleteResult, error) {
	return c.PreRulesListCompleteMatchingPredicate(ctx, id, PreRulesResourceOperationPredicate{})
}

// PreRulesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PreRulesResourcesClient) PreRulesListCompleteMatchingPredicate(ctx context.Context, id GlobalRulestackId, predicate PreRulesResourceOperationPredicate) (result PreRulesListCompleteResult, err error) {
	items := make([]PreRulesResource, 0)

	resp, err := c.PreRulesList(ctx, id)
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

	result = PreRulesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
