package securityrules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DefaultSecurityRulesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SecurityRule
}

type DefaultSecurityRulesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SecurityRule
}

type DefaultSecurityRulesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DefaultSecurityRulesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DefaultSecurityRulesList ...
func (c SecurityRulesClient) DefaultSecurityRulesList(ctx context.Context, id NetworkSecurityGroupId) (result DefaultSecurityRulesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &DefaultSecurityRulesListCustomPager{},
		Path:       fmt.Sprintf("%s/defaultSecurityRules", id.ID()),
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
		Values *[]SecurityRule `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DefaultSecurityRulesListComplete retrieves all the results into a single object
func (c SecurityRulesClient) DefaultSecurityRulesListComplete(ctx context.Context, id NetworkSecurityGroupId) (DefaultSecurityRulesListCompleteResult, error) {
	return c.DefaultSecurityRulesListCompleteMatchingPredicate(ctx, id, SecurityRuleOperationPredicate{})
}

// DefaultSecurityRulesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SecurityRulesClient) DefaultSecurityRulesListCompleteMatchingPredicate(ctx context.Context, id NetworkSecurityGroupId, predicate SecurityRuleOperationPredicate) (result DefaultSecurityRulesListCompleteResult, err error) {
	items := make([]SecurityRule, 0)

	resp, err := c.DefaultSecurityRulesList(ctx, id)
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

	result = DefaultSecurityRulesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
