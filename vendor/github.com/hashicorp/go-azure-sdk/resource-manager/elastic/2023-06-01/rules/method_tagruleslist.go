package rules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagRulesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]MonitoringTagRules
}

type TagRulesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []MonitoringTagRules
}

type TagRulesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *TagRulesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// TagRulesList ...
func (c RulesClient) TagRulesList(ctx context.Context, id MonitorId) (result TagRulesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &TagRulesListCustomPager{},
		Path:       fmt.Sprintf("%s/tagRules", id.ID()),
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
		Values *[]MonitoringTagRules `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// TagRulesListComplete retrieves all the results into a single object
func (c RulesClient) TagRulesListComplete(ctx context.Context, id MonitorId) (TagRulesListCompleteResult, error) {
	return c.TagRulesListCompleteMatchingPredicate(ctx, id, MonitoringTagRulesOperationPredicate{})
}

// TagRulesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RulesClient) TagRulesListCompleteMatchingPredicate(ctx context.Context, id MonitorId, predicate MonitoringTagRulesOperationPredicate) (result TagRulesListCompleteResult, err error) {
	items := make([]MonitoringTagRules, 0)

	resp, err := c.TagRulesList(ctx, id)
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

	result = TagRulesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
