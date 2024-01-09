package localrules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByLocalRulestacksOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LocalRulesResource
}

type ListByLocalRulestacksCompleteResult struct {
	Items []LocalRulesResource
}

// ListByLocalRulestacks ...
func (c LocalRulesClient) ListByLocalRulestacks(ctx context.Context, id LocalRulestackId) (result ListByLocalRulestacksOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
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

// ListByLocalRulestacksComplete retrieves all the results into a single object
func (c LocalRulesClient) ListByLocalRulestacksComplete(ctx context.Context, id LocalRulestackId) (ListByLocalRulestacksCompleteResult, error) {
	return c.ListByLocalRulestacksCompleteMatchingPredicate(ctx, id, LocalRulesResourceOperationPredicate{})
}

// ListByLocalRulestacksCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LocalRulesClient) ListByLocalRulestacksCompleteMatchingPredicate(ctx context.Context, id LocalRulestackId, predicate LocalRulesResourceOperationPredicate) (result ListByLocalRulestacksCompleteResult, err error) {
	items := make([]LocalRulesResource, 0)

	resp, err := c.ListByLocalRulestacks(ctx, id)
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

	result = ListByLocalRulestacksCompleteResult{
		Items: items,
	}
	return
}
