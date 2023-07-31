package tagrules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubAccountTagRulesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]MonitoringTagRules
}

type SubAccountTagRulesListCompleteResult struct {
	Items []MonitoringTagRules
}

// SubAccountTagRulesList ...
func (c TagRulesClient) SubAccountTagRulesList(ctx context.Context, id AccountId) (result SubAccountTagRulesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
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

// SubAccountTagRulesListComplete retrieves all the results into a single object
func (c TagRulesClient) SubAccountTagRulesListComplete(ctx context.Context, id AccountId) (SubAccountTagRulesListCompleteResult, error) {
	return c.SubAccountTagRulesListCompleteMatchingPredicate(ctx, id, MonitoringTagRulesOperationPredicate{})
}

// SubAccountTagRulesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c TagRulesClient) SubAccountTagRulesListCompleteMatchingPredicate(ctx context.Context, id AccountId, predicate MonitoringTagRulesOperationPredicate) (result SubAccountTagRulesListCompleteResult, err error) {
	items := make([]MonitoringTagRules, 0)

	resp, err := c.SubAccountTagRulesList(ctx, id)
	if err != nil {
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

	result = SubAccountTagRulesListCompleteResult{
		Items: items,
	}
	return
}
